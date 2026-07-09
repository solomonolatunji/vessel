package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type ProjectSettingsHandler struct {
	settingsService *services.ProjectSettingsService
}

func NewProjectSettingsHandler(s *services.ProjectSettingsService) *ProjectSettingsHandler {
	return &ProjectSettingsHandler{settingsService: s}
}

func (h *ProjectSettingsHandler) ListWebhooks(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId"})
	}
	list, err := h.settingsService.ListWebhooks(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *ProjectSettingsHandler) CreateWebhook(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId"})
	}
	var req models.Webhook
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	req.ProjectID = projectID
	created, err := h.settingsService.CreateWebhook(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

func (h *ProjectSettingsHandler) DeleteWebhook(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId or id"})
	}
	if err := h.settingsService.DeleteWebhook(c.Request().Context(), id, projectID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ProjectSettingsHandler) ListTokens(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId"})
	}
	list, err := h.settingsService.ListTokens(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *ProjectSettingsHandler) CreateToken(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId"})
	}
	var req models.ProjectToken
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	req.ProjectID = projectID
	token, raw, err := h.settingsService.CreateToken(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"id":        token.ID,
		"name":      token.Name,
		"token":     raw,
		"createdAt": token.CreatedAt,
	})
}

func (h *ProjectSettingsHandler) DeleteToken(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId or id"})
	}
	if err := h.settingsService.DeleteToken(c.Request().Context(), id, projectID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ProjectSettingsHandler) ListMembers(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId"})
	}
	list, err := h.settingsService.ListMembers(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *ProjectSettingsHandler) AddMember(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId"})
	}
	var req models.ProjectMember
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	req.ProjectID = projectID
	added, err := h.settingsService.AddMember(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, added)
}

func (h *ProjectSettingsHandler) RemoveMember(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId or id"})
	}
	if err := h.settingsService.RemoveMember(c.Request().Context(), id, projectID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
