package models

import "time"

// WorkerMessageType defines the kind of message being sent over the WebSocket.
type WorkerMessageType string

const (
	// Control Plane -> Worker Commands
	WorkerMessageTypeDeployApp  WorkerMessageType = "deploy_app"
	WorkerMessageTypeStopApp    WorkerMessageType = "stop_app"
	WorkerMessageTypeRestartApp WorkerMessageType = "restart_app"
	WorkerMessageTypeDeployDB   WorkerMessageType = "deploy_db"

	// Worker -> Control Plane Responses/Events
	WorkerMessageTypeAuth       WorkerMessageType = "auth"
	WorkerMessageTypeAuthResult WorkerMessageType = "auth_result"
	WorkerMessageTypeLogStream  WorkerMessageType = "log_stream"
	WorkerMessageTypeMetrics    WorkerMessageType = "metrics"
	WorkerMessageTypeCommandAck WorkerMessageType = "command_ack"
)

// WorkerMessage is the universal envelope for all WebSocket communication
// between the control plane and the worker daemon.
type WorkerMessage struct {
	ID        string            `json:"id"`
	Type      WorkerMessageType `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	// Payload contains the actual data, which will be unmarshaled based on the Type.
	Payload []byte `json:"payload"`
}

// -----------------------------------------------------------------------------
// Authentication Schemas
// -----------------------------------------------------------------------------

type WorkerAuthPayload struct {
	WorkerToken string `json:"worker_token"`
	Version     string `json:"version"`
}

type WorkerAuthResultPayload struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// -----------------------------------------------------------------------------
// Command Schemas (Control Plane -> Worker)
// -----------------------------------------------------------------------------

// WorkerDeployAppPayload contains everything the worker needs to deploy an app.
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

// -----------------------------------------------------------------------------
// Telemetry Schemas (Worker -> Control Plane)
// -----------------------------------------------------------------------------

type WorkerLogStreamPayload struct {
	ContainerID string `json:"container_id"`
	LogLine     string `json:"log_line"`
	StreamType  string `json:"stream_type"` // "stdout" or "stderr"
}

type WorkerMetricsPayload struct {
	CPUUsagePercentage float64 `json:"cpu_usage_percentage"`
	MemoryUsageBytes   uint64  `json:"memory_usage_bytes"`
	MemoryLimitBytes   uint64  `json:"memory_limit_bytes"`
	DiskUsageBytes     uint64  `json:"disk_usage_bytes"`
	DiskTotalBytes     uint64  `json:"disk_total_bytes"`
}
