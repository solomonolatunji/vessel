package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateUser registers a new authenticated user in SQLite.
func (s *Store) CreateUser(u *types.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err := s.db.Exec(`INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`, u.ID, u.Email, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt)
	return err
}
