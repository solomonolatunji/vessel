package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// initDeploymentsTable initializes the deployments table for the Deployments tab.
func (s *Store) initDeploymentsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS deployments (
		id TEXT PRIMARY KEY,
		service_id TEXT NOT NULL,
		environment_id TEXT NOT NULL,
		project_id TEXT NOT NULL,
		status TEXT NOT NULL,
		commit_hash TEXT DEFAULT '',
		commit_message TEXT DEFAULT '',
		branch TEXT DEFAULT '',
		trigger TEXT DEFAULT 'Manual',
		build_logs TEXT DEFAULT '',
		container_id TEXT DEFAULT '',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		finished_at DATETIME
	);`
	_, err := s.db.Exec(query)
	return err
}

// CreateDeployment records a new deployment attempt.
func (s *Store) CreateDeployment(dep *types.DeploymentRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if dep.ID == "" {
		dep.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	dep.CreatedAt = now
	dep.UpdatedAt = now
	if dep.Status == "" {
		dep.Status = "BUILDING"
	}

	query := `INSERT INTO deployments (
		id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		dep.ID, dep.ServiceID, dep.EnvironmentID, dep.ProjectID, dep.Status, dep.CommitHash,
		dep.CommitMessage, dep.Branch, dep.Trigger, dep.BuildLogs, dep.ContainerID, dep.CreatedAt, dep.UpdatedAt, dep.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create deployment record: %w", err)
	}
	return nil
}

// GetDeployment retrieves a deployment record by ID.
func (s *Store) GetDeployment(id string) (*types.DeploymentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
		FROM deployments WHERE id = ?`

	row := s.db.QueryRow(query, id)
	var dep types.DeploymentRecord
	var finishedAt sql.NullTime
	err := row.Scan(
		&dep.ID, &dep.ServiceID, &dep.EnvironmentID, &dep.ProjectID, &dep.Status, &dep.CommitHash,
		&dep.CommitMessage, &dep.Branch, &dep.Trigger, &dep.BuildLogs, &dep.ContainerID, &dep.CreatedAt, &dep.UpdatedAt, &finishedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("deployment not found: %s", id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to scan deployment: %w", err)
	}
	if finishedAt.Valid {
		dep.FinishedAt = finishedAt.Time
	}
	return &dep, nil
}

// ListDeploymentsByService retrieves all deployment records for a specific service (`Deployments` tab).
func (s *Store) ListDeploymentsByService(serviceID string) ([]*types.DeploymentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
		FROM deployments WHERE service_id = ? ORDER BY created_at DESC`

	rows, err := s.db.Query(query, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query service deployments: %w", err)
	}
	defer rows.Close()

	var deps []*types.DeploymentRecord
	for rows.Next() {
		var dep types.DeploymentRecord
		var finishedAt sql.NullTime
		if err := rows.Scan(
			&dep.ID, &dep.ServiceID, &dep.EnvironmentID, &dep.ProjectID, &dep.Status, &dep.CommitHash,
			&dep.CommitMessage, &dep.Branch, &dep.Trigger, &dep.BuildLogs, &dep.ContainerID, &dep.CreatedAt, &dep.UpdatedAt, &finishedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan deployment row: %w", err)
		}
		if finishedAt.Valid {
			dep.FinishedAt = finishedAt.Time
		}
		deps = append(deps, &dep)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return deps, nil
}

// UpdateDeploymentStatus updates the status, logs, or completion time of a deployment.
func (s *Store) UpdateDeploymentStatus(id, status, buildLogs, containerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	if status == "ACTIVE" || status == "FAILED" || status == "REMOVED" || status == "SLEPT" {
		query := `UPDATE deployments SET status = ?, build_logs = ?, container_id = ?, updated_at = ?, finished_at = ? WHERE id = ?`
		_, err := s.db.Exec(query, status, buildLogs, containerID, now, now, id)
		return err
	}
	query := `UPDATE deployments SET status = ?, build_logs = ?, container_id = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, status, buildLogs, containerID, now, id)
	return err
}
