package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type WorkspaceHandler struct {
	workspaceService *services.WorkspaceService
}

func NewWorkspaceHandler(s *services.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceService: s}
}

func (h *WorkspaceHandler) List(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	wsList, err := h.workspaceService.ListWorkspaces(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, wsList)
}

func (h *WorkspaceHandler) Create(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	ws, err := h.workspaceService.CreateWorkspace(c.Request().Context(), payload.Name, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, ws)
}

func (h *WorkspaceHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	ws, err := h.workspaceService.GetWorkspace(c.Request().Context(), id)
	if err != nil || ws == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "workspace not found"})
	}
	return c.JSON(http.StatusOK, ws)
}

func (h *WorkspaceHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	var ws models.Workspace
	if err := c.Bind(&ws); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	ws.ID = id
	if err := h.workspaceService.UpdateWorkspace(c.Request().Context(), &ws); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, ws)
}

func (h *WorkspaceHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	userID := ExtractUserID(c)
	if id == "" || userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing parameters or unauthorized"})
	}
	if err := h.workspaceService.DeleteWorkspace(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListTrustedDomains(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing teamId parameter"})
	}
	domains, err := h.workspaceService.ListTrustedDomains(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, domains)
}

func (h *WorkspaceHandler) CreateTrustedDomain(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing teamId parameter"})
	}
	var payload struct {
		Domain string `json:"domain"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	td, err := h.workspaceService.AddTrustedDomain(c.Request().Context(), teamID, payload.Domain)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, td)
}

func (h *WorkspaceHandler) DeleteTrustedDomain(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	if err := h.workspaceService.DeleteTrustedDomain(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListSSHKeys(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing teamId parameter"})
	}
	keys, err := h.workspaceService.ListSSHKeys(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, keys)
}

func (h *WorkspaceHandler) CreateSSHKey(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing teamId parameter"})
	}
	var payload struct {
		Name      string `json:"name"`
		PublicKey string `json:"publicKey"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	key, err := h.workspaceService.AddSSHKey(c.Request().Context(), teamID, payload.Name, payload.PublicKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, key)
}

func (h *WorkspaceHandler) DeleteSSHKey(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	if err := h.workspaceService.DeleteSSHKey(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListAuditLogs(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing teamId parameter"})
	}
	logs, err := h.workspaceService.ListAuditLogs(c.Request().Context(), teamID, 100)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, logs)
}
