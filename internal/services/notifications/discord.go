package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"vessl.dev/vessl/internal/models"
)

func SendDiscordNotification(webhookURL string, event *models.NotificationEvent) error {
	payload := map[string]string{
		"content": fmt.Sprintf("**%s**\n%s\n%s", event.Title, event.Message, event.URL),
	}
	body, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	return err
}
