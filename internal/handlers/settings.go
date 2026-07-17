package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type SettingsHandler struct {
	settingsService *services.SettingsService
}

func NewSettingsHandler(s *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: s}
}

// @Summary GetSettings endpoint
// @Description GetSettings endpoint
// @Tags Settings
// @Accept json
// @Produce json
func (h *SettingsHandler) GetSettings(c echo.Context) error {
	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", s)
}

// @Summary GetPublicSettings endpoint
// @Description Get public settings for the frontend (e.g., if registration is enabled)
// @Tags Settings
// @Accept json
// @Produce json
// @Router /system/public [get]
func (h *SettingsHandler) GetPublicSettings(c echo.Context) error {
	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	publicSettings := map[string]any{
		"registrationEnabled": s.RegistrationEnabled,
		"siteName":            s.SiteName,
		"emailEnabled":        s.SMTPEnabled || s.ResendEnabled,
	}
	return utils.Success(c, "Operation successful", publicSettings)
}

// @Summary UpdateSettings endpoint
// @Description UpdateSettings endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body models.ServerSettings true "Payload"
// @Router /settings [put]
func (h *SettingsHandler) UpdateSettings(c echo.Context) error {
	var payload models.ServerSettings
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if err := h.settingsService.UpdateSettings(c.Request().Context(), &payload); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", payload)
}
