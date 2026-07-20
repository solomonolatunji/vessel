package handlers

import (
	"os"
	"os/exec"
	"syscall"

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

func (h *SystemHandler) GetStats(c echo.Context) error {
	stats, err := h.service.GetStats()
	if err != nil {
		return utils.Error(c, 500, "failed to get system stats")
	}
	return utils.Success(c, "System stats", stats)
}

func (h *SystemHandler) Restart(c echo.Context) error {
	go func() {
		if _, err := exec.LookPath("docker"); err == nil {
			exec.Command("docker", "compose", "-f", "/vessl/docker-compose.yml", "restart", "vessl-control-plane").Start()
		} else {
			if p, err := os.FindProcess(os.Getpid()); err == nil {
				_ = p.Signal(syscall.SIGTERM)
			}
		}
	}()
	return utils.Success(c, "Restart initiated", map[string]string{"status": "restarting"})
}

func (h *SystemHandler) Cleanup(c echo.Context) error {
	if err := h.service.Cleanup(); err != nil {
		return utils.Error(c, 500, "Cleanup failed")
	}
	return utils.Success(c, "System cleanup completed successfully", nil)
}
