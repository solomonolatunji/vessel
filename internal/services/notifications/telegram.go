package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"vessl.dev/vessl/internal/models"
)

func SendTelegramNotification(botToken, chatID string, event *models.NotificationEvent) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := map[string]string{
		"chat_id":    chatID,
		"text":       fmt.Sprintf("*%s*\n%s\n%s", event.Title, event.Message, event.URL),
		"parse_mode": "Markdown",
	}
	body, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	return err
}
