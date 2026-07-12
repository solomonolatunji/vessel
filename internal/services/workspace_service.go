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
	repo     repositories.WorkspaceRepository
	userRepo repositories.UserRepository
}

func NewWorkspaceService(r repositories.WorkspaceRepository, userRepo repositories.UserRepository) *WorkspaceService {
	return &WorkspaceService{repo: r, userRepo: userRepo}
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

func (s *WorkspaceService) ListWorkspaces(ctx context.Context, ownerID string, limit, offset int) ([]*models.Workspace, int, error) {
	if ownerID == "" {
		return nil, 0, errors.New("owner id required")
	}
	return s.repo.List(ctx, ownerID, limit, offset)
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

func (s *WorkspaceService) AddTrustedDomain(ctx context.Context, workspaceID, domain string) (*models.TrustedDomain, error) {
	if workspaceID == "" || domain == "" {
		return nil, errors.New("teamId and domain required")
	}
	td := &models.TrustedDomain{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Domain:      domain,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.CreateTrustedDomain(ctx, td); err != nil {
		return nil, err
	}
	return td, nil
}

func (s *WorkspaceService) ListTrustedDomains(ctx context.Context, workspaceID string) ([]*models.TrustedDomain, error) {
	if workspaceID == "" {
		return nil, errors.New("teamId required")
	}
	return s.repo.ListTrustedDomains(ctx, workspaceID)
}

func (s *WorkspaceService) DeleteTrustedDomain(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.DeleteTrustedDomain(ctx, id)
}

func (s *WorkspaceService) AddSSHKey(ctx context.Context, workspaceID, name, pubKey string) (*models.SSHKey, error) {
	if workspaceID == "" || name == "" || pubKey == "" {
		return nil, errors.New("teamId, name, and publicKey required")
	}
	key := &models.SSHKey{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Name:        name,
		PublicKey:   pubKey,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.CreateSSHKey(ctx, key); err != nil {
		return nil, err
	}
	return key, nil
}

func (s *WorkspaceService) ListSSHKeys(ctx context.Context, workspaceID string) ([]*models.SSHKey, error) {
	if workspaceID == "" {
		return nil, errors.New("teamId required")
	}
	return s.repo.ListSSHKeys(ctx, workspaceID)
}

func (s *WorkspaceService) DeleteSSHKey(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.DeleteSSHKey(ctx, id)
}

func (s *WorkspaceService) ListAuditLogs(ctx context.Context, workspaceID string, limit, offset int) ([]*models.AuditLog, int, error) {
	if workspaceID == "" {
		return nil, 0, errors.New("workspaceId required")
	}
	return s.repo.ListAuditLogs(ctx, workspaceID, limit, offset)
}

func (s *WorkspaceService) ListWorkspacesByUser(ctx context.Context, userID string) ([]*models.Workspace, error) {
	if userID == "" {
		return nil, errors.New("userId is required")
	}
	return s.repo.ListWorkspacesByUser(ctx, userID)
}

func (s *WorkspaceService) AddMember(ctx context.Context, workspaceID, userID, role string) error {
	if workspaceID == "" || userID == "" {
		return errors.New("workspaceId and userId are required")
	}
	if role == "" {
		role = "member"
	}
	member := &models.WorkspaceMember{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		JoinedAt:    time.Now(),
	}
	return s.repo.AddMember(ctx, member)
}

func (s *WorkspaceService) RemoveMember(ctx context.Context, workspaceID, userID string) error {
	if workspaceID == "" || userID == "" {
		return errors.New("workspaceId and userId are required")
	}
	return s.repo.RemoveMember(ctx, workspaceID, userID)
}

func (s *WorkspaceService) ListMembers(ctx context.Context, workspaceID string) ([]*models.WorkspaceMember, error) {
	if workspaceID == "" {
		return nil, errors.New("workspaceId is required")
	}
	return s.repo.ListMembers(ctx, workspaceID)
}

func (s *WorkspaceService) InviteMember(ctx context.Context, workspaceID, email, role string) (*models.WorkspaceInvite, error) {
	if workspaceID == "" || email == "" {
		return nil, errors.New("workspaceId and email are required")
	}
	if role == "" {
		role = "member"
	}
	invite := &models.WorkspaceInvite{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Email:       email,
		Role:        role,
		Token:       uuid.New().String(),
		CreatedAt:   time.Now(),
	}
	if err := s.repo.CreateInvite(ctx, invite); err != nil {
		return nil, err
	}
	return invite, nil
}

func (s *WorkspaceService) GetInvite(ctx context.Context, token string) (*models.WorkspaceInvite, error) {
	if token == "" {
		return nil, errors.New("token required")
	}
	return s.repo.GetInviteByToken(ctx, token)
}

func (s *WorkspaceService) AcceptInvite(ctx context.Context, token, userID string) error {
	if token == "" || userID == "" {
		return errors.New("token and userId are required")
	}
	invite, err := s.repo.GetInviteByToken(ctx, token)
	if err != nil || invite == nil {
		return errors.New("invalid or expired invite token")
	}
	if err := s.AddMember(ctx, invite.WorkspaceID, userID, invite.Role); err != nil {
		return err
	}
	return s.repo.DeleteInvite(ctx, invite.ID)
}
