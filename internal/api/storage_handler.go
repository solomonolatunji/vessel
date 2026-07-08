package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/types"
)

// handleListStorage retrieves all managed MinIO and S3 storage instances from the database.
func (s *Server) handleListStorage(w http.ResponseWriter, r *http.Request) {
	storages, err := s.store.ListStorage()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if storages == nil {
		storages = []types.StorageConfig{}
	}
	writeJSON(w, http.StatusOK, storages)
}

// handleCreateStorage registers and spins up a new MinIO object storage container instance.
func (s *Server) handleCreateStorage(w http.ResponseWriter, r *http.Request) {
	var sc types.StorageConfig
	if err := json.NewDecoder(r.Body).Decode(&sc); err != nil {
		writeError(w, http.StatusBadRequest, "invalid storage configuration payload")
		return
	}

	if sc.Name == "" {
		writeError(w, http.StatusBadRequest, "storage name is required")
		return
	}

	if sc.APIPort <= 0 {
		sc.APIPort = 9000
	}
	if sc.ConsolePort <= 0 {
		sc.ConsolePort = 9001
	}
	if sc.AccessKey == "" {
		sc.AccessKey = "vesseladmin"
	}
	if sc.SecretKey == "" {
		sc.SecretKey = "vesselsecretkey123"
	}
	if sc.BucketName == "" {
		sc.BucketName = "vessel-backups"
	}
	if sc.Type == "" {
		sc.Type = "minio"
	}

	if err := s.store.CreateStorage(&sc); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if s.storageDeployer != nil {
		_, _ = s.storageDeployer.SpinUp(r.Context(), &sc)
	}

	writeJSON(w, http.StatusCreated, sc)
}

// handleGetStorage retrieves details and decrypted credentials for an object storage instance.
func (s *Server) handleGetStorage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}

	sc, err := s.store.GetStorage(id)
	if err != nil || sc == nil {
		writeError(w, http.StatusNotFound, "storage record not found")
		return
	}
	writeJSON(w, http.StatusOK, sc)
}

// handleDeleteStorage stops the running object storage container and removes the database record.
func (s *Server) handleDeleteStorage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}

	if s.storageDeployer != nil {
		_ = s.storageDeployer.Stop(r.Context(), id)
	}

	if err := s.store.DeleteStorage(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// handleStartStorage launches or resumes a stopped object storage container.
func (s *Server) handleStartStorage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}

	sc, err := s.store.GetStorage(id)
	if err != nil || sc == nil {
		writeError(w, http.StatusNotFound, "storage record not found")
		return
	}

	if s.storageDeployer == nil {
		writeError(w, http.StatusServiceUnavailable, "storage deployer unavailable")
		return
	}

	if _, err := s.storageDeployer.SpinUp(r.Context(), sc); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sc)
}

// handleStopStorage terminates a running object storage container.
func (s *Server) handleStopStorage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing storage id parameter")
		return
	}

	if s.storageDeployer != nil {
		_ = s.storageDeployer.Stop(r.Context(), id)
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}
