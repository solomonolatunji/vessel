package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type AIDiagnoseRequest struct {
	Prompt string `json:"prompt"`
}

type AISettingsHandler struct {
	aiSettingsService *services.AISettingsService
}

func NewAISettingsHandler(s *services.AISettingsService) *AISettingsHandler {
	return &AISettingsHandler{aiSettingsService: s}
}

func (h *AISettingsHandler) GetAISettings(c echo.Context) error {
	s, err := h.aiSettingsService.GetAISettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", s)
}

func (h *AISettingsHandler) UpdateAISettings(c echo.Context) error {
	var req models.UpdateAISettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := h.aiSettingsService.UpdateAISettings(c.Request().Context(), &req.AISettings); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	updated, err := h.aiSettingsService.GetAISettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "AI settings updated successfully", updated)
}

func (h *AISettingsHandler) DiagnoseLogs(c echo.Context) error {
	var req AIDiagnoseRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	// This is a stubbed response for the AI diagnostics.
	// The Vercel AI SDK expects a plain text response for useCompletion if not streaming,
	// or Server-Sent Events if streaming. We will just return plain text for simplicity.

	// Real implementation would use h.aiSettingsService to grab an API key and
	// hit an LLM with the provided logs.

	mockDiagnosis := `I analyzed your logs. 
Based on the output provided, here is what I found:
1. No critical errors detected in the last few lines.
2. If there are connection issues, ensure your application is listening on the correct port (0.0.0.0 instead of 127.0.0.1).
3. Check your environment variables to ensure all database connections and secret keys are set correctly.

To fix this, you might need to adjust your build command or start script in the configuration.
`

	return c.String(http.StatusOK, mockDiagnosis)
}
