package services

import (
	"context"
	"errors"
	"fmt"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type SettingsService struct {
	settingsRepo     repositories.SettingsRepository
	notificationRepo repositories.NotificationRepository
}

func NewSettingsService(sr repositories.SettingsRepository, nr repositories.NotificationRepository) *SettingsService {
	return &SettingsService{
		settingsRepo:     sr,
		notificationRepo: nr,
	}
}

func (s *SettingsService) GetSettings(ctx context.Context) (*models.ServerSettings, error) {
	return s.settingsRepo.GetServerSettings(ctx)
}

func (s *SettingsService) UpdateSettings(ctx context.Context, cfg *models.ServerSettings) error {
	if cfg == nil {
		return errors.New("server settings cannot be nil")
	}
	return s.settingsRepo.UpdateServerSettings(ctx, cfg)
}

func (s *SettingsService) ListTeamNotificationChannels(ctx context.Context, teamID string) ([]models.TeamNotificationChannel, error) {
	if teamID == "" {
		return nil, errors.New("team id is required")
	}
	return s.notificationRepo.ListChannelsByTeam(ctx, teamID)
}

func (s *SettingsService) SaveTeamNotificationChannel(ctx context.Context, c *models.TeamNotificationChannel) error {
	if c == nil || c.TeamID == "" {
		return errors.New("valid team notification channel is required")
	}
	return s.notificationRepo.SaveChannel(ctx, c)
}

func (s *SettingsService) GetTeamNotificationChannel(ctx context.Context, id string) (*models.TeamNotificationChannel, error) {
	if id == "" {
		return nil, errors.New("channel id is required")
	}
	return s.notificationRepo.GetChannel(ctx, id)
}

func (s *SettingsService) DeleteTeamNotificationChannel(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("channel id is required")
	}
	return s.notificationRepo.DeleteChannel(ctx, id)
}

func (s *SettingsService) CheckMCPEnabled(ctx context.Context) error {
	settings, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil {
		return err
	}
	if settings != nil && !settings.MCPServerEnabled {
		return errors.New("MCP server endpoint is currently disabled by the administrator")
	}
	return nil
}

func (s *SettingsService) ExecuteMCPTool(ctx context.Context, toolName string) ([]map[string]any, error) {
	if err := s.CheckMCPEnabled(ctx); err != nil {
		return nil, err
	}
	switch toolName {
	case "list_projects":
		projects, err := s.settingsRepo.ListProjects(ctx)
		if err != nil {
			return nil, err
		}
		return []map[string]any{
			{"type": "text", "text": fmt.Sprintf("Found %d projects: %+v", len(projects), projects)},
		}, nil
	case "get_system_status":
		return []map[string]any{
			{"type": "text", "text": "Vessel system is healthy and operational."},
		}, nil
	default:
		return nil, fmt.Errorf("Method/Tool not found: %s", toolName)
	}
}
