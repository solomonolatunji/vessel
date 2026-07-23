package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/utils"

	"codedock.dev/codedock/internal/http/middleware"
	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/services"
	"codedock.dev/codedock/internal/telemetry"
)

type ProjectHandler struct {
	projectService         *services.ProjectService
	projectSettingsService *services.ProjectSettingsService
}

func NewProjectHandler(s *services.ProjectService, pss *services.ProjectSettingsService) *ProjectHandler {
	return &ProjectHandler{projectService: s, projectSettingsService: pss}
}

func (h *ProjectHandler) ListProjects(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	projects, total, err := h.projectService.ListProjects(c.Request().Context(), limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Paginated(c, "Operation successful", projects, total, page, limit)
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var req models.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	userClaims, ok := c.Get("user").(*models.UserClaims)
	var p *models.ProjectConfig
	var err error

	if ok && userClaims != nil {
		p, err = h.projectService.CreateProjectWithMemberFromRequest(c.Request().Context(), &req, userClaims.UserID, string(models.MemberPermissionOwner))
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, "failed to create project and assign owner: "+err.Error())
		}
		telemetry.Track(userClaims.Email, "project_created", map[string]interface{}{
			"project_id": p.ID,
			"name":       p.Name,
		})
	} else {
		p, err = h.projectService.CreateProjectFromRequest(c.Request().Context(), &req)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		telemetry.Track("anonymous", "project_created", map[string]interface{}{
			"project_id": p.ID,
			"name":       p.Name,
		})
	}

	return utils.Created(c, "Created successfully", p)
}

func (h *ProjectHandler) GetProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	p, err := h.projectService.GetProject(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
	}
	return utils.Success(c, "Operation successful", p)
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	_, err := h.projectService.GetProject(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
	}
	if err := h.projectService.DeleteProject(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}
