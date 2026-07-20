package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/utils"
)

type Mailer interface {
	SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error
}

type AuthService struct {
	userRepo        repositories.UserRepository
	settingsRepo    repositories.SettingsRepository
	notifRepo       repositories.NotificationSettingsRepository
	projectSettings repositories.ProjectSettingsRepository
	tokenService    *TokenService
	mailer          Mailer
}

func NewAuthService(
	ur repositories.UserRepository,
	sr repositories.SettingsRepository,
	nr repositories.NotificationSettingsRepository,
	psr repositories.ProjectSettingsRepository,
	ts *TokenService,
	mailer Mailer,
) *AuthService {
	return &AuthService{
		userRepo:        ur,
		settingsRepo:    sr,
		notifRepo:       nr,
		projectSettings: psr,
		tokenService:    ts,
		mailer:          mailer,
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
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}
	role := models.UserRoleMember
	if total == 0 {
		role = models.UserRoleOwner
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
	if !utils.CheckPasswordHash(password, u.PasswordHash) {
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

func (a *AuthService) ForgotPassword(ctx context.Context, email string, originUrl string) error {
	if email == "" {
		return errors.New("email is required")
	}

	cfg, err := a.notifRepo.GetNotificationSettings(ctx)
	if err != nil || cfg == nil {
		return errors.New("could not load notification settings")
	}

	if !cfg.SMTPEnabled && !cfg.ResendEnabled {
		return errors.New("your team is yet to set or enable email")
	}

	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return nil
	}

	token, err := a.tokenService.GeneratePasswordResetToken(u.Email)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"ResetUrl": originUrl + "/reset-password?token=" + token,
	}

	err = a.mailer.SendSystemEmail(ctx, "password_reset", u.Email, "Reset Your Password", data)
	if err != nil {
		return err
	}

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
		return utils.NewNotFoundError("User", email)
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashed)
	if err := a.userRepo.UpdateUser(ctx, u); err != nil {
		return err
	}

	_ = a.projectSettings.AcceptAllInvitesForUser(ctx, u.ID)

	return nil
}

func (a *AuthService) InviteUser(ctx context.Context, email string, role models.UserRole, originUrl string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	existing, _ := a.userRepo.GetUserByEmail(ctx, email)
	if existing != nil {
		return existing, nil
	}
	u := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         strings.Split(email, "@")[0],
		PasswordHash: "INVITED_NO_LOGIN_ALLOWED_MUST_RESET",
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := a.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	_ = a.ForgotPassword(ctx, u.Email, originUrl)

	return u, nil
}
