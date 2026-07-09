package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

type DeploymentHandler struct {
	store *store.Store
}

func NewDeploymentHandler(store *store.Store) *DeploymentHandler {
	return &DeploymentHandler{store: store}
}

func (h *DeploymentHandler) ListServiceDeployments(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")

	deps, err := h.store.ListDeploymentsByService(serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deps)
}

func (h *DeploymentHandler) TriggerServiceDeployment(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")

	svc, err := h.store.GetAppService(serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	dep := &types.DeploymentRecord{
		ServiceID:     svc.ID,
		EnvironmentID: svc.EnvironmentID,
		ProjectID:     svc.ProjectID,
		Status:        "BUILDING",
		Branch:        svc.Branch,
		Trigger:       "Manual Deploy",
		BuildLogs:     "Initiating build from " + svc.RepositoryURL + " branch " + svc.Branch + "...\n",
	}

	if err := h.store.CreateDeployment(dep); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(dep)
}

func (h *DeploymentHandler) RollbackDeployment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	targetDep, err := h.store.GetDeployment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	newDep := &types.DeploymentRecord{
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        "BUILDING",
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Rollback to " + targetDep.ID,
		Branch:        targetDep.Branch,
		Trigger:       "Rollback",
		BuildLogs:     "Rolling back container instance to deployment " + targetDep.ID + "...\n",
	}

	if err := h.store.CreateDeployment(newDep); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(newDep)
}

func (h *DeploymentHandler) GetDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	dep, err := h.store.GetDeployment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":        dep.ID,
		"buildLogs": dep.BuildLogs,
		"status":    dep.Status,
	})
}

func (h *DeploymentHandler) GetServiceMetrics(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	metrics := []types.ServiceMetric{
		{Timestamp: now.Add(-4 * time.Minute).Format(time.RFC3339), CPUPercent: 1.2, MemoryMB: 64.5, NetworkRx: 12.4, NetworkTx: 8.1},
		{Timestamp: now.Add(-3 * time.Minute).Format(time.RFC3339), CPUPercent: 2.1, MemoryMB: 66.0, NetworkRx: 15.0, NetworkTx: 10.2},
		{Timestamp: now.Add(-2 * time.Minute).Format(time.RFC3339), CPUPercent: 1.8, MemoryMB: 65.2, NetworkRx: 14.1, NetworkTx: 9.4},
		{Timestamp: now.Add(-1 * time.Minute).Format(time.RFC3339), CPUPercent: 3.4, MemoryMB: 68.1, NetworkRx: 45.2, NetworkTx: 22.0},
		{Timestamp: now.Format(time.RFC3339), CPUPercent: 1.5, MemoryMB: 66.8, NetworkRx: 18.0, NetworkTx: 11.5},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
