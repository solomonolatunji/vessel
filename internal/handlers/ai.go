package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
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

	if req.Prompt == "" {
		return c.String(http.StatusBadRequest, "Prompt is required")
	}

	diagnosis, err := h.aiSettingsService.DiagnoseLogs(c.Request().Context(), req.Prompt)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, diagnosis)
}
