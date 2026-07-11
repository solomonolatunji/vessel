package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(ns *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: ns}
}

// @Summary ListChannels endpoint
// @Description ListChannels endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/notifications [get]
func (h *NotificationHandler) ListChannels(c echo.Context) error {
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
	teamID := c.QueryParam("teamId")
	if teamID == "" {
		teamID = "default"
	}
	channels, err := h.notificationService.ListChannels(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channels)
}

// @Summary SaveChannel endpoint
// @Description SaveChannel endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/notifications [put]
func (h *NotificationHandler) SaveChannel(c echo.Context) error {
	if c.Request().Method != http.MethodPut && c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
	var channel models.TeamNotificationChannel
	if err := c.Bind(&channel); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if channel.TeamID == "" {
		channel.TeamID = "default"
	}
	if err := h.notificationService.SaveChannel(c.Request().Context(), &channel); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channel)
}

// @Summary DeleteChannel endpoint
// @Description DeleteChannel endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/notifications/{id} [delete]
func (h *NotificationHandler) DeleteChannel(c echo.Context) error {
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing channel id"})
	}
	if err := h.notificationService.DeleteChannel(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// @Summary TestNotification endpoint
// @Description TestNotification endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/notifications/test [post]
func (h *NotificationHandler) TestNotification(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
	var req struct {
		ChannelID string `json:"channelId"`
		TeamID    string `json:"teamId"`
		Provider  string `json:"provider"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if req.Provider != "" {
		err := h.notificationService.TestGlobalNotification(c.Request().Context(), req.Provider)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Global test notification queued",
		})
	}

	err := h.notificationService.TestTeamNotification(c.Request().Context(), req.TeamID, req.ChannelID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Test notification queued",
	})
}
