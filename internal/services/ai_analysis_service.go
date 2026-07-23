package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/repositories"
)

type AIAnalysisService struct {
	deploymentRepo repositories.DeploymentRepository
	appRepo        repositories.AppServiceRepository
	aiSettingsRepo repositories.AISettingsRepository
	httpClient     *http.Client
}

func NewAIAnalysisService(d repositories.DeploymentRepository, a repositories.AppServiceRepository, ai repositories.AISettingsRepository) *AIAnalysisService {
	return &AIAnalysisService{
		deploymentRepo: d,
		appRepo:        a,
		aiSettingsRepo: ai,
		httpClient:     &http.Client{Timeout: 60 * time.Second},
	}
}

type AIExplanationResponse struct {
	Summary         string   `json:"summary"`
	Cause           string   `json:"cause"`
	SuggestedFix    string   `json:"suggestedFix"`
	Confidence      string   `json:"confidence"`
	Commands        []string `json:"commands"`
	RelatedLogLines []string `json:"relatedLogLines"`
}

func (s *AIAnalysisService) ExplainDeploymentFailure(ctx context.Context, deploymentID string) (*AIExplanationResponse, error) {
	deployment, err := s.deploymentRepo.GetByID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}
	if deployment.Status != models.DeploymentStatusFailed {
		return nil, errors.New("AI explanations are only available after a deployment fails")
	}

	app, err := s.appRepo.GetByID(ctx, deployment.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	aiSettings, err := s.aiSettingsRepo.GetAISettings(ctx)
	if err != nil || aiSettings == nil {
		return nil, errors.New("failed to load AI settings")
	}

	var baseURL, apiKey, modelName string
	provider := aiSettings.DefaultProvider
	if provider == "" {
		provider = "openai"
	}

	switch provider {
	case "anthropic":
		return nil, errors.New("Anthropic is not fully supported yet in this AI pipeline without a custom proxy")
	case "google":
		return nil, errors.New("Google is not fully supported yet in this AI pipeline without a custom proxy")
	case "groq":
		baseURL = "https://api.groq.com/openai/v1"
		apiKey = aiSettings.GroqKey
		modelName = aiSettings.GroqModel
		if modelName == "" {
			modelName = "llama3-8b-8192"
		}
	case "mistral":
		baseURL = "https://api.mistral.ai/v1"
		apiKey = aiSettings.MistralKey
		modelName = aiSettings.MistralModel
		if modelName == "" {
			modelName = "mistral-large-latest"
		}
	case "deepseek":
		baseURL = "https://api.deepseek.com/v1"
		apiKey = aiSettings.DeepSeekKey
		modelName = aiSettings.DeepSeekModel
		if modelName == "" {
			modelName = "deepseek-chat"
		}
	case "xai":
		baseURL = "https://api.x.ai/v1"
		apiKey = aiSettings.XAIKey
		modelName = aiSettings.XAIModel
		if modelName == "" {
			modelName = "grok-beta"
		}
	case "moonshot":
		baseURL = "https://api.moonshot.cn/v1"
		apiKey = aiSettings.MoonshotKey
		modelName = aiSettings.MoonshotModel
		if modelName == "" {
			modelName = "moonshot-v1-8k"
		}
	case "openai":
		fallthrough
	default:
		baseURL = "https://api.openai.com/v1"
		apiKey = aiSettings.OpenAIKey
		modelName = aiSettings.OpenAIModel
		if modelName == "" {
			modelName = "gpt-4o"
		}
	}

	if apiKey == "" {
		return nil, fmt.Errorf("%s API key is missing in AI settings", provider)
	}

	prompt := buildPrompt(deployment, app)

	return callOpenAI(ctx, s.httpClient, baseURL, apiKey, modelName, prompt)
}

func buildPrompt(d *models.Deployment, s *models.AppService) string {
	maxLogChars := 30000
	logs := d.BuildLogs
	if len(logs) > maxLogChars {
		logs = logs[len(logs)-maxLogChars:]
	}

	return fmt.Sprintf(`You are helping diagnose a failed Codedock deployment.

Explain the failure in plain, specific language for the app owner. Focus on what happened and the smallest likely fix.

Rules:
- Use only the deployment metadata and logs below.
- Do not invent files, commands, ports, packages, or environment variables that are not supported by the logs.
- Prefer concrete fixes over generic advice.
- If evidence is weak, say so and set confidence to low.
- Keep summary, cause, and suggestedFix each under 80 words.
- relatedLogLines should quote the shortest relevant log lines exactly as they appear below.
- commands should contain only commands or setting changes the user can plausibly run or make.

Deployment metadata:
Service: %s
Runtime mode: %s
Repository: %s
Branch: %s
Root directory: %s
Install command: %s
Build command: %s
Start command: %s
Deployment status: %s

Deployment logs:
%s`, s.Name, s.RuntimeMode, s.RepositoryURL, s.Branch, s.RootDirectory, s.InstallCommand, s.BuildCommand, s.StartCommand, d.Status, logs)
}

func callOpenAI(ctx context.Context, client *http.Client, baseURL, apiKey, model, prompt string) (*AIExplanationResponse, error) {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/chat/completions"
	} else {
		baseURL = baseURL + "/chat/completions"
	}

	reqBody := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.2,
		"response_format": map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name": "deployment_failure_explanation",
				"schema": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"summary":         map[string]any{"type": "string"},
						"cause":           map[string]any{"type": "string"},
						"suggestedFix":    map[string]any{"type": "string"},
						"confidence":      map[string]any{"type": "string", "enum": []string{"low", "medium", "high"}},
						"commands":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
						"relatedLogLines": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					},
					"required":             []string{"summary", "cause", "suggestedFix", "confidence", "commands", "relatedLogLines"},
					"additionalProperties": false,
				},
				"strict": true,
			},
		},
	}

	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, errors.New("OpenAI returned no choices")
	}

	var explanation AIExplanationResponse
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &explanation); err != nil {
		return nil, fmt.Errorf("failed to parse explanation JSON: %w", err)
	}

	return &explanation, nil
}
