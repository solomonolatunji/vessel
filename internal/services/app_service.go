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
	appRepo repositories.AppServiceRepository
	varRepo repositories.ServiceVarRepository
}

func NewAppService(ar repositories.AppServiceRepository, vr repositories.ServiceVarRepository) *AppService {
	return &AppService{
		appRepo: ar,
		varRepo: vr,
	}
}

func (s *AppService) CreateAppService(ctx context.Context, svc *models.AppService) (*models.AppService, error) {
	if svc == nil || svc.ProjectID == "" || svc.Name == "" {
		return nil, errors.New("valid service with projectId and name required")
	}
	if svc.ID == "" {
		svc.ID = uuid.New().String()
	}
	if svc.Status == "" {
		svc.Status = models.AppServiceStatusCreated
	}
	now := time.Now()
	if svc.CreatedAt.IsZero() {
		svc.CreatedAt = now
	}
	svc.UpdatedAt = now
	if err := s.appRepo.Create(ctx, svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *AppService) GetAppService(ctx context.Context, id string) (*models.AppService, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.appRepo.GetByID(ctx, id)
}

func (s *AppService) ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	if environmentID == "" {
		return nil, errors.New("environment id required")
	}
	return s.appRepo.ListByEnvironment(ctx, environmentID)
}

func (s *AppService) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	return s.appRepo.ListByProject(ctx, projectID)
}

func (s *AppService) UpdateAppService(ctx context.Context, svc *models.AppService) error {
	if svc == nil || svc.ID == "" {
		return errors.New("valid service required for update")
	}
	svc.UpdatedAt = time.Now()
	return s.appRepo.Update(ctx, svc)
}

func (s *AppService) DeleteAppService(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.appRepo.Delete(ctx, id)
}

func (s *AppService) CreateVariable(ctx context.Context, v *models.Variable) (*models.Variable, error) {
	if v == nil || v.ServiceID == "" || v.Key == "" {
		return nil, errors.New("valid variable with serviceId and key required")
	}
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	now := time.Now()
	if v.CreatedAt.IsZero() {
		v.CreatedAt = now
	}
	v.UpdatedAt = now
	if err := s.varRepo.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *AppService) UpdateVariable(ctx context.Context, v *models.Variable) error {
	if v == nil || v.ID == "" {
		return errors.New("valid variable with id required")
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
		return nil, errors.New("service id required")
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
	if w == nil || w.ServiceID == "" || w.URL == "" {
		return nil, errors.New("valid webhook with serviceId and url required")
	}
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	w.CreatedAt = time.Now()
	if err := s.appRepo.CreateWebhook(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *AppService) ListWebhooks(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	if serviceID == "" {
		return nil, errors.New("serviceId required")
	}
	return s.appRepo.ListWebhooksByService(ctx, serviceID)
}

func (s *AppService) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	if id == "" || serviceID == "" {
		return errors.New("id and serviceId required")
	}
	return s.appRepo.DeleteWebhook(ctx, id, serviceID)
}
