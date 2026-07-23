package notifications

import (
	"codedock.dev/codedock/internal/models"
)

func SendGenericWebhook(webhookURL string, event *models.NotificationEvent) error {
	return postJSON(webhookURL, event)
}
