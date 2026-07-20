package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/telemetry"
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

	domainName := utils.GenerateAppDomain(req.Name, "", "")
	_, _ = h.envService.CreateDomain(c.Request().Context(), &models.DomainConfig{
		ServiceID:  created.ID,
		DomainName: domainName,
	})

	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	distinctID := "anonymous"
	if user != nil {
		distinctID = user.Email
	}
	sourceType := "github"
	if created.ImageRef != "" {
		sourceType = "docker_image"
	}
	telemetry.Track(distinctID, "app_created", map[string]interface{}{
		"app_id": created.ID,
		"name":   created.Name,
		"type":   sourceType,
	})

	return utils.Created(c, "Created successfully", created)
}

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
	existing.CPULimit = req.CPULimit
	existing.MemoryLimit = req.MemoryLimit
	existing.DeployToken = req.DeployToken
	if err := h.appService.UpdateAppService(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", existing)
}

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
