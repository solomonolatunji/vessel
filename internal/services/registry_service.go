package services

import (
	"context"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type RegistryService interface {
	CreateRegistry(ctx context.Context, registry *models.Registry) error
	ListRegistriesByProject(ctx context.Context, projectID string) ([]*models.Registry, error)
	DeleteRegistry(ctx context.Context, id string) error
}

type registryService struct {
	repo repositories.RegistryRepository
}

func NewRegistryService(repo repositories.RegistryRepository) RegistryService {
	return &registryService{repo: repo}
}

func (s *registryService) CreateRegistry(ctx context.Context, registry *models.Registry) error {
	return s.repo.Create(ctx, registry)
}

func (s *registryService) ListRegistriesByProject(ctx context.Context, projectID string) ([]*models.Registry, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *registryService) DeleteRegistry(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
