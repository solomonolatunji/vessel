package types

import "time"

type Workspace struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	AvatarURL       string    `json:"avatarUrl,omitempty"`
	PreferredRegion string    `json:"preferredRegion,omitempty"`
	OwnerID         string    `json:"ownerId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type WorkspaceTrustedDomain struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Domain    string    `json:"domain"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type WorkspaceSSHKey struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Name      string    `json:"name"`
	PublicKey string    `json:"publicKey"`
	CreatedAt time.Time `json:"createdAt"`
}

type WorkspaceAuditLog struct {
	ID            string    `json:"id"`
	TeamID        string    `json:"teamId"`
	ProjectID     string    `json:"projectId,omitempty"`
	EnvironmentID string    `json:"environmentId,omitempty"`
	Action        string    `json:"action"`
	Actor         string    `json:"actor"`
	CreatedAt     time.Time `json:"createdAt"`
}
