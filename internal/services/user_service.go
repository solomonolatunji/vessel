package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(ur repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

func (s *UserService) CreateUser(ctx context.Context, u *models.User) error {
	if u.Email == "" {
		return errors.New("user email is required")
	}
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	return s.userRepo.CreateUser(ctx, u)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.ListUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, u *models.User) error {
	if u == nil || u.ID == "" {
		return errors.New("valid user required for update")
	}
	u.UpdatedAt = time.Now()
	return s.userRepo.UpdateUser(ctx, u)
}

func (s *UserService) CreatePAT(ctx context.Context, userID, name string, expiresAt *time.Time) (*models.PersonalAccessToken, string, error) {
	if userID == "" || name == "" {
		return nil, "", errors.New("userId and name are required")
	}
	rawToken := uuid.New().String() + "-" + uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	pat := &models.PersonalAccessToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		TokenHash: string(hash),
		CreatedAt: time.Now(),
	}
	if expiresAt != nil {
		pat.ExpiresAt = *expiresAt
	}
	if err := s.userRepo.CreatePAT(ctx, pat); err != nil {
		return nil, "", err
	}
	return pat, rawToken, nil
}

func (s *UserService) ListPATs(ctx context.Context, userID string) ([]*models.PersonalAccessToken, error) {
	if userID == "" {
		return nil, errors.New("userId is required")
	}
	return s.userRepo.ListPATs(ctx, userID)
}

func (s *UserService) DeletePAT(ctx context.Context, id, userID string) error {
	if id == "" || userID == "" {
		return errors.New("id and userId are required")
	}
	return s.userRepo.DeletePAT(ctx, id, userID)
}
