package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type WorkspaceHandler struct {
	workspaceService *services.WorkspaceService
}

func NewWorkspaceHandler(s *services.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceService: s}
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

type CreateTrustedDomainRequest struct {
	Domain string `json:"domain"`
}

type CreateSSHKeyRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

// @Summary List endpoint
// @Description List endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
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

// @Summary Create endpoint
// @Description Create endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param request body handlers.CreateWorkspaceRequest true "Payload"
// @Router /workspaces [post]
func (h *WorkspaceHandler) Create(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var payload CreateWorkspaceRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	ws, err := h.workspaceService.CreateWorkspace(c.Request().Context(), payload.Name, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, ws)
}

// @Summary Get endpoint
// @Description Get endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /workspaces/{id} [get]
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

// @Summary Update endpoint
// @Description Update endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param request body models.Workspace true "Payload"
// @Router /workspaces/{id} [put]
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

// @Summary Delete endpoint
// @Description Delete endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /workspaces/{id} [delete]
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

// @Summary ListTrustedDomains endpoint
// @Description ListTrustedDomains endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /teams/{teamId}/trusted-domains [get]
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

// @Summary CreateTrustedDomain endpoint
// @Description CreateTrustedDomain endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Param request body handlers.CreateTrustedDomainRequest true "Payload"
// @Router /teams/{teamId}/trusted-domains [post]
func (h *WorkspaceHandler) CreateTrustedDomain(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing teamId parameter"})
	}
	var payload CreateTrustedDomainRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	td, err := h.workspaceService.AddTrustedDomain(c.Request().Context(), teamID, payload.Domain)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, td)
}

// @Summary DeleteTrustedDomain endpoint
// @Description DeleteTrustedDomain endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
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

// @Summary ListSSHKeys endpoint
// @Description ListSSHKeys endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /teams/{teamId}/ssh-keys [get]
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

// @Summary CreateSSHKey endpoint
// @Description CreateSSHKey endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Param request body handlers.CreateSSHKeyRequest true "Payload"
// @Router /teams/{teamId}/ssh-keys [post]
func (h *WorkspaceHandler) CreateSSHKey(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing teamId parameter"})
	}
	var payload CreateSSHKeyRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	key, err := h.workspaceService.AddSSHKey(c.Request().Context(), teamID, payload.Name, payload.PublicKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, key)
}

// @Summary DeleteSSHKey endpoint
// @Description DeleteSSHKey endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
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

// @Summary ListAuditLogs endpoint
// @Description ListAuditLogs endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /teams/{teamId}/audit-logs [get]
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
