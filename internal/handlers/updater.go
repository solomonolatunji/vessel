package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
)

type UpdaterHandler struct {
	updaterService *services.UpdaterService
}

func NewUpdaterHandler(s *services.UpdaterService) *UpdaterHandler {
	return &UpdaterHandler{updaterService: s}
}

// @Summary GetUpdateStatus endpoint
// @Description GetUpdateStatus endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/updates/status [get]
func (h *UpdaterHandler) GetUpdateStatus(c echo.Context) error {
	if h.updaterService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "updater service not initialized"})
	}
	status := h.updaterService.GetStatus()
	return c.JSON(http.StatusOK, status)
}

// @Summary CheckUpdate endpoint
// @Description CheckUpdate endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/updates/check [post]
func (h *UpdaterHandler) CheckUpdate(c echo.Context) error {
	if h.updaterService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "updater service not initialized"})
	}
	if _, err := h.updaterService.CheckForUpdates(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	status := h.updaterService.GetStatus()
	return c.JSON(http.StatusOK, status)
}

// @Summary DeployUpdate endpoint
// @Description DeployUpdate endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/updates/deploy [post]
func (h *UpdaterHandler) DeployUpdate(c echo.Context) error {
	if h.updaterService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "updater service not initialized"})
	}
	if err := h.updaterService.DeployUpdate(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "update deployment triggered",
	})
}
