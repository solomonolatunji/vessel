package handlers

import (
	"net/http"

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
// @Router /projects [get]
func (h *ProjectHandler) ListProjects(c echo.Context) error {
	projects, err := h.projectService.ListProjects(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
		var filtered []models.ProjectConfig
		for _, p := range projects {
			if p.TeamID == user.UserID {
				filtered = append(filtered, p)
			}
		}
		return utils.Success(c, "Operation successful", filtered)
	}
	return utils.Success(c, "Operation successful", projects)
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
		req.TeamID = user.UserID
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
	if user != nil && user.Role != "admin" && p.TeamID != user.UserID {
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
	if user != nil && user.Role != "admin" && p.TeamID != user.UserID {
		return utils.Error(c, http.StatusForbidden, "access denied")
	}
	if err := h.projectService.DeleteProject(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}
