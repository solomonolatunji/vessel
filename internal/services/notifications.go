package services

import (
	"context"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type NotificationSettingsService struct {
	repo repositories.NotificationSettingsRepository
}

func NewNotificationSettingsService(repo repositories.NotificationSettingsRepository) *NotificationSettingsService {
	return &NotificationSettingsService{repo: repo}
}

func (s *NotificationSettingsService) GetNotificationSettings(ctx context.Context) (*models.NotificationSettings, error) {
	return s.repo.GetNotificationSettings(ctx)
}

func (s *NotificationSettingsService) UpdateNotificationSettings(ctx context.Context, cfg *models.NotificationSettings) error {
	return s.repo.UpdateNotificationSettings(ctx, cfg)
}
