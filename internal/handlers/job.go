package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(s *services.JobService) *JobHandler {
	return &JobHandler{jobService: s}
}

// @Summary ListProjectJobs endpoint
// @Description ListProjectJobs endpoint
// @Tags Jobs
// @Accept json
// @Produce json
// @Router /api/jobs [get]
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

// @Summary Create endpoint
// @Description Create endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param request body models.Job true "Payload"
// @Router /api/jobs [post]
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

// @Summary Get endpoint
// @Description Get endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Summary Get Job
// @Description Get Job
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Router /api/jobs/{id} [get]
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

// @Summary Delete endpoint
// @Description Delete endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Summary Delete Job
// @Description Delete Job
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
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

// @Summary Run endpoint
// @Description Run endpoint
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/jobs/{id}/trigger [post]
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
