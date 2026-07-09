package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/services"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(s *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: s}
}

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
