package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/utils"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/services"
)

type DomainHandler struct {
	envService *services.EnvironmentService
}

func NewDomainHandler(s *services.EnvironmentService) *DomainHandler {
	return &DomainHandler{envService: s}
}

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
