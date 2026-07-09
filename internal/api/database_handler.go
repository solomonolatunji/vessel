package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/solomonolatunji/vessel/internal/types"
)

func (s *Server) handleListDatabases(w http.ResponseWriter, r *http.Request) {
	databases, err := s.store.ListDatabases()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if databases == nil {
		databases = []types.DatabaseConfig{}
	}
	writeJSON(w, http.StatusOK, databases)
}

func (s *Server) handleCreateDatabase(w http.ResponseWriter, r *http.Request) {
	var db types.DatabaseConfig
	if err := json.NewDecoder(r.Body).Decode(&db); err != nil {
		writeError(w, http.StatusBadRequest, "invalid database configuration payload")
		return
	}

	if db.Name == "" || db.Engine == "" {
		writeError(w, http.StatusBadRequest, "name and engine fields are required")
		return
	}

	if db.Port <= 0 {
		switch strings.ToLower(db.Engine) {
		case "postgres", "postgresql":
			db.Port = 5432
		case "mysql":
			db.Port = 3306
		case "redis":
			db.Port = 6379
		case "mongodb", "mongo":
			db.Port = 27017
		default:
			db.Port = 5432
		}
	}
	if db.Username == "" && strings.ToLower(db.Engine) != "redis" {
		db.Username = "vesseladmin"
	}
	if db.DatabaseName == "" {
		db.DatabaseName = "appdb"
	}

	if err := s.store.CreateDatabase(&db); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if s.dbDeployer != nil {
		_, _ = s.dbDeployer.SpinUp(r.Context(), &db)
	}

	writeJSON(w, http.StatusCreated, db)
}

func (s *Server) handleGetDatabase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}

	db, err := s.store.GetDatabase(id)
	if err != nil || db == nil {
		writeError(w, http.StatusNotFound, "database not found")
		return
	}
	writeJSON(w, http.StatusOK, db)
}

func (s *Server) handleDeleteDatabase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}

	if s.dbDeployer != nil {
		_ = s.dbDeployer.Stop(r.Context(), id)
	}

	if err := s.store.DeleteDatabase(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleStartDatabase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}

	db, err := s.store.GetDatabase(id)
	if err != nil || db == nil {
		writeError(w, http.StatusNotFound, "database not found")
		return
	}

	if s.dbDeployer == nil {
		writeError(w, http.StatusServiceUnavailable, "database deployer unavailable")
		return
	}

	if _, err := s.dbDeployer.SpinUp(r.Context(), db); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, db)
}

func (s *Server) handleStopDatabase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}

	if s.dbDeployer != nil {
		_ = s.dbDeployer.Stop(r.Context(), id)
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}
