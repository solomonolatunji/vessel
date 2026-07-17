package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type OnboardingHandler struct {
	userService    *services.UserService
	authService    *services.AuthService
	settingsRepo   *services.SettingsService
	gitAppsService *services.GitAppsService
	backupService  *services.BackupService
}

func NewOnboardingHandler(
	userService *services.UserService,
	authService *services.AuthService,
	settingsRepo *services.SettingsService,
	gitAppsService *services.GitAppsService,
	backupService *services.BackupService,
) *OnboardingHandler {
	return &OnboardingHandler{
		userService:    userService,
		authService:    authService,
		settingsRepo:   settingsRepo,
		gitAppsService: gitAppsService,
		backupService:  backupService,
	}
}

// @Summary Check if onboarding is required
// @Description Returns true if no users exist in the system, indicating setup is needed
// @Tags System
// @Produce json
// @Success 200 {object} map[string]any
// @Router /system/setup-status [get]
func (h *OnboardingHandler) SetupStatus(c echo.Context) error {
	count, err := h.userService.CountUsers(c.Request().Context())
	if err != nil {
		return utils.Error(c, 500, "failed to check user count")
	}
	return utils.Success(c, "Setup status", map[string]bool{
		"setupRequired": count == 0,
	})
}

// RegisterRequest defines the expected payload for setup
type SetupRequest struct {
	// User
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	// Github Integration (optional)
	GithubAppID         string `json:"githubAppId,omitempty"`
	GithubClientID      string `json:"githubClientId,omitempty"`
	GithubClientSecret  string `json:"githubClientSecret,omitempty"`
	GithubPrivateKey    string `json:"githubPrivateKey,omitempty"`
	GithubWebhookSecret string `json:"githubWebhookSecret,omitempty"`
	GithubAppName       string `json:"githubAppName,omitempty"`

	// Domain (optional)
	DefaultWildcardDomain string `json:"defaultWildcardDomain,omitempty"`

	// Backups (optional)
	S3AccountID       string `json:"s3AccountId,omitempty"`
	S3Bucket          string `json:"s3Bucket,omitempty"`
	S3AccessKeyID     string `json:"s3AccessKeyId,omitempty"`
	S3SecretAccessKey string `json:"s3SecretAccessKey,omitempty"`
	S3Skip            bool   `json:"s3Skip,omitempty"`
}

// @Summary Complete onboarding setup
// @Description Creates the first user and optionally configures initial settings
// @Tags System
// @Accept json
// @Produce json
// @Param request body SetupRequest true "Setup details"
// @Success 200 {object} map[string]any
// @Router /system/setup [post]
func (h *OnboardingHandler) Setup(c echo.Context) error {
	ctx := c.Request().Context()

	count, err := h.userService.CountUsers(ctx)
	if err != nil {
		return utils.Error(c, 500, "failed to check user count")
	}
	if count > 0 {
		return utils.Error(c, 403, "Setup has already been completed")
	}

	var req SetupRequest
	if err := c.Bind(&req); err != nil {
		fmt.Printf("Setup Error: Failed to bind request: %v\n", err)
		return utils.Error(c, 400, "invalid request")
	}

	u, token, err := h.authService.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return utils.Error(c, 400, err.Error())
	}

	// Update settings
	settings, err := h.settingsRepo.GetSettings(ctx)
	if err == nil && settings != nil {
		updated := false
		if req.DefaultWildcardDomain != "" {
			settings.DefaultWildcardDomain = req.DefaultWildcardDomain
			updated = true
		}
		if updated {
			_ = h.settingsRepo.UpdateSettings(ctx, settings)
		}
	}

	// Save Github App
	if req.GithubAppID != "" && req.GithubClientID != "" && req.GithubClientSecret != "" && req.GithubPrivateKey != "" {
		appName := req.GithubAppName
		if appName == "" {
			appName = "Vessl Setup App"
		}
		_ = h.gitAppsService.SaveGithubApp(ctx, &models.GithubApp{
			Name:          appName,
			AppID:         req.GithubAppID,
			ClientID:      req.GithubClientID,
			ClientSecret:  req.GithubClientSecret,
			WebhookSecret: req.GithubWebhookSecret,
			PrivateKey:    req.GithubPrivateKey,
			IsPublic:      false,
		})
	}

	// Save S3 Destination
	if !req.S3Skip && req.S3AccountID != "" && req.S3Bucket != "" && req.S3AccessKeyID != "" && req.S3SecretAccessKey != "" {
		endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", req.S3AccountID)
		_ = h.backupService.CreateS3Destination(ctx, &models.S3Destination{
			Name:            "Default R2 Backup",
			Endpoint:        endpoint,
			Bucket:          req.S3Bucket,
			Region:          "auto",
			AccessKeyID:     req.S3AccessKeyID,
			SecretAccessKey: req.S3SecretAccessKey,
		})
	}

	res := map[string]any{
		"user":  u,
		"token": token,
	}

	return utils.Success(c, "Setup completed successfully", res)
}
