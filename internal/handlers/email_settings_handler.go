package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if settings == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"configured": false})
	}

	settings.SMTPPassword = ""
	settings.ResendAPIKey = ""
	return c.JSON(http.StatusOK, settings)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	req.TeamID = teamID
	if err := h.svc.SaveTeamEmailSettings(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
