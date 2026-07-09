package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

type ProjectSettingsHandler struct {
	store *store.Store
}

func NewProjectSettingsHandler(store *store.Store) *ProjectSettingsHandler {
	return &ProjectSettingsHandler{store: store}
}

func (h *ProjectSettingsHandler) GetProjectBilling(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"projectId":      projectID,
		"plan":           "Self-Hosted / Open-Source (Community Edition)",
		"billingEnabled": false,
		"commercialSide": true,
		"message":        "Usage tracking and automated billing are managed in Vessel Commercial Edition.",
	})
}

// --- Webhooks endpoints ---

func (h *ProjectSettingsHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	list, err := h.store.ListProjectWebhooks(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *ProjectSettingsHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var webhook types.ProjectWebhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	webhook.ProjectID = projectID

	if err := h.store.CreateProjectWebhook(&webhook); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(webhook)
}

func (h *ProjectSettingsHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")

	if err := h.store.DeleteProjectWebhook(id, projectID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- API Tokens endpoints ---

func (h *ProjectSettingsHandler) ListTokens(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	list, err := h.store.ListProjectTokens(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *ProjectSettingsHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var tok types.ProjectToken
	if err := json.NewDecoder(r.Body).Decode(&tok); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tok.ProjectID = projectID

	fullSecret, err := h.store.CreateProjectToken(&tok)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      tok,
		"secretToken": fullSecret, // Displayed only once
	})
}

func (h *ProjectSettingsHandler) DeleteToken(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")

	if err := h.store.DeleteProjectToken(id, projectID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Members & Workspace Roles endpoints ---

func (h *ProjectSettingsHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	list, err := h.store.ListProjectMembers(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *ProjectSettingsHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var member types.ProjectMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	member.ProjectID = projectID

	if err := h.store.CreateOrInviteProjectMember(&member); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
}

func (h *ProjectSettingsHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")

	if err := h.store.RemoveProjectMember(id, projectID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
