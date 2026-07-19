package handlers

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/services"
)

type SettingsHandler struct {
	settingsService      *services.SettingsService
	notifSettingsService *services.NotificationSettingsService
	mu                   sync.Mutex
}

func NewSettingsHandler(s *services.SettingsService, ns *services.NotificationSettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: s, notifSettingsService: ns}
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

	masked := *s
	if masked.CloudflareAPIToken != "" {
		masked.CloudflareAPIToken = "********"
	}
	if masked.NamecheapAPIKey != "" {
		masked.NamecheapAPIKey = "********"
	}
	if masked.SpaceshipAPIKey != "" {
		masked.SpaceshipAPIKey = "********"
	}

	return utils.Success(c, "Operation successful", masked)
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

	notif, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	publicSettings := map[string]any{
		"registrationEnabled": s.RegistrationEnabled,
		"siteName":            s.SiteName,
		"emailEnabled":        notif.SMTPEnabled || notif.ResendEnabled,
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
	h.mu.Lock()
	defer h.mu.Unlock()

	existing, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to fetch existing settings")
	}

	realCloudflare := existing.CloudflareAPIToken
	realNamecheap := existing.NamecheapAPIKey
	realSpaceship := existing.SpaceshipAPIKey

	if err := c.Bind(existing); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if existing.CloudflareAPIToken == "********" {
		existing.CloudflareAPIToken = realCloudflare
	}
	if existing.NamecheapAPIKey == "********" {
		existing.NamecheapAPIKey = realNamecheap
	}
	if existing.SpaceshipAPIKey == "********" {
		existing.SpaceshipAPIKey = realSpaceship
	}

	if err := h.settingsService.UpdateSettings(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if existing.CloudflareAPIToken != "" {
		existing.CloudflareAPIToken = "********"
	}
	if existing.NamecheapAPIKey != "" {
		existing.NamecheapAPIKey = "********"
	}
	if existing.SpaceshipAPIKey != "" {
		existing.SpaceshipAPIKey = "********"
	}

	return utils.Success(c, "Operation successful", existing)
}
