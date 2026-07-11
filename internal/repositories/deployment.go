package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

type DeploymentRepository interface {
	Create(ctx context.Context, d *models.Deployment) error
	GetByID(ctx context.Context, id string) (*models.Deployment, error)
	ListByService(ctx context.Context, serviceID string) ([]*models.Deployment, error)
	Update(ctx context.Context, d *models.Deployment) error
	UpdateStatus(ctx context.Context, id, status, buildLogs, containerID string) error
}

type DeploymentSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewDeploymentSQLiteRepository(db *sql.DB) *DeploymentSQLiteRepository {
	return &DeploymentSQLiteRepository{db: db}
}

func (r *DeploymentSQLiteRepository) Create(_ context.Context, d *models.Deployment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.Status == "" {
		d.Status = "BUILDING"
	}
	_, err := r.db.Exec(`INSERT INTO deployments (
		id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ServiceID, d.EnvironmentID, d.ProjectID, d.Status, d.CommitHash,
		d.CommitMessage, d.Branch, d.Trigger, d.BuildLogs, d.ContainerID, d.CreatedAt, d.UpdatedAt, d.FinishedAt)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}
	return nil
}

func (r *DeploymentSQLiteRepository) GetByID(_ context.Context, id string) (*models.Deployment, error) {
	row := r.db.QueryRow(`SELECT id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
		FROM deployments WHERE id = ?`, id)
	var d models.Deployment
	var finishedAt sql.NullTime
	err := row.Scan(
		&d.ID, &d.ServiceID, &d.EnvironmentID, &d.ProjectID, &d.Status, &d.CommitHash,
		&d.CommitMessage, &d.Branch, &d.Trigger, &d.BuildLogs, &d.ContainerID, &d.CreatedAt, &d.UpdatedAt, &finishedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("deployment not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan deployment: %w", err)
	}
	if finishedAt.Valid {
		d.FinishedAt = finishedAt.Time
	}
	return &d, nil
}

func (r *DeploymentSQLiteRepository) ListByService(_ context.Context, serviceID string) ([]*models.Deployment, error) {
	rows, err := r.db.Query(`SELECT id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
		FROM deployments WHERE service_id = ? ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query service deployments: %w", err)
	}
	defer rows.Close()
	var deps []*models.Deployment
	for rows.Next() {
		var d models.Deployment
		var finishedAt sql.NullTime
		if err := rows.Scan(
			&d.ID, &d.ServiceID, &d.EnvironmentID, &d.ProjectID, &d.Status, &d.CommitHash,
			&d.CommitMessage, &d.Branch, &d.Trigger, &d.BuildLogs, &d.ContainerID, &d.CreatedAt, &d.UpdatedAt, &finishedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan deployment row: %w", err)
		}
		if finishedAt.Valid {
			d.FinishedAt = finishedAt.Time
		}
		deps = append(deps, &d)
	}
	return deps, rows.Err()
}

func (r *DeploymentSQLiteRepository) Update(_ context.Context, d *models.Deployment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d.UpdatedAt = time.Now().UTC()
	_, err := r.db.Exec(`UPDATE deployments SET status = ?, commit_hash = ?, commit_message = ?,
		branch = ?, trigger = ?, build_logs = ?, container_id = ?, updated_at = ?, finished_at = ? WHERE id = ?`,
		d.Status, d.CommitHash, d.CommitMessage, d.Branch, d.Trigger, d.BuildLogs, d.ContainerID, d.UpdatedAt, d.FinishedAt, d.ID)
	return err
}

func (r *DeploymentSQLiteRepository) UpdateStatus(_ context.Context, id, status, buildLogs, containerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	if status == "ACTIVE" || status == "FAILED" || status == "REMOVED" || status == "SLEPT" {
		_, err := r.db.Exec(`UPDATE deployments SET status = ?, build_logs = ?, container_id = ?, updated_at = ?, finished_at = ? WHERE id = ?`,
			status, buildLogs, containerID, now, now, id)
		return err
	}
	_, err := r.db.Exec(`UPDATE deployments SET status = ?, build_logs = ?, container_id = ?, updated_at = ? WHERE id = ?`,
		status, buildLogs, containerID, now, id)
	return err
}
