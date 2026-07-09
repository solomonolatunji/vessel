package types

import "time"

// WorkspaceTrustedDomain represents a trusted domain or SSO domain for a team workspace.
type WorkspaceTrustedDomain struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Domain    string    `json:"domain"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// WorkspaceSSHKey represents an SSH public key registered to a team workspace for git clone/deployment authentication.
type WorkspaceSSHKey struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Name      string    `json:"name"`
	PublicKey string    `json:"publicKey"`
	CreatedAt time.Time `json:"createdAt"`
}

// WorkspaceAuditLog records security and lifecycle events within a workspace or project.
type WorkspaceAuditLog struct {
	ID            string    `json:"id"`
	TeamID        string    `json:"teamId"`
	ProjectID     string    `json:"projectId,omitempty"`
	EnvironmentID string    `json:"environmentId,omitempty"`
	Action        string    `json:"action"`
	Actor         string    `json:"actor"`
	CreatedAt     time.Time `json:"createdAt"`
}
