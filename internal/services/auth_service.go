package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type AuthService struct {
	userRepo     repositories.UserRepository
	settingsRepo repositories.SettingsRepository
	tokenService *TokenService
}

func NewAuthService(ur repositories.UserRepository, sr repositories.SettingsRepository, ts *TokenService) *AuthService {
	return &AuthService{
		userRepo:     ur,
		settingsRepo: sr,
		tokenService: ts,
	}
}

func (a *AuthService) Register(ctx context.Context, email, password string) (*models.User, string, error) {
	if email == "" || password == "" {
		return nil, "", errors.New("email and password are required")
	}
	users, _ := a.userRepo.ListUsers(ctx)
	isInitial := len(users) == 0
	cfg, _ := a.settingsRepo.GetServerSettings(ctx)
	if cfg != nil && !cfg.RegistrationEnabled && !isInitial {
		return nil, "", errors.New("user registration is disabled on this server")
	}
	if cfg != nil && !isInitial && strings.TrimSpace(cfg.RegistrationDomainAllowlist) != "" {
		allowed := false
		for _, d := range strings.Split(cfg.RegistrationDomainAllowlist, ",") {
			d = strings.TrimSpace(d)
			if d != "" && strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(d)) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, "", errors.New("email domain is not allowed on this server")
		}
	}
	existing, _ := a.userRepo.GetUserByEmail(ctx, email)
	if existing != nil {
		return nil, "", errors.New("user already exists with that email")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	role := "member"
	if isInitial {
		role = "admin"
	}
	u := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashed),
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := a.userRepo.CreateUser(ctx, u); err != nil {
		return nil, "", err
	}
	token, err := a.tokenService.GenerateToken(u)
	if err != nil {
		return nil, "", err
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	return &uCopy, token, nil
}

func (a *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	if email == "" || password == "" {
		return nil, "", errors.New("email and password are required")
	}
	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return nil, "", errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}
	token, err := a.tokenService.GenerateToken(u)
	if err != nil {
		return nil, "", err
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	return &uCopy, token, nil
}
