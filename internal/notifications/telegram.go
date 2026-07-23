package notifications

import (
	"fmt"

	"codedock.dev/codedock/internal/models"
)

func SendTelegramNotification(botToken, chatID string, event *models.NotificationEvent) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	return postJSON(url, map[string]string{
		"chat_id":    chatID,
		"text":       fmt.Sprintf("*%s*\n%s\n%s", event.Title, event.Message, event.URL),
		"parse_mode": "Markdown",
	})
}
