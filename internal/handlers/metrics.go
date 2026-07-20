package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type MetricsHandler struct {
	metricsService *services.MetricsService
}

func NewMetricsHandler(ms *services.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsService: ms}
}

func (h *MetricsHandler) GetHistoricalMetrics(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	rangeParam := c.QueryParam("range")
	if rangeParam == "" {
		rangeParam = "24h"
	}

	end := time.Now()
	var start time.Time
	var step string

	switch rangeParam {
	case "7d":
		start = end.Add(-7 * 24 * time.Hour)
		step = "1h"
	case "24h":
		start = end.Add(-24 * time.Hour)
		step = "5m"
	case "1h":
		start = end.Add(-1 * time.Hour)
		step = "1m"
	default:
		start = end.Add(-24 * time.Hour)
		step = "5m"
	}

	opts := services.ServiceMetricsOpts{
		ServiceID: serviceID,
		Start:     start,
		End:       end,
		Step:      step,
	}
	data, err := h.metricsService.GetServiceMetrics(c.Request().Context(), opts)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Historical metrics fetched", data)
}
