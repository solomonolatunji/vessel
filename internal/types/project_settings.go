package types

import "time"

// ServiceVariable represents an environment variable scoped specifically to a Service inside an Environment.
type ServiceVariable struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"serviceId"`
	EnvironmentID string    `json:"environmentId"`
	Key           string    `json:"key"`
	Value         string    `json:"value"` // Decrypted plaintext when returned to API if authorized
	IsSecret      bool      `json:"isSecret"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ProjectWebhook represents a webhook configured in Project Settings -> Webhooks tab.
type ProjectWebhook struct {
	ID                    string    `json:"id"`
	ProjectID             string    `json:"projectId"`
	URL                   string    `json:"url"`
	EventTypes            []string  `json:"eventTypes"` // e.g. ["Deployment Deployed", "Deployment Failed", "VolumeAlert Triggered"]
	IncludePREnvironments bool      `json:"includePrEnvironments"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// ProjectToken represents an API token created inside Project Settings -> Tokens tab.
type ProjectToken struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	TokenPrefix   string    `json:"tokenPrefix"` // First 8 chars for display
	CreatedAt     time.Time `json:"createdAt"`
}

// ProjectMember represents an invited or active collaborator inside Project Settings -> Members tab.
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
