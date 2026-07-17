package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

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
