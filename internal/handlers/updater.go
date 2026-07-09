package handlers

import (
	"net/http"

	"vessel.dev/vessel/internal/services"
)

type UpdaterHandler struct {
	updaterService *services.UpdaterService
}

func NewUpdaterHandler(s *services.UpdaterService) *UpdaterHandler {
	return &UpdaterHandler{updaterService: s}
}

func (h *UpdaterHandler) GetUpdateStatus(w http.ResponseWriter, r *http.Request) {
	if h.updaterService == nil {
		WriteError(w, http.StatusInternalServerError, "updater service not initialized")
		return
	}
	status := h.updaterService.GetStatus()
	WriteJSON(w, http.StatusOK, status)
}

func (h *UpdaterHandler) CheckUpdate(w http.ResponseWriter, r *http.Request) {
	if h.updaterService == nil {
		WriteError(w, http.StatusInternalServerError, "updater service not initialized")
		return
	}
	if _, err := h.updaterService.CheckForUpdates(r.Context()); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	status := h.updaterService.GetStatus()
	WriteJSON(w, http.StatusOK, status)
}

func (h *UpdaterHandler) DeployUpdate(w http.ResponseWriter, r *http.Request) {
	if h.updaterService == nil {
		WriteError(w, http.StatusInternalServerError, "updater service not initialized")
		return
	}
	if err := h.updaterService.DeployUpdate(r.Context()); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusAccepted, map[string]string{
		"message": "update deployment triggered",
	})
}
