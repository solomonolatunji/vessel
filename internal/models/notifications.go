package models

type NotificationSettings struct {
	ID                    string `json:"id"`
	DiscordWebhookURL     string `json:"discordWebhookUrl,omitempty"`
	DiscordPingEnabled    bool   `json:"discordPingEnabled"`
	DiscordEnabled        bool   `json:"discordEnabled"`
	SlackWebhookURL       string `json:"slackWebhookUrl,omitempty"`
	SlackEnabled          bool   `json:"slackEnabled"`
	TelegramBotToken      string `json:"telegramBotToken,omitempty"`
	TelegramChatID        string `json:"telegramChatId,omitempty"`
	TelegramEnabled       bool   `json:"telegramEnabled"`
	SMTPHost              string `json:"smtpHost,omitempty"`
	SMTPPort              int    `json:"smtpPort,omitempty"`
	SMTPUser              string `json:"smtpUser,omitempty"`
	SMTPPassword          string `json:"smtpPassword,omitempty"`
	SMTPFromName          string `json:"smtpFromName,omitempty"`
	SMTPFromAddress       string `json:"smtpFromAddress,omitempty"`
	SMTPEnabled           bool   `json:"smtpEnabled"`
	ResendAPIKey          string `json:"resendApiKey,omitempty"`
	ResendEnabled         bool   `json:"resendEnabled"`
	PushoverUserKey       string `json:"pushoverUserKey,omitempty"`
	PushoverAPIToken      string `json:"pushoverApiToken,omitempty"`
	PushoverEnabled       bool   `json:"pushoverEnabled"`
	GenericWebhookURL     string `json:"genericWebhookUrl,omitempty"`
	GenericWebhookEnabled bool   `json:"genericWebhookEnabled"`
	NotificationAlerts    bool   `json:"notificationAlerts"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

type UpdateNotificationSettingsRequest struct {
	NotificationSettings
}
