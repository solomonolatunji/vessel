package handlers

import (
	"fmt"
	"net/http"
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
// @Router /services/{serviceId}/deployments [get]
func (h *DeploymentHandler) ListServiceDeployments(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	deps, err := h.deploymentService.ListByService(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", deps)
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
	return c.JSON(http.StatusAccepted, created)
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
	return c.JSON(http.StatusAccepted, created)
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
	return c.JSON(http.StatusOK, map[string]string{
		"id":        dep.ID,
		"buildLogs": dep.BuildLogs,
		"status":    dep.Status,
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
	now := time.Now().UTC()
	metrics := []map[string]any{
		{"timestamp": now.Add(-4 * time.Minute).Format(time.RFC3339), "cpuPercent": 1.2, "memoryMB": 64.5, "networkRx": 12.4, "networkTx": 8.1},
		{"timestamp": now.Add(-3 * time.Minute).Format(time.RFC3339), "cpuPercent": 2.1, "memoryMB": 66.0, "networkRx": 15.0, "networkTx": 10.2},
		{"timestamp": now.Add(-2 * time.Minute).Format(time.RFC3339), "cpuPercent": 1.8, "memoryMB": 65.2, "networkRx": 14.1, "networkTx": 9.4},
		{"timestamp": now.Add(-1 * time.Minute).Format(time.RFC3339), "cpuPercent": 3.4, "memoryMB": 68.1, "networkRx": 45.2, "networkTx": 22.0},
		{"timestamp": now.Format(time.RFC3339), "cpuPercent": 1.5, "memoryMB": 66.8, "networkRx": 18.0, "networkTx": 11.5},
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
	return c.JSON(http.StatusOK, map[string]string{
		"status":       "deployed",
		"container_id": containerID,
	})
}
