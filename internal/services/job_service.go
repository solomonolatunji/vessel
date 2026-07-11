package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type JobService struct {
	repo    repositories.JobRepository
	manager *engine.CronManager
}

func NewJobService(r repositories.JobRepository, m *engine.CronManager) *JobService {
	return &JobService{
		repo:    r,
		manager: m,
	}
}

func (s *JobService) CreateJob(ctx context.Context, j *models.Job) (*models.Job, error) {
	if j == nil || j.ProjectID == "" || j.Name == "" {
		return nil, errors.New("valid job with projectId and name required")
	}
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	if j.Schedule == "" {
		j.Schedule = "0 0 * * *"
	}
	if j.Status == "" {
		j.Status = "active"
	}
	now := time.Now()
	j.CreatedAt = now
	j.UpdatedAt = now
	if err := s.repo.Create(ctx, j); err != nil {
		return nil, err
	}
	if s.manager != nil && j.Status == "active" {
		_ = s.manager.RegisterJob(j)
	}
	return j, nil
}

func (s *JobService) GetJob(ctx context.Context, id string) (*models.Job, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *JobService) ListAllJobs(ctx context.Context) ([]models.Job, error) {
	return s.repo.ListAll(ctx)
}

func (s *JobService) ListJobsByProject(ctx context.Context, projectID string) ([]models.Job, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	return s.repo.ListByProject(ctx, projectID)
}

func (s *JobService) UpdateJob(ctx context.Context, j *models.Job) error {
	if j == nil || j.ID == "" {
		return errors.New("valid job required for update")
	}
	j.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, j); err != nil {
		return err
	}
	if s.manager != nil {
		s.manager.UnregisterJob(j.ID)
		if j.Status == "active" {
			_ = s.manager.RegisterJob(j)
		}
	}
	return nil
}

func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	if s.manager != nil {
		s.manager.UnregisterJob(id)
	}
	return s.repo.Delete(ctx, id)
}

func (s *JobService) ExecuteJob(ctx context.Context, id string) (string, error) {
	if s.manager == nil {
		return "", errors.New("cron manager not available")
	}
	return s.manager.ExecuteJob(ctx, id)
}
