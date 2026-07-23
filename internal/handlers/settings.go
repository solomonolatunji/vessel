package handlers

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/utils"

	"codedock.dev/codedock/internal/services"
)

type SettingsHandler struct {
	settingsService      *services.SettingsService
	notifSettingsService *services.NotificationSettingsService
	mu                   sync.Mutex
}

func NewSettingsHandler(s *services.SettingsService, ns *services.NotificationSettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: s, notifSettingsService: ns}
}

func maskSettingsSecrets(s *models.ServerSettings) {
	if s.CloudflareAPIToken != "" {
		s.CloudflareAPIToken = "********"
	}
	if s.NamecheapAPIKey != "" {
		s.NamecheapAPIKey = "********"
	}
	if s.SpaceshipAPIKey != "" {
		s.SpaceshipAPIKey = "********"
	}
}

func (h *SettingsHandler) GetSettings(c echo.Context) error {
	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	masked := *s
	maskSettingsSecrets(&masked)
	return utils.Success(c, "Operation successful", masked)
}

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

	maskSettingsSecrets(existing)
	return utils.Success(c, "Operation successful", existing)
}
