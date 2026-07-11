package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type AISettingsService struct {
	repo repositories.TeamAISettingsRepository
}

func NewAISettingsService(repo repositories.TeamAISettingsRepository) *AISettingsService {
	return &AISettingsService{repo: repo}
}

func (s *AISettingsService) Get(ctx context.Context, teamID string) (*models.TeamAISettings, error) {
	if teamID == "" {
		return nil, errors.New("team ID is required")
	}
	return s.repo.Get(ctx, teamID)
}

func (s *AISettingsService) Save(ctx context.Context, settings *models.TeamAISettings) error {
	if settings == nil {
		return errors.New("settings are required")
	}
	if settings.TeamID == "" {
		return errors.New("team ID is required")
	}
	if settings.ID == "" {
		settings.ID = uuid.NewString()
	}
	return s.repo.Save(ctx, settings)
}
