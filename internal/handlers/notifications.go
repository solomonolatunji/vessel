package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type NotificationSettingsHandler struct {
	notifSettingsService *services.NotificationSettingsService
}

func NewNotificationSettingsHandler(s *services.NotificationSettingsService) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{notifSettingsService: s}
}

func (h *NotificationSettingsHandler) GetNotificationSettings(c echo.Context) error {
	s, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	
	masked := *s
	if masked.SMTPPassword != "" {
		masked.SMTPPassword = "********"
	}
	if masked.ResendAPIKey != "" {
		masked.ResendAPIKey = "********"
	}
	if masked.TelegramBotToken != "" {
		masked.TelegramBotToken = "********"
	}
	if masked.PushoverAPIToken != "" {
		masked.PushoverAPIToken = "********"
	}

	return utils.Success(c, "Operation successful", masked)
}

func (h *NotificationSettingsHandler) UpdateNotificationSettings(c echo.Context) error {
	existing, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to fetch existing notification settings")
	}

	realSMTP := existing.SMTPPassword
	realResend := existing.ResendAPIKey
	realTelegram := existing.TelegramBotToken
	realPushover := existing.PushoverAPIToken

	if err := c.Bind(existing); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	if existing.SMTPPassword == "********" {
		existing.SMTPPassword = realSMTP
	}
	if existing.ResendAPIKey == "********" {
		existing.ResendAPIKey = realResend
	}
	if existing.TelegramBotToken == "********" {
		existing.TelegramBotToken = realTelegram
	}
	if existing.PushoverAPIToken == "********" {
		existing.PushoverAPIToken = realPushover
	}

	if err := h.notifSettingsService.UpdateNotificationSettings(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Notification settings updated successfully", existing)
}
