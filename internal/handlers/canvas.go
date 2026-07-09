package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/services"
)

type CanvasHandler struct {
	canvasService *services.CanvasService
}

func NewCanvasHandler(s *services.CanvasService) *CanvasHandler {
	return &CanvasHandler{canvasService: s}
}

func (h *CanvasHandler) ListCanvasSummaries(c echo.Context) error {
	summaries, err := h.canvasService.ListSummaries(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, summaries)
}

func (h *CanvasHandler) GetCanvasSummary(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	summary, err := h.canvasService.GetSummary(c.Request().Context(), id)
	if err != nil || summary == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "canvas summary not found"})
	}
	return c.JSON(http.StatusOK, summary)
}

func (h *CanvasHandler) GetEnvironmentCanvas(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	canvas, err := h.canvasService.GetEnvironmentCanvas(c.Request().Context(), id)
	if err != nil || canvas == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "environment canvas not found"})
	}
	return c.JSON(http.StatusOK, canvas)
}
