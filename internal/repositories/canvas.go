package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"vessl.dev/vessl/internal/models"
)

type CanvasRepository interface {
	ListCanvasSummaries(ctx context.Context) ([]models.CanvasSummary, error)
	GetCanvasSummary(ctx context.Context, id string) (*models.CanvasSummary, error)
	GetEnvironmentCanvas(ctx context.Context, id string) (*models.EnvironmentCanvas, error)
}

type CanvasSQLiteRepository struct {
	db           *sql.DB
	mu           sync.Mutex
	environments EnvironmentRepository
}

func NewCanvasSQLiteRepository(db *sql.DB, envRepo EnvironmentRepository) *CanvasSQLiteRepository {
	return &CanvasSQLiteRepository{db: db, environments: envRepo}
}

func (r *CanvasSQLiteRepository) ListCanvasSummaries(ctx context.Context) ([]models.CanvasSummary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	projects, err := r.listAllProjects()
	if err != nil {
		return nil, err
	}
	allEnvs, err := r.listAllEnvironments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all environments: %w", err)
	}
	allApps, err := r.listAllAppServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list all app services: %w", err)
	}
	allDbs, err := r.listAllDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list all databases: %w", err)
	}
	allStorage, err := r.listAllStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to list all storage: %w", err)
	}
	envsByProject := make(map[string][]*models.EnvironmentConfig)
	for _, e := range allEnvs {
		envsByProject[e.ProjectID] = append(envsByProject[e.ProjectID], e)
	}
	appsByProject := make(map[string][]*models.AppService)
	for _, a := range allApps {
		appsByProject[a.ProjectID] = append(appsByProject[a.ProjectID], a)
	}
	dbsByProject := make(map[string][]models.Database)
	for _, d := range allDbs {
		dbsByProject[d.ProjectID] = append(dbsByProject[d.ProjectID], d)
	}
	storageByProject := make(map[string][]models.Storage)
	for _, st := range allStorage {
		storageByProject[st.ProjectID] = append(storageByProject[st.ProjectID], st)
	}
	var summaries []models.CanvasSummary
	for _, project := range projects {
		envs := envsByProject[project.ID]
		apps := appsByProject[project.ID]
		dbs := dbsByProject[project.ID]
		storageItems := storageByProject[project.ID]
		var defaultEnv *models.EnvironmentConfig
		if len(envs) > 0 {
			for _, e := range envs {
				if e.IsDefault {
					defaultEnv = e
					break
				}
			}
			if defaultEnv == nil {
				defaultEnv = envs[0]
			}
		}
		summary := models.CanvasSummary{
			ID:                 project.ID,
			WorkspaceID:        project.WorkspaceID,
			Name:               project.Name,
			Description:        project.Description,
			CreatedAt:          project.CreatedAt,
			UpdatedAt:          project.UpdatedAt,
			EnvironmentsCount:  len(envs),
			AppsCount:          len(apps),
			DatabasesCount:     len(dbs),
			StorageCount:       len(storageItems),
			TotalServices:      len(apps) + len(dbs) + len(storageItems),
			DefaultEnvironment: defaultEnv,
			ServiceIcons:       make([]string, 0),
		}
		onlineCount := 0
		for _, app := range apps {
			if app.Status == "running" {
				onlineCount++
			}
			summary.ServiceIcons = append(summary.ServiceIcons, "github")
		}
		for _, db := range dbs {
			if db.Status == "running" {
				onlineCount++
			}
			summary.ServiceIcons = append(summary.ServiceIcons, db.Engine)
		}
		for _, st := range storageItems {
			if st.Status == "running" {
				onlineCount++
			}
			summary.ServiceIcons = append(summary.ServiceIcons, st.Type)
		}
		summary.OnlineServices = onlineCount
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (r *CanvasSQLiteRepository) GetCanvasSummary(ctx context.Context, id string) (*models.CanvasSummary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	project, err := r.getProject(id)
	if err != nil || project == nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	envs, _ := r.environments.ListByProject(ctx, id)
	apps, _ := r.listAppServicesByProject(id)
	dbs, _ := r.listDatabasesByProject(id)
	storageItems, _ := r.listStorageByProject(id)
	var defaultEnv *models.EnvironmentConfig
	if len(envs) > 0 {
		for _, e := range envs {
			e := e
			if e.IsDefault {
				defaultEnv = &e
				break
			}
		}
		if defaultEnv == nil {
			defaultEnv = &envs[0]
		}
	}
	summary := &models.CanvasSummary{
		ID:                 project.ID,
		WorkspaceID:        project.WorkspaceID,
		Name:               project.Name,
		Description:        project.Description,
		CreatedAt:          project.CreatedAt,
		UpdatedAt:          project.UpdatedAt,
		EnvironmentsCount:  len(envs),
		AppsCount:          len(apps),
		DatabasesCount:     len(dbs),
		StorageCount:       len(storageItems),
		TotalServices:      len(apps) + len(dbs) + len(storageItems),
		DefaultEnvironment: defaultEnv,
		ServiceIcons:       make([]string, 0),
	}
	onlineCount := 0
	for _, app := range apps {
		if app.Status == "running" {
			onlineCount++
		}
		summary.ServiceIcons = append(summary.ServiceIcons, "github")
	}
	for _, db := range dbs {
		if db.Status == "running" {
			onlineCount++
		}
		summary.ServiceIcons = append(summary.ServiceIcons, db.Engine)
	}
	for _, st := range storageItems {
		if st.Status == "running" {
			onlineCount++
		}
		summary.ServiceIcons = append(summary.ServiceIcons, st.Type)
	}
	summary.OnlineServices = onlineCount
	return summary, nil
}
