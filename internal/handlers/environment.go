package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type EnvironmentHandler struct {
	envService *services.EnvironmentService
}

func NewEnvironmentHandler(s *services.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{envService: s}
}

// @Summary ListByProject endpoint
// @Description ListByProject endpoint
// @Tags Environments
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /projects/{id}/environments [get]
func (h *EnvironmentHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	envs, err := h.envService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, envs)
}

// @Summary Create endpoint
// @Description Create endpoint
// @Tags Environments
// @Accept json
// @Produce json
// @Param request body models.EnvironmentConfig true "Payload"
// @Router /projects/{id}/environments [post]
func (h *EnvironmentHandler) Create(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	var env models.EnvironmentConfig
	if err := c.Bind(&env); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	env.ProjectID = projectID
	if env.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "environment name is required"})
	}
	created, err := h.envService.CreateEnvironment(c.Request().Context(), &env)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

// @Summary Delete Environment
// @Description Delete Environment
// @Tags Environments
// @Accept json
// @Produce json
// @Param id path string true "Environment ID"
// @Router /environments/{id} [delete]
func (h *EnvironmentHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	if err := h.envService.DeleteEnvironment(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
