package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

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
// @Tags Vercel
// @Accept json
// @Produce json
// @Router /api/vercel/projects [get]
func (h *ProjectHandler) ListProjects(c echo.Context) error {
	projects, err := h.projectService.ListProjects(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, projects)
}

// @Summary CreateProject endpoint
// @Description CreateProject endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Router /api/projects [post]
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var req models.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	p, err := h.projectService.CreateProjectFromRequest(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, p)
}

// @Summary GetProject endpoint
// @Description GetProject endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	p, err := h.projectService.GetProject(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
	}
	return c.JSON(http.StatusOK, p)
}

// @Summary DeleteProject endpoint
// @Description DeleteProject endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	if err := h.projectService.DeleteProject(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
