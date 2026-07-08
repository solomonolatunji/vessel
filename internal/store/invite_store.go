package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateInvite issues a new workspace invitation token with a 7-day expiration.
func (s *Store) CreateInvite(inv *types.Invite) error {
	if inv.ID == "" {
		inv.ID = uuid.NewString()
	}
	if inv.Token == "" {
		inv.Token = uuid.NewString()
	}
	inv.CreatedAt = time.Now()
	if inv.ExpiresAt.IsZero() {
		inv.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	}

	_, err := s.db.Exec(`INSERT INTO invites (id, email, role, token, invited_by, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, inv.ID, inv.Email, inv.Role, inv.Token, inv.InvitedBy, inv.ExpiresAt, inv.CreatedAt)
	return err
}
