package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

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
// @Router /teams/{teamId}/ai_settings [get]
func (h *AISettingsHandler) Get(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return utils.Error(c, http.StatusBadRequest, "team ID is required")
	}

	settings, err := h.aiService.Get(c.Request().Context(), teamID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if settings == nil {
		return utils.Error(c, http.StatusNotFound, "Settings not found")
	}
	return utils.Success(c, "Operation successful", settings)
}

// @Summary Save endpoint
// @Description Save endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Param request body models.TeamAISettings true "Payload"
// @Router /teams/{teamId}/ai_settings [put]
func (h *AISettingsHandler) Save(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return utils.Error(c, http.StatusBadRequest, "team ID is required")
	}

	var settings models.TeamAISettings
	if err := c.Bind(&settings); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	settings.TeamID = teamID

	if err := h.aiService.Save(c.Request().Context(), &settings); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", settings)
}
