package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	envs, err := h.envService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, envs)
}

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

type DomainHandler struct {
	envService *services.EnvironmentService
}

func NewDomainHandler(s *services.EnvironmentService) *DomainHandler {
	return &DomainHandler{envService: s}
}

func (h *DomainHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	domains, err := h.envService.ListDomainsByProject(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, domains)
}

func (h *DomainHandler) Create(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing project id parameter"})
	}
	var d models.DomainConfig
	if err := c.Bind(&d); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	d.ProjectID = projectID
	if d.DomainName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "domainName is required"})
	}
	created, err := h.envService.CreateDomain(c.Request().Context(), &d)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

func (h *DomainHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	if err := h.envService.DeleteDomain(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

type ProjectEnvHandler struct {
	envService *services.EnvironmentService
}

func NewProjectEnvHandler(s *services.EnvironmentService) *ProjectEnvHandler {
	return &ProjectEnvHandler{envService: s}
}

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
