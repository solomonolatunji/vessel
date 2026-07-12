package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type WebhookHandler struct {
	gitService        *services.GitService
	projectService    *services.ProjectService
	appService        *services.AppService
	deploymentService *services.DeploymentService
	prPreviewService  *services.PRPreviewService
}

func NewWebhookHandler(
	gitService *services.GitService,
	projectService *services.ProjectService,
	appService *services.AppService,
	deploymentService *services.DeploymentService,
	prPreviewService *services.PRPreviewService,
) *WebhookHandler {
	return &WebhookHandler{
		gitService:        gitService,
		projectService:    projectService,
		appService:        appService,
		deploymentService: deploymentService,
		prPreviewService:  prPreviewService,
	}
}

type GithubWebhookPayload struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		Head struct {
			Ref string `json:"ref"`
			Sha string `json:"sha"`
		} `json:"head"`
	} `json:"pull_request"`
}

// @Summary HandleProjectWebhook endpoint
// @Description HandleProjectWebhook endpoint
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param projectId path string true "projectId"
// @Router /webhooks/git/{projectId} [post]
func (h *WebhookHandler) HandleProjectWebhook(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId parameter"})
	}
	project, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil || project == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
	}
	go func() {
		ctx := context.Background()
		sourceDir := filepath.Join("data", "builds", project.ID)
		_, _ = h.deploymentService.DeployProject(ctx, project.ID, sourceDir, nil)
	}()
	return c.JSON(http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & deployment for %s", project.Name),
	})
}

// @Summary HandleServiceWebhook endpoint
// @Description HandleServiceWebhook endpoint
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param serviceId path string true "serviceId"
// @Router /webhooks/git/services/{serviceId} [post]
func (h *WebhookHandler) HandleServiceWebhook(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing serviceId parameter"})
	}
	appSvc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || appSvc == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "service not found"})
	}
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
		_, _ = h.deploymentService.CreateDeployment(ctx, dep)
		sourceDir := filepath.Join("data", "builds", "services", appSvc.ID)
		if h.gitService != nil && appSvc.RepositoryURL != "" {
			if err := h.gitService.CloneOrPullAppRepository(ctx, appSvc, sourceDir, nil); err != nil {
				log.Printf("[ServiceGitWebhook] Git clone/pull failed for service %s (%s): %v", appSvc.Name, appSvc.ID, err)
				_ = h.deploymentService.UpdateStatus(ctx, dep.ID, "FAILED", dep.BuildLogs+fmt.Sprintf("Error cloning repository: %v\n", err), "")
				return
			}
		}
		_ = h.deploymentService.UpdateStatus(ctx, dep.ID, "ACTIVE", dep.BuildLogs+"Deployment rollout triggered via Webhook.\n", appSvc.ContainerID)
	}()
	return c.JSON(http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & rollout for service %s", appSvc.Name),
	})
}

// @Summary HandleGitHubWebhook endpoint
// @Description HandleGitHubWebhook endpoint
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param serviceId path string true "serviceId"
// @Param request body handlers.GithubWebhookPayload true "Payload"
// @Router /webhooks/github/services/{serviceId} [post]
func (h *WebhookHandler) HandleGitHubWebhook(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing serviceId parameter"})
	}
	event := c.Request().Header.Get("X-GitHub-Event")
	if event == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing X-GitHub-Event header"})
	}
	var payload GithubWebhookPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if event == "pull_request" {
		if payload.Action == "opened" || payload.Action == "synchronize" {
			go func() {
				ctx := context.Background()
				_, _ = h.prPreviewService.DeployPRPreview(ctx, serviceID, payload.Number, payload.PullRequest.Head.Sha, payload.PullRequest.Head.Ref)
			}()
			return c.JSON(http.StatusAccepted, map[string]string{"message": "Deploying PR preview"})
		} else if payload.Action == "closed" {
			go func() {
				ctx := context.Background()
				_ = h.prPreviewService.DestroyPRPreview(ctx, serviceID, payload.Number)
			}()
			return c.JSON(http.StatusAccepted, map[string]string{"message": "Destroying PR preview"})
		}
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Event ignored"})
}
