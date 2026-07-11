package handlers

import (
	"net/http"

	"os"

	"github.com/labstack/echo/v4"

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

// @Summary GetSettings endpoint
// @Description GetSettings endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings [get]
func (h *SettingsHandler) GetSettings(c echo.Context) error {
	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

// @Summary UpdateSettings endpoint
// @Description UpdateSettings endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings [put]
func (h *SettingsHandler) UpdateSettings(c echo.Context) error {
	var payload models.ServerSettings
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.settingsService.UpdateSettings(c.Request().Context(), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, payload)
}

func (h *SettingsHandler) ListTeamNotificationChannels(c echo.Context) error {
	teamID := c.QueryParam("teamId")
	if teamID == "" {
		teamID = "default"
	}
	channels, err := h.settingsService.ListTeamNotificationChannels(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channels)
}

func (h *SettingsHandler) SaveTeamNotificationChannel(c echo.Context) error {
	var payload models.TeamNotificationChannel
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if payload.TeamID == "" {
		payload.TeamID = "default"
	}
	if err := h.settingsService.SaveTeamNotificationChannel(c.Request().Context(), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, payload)
}

// @Summary GetTeamNotificationChannel endpoint
// @Description GetTeamNotificationChannel endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/notifications/{id} [get]
func (h *SettingsHandler) GetTeamNotificationChannel(c echo.Context) error {
	id := c.Param("id")
	channel, err := h.settingsService.GetTeamNotificationChannel(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channel)
}

func (h *SettingsHandler) DeleteTeamNotificationChannel(c echo.Context) error {
	id := c.Param("id")
	if err := h.settingsService.DeleteTeamNotificationChannel(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// @Summary Activate License endpoint
// @Description Activates offline license key
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/license [post]
func (h *SettingsHandler) ActivateLicense(c echo.Context) error {
	var payload struct {
		LicenseKey string `json:"license_key"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	pubKey := os.Getenv("LICENSE_PUBLIC_KEY")
	if pubKey == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "license public key not configured on this instance"})
	}

	claims, err := license.VerifyLicense(pubKey, payload.LicenseKey)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load settings"})
	}

	s.LicenseKey = payload.LicenseKey
	s.Plan = claims.Plan
	s.MaxSeats = claims.MaxSeats

	if err := h.settingsService.UpdateSettings(c.Request().Context(), s); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save license to settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "activated",
		"plan":   claims.Plan,
	})
}
