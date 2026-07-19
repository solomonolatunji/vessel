package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

func (r *CanvasRepo) GetEnvironmentCanvas(ctx context.Context, environmentID string) (*models.EnvironmentCanvas, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var env models.EnvironmentConfig
	err := r.db.GetContext(ctx, &env, `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = ?`, environmentID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("environment not found: %s", environmentID)
	}
	if err != nil {
		return nil, err
	}
	apps, _ := r.listAppServicesByEnvironment(ctx, environmentID)
	dbs, _ := r.listDatabasesByEnvironment(ctx, environmentID)
	var dbsPtrs []*models.Database
	for i := range dbs {
		dbsPtrs = append(dbsPtrs, &dbs[i])
	}
	return &models.EnvironmentCanvas{
		Environment: &env,
		Apps:        apps,
		Databases:   dbsPtrs,
	}, nil
}

type projectRow struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r *CanvasRepo) listAllProjects() ([]projectRow, error) {
	var projects []projectRow
	err := r.db.Select(&projects, `SELECT id, name, COALESCE(description,'') as description, created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	if projects == nil {
		projects = make([]projectRow, 0)
	}
	return projects, nil
}

func (r *CanvasRepo) getProject(id string) (*projectRow, error) {
	var p projectRow
	err := r.db.Get(&p, `SELECT id, name, COALESCE(description,'') as description, created_at, updated_at FROM projects WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("CanvasEnvironment", id)
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *CanvasRepo) listAllEnvironments(ctx context.Context) ([]*models.EnvironmentConfig, error) {
	envs, err := r.environments.ListByProject(ctx, "")
	if err == nil && len(envs) > 0 {
		var result []*models.EnvironmentConfig
		for i := range envs {
			result = append(result, &envs[i])
		}
		return result, nil
	}
	var result []*models.EnvironmentConfig
	err = r.db.SelectContext(ctx, &result, `SELECT id, project_id, name, is_default, created_at, updated_at FROM environments ORDER BY is_default DESC, created_at ASC`)
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = make([]*models.EnvironmentConfig, 0)
	}
	return result, nil
}

func (r *CanvasRepo) listAllAppServices() ([]*models.AppService, error) {
	return r.scanAppServices(context.Background(), `SELECT id, project_id, environment_id, name, COALESCE(repository_url,'') as repository_url, COALESCE(image_ref,'') as image_ref, COALESCE(branch,'') as branch, internal_port, COALESCE(domain,'') as domain, COALESCE(container_id,'') as container_id, status, replicas, created_at, updated_at FROM app_services ORDER BY created_at DESC`)
}

func (r *CanvasRepo) listAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return r.scanAppServices(context.Background(), `SELECT id, project_id, environment_id, name, COALESCE(repository_url,'') as repository_url, COALESCE(image_ref,'') as image_ref, COALESCE(branch,'') as branch, internal_port, COALESCE(domain,'') as domain, COALESCE(container_id,'') as container_id, status, replicas, created_at, updated_at FROM app_services WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *CanvasRepo) listAppServicesByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	return r.scanAppServices(ctx, `SELECT id, project_id, environment_id, name, COALESCE(repository_url,'') as repository_url, COALESCE(image_ref,'') as image_ref, COALESCE(branch,'') as branch, internal_port, COALESCE(domain,'') as domain, COALESCE(container_id,'') as container_id, status, replicas, created_at, updated_at FROM app_services WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *CanvasRepo) scanAppServices(ctx context.Context, query string, args ...any) ([]*models.AppService, error) {
	var apps []*models.AppService
	err := r.db.SelectContext(ctx, &apps, query, args...)
	if err != nil {
		return nil, err
	}
	if apps == nil {
		apps = make([]*models.AppService, 0)
	}
	return apps, nil
}

func (r *CanvasRepo) listAllDatabases() ([]models.Database, error) {
	return r.scanDatabases(context.Background(), `SELECT id, COALESCE(project_id,'') as project_id, COALESCE(environment_id,'') as environment_id, name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,'') as container_id, status, COALESCE(internal_dns,'') as internal_dns, COALESCE(external_dns,'') as external_dns, created_at, updated_at FROM databases ORDER BY created_at DESC`)
}

func (r *CanvasRepo) listDatabasesByProject(projectID string) ([]models.Database, error) {
	return r.scanDatabases(context.Background(), `SELECT id, COALESCE(project_id,'') as project_id, COALESCE(environment_id,'') as environment_id, name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,'') as container_id, status, COALESCE(internal_dns,'') as internal_dns, COALESCE(external_dns,'') as external_dns, created_at, updated_at FROM databases WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *CanvasRepo) listDatabasesByEnvironment(ctx context.Context, environmentID string) ([]models.Database, error) {
	return r.scanDatabases(ctx, `SELECT id, COALESCE(project_id,'') as project_id, COALESCE(environment_id,'') as environment_id, name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,'') as container_id, status, COALESCE(internal_dns,'') as internal_dns, COALESCE(external_dns,'') as external_dns, created_at, updated_at FROM databases WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *CanvasRepo) scanDatabases(ctx context.Context, query string, args ...any) ([]models.Database, error) {
	var dbs []models.Database
	err := r.db.SelectContext(ctx, &dbs, query, args...)
	if err != nil {
		return nil, err
	}
	if dbs == nil {
		dbs = make([]models.Database, 0)
	}
	return dbs, nil
}
