package store

import (
	"context"

	"vessel.dev/vessel/internal/domain"
	"vessel.dev/vessel/internal/env"
	"vessel.dev/vessel/internal/environment"
	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/types"
	"vessel.dev/vessel/internal/user"
)

// ── Settings ─────────────────────────────────────────────────────────────────

// GetServerSettings delegates to the modular settings repository.
func (s *Store) GetServerSettings() (*settings.ServerSettings, error) {
	return settings.NewSQLiteRepository(s.db).GetServerSettings(context.Background())
}

// UpdateServerSettings delegates to the modular settings repository.
func (s *Store) UpdateServerSettings(cfg *settings.ServerSettings) error {
	return settings.NewSQLiteRepository(s.db).UpdateServerSettings(context.Background(), cfg)
}

// ── User ─────────────────────────────────────────────────────────────────────

// CreateUser delegates to the modular user repository.
func (s *Store) CreateUser(u *user.User) error {
	return user.NewSQLiteRepository(s.db).CreateUser(context.Background(), u)
}

// GetUserByEmail delegates to the modular user repository.
func (s *Store) GetUserByEmail(email string) (*user.User, error) {
	return user.NewSQLiteRepository(s.db).GetUserByEmail(context.Background(), email)
}

// GetUserByID delegates to the modular user repository.
func (s *Store) GetUserByID(id string) (*user.User, error) {
	return user.NewSQLiteRepository(s.db).GetUserByID(context.Background(), id)
}

// ListUsers delegates to the modular user repository.
func (s *Store) ListUsers() ([]user.User, error) {
	return user.NewSQLiteRepository(s.db).ListUsers(context.Background())
}

// UpdateUser delegates to the modular user repository.
func (s *Store) UpdateUser(u *user.User) error {
	return user.NewSQLiteRepository(s.db).UpdateUser(context.Background(), u)
}

// CreatePersonalAccessToken delegates to the modular user repository.
func (s *Store) CreatePersonalAccessToken(pat *user.PersonalAccessToken) error {
	return user.NewSQLiteRepository(s.db).CreatePAT(context.Background(), pat)
}

// ListPersonalAccessTokens delegates to the modular user repository.
func (s *Store) ListPersonalAccessTokens(userID string) ([]*user.PersonalAccessToken, error) {
	return user.NewSQLiteRepository(s.db).ListPATs(context.Background(), userID)
}

// DeletePersonalAccessToken delegates to the modular user repository.
func (s *Store) DeletePersonalAccessToken(id, userID string) error {
	return user.NewSQLiteRepository(s.db).DeletePAT(context.Background(), id, userID)
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func (s *Store) projectRepo() project.Repository {
	return project.NewSQLiteRepository(s.db, environment.NewSQLiteRepository(s.db))
}

func (s *Store) environmentRepo() environment.Repository {
	return environment.NewSQLiteRepository(s.db)
}

func (s *Store) domainRepo() domain.Repository {
	return domain.NewSQLiteRepository(s.db)
}

func (s *Store) projectEnvRepo() env.Repository {
	return env.NewSQLiteRepository(s.db, s.vault)
}

// ── Projects ─────────────────────────────────────────────────────────────────

// ListProjects delegates to the modular project repository.
func (s *Store) ListProjects() ([]types.ProjectConfig, error) {
	projects, err := s.projectRepo().List(context.Background())
	if err != nil {
		return nil, err
	}
	result := make([]types.ProjectConfig, len(projects))
	for i, p := range projects {
		result[i] = toTypesProjectConfig(p)
	}
	return result, nil
}

// ListProjectsByWorkspace delegates to the modular project repository.
func (s *Store) ListProjectsByWorkspace(workspaceID string) ([]types.ProjectConfig, error) {
	all, err := s.projectRepo().List(context.Background())
	if err != nil {
		return nil, err
	}
	var result []types.ProjectConfig
	for _, p := range all {
		if p.WorkspaceID == workspaceID {
			result = append(result, toTypesProjectConfig(p))
		}
	}
	return result, nil
}

// GetProject delegates to the modular project repository.
func (s *Store) GetProject(id string) (*types.ProjectConfig, error) {
	p, err := s.projectRepo().Get(context.Background(), id)
	if err != nil || p == nil {
		return nil, err
	}
	cfg := toTypesProjectConfig(*p)
	return &cfg, nil
}

// CreateProject delegates to the modular project repository.
func (s *Store) CreateProject(p *types.ProjectConfig) error {
	cfg := project.ProjectConfig{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		TeamID:      p.TeamID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if err := s.projectRepo().Create(context.Background(), &cfg); err != nil {
		return err
	}
	p.ID = cfg.ID
	return nil
}

// DeleteProject delegates to the modular project repository.
func (s *Store) DeleteProject(id string) error {
	return s.projectRepo().Delete(context.Background(), id)
}

func toTypesProjectConfig(p project.ProjectConfig) types.ProjectConfig {
	return types.ProjectConfig{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		TeamID:      p.TeamID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ── Environments ─────────────────────────────────────────────────────────────

func toTypesEnvironmentConfig(e environment.Config) types.EnvironmentConfig {
	return types.EnvironmentConfig{
		ID:        e.ID,
		ProjectID: e.ProjectID,
		Name:      e.Name,
		IsDefault: e.IsDefault,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ListEnvironments delegates to the modular environment repository.
func (s *Store) ListEnvironments(projectID string) ([]types.EnvironmentConfig, error) {
	envs, err := s.environmentRepo().ListByProject(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	result := make([]types.EnvironmentConfig, len(envs))
	for i, e := range envs {
		result[i] = toTypesEnvironmentConfig(e)
	}
	return result, nil
}

// GetEnvironment delegates to the modular environment repository.
func (s *Store) GetEnvironment(id string) (*types.EnvironmentConfig, error) {
	env, err := s.environmentRepo().Get(context.Background(), id)
	if err != nil || env == nil {
		return nil, err
	}
	cfg := toTypesEnvironmentConfig(*env)
	return &cfg, nil
}

// CreateEnvironment delegates to the modular environment repository.
func (s *Store) CreateEnvironment(env *types.EnvironmentConfig) error {
	cfg := environment.Config{
		ID:        env.ID,
		ProjectID: env.ProjectID,
		Name:      env.Name,
		IsDefault: env.IsDefault,
		CreatedAt: env.CreatedAt,
		UpdatedAt: env.UpdatedAt,
	}
	if err := s.environmentRepo().Create(context.Background(), &cfg); err != nil {
		return err
	}
	env.ID = cfg.ID
	return nil
}

// DeleteEnvironment delegates to the modular environment repository.
func (s *Store) DeleteEnvironment(id string) error {
	return s.environmentRepo().Delete(context.Background(), id)
}

// ── Domains ──────────────────────────────────────────────────────────────────

func toTypesDomainConfig(d domain.Config) types.DomainConfig {
	return types.DomainConfig{
		ID:            d.ID,
		ProjectID:     d.ProjectID,
		DomainName:    d.DomainName,
		RedirectTo:    d.RedirectTo,
		SSLCertStatus: d.SSLCertStatus,
		PathPrefix:    d.PathPrefix,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

// ListDomains delegates to the modular domain repository.
func (s *Store) ListDomains(projectID string) ([]types.DomainConfig, error) {
	domains, err := s.domainRepo().ListByProject(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	result := make([]types.DomainConfig, len(domains))
	for i, d := range domains {
		result[i] = toTypesDomainConfig(d)
	}
	return result, nil
}

// ListAllDomains returns every custom domain across all projects.
func (s *Store) ListAllDomains() ([]types.DomainConfig, error) {
	domains, err := s.domainRepo().ListAll(context.Background())
	if err != nil {
		return nil, err
	}
	result := make([]types.DomainConfig, len(domains))
	for i, d := range domains {
		result[i] = toTypesDomainConfig(d)
	}
	return result, nil
}

// AddDomain delegates to the modular domain repository.
func (s *Store) AddDomain(d *types.DomainConfig) error {
	cfg := domain.Config{
		ID:            d.ID,
		ProjectID:     d.ProjectID,
		DomainName:    d.DomainName,
		RedirectTo:    d.RedirectTo,
		SSLCertStatus: d.SSLCertStatus,
		PathPrefix:    d.PathPrefix,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
	if err := s.domainRepo().Create(context.Background(), &cfg); err != nil {
		return err
	}
	d.ID = cfg.ID
	return nil
}

// DeleteDomain delegates to the modular domain repository.
func (s *Store) DeleteDomain(id string) error {
	return s.domainRepo().Delete(context.Background(), id)
}

// ── Env Vars ─────────────────────────────────────────────────────────────────

// GetEnvVars delegates to the modular project_env repository.
func (s *Store) GetEnvVars(projectID string) (map[string]string, error) {
	return s.projectEnvRepo().GetVars(context.Background(), projectID)
}

// SetEnvVar delegates to the modular project_env repository.
func (s *Store) SetEnvVar(projectID, key, value string) error {
	return s.projectEnvRepo().SetVar(context.Background(), projectID, key, value)
}
