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

func (r *CanvasSQLiteRepository) GetEnvironmentCanvas(_ context.Context, environmentID string) (*models.EnvironmentCanvas, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	row := r.db.QueryRow(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = ?`, environmentID,
	)
	var env models.EnvironmentConfig
	var isDefault int
	err := row.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("environment not found: %s", environmentID)
	}
	if err != nil {
		return nil, err
	}
	env.IsDefault = isDefault == 1
	apps, _ := r.listAppServicesByEnvironment(environmentID)
	dbs, _ := r.listDatabasesByEnvironment(environmentID)
	storageItems, _ := r.listStorageByEnvironment(environmentID)
	var dbsPtrs []*models.Database
	for i := range dbs {
		dbsPtrs = append(dbsPtrs, &dbs[i])
	}
	var storagePtrs []*models.Storage
	for i := range storageItems {
		storagePtrs = append(storagePtrs, &storageItems[i])
	}
	return &models.EnvironmentCanvas{
		Environment: &env,
		Apps:        apps,
		Databases:   dbsPtrs,
		Storage:     storagePtrs,
	}, nil
}

type projectRow struct {
	ID          string
	WorkspaceID string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *CanvasSQLiteRepository) listAllProjects() ([]projectRow, error) {
	rows, err := r.db.Query(`SELECT id, COALESCE(workspace_id, ''), name, COALESCE(description,''), created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []projectRow
	for rows.Next() {
		var p projectRow
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *CanvasSQLiteRepository) getProject(id string) (*projectRow, error) {
	row := r.db.QueryRow(`SELECT id, COALESCE(workspace_id, ''), name, COALESCE(description,''), created_at, updated_at FROM projects WHERE id = ?`, id)
	var p projectRow
	err := row.Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("CanvasEnvironment", id)
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *CanvasSQLiteRepository) listAllEnvironments(ctx context.Context) ([]*models.EnvironmentConfig, error) {
	envs, err := r.environments.ListByProject(ctx, "")
	if err == nil && len(envs) > 0 {
		var result []*models.EnvironmentConfig
		for i := range envs {
			result = append(result, &envs[i])
		}
		return result, nil
	}
	rows, err := r.db.Query(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments ORDER BY is_default DESC, created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*models.EnvironmentConfig
	for rows.Next() {
		var env models.EnvironmentConfig
		var isDefault int
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, err
		}
		env.IsDefault = isDefault == 1
		result = append(result, &env)
	}
	return result, rows.Err()
}

func (r *CanvasSQLiteRepository) listAllAppServices() ([]*models.AppService, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(repository_url,''), COALESCE(branch,''), internal_port, COALESCE(domain,''), COALESCE(container_id,''), status, created_at, updated_at FROM app_services ORDER BY created_at DESC`)
}

func (r *CanvasSQLiteRepository) listAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(repository_url,''), COALESCE(branch,''), internal_port, COALESCE(domain,''), COALESCE(container_id,''), status, created_at, updated_at FROM app_services WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *CanvasSQLiteRepository) listAppServicesByEnvironment(environmentID string) ([]*models.AppService, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(repository_url,''), COALESCE(branch,''), internal_port, COALESCE(domain,''), COALESCE(container_id,''), status, created_at, updated_at FROM app_services WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *CanvasSQLiteRepository) scanAppServices(query string, args ...any) ([]*models.AppService, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []*models.AppService
	for rows.Next() {
		var a models.AppService
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.EnvironmentID, &a.Name,
			&a.RepositoryURL, &a.Branch, &a.InternalPort,
			&a.Domain, &a.ContainerID, &a.Status, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		apps = append(apps, &a)
	}
	return apps, rows.Err()
}

func (r *CanvasSQLiteRepository) listAllDatabases() ([]models.Database, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases ORDER BY created_at DESC`)
}

func (r *CanvasSQLiteRepository) listDatabasesByProject(projectID string) ([]models.Database, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *CanvasSQLiteRepository) listDatabasesByEnvironment(environmentID string) ([]models.Database, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *CanvasSQLiteRepository) scanDatabases(query string, args ...any) ([]models.Database, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dbs []models.Database
	for rows.Next() {
		var d models.Database
		if err := rows.Scan(
			&d.ID, &d.ProjectID, &d.EnvironmentID, &d.Name, &d.Engine, &d.Version, &d.Port,
			&d.Username, &d.DatabaseName, &d.VolumePath, &d.ContainerID, &d.Status,
			&d.InternalDNS, &d.ExternalDNS, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		dbs = append(dbs, d)
	}
	return dbs, rows.Err()
}

func (r *CanvasSQLiteRepository) listAllStorage() ([]models.Storage, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage ORDER BY created_at DESC`)
}

func (r *CanvasSQLiteRepository) listStorageByProject(projectID string) ([]models.Storage, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *CanvasSQLiteRepository) listStorageByEnvironment(environmentID string) ([]models.Storage, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *CanvasSQLiteRepository) scanStorage(query string, args ...any) ([]models.Storage, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Storage
	for rows.Next() {
		var s models.Storage
		if err := rows.Scan(
			&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Name, &s.Type, &s.APIPort, &s.ConsolePort,
			&s.AccessKey, &s.BucketName, &s.VolumePath, &s.ContainerID, &s.Status,
			&s.InternalDNS, &s.ExternalDNS, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}
