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

type AISettingsService struct {
	repo       repositories.AISettingsRepository
	httpClient *http.Client
}

func NewAISettingsService(repo repositories.AISettingsRepository) *AISettingsService {
	return &AISettingsService{
		repo:       repo,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *AISettingsService) GetAISettings(ctx context.Context) (*models.AISettings, error) {
	return s.repo.GetAISettings(ctx)
}

func (s *AISettingsService) UpdateAISettings(ctx context.Context, cfg *models.AISettings) error {
	return s.repo.UpdateAISettings(ctx, cfg)
}

func (s *AISettingsService) DiagnoseLogs(ctx context.Context, logs string) (string, error) {
	settings, err := s.repo.GetAISettings(ctx)
	if err != nil || settings == nil {
		return "", errors.New("failed to load AI settings")
	}

	var baseURL, apiKey, modelName string
	provider := settings.DefaultProvider
	if provider == "" {
		provider = "openai"
	}

	switch provider {
	case "anthropic":
		return "", errors.New("Anthropic is not fully supported for this endpoint yet")
	case "google":
		return "", errors.New("Google is not fully supported for this endpoint yet")
	case "groq":
		baseURL = "https://api.groq.com/openai/v1"
		apiKey = settings.GroqKey
		modelName = settings.GroqModel
		if modelName == "" {
			modelName = "llama3-8b-8192"
		}
	case "mistral":
		baseURL = "https://api.mistral.ai/v1"
		apiKey = settings.MistralKey
		modelName = settings.MistralModel
		if modelName == "" {
			modelName = "mistral-large-latest"
		}
	case "deepseek":
		baseURL = "https://api.deepseek.com/v1"
		apiKey = settings.DeepSeekKey
		modelName = settings.DeepSeekModel
		if modelName == "" {
			modelName = "deepseek-chat"
		}
	case "xai":
		baseURL = "https://api.x.ai/v1"
		apiKey = settings.XAIKey
		modelName = settings.XAIModel
		if modelName == "" {
			modelName = "grok-beta"
		}
	case "moonshot":
		baseURL = "https://api.moonshot.cn/v1"
		apiKey = settings.MoonshotKey
		modelName = settings.MoonshotModel
		if modelName == "" {
			modelName = "moonshot-v1-8k"
		}
	case "openai":
		fallthrough
	default:
		baseURL = "https://api.openai.com/v1"
		apiKey = settings.OpenAIKey
		modelName = settings.OpenAIModel
		if modelName == "" {
			modelName = "gpt-4o"
		}
	}

	if apiKey == "" {
		return "", fmt.Errorf("%s API key is missing in AI settings", provider)
	}

	prompt := fmt.Sprintf("You are a DevOps AI assistant. Diagnose these logs and explain the error in plain English. Provide a fix if possible:\n\n%s", logs)

	url := baseURL + "/chat/completions"
	reqBody := map[string]any{
		"model": modelName,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.2,
	}

	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("AI API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", errors.New("AI returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}
