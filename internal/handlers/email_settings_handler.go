package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type EmailSettingsHandler struct {
	svc *services.EmailSettingsService
}

func NewEmailSettingsHandler(svc *services.EmailSettingsService) *EmailSettingsHandler {
	return &EmailSettingsHandler{svc: svc}
}

// @Summary Get Team Email Settings
// @Description Get Team Email Settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param workspaceId path string true "Team ID"
// @Router /workspaces/{workspaceId}/email_settings [get]
func (h *EmailSettingsHandler) GetWorkspaceEmailSettings(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	settings, err := h.svc.GetWorkspaceEmailSettings(c.Request().Context(), workspaceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if settings == nil {
		return utils.Success(c, "Operation successful", map[string]interface{}{"configured": false})
	}

	settings.SMTPPassword = ""
	settings.ResendAPIKey = ""
	return utils.Success(c, "Operation successful", settings)
}

// @Summary Save Team Email Settings
// @Description Save Team Email Settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param workspaceId path string true "Team ID"
// @Param request body models.WorkspaceEmailSettings true "Payload"
// @Router /workspaces/{workspaceId}/email_settings [put]
func (h *EmailSettingsHandler) SaveWorkspaceEmailSettings(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	var req models.WorkspaceEmailSettings
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	req.WorkspaceID = workspaceID
	if err := h.svc.SaveWorkspaceEmailSettings(c.Request().Context(), &req); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", map[string]string{"status": "ok"})
}
