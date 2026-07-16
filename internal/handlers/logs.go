package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type LogHandler struct {
	logService *services.LogService
}

func NewLogHandler(ls *services.LogService) *LogHandler {
	return &LogHandler{logService: ls}
}

// @Summary Get Historical Logs
// @Description Fetches historical logs from Loki for a specific service
// @Tags Logs
// @Accept json
// @Produce json
// @Param serviceId path string true "Service ID"
// @Param range query string false "Time range (e.g., 24h, 7d)"
// @Param limit query int false "Max lines to fetch"
// @Router /services/{serviceId}/logs/historical [get]
func (h *LogHandler) GetHistoricalLogs(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	rangeParam := c.QueryParam("range")
	if rangeParam == "" {
		rangeParam = "24h"
	}

	limitParam := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 1000 // default limit
	}

	end := time.Now()
	var start time.Time

	switch rangeParam {
	case "7d":
		start = end.Add(-7 * 24 * time.Hour)
	case "24h":
		start = end.Add(-24 * time.Hour)
	case "1h":
		start = end.Add(-1 * time.Hour)
	default:
		start = end.Add(-24 * time.Hour)
	}

	logs, err := h.logService.GetHistoricalLogs(c.Request().Context(), serviceID, start, end, limit)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Historical logs fetched", logs)
}
