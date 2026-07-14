package services

import (
	"context"
	"errors"
	"os"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type NotificationDispatcher interface {
	Send(event *models.NotificationEvent) error
}

type NotificationService struct {
	repo       repositories.NotificationRepository
	dispatcher NotificationDispatcher
}

func NewNotificationService(repo repositories.NotificationRepository, dispatcher NotificationDispatcher) *NotificationService {
	return &NotificationService{repo: repo, dispatcher: dispatcher}
}

func (s *NotificationService) ListChannels(ctx context.Context, workspaceID string) ([]models.WorkspaceNotificationChannel, error) {
	if workspaceID == "" {
		return nil, errors.New("workspaceId required")
	}
	return s.repo.ListChannelsByTeam(ctx, workspaceID)
}

func (s *NotificationService) GetChannel(ctx context.Context, id string) (*models.WorkspaceNotificationChannel, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetChannel(ctx, id)
}

func (s *NotificationService) SaveChannel(ctx context.Context, c *models.WorkspaceNotificationChannel) error {
	if c == nil || c.WorkspaceID == "" {
		return errors.New("valid channel with workspaceId required")
	}
	return s.repo.SaveChannel(ctx, c)
}

func (s *NotificationService) DeleteChannel(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.DeleteChannel(ctx, id)
}

func (s *NotificationService) TestGlobalNotification(ctx context.Context, provider string) error {
	if s.dispatcher == nil {
		return errors.New("dispatcher unavailable")
	}
	dashboardURL := os.Getenv("VESSL_DASHBOARD_URL")
	return s.dispatcher.Send(&models.NotificationEvent{
		WorkspaceID: "global_test",
		EventType:   "test_global_" + provider,
		Title:       "Global Test Notification from Vessl",
		Message:     "If you see this, your global integration is working correctly!",
		Level:       "info",
		URL:         dashboardURL + "/settings/notifications",
	})
}

func (s *NotificationService) TestTeamNotification(ctx context.Context, workspaceID, channelID string) error {
	if s.dispatcher == nil {
		return errors.New("dispatcher unavailable")
	}
	dashboardURL := os.Getenv("VESSL_DASHBOARD_URL")
	return s.dispatcher.Send(&models.NotificationEvent{
		WorkspaceID: workspaceID,
		EventType:   "test_channel_" + channelID,
		Title:       "Test Notification from Vessl",
		Message:     "If you see this, your team integration is working correctly!",
		Level:       "info",
		URL:         dashboardURL + "/settings/notifications",
	})
}
