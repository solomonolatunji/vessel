package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(s *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: s}
}

// @Summary List endpoint
// @Description List endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Router /api/workspaces [get]
func (h *TeamHandler) List(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	teams, err := h.teamService.ListTeamsByUser(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, teams)
}

// @Summary Create endpoint
// @Description Create endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Router /api/workspaces [post]
func (h *TeamHandler) Create(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "valid team name required"})
	}
	team, err := h.teamService.CreateTeam(c.Request().Context(), req.Name, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, team)
}

// @Summary Get endpoint
// @Description Get endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /api/teams/{teamId}/ai_settings [get]
func (h *TeamHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing team id"})
	}
	team, err := h.teamService.GetTeam(c.Request().Context(), id)
	if err != nil || team == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "team not found"})
	}
	return c.JSON(http.StatusOK, team)
}

// @Summary Delete endpoint
// @Description Delete endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/workspaces/{id} [delete]
func (h *TeamHandler) Delete(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing team id"})
	}
	if err := h.teamService.DeleteTeam(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary ListMembers endpoint
// @Description ListMembers endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/teams/{id}/members [get]
func (h *TeamHandler) ListMembers(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing team id"})
	}
	members, err := h.teamService.ListMembers(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, members)
}

// @Summary InviteMember endpoint
// @Description InviteMember endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/teams/{id}/invite [post]
func (h *TeamHandler) InviteMember(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing team id"})
	}
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "valid email required"})
	}
	inv, err := h.teamService.InviteMember(c.Request().Context(), id, req.Email, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, inv)
}

// @Summary RemoveMember endpoint
// @Description RemoveMember endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param userId path string true "userId"
// @Router /api/teams/{id}/members/{userId} [delete]
func (h *TeamHandler) RemoveMember(c echo.Context) error {
	id := c.Param("id")
	targetUserID := c.Param("userId")
	if id == "" || targetUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing team id or userId"})
	}
	if err := h.teamService.RemoveMember(c.Request().Context(), id, targetUserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary GetInvite endpoint
// @Description GetInvite endpoint
// @Tags Team-invites
// @Accept json
// @Produce json
// @Param token path string true "token"
// @Router /api/team-invites/{token} [get]
func (h *TeamHandler) GetInvite(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing invite token"})
	}
	inv, err := h.teamService.GetInvite(c.Request().Context(), token)
	if err != nil || inv == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invite not found or expired"})
	}
	return c.JSON(http.StatusOK, inv)
}

// @Summary AcceptInvite endpoint
// @Description AcceptInvite endpoint
// @Tags Team-invites
// @Accept json
// @Produce json
// @Param token path string true "token"
// @Router /api/team-invites/{token}/accept [post]
func (h *TeamHandler) AcceptInvite(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	token := c.Param("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing invite token"})
	}
	if err := h.teamService.AcceptInvite(c.Request().Context(), token, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "accepted"})
}
