package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateWorkspaceTrustedDomain inserts a new trusted domain entry for a workspace team.
func (s *Store) CreateWorkspaceTrustedDomain(item *types.WorkspaceTrustedDomain) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}
	if item.Role == "" {
		item.Role = "developer"
	}

	query := `INSERT INTO workspace_trusted_domains (id, team_id, domain, role, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, item.ID, item.TeamID, item.Domain, item.Role, item.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create workspace trusted domain: %w", err)
	}
	return nil
}

// ListWorkspaceTrustedDomains returns all trusted domains for a given team workspace.
func (s *Store) ListWorkspaceTrustedDomains(teamID string) ([]*types.WorkspaceTrustedDomain, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, team_id, domain, role, created_at FROM workspace_trusted_domains WHERE team_id = ? ORDER BY created_at DESC`
	rows, err := s.db.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace trusted domains: %w", err)
	}
	defer rows.Close()

	var list []*types.WorkspaceTrustedDomain
	for rows.Next() {
		var item types.WorkspaceTrustedDomain
		var createdAtStr string
		if err := rows.Scan(&item.ID, &item.TeamID, &item.Domain, &item.Role, &createdAtStr); err != nil {
			return nil, err
		}
		item.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &item)
	}
	return list, nil
}

// DeleteWorkspaceTrustedDomain removes a trusted domain by ID.
func (s *Store) DeleteWorkspaceTrustedDomain(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM workspace_trusted_domains WHERE id = ?`, id)
	return err
}

// CreateWorkspaceSSHKey inserts a new SSH public key entry for a workspace team.
func (s *Store) CreateWorkspaceSSHKey(item *types.WorkspaceSSHKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}

	query := `INSERT INTO workspace_ssh_keys (id, team_id, name, public_key, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, item.ID, item.TeamID, item.Name, item.PublicKey, item.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create workspace ssh key: %w", err)
	}
	return nil
}

// ListWorkspaceSSHKeys returns all SSH public keys for a given team workspace.
func (s *Store) ListWorkspaceSSHKeys(teamID string) ([]*types.WorkspaceSSHKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, team_id, name, public_key, created_at FROM workspace_ssh_keys WHERE team_id = ? ORDER BY created_at DESC`
	rows, err := s.db.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace ssh keys: %w", err)
	}
	defer rows.Close()

	var list []*types.WorkspaceSSHKey
	for rows.Next() {
		var item types.WorkspaceSSHKey
		var createdAtStr string
		if err := rows.Scan(&item.ID, &item.TeamID, &item.Name, &item.PublicKey, &createdAtStr); err != nil {
			return nil, err
		}
		item.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &item)
	}
	return list, nil
}

// DeleteWorkspaceSSHKey removes an SSH key by ID.
func (s *Store) DeleteWorkspaceSSHKey(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM workspace_ssh_keys WHERE id = ?`, id)
	return err
}

// CreateWorkspaceAuditLog inserts an audit log event into the workspace history.
func (s *Store) CreateWorkspaceAuditLog(log *types.WorkspaceAuditLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}

	query := `INSERT INTO workspace_audit_logs (id, team_id, project_id, environment_id, action, actor, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, log.ID, log.TeamID, log.ProjectID, log.EnvironmentID, log.Action, log.Actor, log.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create workspace audit log: %w", err)
	}
	return nil
}

// ListWorkspaceAuditLogs returns the audit history for a given team workspace.
func (s *Store) ListWorkspaceAuditLogs(teamID string, limit int) ([]*types.WorkspaceAuditLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}

	query := `SELECT id, team_id, COALESCE(project_id, ''), COALESCE(environment_id, ''), action, actor, created_at
		FROM workspace_audit_logs WHERE team_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := s.db.Query(query, teamID, limit)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to query workspace audit logs: %w", err)
	}
	if rows == nil {
		return []*types.WorkspaceAuditLog{}, nil
	}
	defer rows.Close()

	var list []*types.WorkspaceAuditLog
	for rows.Next() {
		var log types.WorkspaceAuditLog
		var createdAtStr string
		if err := rows.Scan(&log.ID, &log.TeamID, &log.ProjectID, &log.EnvironmentID, &log.Action, &log.Actor, &createdAtStr); err != nil {
			return nil, err
		}
		log.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &log)
	}
	return list, nil
}

// CreateWorkspace inserts a new top-level Workspace into the database.
func (s *Store) CreateWorkspace(item *types.Workspace) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = time.Now().UTC()
	}
	if item.PreferredRegion == "" {
		item.PreferredRegion = "local"
	}

	query := `INSERT INTO workspaces (id, name, avatar_url, preferred_region, owner_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, item.ID, item.Name, item.AvatarURL, item.PreferredRegion, item.OwnerID, item.CreatedAt.Format(time.RFC3339), item.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	return nil
}

// GetWorkspace retrieves a Workspace by ID.
func (s *Store) GetWorkspace(id string) (*types.Workspace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var item types.Workspace
	var createdAtStr, updatedAtStr string
	if err := row.Scan(&item.ID, &item.Name, &item.AvatarURL, &item.PreferredRegion, &item.OwnerID, &createdAtStr, &updatedAtStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}
	item.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	item.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
	return &item, nil
}

// ListWorkspaces lists all workspaces owned by or accessible to ownerID.
func (s *Store) ListWorkspaces(ownerID string) ([]*types.Workspace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE owner_id = ? ORDER BY created_at DESC`
	rows, err := s.db.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", err)
	}
	defer rows.Close()

	var list []*types.Workspace
	for rows.Next() {
		var item types.Workspace
		var createdAtStr, updatedAtStr string
		if err := rows.Scan(&item.ID, &item.Name, &item.AvatarURL, &item.PreferredRegion, &item.OwnerID, &createdAtStr, &updatedAtStr); err != nil {
			return nil, err
		}
		item.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		item.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		list = append(list, &item)
	}
	return list, nil
}

// UpdateWorkspace updates workspace details.
func (s *Store) UpdateWorkspace(item *types.Workspace) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item.UpdatedAt = time.Now().UTC()
	query := `UPDATE workspaces SET name = ?, avatar_url = ?, preferred_region = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, item.Name, item.AvatarURL, item.PreferredRegion, item.UpdatedAt.Format(time.RFC3339), item.ID)
	if err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}
	return nil
}

// DeleteWorkspace removes a workspace if it is not the last workspace owned by the user.
func (s *Store) DeleteWorkspace(id string, ownerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	_ = s.db.QueryRow("SELECT count(*) FROM workspaces WHERE owner_id = ?", ownerID).Scan(&count)
	if count <= 1 {
		return errors.New("You cannot delete your last workspace. To delete your account, visit Account Settings")
	}

	_, err := s.db.Exec("DELETE FROM workspaces WHERE id = ? AND owner_id = ?", id, ownerID)
	return err
}

