package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type NotificationSettingsHandler struct {
	notifSettingsService *services.NotificationSettingsService
}

func NewNotificationSettingsHandler(s *services.NotificationSettingsService) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{notifSettingsService: s}
}

func (h *NotificationSettingsHandler) GetNotificationSettings(c echo.Context) error {
	s, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", s)
}

func (h *NotificationSettingsHandler) UpdateNotificationSettings(c echo.Context) error {
	var req models.UpdateNotificationSettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := h.notifSettingsService.UpdateNotificationSettings(c.Request().Context(), &req.NotificationSettings); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	updated, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Notification settings updated successfully", updated)
}
