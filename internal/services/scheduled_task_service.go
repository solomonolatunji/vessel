package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type ScheduledTaskService struct {
	repo    repositories.ScheduledTaskRepository
	manager *engine.CronManager
}

func NewScheduledTaskService(r repositories.ScheduledTaskRepository, m *engine.CronManager) *ScheduledTaskService {
	return &ScheduledTaskService{
		repo:    r,
		manager: m,
	}
}

func (s *ScheduledTaskService) CreateScheduledTask(ctx context.Context, j *models.ScheduledTask) (*models.ScheduledTask, error) {
	if j == nil || j.ServiceID == "" || j.Name == "" {
		return nil, errors.New("valid scheduledTask with serviceId and name required")
	}
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	if j.Schedule == "" {
		j.Schedule = "0 0 * * *"
	}
	if j.Status == "" {
		j.Status = string(models.ScheduledTaskStatusActive)
	}
	now := time.Now()
	j.CreatedAt = now
	j.UpdatedAt = now
	if err := s.repo.Create(ctx, j); err != nil {
		return nil, err
	}
	if s.manager != nil && j.Status == string(models.ScheduledTaskStatusActive) {
		_ = s.manager.RegisterScheduledTask(j)
	}
	return j, nil
}

func (s *ScheduledTaskService) GetScheduledTask(ctx context.Context, id string) (*models.ScheduledTask, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ScheduledTaskService) ListAllScheduledTasks(ctx context.Context) ([]models.ScheduledTask, error) {
	return s.repo.ListAll(ctx)
}

func (s *ScheduledTaskService) ListScheduledTasksByProject(ctx context.Context, projectID string) ([]models.ScheduledTask, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	return s.repo.ListByProject(ctx, projectID)
}

func (s *ScheduledTaskService) ListScheduledTasksByService(ctx context.Context, serviceID string) ([]models.ScheduledTask, error) {
	if serviceID == "" {
		return nil, errors.New("service id required")
	}
	return s.repo.ListByService(ctx, serviceID)
}

func (s *ScheduledTaskService) UpdateScheduledTask(ctx context.Context, j *models.ScheduledTask) error {
	if j == nil || j.ID == "" {
		return errors.New("valid scheduledTask required for update")
	}
	j.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, j); err != nil {
		return err
	}
	if s.manager != nil {
		s.manager.UnregisterScheduledTask(j.ID)
		if j.Status == string(models.ScheduledTaskStatusActive) {
			_ = s.manager.RegisterScheduledTask(j)
		}
	}
	return nil
}

func (s *ScheduledTaskService) DeleteScheduledTask(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	if s.manager != nil {
		s.manager.UnregisterScheduledTask(id)
	}
	return s.repo.Delete(ctx, id)
}

func (s *ScheduledTaskService) ExecuteScheduledTask(ctx context.Context, id string) (string, error) {
	if s.manager == nil {
		return "", errors.New("cron manager not available")
	}
	return s.manager.ExecuteScheduledTask(ctx, id)
}
