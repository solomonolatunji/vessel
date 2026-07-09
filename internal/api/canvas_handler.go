package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ListProjectCanvasSummaries(w http.ResponseWriter, r *http.Request) {
	summaries, err := s.store.ListProjectCanvasSummaries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

func (s *Server) GetProjectCanvasSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	summary, err := s.store.GetProjectCanvasSummary(id)
	if err != nil {
		http.Error(w, "Project not found or summary calculation failed: "+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (s *Server) GetEnvironmentCanvas(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	canvas, err := s.store.GetEnvironmentCanvas(id)
	if err != nil {
		http.Error(w, "Environment canvas not found: "+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(canvas)
}
