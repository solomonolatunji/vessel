package services

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type DeploymentService struct {
	repo        repositories.DeploymentRepository
	appRepo     repositories.AppServiceRepository
	projectRepo repositories.ProjectRepository
	deployer    *engine.Deployer
}

func NewDeploymentService(
	r repositories.DeploymentRepository,
	ar repositories.AppServiceRepository,
	pr repositories.ProjectRepository,
	d *engine.Deployer,
) *DeploymentService {
	return &DeploymentService{
		repo:        r,
		appRepo:     ar,
		projectRepo: pr,
		deployer:    d,
	}
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, d *models.Deployment) (*models.Deployment, error) {
	if d == nil || d.ServiceID == "" {
		return nil, errors.New("valid deployment with serviceId required")
	}
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	if d.Status == "" {
		d.Status = "pending"
	}
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DeploymentService) GetDeployment(ctx context.Context, id string) (*models.Deployment, error) {
	if id == "" {
		return nil, errors.New("deployment id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DeploymentService) ListByService(ctx context.Context, serviceID string) ([]*models.Deployment, error) {
	if serviceID == "" {
		return nil, errors.New("service id required")
	}
	return s.repo.ListByService(ctx, serviceID)
}

func (s *DeploymentService) UpdateDeployment(ctx context.Context, d *models.Deployment) error {
	if d == nil || d.ID == "" {
		return errors.New("valid deployment required for update")
	}
	d.UpdatedAt = time.Now()
	return s.repo.Update(ctx, d)
}

func (s *DeploymentService) UpdateStatus(ctx context.Context, id, status, buildLogs, containerID string) error {
	if id == "" {
		return errors.New("deployment id required")
	}
	return s.repo.UpdateStatus(ctx, id, status, buildLogs, containerID)
}

func (s *DeploymentService) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	if s.deployer == nil || s.appRepo == nil {
		return "", errors.New("deployer or app repo not available")
	}
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return "", err
	}
	containerID, err := s.deployer.DeployAppService(ctx, app, sourceDir, logWriter)
	if err == nil && containerID != "" {
		app.ContainerID = containerID
		_ = s.appRepo.Update(ctx, app)
	}
	return containerID, err
}

func (s *DeploymentService) DeployProject(ctx context.Context, projectID, sourceDir string, logWriter io.Writer) (string, error) {
	if s.deployer == nil || s.projectRepo == nil {
		return "", errors.New("deployer or project repo not available")
	}
	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return "", err
	}
	containerID, err := s.deployer.Deploy(ctx, project, sourceDir, logWriter)
	if err == nil && containerID != "" {
		apps, appErr := s.appRepo.ListByProject(ctx, projectID)
		if appErr == nil && len(apps) > 0 {
			apps[0].ContainerID = containerID
			_ = s.appRepo.Update(ctx, apps[0])
		}
	}
	return containerID, err
}
