package services

import (
	"context"
	"errors"

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

func (s *NotificationService) GetIntegration(ctx context.Context) (*models.NotificationIntegration, error) {
	return s.repo.GetIntegration(ctx)
}

func (s *NotificationService) SaveIntegration(ctx context.Context, n *models.NotificationIntegration) error {
	if n == nil {
		return errors.New("integration required")
	}
	return s.repo.SaveIntegration(ctx, n)
}

func (s *NotificationService) GetProjectPref(ctx context.Context, projectID string) (*models.ProjectNotificationPref, error) {
	if projectID == "" {
		return nil, errors.New("projectId required")
	}
	return s.repo.GetProjectPref(ctx, projectID)
}

func (s *NotificationService) SaveProjectPref(ctx context.Context, pref *models.ProjectNotificationPref) error {
	if pref == nil || pref.ProjectID == "" {
		return errors.New("valid preference with projectId required")
	}
	return s.repo.SaveProjectPref(ctx, pref)
}

func (s *NotificationService) SendTest(channel, projectID string) error {
	if s.dispatcher == nil {
		return errors.New("dispatcher unavailable")
	}
	return s.dispatcher.Send(&models.NotificationEvent{
		Title:     "Test Notification from Vessel",
		Message:   "If you see this, your " + channel + " integration is working correctly!",
		Level:     "info",
		ProjectID: projectID,
		URL:       "http://localhost:3000/settings/notifications",
	})
}
