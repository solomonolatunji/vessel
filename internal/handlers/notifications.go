package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type NotificationSettingsHandler struct {
	notifSettingsService *services.NotificationSettingsService
}

func NewNotificationSettingsHandler(s *services.NotificationSettingsService) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{notifSettingsService: s}
}

func maskNotificationSecrets(s *models.NotificationSettings) {
	if s.SMTPPassword != "" {
		s.SMTPPassword = "********"
	}
	if s.ResendAPIKey != "" {
		s.ResendAPIKey = "********"
	}
	if s.TelegramBotToken != "" {
		s.TelegramBotToken = "********"
	}
	if s.PushoverAPIToken != "" {
		s.PushoverAPIToken = "********"
	}
	if s.SlackWebhookURL != "" {
		s.SlackWebhookURL = "********"
	}
}

func (h *NotificationSettingsHandler) GetNotificationSettings(c echo.Context) error {
	s, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	masked := *s
	maskNotificationSecrets(&masked)
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
	realSlack := existing.SlackWebhookURL

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
	if existing.SlackWebhookURL == "********" {
		existing.SlackWebhookURL = realSlack
	}

	if err := h.notifSettingsService.UpdateNotificationSettings(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	maskNotificationSecrets(existing)
	return utils.Success(c, "Notification settings updated successfully", existing)
}
