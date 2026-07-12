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
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Router /workspaces/{workspaceId}/ai_settings [get]
func (h *AISettingsHandler) Get(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "team ID is required")
	}

	settings, err := h.aiService.Get(c.Request().Context(), workspaceID)
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
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Param request body models.WorkspaceAISettings true "Payload"
// @Router /workspaces/{workspaceId}/ai_settings [put]
func (h *AISettingsHandler) Save(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "team ID is required")
	}

	var settings models.WorkspaceAISettings
	if err := c.Bind(&settings); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	settings.WorkspaceID = workspaceID

	if err := h.aiService.Save(c.Request().Context(), &settings); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", settings)
}
