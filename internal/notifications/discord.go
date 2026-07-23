package notifications

import (
	"fmt"

	"codedock.dev/codedock/internal/models"
)

func SendDiscordNotification(webhookURL string, event *models.NotificationEvent) error {
	return postJSON(webhookURL, map[string]string{
		"content": fmt.Sprintf("**%s**\n%s\n%s", event.Title, event.Message, event.URL),
	})
}
