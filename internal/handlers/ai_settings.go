package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type AISettingsHandler struct {
	aiService *services.AISettingsService
}

func NewAISettingsHandler(s *services.AISettingsService) *AISettingsHandler {
	return &AISettingsHandler{aiService: s}
}

// @Summary Get endpoint
// @Description Get endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /api/teams/{teamId}/ai_settings [get]
func (h *AISettingsHandler) Get(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "team ID is required"})
	}

	settings, err := h.aiService.Get(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if settings == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Settings not found"})
	}
	return c.JSON(http.StatusOK, settings)
}

// @Summary Save endpoint
// @Description Save endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /api/teams/{teamId}/ai_settings [put]
func (h *AISettingsHandler) Save(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "team ID is required"})
	}

	var settings models.TeamAISettings
	if err := c.Bind(&settings); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	settings.TeamID = teamID

	if err := h.aiService.Save(c.Request().Context(), &settings); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, settings)
}
