package services

import (
	"context"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type EmailSettingsService struct {
	repo repositories.WorkspaceEmailSettingsRepository
}

func NewEmailSettingsService(repo repositories.WorkspaceEmailSettingsRepository) *EmailSettingsService {
	return &EmailSettingsService{repo: repo}
}

func (s *EmailSettingsService) GetWorkspaceEmailSettings(ctx context.Context, workspaceID string) (*models.WorkspaceEmailSettings, error) {
	return s.repo.GetByWorkspaceID(ctx, workspaceID)
}

func (s *EmailSettingsService) SaveWorkspaceEmailSettings(ctx context.Context, settings *models.WorkspaceEmailSettings) error {
	return s.repo.Save(ctx, settings)
}
