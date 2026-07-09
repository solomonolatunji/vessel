package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(ns *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: ns}
}

func (h *NotificationHandler) GetIntegrations(c echo.Context) error {
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	integ, err := h.notificationService.GetIntegration(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, integ)
}

func (h *NotificationHandler) SaveIntegrations(c echo.Context) error {
	if c.Request().Method != http.MethodPut && c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	var integ models.NotificationIntegration
	if err := c.Bind(&integ); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if err := h.notificationService.SaveIntegration(c.Request().Context(), &integ); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, integ)
}

func (h *NotificationHandler) TestNotification(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	var req struct {
		Channel   string `json:"channel"`
		ProjectID string `json:"projectId,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if err := h.notificationService.SendTest(req.Channel, req.ProjectID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Test notification sent successfully over " + req.Channel,
	})
}

func (h *NotificationHandler) GetProjectPreferences(c echo.Context) error {
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing project id parameter"})
	}

	pref, err := h.notificationService.GetProjectPref(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, pref)
}

func (h *NotificationHandler) SaveProjectPreferences(c echo.Context) error {
	if c.Request().Method != http.MethodPut && c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing project id parameter"})
	}

	var pref models.ProjectNotificationPref
	if err := c.Bind(&pref); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	pref.ProjectID = projectID

	if err := h.notificationService.SaveProjectPref(c.Request().Context(), &pref); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, pref)
}
