package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type AppHandler struct {
	appService *services.AppService
}

func NewAppHandler(s *services.AppService) *AppHandler {
	return &AppHandler{appService: s}
}

// @Summary Create endpoint
// @Description Create endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Router /api/workspaces [post]
func (h *AppHandler) Create(c echo.Context) error {
	envID := c.Param("id")
	var req models.AppService
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "app service name is required"})
	}
	req.EnvironmentID = envID
	if req.InternalPort == 0 {
		req.InternalPort = 3000
	}
	created, err := h.appService.CreateAppService(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

// @Summary ListByEnvironment endpoint
// @Description ListByEnvironment endpoint
// @Tags Environments
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/environments/{id}/apps [get]
func (h *AppHandler) ListByEnvironment(c echo.Context) error {
	envID := c.Param("id")
	apps, err := h.appService.ListByEnvironment(c.Request().Context(), envID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, apps)
}

// @Summary ListByProject endpoint
// @Description ListByProject endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/projects/{id}/apps [get]
func (h *AppHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	apps, err := h.appService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, apps)
}

// @Summary Get endpoint
// @Description Get endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /api/teams/{teamId}/ai_settings [get]
func (h *AppHandler) Get(c echo.Context) error {
	id := c.Param("id")
	svc, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || svc == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "app service not found"})
	}
	return c.JSON(http.StatusOK, svc)
}

// @Summary Update endpoint
// @Description Update endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/workspaces/{id} [put]
func (h *AppHandler) Update(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "app service not found"})
	}
	var req models.AppService
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	existing.Name = req.Name
	existing.RepositoryURL = req.RepositoryURL
	existing.Branch = req.Branch
	existing.RootDirectory = req.RootDirectory
	existing.BuildCommand = req.BuildCommand
	existing.StartCommand = req.StartCommand
	existing.DockerfilePath = req.DockerfilePath
	existing.BuildEngine = req.BuildEngine
	existing.InternalPort = req.InternalPort
	existing.Domain = req.Domain
	existing.HealthCheckPath = req.HealthCheckPath
	existing.ContainerID = req.ContainerID
	existing.Status = req.Status
	if err := h.appService.UpdateAppService(c.Request().Context(), existing); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, existing)
}

// @Summary Delete endpoint
// @Description Delete endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/workspaces/{id} [delete]
func (h *AppHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.appService.DeleteAppService(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
