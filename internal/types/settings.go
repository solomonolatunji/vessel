package types

import "time"

type ServerSettings struct {
	ID                 string `json:"id"` // singleton "global"
	CaddyWildcardIP    string `json:"caddyWildcardIp"`
	DiscordWebhookURL  string `json:"discordWebhookUrl,omitempty"`
	SlackWebhookURL    string `json:"slackWebhookUrl,omitempty"`
	TelegramBotToken   string `json:"telegramBotToken,omitempty"`
	TelegramChatID     string `json:"telegramChatId,omitempty"`
	SMTPHost           string `json:"smtpHost,omitempty"`
	SMTPPort           int    `json:"smtpPort,omitempty"`
	SMTPUser           string `json:"smtpUser,omitempty"`
	SMTPPassword       string `json:"smtpPassword,omitempty"`
	SMTPFromName       string `json:"smtpFromName,omitempty"`
	SMTPFromAddress    string `json:"smtpFromAddress,omitempty"`
	NotificationAlerts bool   `json:"notificationAlerts"` // enable/disable global push alerts

	// Advanced Server Settings
	RegistrationEnabled  bool   `json:"registrationEnabled"`
	CustomDNSResolvers   string `json:"customDnsResolvers"` // comma-separated e.g. "1.1.1.1,8.8.8.8"
	DNSValidationEnabled bool   `json:"dnsValidationEnabled"`
	IPAllowlist          string `json:"ipAllowlist"`      // comma-separated IPs or CIDRs e.g. "192.168.1.100,10.0.0.0/8"
	MCPServerEnabled     bool   `json:"mcpServerEnabled"` // toggle for AI agent MCP integrations

	// Update Management
	UpdateCheckCron   string `json:"updateCheckCron"`   // cron expression e.g. "0 * * * *"
	AutoUpdateEnabled bool   `json:"autoUpdateEnabled"` // toggle for auto-deploying updates
	CurrentVersion    string `json:"currentVersion"`
	LatestVersion     string `json:"latestVersion"`
	LastUpdateCheck   string `json:"lastUpdateCheck"`

	UpdatedAt string `json:"updatedAt"`
}

type PersonalAccessToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	TokenHash string    `json:"-"`
	Prefix    string    `json:"prefix"` // e.g. vsl_user_
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}
