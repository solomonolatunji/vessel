package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateEnvironment provisions a new isolated runtime environment inside a project workspace canvas.
func (s *Server) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	var env types.EnvironmentConfig
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	env.ProjectID = projectID
	if env.Name == "" {
		http.Error(w, "Environment name is required (e.g. production, staging)", http.StatusBadRequest)
		return
	}

	if err := s.store.CreateEnvironment(&env); err != nil {
		http.Error(w, "Failed to create environment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(env)
}

// ListEnvironments returns all environments belonging to a project canvas workspace.
func (s *Server) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	envs, err := s.store.ListEnvironments(projectID)
	if err != nil {
		http.Error(w, "Failed to list environments: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(envs)
}

// DeleteEnvironment removes an environment from the project canvas.
func (s *Server) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := s.store.DeleteEnvironment(id); err != nil {
		http.Error(w, "Failed to delete environment: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
