package services

import (
	"context"
	"errors"
	"os"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
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

func (s *NotificationService) ListChannels(ctx context.Context, teamID string) ([]models.TeamNotificationChannel, error) {
	if teamID == "" {
		return nil, errors.New("teamId required")
	}
	return s.repo.ListChannelsByTeam(ctx, teamID)
}

func (s *NotificationService) GetChannel(ctx context.Context, id string) (*models.TeamNotificationChannel, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetChannel(ctx, id)
}

func (s *NotificationService) SaveChannel(ctx context.Context, c *models.TeamNotificationChannel) error {
	if c == nil || c.TeamID == "" {
		return errors.New("valid channel with teamId required")
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
	dashboardURL := os.Getenv("VESSEL_DASHBOARD_URL")
	if dashboardURL == "" {
		dashboardURL = "http://localhost:3000"
	}
	return s.dispatcher.Send(&models.NotificationEvent{
		TeamID:    "global_test",
		EventType: "test_global_" + provider,
		Title:     "Global Test Notification from Vessel",
		Message:   "If you see this, your global integration is working correctly!",
		Level:     "info",
		URL:       dashboardURL + "/settings/notifications",
	})
}

func (s *NotificationService) TestTeamNotification(ctx context.Context, teamID, channelID string) error {
	if s.dispatcher == nil {
		return errors.New("dispatcher unavailable")
	}
	dashboardURL := os.Getenv("VESSEL_DASHBOARD_URL")
	if dashboardURL == "" {
		dashboardURL = "http://localhost:3000"
	}
	return s.dispatcher.Send(&models.NotificationEvent{
		TeamID:    teamID,
		EventType: "test_channel_" + channelID,
		Title:     "Test Notification from Vessel",
		Message:   "If you see this, your team integration is working correctly!",
		Level:     "info",
		URL:       dashboardURL + "/settings/notifications",
	})
}
