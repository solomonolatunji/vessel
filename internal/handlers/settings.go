package handlers

import (
	"net/http"

	"os"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/license"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type SettingsHandler struct {
	settingsService *services.SettingsService
}

func NewSettingsHandler(s *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: s}
}

type ActivateLicenseRequest struct {
	LicenseKey string `json:"license_key"`
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

func (h *SettingsHandler) ListTeamNotificationChannels(c echo.Context) error {
	teamID := c.QueryParam("teamId")
	if teamID == "" {
		teamID = "default"
	}
	channels, err := h.settingsService.ListTeamNotificationChannels(c.Request().Context(), teamID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", channels)
}

func (h *SettingsHandler) SaveTeamNotificationChannel(c echo.Context) error {
	var payload models.TeamNotificationChannel
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if payload.TeamID == "" {
		payload.TeamID = "default"
	}
	if err := h.settingsService.SaveTeamNotificationChannel(c.Request().Context(), &payload); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", payload)
}

// @Summary GetTeamNotificationChannel endpoint
// @Description GetTeamNotificationChannel endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /settings/notifications/{id} [get]
func (h *SettingsHandler) GetTeamNotificationChannel(c echo.Context) error {
	id := c.Param("id")
	channel, err := h.settingsService.GetTeamNotificationChannel(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", channel)
}

func (h *SettingsHandler) DeleteTeamNotificationChannel(c echo.Context) error {
	id := c.Param("id")
	if err := h.settingsService.DeleteTeamNotificationChannel(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}

// @Summary Activate License endpoint
// @Description Activates offline license key
// @Tags Settings
// @Accept json
// @Produce json
// @Summary Activate License
// @Description Activate License
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body handlers.ActivateLicenseRequest true "Payload"
// @Router /settings/license [post]
func (h *SettingsHandler) ActivateLicense(c echo.Context) error {
	var payload ActivateLicenseRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	pubKey := os.Getenv("LICENSE_PUBLIC_KEY")
	if pubKey == "" {
		return utils.Error(c, http.StatusInternalServerError, "license public key not configured on this instance")
	}

	claims, err := license.VerifyLicense(pubKey, payload.LicenseKey)
	if err != nil {
		return utils.Error(c, http.StatusForbidden, err.Error())
	}

	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to load settings")
	}

	s.LicenseKey = payload.LicenseKey
	s.Plan = claims.Plan
	s.MaxSeats = claims.MaxSeats

	if err := h.settingsService.UpdateSettings(c.Request().Context(), s); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to save license to settings")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "activated",
		"plan":   claims.Plan,
	})
}
