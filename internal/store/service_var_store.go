package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// initServiceVarsTable initializes the service_vars table.
func (s *Store) initServiceVarsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS service_vars (
		id TEXT PRIMARY KEY,
		service_id TEXT NOT NULL,
		environment_id TEXT NOT NULL,
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		is_secret BOOLEAN DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(service_id, key)
	);`
	_, err := s.db.Exec(query)
	return err
}

// SetServiceVariable sets or updates an environment variable for a specific service.
func (s *Store) SetServiceVariable(v *types.ServiceVariable) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	v.CreatedAt = now
	v.UpdatedAt = now

	query := `INSERT INTO service_vars (id, service_id, environment_id, key, value, is_secret, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(service_id, key) DO UPDATE SET
		value = excluded.value,
		is_secret = excluded.is_secret,
		updated_at = excluded.updated_at`

	_, err := s.db.Exec(query, v.ID, v.ServiceID, v.EnvironmentID, v.Key, v.Value, v.IsSecret, v.CreatedAt, v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to set service variable: %w", err)
	}

	// Update env_vars_count on app_services table
	_, _ = s.db.Exec(`UPDATE app_services SET env_vars_count = (SELECT COUNT(*) FROM service_vars WHERE service_id = ?) WHERE id = ?`, v.ServiceID, v.ServiceID)
	return nil
}

// ListServiceVariables retrieves all variables configured for a specific service (`Variables` tab).
func (s *Store) ListServiceVariables(serviceID string) ([]*types.ServiceVariable, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, service_id, environment_id, key, value, is_secret, created_at, updated_at
		FROM service_vars WHERE service_id = ? ORDER BY key ASC`

	rows, err := s.db.Query(query, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list service variables: %w", err)
	}
	defer rows.Close()

	var vars []*types.ServiceVariable
	for rows.Next() {
		var v types.ServiceVariable
		var isSecret int
		if err := rows.Scan(&v.ID, &v.ServiceID, &v.EnvironmentID, &v.Key, &v.Value, &isSecret, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan service var row: %w", err)
		}
		v.IsSecret = isSecret == 1
		vars = append(vars, &v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return vars, nil
}

// DeleteServiceVariable removes a specific environment variable by ID.
func (s *Store) DeleteServiceVariable(id, serviceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(`DELETE FROM service_vars WHERE id = ? AND service_id = ?`, id, serviceID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("service variable not found")
	}

	_, _ = s.db.Exec(`UPDATE app_services SET env_vars_count = (SELECT COUNT(*) FROM service_vars WHERE service_id = ?) WHERE id = ?`, serviceID, serviceID)
	return nil
}

// BulkSetServiceVariables replaces all variables for a service (supports Raw Editor tab in UI).
func (s *Store) BulkSetServiceVariables(serviceID, environmentID string, vars []*types.ServiceVariable) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM service_vars WHERE service_id = ?`, serviceID); err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, v := range vars {
		if v.ID == "" {
			v.ID = uuid.NewString()
		}
		v.ServiceID = serviceID
		v.EnvironmentID = environmentID
		v.CreatedAt = now
		v.UpdatedAt = now
		_, err := tx.Exec(`INSERT INTO service_vars (id, service_id, environment_id, key, value, is_secret, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			v.ID, v.ServiceID, v.EnvironmentID, v.Key, v.Value, v.IsSecret, v.CreatedAt, v.UpdatedAt)
		if err != nil {
			return err
		}
	}

	if _, err := tx.Exec(`UPDATE app_services SET env_vars_count = ? WHERE id = ?`, len(vars), serviceID); err != nil {
		return err
	}

	return tx.Commit()
}
