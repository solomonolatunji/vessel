package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

type WorkspaceHandler struct {
	store *store.Store
}

func NewWorkspaceHandler(store *store.Store) *WorkspaceHandler {
	return &WorkspaceHandler{store: store}
}

func (h *WorkspaceHandler) ListTrustedDomains(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	list, err := h.store.ListWorkspaceTrustedDomains(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *WorkspaceHandler) CreateTrustedDomain(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	var item types.WorkspaceTrustedDomain
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item.TeamID = teamID

	if err := h.store.CreateWorkspaceTrustedDomain(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *WorkspaceHandler) DeleteTrustedDomain(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteWorkspaceTrustedDomain(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListSSHKeys(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	list, err := h.store.ListWorkspaceSSHKeys(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *WorkspaceHandler) CreateSSHKey(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	var item types.WorkspaceSSHKey
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item.TeamID = teamID

	if err := h.store.CreateWorkspaceSSHKey(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *WorkspaceHandler) DeleteSSHKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteWorkspaceSSHKey(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	list, err := h.store.ListWorkspaceAuditLogs(teamID, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *WorkspaceHandler) ListWorkspaces(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	ownerID := "default"
	if claims != nil {
		ownerID = claims.UserID
	}

	list, err := h.store.ListWorkspaces(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *WorkspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	ownerID := "default"
	if claims != nil {
		ownerID = claims.UserID
	}

	var item types.Workspace
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if item.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	item.OwnerID = ownerID
	if err := h.store.CreateWorkspace(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *WorkspaceHandler) GetWorkspace(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	item, err := h.store.GetWorkspace(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, "workspace not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *WorkspaceHandler) UpdateWorkspace(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	ownerID := "default"
	if claims != nil {
		ownerID = claims.UserID
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	existing, err := h.store.GetWorkspace(id)
	if err != nil || existing == nil {
		http.Error(w, "workspace not found", http.StatusNotFound)
		return
	}
	if existing.OwnerID != ownerID && ownerID != "default" {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}

	var item types.Workspace
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item.ID = id
	item.OwnerID = existing.OwnerID
	if item.Name == "" {
		item.Name = existing.Name
	}
	if err := h.store.UpdateWorkspace(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *WorkspaceHandler) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	ownerID := "default"
	if claims != nil {
		ownerID = claims.UserID
	}
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteWorkspace(id, ownerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListWorkspaceProjects(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	projects, err := h.store.ListProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var filtered []types.ProjectConfig
	for _, p := range projects {
		if p.WorkspaceID == id || p.TeamID == id {
			filtered = append(filtered, p)
		}
	}
	if filtered == nil {
		filtered = []types.ProjectConfig{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}
