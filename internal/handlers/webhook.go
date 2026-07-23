package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/utils"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/services"
)

type WebhookHandler struct {
	gitService        *services.GitService
	projectService    *services.ProjectService
	appService        *services.AppService
	deploymentService *services.DeploymentService
	prPreviewService  *services.PRPreviewService
	gitAppsService    *services.GitAppsService
}

func NewWebhookHandler(
	gitService *services.GitService,
	projectService *services.ProjectService,
	appService *services.AppService,
	deploymentService *services.DeploymentService,
	prPreviewService *services.PRPreviewService,
	gitAppsService *services.GitAppsService,
) *WebhookHandler {
	return &WebhookHandler{
		gitService:        gitService,
		projectService:    projectService,
		appService:        appService,
		deploymentService: deploymentService,
		prPreviewService:  prPreviewService,
		gitAppsService:    gitAppsService,
	}
}

func (h *WebhookHandler) HandleServiceWebhook(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	appSvc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || appSvc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	token := c.QueryParam("token")
	if appSvc.DeployToken == "" || token != appSvc.DeployToken {
		return utils.Error(c, http.StatusUnauthorized, "invalid or missing deploy token")
	}

	h.deployServiceAsync(appSvc)
	return utils.Accepted(c, fmt.Sprintf("triggering background build & rollout for service %s", appSvc.Name), nil)
}

func (h *WebhookHandler) deployServiceAsync(appSvc *models.AppService) {
	go func() {
		ctx := context.Background()
		dep := &models.Deployment{
			ID:            uuid.NewString(),
			ServiceID:     appSvc.ID,
			EnvironmentID: appSvc.EnvironmentID,
			ProjectID:     appSvc.ProjectID,
			Status:        "BUILDING",
			Branch:        appSvc.Branch,
			Trigger:       "Git Webhook Push",
			BuildLogs:     fmt.Sprintf("Initiating automated build from %s branch %s...\n", appSvc.RepositoryURL, appSvc.Branch),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		dep, _ = h.deploymentService.CreateDeployment(ctx, dep)
		h.deploymentService.ExecuteDeploymentAsync(dep)
	}()
}

func verifyHMAC(payload []byte, secret, signature string) bool {
	if secret == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}

func (h *WebhookHandler) HandleGitHubWebhook(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	event := c.Request().Header.Get("X-GitHub-Event")
	if event == "" {
		return utils.Error(c, http.StatusBadRequest, "missing X-GitHub-Event header")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "failed to read body")
	}

	signature := c.Request().Header.Get("X-Hub-Signature-256")
	if signature != "" {
		apps, err := h.gitAppsService.ListGithubApps(c.Request().Context())
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, "failed to check webhook secrets")
		}

		valid := false
		for _, app := range apps {
			if verifyHMAC(bodyBytes, app.WebhookSecret, signature) {
				valid = true
				break
			}
		}
		if !valid {
			return utils.Error(c, http.StatusUnauthorized, "invalid webhook signature")
		}
	}

	var payload models.GithubWebhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if event == "push" {
		appSvc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
		if err != nil || appSvc == nil {
			return utils.Error(c, http.StatusNotFound, "service not found")
		}
		expectedRef := "refs/heads/" + appSvc.Branch
		if payload.Ref == expectedRef {
			h.deployServiceAsync(appSvc)
			return utils.Accepted(c, fmt.Sprintf("triggering background build for branch %s", appSvc.Branch), nil)
		}
		return utils.Success(c, "Push event ignored (branch mismatch)", nil)
	}

	if event == "pull_request" {
		appSvc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
		if err != nil || appSvc == nil {
			return utils.Error(c, http.StatusNotFound, "service not found")
		}
		switch payload.Action {
		case "opened", "synchronize", "reopened":
			if !appSvc.EnablePRPreviews {
				return utils.Success(c, "PR previews are disabled for this service", nil)
			}
			h.deployPRPreviewAsync(serviceID, payload)
			return utils.Accepted(c, "Deploying PR preview", nil)
		case "closed":
			h.destroyPRPreviewAsync(serviceID, payload.Number)
			return utils.Accepted(c, "Destroying PR preview", nil)
		}
	}
	return utils.Success(c, "Operation successful", map[string]string{"message": "Event ignored"})
}

func (h *WebhookHandler) deployPRPreviewAsync(serviceID string, payload models.GithubWebhookPayload) {
	go func() {
		ctx := context.Background()
		opts := services.DeployPRPreviewOpts{
			AppID:      serviceID,
			PRNumber:   payload.Number,
			CommitHash: payload.PullRequest.Head.Sha,
			Branch:     payload.PullRequest.Head.Ref,
		}
		_, _ = h.prPreviewService.DeployPRPreview(ctx, opts)
	}()
}

func (h *WebhookHandler) destroyPRPreviewAsync(serviceID string, prNumber int) {
	go func() {
		ctx := context.Background()
		_ = h.prPreviewService.DestroyPRPreview(ctx, serviceID, prNumber)
	}()
}
