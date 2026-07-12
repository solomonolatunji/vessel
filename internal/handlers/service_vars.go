package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type ServiceVarHandler struct {
	appService *services.AppService
}

func NewServiceVarHandler(s *services.AppService) *ServiceVarHandler {
	return &ServiceVarHandler{appService: s}
}

// @Summary List endpoint
// @Description List endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Router /api/workspaces [get]
// @Summary List Service Variables
// @Description List Service Variables
// @Tags ServiceVariables
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Router /api/services/{serviceId}/variables [get]
func (h *ServiceVarHandler) List(c echo.Context) error {
	serviceID := c.Param("serviceId")
	list, err := h.appService.ListVariablesByService(c.Request().Context(), serviceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

// @Summary Create endpoint
// @Description Create endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Router /api/workspaces [post]
// @Summary Create Service Variable
// @Description Create Service Variable
// @Tags ServiceVariables
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Router /api/services/{serviceId}/variables [post]
func (h *ServiceVarHandler) Create(c echo.Context) error {
	serviceID := c.Param("serviceId")
	var req models.Variable
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "service not found"})
	}
	req.ServiceID = serviceID
	req.ProjectID = svc.ProjectID
	req.EnvironmentID = svc.EnvironmentID
	created, err := h.appService.CreateVariable(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

// @Summary Update endpoint
// @Description Update endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/workspaces/{id} [put]
// @Summary Update Service Variable
// @Description Update Service Variable
// @Tags ServiceVariables
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Param id path string true "Variable ID"
// @Router /api/services/{serviceId}/variables/{id} [put]
func (h *ServiceVarHandler) Update(c echo.Context) error {
	serviceID := c.Param("serviceId")
	id := c.Param("id")
	var req models.Variable
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	req.ID = id
	req.ServiceID = serviceID
	if err := h.appService.UpdateVariable(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, req)
}

// @Summary Delete endpoint
// @Description Delete endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/workspaces/{id} [delete]
// @Summary Delete Service Variable
// @Description Delete Service Variable
// @Tags ServiceVariables
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Param id path string true "Variable ID"
// @Router /api/services/{serviceId}/variables/{id} [delete]
func (h *ServiceVarHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.appService.DeleteVariable(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
