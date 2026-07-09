package types

import "time"

type AppServiceConfig struct {
	ID                string          `json:"id"`
	ProjectID         string          `json:"projectId"`
	EnvironmentID     string          `json:"environmentId"`
	Name              string          `json:"name"`
	Icon              string          `json:"icon,omitempty"`
	RepositoryURL     string          `json:"repositoryUrl,omitempty"`
	Branch            string          `json:"branch,omitempty"`
	RootDirectory     string          `json:"rootDirectory,omitempty"`
	BuildCommand      string          `json:"buildCommand,omitempty"`
	StartCommand      string          `json:"startCommand,omitempty"`
	DockerfilePath    string          `json:"dockerfilePath,omitempty"`
	InternalPort      int             `json:"internalPort"`
	Domain            string          `json:"domain,omitempty"`
	EnvVarsCount      int             `json:"envVarsCount"`
	AutoDeployWebhook bool            `json:"autoDeployWebhook"`
	CPURequest        float64         `json:"cpuRequest,omitempty"`
	MemoryLimitMB     int             `json:"memoryLimitMB,omitempty"`
	Replicas          int             `json:"replicas,omitempty"`
	RestartPolicy     string          `json:"restartPolicy,omitempty"`
	TeardownTimeout   int             `json:"teardownTimeout,omitempty"`
	Serverless        bool            `json:"serverless,omitempty"`
	CronSchedule      string          `json:"cronSchedule,omitempty"`
	HealthCheckPath   string          `json:"healthCheckPath,omitempty"`
	Health            ContainerHealth `json:"health"`
	Status            string          `json:"status"`
	ContainerID       string          `json:"containerId"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}
