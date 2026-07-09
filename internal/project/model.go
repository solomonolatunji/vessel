package project

import (
	"time"

	"vessel.dev/vessel/internal/domain"
)

// ProjectConfig holds the core metadata for a Vessel project.
type ProjectConfig struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId,omitempty"`
	TeamID      string    `json:"teamId,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// DomainConfig is an alias for domain.Config kept for internal convenience.
type DomainConfig = domain.Config
