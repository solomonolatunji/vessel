package notifications

import (
	"bytes"
	"encoding/json"
	"net/http"

	"vessl.dev/vessl/internal/models"
)

func SendGenericWebhook(webhookURL string, event *models.NotificationEvent) error {
	body, _ := json.Marshal(event)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	return err
}
