package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

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
		return utils.Error(c, http.StatusBadRequest, "missing projectId parameter")
	}
	project, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil || project == nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	go func() {
		ctx := context.Background()
		sourceDir := filepath.Join(utils.GetDataDir(), "builds", project.ID)
		_, _ = h.deploymentService.DeployProject(ctx, project.ID, sourceDir, nil)
	}()
	return utils.Accepted(c, fmt.Sprintf("triggering background build & deployment for %s", project.Name), nil)
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
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	appSvc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || appSvc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
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
		dep, _ = h.deploymentService.CreateDeployment(ctx, dep)
		h.deploymentService.ExecuteDeploymentAsync(dep)
	}()
	return utils.Accepted(c, fmt.Sprintf("triggering background build & rollout for service %s", appSvc.Name), nil)
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
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	event := c.Request().Header.Get("X-GitHub-Event")
	if event == "" {
		return utils.Error(c, http.StatusBadRequest, "missing X-GitHub-Event header")
	}
	var payload GithubWebhookPayload
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if event == "pull_request" {
		if payload.Action == "opened" || payload.Action == "synchronize" {
			go func() {
				ctx := context.Background()
				_, _ = h.prPreviewService.DeployPRPreview(ctx, serviceID, payload.Number, payload.PullRequest.Head.Sha, payload.PullRequest.Head.Ref)
			}()
			return utils.Accepted(c, "Deploying PR preview", nil)
		} else if payload.Action == "closed" {
			go func() {
				ctx := context.Background()
				_ = h.prPreviewService.DestroyPRPreview(ctx, serviceID, payload.Number)
			}()
			return utils.Accepted(c, "Destroying PR preview", nil)
		}
	}
	return utils.Success(c, "Operation successful", map[string]string{"message": "Event ignored"})
}
