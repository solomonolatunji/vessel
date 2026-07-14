package handlers

import (
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type SystemHandler struct {
	service *services.SystemService
}

func NewSystemHandler(s *services.SystemService) *SystemHandler {
	return &SystemHandler{service: s}
}

// @Summary Get system stats
// @Description Returns CPU, memory, disk usage, and uptime for the host server
// @Tags System
// @Produce json
// @Success 200 {object} models.SystemStats
// @Router /system/stats [get]
func (h *SystemHandler) GetStats(c echo.Context) error {
	stats, err := h.service.GetStats()
	if err != nil {
		return utils.Error(c, 500, "failed to get system stats")
	}
	return utils.Success(c, "System stats", stats)
}

// @Summary Restart the Vessl daemon
// @Description Triggers a restart of the Vessl daemon (admin only)
// @Tags System
// @Produce json
// @Success 200 {object} map[string]any
// @Router /system/restart [post]
func (h *SystemHandler) Restart(c echo.Context) error {
	go func() {
		if _, err := exec.LookPath("docker"); err == nil {
			exec.Command("docker", "compose", "-f", "/vessl/docker-compose.yml", "restart", "vessl-control-plane").Start()
		} else {
			os.Exit(0)
		}
	}()
	return utils.Success(c, "Restart initiated", map[string]string{"status": "restarting"})
}
