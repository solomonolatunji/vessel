package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

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
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
func (h *WorkspaceHandler) List(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	wsList, total, err := h.workspaceService.ListWorkspaces(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Paginated(c, "Operation successful", wsList, total, page, limit)
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
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var payload CreateWorkspaceRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	ws, err := h.workspaceService.CreateWorkspace(c.Request().Context(), payload.Name, userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", ws)
}

// @Summary Get endpoint
// @Description Get endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /workspaces/{id} [get]
func (h *WorkspaceHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	ws, err := h.workspaceService.GetWorkspace(c.Request().Context(), id)
	if err != nil || ws == nil {
		return utils.Error(c, http.StatusNotFound, "workspace not found")
	}
	return utils.Success(c, "Operation successful", ws)
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
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	var ws models.Workspace
	if err := c.Bind(&ws); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	ws.ID = id
	if err := h.workspaceService.UpdateWorkspace(c.Request().Context(), &ws); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", ws)
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
		return utils.Error(c, http.StatusBadRequest, "missing parameters or unauthorized")
	}
	if err := h.workspaceService.DeleteWorkspace(c.Request().Context(), id, userID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary ListTrustedDomains endpoint
// @Description ListTrustedDomains endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Router /workspaces/{workspaceId}/trusted-domains [get]
func (h *WorkspaceHandler) ListTrustedDomains(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}
	domains, err := h.workspaceService.ListTrustedDomains(c.Request().Context(), workspaceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", domains)
}

// @Summary CreateTrustedDomain endpoint
// @Description CreateTrustedDomain endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Param request body handlers.CreateTrustedDomainRequest true "Payload"
// @Router /workspaces/{workspaceId}/trusted-domains [post]
func (h *WorkspaceHandler) CreateTrustedDomain(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}
	var payload CreateTrustedDomainRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	td, err := h.workspaceService.AddTrustedDomain(c.Request().Context(), workspaceID, payload.Domain)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", td)
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
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.workspaceService.DeleteTrustedDomain(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary ListSSHKeys endpoint
// @Description ListSSHKeys endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Router /workspaces/{workspaceId}/ssh-keys [get]
func (h *WorkspaceHandler) ListSSHKeys(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}
	keys, err := h.workspaceService.ListSSHKeys(c.Request().Context(), workspaceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", keys)
}

// @Summary CreateSSHKey endpoint
// @Description CreateSSHKey endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Param request body handlers.CreateSSHKeyRequest true "Payload"
// @Router /workspaces/{workspaceId}/ssh-keys [post]
func (h *WorkspaceHandler) CreateSSHKey(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}
	var payload CreateSSHKeyRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	key, err := h.workspaceService.AddSSHKey(c.Request().Context(), workspaceID, payload.Name, payload.PublicKey)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", key)
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
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if err := h.workspaceService.DeleteSSHKey(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary ListAuditLogs endpoint
// @Description ListAuditLogs endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param workspaceId path string true "workspaceId"
// @Router /workspaces/{workspaceId}/audit-logs [get]
func (h *WorkspaceHandler) ListAuditLogs(c echo.Context) error {
	workspaceID := c.Param("workspaceId")
	if workspaceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspaceId parameter")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	logs, total, err := h.workspaceService.ListAuditLogs(c.Request().Context(), workspaceID, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Paginated(c, "Operation successful", logs, total, page, limit)
}

func (h *WorkspaceHandler) ListMembers(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspace id")
	}
	members, err := h.workspaceService.ListMembers(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", members)
}

func (h *WorkspaceHandler) InviteMember(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspace id")
	}
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "valid email required")
	}
	inv, err := h.workspaceService.InviteMember(c.Request().Context(), id, req.Email, req.Role)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", inv)
}

func (h *WorkspaceHandler) RemoveMember(c echo.Context) error {
	id := c.Param("id")
	targetUserID := c.Param("userId")
	if id == "" || targetUserID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing workspace id or userId")
	}
	if err := h.workspaceService.RemoveMember(c.Request().Context(), id, targetUserID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *WorkspaceHandler) GetInvite(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "missing invite token")
	}
	inv, err := h.workspaceService.GetInvite(c.Request().Context(), token)
	if err != nil || inv == nil {
		return utils.Error(c, http.StatusNotFound, "invite not found or expired")
	}
	return utils.Success(c, "Operation successful", inv)
}

func (h *WorkspaceHandler) AcceptInvite(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	token := c.Param("token")
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "missing invite token")
	}
	if err := h.workspaceService.AcceptInvite(c.Request().Context(), token, userID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "accepted"})
}
