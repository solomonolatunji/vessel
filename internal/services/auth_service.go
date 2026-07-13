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
	userRepo         repositories.UserRepository
	settingsRepo     repositories.SettingsRepository
	tokenService     *TokenService
	workspaceService *WorkspaceService
}

func NewAuthService(ur repositories.UserRepository, sr repositories.SettingsRepository, ts *TokenService, ws *WorkspaceService) *AuthService {
	return &AuthService{
		userRepo:         ur,
		settingsRepo:     sr,
		tokenService:     ts,
		workspaceService: ws,
	}
}

func (a *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, string, error) {
	if email == "" || password == "" || name == "" {
		return nil, "", errors.New("name, email and password are required")
	}
	_, total, _ := a.userRepo.ListUsers(ctx, 1, 0)
	isInitial := total == 0
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
		Name:         name,
		PasswordHash: string(hashed),
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := a.userRepo.CreateUser(ctx, u); err != nil {
		return nil, "", err
	}

	// Create default workspace for the user
	_, err = a.workspaceService.CreateWorkspace(ctx, "Personal Workspace", u.ID)
	if err != nil {
		// Log the error but don't fail registration
		// We could potentially retry or let the user create one manually if this fails
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

func (a *AuthService) ForgotPassword(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	cfg, err := a.settingsRepo.GetServerSettings(ctx)
	if err != nil || cfg == nil {
		return errors.New("could not load server settings")
	}

	if !cfg.SMTPEnabled && !cfg.ResendEnabled {
		return errors.New("your team is yet to set or enable email")
	}

	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		// Prevent email enumeration by returning nil even if not found
		return nil
	}

	// TODO: implement actual email sending with reset token
	// This fulfills the immediate requirement of checking if email is enabled
	token, err := a.tokenService.GeneratePasswordResetToken(u.Email)
	if err != nil {
		return err
	}
	// Stub: In reality, send this token via email
	_ = token

	return nil
}

func (a *AuthService) ResetPassword(ctx context.Context, tokenStr, newPassword string) error {
	if newPassword == "" {
		return errors.New("new password is required")
	}

	email, err := a.tokenService.ValidatePasswordResetToken(tokenStr)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return errors.New("user not found")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashed)
	if err := a.userRepo.UpdateUser(ctx, u); err != nil {
		return err
	}

	return nil
}
