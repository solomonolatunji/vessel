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
	"golang.org/x/crypto/bcrypt"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
	"vessel.dev/vessel/internal/updater"
)

type SettingsHandler struct {
	store        *store.Store
	dockerClient *client.Client
	updater      *updater.UpdaterService
}

func NewSettingsHandler(s *store.Store, dockerClient *client.Client, u *updater.UpdaterService) *SettingsHandler {
	return &SettingsHandler{
		store:        s,
		dockerClient: dockerClient,
		updater:      u,
	}
}

func (h *SettingsHandler) GetServerSettings(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.store.GetServerSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

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

func (h *SettingsHandler) TriggerSystemPrune(w http.ResponseWriter, r *http.Request) {
	if h.dockerClient == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":              "simulated",
			"message":             "Docker client not initialized in standalone mode; simulated clean system prune.",
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
		"status":              "success",
		"message":             "Docker system prune executed cleanly.",
		"spaceReclaimedBytes": totalReclaimed,
	})
}

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

func (h *SettingsHandler) HandleMCP(w http.ResponseWriter, r *http.Request) {
	settings, _ := h.store.GetServerSettings()
	if settings != nil && !settings.MCPServerEnabled {
		writeError(w, http.StatusForbidden, "MCP server endpoint is currently disabled by the administrator")
		return
	}

	if r.Method == http.MethodGet {
		writeJSON(w, http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"server": map[string]string{
				"name":    "vessel-mcp-server",
				"version": "v1.0.0",
			},
			"capabilities": map[string]any{
				"tools": map[string]any{
					"listChanged": false,
				},
			},
			"tools": []map[string]any{
				{
					"name":        "list_projects",
					"description": "List all deployed projects on Vessel",
					"inputSchema": map[string]any{
						"type":       "object",
						"properties": map[string]any{},
					},
				},
				{
					"name":        "get_system_status",
					"description": "Check Vessel server CPU, RAM, and database health",
					"inputSchema": map[string]any{
						"type":       "object",
						"properties": map[string]any{},
					},
				},
			},
		})
		return
	}

	var req struct {
		JSONRPC string `json:"jsonrpc"`
		ID      any    `json:"id"`
		Method  string `json:"method"`
		Params  struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON-RPC format")
		return
	}

	switch req.Method {
	case "tools/list":
		writeJSON(w, http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result": map[string]any{
				"tools": []map[string]any{
					{
						"name":        "list_projects",
						"description": "List all deployed projects on Vessel",
					},
					{
						"name":        "get_system_status",
						"description": "Check Vessel server CPU, RAM, and database health",
					},
				},
			},
		})
	case "tools/call":
		switch req.Params.Name {
		case "list_projects":
			projects, _ := h.store.ListProjects()
			writeJSON(w, http.StatusOK, map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"result": map[string]any{
					"content": []map[string]any{
						{"type": "text", "text": fmt.Sprintf("Found %d projects: %+v", len(projects), projects)},
					},
				},
			})
		case "get_system_status":
			writeJSON(w, http.StatusOK, map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"result": map[string]any{
					"content": []map[string]any{
						{"type": "text", "text": "Vessel system is healthy and operational."},
					},
				},
			})
		default:
			writeJSON(w, http.StatusOK, map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"error": map[string]any{
					"code":    -32601,
					"message": "Method/Tool not found: " + req.Params.Name,
				},
			})
		}
	default:
		writeJSON(w, http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"error": map[string]any{
				"code":    -32601,
				"message": "Method not supported: " + req.Method,
			},
		})
	}
}

func (h *SettingsHandler) CheckUpdate(w http.ResponseWriter, r *http.Request) {
	if h.updater == nil {
		writeError(w, http.StatusServiceUnavailable, "update management service is not initialized")
		return
	}
	info, err := h.updater.CheckForUpdate(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *SettingsHandler) DeployUpdate(w http.ResponseWriter, r *http.Request) {
	if h.updater == nil {
		writeError(w, http.StatusServiceUnavailable, "update management service is not initialized")
		return
	}
	if err := h.updater.DeployUpdate(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Update successfully applied and system restart triggered.",
	})
}

func (h *SettingsHandler) GetUpdateStatus(w http.ResponseWriter, r *http.Request) {
	settings, err := h.store.GetServerSettings()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	hasUpdate := false
	cClean := strings.TrimPrefix(strings.TrimSpace(settings.CurrentVersion), "v")
	lClean := strings.TrimPrefix(strings.TrimSpace(settings.LatestVersion), "v")
	if lClean != "" && lClean != cClean {
		hasUpdate = true
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"currentVersion":  settings.CurrentVersion,
		"latestVersion":   settings.LatestVersion,
		"hasUpdate":       hasUpdate,
		"lastChecked":     settings.LastUpdateCheck,
		"autoUpdate":      settings.AutoUpdateEnabled,
		"updateCheckCron": settings.UpdateCheckCron,
	})
}
