package notifications

import (
	"codedock.dev/codedock/internal/models"
)

func SendPushoverNotification(appToken, userKey string, event *models.NotificationEvent) error {
	return postJSON("https://api.pushover.net/1/messages.json", map[string]string{
		"token":   appToken,
		"user":    userKey,
		"title":   event.Title,
		"message": event.Message,
		"url":     event.URL,
	})
}
