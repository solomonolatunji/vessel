package services

import (
	"context"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type AISettingsService struct {
	repo repositories.AISettingsRepository
}

func NewAISettingsService(repo repositories.AISettingsRepository) *AISettingsService {
	return &AISettingsService{repo: repo}
}

func (s *AISettingsService) GetAISettings(ctx context.Context) (*models.AISettings, error) {
	return s.repo.GetAISettings(ctx)
}

func (s *AISettingsService) UpdateAISettings(ctx context.Context, cfg *models.AISettings) error {
	return s.repo.UpdateAISettings(ctx, cfg)
}
