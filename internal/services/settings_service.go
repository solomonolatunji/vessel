package services

import (
	"context"
	"errors"
	"fmt"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/repositories"
)

type SettingsService struct {
	settingsRepo repositories.SettingsRepository
}

func NewSettingsService(sr repositories.SettingsRepository) *SettingsService {
	return &SettingsService{
		settingsRepo: sr,
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
			{"type": "text", "text": "Codedock system is healthy and operational."},
		}, nil
	default:
		return nil, fmt.Errorf("Method/Tool not found: %s", toolName)
	}
}
