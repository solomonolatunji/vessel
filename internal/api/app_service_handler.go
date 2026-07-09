package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/types"
)

func (s *Server) CreateAppService(w http.ResponseWriter, r *http.Request) {
	envID := r.PathValue("id")

	env, err := s.store.GetEnvironment(envID)
	if err != nil || env == nil {
		http.Error(w, "Target environment not found: "+err.Error(), http.StatusNotFound)
		return
	}

	var app types.AppServiceConfig
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	app.EnvironmentID = envID
	app.ProjectID = env.ProjectID
	if app.Name == "" {
		http.Error(w, "App service name is required (e.g. recovery, wallet-bot)", http.StatusBadRequest)
		return
	}
	if app.InternalPort == 0 {
		app.InternalPort = 3000
	}
	if app.Status == "" {
		app.Status = "building"
	}

	if err := s.store.CreateAppService(&app); err != nil {
		http.Error(w, "Failed to register app service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

func (s *Server) ListAppServicesByEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := r.PathValue("id")

	apps, err := s.store.ListAppServicesByEnvironment(envID)
	if err != nil {
		http.Error(w, "Failed to list app services: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func (s *Server) ListAppServicesByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		projectID = r.PathValue("projectId")
	}

	apps, err := s.store.ListAppServicesByProject(projectID)
	if err != nil {
		http.Error(w, "Failed to list app services for project: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func (s *Server) GetAppService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	app, err := s.store.GetAppService(id)
	if err != nil {
		http.Error(w, "App service not found: "+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func (s *Server) UpdateAppService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := s.store.GetAppService(id)
	if err != nil || existing == nil {
		http.Error(w, "App service not found: "+err.Error(), http.StatusNotFound)
		return
	}

	var payload types.AppServiceConfig
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	payload.ID = id
	payload.ProjectID = existing.ProjectID
	payload.EnvironmentID = existing.EnvironmentID

	if err := s.store.UpdateAppService(&payload); err != nil {
		http.Error(w, "Failed to update app service settings: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func (s *Server) DeleteAppService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := s.store.DeleteAppService(id); err != nil {
		http.Error(w, "Failed to delete app service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
