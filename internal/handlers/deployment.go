package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type DeploymentHandler struct {
	deploymentService *services.DeploymentService
	appService        *services.AppService
}

func NewDeploymentHandler(ds *services.DeploymentService, as *services.AppService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService: ds,
		appService:        as,
	}
}

// @Summary ListServiceDeployments endpoint
// @Description ListServiceDeployments endpoint
// @Tags AppServices
// @Accept json
// @Produce json
// @Param serviceId path string true "serviceId"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Router /services/{serviceId}/deployments [get]
func (h *DeploymentHandler) ListServiceDeployments(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	deps, total, err := h.deploymentService.ListByService(c.Request().Context(), serviceID, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Paginated(c, "Deployments retrieved", deps, total, page, limit)
}

// @Summary Trigger Deployment
// @Description Trigger Deployment
// @Tags Deployments
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Router /services/{serviceId}/deploy [post]
func (h *DeploymentHandler) Trigger(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}
	dep := &models.Deployment{
		ServiceID:     serviceID,
		EnvironmentID: svc.EnvironmentID,
		ProjectID:     svc.ProjectID,
		Status:        "BUILDING",
		Branch:        svc.Branch,
		Trigger:       "Manual Deploy",
		BuildLogs:     "Initiating build...\n",
	}
	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), dep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)

	return utils.Accepted(c, "Deployment created", created)
}

// @Summary Rollback endpoint
// @Description Rollback endpoint
// @Tags Deployments
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /deployments/{id}/rollback [post]
func (h *DeploymentHandler) Rollback(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	targetDep, err := h.deploymentService.GetDeployment(c.Request().Context(), id)
	if err != nil || targetDep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}
	newDep := &models.Deployment{
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        "BUILDING",
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Rollback to " + targetDep.ID,
		Branch:        targetDep.Branch,
		Trigger:       "Rollback",
		BuildLogs:     "Rolling back to deployment " + targetDep.ID + "...\n",
	}
	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), newDep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)

	return utils.Accepted(c, "Rollback created", created)
}

// @Summary GetLogs endpoint
// @Description GetLogs endpoint
// @Tags Deployments
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /deployments/{id}/logs [get]
func (h *DeploymentHandler) GetLogs(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	dep, err := h.deploymentService.GetDeployment(c.Request().Context(), id)
	if err != nil || dep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}
	return utils.Success(c, "Logs fetched successfully", map[string]string{
		"id":        dep.ID,
		"buildLogs": dep.BuildLogs,
		"status":    string(dep.Status),
	})
}

// @Summary GetMetrics endpoint
// @Description GetMetrics endpoint
// @Tags AppServices
// @Accept json
// @Produce json
// @Param serviceId path string true "serviceId"
// @Router /services/{serviceId}/metrics [get]
func (h *DeploymentHandler) GetMetrics(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	health, err := h.deploymentService.GetMetrics(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	// For the dashboard, we return an array. We can simulate the 5-minute window or just return the current one.
	now := time.Now().UTC()
	metrics := []map[string]any{
		{
			"timestamp":  now.Format(time.RFC3339),
			"cpuPercent": health.CPUUsagePercentage,
			"memoryMB":   float64(health.MemoryUsageBytes) / 1024 / 1024,
			"status":     health.Status,
			"uptime":     health.UptimeSeconds,
		},
	}
	return utils.Success(c, "Operation successful", metrics)
}

type DeployRequest struct {
	DryRun bool `json:"dry_run"`
}

// @Summary DeployProject endpoint
// @Description DeployProject endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /projects/{id}/deploy [post]
func (h *DeploymentHandler) DeployProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}

	sourceDir := fmt.Sprintf("data/builds/%s", id)
	containerID, err := h.deploymentService.DeployProject(c.Request().Context(), id, sourceDir, nil)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Project deployed successfully", map[string]string{
		"status":       "deployed",
		"container_id": containerID,
	})
}
