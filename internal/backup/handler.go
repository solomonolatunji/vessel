package backup

import (
	"context"
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
)

type BackupManager interface {
	RegisterBackup(cfg *models.BackupConfig) error
	UnregisterBackup(backupConfigID string)
	TriggerBackup(ctx context.Context, backupConfigID string) (*models.BackupRecord, error)
}

type Handler struct {
	repo          Repository
	backupManager BackupManager
}

func NewHandler(repo Repository, backupManager BackupManager) *Handler {
	return &Handler{repo: repo, backupManager: backupManager}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.PathValue("projectId")
	}

	var (
		list []*BackupConfig
		err  error
	)
	if projectID != "" {
		list, err = h.repo.ListConfigsByProject(r.Context(), projectID)
	} else {
		list, err = h.repo.ListAllActiveConfigs(r.Context())
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var cfg BackupConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.repo.CreateConfig(r.Context(), &cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.backupManager != nil {
		_ = h.backupManager.RegisterBackup(toModelBackupConfig(&cfg))
	}

	writeJSON(w, http.StatusCreated, cfg)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cfg, err := h.repo.GetConfigByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if cfg == nil {
		writeError(w, http.StatusNotFound, "backup schedule not found")
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cfg, err := h.repo.GetConfigByID(r.Context(), id)
	if err != nil || cfg == nil {
		writeError(w, http.StatusNotFound, "backup schedule not found")
		return
	}

	if h.backupManager != nil {
		h.backupManager.UnregisterBackup(id)
	}

	if err := h.repo.DeleteConfig(r.Context(), id, cfg.ProjectID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Trigger(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if h.backupManager == nil {
		writeError(w, http.StatusServiceUnavailable, "backup manager not initialized")
		return
	}

	rec, err := h.backupManager.TriggerBackup(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, rec)
}

func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	list, err := h.repo.ListRecordsByConfig(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) ListS3Destinations(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.PathValue("projectId")
	}
	list, err := h.repo.ListS3Destinations(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) CreateS3Destination(w http.ResponseWriter, r *http.Request) {
	var dest S3Destination
	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.CreateS3Destination(r.Context(), &dest); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, dest)
}

func (h *Handler) DeleteS3Destination(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		dest, _ := h.repo.GetS3Destination(r.Context(), id)
		if dest != nil {
			projectID = dest.ProjectID
		}
	}
	if err := h.repo.DeleteS3Destination(r.Context(), id, projectID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
