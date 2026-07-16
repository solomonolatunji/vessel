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

type AuditActionOpts struct {
	UserID    string
	Action    string
	Resource  string
	IPAddress string
	Details   any
}

func (s *AuditService) LogAction(ctx context.Context, opts AuditActionOpts) {
	var detailsStr string
	if opts.Details != nil {
		b, err := json.Marshal(opts.Details)
		if err == nil {
			detailsStr = string(b)
		}
	}

	log := &models.AuditLog{
		ID:        uuid.New().String(),
		UserID:    opts.UserID,
		Action:    opts.Action,
		Resource:  opts.Resource,
		Details:   detailsStr,
		IPAddress: opts.IPAddress,
	}

	// We log asynchronously so we don't block the caller
	go func() {
		err := s.repo.Create(context.Background(), log)
		if err != nil {
			slog.Error("failed to write audit log", "err", err, "action", opts.Action)
		}
	}()
}

func (s *AuditService) ListLogs(ctx context.Context, limit, offset int) ([]models.AuditLog, error) {
	return s.repo.List(ctx, limit, offset)
}
