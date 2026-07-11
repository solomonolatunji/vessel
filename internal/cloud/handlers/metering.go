package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"vessel.dev/vessel/internal/cloud/services"
)

type MeteringHandler struct {
	meteringService services.MeteringService
}

func NewMeteringHandler(meteringService services.MeteringService) *MeteringHandler {
	return &MeteringHandler{
		meteringService: meteringService,
	}
}

type UsageReport struct {
	TeamID         uint `json:"team_id"`
	Deployments    int  `json:"deployments"`
	ContainerHours int  `json:"container_hours"`
	BandwidthGB    int  `json:"bandwidth_gb"`
}

// ReportUsage handles incoming usage metrics from connected agents
// @Summary Report Usage Metrics
// @Description Receives telemetry from Vossel Daemons for billing
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/billing/usage/report [post]
func (h *MeteringHandler) ReportUsage(c echo.Context) error {
	var req UsageReport
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid usage report"})
	}

	err := h.meteringService.RecordUsage(req.TeamID, req.Deployments, req.ContainerHours, req.BandwidthGB)
	if err != nil {
		log.Printf("Error recording usage for team %d: %v", req.TeamID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to record usage"})
	}

	log.Printf("Successfully recorded usage report for team %d: %d deploys, %d hours", req.TeamID, req.Deployments, req.ContainerHours)

	return c.JSON(http.StatusOK, map[string]string{"status": "recorded"})
}
