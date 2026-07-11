package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
)

type ProjectEnvHandler struct {
	envService *services.EnvironmentService
}

func NewProjectEnvHandler(s *services.EnvironmentService) *ProjectEnvHandler {
	return &ProjectEnvHandler{envService: s}
}

// @Summary GetVars endpoint
// @Description GetVars endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/projects/{id}/env [get]
func (h *ProjectEnvHandler) GetVars(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	vars, err := h.envService.GetVars(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if vars == nil {
		vars = map[string]string{}
	}
	return c.JSON(http.StatusOK, vars)
}

// @Summary SetVars endpoint
// @Description SetVars endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/projects/{id}/env [put]
func (h *ProjectEnvHandler) SetVars(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	var vars map[string]string
	if err := c.Bind(&vars); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	for k, v := range vars {
		if err := h.envService.SetVar(c.Request().Context(), projectID, k, v); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	return c.JSON(http.StatusOK, vars)
}
