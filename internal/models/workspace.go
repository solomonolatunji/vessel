package models

import "time"

type Workspace struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	AvatarURL       string    `json:"avatarUrl,omitempty" db:"avatar_url"`
	PreferredRegion string    `json:"preferredRegion,omitempty" db:"preferred_region"`
	OwnerID         string    `json:"ownerId" db:"owner_id"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type WorkspaceMember struct {
	ID          string    `json:"id" db:"id"`
	WorkspaceID string    `json:"workspaceId" db:"workspace_id"`
	UserID      string    `json:"userId" db:"user_id"`
	UserEmail   string    `json:"userEmail" db:"user_email"`
	Role        string    `json:"role" db:"role"`
	JoinedAt    time.Time `json:"joinedAt" db:"joined_at"`
}

type WorkspaceInvite struct {
	ID          string    `json:"id" db:"id"`
	WorkspaceID string    `json:"workspaceId" db:"workspace_id"`
	Email       string    `json:"email" db:"email"`
	Role        string    `json:"role" db:"role"`
	Token       string    `json:"token" db:"token"`
	InvitedBy   string    `json:"invitedBy" db:"invited_by"`
	ExpiresAt   time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type CreateWorkspaceRequest struct {
	Name            string `json:"name"`
	AvatarURL       string `json:"avatarUrl,omitempty"`
	PreferredRegion string `json:"preferredRegion,omitempty"`
}

type InviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AcceptInviteRequest struct {
	Token string `json:"token"`
}

type GetWorkspaceResponse struct {
	Workspace *Workspace         `json:"workspace"`
	Members   []*WorkspaceMember `json:"members"`
}

type TrustedDomain struct {
	ID          string    `json:"id" db:"id"`
	WorkspaceID string    `json:"workspaceId" db:"workspace_id"`
	Domain      string    `json:"domain" db:"domain"`
	Role        string    `json:"role" db:"role"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type SSHKey struct {
	ID          string    `json:"id" db:"id"`
	WorkspaceID string    `json:"workspaceId" db:"workspace_id"`
	Name        string    `json:"name" db:"name"`
	PublicKey   string    `json:"publicKey" db:"public_key"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type AuditLog struct {
	ID            string    `json:"id" gorm:"primarykey" db:"id"`
	WorkspaceID   string    `json:"workspaceId" db:"workspace_id"`
	UserID        string    `json:"userId,omitempty" db:"user_id"`
	ProjectID     string    `json:"projectId,omitempty" db:"project_id"`
	EnvironmentID string    `json:"environmentId,omitempty" db:"environment_id"`
	Action        string    `json:"action" db:"action"`
	Resource      string    `json:"resource,omitempty" db:"resource"`
	Actor         string    `json:"actor" db:"actor"`
	IPAddress     string    `json:"ipAddress,omitempty" db:"ip_address"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	Timestamp     time.Time `json:"timestamp,omitempty" db:"timestamp"`
}

type UpdateWorkspaceRequest struct {
	Name            string `json:"name,omitempty"`
	AvatarURL       string `json:"avatarUrl,omitempty"`
	PreferredRegion string `json:"preferredRegion,omitempty"`
}

type CreateTrustedDomainRequest struct {
	Domain string `json:"domain"`
	Role   string `json:"role"`
}

type CreateSSHKeyRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}
