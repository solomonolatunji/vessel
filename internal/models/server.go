package models

import "time"

// ServerStatus represents the connection health of a worker node.
type ServerStatus string

const (
	ServerStatusOnline       ServerStatus = "online"
	ServerStatusOffline      ServerStatus = "offline"
	ServerStatusProvisioning ServerStatus = "provisioning"
)

// Server represents a physical or virtual machine running the codedock-worker daemon.
type Server struct {
	ID          string       `json:"id" db:"id"`
	UserID      string       `json:"userId" db:"user_id"` // Owner of the server (for multi-tenancy)
	Name        string       `json:"name" db:"name"`
	IPAddress   string       `json:"ipAddress" db:"ip_address"`
	Status      ServerStatus `json:"status" db:"status"`
	WorkerToken string       `json:"workerToken" db:"worker_token"`
	LastSeenAt  *time.Time   `json:"lastSeenAt" db:"last_seen_at"`

	// Metrics contains the latest JSON payload from the worker (CPU, RAM, Disk).
	// This is stored as a raw JSON string in the database and unmarshaled on demand.
	Metrics []byte `json:"metrics,omitempty" db:"metrics"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
