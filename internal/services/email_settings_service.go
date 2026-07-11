package services

import (
	"context"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type EmailSettingsService struct {
	repo repositories.TeamEmailSettingsRepository
}

func NewEmailSettingsService(repo repositories.TeamEmailSettingsRepository) *EmailSettingsService {
	return &EmailSettingsService{repo: repo}
}

func (s *EmailSettingsService) GetTeamEmailSettings(ctx context.Context, teamID string) (*models.TeamEmailSettings, error) {
	return s.repo.GetByTeamID(ctx, teamID)
}

func (s *EmailSettingsService) SaveTeamEmailSettings(ctx context.Context, settings *models.TeamEmailSettings) error {
	return s.repo.Save(ctx, settings)
}
