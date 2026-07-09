package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProjectStore interface {
	GetByID(ctx context.Context, id string) (*ProjectConfig, error)
}

type ProjectConfig struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TeamID      string `json:"teamId"`
}

type ProjectDeployer interface {
	CloneOrPullRepository(ctx context.Context, projectID, sourceDir string) error
	DeployProject(ctx context.Context, project *ProjectConfig, sourceDir string) (string, error)
	ReloadProxy(ctx context.Context) error
}

type Handler struct {
	repo          Repository
	serviceRepo   ServiceRepository
	projectStore  ProjectStore
	projectDeploy ProjectDeployer
}

func NewHandler(repo Repository, serviceRepo ServiceRepository, projectStore ProjectStore, projectDeploy ProjectDeployer) *Handler {
	return &Handler{repo: repo, serviceRepo: serviceRepo, projectStore: projectStore, projectDeploy: projectDeploy}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) ListServiceDeployments(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	deps, err := h.repo.ListByService(r.Context(), serviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, deps)
}

func (h *Handler) Trigger(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")

	svc, err := h.serviceRepo.GetByID(r.Context(), serviceID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	svcMap, ok := svc.(map[string]any)
	if !ok {
		writeError(w, http.StatusInternalServerError, "unexpected service type")
		return
	}

	dep := &Deployment{
		ServiceID:     serviceID,
		EnvironmentID: stringOrDefault(svcMap, "environmentId"),
		ProjectID:     stringOrDefault(svcMap, "projectId"),
		Status:        "BUILDING",
		Branch:        stringOrDefault(svcMap, "branch"),
		Trigger:       "Manual Deploy",
		BuildLogs:     "Initiating build...\n",
	}

	if err := h.repo.Create(r.Context(), dep); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, http.StatusAccepted, dep)
}

func (h *Handler) Rollback(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	targetDep, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	newDep := &Deployment{
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        "BUILDING",
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Rollback to " + targetDep.ID,
		Branch:        targetDep.Branch,
		Trigger:       "Rollback",
		BuildLogs:     "Rolling back to deployment " + targetDep.ID + "...\n",
	}

	if err := h.repo.Create(r.Context(), newDep); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, http.StatusAccepted, newDep)
}

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	dep, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":        dep.ID,
		"buildLogs": dep.BuildLogs,
		"status":    dep.Status,
	})
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	metrics := []ServiceMetric{
		{Timestamp: now.Add(-4 * time.Minute).Format(time.RFC3339), CPUPercent: 1.2, MemoryMB: 64.5, NetworkRx: 12.4, NetworkTx: 8.1},
		{Timestamp: now.Add(-3 * time.Minute).Format(time.RFC3339), CPUPercent: 2.1, MemoryMB: 66.0, NetworkRx: 15.0, NetworkTx: 10.2},
		{Timestamp: now.Add(-2 * time.Minute).Format(time.RFC3339), CPUPercent: 1.8, MemoryMB: 65.2, NetworkRx: 14.1, NetworkTx: 9.4},
		{Timestamp: now.Add(-1 * time.Minute).Format(time.RFC3339), CPUPercent: 3.4, MemoryMB: 68.1, NetworkRx: 45.2, NetworkTx: 22.0},
		{Timestamp: now.Format(time.RFC3339), CPUPercent: 1.5, MemoryMB: 66.8, NetworkRx: 18.0, NetworkTx: 11.5},
	}
	writeJSON(w, http.StatusOK, metrics)
}

func (h *Handler) DeployProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	project, err := h.projectStore.GetByID(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	sourceDir := fmt.Sprintf("data/builds/%s", id)
	if h.projectDeploy != nil {
		if err := h.projectDeploy.CloneOrPullRepository(r.Context(), id, sourceDir); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("git operation failed: %v", err))
			return
		}
		containerID, err := h.projectDeploy.DeployProject(r.Context(), project, sourceDir)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("deployment rollout failed: %v", err))
			return
		}
		_ = h.projectDeploy.ReloadProxy(r.Context())
		writeJSON(w, http.StatusOK, map[string]string{
			"status":       "deployed",
			"container_id": containerID,
		})
		return
	}

	writeError(w, http.StatusServiceUnavailable, "deployer not available")
}

func stringOrDefault(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
