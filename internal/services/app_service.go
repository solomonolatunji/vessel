package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type AppService struct {
	repo       repositories.AppServiceRepository
	varRepo    repositories.ServiceVarRepository
	volumeRepo repositories.ServiceVolumeRepository
}

func NewAppService(
	repo repositories.AppServiceRepository,
	varRepo repositories.ServiceVarRepository,
	volumeRepo repositories.ServiceVolumeRepository,
) *AppService {
	return &AppService{
		repo:       repo,
		varRepo:    varRepo,
		volumeRepo: volumeRepo,
	}
}

func (s *AppService) CreateAppService(ctx context.Context, svc *models.AppService) (*models.AppService, error) {
	if svc == nil {
		return nil, errors.New("app service is nil")
	}
	if svc.ProjectID == "" {
		return nil, errors.New("project id is required")
	}
	if svc.Name == "" {
		return nil, errors.New("name is required")
	}
	if svc.ID == "" {
		svc.ID = uuid.New().String()
	}
	now := time.Now()
	svc.CreatedAt = now
	svc.UpdatedAt = now
	if svc.Status == "" {
		svc.Status = models.AppServiceStatusCreated
	}
	if err := s.repo.Create(ctx, svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *AppService) GetAppService(ctx context.Context, id string) (*models.AppService, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load volumes
	if s.volumeRepo != nil {
		volumes, _ := s.volumeRepo.ListByService(ctx, app.ID)
		app.Volumes = volumes
	}

	return app, nil
}

func (s *AppService) ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	if environmentID == "" {
		return nil, errors.New("environment id required")
	}
	return s.repo.ListByEnvironment(ctx, environmentID)
}

func (s *AppService) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	apps, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if s.volumeRepo != nil {
		for _, app := range apps {
			volumes, _ := s.volumeRepo.ListByService(ctx, app.ID)
			app.Volumes = volumes
		}
	}
	return apps, nil
}

func (s *AppService) UpdateAppService(ctx context.Context, svc *models.AppService) error {
	if svc == nil {
		return errors.New("app service is nil")
	}
	if svc.ID == "" {
		return errors.New("valid service required for update")
	}
	svc.UpdatedAt = time.Now()
	return s.repo.Update(ctx, svc)
}

func (s *AppService) DeleteAppService(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *AppService) CreateVariable(ctx context.Context, v *models.Variable) (*models.Variable, error) {
	if v.ServiceID == "" || v.Key == "" {
		return nil, errors.New("serviceId and key required")
	}
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	now := time.Now()
	v.CreatedAt = now
	v.UpdatedAt = now
	if err := s.varRepo.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *AppService) UpdateVariable(ctx context.Context, v *models.Variable) error {
	if v.ID == "" {
		return errors.New("variable ID required")
	}
	v.UpdatedAt = time.Now()
	return s.varRepo.Update(ctx, v)
}

func (s *AppService) GetVariable(ctx context.Context, id string) (*models.Variable, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.varRepo.GetByID(ctx, id)
}

func (s *AppService) ListVariablesByService(ctx context.Context, serviceID string) ([]*models.Variable, error) {
	if serviceID == "" {
		return nil, errors.New("service ID required")
	}
	return s.varRepo.ListByService(ctx, serviceID)
}

func (s *AppService) DeleteVariable(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.varRepo.Delete(ctx, id)
}

func (s *AppService) CreateWebhook(ctx context.Context, w *models.Webhook) (*models.Webhook, error) {
	if w == nil {
		return nil, errors.New("webhook is nil")
	}
	if w.ServiceID == "" {
		return nil, errors.New("serviceId is required")
	}
	if w.URL == "" {
		return nil, errors.New("URL is required")
	}

	if w.ID == "" {
		w.ID = uuid.New().String()
	}

	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()

	if err := s.repo.CreateWebhook(ctx, w); err != nil {
		return nil, err
	}

	return w, nil
}

func (s *AppService) ListWebhooks(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	if serviceID == "" {
		return nil, errors.New("serviceId is required")
	}
	return s.repo.ListWebhooksByService(ctx, serviceID)
}

func (s *AppService) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	if id == "" || serviceID == "" {
		return errors.New("id and serviceId are required")
	}
	return s.repo.DeleteWebhook(ctx, id, serviceID)
}

func (s *AppService) CreateVolume(ctx context.Context, vol *models.ServiceVolume) (*models.ServiceVolume, error) {
	if vol.ServiceID == "" || vol.HostPath == "" || vol.ContainerPath == "" {
		return nil, errors.New("serviceId, hostPath and containerPath are required")
	}
	if vol.ID == "" {
		vol.ID = uuid.New().String()
	}
	vol.CreatedAt = time.Now()
	if err := s.volumeRepo.Create(ctx, vol); err != nil {
		return nil, err
	}
	return vol, nil
}

func (s *AppService) ListVolumes(ctx context.Context, serviceID string) ([]models.ServiceVolume, error) {
	if serviceID == "" {
		return nil, errors.New("serviceId required")
	}
	return s.volumeRepo.ListByService(ctx, serviceID)
}

func (s *AppService) DeleteVolume(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.volumeRepo.Delete(ctx, id)
}
