package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/utils"
)

type StorageService struct {
	repo     repositories.StorageRepository
	deployer *engine.StorageDeployer
}

func NewStorageService(r repositories.StorageRepository, d *engine.StorageDeployer) *StorageService {
	return &StorageService{
		repo:     r,
		deployer: d,
	}
}

func (s *StorageService) CreateStorage(ctx context.Context, st *models.Storage) (*models.Storage, error) {
	if st == nil || st.Name == "" {
		return nil, errors.New("storage name required")
	}
	if st.ID == "" {
		st.ID = uuid.New().String()
	}
	if st.Status == "" {
		st.Status = "stopped"
	}
	now := time.Now()
	if st.CreatedAt.IsZero() {
		st.CreatedAt = now
	}
	st.UpdatedAt = now
	if err := s.repo.Create(ctx, st); err != nil {
		return nil, err
	}
	if s.deployer != nil {
		containerID, err := s.deployer.SpinUp(ctx, st)
		if err == nil && containerID != "" {
			st.ContainerID = containerID
			st.Status = "running"
			_ = s.repo.Update(ctx, st)
		} else if err != nil {
			st.Status = "error"
			_ = s.repo.Update(ctx, st)
		}
	}
	return st, nil
}

func (s *StorageService) GetStorage(ctx context.Context, id string) (*models.Storage, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *StorageService) ListStorage(ctx context.Context) ([]*models.Storage, error) {
	return s.repo.List(ctx)
}

func (s *StorageService) ListStorageByProject(ctx context.Context, projectID string) ([]*models.Storage, error) {
	if projectID == "" {
		return nil, errors.New("project id is required")
	}
	return s.repo.ListByProject(ctx, projectID)
}

func (s *StorageService) UpdateStorage(ctx context.Context, st *models.Storage) error {
	if st == nil || st.ID == "" {
		return errors.New("valid storage required for update")
	}
	st.UpdatedAt = time.Now()
	return s.repo.Update(ctx, st)
}

func (s *StorageService) DeleteStorage(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if s.deployer != nil {
		_ = s.deployer.Stop(ctx, id)
	}
	return s.repo.Delete(ctx, id)
}

func (s *StorageService) CreateStorageWithDefaults(ctx context.Context, st *models.Storage) (*models.Storage, error) {
	if st == nil || st.Name == "" {
		return nil, errors.New("storage name is required")
	}
	if st.APIPort <= 0 {
		st.APIPort = 9000
	}
	if st.ConsolePort <= 0 {
		st.ConsolePort = 9001
	}
	if st.AccessKey == "" {
		st.AccessKey = uuid.New().String()[:16]
	}
	if st.SecretKey == "" {
		st.SecretKey = uuid.New().String()[:24]
	}
	if st.BucketName == "" {
		st.BucketName = "vessl-backups"
	}
	if st.Type == "" {
		st.Type = "minio"
	}
	return s.CreateStorage(ctx, st)
}

func (s *StorageService) StartStorage(ctx context.Context, id string) (*models.Storage, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	st, err := s.repo.GetByID(ctx, id)
	if err != nil || st == nil {
		return nil, utils.NewNotFoundError("Storage", id)
	}
	if s.deployer == nil {
		return nil, errors.New("storage deployer unavailable")
	}
	containerID, err := s.deployer.SpinUp(ctx, st)
	if err != nil {
		st.Status = "error"
		_ = s.repo.Update(ctx, st)
		return nil, err
	}
	if containerID != "" {
		st.ContainerID = containerID
	}
	st.Status = "running"
	st.UpdatedAt = time.Now()
	_ = s.repo.Update(ctx, st)
	return st, nil
}

func (s *StorageService) StopStorage(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	st, err := s.repo.GetByID(ctx, id)
	if err != nil || st == nil {
		return errors.New("storage record not found")
	}
	if s.deployer != nil {
		_ = s.deployer.Stop(ctx, id)
	}
	st.Status = "stopped"
	st.UpdatedAt = time.Now()
	return s.repo.Update(ctx, st)
}
