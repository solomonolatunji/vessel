package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type AppHandler struct {
	appService        *services.AppService
	projectService    *services.ProjectService
	deployer          *engine.Deployer
	deploymentService *services.DeploymentService
}

func NewAppHandler(s *services.AppService, ps *services.ProjectService, d *engine.Deployer, ds *services.DeploymentService) *AppHandler {
	return &AppHandler{
		appService:        s,
		projectService:    ps,
		deployer:          d,
		deploymentService: ds,
	}
}

func (h *AppHandler) verifyProjectOwnership(c echo.Context, projectID string) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role == "api" {
		tokenProjectID, ok := c.Get("project_id").(string)
		if ok && tokenProjectID != projectID {
			return utils.Error(c, http.StatusForbidden, "token does not have access to this project")
		}
	}
	_, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	return nil
}

// @Summary Create endpoint
// @Description Create endpoint
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "Environment ID"
// @Param request body models.AppService true "Payload"
// @Router /environments/{id}/apps [post]
func (h *AppHandler) Create(c echo.Context) error {
	envID := c.Param("id")
	var req models.AppService
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "app service name is required")
	}
	if err := h.verifyProjectOwnership(c, req.ProjectID); err != nil {
		return err
	}
	req.EnvironmentID = envID
	if req.InternalPort == 0 {
		req.InternalPort = 3000
	}
	if req.RuntimeMode == "" {
		req.RuntimeMode = models.RuntimeModeWeb
	}
	created, err := h.appService.CreateAppService(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

// @Summary ListByEnvironment endpoint
// @Description ListByEnvironment endpoint
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /environments/{id}/apps [get]
func (h *AppHandler) ListByEnvironment(c echo.Context) error {
	envID := c.Param("id")
	apps, err := h.appService.ListByEnvironment(c.Request().Context(), envID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
		var filtered []*models.AppService
		for _, app := range apps {
			_, err := h.projectService.GetProject(c.Request().Context(), app.ProjectID)
			if err == nil {
				filtered = append(filtered, app)
			}
		}
		return utils.Success(c, "Operation successful", filtered)
	}
	return utils.Success(c, "Operation successful", apps)
}

// @Summary ListByProject endpoint
// @Description ListByProject endpoint
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /projects/{id}/apps [get]
func (h *AppHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if err := h.verifyProjectOwnership(c, projectID); err != nil {
		return err
	}
	apps, err := h.appService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", apps)
}

// @Summary Get App Service
// @Description Get App Service
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "App ID"
// @Router /apps/{id} [get]
func (h *AppHandler) Get(c echo.Context) error {
	id := c.Param("id")
	svc, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, svc.ProjectID); err != nil {
		return err
	}
	return utils.Success(c, "Operation successful", svc)
}

// @Summary Update App Service
// @Description Update App Service
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "App ID"
// @Param request body models.AppService true "Payload"
// @Router /apps/{id} [put]
func (h *AppHandler) Update(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	var req models.AppService
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	existing.Name = req.Name
	existing.RepositoryURL = req.RepositoryURL
	existing.Branch = req.Branch
	existing.RootDirectory = req.RootDirectory
	existing.BuildCommand = req.BuildCommand
	existing.StartCommand = req.StartCommand
	existing.InstallCommand = req.InstallCommand
	existing.DockerfilePath = req.DockerfilePath
	existing.BuildEngine = req.BuildEngine
	existing.InternalPort = req.InternalPort
	existing.RuntimeMode = req.RuntimeMode
	existing.Domain = req.Domain
	existing.StaticOutput = req.StaticOutput
	existing.HealthCheckPath = req.HealthCheckPath
	existing.ContainerID = req.ContainerID
	existing.Status = req.Status
	if err := h.appService.UpdateAppService(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", existing)
}

// @Summary Delete App Service
// @Description Delete App Service
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "App ID"
// @Router /apps/{id} [delete]
func (h *AppHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	if err := h.appService.DeleteAppService(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary Stop App Service
// @Description Stops all containers belonging to this app service
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "App ID"
// @Router /apps/{id}/stop [post]
func (h *AppHandler) StopService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	if err := h.deployer.StopAppService(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	existing.Status = models.AppServiceStatusStopped
	_ = h.appService.UpdateAppService(c.Request().Context(), existing)
	return utils.Success(c, "Service stopped successfully", existing)
}

// @Summary Redeploy App Service
// @Description Creates a new deployment for this app service using the same branch/commit as the last deployment
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "App ID"
// @Router /apps/{id}/redeploy [post]
func (h *AppHandler) RedeployService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}

	newDep := &models.Deployment{
		ServiceID:     existing.ID,
		EnvironmentID: existing.EnvironmentID,
		ProjectID:     existing.ProjectID,
		Status:        "BUILDING",
		CommitMessage: "Manual Redeploy",
		Branch:        existing.Branch,
		Trigger:       "Manual Redeploy",
	}
	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), newDep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)
	return utils.Accepted(c, "Redeployment triggered", created)
}

// @Summary Restart App Service
// @Description Restarts all containers belonging to this app service
// @Tags AppServices
// @Accept json
// @Produce json
// @Param id path string true "App ID"
// @Router /apps/{id}/restart [post]
func (h *AppHandler) RestartService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	if err := h.deployer.RestartAppService(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	existing.Status = models.AppServiceStatusRunning
	_ = h.appService.UpdateAppService(c.Request().Context(), existing)
	return utils.Success(c, "Service restarted successfully", existing)
}

// @Summary ListWebhooks endpoint
// @Description ListWebhooks endpoint
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /apps/{id}/webhooks [get]
func (h *AppHandler) ListWebhooks(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId")
	}
	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to look up app service")
		}
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	list, err := h.appService.ListWebhooks(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

// @Summary CreateWebhook endpoint
// @Description CreateWebhook endpoint
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param request body models.Webhook true "Payload"
// @Router /apps/{id}/webhooks [post]
func (h *AppHandler) CreateWebhook(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId")
	}
	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to look up app service")
		}
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	var req models.CreateWebhookRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		return utils.Error(c, http.StatusBadRequest, "missing url")
	}
	parsedURL, err := url.Parse(req.URL)
	if err != nil || !parsedURL.IsAbs() || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return utils.Error(c, http.StatusBadRequest, "invalid webhook url: must be an absolute http/https url")
	}
	for _, et := range req.EventTypes {
		if strings.Contains(et, ",") {
			return utils.Error(c, http.StatusBadRequest, "event type cannot contain commas")
		}
	}
	webhook := models.Webhook{
		ServiceID:             serviceID,
		URL:                   req.URL,
		EventTypes:            req.EventTypes,
		IncludePREnvironments: req.IncludePREnvironments,
	}
	created, err := h.appService.CreateWebhook(c.Request().Context(), &webhook)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

// @Summary DeleteWebhook endpoint
// @Description DeleteWebhook endpoint
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param webhookId path string true "webhookId"
// @Router /apps/{id}/webhooks/{webhookId} [delete]
func (h *AppHandler) DeleteWebhook(c echo.Context) error {
	serviceID := c.Param("id")
	webhookID := c.Param("webhookId")
	if serviceID == "" || webhookID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId or webhookId")
	}
	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to look up app service")
		}
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	if err := h.appService.DeleteWebhook(c.Request().Context(), webhookID, serviceID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
