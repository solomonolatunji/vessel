package models

import "time"

// WorkerMessageType defines the kind of message being sent over the WebSocket.
type WorkerMessageType string

const (
	WorkerMessageTypeDeployApp  WorkerMessageType = "deploy_app"
	WorkerMessageTypeStopApp    WorkerMessageType = "stop_app"
	WorkerMessageTypeRestartApp WorkerMessageType = "restart_app"
	WorkerMessageTypeDeployDB   WorkerMessageType = "deploy_db"

	WorkerMessageTypeAuth       WorkerMessageType = "auth"
	WorkerMessageTypeAuthResult WorkerMessageType = "auth_result"
	WorkerMessageTypeLogStream  WorkerMessageType = "log_stream"
	WorkerMessageTypeMetrics    WorkerMessageType = "metrics"
	WorkerMessageTypeCommandAck WorkerMessageType = "command_ack"
)

type WorkerMessage struct {
	ID        string            `json:"id"`
	Type      WorkerMessageType `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	Payload   []byte            `json:"payload"`
}

type WorkerAuthPayload struct {
	WorkerToken string `json:"worker_token"`
	Version     string `json:"version"`
}

type WorkerAuthResultPayload struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type WorkerDeployAppPayload struct {
	AppID           string            `json:"app_id"`
	DeploymentID    string            `json:"deployment_id,omitempty"`
	Image           string            `json:"image,omitempty"`
	Env             map[string]string `json:"env"`
	Ports           []string          `json:"ports"`
	Volumes         []string          `json:"volumes"`
	Network         string            `json:"network"`
	Registry        *RegistryConfig   `json:"registry,omitempty"`
	
	// Build specific (Decentralized Builds)
	GitRepoURL      string `json:"git_repo_url,omitempty"`
	GitBranch       string `json:"git_branch,omitempty"`
	GitCommitHash   string `json:"git_commit_hash,omitempty"`
	GitAuthToken    string `json:"git_auth_token,omitempty"`
	BuildCommand    string `json:"build_command,omitempty"`
	InstallCommand  string `json:"install_command,omitempty"`
	StartCommand    string `json:"start_command,omitempty"`
	BaseDirectory   string `json:"base_directory,omitempty"`
	NixpacksVersion string `json:"nixpacks_version,omitempty"`
	MemoryLimitMB   int    `json:"memory_limit_mb,omitempty"`
	CPURequest      int    `json:"cpu_request,omitempty"`
	
	Domain          string `json:"domain,omitempty"`
	RuntimeMode     string `json:"runtime_mode,omitempty"`
	HealthCheckPath string `json:"health_check_path,omitempty"`
}

type RegistryConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type WorkerCommandAckPayload struct {
	CommandID string `json:"command_id"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

type WorkerLogStreamPayload struct {
	ContainerID string `json:"container_id"`
	LogLine     string `json:"log_line"`
	StreamType  string `json:"stream_type"`
}

type WorkerMetricsPayload struct {
	CPUUsagePercentage float64 `json:"cpu_usage_percentage"`
	MemoryUsageBytes   uint64  `json:"memory_usage_bytes"`
	MemoryLimitBytes   uint64  `json:"memory_limit_bytes"`
	DiskUsageBytes     uint64  `json:"disk_usage_bytes"`
	DiskTotalBytes     uint64  `json:"disk_total_bytes"`
}
