package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
)

type ScheduledTaskHandler struct {
	scheduledTaskService *services.ScheduledTaskService
}

func NewScheduledTaskHandler(s *services.ScheduledTaskService) *ScheduledTaskHandler {
	return &ScheduledTaskHandler{scheduledTaskService: s}
}

func (h *ScheduledTaskHandler) ListProjectScheduledTasks(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	serviceID := c.QueryParam("serviceId")

	var tasks []models.ScheduledTask
	var err error

	if serviceID != "" {
		tasks, err = h.scheduledTaskService.ListScheduledTasksByService(c.Request().Context(), serviceID)
	} else {
		tasks, err = h.scheduledTaskService.ListScheduledTasksByProject(c.Request().Context(), projectID)
	}

	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", tasks)
}

func (h *ScheduledTaskHandler) Create(c echo.Context) error {
	var j models.ScheduledTask
	if err := c.Bind(&j); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	created, err := h.scheduledTaskService.CreateScheduledTask(c.Request().Context(), &j)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *ScheduledTaskHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	j, err := h.scheduledTaskService.GetScheduledTask(c.Request().Context(), id)
	if err != nil || j == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}
	return utils.Success(c, "Operation successful", j)
}

func (h *ScheduledTaskHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.scheduledTaskService.DeleteScheduledTask(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ScheduledTaskHandler) Run(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	out, err := h.scheduledTaskService.ExecuteScheduledTask(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "executed", "output": out})
}
