package types

import "time"

// Team represents an organization or collaboration workspace.
type Team struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TeamMember represents a user's membership and role within a team (`Owner`, `Admin`, or `Member`).
type TeamMember struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	UserID    string    `json:"userId"`
	UserEmail string    `json:"userEmail"`
	Role      string    `json:"role"` // Owner, Admin, Member
	JoinedAt  time.Time `json:"joinedAt"`
}

// TeamInvite represents a pending invitation to join a team.
type TeamInvite struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Token     string    `json:"token"`
	InvitedBy string    `json:"invitedBy"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}
