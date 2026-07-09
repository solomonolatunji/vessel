package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

type ServiceVarHandler struct {
	store *store.Store
}

func NewServiceVarHandler(store *store.Store) *ServiceVarHandler {
	return &ServiceVarHandler{store: store}
}

// ListServiceVariables retrieves all variables for a service (`Variables` tab).
func (h *ServiceVarHandler) ListServiceVariables(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")

	list, err := h.store.ListServiceVariables(serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// SetServiceVariable creates or updates a variable (`Variables` tab -> `Add Variable`).
func (h *ServiceVarHandler) SetServiceVariable(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")

	var v types.ServiceVariable
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v.ServiceID = serviceID

	if v.EnvironmentID == "" {
		svc, err := h.store.GetAppService(serviceID)
		if err == nil && svc != nil {
			v.EnvironmentID = svc.EnvironmentID
		}
	}

	if err := h.store.SetServiceVariable(&v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

// BulkSetServiceVariables handles the `RAW Editor` tab where multiple KEY=VALUE lines are saved at once.
func (h *ServiceVarHandler) BulkSetServiceVariables(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")

	var req struct {
		Variables     []*types.ServiceVariable `json:"variables"`
		EnvironmentID string                   `json:"environmentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.EnvironmentID == "" {
		svc, err := h.store.GetAppService(serviceID)
		if err == nil && svc != nil {
			req.EnvironmentID = svc.EnvironmentID
		}
	}

	if err := h.store.BulkSetServiceVariables(serviceID, req.EnvironmentID, req.Variables); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req.Variables)
}

// DeleteServiceVariable removes a service variable.
func (h *ServiceVarHandler) DeleteServiceVariable(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	id := r.PathValue("id")

	if err := h.store.DeleteServiceVariable(id, serviceID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
