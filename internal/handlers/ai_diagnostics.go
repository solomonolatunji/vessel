package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"

	"vessl.dev/vessl/internal/services"
)

const (
	ProviderOpenAI   = "openai"
	ProviderDeepseek = "deepseek"
	ProviderGroq     = "groq"

	BaseURLDeepseek = "https://api.deepseek.com"
	BaseURLGroq     = "https://api.groq.com/openai/v1"

	DefaultModelOpenAI   = openai.GPT4oMini
	DefaultModelDeepseek = "deepseek-chat"
	DefaultModelGroq     = "llama3-8b-8192"
)

type AIDiagnosticsHandler struct {
	aiService         *services.AISettingsService
	deploymentService *services.DeploymentService
	projectService    *services.ProjectService
}

func NewAIDiagnosticsHandler(ai *services.AISettingsService, ds *services.DeploymentService, ps *services.ProjectService) *AIDiagnosticsHandler {
	return &AIDiagnosticsHandler{
		aiService:         ai,
		deploymentService: ds,
		projectService:    ps,
	}
}

// @Summary Analyze endpoint
// @Description Analyze endpoint
// @Tags Deployments
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/deployments/{id}/diagnostics [post]
func (h *AIDiagnosticsHandler) Analyze(c echo.Context) error {
	deploymentID := c.Param("id")
	if deploymentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "deployment ID required"})
	}

	dep, err := h.deploymentService.GetDeployment(c.Request().Context(), deploymentID)
	if err != nil || dep == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "deployment not found"})
	}

	project, err := h.projectService.GetProject(c.Request().Context(), dep.ProjectID)
	if err != nil || project == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
	}

	settings, err := h.aiService.Get(c.Request().Context(), project.TeamID)
	if err != nil || settings == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "AI settings not configured for this team"})
	}

	if settings.APIKey == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "AI API Key is empty"})
	}

	client, model, err := h.buildAIClientAndModel(settings.Provider, settings.APIKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	prompt := fmt.Sprintf(
		"Analyze the following build/deployment logs and explain why it failed. Be concise and provide actionable steps to fix it:\n\n%s",
		dep.BuildLogs,
	)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert DevOps AI assistant. Your task is to diagnose server build failures.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(c.Request().Context(), req)
	if err != nil {
		c.Response().Write([]byte(fmt.Sprintf("0:\"Error connecting to AI Provider: %v\"\n", err)))
		return nil
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			break
		}

		chunk := response.Choices[0].Delta.Content
		if chunk != "" {
			c.Response().Write([]byte(fmt.Sprintf("0:%q\n", chunk)))
			c.Response().Flush()
		}
	}

	return nil
}

func (h *AIDiagnosticsHandler) buildAIClientAndModel(provider, apiKey string) (*openai.Client, string, error) {
	var client *openai.Client
	var model string

	switch provider {
	case ProviderOpenAI:
		client = openai.NewClient(apiKey)
		model = DefaultModelOpenAI
	case ProviderDeepseek:
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = BaseURLDeepseek
		client = openai.NewClientWithConfig(config)
		model = DefaultModelDeepseek
	case ProviderGroq:
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = BaseURLGroq
		client = openai.NewClientWithConfig(config)
		model = DefaultModelGroq
	default:
		return nil, "", fmt.Errorf("unsupported provider for streaming yet: %s", provider)
	}

	return client, model, nil
}
