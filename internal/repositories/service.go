package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/utils"
)

type AppServiceRepository interface {
	Create(ctx context.Context, svc *models.AppService) error
	GetByID(ctx context.Context, id string) (*models.AppService, error)
	ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error)
	ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error)
	ListAll(ctx context.Context) ([]*models.AppService, error)
	Update(ctx context.Context, svc *models.AppService) error
	Delete(ctx context.Context, id string) error
	CreateWebhook(ctx context.Context, w *models.Webhook) error
	ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error)
	DeleteWebhook(ctx context.Context, id, serviceID string) error
	CreateLogDrain(ctx context.Context, d *models.LogDrain) error
	ListLogDrainsByService(ctx context.Context, serviceID string) ([]*models.LogDrain, error)
	DeleteLogDrain(ctx context.Context, id, serviceID string) error
}

type AppServiceRepo struct {
	mu sync.RWMutex
	db *sqlx.DB
}

func NewAppServiceRepo(db *sql.DB) *AppServiceRepo {
	return &AppServiceRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *AppServiceRepo) Create(_ context.Context, svc *models.AppService) error {
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
		`INSERT INTO app_services (id, project_id, environment_id, name, repository_url, image_ref, branch, root_directory, icon, runtime_mode, install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, static_output, health_check_path, container_id, status, replicas, cpu_limit, memory_limit, deploy_token, enable_pr_previews, maintenance_mode, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		svc.ID, svc.ProjectID, svc.EnvironmentID, svc.Name, svc.RepositoryURL, svc.ImageRef, svc.Branch,
		svc.RootDirectory, svc.Icon, svc.RuntimeMode, svc.InstallCommand, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine,
		svc.InternalPort, svc.Domain, svc.StaticOutput, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.Replicas, svc.CPULimit, svc.MemoryLimit, svc.DeployToken, svc.EnablePRPreviews, svc.MaintenanceMode, svc.CreatedAt, svc.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create app service: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) GetByID(ctx context.Context, id string) (*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var svc models.AppService
	err := r.db.GetContext(ctx, &svc,
		`SELECT id, project_id, environment_id, name, repository_url, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, 0) AS enable_pr_previews, COALESCE(maintenance_mode, 0) AS maintenance_mode, created_at, updated_at
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

func (r *AppServiceRepo) ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT id, project_id, environment_id, name, repository_url, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, 0) AS enable_pr_previews, COALESCE(maintenance_mode, 0) AS maintenance_mode, created_at, updated_at
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

func (r *AppServiceRepo) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT id, project_id, environment_id, name, repository_url, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, 0) AS enable_pr_previews, COALESCE(maintenance_mode, 0) AS maintenance_mode, created_at, updated_at
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

func (r *AppServiceRepo) ListAll(ctx context.Context) ([]*models.AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*models.AppService
	err := r.db.SelectContext(ctx, &list,
		`SELECT id, project_id, environment_id, name, repository_url, COALESCE(image_ref,'') AS image_ref, branch, root_directory, COALESCE(icon,'git') AS icon, runtime_mode, COALESCE(install_command,'') AS install_command, build_command, start_command, dockerfile_path, build_engine, internal_port, domain, COALESCE(static_output,'') AS static_output, health_check_path, container_id, status, replicas, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, COALESCE(deploy_token,'') AS deploy_token, COALESCE(enable_pr_previews, 0) AS enable_pr_previews, COALESCE(maintenance_mode, 0) AS maintenance_mode, created_at, updated_at
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

func (r *AppServiceRepo) Update(_ context.Context, svc *models.AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc.UpdatedAt = time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE app_services SET
		name = ?, repository_url = ?, image_ref = ?, branch = ?, root_directory = ?, icon = ?, runtime_mode = ?,
		install_command = ?, build_command = ?, start_command = ?, dockerfile_path = ?, build_engine = ?,
		internal_port = ?, domain = ?, static_output = ?, health_check_path = ?, container_id = ?, status = ?, replicas = ?, cpu_limit = ?, memory_limit = ?, deploy_token = ?, enable_pr_previews = ?, maintenance_mode = ?, updated_at = ?
		WHERE id = ?`,
		svc.Name, svc.RepositoryURL, svc.ImageRef, svc.Branch, svc.RootDirectory, svc.Icon, svc.RuntimeMode,
		svc.InstallCommand, svc.BuildCommand, svc.StartCommand, svc.DockerfilePath, svc.BuildEngine,
		svc.InternalPort, svc.Domain, svc.StaticOutput, svc.HealthCheckPath, svc.ContainerID, svc.Status, svc.Replicas, svc.CPULimit, svc.MemoryLimit, svc.DeployToken, svc.EnablePRPreviews, svc.MaintenanceMode, svc.UpdatedAt, svc.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update app service: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM app_services WHERE id = ?`, id)
	return err
}

func (r *AppServiceRepo) CreateWebhook(ctx context.Context, w *models.Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	eventTypesStr := strings.Join(w.EventTypes, ",")
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO service_webhooks (id, service_id, url, event_types, include_pr_environments, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		w.ID, w.ServiceID, w.URL, eventTypesStr, w.IncludePREnvironments, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, service_id, url, event_types, include_pr_environments, created_at, updated_at
		 FROM service_webhooks WHERE service_id = ? ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list webhooks: %w", err)
	}
	defer rows.Close()
	var out []*models.Webhook
	for rows.Next() {
		var w models.Webhook
		var eventsStr string
		var includePr int
		if err := rows.Scan(&w.ID, &w.ServiceID, &w.URL, &eventsStr, &includePr, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan webhook: %w", err)
		}
		if eventsStr != "" {
			w.EventTypes = strings.Split(eventsStr, ",")
		} else {
			w.EventTypes = []string{}
		}
		w.IncludePREnvironments = includePr == 1
		out = append(out, &w)
	}
	return out, rows.Err()
}

func (r *AppServiceRepo) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := "DELETE FROM service_webhooks WHERE id = ? AND service_id = ?"
	res, err := r.db.ExecContext(ctx, query, id, serviceID)
	if err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return utils.NewNotFoundError("Webhook", id)
	}
	return nil
}

func (r *AppServiceRepo) CreateLogDrain(ctx context.Context, d *models.LogDrain) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO log_drains (id, service_id, project_id, drain_type, endpoint_url, auth_token, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ServiceID, d.ProjectID, d.DrainType, d.EndpointURL, d.AuthToken, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create log drain: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) ListLogDrainsByService(ctx context.Context, serviceID string) ([]*models.LogDrain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, service_id, project_id, drain_type, endpoint_url, auth_token, created_at, updated_at
		 FROM log_drains WHERE service_id = ? ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list log drains: %w", err)
	}
	defer rows.Close()
	var out []*models.LogDrain
	for rows.Next() {
		var d models.LogDrain
		if err := rows.Scan(&d.ID, &d.ServiceID, &d.ProjectID, &d.DrainType, &d.EndpointURL, &d.AuthToken, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan log drain: %w", err)
		}
		out = append(out, &d)
	}
	return out, rows.Err()
}

func (r *AppServiceRepo) DeleteLogDrain(ctx context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := "DELETE FROM log_drains WHERE id = ? AND service_id = ?"
	res, err := r.db.ExecContext(ctx, query, id, serviceID)
	if err != nil {
		return fmt.Errorf("delete log drain: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return utils.NewNotFoundError("LogDrain", id)
	}
	return nil
}
