package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services/notifications"
)

type DispatcherService struct {
	notifRepo    repositories.NotificationRepository
	settingsRepo repositories.SettingsRepository
	mailer       interface {
		SendTeamEmail(ctx context.Context, teamID, templateName string, toAddress string, subject string, data any) error
	}
}

func NewDispatcherService(notifRepo repositories.NotificationRepository, settingsRepo repositories.SettingsRepository, mailer interface {
	SendTeamEmail(ctx context.Context, teamID, templateName string, toAddress string, subject string, data any) error
}) *DispatcherService {
	return &DispatcherService{notifRepo: notifRepo, settingsRepo: settingsRepo, mailer: mailer}
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
				_ = notifications.SendSlackNotification(cfg.WebhookURL, event)
			}
		case "discord":
			var cfg struct {
				WebhookURL string `json:"webhookUrl"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.WebhookURL != "" {
				_ = notifications.SendDiscordNotification(cfg.WebhookURL, event)
			}
		case "telegram":
			var cfg struct {
				BotToken string `json:"botToken"`
				ChatID   string `json:"chatId"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.BotToken != "" && cfg.ChatID != "" {
				_ = notifications.SendTelegramNotification(cfg.BotToken, cfg.ChatID, event)
			}
		case "pushover":
			var cfg struct {
				AppToken string `json:"appToken"`
				UserKey  string `json:"userKey"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.AppToken != "" && cfg.UserKey != "" {
				_ = notifications.SendPushoverNotification(cfg.AppToken, cfg.UserKey, event)
			}
		case "generic":
			var cfg struct {
				WebhookURL string `json:"webhookUrl"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.WebhookURL != "" {
				_ = notifications.SendGenericWebhook(cfg.WebhookURL, event)
			}
		case "smtp":
			var cfg struct {
				ToEmail string `json:"toEmail"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.ToEmail != "" && d.mailer != nil {
				_ = d.mailer.SendTeamEmail(context.Background(), event.TeamID, "notification", cfg.ToEmail, event.Title, map[string]string{
					"Message": event.Message,
					"URL":     event.URL,
				})
			}
		}
	}
	return nil
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
			return notifications.SendDiscordNotification(settings.DiscordWebhookURL, event)
		}
	case "slack":
		if settings.SlackEnabled && settings.SlackWebhookURL != "" {
			return notifications.SendSlackNotification(settings.SlackWebhookURL, event)
		}
	case "telegram":
		if settings.TelegramEnabled && settings.TelegramBotToken != "" && settings.TelegramChatID != "" {
			return notifications.SendTelegramNotification(settings.TelegramBotToken, settings.TelegramChatID, event)
		}
	case "generic":
		if settings.GenericWebhookEnabled && settings.GenericWebhookURL != "" {
			return notifications.SendGenericWebhook(settings.GenericWebhookURL, event)
		}
	case "pushover":
		if settings.PushoverEnabled && settings.PushoverAPIToken != "" && settings.PushoverUserKey != "" {
			return notifications.SendPushoverNotification(settings.PushoverAPIToken, settings.PushoverUserKey, event)
		}
	}

	return fmt.Errorf("provider not enabled or configured")
}
