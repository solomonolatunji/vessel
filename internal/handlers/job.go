package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

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
// @Router /jobs [get]
func (h *JobHandler) ListProjectJobs(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId query parameter")
	}
	jobs, err := h.jobService.ListJobsByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", jobs)
}

// @Summary Create endpoint
// @Description Create endpoint
// @Tags Jobs
// @Accept json
// @Produce json
// @Param request body models.Job true "Payload"
// @Router /jobs [post]
func (h *JobHandler) Create(c echo.Context) error {
	var j models.Job
	if err := c.Bind(&j); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	created, err := h.jobService.CreateJob(c.Request().Context(), &j)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

// @Summary Get Job
// @Description Get Job
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Router /jobs/{id} [get]
func (h *JobHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	j, err := h.jobService.GetJob(c.Request().Context(), id)
	if err != nil || j == nil {
		return utils.Error(c, http.StatusNotFound, "job not found")
	}
	return utils.Success(c, "Operation successful", j)
}

// @Summary Delete Job
// @Description Delete Job
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Router /jobs/{id} [delete]
func (h *JobHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.jobService.DeleteJob(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary Run endpoint
// @Description Run endpoint
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /jobs/{id}/trigger [post]
func (h *JobHandler) Run(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	out, err := h.jobService.ExecuteJob(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "executed", "output": out})
}
