package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(ns *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: ns}
}

func (h *NotificationHandler) GetIntegrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	integ, err := h.notificationService.GetIntegration(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, integ)
}

func (h *NotificationHandler) SaveIntegrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var integ models.NotificationIntegration
	if err := json.NewDecoder(r.Body).Decode(&integ); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.notificationService.SaveIntegration(r.Context(), &integ); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, integ)
}

func (h *NotificationHandler) TestNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Channel   string `json:"channel"`
		ProjectID string `json:"projectId,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.notificationService.SendTest(req.Channel, req.ProjectID); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Test notification sent successfully over " + req.Channel,
	})
}

func (h *NotificationHandler) GetProjectPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	projectID := r.PathValue("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "Missing project id parameter")
		return
	}

	pref, err := h.notificationService.GetProjectPref(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, pref)
}

func (h *NotificationHandler) SaveProjectPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	projectID := r.PathValue("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "Missing project id parameter")
		return
	}

	var pref models.ProjectNotificationPref
	if err := json.NewDecoder(r.Body).Decode(&pref); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	pref.ProjectID = projectID

	if err := h.notificationService.SaveProjectPref(r.Context(), &pref); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, pref)
}
