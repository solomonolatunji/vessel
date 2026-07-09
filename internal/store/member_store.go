package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// initProjectMembersTable initializes the project_members table.
func (s *Store) initProjectMembersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS project_members (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		user_id TEXT DEFAULT '',
		email TEXT NOT NULL,
		permission TEXT NOT NULL,
		status TEXT NOT NULL,
		invited_at DATETIME NOT NULL,
		accepted_at DATETIME,
		UNIQUE(project_id, email)
	);`
	_, err := s.db.Exec(query)
	return err
}

// CreateOrInviteProjectMember adds a member or sends an invite (`Project Settings` -> `Members`).
func (s *Store) CreateOrInviteProjectMember(m *types.ProjectMember) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	m.InvitedAt = now
	if m.Status == "" {
		m.Status = "pending"
	}
	if m.Permission == "" {
		m.Permission = "Can Edit"
	}

	query := `INSERT INTO project_members (id, project_id, user_id, email, permission, status, invited_at, accepted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(project_id, email) DO UPDATE SET
		permission = excluded.permission,
		status = excluded.status`

	_, err := s.db.Exec(query, m.ID, m.ProjectID, m.UserID, m.Email, m.Permission, m.Status, m.InvitedAt, m.AcceptedAt)
	if err != nil {
		return fmt.Errorf("failed to invite/add project member: %w", err)
	}
	return nil
}

// ListProjectMembers lists all members and invitations for a project (`Project Settings` -> `Members`).
func (s *Store) ListProjectMembers(projectID string) ([]*types.ProjectMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, user_id, email, permission, status, invited_at, accepted_at
		FROM project_members WHERE project_id = ? ORDER BY invited_at ASC`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list project members: %w", err)
	}
	defer rows.Close()

	var members []*types.ProjectMember
	for rows.Next() {
		var m types.ProjectMember
		var acceptedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Email, &m.Permission, &m.Status, &m.InvitedAt, &acceptedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project member row: %w", err)
		}
		if acceptedAt.Valid {
			m.AcceptedAt = acceptedAt.Time
		}
		members = append(members, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

// RemoveProjectMember removes a collaborator or cancels an invitation.
func (s *Store) RemoveProjectMember(id, projectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(`DELETE FROM project_members WHERE id = ? AND project_id = ?`, id, projectID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("project member not found")
	}
	return nil
}
