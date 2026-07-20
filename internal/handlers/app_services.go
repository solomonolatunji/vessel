package handlers

import (
	"net/http"

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
	envService        *services.EnvironmentService
}

func NewAppHandler(s *services.AppService, ps *services.ProjectService, d *engine.Deployer, ds *services.DeploymentService, es *services.EnvironmentService) *AppHandler {
	return &AppHandler{
		appService:        s,
		projectService:    ps,
		deployer:          d,
		deploymentService: ds,
		envService:        es,
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

	// Auto-generate default domain
	domainName := utils.GenerateAppDomain(req.Name, "", "")
	_, _ = h.envService.CreateDomain(c.Request().Context(), &models.DomainConfig{
		ServiceID:  created.ID,
		DomainName: domainName,
	})

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
