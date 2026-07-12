package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/services"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(s *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: s}
}

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type InviteTeamMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// @Summary List Teams
// @Description List Teams
// @Tags Teams
// @Accept json
// @Produce json
// @Router /teams [get]
func (h *TeamHandler) List(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	teams, err := h.teamService.ListTeamsByUser(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", teams)
}

// @Summary Create Team
// @Description Create Team
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body handlers.CreateTeamRequest true "Payload"
// @Router /teams [post]
func (h *TeamHandler) Create(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "valid team name required")
	}
	team, err := h.teamService.CreateTeam(c.Request().Context(), req.Name, userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", team)
}

// @Summary Get Team
// @Description Get Team
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Router /teams/{id} [get]
func (h *TeamHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing team id")
	}
	team, err := h.teamService.GetTeam(c.Request().Context(), id)
	if err != nil || team == nil {
		return utils.Error(c, http.StatusNotFound, "team not found")
	}
	return utils.Success(c, "Operation successful", team)
}

// @Summary Delete Team
// @Description Delete Team
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Router /teams/{id} [delete]
func (h *TeamHandler) Delete(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing team id")
	}
	if err := h.teamService.DeleteTeam(c.Request().Context(), id, userID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary ListMembers endpoint
// @Description ListMembers endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /teams/{id}/members [get]
func (h *TeamHandler) ListMembers(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing team id")
	}
	members, err := h.teamService.ListMembers(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", members)
}

// @Summary InviteMember endpoint
// @Description InviteMember endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param request body handlers.InviteTeamMemberRequest true "Payload"
// @Router /teams/{id}/invite [post]
func (h *TeamHandler) InviteMember(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing team id")
	}
	var req InviteTeamMemberRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "valid email required")
	}
	inv, err := h.teamService.InviteMember(c.Request().Context(), id, req.Email, req.Role)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", inv)
}

// @Summary RemoveMember endpoint
// @Description RemoveMember endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param userId path string true "userId"
// @Router /teams/{id}/members/{userId} [delete]
func (h *TeamHandler) RemoveMember(c echo.Context) error {
	id := c.Param("id")
	targetUserID := c.Param("userId")
	if id == "" || targetUserID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing team id or userId")
	}
	if err := h.teamService.RemoveMember(c.Request().Context(), id, targetUserID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary GetInvite endpoint
// @Description GetInvite endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param token path string true "token"
// @Router /team-invites/{token} [get]
func (h *TeamHandler) GetInvite(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "missing invite token")
	}
	inv, err := h.teamService.GetInvite(c.Request().Context(), token)
	if err != nil || inv == nil {
		return utils.Error(c, http.StatusNotFound, "invite not found or expired")
	}
	return utils.Success(c, "Operation successful", inv)
}

// @Summary AcceptInvite endpoint
// @Description AcceptInvite endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param token path string true "token"
// @Router /team-invites/{token}/accept [post]
func (h *TeamHandler) AcceptInvite(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	token := c.Param("token")
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "missing invite token")
	}
	if err := h.teamService.AcceptInvite(c.Request().Context(), token, userID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "accepted"})
}
