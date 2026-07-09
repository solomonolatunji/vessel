package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type WebhookHandler struct {
	gitService        *services.GitService
	projectService    *services.ProjectService
	appService        *services.AppService
	deploymentService *services.DeploymentService
}

func NewWebhookHandler(
	gitService *services.GitService,
	projectService *services.ProjectService,
	appService *services.AppService,
	deploymentService *services.DeploymentService,
) *WebhookHandler {
	return &WebhookHandler{
		gitService:        gitService,
		projectService:    projectService,
		appService:        appService,
		deploymentService: deploymentService,
	}
}

func (h *WebhookHandler) HandleProjectWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId parameter")
		return
	}

	project, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil || project == nil {
		WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	WriteJSON(w, http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & deployment for %s", project.Name),
	})

	go func() {
		ctx := context.Background()
		sourceDir := filepath.Join("data", "builds", project.ID)
		_, _ = h.deploymentService.DeployProject(ctx, project.ID, sourceDir, nil)
	}()
}

func (h *WebhookHandler) HandleServiceWebhook(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "missing serviceId parameter")
		return
	}

	appSvc, err := h.appService.GetAppService(r.Context(), serviceID)
	if err != nil || appSvc == nil {
		WriteError(w, http.StatusNotFound, "service not found")
		return
	}

	WriteJSON(w, http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & rollout for service %s", appSvc.Name),
	})

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
}
