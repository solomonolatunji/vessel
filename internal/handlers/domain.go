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

// @Summary List domains by service
// @Description List domains by service
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Router /services/{id}/domains [get]
func (h *DomainHandler) ListByService(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id parameter")
	}
	domains, err := h.envService.ListDomainsByService(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", domains)
}

// @Summary Create domain
// @Description Create domain
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param request body models.DomainConfig true "Payload"
// @Router /services/{id}/domains [post]
func (h *DomainHandler) Create(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id parameter")
	}
	var d models.DomainConfig
	if err := c.Bind(&d); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	d.ServiceID = serviceID
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
