package services

import (
	"context"
	"log"
	"time"

	"vessel.dev/vessel/internal/cloud/models"
	"vessel.dev/vessel/internal/cloud/repos"
)

type AuditEvent struct {
	ID        string
	TeamID    string
	UserID    string
	Action    string
	Resource  string
	IPAddress string
	Timestamp time.Time
}

type AuditService struct {
	repo repos.CloudRepo
}

func NewAuditService(repo repos.CloudRepo) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) LogEvent(ctx context.Context, event AuditEvent) error {
	event.Timestamp = time.Now()

	log.Printf("[AUDIT] Team: %s | User: %s | Action: %s | Resource: %s",
		event.TeamID, event.UserID, event.Action, event.Resource)

	return s.repo.InsertAuditLog(ctx, &models.AuditLog{
		TeamID:    event.TeamID,
		UserID:    event.UserID,
		Action:    event.Action,
		Resource:  event.Resource,
		IPAddress: event.IPAddress,
		Timestamp: event.Timestamp,
	})
}
