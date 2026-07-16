package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type AuditService struct {
	repo repositories.AuditLogRepository
}

func NewAuditService(repo repositories.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) LogAction(ctx context.Context, userID, action, resource, ipAddress string, details any) {
	var detailsStr string
	if details != nil {
		b, err := json.Marshal(details)
		if err == nil {
			detailsStr = string(b)
		}
	}

	log := &models.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Details:   detailsStr,
		IPAddress: ipAddress,
	}

	// We log asynchronously so we don't block the caller
	go func() {
		err := s.repo.Create(context.Background(), log)
		if err != nil {
			slog.Error("failed to write audit log", "err", err, "action", action)
		}
	}()
}

func (s *AuditService) ListLogs(ctx context.Context, limit, offset int) ([]models.AuditLog, error) {
	return s.repo.List(ctx, limit, offset)
}
