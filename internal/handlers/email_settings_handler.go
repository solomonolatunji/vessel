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
// @Param teamId path string true "Team ID"
// @Router /teams/{teamId}/email_settings [get]
func (h *EmailSettingsHandler) GetTeamEmailSettings(c echo.Context) error {
	teamID := c.Param("teamId")
	settings, err := h.svc.GetTeamEmailSettings(c.Request().Context(), teamID)
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
// @Param teamId path string true "Team ID"
// @Param request body models.TeamEmailSettings true "Payload"
// @Router /teams/{teamId}/email_settings [put]
func (h *EmailSettingsHandler) SaveTeamEmailSettings(c echo.Context) error {
	teamID := c.Param("teamId")
	var req models.TeamEmailSettings
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	req.TeamID = teamID
	if err := h.svc.SaveTeamEmailSettings(c.Request().Context(), &req); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", map[string]string{"status": "ok"})
}
