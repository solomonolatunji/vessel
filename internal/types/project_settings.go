package types

import "time"

type EnvVar struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"projectId"`
	Key            string    `json:"key"`
	EncryptedValue string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ProjectWebhook struct {
	ID                    string    `json:"id"`
	ProjectID             string    `json:"projectId"`
	URL                   string    `json:"url"`
	EventTypes            []string  `json:"eventTypes"` // e.g. ["Deployment Deployed", "Deployment Failed", "VolumeAlert Triggered"]
	IncludePREnvironments bool      `json:"includePrEnvironments"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type ProjectToken struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	TokenPrefix   string    `json:"tokenPrefix"` // First 8 chars for display
	CreatedAt     time.Time `json:"createdAt"`
}

type ProjectMember struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"projectId"`
	UserID     string    `json:"userId,omitempty"`
	Email      string    `json:"email"`
	Permission string    `json:"permission"` // e.g. "Can Edit", "Can View", "Admin"
	Status     string    `json:"status"`     // "active", "pending"
	InvitedAt  time.Time `json:"invitedAt"`
	AcceptedAt time.Time `json:"acceptedAt,omitempty"`
}
