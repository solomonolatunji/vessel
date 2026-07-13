package handlers

import (
	"log"
	"net/http"

	"strings"

	"github.com/labstack/echo/v4"
	repos "vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/cloud/services"
	"vessl.dev/vessl/internal/utils"
)

type MeteringHandler struct {
	repo            repos.CloudRepo
	meteringService services.MeteringService
}

func NewMeteringHandler(repo repos.CloudRepo, meteringService services.MeteringService) *MeteringHandler {
	return &MeteringHandler{
		repo:            repo,
		meteringService: meteringService,
	}
}

type UsageReport struct {
	Deployments    int `json:"deployments"`
	ContainerHours int `json:"container_hours"`
	BandwidthGB    int `json:"bandwidth_gb"`
}

// ReportUsage handles incoming usage metrics from connected agents
// @Summary Report Usage Metrics
// @Description Receives telemetry from Vossel Daemons for billing
// @Tags Cloud-Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /billing/usage/report [post]
func (h *MeteringHandler) ReportUsage(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return utils.Error(c, http.StatusUnauthorized, "Missing Authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	var req UsageReport
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid usage report")
	}

	server, err := h.repo.GetServerByToken(token)
	if err != nil || server == nil {
		return utils.Error(c, http.StatusUnauthorized, "Invalid or unknown Agent token")
	}

	err = h.meteringService.RecordUsage(server.WorkspaceID, req.Deployments, req.ContainerHours, req.BandwidthGB)
	if err != nil {
		log.Printf("Error recording usage for team %d: %v", server.WorkspaceID, err)
		return utils.Error(c, http.StatusInternalServerError, "Failed to record usage")
	}

	log.Printf("Successfully recorded usage report for team %d: %d deploys, %d hours", server.WorkspaceID, req.Deployments, req.ContainerHours)

	return utils.Success(c, "recorded", nil)
}
