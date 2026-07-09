package storage

import (
	"context"
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
)

type Deployer interface {
	SpinUp(ctx context.Context, s *models.Storage) (string, error)
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
	storages, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if storages == nil {
		storages = []*Storage{}
	}
	writeJSON(w, http.StatusOK, storages)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var s Storage
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid storage configuration payload")
		return
	}

	if s.Name == "" {
		writeError(w, http.StatusBadRequest, "storage name is required")
		return
	}
	if s.APIPort <= 0 {
		s.APIPort = 9000
	}
	if s.ConsolePort <= 0 {
		s.ConsolePort = 9001
	}
	if s.AccessKey == "" {
		s.AccessKey = "vesseladmin"
	}
	if s.SecretKey == "" {
		s.SecretKey = "vesselsecretkey123"
	}
	if s.BucketName == "" {
		s.BucketName = "vessel-backups"
	}
	if s.Type == "" {
		s.Type = "minio"
	}

	if err := h.repo.Create(r.Context(), &s); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.deployer != nil {
		_, _ = h.deployer.SpinUp(r.Context(), toModelStorage(&s))
	}

	writeJSON(w, http.StatusCreated, s)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}

	s, err := h.repo.GetByID(r.Context(), id)
	if err != nil || s == nil {
		writeError(w, http.StatusNotFound, "storage record not found")
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing storage id parameter")
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
		writeError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}

	s, err := h.repo.GetByID(r.Context(), id)
	if err != nil || s == nil {
		writeError(w, http.StatusNotFound, "storage record not found")
		return
	}

	if h.deployer == nil {
		writeError(w, http.StatusServiceUnavailable, "storage deployer unavailable")
		return
	}

	if _, err := h.deployer.SpinUp(r.Context(), toModelStorage(s)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing storage id parameter")
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
