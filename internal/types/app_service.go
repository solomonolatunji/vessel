package types

import "time"

type AppServiceConfig struct {
	ID                string          `json:"id"`
	ProjectID         string          `json:"projectId"`
	EnvironmentID     string          `json:"environmentId"`
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
	Status            string          `json:"status"`
	ContainerID       string          `json:"containerId"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}
