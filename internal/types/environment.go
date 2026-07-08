package types

import "time"

// EnvironmentConfig represents an isolated runtime environment (e.g., production, staging, preview) under a Project canvas.
type EnvironmentConfig struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Name      string    `json:"name"`      // e.g., "production", "staging"
	IsDefault bool      `json:"isDefault"` // true for initial "production" environment
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
