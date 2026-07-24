package services

import (
	"context"
	"errors"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type CanvasService struct {
	repo repositories.CanvasRepository
}

func NewCanvasService(r repositories.CanvasRepository) *CanvasService {
	return &CanvasService{repo: r}
}

func (s *CanvasService) ListSummaries(ctx context.Context) ([]models.CanvasSummary, error) {
	return s.repo.ListCanvasSummaries(ctx)
}

func (s *CanvasService) GetSummary(ctx context.Context, id string) (*models.CanvasSummary, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetCanvasSummary(ctx, id)
}

func (s *CanvasService) GetEnvironmentCanvas(ctx context.Context, id string) (*models.EnvironmentCanvas, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetEnvironmentCanvas(ctx, id)
}
