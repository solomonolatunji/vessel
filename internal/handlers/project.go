package handlers

import (
	"github.com/labstack/echo/v4"

	"context"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type ProxyReloader interface {
	Reload(ctx context.Context) error
}

type ProjectHandler struct {
	projectService *services.ProjectService
	proxy          ProxyReloader
}

func NewProjectHandler(s *services.ProjectService, p ProxyReloader) *ProjectHandler {
	return &ProjectHandler{projectService: s, proxy: p}
}

func (h *ProjectHandler) ListProjects(c echo.Context) error {
	projects, err := h.projectService.ListProjects(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var req models.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	p, err := h.projectService.CreateProjectFromRequest(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if h.proxy != nil {
		_ = h.proxy.Reload(c.Request().Context())
	}
	return c.JSON(http.StatusCreated, p)
}

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

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	if err := h.projectService.DeleteProject(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if h.proxy != nil {
		_ = h.proxy.Reload(c.Request().Context())
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
