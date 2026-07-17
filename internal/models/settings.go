package models

type ServerSettings struct {
	ID                          string `json:"id"`
	TraefikWildcardIP           string `json:"traefikWildcardIp"`
	DiscordWebhookURL           string `json:"discordWebhookUrl,omitempty"`
	RegistrationEnabled         bool   `json:"registrationEnabled"`
	RegistrationDomainAllowlist string `json:"registrationDomainAllowlist,omitempty"`
	CustomDNSResolvers          string `json:"customDnsResolvers"`
	DNSValidationEnabled        bool   `json:"dnsValidationEnabled"`
	IPAllowlist                 string `json:"ipAllowlist"`
	MCPServerEnabled            bool   `json:"mcpServerEnabled"`
	DefaultWildcardDomain       string `json:"defaultWildcardDomain,omitempty"`
	PanelDomain                 string `json:"dashboardDomain,omitempty"`
	SiteName                    string `json:"siteName,omitempty"`
	PublicIPv4                  string `json:"publicIpv4,omitempty"`
	PublicIPv6                  string `json:"publicIpv6,omitempty"`
	ShowSponsorshipPopup        bool   `json:"showSponsorshipPopup"`
	DisableTwoStepConfirmation  bool   `json:"disableTwoStepConfirmation"`
	CloudflareAPIToken          string `json:"cloudflareApiToken,omitempty"`
	NamecheapAPIUser            string `json:"namecheapApiUser,omitempty"`
	NamecheapAPIKey             string `json:"namecheapApiKey,omitempty"`
	NamecheapClientIP           string `json:"namecheapClientIp,omitempty"`
	SpaceshipAPIKey             string `json:"spaceshipApiKey,omitempty"`
	UpdateCheckCron             string `json:"updateCheckCron"`
	AutoUpdateEnabled           bool   `json:"autoUpdateEnabled"`
	TelemetryEnabled            bool   `json:"telemetryEnabled"`
	ConcurrentBuilds            int    `json:"concurrentBuilds"`
	DeploymentTimeout           int    `json:"deploymentTimeout"`
	ServerTimezone              string `json:"serverTimezone"`
	DockerCleanupCron           string `json:"dockerCleanupCron"`
	DiskUsageThreshold          int    `json:"diskUsageThreshold"`
	DiskUsageCron               string `json:"diskUsageCron"`
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
