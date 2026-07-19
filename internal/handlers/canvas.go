package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/services"
)

type CanvasHandler struct {
	canvasService *services.CanvasService
}

func NewCanvasHandler(s *services.CanvasService) *CanvasHandler {
	return &CanvasHandler{canvasService: s}
}

// @Summary ListCanvasSummaries endpoint
// @Description ListCanvasSummaries endpoint
// @Tags Canvas
// @Accept json
// @Produce json
func (h *CanvasHandler) ListCanvasSummaries(c echo.Context) error {
	summaries, err := h.canvasService.ListSummaries(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if summaries == nil {
		summaries = make([]models.CanvasSummary, 0)
	}
	return utils.Success(c, "Operation successful", summaries)
}

// @Summary GetCanvasSummary endpoint
// @Description GetCanvasSummary endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /projects/{id}/summary [get]
func (h *CanvasHandler) GetCanvasSummary(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	summary, err := h.canvasService.GetSummary(c.Request().Context(), id)
	if err != nil || summary == nil {
		return utils.Error(c, http.StatusNotFound, "canvas summary not found")
	}
	return utils.Success(c, "Operation successful", summary)
}

// @Summary GetEnvironmentCanvas endpoint
// @Description GetEnvironmentCanvas endpoint
// @Tags Environments
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /environments/{id}/canvas [get]
func (h *CanvasHandler) GetEnvironmentCanvas(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	canvas, err := h.canvasService.GetEnvironmentCanvas(c.Request().Context(), id)
	if err != nil || canvas == nil {
		return utils.Error(c, http.StatusNotFound, "environment canvas not found")
	}
	return utils.Success(c, "Operation successful", canvas)
}
