package handlers

import (
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// @Summary Restart the Vessl daemon
// @Description Triggers a restart of the Vessl daemon (admin only)
// @Tags System
// @Produce json
// @Success 200 {object} utils.BaseResponse
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
