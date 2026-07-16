package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type ProjectSettingsService struct {
	repo        repositories.ProjectSettingsRepository
	userRepo    repositories.UserRepository
	authService *AuthService
}

func NewProjectSettingsService(r repositories.ProjectSettingsRepository, ur repositories.UserRepository, authService *AuthService) *ProjectSettingsService {
	return &ProjectSettingsService{
		repo:        r,
		userRepo:    ur,
		authService: authService,
	}
}

func (s *ProjectSettingsService) CreateWebhook(ctx context.Context, w *models.Webhook) (*models.Webhook, error) {
	if w == nil || w.ProjectID == "" || w.URL == "" {
		return nil, errors.New("valid webhook with projectId and url required")
	}
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	w.CreatedAt = time.Now()
	if err := s.repo.CreateWebhook(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *ProjectSettingsService) ListWebhooks(ctx context.Context, projectID string) ([]*models.Webhook, error) {
	if projectID == "" {
		return nil, errors.New("projectId required")
	}
	return s.repo.ListWebhooksByProject(ctx, projectID)
}

func (s *ProjectSettingsService) DeleteWebhook(ctx context.Context, id, projectID string) error {
	if id == "" || projectID == "" {
		return errors.New("id and projectId required")
	}
	return s.repo.DeleteWebhook(ctx, id, projectID)
}

func (s *ProjectSettingsService) CreateToken(ctx context.Context, t *models.ProjectToken) (*models.ProjectToken, string, error) {
	if t == nil || t.ProjectID == "" || t.Name == "" {
		return nil, "", errors.New("valid token with projectId and name required")
	}
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	t.CreatedAt = time.Now().UTC()

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, "", fmt.Errorf("generate token bytes: %w", err)
	}
	rawSecret := hex.EncodeToString(randomBytes)
	fullToken := fmt.Sprintf("vsl_tok_%s", rawSecret)
	t.TokenPrefix = fullToken[:16]

	err := s.repo.CreateToken(ctx, t, fullToken)
	if err != nil {
		return nil, "", err
	}
	return t, fullToken, nil
}

func (s *ProjectSettingsService) GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error) {
	return s.repo.GetTokenByHash(ctx, tokenHash)
}

func (s *ProjectSettingsService) UpdateTokenLastUsed(ctx context.Context, id string) error {
	return s.repo.UpdateTokenLastUsed(ctx, id)
}

func (s *ProjectSettingsService) ListTokens(ctx context.Context, projectID string) ([]*models.ProjectToken, error) {
	if projectID == "" {
		return nil, errors.New("projectId required")
	}
	return s.repo.ListTokensByProject(ctx, projectID)
}

func (s *ProjectSettingsService) DeleteToken(ctx context.Context, id, projectID string) error {
	if id == "" || projectID == "" {
		return errors.New("id and projectId required")
	}
	return s.repo.DeleteToken(ctx, id, projectID)
}

func (s *ProjectSettingsService) AddMemberByEmail(ctx context.Context, projectID, email string, permission models.MemberPermission, originUrl string) (*models.ProjectMember, error) {
	if projectID == "" || email == "" {
		return nil, errors.New("valid member with projectId and email required")
	}

	u, err := s.authService.InviteUser(ctx, email, originUrl)
	if err != nil {
		return nil, err
	}

	m := &models.ProjectMember{
		ID:         uuid.New().String(),
		ProjectID:  projectID,
		UserID:     u.ID,
		Email:      email,
		Permission: permission,
		Status:     models.MemberStatusPending,
		InvitedAt:  time.Now(),
	}

	if err := s.repo.AddMember(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *ProjectSettingsService) ListMembers(ctx context.Context, projectID string) ([]map[string]any, error) {
	if projectID == "" {
		return nil, errors.New("projectId required")
	}
	members, err := s.repo.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var res []map[string]any
	for _, m := range members {
		item := map[string]any{
			"id":         m.ID,
			"projectId":  m.ProjectID,
			"userId":     m.UserID,
			"permission": m.Permission,
			"invitedAt":  m.InvitedAt,
		}
		if s.userRepo != nil {
			if u, err := s.userRepo.GetUserByID(ctx, m.UserID); err == nil && u != nil {
				item["email"] = u.Email
			}
		}
		res = append(res, item)
	}
	return res, nil
}

func (s *ProjectSettingsService) RemoveMember(ctx context.Context, id, projectID string) error {
	if id == "" || projectID == "" {
		return errors.New("id and projectId required")
	}
	return s.repo.RemoveMember(ctx, id, projectID)
}
