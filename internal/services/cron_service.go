package services

import (
	"context"
	"errors"
	"fmt"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type CronService struct {
	jobs        repositories.JobRepository
	projects    repositories.ProjectRepository
	cronManager *engine.CronManager
}

func NewCronService(js repositories.JobRepository, ps repositories.ProjectRepository, cm *engine.CronManager) *CronService {
	return &CronService{
		jobs:        js,
		projects:    ps,
		cronManager: cm,
	}
}

func (cs *CronService) CreateJob(ctx context.Context, j *models.Job) error {
	if j.ProjectID == "" {
		return errors.New("projectId is required when creating a scheduled job")
	}
	if j.Schedule == "" {
		return errors.New("schedule cron expression is required")
	}
	if j.Command == "" {
		return errors.New("command is required")
	}
	project, err := cs.projects.Get(ctx, j.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to verify project existence: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project with ID %s not found", j.ProjectID)
	}
	if err := cs.jobs.Create(ctx, j); err != nil {
		return err
	}
	return cs.cronManager.RegisterJob(j)
}

func (cs *CronService) GetJob(ctx context.Context, id string) (*models.Job, error) {
	return cs.jobs.GetByID(ctx, id)
}

func (cs *CronService) ListJobs(ctx context.Context, projectID string) ([]models.Job, error) {
	if projectID != "" {
		return cs.jobs.ListByProject(ctx, projectID)
	}
	return cs.jobs.ListAll(ctx)
}

func (cs *CronService) TriggerJobImmediately(ctx context.Context, jobID string) (string, error) {
	return cs.cronManager.ExecuteJob(ctx, jobID)
}

func (cs *CronService) DeleteJob(ctx context.Context, id string) error {
	cs.cronManager.UnregisterJob(id)
	return cs.jobs.Delete(ctx, id)
}
