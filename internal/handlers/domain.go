package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type DomainHandler struct {
	envService *services.EnvironmentService
}

func NewDomainHandler(s *services.EnvironmentService) *DomainHandler {
	return &DomainHandler{envService: s}
}

// @Summary List domains by project
// @Description List domains by project
// @Tags Domains
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Router /api/projects/{id}/domains [get]
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

// @Summary Create domain
// @Description Create domain
// @Tags Domains
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body models.DomainConfig true "Payload"
// @Router /api/projects/{id}/domains [post]
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

// @Summary Delete Domain
// @Description Delete Domain
// @Tags Domains
// @Accept json
// @Produce json
// @Param id path string true "Domain ID"
// @Router /api/domains/{id} [delete]
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
