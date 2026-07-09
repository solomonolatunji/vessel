package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/orchestrator"
	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

type BackupHandler struct {
	store         *store.Store
	backupManager *orchestrator.BackupManager
}

func NewBackupHandler(s *store.Store, bm *orchestrator.BackupManager) *BackupHandler {
	return &BackupHandler{
		store:         s,
		backupManager: bm,
	}
}

func (h *BackupHandler) ListBackups(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.PathValue("projectId")
	}

	if projectID != "" {
		list, err := h.store.ListBackupConfigs(projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
		return
	}

	list, err := h.store.ListAllActiveBackupConfigs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *BackupHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	var cfg types.BackupConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.CreateBackupConfig(&cfg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.backupManager != nil {
		_ = h.backupManager.RegisterBackup(&cfg)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cfg)
}

func (h *BackupHandler) GetBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cfg, err := h.store.GetBackupConfig(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cfg == nil {
		http.Error(w, "backup schedule not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (h *BackupHandler) DeleteBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cfg, err := h.store.GetBackupConfig(id)
	if err != nil || cfg == nil {
		http.Error(w, "backup schedule not found", http.StatusNotFound)
		return
	}

	if h.backupManager != nil {
		h.backupManager.UnregisterBackup(id)
	}

	if err := h.store.DeleteBackupConfig(id, cfg.ProjectID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BackupHandler) TriggerBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if h.backupManager == nil {
		http.Error(w, "backup manager not initialized", http.StatusServiceUnavailable)
		return
	}

	rec, err := h.backupManager.TriggerBackup(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(rec)
}

func (h *BackupHandler) ListBackupRecords(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	list, err := h.store.ListBackupRecords(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *BackupHandler) ListS3Destinations(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.PathValue("projectId")
	}
	list, err := h.store.ListS3Destinations(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *BackupHandler) CreateS3Destination(w http.ResponseWriter, r *http.Request) {
	var dest types.S3Destination
	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.store.CreateS3Destination(&dest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dest)
}

func (h *BackupHandler) DeleteS3Destination(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		dest, _ := h.store.GetS3Destination(id)
		if dest != nil {
			projectID = dest.ProjectID
		}
	}
	if err := h.store.DeleteS3Destination(id, projectID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
