package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

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
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Router /projects/{id}/domains [get]
func (h *DomainHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	domains, err := h.envService.ListDomainsByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", domains)
}

// @Summary Create domain
// @Description Create domain
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body models.DomainConfig true "Payload"
// @Router /projects/{id}/domains [post]
func (h *DomainHandler) Create(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	var d models.DomainConfig
	if err := c.Bind(&d); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	d.ProjectID = projectID
	if d.DomainName == "" {
		return utils.Error(c, http.StatusBadRequest, "domainName is required")
	}
	created, err := h.envService.CreateDomain(c.Request().Context(), &d)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

// @Summary Delete Domain
// @Description Delete Domain
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Domain ID"
// @Router /domains/{id} [delete]
func (h *DomainHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.envService.DeleteDomain(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
