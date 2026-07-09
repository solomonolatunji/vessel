package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleGetEnvVars(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	envVars, err := s.store.GetEnvVars(projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if envVars == nil {
		envVars = make(map[string]string)
	}
	writeJSON(w, http.StatusOK, envVars)
}

func (s *Server) handleSetEnvVars(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid environment variable dictionary payload")
		return
	}

	for key, value := range payload {
		if err := s.store.SetEnvVar(projectID, key, value); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
