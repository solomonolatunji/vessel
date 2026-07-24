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
	AppID    string            `json:"app_id"`
	Image    string            `json:"image"`
	Env      map[string]string `json:"env"`
	Ports    []string          `json:"ports"`
	Volumes  []string          `json:"volumes"`
	Network  string            `json:"network"`
	Registry *RegistryConfig   `json:"registry,omitempty"`
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
