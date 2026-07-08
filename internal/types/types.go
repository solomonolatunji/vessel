package types

import "time"

// ContainerHealth tracks real-time resource utilization and status for running containers.
type ContainerHealth struct {
	Status             string  `json:"status"`
	CPUUsagePercentage float64 `json:"cpuUsagePercentage"`
	MemoryUsageBytes   int64   `json:"memoryUsageBytes"`
	MemoryLimitBytes   int64   `json:"memoryLimitBytes"`
	UptimeSeconds      int64   `json:"uptimeSeconds"`
}

// ProjectConfig stores the core application configuration, build rules, and runtime settings.
type ProjectConfig struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	RepositoryURL     string          `json:"repositoryUrl,omitempty"`
	Branch            string          `json:"branch,omitempty"`
	BuildCommand      string          `json:"buildCommand,omitempty"`
	StartCommand      string          `json:"startCommand,omitempty"`
	DockerfilePath    string          `json:"dockerfilePath,omitempty"`
	InternalPort      int             `json:"internalPort"`
	Domain            string          `json:"domain,omitempty"`
	EnvVarsCount      int             `json:"envVarsCount"`
	AutoDeployWebhook bool            `json:"autoDeployWebhook"`
	CPURequest        float64         `json:"cpuRequest,omitempty"`
	MemoryLimitMB     int             `json:"memoryLimitMB,omitempty"`
	HealthCheckPath   string          `json:"healthCheckPath,omitempty"`
	Health            ContainerHealth `json:"health"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}

// DomainConfig manages custom domain routing, SSL certificate issuance state, and Caddy integration.
type DomainConfig struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	DomainName    string    `json:"domainName"`
	RedirectTo    string    `json:"redirectTo,omitempty"`
	SSLCertStatus string    `json:"sslCertStatus"`
	PathPrefix    string    `json:"pathPrefix"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// User represents an authenticated workspace member of the Vessel control plane.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Invite tracks pending workspace role invitations and expiration metadata.
type Invite struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	Token      string     `json:"token"`
	InvitedBy  string     `json:"invitedBy"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	AcceptedAt *time.Time `json:"acceptedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}

// EnvVar represents an encrypted environment variable record stored in SQLite.
type EnvVar struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"projectId"`
	Key            string    `json:"key"`
	EncryptedValue string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// SystemInfo delivers host health, Docker daemon metrics, and self-update alerts.
type SystemInfo struct {
	Version         string `json:"version"`
	GoVersion       string `json:"goVersion"`
	DockerVersion   string `json:"dockerVersion"`
	CaddyVersion    string `json:"caddyVersion"`
	OS              string `json:"os"`
	Arch            string `json:"arch"`
	TotalMemoryMB   int64  `json:"totalMemoryMB"`
	FreeMemoryMB    int64  `json:"freeMemoryMB"`
	CPUCores        int    `json:"cpuCores"`
	UpdateAvailable bool   `json:"updateAvailable"`
	LatestVersion   string `json:"latestVersion,omitempty"`
}
