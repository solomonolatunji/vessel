package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type OAuthService struct {
	oauthRepo    repositories.OAuthRepository
	userRepo        repositories.UserRepository
	tokenService    *TokenService
	pendingTOTP     sync.Map
	pendingRecovery sync.Map
}

func NewOAuthService(or repositories.OAuthRepository, ur repositories.UserRepository, ts *TokenService) *OAuthService {
	return &OAuthService{
		oauthRepo:    or,
		userRepo:     ur,
		tokenService: ts,
	}
}

func (s *OAuthService) ListProviders(ctx context.Context) ([]models.OAuthProviderConfig, error) {
	return s.oauthRepo.ListProviders(ctx)
}

func (s *OAuthService) ListEnabledProviders(ctx context.Context) ([]models.OAuthProviderConfig, error) {
	allProviders, err := s.oauthRepo.ListProviders(ctx)
	if err != nil {
		return nil, err
	}
	var enabledProviders []models.OAuthProviderConfig
	for _, p := range allProviders {
		if p.Enabled {
			p.ClientSecret = ""
			enabledProviders = append(enabledProviders, p)
		}
	}
	return enabledProviders, nil
}

func (s *OAuthService) GetProvider(ctx context.Context, idOrName string) (*models.OAuthProviderConfig, error) {
	if idOrName == "" {
		return nil, errors.New("provider id or name required")
	}
	return s.oauthRepo.GetProvider(ctx, idOrName)
}

func (s *OAuthService) SaveProvider(ctx context.Context, p *models.OAuthProviderConfig) error {
	if p == nil || p.ProviderName == "" {
		return errors.New("valid provider required")
	}
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
	return s.oauthRepo.SaveProvider(ctx, p)
}

func (s *OAuthService) GetUserTOTP(ctx context.Context, userID string) (string, []string, error) {
	if userID == "" {
		return "", nil, errors.New("user id required")
	}
	return s.oauthRepo.GetUserTOTPSecret(ctx, userID)
}

func (s *OAuthService) UpdateUserTOTP(ctx context.Context, userID string, enabled bool, secret string, recoveryCodes []string) error {
	if userID == "" {
		return errors.New("user id required")
	}
	return s.oauthRepo.UpdateUserTOTP(ctx, userID, enabled, secret, recoveryCodes)
}

func (s *OAuthService) HandleCallback(ctx context.Context, providerName, code string) (string, *models.User, error) {
	p, err := s.oauthRepo.GetProvider(ctx, providerName)
	if err != nil || p == nil {
		return "", nil, errors.New("oauth provider not found: " + providerName)
	}
	email, err := ExchangeCode(p, code)
	if err != nil || email == "" {
		return "", nil, errors.New("failed oauth code exchange")
	}
	u, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if u == nil {
		u = &models.User{
			ID:           uuid.New().String(),
			Email:        email,
			PasswordHash: "oauth-login-no-password",
			Role:         "member",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := s.userRepo.CreateUser(ctx, u); err != nil {
			return "", nil, errors.New("failed to create user account from oauth: " + err.Error())
		}

	}
	token, err := s.tokenService.GenerateToken(u)
	if err != nil {
		return "", nil, errors.New("failed generating token")
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	return token, &uCopy, nil
}

func (s *OAuthService) Setup2FA(ctx context.Context, userID, email string) (*models.TwoFASetupResponse, error) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		return nil, errors.New("failed generating totp secret")
	}
	recoveryCodes, err := GenerateRecoveryCodes(8)
	if err != nil {
		return nil, errors.New("failed generating recovery codes")
	}
	s.pendingTOTP.Store(userID, secret)
	s.pendingRecovery.Store(userID, recoveryCodes)
	
	return &models.TwoFASetupResponse{
		QRCodeURI:     GenerateTOTPQRUri(email, secret),
	}, nil
}

func (s *OAuthService) Verify2FA(ctx context.Context, userID, passcode string) error {
	secretAny, ok := s.pendingTOTP.Load(userID)
	if !ok {
		return errors.New("totp setup has not been initiated or has expired")
	}
	secret := secretAny.(string)

	if !ValidateTOTP(secret, passcode) {
		return errors.New("invalid 6-digit totp verification code")
	}

	recoveryAny, _ := s.pendingRecovery.Load(userID)
	recoveryCodes := recoveryAny.([]string)

	err := s.oauthRepo.UpdateUserTOTP(ctx, userID, true, secret, recoveryCodes)
	if err == nil {
		s.pendingTOTP.Delete(userID)
		s.pendingRecovery.Delete(userID)
	}
	return err
}

func (s *OAuthService) Disable2FA(ctx context.Context, userID string) error {
	return s.oauthRepo.UpdateUserTOTP(ctx, userID, false, "", nil)
}
