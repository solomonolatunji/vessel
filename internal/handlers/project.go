package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler(s *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: s}
}

// @Summary ListProjects endpoint
// @Description ListProjects endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Router /projects [get]
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

	var workspaceID string
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
		workspaceID = user.UserID
	}

	projects, total, err := h.projectService.ListProjects(c.Request().Context(), workspaceID, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Paginated(c, "Operation successful", projects, total, page, limit)
}

// @Summary CreateProject endpoint
// @Description CreateProject endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param request body models.CreateProjectRequest true "Payload"
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var req models.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil {
		req.WorkspaceID = user.UserID
	}
	p, err := h.projectService.CreateProjectFromRequest(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", p)
}

// @Summary GetProject endpoint
// @Description GetProject endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /projects/{id} [get]
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
	if user != nil && user.Role != "admin" && p.WorkspaceID != user.UserID {
		return utils.Error(c, http.StatusForbidden, "access denied")
	}
	return utils.Success(c, "Operation successful", p)
}

// @Summary DeleteProject endpoint
// @Description DeleteProject endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	p, err := h.projectService.GetProject(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" && p.WorkspaceID != user.UserID {
		return utils.Error(c, http.StatusForbidden, "access denied")
	}
	if err := h.projectService.DeleteProject(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}
