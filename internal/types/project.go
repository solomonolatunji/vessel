package types

import "time"

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
