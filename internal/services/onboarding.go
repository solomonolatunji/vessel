package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"codedock.run/codedock/internal/models"
)

type SetupEnv struct {
	JWTSecret    string `json:"jwtSecret"`
	DataDir      string `json:"dataDir"`
	DashboardURL string `json:"dashboardUrl"`
	Port         int    `json:"port"`
}

type SetupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Env SetupEnv `json:"env"`

	DefaultWildcardDomain string `json:"defaultWildcardDomain,omitempty"`
}

type OnboardingService struct {
	userService    *UserService
	authService    *AuthService
	settingsRepo   *SettingsService
	gitAppsService *GitAppsService
	backupService  *BackupService
}

func NewOnboardingService(
	userService *UserService,
	authService *AuthService,
	settingsRepo *SettingsService,
	gitAppsService *GitAppsService,
	backupService *BackupService,
) *OnboardingService {
	return &OnboardingService{
		userService:    userService,
		authService:    authService,
		settingsRepo:   settingsRepo,
		gitAppsService: gitAppsService,
		backupService:  backupService,
	}
}

func (s *OnboardingService) CompleteSetup(ctx context.Context, req SetupRequest) (*models.User, string, string, error) {
	count, err := s.userService.CountUsers(ctx)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to check user count: %w", err)
	}
	if count > 0 {
		return nil, "", "", fmt.Errorf("setup has already been completed")
	}

	u, token, refreshToken, err := s.authService.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, "", "", fmt.Errorf("registration failed: %w", err)
	}

	if req.Env.JWTSecret == "" {
		req.Env.JWTSecret = s.generateJWTSecret()
	}

	if err := s.writeEnvFile(req.Env); err != nil {
		fmt.Printf("Warning: failed to write .env file: %v\n", err)
	}

	if req.DefaultWildcardDomain != "" {
		if err := s.updateDefaultDomain(ctx, req.DefaultWildcardDomain); err != nil {
			fmt.Printf("Warning: failed to update default domain: %v\n", err)
		}
	}

	return u, token, refreshToken, nil
}

func (s *OnboardingService) generateJWTSecret() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *OnboardingService) writeEnvFile(env SetupEnv) error {
	envBytes, err := os.ReadFile(".env.example")
	if err != nil {
		envContent := fmt.Sprintf("CODEDOCK_JWT_SECRET=%s\nCODEDOCK_DATA_DIR=%s\nCODEDOCK_DASHBOARD_URL=%s\nPORT=%d\n",
			env.JWTSecret, env.DataDir, env.DashboardURL, env.Port)
		return os.WriteFile(".env", []byte(envContent), 0644)
	}

	envStr := string(envBytes)
	envStr = strings.ReplaceAll(envStr, "CODEDOCK_JWT_SECRET=change-this-to-a-secure-random-secret-in-prod", "CODEDOCK_JWT_SECRET="+env.JWTSecret)
	if env.DataDir != "" {
		envStr = strings.ReplaceAll(envStr, "CODEDOCK_DATA_DIR=./data", "CODEDOCK_DATA_DIR="+env.DataDir)
	}
	if env.DashboardURL != "" {
		envStr = strings.ReplaceAll(envStr, "CODEDOCK_DASHBOARD_URL=http://localhost:3000", "CODEDOCK_DASHBOARD_URL="+env.DashboardURL)
	}
	if env.Port != 0 {
		envStr = strings.ReplaceAll(envStr, "PORT=8080", fmt.Sprintf("PORT=%d", env.Port))
	}
	return os.WriteFile(".env", []byte(envStr), 0644)
}

func (s *OnboardingService) updateDefaultDomain(ctx context.Context, domain string) error {
	settings, err := s.settingsRepo.GetSettings(ctx)
	if err != nil || settings == nil {
		return err
	}
	settings.DefaultWildcardDomain = domain
	return s.settingsRepo.UpdateSettings(ctx, settings)
}
