package handlers

import (
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
	auditService      *services.AuditService
	aiAnalysis        *services.AIAnalysisService
	prPreviewService  *services.PRPreviewService
}

func NewDeploymentHandler(ds *services.DeploymentService, as *services.AppService, audit *services.AuditService, aiAnalysis *services.AIAnalysisService, prp *services.PRPreviewService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService: ds,
		appService:        as,
		auditService:      audit,
		aiAnalysis:        aiAnalysis,
		prPreviewService:  prp,
	}
}

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

	h.auditService.LogAction(c.Request().Context(), services.AuditActionOpts{
		UserID:    "system",
		Action:    "deployment.trigger",
		Resource:  serviceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"deploymentId": created.ID,
		},
	})

	return utils.Accepted(c, "Deployment created", created)
}

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

	h.auditService.LogAction(c.Request().Context(), services.AuditActionOpts{
		UserID:    "system",
		Action:    "deployment.rollback",
		Resource:  newDep.ServiceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"deploymentId": created.ID,
		},
	})

	return utils.Accepted(c, "Rollback created", created)
}

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

func (h *DeploymentHandler) GetMetrics(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	health, err := h.deploymentService.GetMetrics(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

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

func (h *DeploymentHandler) ExplainFailure(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	explanation, err := h.aiAnalysis.ExplainDeploymentFailure(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "AI Analysis completed", explanation)
}

func (h *DeploymentHandler) ListPRPreviews(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	previews, err := h.prPreviewService.ListByApp(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", previews)
}
