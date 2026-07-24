package services

import (
	"context"
	"errors"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type ServerlessService interface {
	SaveCode(ctx context.Context, serviceID, runtime, codeContent string) (*models.ServerlessFunctionCode, error)
	GetCode(ctx context.Context, serviceID string) (*models.ServerlessFunctionCode, error)
}

type serverlessService struct {
	repo repositories.ServerlessRepository
}

func NewServerlessService(repo repositories.ServerlessRepository) ServerlessService {
	return &serverlessService{repo: repo}
}

func (s *serverlessService) SaveCode(ctx context.Context, serviceID, runtime, codeContent string) (*models.ServerlessFunctionCode, error) {
	if serviceID == "" || runtime == "" || codeContent == "" {
		return nil, errors.New("serviceID, runtime, and codeContent are required")
	}
	return s.repo.SaveCode(ctx, serviceID, runtime, codeContent)
}

func (s *serverlessService) GetCode(ctx context.Context, serviceID string) (*models.ServerlessFunctionCode, error) {
	if serviceID == "" {
		return nil, errors.New("serviceID is required")
	}
	return s.repo.GetCodeByServiceID(ctx, serviceID)
}
