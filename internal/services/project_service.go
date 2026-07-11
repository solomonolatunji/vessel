package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/utils"
)

type ProjectService struct {
	projectRepo    repositories.ProjectRepository
	envRepo        repositories.EnvironmentRepository
	appServiceRepo repositories.AppServiceRepository
	serviceVarRepo repositories.ServiceVarRepository
}

func NewProjectService(pr repositories.ProjectRepository, er repositories.EnvironmentRepository, ar repositories.AppServiceRepository, svr repositories.ServiceVarRepository) *ProjectService {
	return &ProjectService{
		projectRepo:    pr,
		envRepo:        er,
		appServiceRepo: ar,
		serviceVarRepo: svr,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, name, workspaceID, teamID, description string) (*models.ProjectConfig, error) {
	if name == "" {
		return nil, errors.New("project name is required")
	}
	p := &models.ProjectConfig{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		TeamID:      teamID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.projectRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) CreateProjectFromRequest(ctx context.Context, req *models.CreateProjectRequest) (*models.ProjectConfig, error) {
	if req.Name == "" {
		req.Name = fmt.Sprintf("project-%s", uuid.NewString()[:8])
	}
	id := req.ID
	if id == "" {
		id = uuid.NewString()
	}
	p := &models.ProjectConfig{
		ID:          id,
		TeamID:      req.TeamID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.projectRepo.Create(ctx, p); err != nil {
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
		domain = utils.GenerateSslipDomain(req.Name, "")
	}
	branch := req.Branch
	if branch == "" {
		branch = "main"
	}
	app := &models.AppService{
		ID:            uuid.NewString(),
		ProjectID:     p.ID,
		EnvironmentID: "env-prod",
		Name:          req.Name,
		RepositoryURL: repo,
		Branch:        branch,
		InternalPort:  port,
		Domain:        domain,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_ = s.appServiceRepo.Create(ctx, app)
	return p, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*models.ProjectConfig, error) {
	if id == "" {
		return nil, errors.New("project id is required")
	}
	return s.projectRepo.Get(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context) ([]models.ProjectConfig, error) {
	return s.projectRepo.List(ctx)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("project id is required")
	}
	return s.projectRepo.Delete(ctx, id)
}

func (s *ProjectService) CreateEnvironment(ctx context.Context, projectID, name string) (*models.EnvironmentConfig, error) {
	if projectID == "" || name == "" {
		return nil, errors.New("project id and environment name required")
	}
	env := &models.EnvironmentConfig{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Name:      name,
	}
	if err := s.envRepo.Create(ctx, env); err != nil {
		return nil, err
	}
	return env, nil
}

func (s *ProjectService) ListEnvironments(ctx context.Context, projectID string) ([]models.EnvironmentConfig, error) {
	if projectID == "" {
		return nil, errors.New("project id is required")
	}
	return s.envRepo.ListByProject(ctx, projectID)
}

func (s *ProjectService) DeleteEnvironment(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("environment id is required")
	}
	return s.envRepo.Delete(ctx, id)
}

func (s *ProjectService) CreateAppService(ctx context.Context, svc *models.AppService) error {
	if svc == nil || svc.Name == "" || svc.ProjectID == "" {
		return errors.New("valid app service with name and projectId required")
	}
	if svc.ID == "" {
		svc.ID = uuid.New().String()
	}
	if svc.Status == "" {
		svc.Status = "stopped"
	}
	if svc.InternalPort == 0 {
		svc.InternalPort = 3000
	}
	svc.CreatedAt = time.Now()
	svc.UpdatedAt = time.Now()
	return s.appServiceRepo.Create(ctx, svc)
}

func (s *ProjectService) GetAppService(ctx context.Context, id string) (*models.AppService, error) {
	if id == "" {
		return nil, errors.New("service id is required")
	}
	return s.appServiceRepo.GetByID(ctx, id)
}

func (s *ProjectService) ListAppServicesByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	if projectID == "" {
		return nil, errors.New("project id is required")
	}
	return s.appServiceRepo.ListByProject(ctx, projectID)
}

func (s *ProjectService) UpdateAppService(ctx context.Context, svc *models.AppService) error {
	if svc == nil || svc.ID == "" {
		return errors.New("valid app service required for update")
	}
	svc.UpdatedAt = time.Now()
	return s.appServiceRepo.Update(ctx, svc)
}

func (s *ProjectService) DeleteAppService(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("service id is required")
	}
	return s.appServiceRepo.Delete(ctx, id)
}

func (s *ProjectService) CreateServiceVariable(ctx context.Context, v *models.Variable) error {
	if v == nil || v.ServiceID == "" || v.Key == "" {
		return errors.New("valid variable required")
	}
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	return s.serviceVarRepo.Create(ctx, v)
}

func (s *ProjectService) ListServiceVariables(ctx context.Context, serviceID string) ([]*models.Variable, error) {
	if serviceID == "" {
		return nil, errors.New("service id is required")
	}
	return s.serviceVarRepo.ListByService(ctx, serviceID)
}

func (s *ProjectService) DeleteServiceVariable(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("variable id is required")
	}
	return s.serviceVarRepo.Delete(ctx, id)
}
