package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type EnvironmentHandler struct {
	envService *services.EnvironmentService
}

func NewEnvironmentHandler(s *services.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{envService: s}
}

func (h *EnvironmentHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	envs, err := h.envService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", envs)
}

func (h *EnvironmentHandler) Create(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	var env models.EnvironmentConfig
	if err := c.Bind(&env); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	env.ProjectID = projectID
	if env.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "environment name is required")
	}
	created, err := h.envService.CreateEnvironment(c.Request().Context(), &env)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *EnvironmentHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.envService.DeleteEnvironment(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
