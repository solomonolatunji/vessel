package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type AISettingsService struct {
	repo repositories.WorkspaceAISettingsRepository
}

func NewAISettingsService(repo repositories.WorkspaceAISettingsRepository) *AISettingsService {
	return &AISettingsService{repo: repo}
}

func (s *AISettingsService) Get(ctx context.Context, workspaceID string) (*models.WorkspaceAISettings, error) {
	if workspaceID == "" {
		return nil, errors.New("team ID is required")
	}
	return s.repo.Get(ctx, workspaceID)
}

func (s *AISettingsService) Save(ctx context.Context, settings *models.WorkspaceAISettings) error {
	if settings == nil {
		return errors.New("settings are required")
	}
	if settings.WorkspaceID == "" {
		return errors.New("team ID is required")
	}
	if settings.ID == "" {
		settings.ID = uuid.NewString()
	}
	return s.repo.Save(ctx, settings)
}
