package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// initEnvironmentTable creates the environments table if it doesn't already exist.
func (s *Store) initEnvironmentTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS environments (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		name TEXT NOT NULL,
		is_default BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(project_id, name)
	);`
	_, err := s.db.Exec(query)
	return err
}

// CreateEnvironment adds a new environment (e.g., production, staging) to a project workspace canvas.
func (s *Store) CreateEnvironment(env *types.EnvironmentConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if env.ID == "" {
		env.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	env.CreatedAt = now
	env.UpdatedAt = now

	query := `INSERT INTO environments (id, project_id, name, is_default, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, env.ID, env.ProjectID, env.Name, env.IsDefault, env.CreatedAt, env.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}
	return nil
}

// GetEnvironment retrieves an environment configuration by its unique ID.
func (s *Store) GetEnvironment(id string) (*types.EnvironmentConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var env types.EnvironmentConfig
	var isDefault int
	err := row.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("environment not found: %s", id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}
	env.IsDefault = isDefault == 1
	return &env, nil
}

// GetDefaultEnvironment retrieves the default ("production") environment for a project workspace.
func (s *Store) GetDefaultEnvironment(projectID string) (*types.EnvironmentConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = ? AND is_default = 1 LIMIT 1`
	row := s.db.QueryRow(query, projectID)

	var env types.EnvironmentConfig
	var isDefault int
	err := row.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		// Fallback check: if no default flag is set, return the first created environment
		fallbackQuery := `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY created_at ASC LIMIT 1`
		fallbackRow := s.db.QueryRow(fallbackQuery, projectID)
		err = fallbackRow.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no environments found for project: %s", projectID)
		} else if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get default environment: %w", err)
	}
	env.IsDefault = isDefault == 1
	return &env, nil
}

// ListEnvironments returns all environments belonging to a project canvas.
func (s *Store) ListEnvironments(projectID string) ([]*types.EnvironmentConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY is_default DESC, created_at ASC`
	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}
	defer rows.Close()

	var envs []*types.EnvironmentConfig
	for rows.Next() {
		var env types.EnvironmentConfig
		var isDefault int
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan environment row: %w", err)
		}
		env.IsDefault = isDefault == 1
		envs = append(envs, &env)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return envs, nil
}

// DeleteEnvironment deletes an environment and removes associated services if required.
func (s *Store) DeleteEnvironment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM environments WHERE id = ?`, id)
	return err
}
