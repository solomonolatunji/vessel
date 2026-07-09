package database

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"vessel.dev/vessel/internal/models"
)

type Deployer interface {
	SpinUp(ctx context.Context, db *models.Database) (string, error)
	Stop(ctx context.Context, id string) error
}

type Handler struct {
	repo     Repository
	deployer Deployer
}

func NewHandler(repo Repository, deployer Deployer) *Handler {
	return &Handler{repo: repo, deployer: deployer}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	databases, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if databases == nil {
		databases = []*Database{}
	}
	writeJSON(w, http.StatusOK, databases)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid database configuration payload")
		return
	}

	if req.Name == "" || req.Engine == "" {
		writeError(w, http.StatusBadRequest, "name and engine fields are required")
		return
	}

	if req.Port <= 0 {
		switch strings.ToLower(req.Engine) {
		case "postgres", "postgresql":
			req.Port = 5432
		case "mysql":
			req.Port = 3306
		case "redis":
			req.Port = 6379
		case "mongodb", "mongo":
			req.Port = 27017
		default:
			req.Port = 5432
		}
	}
	if req.Username == "" && strings.ToLower(req.Engine) != "redis" {
		req.Username = "vesseladmin"
	}
	if req.DatabaseName == "" {
		req.DatabaseName = "appdb"
	}

	db := &Database{
		ProjectID:     req.ProjectID,
		EnvironmentID: req.EnvironmentID,
		Name:          req.Name,
		Engine:        req.Engine,
		Version:       req.Version,
		Port:          req.Port,
		Username:      req.Username,
		Password:      req.Password,
		DatabaseName:  req.DatabaseName,
		VolumePath:    req.VolumePath,
		Status:        "stopped",
	}

	if err := h.repo.Create(r.Context(), db); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.deployer != nil {
		_, _ = h.deployer.SpinUp(r.Context(), toModelDatabase(db))
	}

	writeJSON(w, http.StatusCreated, db)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}

	db, err := h.repo.GetByID(r.Context(), id)
	if err != nil || db == nil {
		writeError(w, http.StatusNotFound, "database not found")
		return
	}
	writeJSON(w, http.StatusOK, db)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}

	if h.deployer != nil {
		_ = h.deployer.Stop(r.Context(), id)
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}

	db, err := h.repo.GetByID(r.Context(), id)
	if err != nil || db == nil {
		writeError(w, http.StatusNotFound, "database not found")
		return
	}

	if h.deployer == nil {
		writeError(w, http.StatusServiceUnavailable, "database deployer unavailable")
		return
	}

	if _, err := h.deployer.SpinUp(r.Context(), toModelDatabase(db)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, db)
}

func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing database id parameter")
		return
	}

	if h.deployer != nil {
		_ = h.deployer.Stop(r.Context(), id)
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
