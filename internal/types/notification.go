package types

import "time"

type NotificationIntegration struct {
	ID                 string `json:"id"`
	SMTPEnabled        bool   `json:"smtpEnabled"`
	SMTPHost           string `json:"smtpHost,omitempty"`
	SMTPPort           int    `json:"smtpPort,omitempty"`
	SMTPUser           string `json:"smtpUser,omitempty"`
	SMTPPassword       string `json:"smtpPassword,omitempty"`
	SMTPFromName       string `json:"smtpFromName,omitempty"`
	SMTPFromAddress    string `json:"smtpFromAddress,omitempty"`
	ResendEnabled      bool   `json:"resendEnabled"`
	ResendAPIKey       string `json:"resendApiKey,omitempty"`
	SlackEnabled       bool   `json:"slackEnabled"`
	SlackWebhookURL    string `json:"slackWebhookUrl,omitempty"`
	DiscordEnabled     bool   `json:"discordEnabled"`
	DiscordWebhookURL  string `json:"discordWebhookUrl,omitempty"`
	DiscordPingEnabled bool   `json:"discordPingEnabled"`
	TelegramEnabled    bool   `json:"telegramEnabled"`
	TelegramBotToken   string `json:"telegramBotToken,omitempty"`
	TelegramChatID     string `json:"telegramChatId,omitempty"`
	PushoverEnabled    bool   `json:"pushoverEnabled"`
	PushoverUserKey    string `json:"pushoverUserKey,omitempty"`
	PushoverAPIToken   string `json:"pushoverApiToken,omitempty"`
	WebhookEnabled     bool   `json:"webhookEnabled"`
	WebhookURL         string `json:"webhookUrl,omitempty"`
	UpdatedAt          string `json:"updatedAt"`
}

type ProjectNotificationPref struct {
	ProjectID       string    `json:"projectId"`
	EmailEnabled    bool      `json:"emailEnabled"`
	SlackEnabled    bool      `json:"slackEnabled"`
	DiscordEnabled  bool      `json:"discordEnabled"`
	TelegramEnabled bool      `json:"telegramEnabled"`
	PushoverEnabled bool      `json:"pushoverEnabled"`
	WebhookEnabled  bool      `json:"webhookEnabled"`
	Events          string    `json:"events"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type NotificationEvent struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	Level     string `json:"level"`
	ProjectID string `json:"projectId,omitempty"`
	URL       string `json:"url,omitempty"`
}
