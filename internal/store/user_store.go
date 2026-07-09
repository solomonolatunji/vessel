package store

import (
	"database/sql"
	"errors"
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

// GetUserByEmail queries a user account by email address.
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow(`SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE email = ?`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByID queries a user account by unique ID.
func (s *Store) GetUserByID(id string) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow(`SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ListUsers retrieves all registered user accounts from SQLite.
func (s *Store) ListUsers() ([]types.User, error) {
	rows, err := s.db.Query(`SELECT id, email, password_hash, role, created_at, updated_at FROM users ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var u types.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// UpdateUser updates an existing user account in SQLite.
func (s *Store) UpdateUser(u *types.User) error {
	u.UpdatedAt = time.Now()
	_, err := s.db.Exec(`UPDATE users SET email = ?, password_hash = ?, role = ?, updated_at = ? WHERE id = ?`,
		u.Email, u.PasswordHash, u.Role, u.UpdatedAt, u.ID)
	return err
}

