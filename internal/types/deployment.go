package types

import "time"

// DeploymentRecord represents a specific deployment run for a service inside an environment.
type DeploymentRecord struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"serviceId"`
	EnvironmentID string    `json:"environmentId"`
	ProjectID     string    `json:"projectId"`
	Status        string    `json:"status"` // "ACTIVE", "BUILDING", "FAILED", "SLEPT", "REMOVED", "QUEUED", "NEEDS_APPROVAL"
	CommitHash    string    `json:"commitHash,omitempty"`
	CommitMessage string    `json:"commitMessage,omitempty"`
	Branch        string    `json:"branch,omitempty"`
	Trigger       string    `json:"trigger,omitempty"` // e.g. "Git Push", "Manual", "API"
	BuildLogs     string    `json:"buildLogs,omitempty"`
	ContainerID   string    `json:"containerId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	FinishedAt    time.Time `json:"finishedAt,omitempty"`
}

// ServiceMetric represents time-series telemetry data for a deployed container instance.
type ServiceMetric struct {
	Timestamp  string  `json:"timestamp"`
	CPUPercent float64 `json:"cpuPercent"`
	MemoryMB   float64 `json:"memoryMB"`
	NetworkRx  float64 `json:"networkRxKB"`
	NetworkTx  float64 `json:"networkTxKB"`
}
