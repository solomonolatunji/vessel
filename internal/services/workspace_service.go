package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type WorkspaceService struct {
	repo repositories.WorkspaceRepository
}

func NewWorkspaceService(r repositories.WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{repo: r}
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, name, ownerID string) (*models.Workspace, error) {
	if name == "" || ownerID == "" {
		return nil, errors.New("workspace name and ownerId required")
	}
	ws := &models.Workspace{
		ID:        uuid.New().String(),
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, ws); err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, id string) (*models.Workspace, error) {
	if id == "" {
		return nil, errors.New("workspace id required")
	}
	return s.repo.Get(ctx, id)
}

func (s *WorkspaceService) ListWorkspaces(ctx context.Context, ownerID string) ([]*models.Workspace, error) {
	if ownerID == "" {
		return nil, errors.New("owner id required")
	}
	return s.repo.List(ctx, ownerID)
}

func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, ws *models.Workspace) error {
	if ws == nil || ws.ID == "" {
		return errors.New("valid workspace required for update")
	}
	ws.UpdatedAt = time.Now()
	return s.repo.Update(ctx, ws)
}

func (s *WorkspaceService) DeleteWorkspace(ctx context.Context, id, ownerID string) error {
	if id == "" || ownerID == "" {
		return errors.New("workspace id and owner id required")
	}
	return s.repo.Delete(ctx, id, ownerID)
}

func (s *WorkspaceService) AddTrustedDomain(ctx context.Context, teamID, domain string) (*models.TrustedDomain, error) {
	if teamID == "" || domain == "" {
		return nil, errors.New("teamId and domain required")
	}
	td := &models.TrustedDomain{
		ID:        uuid.New().String(),
		TeamID:    teamID,
		Domain:    domain,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateTrustedDomain(ctx, td); err != nil {
		return nil, err
	}
	return td, nil
}

func (s *WorkspaceService) ListTrustedDomains(ctx context.Context, teamID string) ([]*models.TrustedDomain, error) {
	if teamID == "" {
		return nil, errors.New("teamId required")
	}
	return s.repo.ListTrustedDomains(ctx, teamID)
}

func (s *WorkspaceService) DeleteTrustedDomain(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.DeleteTrustedDomain(ctx, id)
}

func (s *WorkspaceService) AddSSHKey(ctx context.Context, teamID, name, pubKey string) (*models.SSHKey, error) {
	if teamID == "" || name == "" || pubKey == "" {
		return nil, errors.New("teamId, name, and publicKey required")
	}
	key := &models.SSHKey{
		ID:        uuid.New().String(),
		TeamID:    teamID,
		Name:      name,
		PublicKey: pubKey,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateSSHKey(ctx, key); err != nil {
		return nil, err
	}
	return key, nil
}

func (s *WorkspaceService) ListSSHKeys(ctx context.Context, teamID string) ([]*models.SSHKey, error) {
	if teamID == "" {
		return nil, errors.New("teamId required")
	}
	return s.repo.ListSSHKeys(ctx, teamID)
}

func (s *WorkspaceService) DeleteSSHKey(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.DeleteSSHKey(ctx, id)
}

func (s *WorkspaceService) ListAuditLogs(ctx context.Context, teamID string, limit int) ([]*models.AuditLog, error) {
	if teamID == "" {
		return nil, errors.New("teamId required")
	}
	return s.repo.ListAuditLogs(ctx, teamID, limit)
}
