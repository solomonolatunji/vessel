package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(s *services.JobService) *JobHandler {
	return &JobHandler{jobService: s}
}

func (h *JobHandler) ListProjectJobs(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId query parameter"})
	}
	jobs, err := h.jobService.ListJobsByProject(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) Create(c echo.Context) error {
	var j models.Job
	if err := c.Bind(&j); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	created, err := h.jobService.CreateJob(c.Request().Context(), &j)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

func (h *JobHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	j, err := h.jobService.GetJob(c.Request().Context(), id)
	if err != nil || j == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "job not found"})
	}
	return c.JSON(http.StatusOK, j)
}

func (h *JobHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	if err := h.jobService.DeleteJob(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *JobHandler) Run(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	out, err := h.jobService.ExecuteJob(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "executed", "output": out})
}
