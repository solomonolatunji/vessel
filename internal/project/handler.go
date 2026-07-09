package project

import (
	"context"
	"encoding/json"
	"net/http"
)

// ProxyReloader is the minimal proxy surface used by project handlers.
type ProxyReloader interface {
	Reload(ctx context.Context) error
}

// Handler serves HTTP requests for the project domain.
type Handler struct {
	service     *Service
	proxy       ProxyReloader
	extractUser func(r *http.Request) string
}

// NewHandler creates a new project Handler.
func NewHandler(service *Service, proxy ProxyReloader, extractUser func(r *http.Request) string) *Handler {
	return &Handler{
		service:     service,
		proxy:       proxy,
		extractUser: extractUser,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// List handles GET /api/projects.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

// Create handles POST /api/projects.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid project configuration payload")
		return
	}

	p, err := h.service.Create(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.proxy != nil {
		_ = h.proxy.Reload(r.Context())
	}
	writeJSON(w, http.StatusCreated, p)
}

// Get handles GET /api/projects/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	p, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// Delete handles DELETE /api/projects/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.proxy != nil {
		_ = h.proxy.Reload(r.Context())
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
