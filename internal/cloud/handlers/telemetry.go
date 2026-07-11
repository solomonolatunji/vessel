package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"vessel.dev/vessel/internal/cloud/repos"
	"vessel.dev/vessel/internal/models"
)

type TelemetryHandler struct {
	repo repos.CloudRepo
}

func NewTelemetryHandler(repo repos.CloudRepo) *TelemetryHandler {
	return &TelemetryHandler{repo: repo}
}

type TelemetryPayload struct {
	InstanceID    string `json:"instance_id"`
	Version       string `json:"version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	ActiveServers int    `json:"active_servers"`
	ActiveApps    int    `json:"active_apps"`
}

// @Summary Receive Telemetry Ping
// @Description Receives anonymous telemetry pings from OSS instances
// @Tags Cloud-Telemetry
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/telemetry/ping [post]
func (h *TelemetryHandler) ReceivePing(c echo.Context) error {
	var payload TelemetryPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid telemetry payload"})
	}

	if payload.InstanceID == "" || payload.Version == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	logEntry := &models.CloudTelemetryLog{
		InstanceID:    payload.InstanceID,
		Version:       payload.Version,
		OS:            payload.OS,
		Arch:          payload.Arch,
		ActiveServers: payload.ActiveServers,
		ActiveApps:    payload.ActiveApps,
		ReportedAt:    time.Now(),
	}

	if err := h.repo.LogTelemetry(logEntry); err != nil {
		log.Printf("Failed to store telemetry ping: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
