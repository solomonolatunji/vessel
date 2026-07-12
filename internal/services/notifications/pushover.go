package notifications

import (
	"bytes"
	"encoding/json"
	"net/http"

	"vessl.dev/vessl/internal/models"
)

func SendPushoverNotification(appToken, userKey string, event *models.NotificationEvent) error {
	url := "https://api.pushover.net/1/messages.json"
	payload := map[string]string{
		"token":   appToken,
		"user":    userKey,
		"title":   event.Title,
		"message": event.Message,
		"url":     event.URL,
	}
	body, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	return err
}
