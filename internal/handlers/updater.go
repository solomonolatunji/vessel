package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/services"
)

type UpdaterHandler struct {
	updaterService *services.UpdaterService
}

func NewUpdaterHandler(s *services.UpdaterService) *UpdaterHandler {
	return &UpdaterHandler{updaterService: s}
}

func (h *UpdaterHandler) GetUpdateStatus(c echo.Context) error {
	if h.updaterService == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "updater service not initialized"})
	}
	status := h.updaterService.GetStatus()
	return c.JSON(http.StatusOK, status)
}

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
