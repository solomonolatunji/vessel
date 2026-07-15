package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type DeploymentRepository interface {
	Create(ctx context.Context, d *models.Deployment) error
	GetByID(ctx context.Context, id string) (*models.Deployment, error)
	ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*models.Deployment, int, error)
	Update(ctx context.Context, d *models.Deployment) error
	UpdateStatus(ctx context.Context, id string, status models.DeploymentStatus, buildLogs, containerID string) error
}

type DeploymentSQLiteRepository struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewDeploymentSQLiteRepository(db *sql.DB) *DeploymentSQLiteRepository {
	return &DeploymentSQLiteRepository{db: sqlx.NewDb(db, "sqlite")}
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
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ServiceID, d.EnvironmentID, d.ProjectID, d.Status, d.CommitHash,
		d.CommitMessage, d.Branch, d.Trigger, d.BuildLogs, d.ContainerID, d.CreatedAt, d.UpdatedAt, d.FinishedAt)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}
	return nil
}

func (r *DeploymentSQLiteRepository) GetByID(ctx context.Context, id string) (*models.Deployment, error) {
	var d models.Deployment
	err := r.db.GetContext(ctx, &d, `SELECT id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
		FROM deployments WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Deployment", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan deployment: %w", err)
	}
	return &d, nil
}

func (r *DeploymentSQLiteRepository) ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*models.Deployment, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM deployments WHERE service_id = ?`, serviceID); err != nil {
		return nil, 0, err
	}

	var deps []*models.Deployment
	err := r.db.SelectContext(ctx, &deps, `SELECT id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
		FROM deployments WHERE service_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, serviceID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query service deployments: %w", err)
	}
	if deps == nil {
		deps = make([]*models.Deployment, 0)
	}
	return deps, total, nil
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

func (r *DeploymentSQLiteRepository) UpdateStatus(_ context.Context, id string, status models.DeploymentStatus, buildLogs, containerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	if status == models.DeploymentStatusActive || status == models.DeploymentStatusFailed || status == models.DeploymentStatusRemoved || status == models.DeploymentStatusSlept {
		_, err := r.db.Exec(`UPDATE deployments SET status = ?, build_logs = ?, container_id = ?, updated_at = ?, finished_at = ? WHERE id = ?`,
			status, buildLogs, containerID, now, now, id)
		return err
	}
	_, err := r.db.Exec(`UPDATE deployments SET status = ?, build_logs = ?, container_id = ?, updated_at = ? WHERE id = ?`,
		status, buildLogs, containerID, now, id)
	return err
}
