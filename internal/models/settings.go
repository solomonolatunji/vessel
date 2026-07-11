package models

import "time"

type TeamAISettings struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Provider  string    `json:"provider"`
	APIKey    string    `json:"apiKey,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TeamEmailSettings struct {
	ID              string    `json:"id"`
	TeamID          string    `json:"teamId"`
	SMTPHost        string    `json:"smtpHost,omitempty"`
	SMTPPort        int       `json:"smtpPort,omitempty"`
	SMTPUser        string    `json:"smtpUser,omitempty"`
	SMTPPassword    string    `json:"smtpPassword,omitempty"`
	SMTPFromName    string    `json:"smtpFromName,omitempty"`
	SMTPFromAddress string    `json:"smtpFromAddress,omitempty"`
	ResendAPIKey    string    `json:"resendApiKey,omitempty"`
	UseResend       bool      `json:"useResend"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type ServerSettings struct {
	ID                          string `json:"id"`
	CaddyWildcardIP             string `json:"caddyWildcardIp"`
	DiscordWebhookURL           string `json:"discordWebhookUrl,omitempty"`
	DiscordPingEnabled          bool   `json:"discordPingEnabled"`
	DiscordEnabled              bool   `json:"discordEnabled"`
	SlackWebhookURL             string `json:"slackWebhookUrl,omitempty"`
	SlackEnabled                bool   `json:"slackEnabled"`
	TelegramBotToken            string `json:"telegramBotToken,omitempty"`
	TelegramChatID              string `json:"telegramChatId,omitempty"`
	TelegramEnabled             bool   `json:"telegramEnabled"`
	SMTPHost                    string `json:"smtpHost,omitempty"`
	SMTPPort                    int    `json:"smtpPort,omitempty"`
	SMTPUser                    string `json:"smtpUser,omitempty"`
	SMTPPassword                string `json:"smtpPassword,omitempty"`
	SMTPFromName                string `json:"smtpFromName,omitempty"`
	SMTPFromAddress             string `json:"smtpFromAddress,omitempty"`
	SMTPEnabled                 bool   `json:"smtpEnabled"`
	ResendAPIKey                string `json:"resendApiKey,omitempty"`
	ResendEnabled               bool   `json:"resendEnabled"`
	PushoverUserKey             string `json:"pushoverUserKey,omitempty"`
	PushoverAPIToken            string `json:"pushoverApiToken,omitempty"`
	PushoverEnabled             bool   `json:"pushoverEnabled"`
	GenericWebhookURL           string `json:"genericWebhookUrl,omitempty"`
	GenericWebhookEnabled       bool   `json:"genericWebhookEnabled"`
	NotificationAlerts          bool   `json:"notificationAlerts"`
	RegistrationEnabled         bool   `json:"registrationEnabled"`
	RegistrationDomainAllowlist string `json:"registrationDomainAllowlist,omitempty"`
	CustomDNSResolvers          string `json:"customDnsResolvers"`
	DNSValidationEnabled        bool   `json:"dnsValidationEnabled"`
	IPAllowlist                 string `json:"ipAllowlist"`
	MCPServerEnabled            bool   `json:"mcpServerEnabled"`
	DefaultWildcardDomain       string `json:"defaultWildcardDomain,omitempty"`
	DefaultOpenAIKey            string `json:"defaultOpenAIKey,omitempty"`
	DefaultAnthropicKey         string `json:"defaultAnthropicKey,omitempty"`
	UpdateCheckCron             string `json:"updateCheckCron"`
	AutoUpdateEnabled           bool   `json:"autoUpdateEnabled"`
	TelemetryEnabled            bool   `json:"telemetryEnabled"`
	LicenseKey                  string `json:"licenseKey,omitempty"`
	Plan                        string `json:"plan"`
	MaxSeats                    int    `json:"maxSeats"`
	CurrentVersion              string `json:"currentVersion"`
	LatestVersion               string `json:"latestVersion"`
	LastUpdateCheck             string `json:"lastUpdateCheck"`
	UpdatedAt                   string `json:"updatedAt"`
}

type UpdateSettingsRequest struct {
	ServerSettings
}

type PruneResponse struct {
	Status              string `json:"status"`
	Message             string `json:"message"`
	SpaceReclaimedBytes uint64 `json:"spaceReclaimedBytes"`
}

type MCPResponse struct {
	JSONRPC      string           `json:"jsonrpc"`
	ID           any              `json:"id,omitempty"`
	Result       any              `json:"result,omitempty"`
	Error        *MCPError        `json:"error,omitempty"`
	Server       map[string]any   `json:"server,omitempty"`
	Tools        []map[string]any `json:"tools,omitempty"`
	Capabilities map[string]any   `json:"capabilities,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
