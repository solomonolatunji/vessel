package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type DispatcherService struct {
	notifRepo    repositories.NotificationRepository
	settingsRepo repositories.SettingsRepository
}

func NewDispatcherService(notifRepo repositories.NotificationRepository, settingsRepo repositories.SettingsRepository) *DispatcherService {
	return &DispatcherService{notifRepo: notifRepo, settingsRepo: settingsRepo}
}

func (d *DispatcherService) Dispatch(event *models.NotificationEvent) {
	go func() {
		if err := d.Send(event); err != nil {
			log.Printf("[Dispatcher] Failed to dispatch event '%s': %v", event.Title, err)
		}
	}()
}

func (d *DispatcherService) Send(event *models.NotificationEvent) error {
	if event.TeamID == "" {
		return fmt.Errorf("TeamID is required for dispatch")
	}

	if event.TeamID == "global_test" {
		return d.sendGlobalTest(event)
	}

	channels, err := d.notifRepo.ListChannelsByTeam(context.Background(), event.TeamID)
	if err != nil {
		return fmt.Errorf("failed to list channels for team %s: %w", event.TeamID, err)
	}

	for _, c := range channels {
		if !c.IsEnabled {
			continue
		}

		var events []string
		if len(c.Events) > 0 {
			if err := json.Unmarshal(c.Events, &events); err == nil {
				matches := false
				for _, e := range events {
					if e == event.EventType || e == "*" {
						matches = true
						break
					}
				}
				if !matches && len(events) > 0 {
					continue
				}
			}
		}

		switch c.Provider {
		case "slack":
			var cfg struct {
				WebhookURL string `json:"webhookUrl"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.WebhookURL != "" {
				_ = d.sendWebhook(cfg.WebhookURL, event)
			}
		case "discord":
			var cfg struct {
				WebhookURL string `json:"webhookUrl"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.WebhookURL != "" {
				_ = d.sendWebhook(cfg.WebhookURL, event)
			}
		case "smtp":
		}
	}
	return nil
}

func (d *DispatcherService) sendWebhook(webhookURL string, event *models.NotificationEvent) error {
	payload := map[string]string{
		"content": fmt.Sprintf("**%s**\n%s\n%s", event.Title, event.Message, event.URL),
	}
	body, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	return err
}

func (d *DispatcherService) sendGlobalTest(event *models.NotificationEvent) error {
	ctx := context.Background()
	settings, err := d.settingsRepo.GetServerSettings(ctx)
	if err != nil || settings == nil {
		return fmt.Errorf("could not fetch server settings: %v", err)
	}

	provider := event.EventType[len("test_global_"):]

	switch provider {
	case "discord":
		if settings.DiscordEnabled && settings.DiscordWebhookURL != "" {
			return d.sendWebhook(settings.DiscordWebhookURL, event)
		}
	case "slack":
		if settings.SlackEnabled && settings.SlackWebhookURL != "" {
			return d.sendWebhook(settings.SlackWebhookURL, event)
		}
	case "telegram":
		if settings.TelegramEnabled && settings.TelegramBotToken != "" && settings.TelegramChatID != "" {
			url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", settings.TelegramBotToken)
			payload := map[string]string{
				"chat_id": settings.TelegramChatID,
				"text":    fmt.Sprintf("**%s**\n%s\n%s", event.Title, event.Message, event.URL),
			}
			body, _ := json.Marshal(payload)
			_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
			return err
		}
	case "generic":
		if settings.GenericWebhookEnabled && settings.GenericWebhookURL != "" {
			return d.sendWebhook(settings.GenericWebhookURL, event)
		}
	}

	return fmt.Errorf("provider not enabled or configured")
}
