package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type BackupService struct {
	backupRepo repositories.BackupRepository
	s3Repo     repositories.S3DestinationRepository
	manager    *engine.BackupManager
}

func NewBackupService(br repositories.BackupRepository, sr repositories.S3DestinationRepository, m *engine.BackupManager) *BackupService {
	return &BackupService{
		backupRepo: br,
		s3Repo:     sr,
		manager:    m,
	}
}

func (s *BackupService) CreateConfig(ctx context.Context, cfg *models.BackupConfig) error {
	if cfg == nil || cfg.ProjectID == "" {
		return errors.New("valid backup config with projectId required")
	}
	if cfg.ID == "" {
		cfg.ID = uuid.New().String()
	}
	if cfg.Schedule == "" {
		cfg.Schedule = "0 2 * * *"
	}
	if cfg.RetentionDays == 0 {
		cfg.RetentionDays = 7
	}
	cfg.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	cfg.UpdatedAt = cfg.CreatedAt
	if err := s.backupRepo.CreateConfig(ctx, cfg); err != nil {
		return err
	}
	if s.manager != nil {
		_ = s.manager.RegisterBackup(cfg)
	}
	return nil
}

func (s *BackupService) GetConfig(ctx context.Context, id string) (*models.BackupConfig, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.backupRepo.GetConfigByID(ctx, id)
}

func (s *BackupService) ListConfigsByProject(ctx context.Context, projectID string) ([]*models.BackupConfig, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	return s.backupRepo.ListConfigsByProject(ctx, projectID)
}

func (s *BackupService) DeleteConfig(ctx context.Context, id, projectID string) error {
	if id == "" || projectID == "" {
		return errors.New("id and projectId required")
	}
	if s.manager != nil {
		s.manager.UnregisterBackup(id)
	}
	return s.backupRepo.DeleteConfig(ctx, id, projectID)
}

func (s *BackupService) CreateS3Destination(ctx context.Context, dest *models.S3Destination) error {
	if dest == nil || dest.ProjectID == "" || dest.Bucket == "" {
		return errors.New("valid s3 destination required")
	}
	if dest.ID == "" {
		dest.ID = uuid.New().String()
	}
	dest.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := engine.EnsureS3Bucket(ctx, dest); err != nil {
		return fmt.Errorf("failed to verify or create bucket: %w", err)
	}

	return s.s3Repo.CreateS3Destination(ctx, dest)
}

func (s *BackupService) ListS3Destinations(ctx context.Context, projectID string) ([]*models.S3Destination, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	return s.s3Repo.ListS3Destinations(ctx, projectID)
}

func (s *BackupService) DeleteS3Destination(ctx context.Context, id, projectID string) error {
	if id == "" || projectID == "" {
		return errors.New("id and projectId required")
	}
	return s.s3Repo.DeleteS3Destination(ctx, id, projectID)
}

func (s *BackupService) TriggerBackup(ctx context.Context, configID string) (*models.BackupRecord, error) {
	if s.manager == nil {
		return nil, errors.New("backup manager not available")
	}
	return s.manager.TriggerBackup(ctx, configID)
}

func (s *BackupService) ListRecordsByConfig(ctx context.Context, configID string) ([]*models.BackupRecord, error) {
	if configID == "" {
		return nil, errors.New("config id required")
	}
	return s.backupRepo.ListRecordsByConfig(ctx, configID)
}

func (s *BackupService) RestoreBackup(ctx context.Context, recordID string) error {
	if s.manager == nil {
		return errors.New("backup manager not available")
	}
	return s.manager.RestoreBackup(ctx, recordID)
}
