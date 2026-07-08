package store

import (
	"fmt"

	"github.com/solomonolatunji/vessel/internal/types"
)

// GetProjectCanvasSummary calculates aggregated counts and status icons for a single project canvas.
func (s *Store) GetProjectCanvasSummary(projectID string) (*types.ProjectCanvasSummary, error) {
	project, err := s.GetProject(projectID)
	if err != nil || project == nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	envs, _ := s.ListEnvironments(projectID)
	apps, _ := s.ListAppServicesByProject(projectID)
	dbs, _ := s.ListDatabasesByProject(projectID)
	storage, _ := s.ListStorageByProject(projectID)

	defaultEnv, _ := s.GetDefaultEnvironment(projectID)

	summary := &types.ProjectCanvasSummary{
		ProjectConfig:      *project,
		EnvironmentsCount:  len(envs),
		AppsCount:          len(apps),
		DatabasesCount:     len(dbs),
		StorageCount:       len(storage),
		TotalServices:      len(apps) + len(dbs) + len(storage),
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
	for _, st := range storage {
		if st.Status == "running" {
			onlineCount++
		}
		summary.ServiceIcons = append(summary.ServiceIcons, st.Type)
	}
	summary.OnlineServices = onlineCount

	return summary, nil
}

// ListProjectCanvasSummaries returns canvas summaries for all projects on the dashboard.
func (s *Store) ListProjectCanvasSummaries() ([]*types.ProjectCanvasSummary, error) {
	projects, err := s.ListProjects()
	if err != nil {
		return nil, err
	}

	var summaries []*types.ProjectCanvasSummary
	for _, p := range projects {
		summary, err := s.GetProjectCanvasSummary(p.ID)
		if err == nil && summary != nil {
			summaries = append(summaries, summary)
		}
	}
	return summaries, nil
}

// GetEnvironmentCanvas retrieves all Git applications, databases, and storage buckets inside a specific environment.
func (s *Store) GetEnvironmentCanvas(environmentID string) (*types.EnvironmentCanvas, error) {
	env, err := s.GetEnvironment(environmentID)
	if err != nil || env == nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	apps, _ := s.ListAppServicesByEnvironment(environmentID)
	dbs, _ := s.ListDatabasesByEnvironment(environmentID)
	storage, _ := s.ListStorageByEnvironment(environmentID)

	var dbsPtrs []*types.DatabaseConfig
	for i := range dbs {
		dbsPtrs = append(dbsPtrs, &dbs[i])
	}
	var storagePtrs []*types.StorageConfig
	for i := range storage {
		storagePtrs = append(storagePtrs, &storage[i])
	}

	canvas := &types.EnvironmentCanvas{
		Environment: env,
		Apps:        apps,
		Databases:   dbsPtrs,
		Storage:     storagePtrs,
	}
	return canvas, nil
}
