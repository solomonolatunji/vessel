package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type AppHandler struct {
	appService *services.AppService
}

func NewAppHandler(s *services.AppService) *AppHandler {
	return &AppHandler{appService: s}
}

func (h *AppHandler) Create(c echo.Context) error {
	envID := c.Param("id")
	var req models.AppService
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "app service name is required"})
	}
	req.EnvironmentID = envID
	if req.InternalPort == 0 {
		req.InternalPort = 3000
	}
	created, err := h.appService.CreateAppService(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

func (h *AppHandler) ListByEnvironment(c echo.Context) error {
	envID := c.Param("id")
	apps, err := h.appService.ListByEnvironment(c.Request().Context(), envID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, apps)
}

func (h *AppHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	apps, err := h.appService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, apps)
}

func (h *AppHandler) Get(c echo.Context) error {
	id := c.Param("id")
	svc, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || svc == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "app service not found"})
	}
	return c.JSON(http.StatusOK, svc)
}

func (h *AppHandler) Update(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "app service not found"})
	}
	var req models.AppService
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	existing.Name = req.Name
	existing.RepositoryURL = req.RepositoryURL
	existing.Branch = req.Branch
	existing.InternalPort = req.InternalPort
	existing.Domain = req.Domain
	existing.ContainerID = req.ContainerID
	existing.Status = req.Status
	if err := h.appService.UpdateAppService(c.Request().Context(), existing); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, existing)
}

func (h *AppHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.appService.DeleteAppService(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

type ServiceVarHandler struct {
	appService *services.AppService
}

func NewServiceVarHandler(s *services.AppService) *ServiceVarHandler {
	return &ServiceVarHandler{appService: s}
}

func (h *ServiceVarHandler) List(c echo.Context) error {
	serviceID := c.Param("serviceId")
	list, err := h.appService.ListVariablesByService(c.Request().Context(), serviceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

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

func (h *ServiceVarHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.appService.DeleteVariable(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
