package models

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

type WorkspaceMember struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	UserID      string    `json:"userId"`
	UserEmail   string    `json:"userEmail"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joinedAt"`
}

type WorkspaceInvite struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Token       string    `json:"token"`
	InvitedBy   string    `json:"invitedBy"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CreatedAt   time.Time `json:"createdAt"`
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
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Domain      string    `json:"domain"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SSHKey struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Name        string    `json:"name"`
	PublicKey   string    `json:"publicKey"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AuditLog struct {
	ID            string    `json:"id" gorm:"primarykey"`
	WorkspaceID   string    `json:"workspaceId"`
	UserID        string    `json:"userId,omitempty"`
	ProjectID     string    `json:"projectId,omitempty"`
	EnvironmentID string    `json:"environmentId,omitempty"`
	Action        string    `json:"action"`
	Resource      string    `json:"resource,omitempty"`
	Actor         string    `json:"actor"`
	IPAddress     string    `json:"ipAddress,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	Timestamp     time.Time `json:"timestamp,omitempty"`
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
