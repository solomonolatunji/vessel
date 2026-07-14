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

type AppServiceRepository interface {
	Create(ctx context.Context, svc *models.AppService) error
	GetByID(ctx context.Context, id string) (*models.AppService, error)
	ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error)
	ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error)
	ListAll(ctx context.Context) ([]*models.AppService, error)
	Update(ctx context.Context, svc *models.AppService) error
	Delete(ctx context.Context, id string) error
}

type AppServiceSQLiteRepository struct {
	mu sync.RWMutex
	db *sqlx.DB
}

func NewAppServiceSQLiteRepository(db *sql.DB) *AppServiceSQLiteRepository {
	return &AppServiceSQLiteRepository{db: sqlx.NewDb(db, "sqlite")}
}

func (r *AppServiceSQLiteRepository) Create(_ context.Context, svc *models.AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if svc.ID == "" {
		svc.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	svc.CreatedAt = now
	svc.UpdatedAt = now
	if svc.Status == "" {
		svc.Status = "building"
	}
	if svc.InternalPort == 0 {
		svc.InternalPort = 3000
	}
	_, err := r.db.Exec(
		`INSERT INTO app_services (id, project_id, environment_id, name, repository_url, image_ref, branch, root_directory, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, health_check_path, container_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		svc.ID, svc.ProjectID, svc.EnvironmentID, svc.Name, svc.RepositoryURL, svc.ImageRef, svc.Branch,
		svc.RootDirectory, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine,
		svc.InternalPort, svc.Domain, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.CreatedAt, svc.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create app service: %w", err)
	}
	return nil
}

func (r *AppServiceSQLiteRepository) GetByID(ctx context.Context, id string) (*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var svc models.AppService
	err := r.db.GetContext(ctx, &svc,
		`SELECT id, project_id, environment_id, name, repository_url, COALESCE(image_ref,'') AS image_ref, branch, root_directory, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, health_check_path, container_id, status, created_at, updated_at
		FROM app_services WHERE id = ?`, id,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("AppService", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get app service: %w", err)
	}
	return &svc, nil
}

func (r *AppServiceSQLiteRepository) ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT id, project_id, environment_id, name, repository_url, COALESCE(image_ref,'') AS image_ref, branch, root_directory, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, health_check_path, container_id, status, created_at, updated_at
		FROM app_services WHERE environment_id = ? ORDER BY created_at ASC`, environmentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by environment: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceSQLiteRepository) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT id, project_id, environment_id, name, repository_url, COALESCE(image_ref,'') AS image_ref, branch, root_directory, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, health_check_path, container_id, status, created_at, updated_at
		FROM app_services WHERE project_id = ? ORDER BY created_at ASC`, projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by project: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceSQLiteRepository) ListAll(ctx context.Context) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT id, project_id, environment_id, name, repository_url, COALESCE(image_ref,'') AS image_ref, branch, root_directory, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, health_check_path, container_id, status, created_at, updated_at
		FROM app_services ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list all app services: %w", err)
	}
	if list == nil {
		list = make([]*models.AppService, 0)
	}
	return list, nil
}

func (r *AppServiceSQLiteRepository) Update(_ context.Context, svc *models.AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc.UpdatedAt = time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE app_services SET
			name = ?, repository_url = ?, image_ref = ?, branch = ?, root_directory = ?, build_command = ?, start_command = ?, dockerfile_path = ?, build_engine = ?, internal_port = ?, domain = ?, health_check_path = ?, container_id = ?, status = ?, updated_at = ?
		WHERE id = ?`,
		svc.Name, svc.RepositoryURL, svc.ImageRef, svc.Branch, svc.RootDirectory, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine, svc.InternalPort, svc.Domain, svc.HealthCheckPath,
		svc.ContainerID, svc.Status, svc.UpdatedAt, svc.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update app service: %w", err)
	}
	return nil
}

func (r *AppServiceSQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM app_services WHERE id = ?`, id)
	return err
}
