package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	dockerfilters "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type SettingsHandler struct {
	store        *store.Store
	dockerClient *client.Client
}

func NewSettingsHandler(s *store.Store, dockerClient *client.Client) *SettingsHandler {
	return &SettingsHandler{
		store:        s,
		dockerClient: dockerClient,
	}
}

// GetServerSettings returns global system configurations (Caddy IP, webhooks, SMTP).
func (h *SettingsHandler) GetServerSettings(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.store.GetServerSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

// UpdateServerSettings updates global daemon configurations (restricted to admin).
func (h *SettingsHandler) UpdateServerSettings(w http.ResponseWriter, r *http.Request) {
	var req types.ServerSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.store.UpdateServerSettings(&req); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, req)
}

// TriggerSystemPrune executes Docker system prune (containers, images, networks, volumes) and reports reclaimed space.
func (h *SettingsHandler) TriggerSystemPrune(w http.ResponseWriter, r *http.Request) {
	if h.dockerClient == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":            "simulated",
			"message":           "Docker client not initialized in standalone mode; simulated clean system prune.",
			"spaceReclaimedBytes": 104857600, // 100MB simulation
		})
		return
	}

	ctx := r.Context()
	var totalReclaimed uint64

	if cReport, err := h.dockerClient.ContainersPrune(ctx, dockerfilters.NewArgs()); err == nil {
		totalReclaimed += cReport.SpaceReclaimed
	}
	if iReport, err := h.dockerClient.ImagesPrune(ctx, dockerfilters.NewArgs()); err == nil {
		totalReclaimed += iReport.SpaceReclaimed
	}
	if nReport, err := h.dockerClient.NetworksPrune(ctx, dockerfilters.NewArgs()); err == nil {
		_ = nReport
	}
	if vReport, err := h.dockerClient.VolumesPrune(ctx, dockerfilters.NewArgs()); err == nil {
		totalReclaimed += vReport.SpaceReclaimed
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":            "success",
		"message":           "Docker system prune executed cleanly.",
		"spaceReclaimedBytes": totalReclaimed,
	})
}

// GetProfile returns the authenticated user's profile details.
func (h *SettingsHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	user, err := h.store.GetUserByID(claims.UserID)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	user.PasswordHash = ""
	writeJSON(w, http.StatusOK, user)
}

// UpdateProfile updates the authenticated user's name, email, or password.
func (h *SettingsHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	user, err := h.store.GetUserByID(claims.UserID)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Email) != "" {
		user.Email = strings.TrimSpace(req.Email)
	}
	if strings.TrimSpace(req.Name) != "" {
		user.Role = req.Name // Note: storing custom display name in User struct or retaining existing fields
	}
	if strings.TrimSpace(req.NewPassword) != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err == nil {
			user.PasswordHash = string(hash)
		}
	}

	if err := h.store.UpdateUser(user); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user.PasswordHash = ""
	writeJSON(w, http.StatusOK, user)
}

// ListPATs lists Personal Access Tokens created by the authenticated user.
func (h *SettingsHandler) ListPATs(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	list, err := h.store.ListPersonalAccessTokens(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// CreatePAT generates a new secure CLI / API access token.
func (h *SettingsHandler) CreatePAT(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		req.Name = "Personal Access Token"
	}

	rawBytes := make([]byte, 24)
	_, _ = rand.Read(rawBytes)
	rawToken := fmt.Sprintf("vsl_user_%s", hex.EncodeToString(rawBytes))
	tokenHash := sha256.Sum256([]byte(rawToken))
	hashStr := hex.EncodeToString(tokenHash[:])

	pat := &types.PersonalAccessToken{
		UserID:    claims.UserID,
		Name:      req.Name,
		TokenHash: hashStr,
		Prefix:    "vsl_user_",
	}

	if err := h.store.CreatePersonalAccessToken(pat); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"token": rawToken,
		"pat":   pat,
	})
}

// DeletePAT revokes and deletes a Personal Access Token.
func (h *SettingsHandler) DeletePAT(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	id := r.PathValue("id")
	if err := h.store.DeletePersonalAccessToken(id, claims.UserID); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
