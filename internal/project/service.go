package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/types"
	"vessel.dev/vessel/internal/utils"
)

// AppServiceRepository is the minimal surface project.Service needs from the app-service domain.
type AppServiceRepository interface {
	CreateAppService(ctx context.Context, app *types.AppServiceConfig) error
}

// Service implements the project domain business logic.
type Service struct {
	repo     Repository
	apps     AppServiceRepository
	domainFn func(name string) string
}

// NewService creates a new project Service.
func NewService(repo Repository, apps AppServiceRepository) *Service {
	return &Service{
		repo:     repo,
		apps:     apps,
		domainFn: func(name string) string { return utils.GenerateSslipDomain(name, "") },
	}
}

// List returns all projects.
func (s *Service) List(ctx context.Context) ([]ProjectConfig, error) {
	return s.repo.List(ctx)
}

// Get returns a single project by ID.
func (s *Service) Get(ctx context.Context, id string) (*ProjectConfig, error) {
	return s.repo.Get(ctx, id)
}

// Create creates a project and its default application service.
func (s *Service) Create(ctx context.Context, req *CreateProjectRequest) (*ProjectConfig, error) {
	if req.Name == "" {
		req.Name = fmt.Sprintf("project-%s", uuid.NewString()[:8])
	}

	p := &ProjectConfig{
		ID:          req.ID,
		TeamID:      req.TeamID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	port := req.InternalPort
	if port <= 0 {
		port = req.InternalPortSnake
	}
	if port <= 0 {
		port = 3000
	}

	repo := req.RepositoryURL
	if repo == "" {
		repo = req.RepositoryURLSnake
	}

	domain := req.Domain
	if domain == "" {
		domain = s.domainFn(req.Name)
	}

	branch := req.Branch
	if branch == "" {
		branch = "main"
	}

	app := &types.AppServiceConfig{
		ProjectID:     p.ID,
		EnvironmentID: "env-prod",
		Name:          req.Name,
		RepositoryURL: repo,
		Branch:        branch,
		InternalPort:  port,
		Domain:        domain,
	}
	_ = s.apps.CreateAppService(ctx, app)

	return p, nil
}

// Delete removes a project by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
