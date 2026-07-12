package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(ns *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: ns}
}

type TestNotificationRequest struct {
	ChannelID   string `json:"channelId"`
	WorkspaceID string `json:"workspaceId"`
	Provider    string `json:"provider"`
}

// @Summary ListChannels endpoint
// @Description ListChannels endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /settings/notifications [get]
func (h *NotificationHandler) ListChannels(c echo.Context) error {
	if c.Request().Method != http.MethodGet {
		return utils.Error(c, http.StatusMethodNotAllowed, "Method not allowed")
	}
	workspaceID := c.QueryParam("workspaceId")
	if workspaceID == "" {
		workspaceID = "default"
	}
	channels, err := h.notificationService.ListChannels(c.Request().Context(), workspaceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", channels)
}

// @Summary SaveChannel endpoint
// @Description SaveChannel endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body models.WorkspaceNotificationChannel true "Payload"
// @Router /settings/notifications [put]
func (h *NotificationHandler) SaveChannel(c echo.Context) error {
	if c.Request().Method != http.MethodPut && c.Request().Method != http.MethodPost {
		return utils.Error(c, http.StatusMethodNotAllowed, "Method not allowed")
	}
	var channel models.WorkspaceNotificationChannel
	if err := c.Bind(&channel); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if channel.WorkspaceID == "" {
		channel.WorkspaceID = "default"
	}
	if err := h.notificationService.SaveChannel(c.Request().Context(), &channel); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", channel)
}

// @Summary DeleteChannel endpoint
// @Description DeleteChannel endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /settings/notifications/{id} [delete]
func (h *NotificationHandler) DeleteChannel(c echo.Context) error {
	if c.Request().Method != http.MethodDelete {
		return utils.Error(c, http.StatusMethodNotAllowed, "Method not allowed")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "Missing channel id")
	}
	if err := h.notificationService.DeleteChannel(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}

// @Summary TestNotification endpoint
// @Description TestNotification endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body handlers.TestNotificationRequest true "Payload"
// @Router /settings/notifications/test [post]
func (h *NotificationHandler) TestNotification(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return utils.Error(c, http.StatusMethodNotAllowed, "Method not allowed")
	}
	var req TestNotificationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if req.Provider != "" {
		err := h.notificationService.TestGlobalNotification(c.Request().Context(), req.Provider)
		if err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Global test notification queued",
		})
	}

	err := h.notificationService.TestTeamNotification(c.Request().Context(), req.WorkspaceID, req.ChannelID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Test notification queued",
	})
}
