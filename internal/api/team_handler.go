package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

type TeamHandler struct {
	store *store.Store
}

func NewTeamHandler(s *store.Store) *TeamHandler {
	return &TeamHandler{store: s}
}

func (h *TeamHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teams, err := h.store.ListUserTeams(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, teams)
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	var req types.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "team name is required")
		return
	}

	req.OwnerID = claims.UserID
	if err := h.store.CreateTeam(&req); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, req)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teamID := r.PathValue("id")
	team, err := h.store.GetTeam(teamID)
	if err != nil || team == nil {
		writeError(w, http.StatusNotFound, "team not found")
		return
	}

	member, _ := h.store.GetTeamMember(teamID, claims.UserID)
	if member == nil && claims.Role != "admin" {
		writeError(w, http.StatusForbidden, "not a member of this team")
		return
	}

	members, _ := h.store.ListTeamMembers(teamID)
	writeJSON(w, http.StatusOK, map[string]any{
		"team":    team,
		"members": members,
	})
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teamID := r.PathValue("id")
	if err := h.store.DeleteTeam(teamID, claims.UserID); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("id")
	members, err := h.store.ListTeamMembers(teamID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, members)
}

func (h *TeamHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teamID := r.PathValue("id")
	callerMember, _ := h.store.GetTeamMember(teamID, claims.UserID)
	if callerMember == nil || (callerMember.Role != "Owner" && callerMember.Role != "Admin") {
		if claims.Role != "admin" {
			writeError(w, http.StatusForbidden, "only Owner or Admin can invite members")
			return
		}
	}

	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid invite payload")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Role == "" {
		req.Role = "Member"
	}

	existingUser, _ := h.store.GetUserByEmail(req.Email)
	if existingUser != nil {
		member := &types.TeamMember{
			TeamID:    teamID,
			UserID:    existingUser.ID,
			UserEmail: existingUser.Email,
			Role:      req.Role,
		}
		if err := h.store.AddTeamMember(member); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"status": "added_to_team",
			"member": member,
		})
		return
	}

	invite := &types.TeamInvite{
		TeamID:    teamID,
		Email:     req.Email,
		Role:      req.Role,
		InvitedBy: claims.Email,
	}
	if err := h.store.CreateTeamInvite(invite); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"status": "invitation_sent",
		"invite": invite,
	})
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	teamID := r.PathValue("id")
	targetUserID := r.PathValue("userId")

	callerMember, _ := h.store.GetTeamMember(teamID, claims.UserID)
	if callerMember == nil || (callerMember.Role != "Owner" && callerMember.Role != "Admin" && claims.UserID != targetUserID) {
		if claims.Role != "admin" {
			writeError(w, http.StatusForbidden, "unauthorized to remove this member")
			return
		}
	}

	if err := h.store.RemoveTeamMember(teamID, targetUserID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) GetInvite(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	inv, err := h.store.GetTeamInviteByToken(token)
	if err != nil || inv == nil {
		writeError(w, http.StatusNotFound, "invitation not found or expired")
		return
	}
	writeJSON(w, http.StatusOK, inv)
}

func (h *TeamHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	token := r.PathValue("token")
	inv, err := h.store.GetTeamInviteByToken(token)
	if err != nil || inv == nil {
		writeError(w, http.StatusNotFound, "invitation not found or expired")
		return
	}

	member := &types.TeamMember{
		TeamID:    inv.TeamID,
		UserID:    claims.UserID,
		UserEmail: claims.Email,
		Role:      inv.Role,
	}
	if err := h.store.AddTeamMember(member); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = h.store.DeleteTeamInvite(inv.ID)

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "accepted",
		"member": member,
	})
}
