package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type DeploymentListeners struct {
	dispatcher *DispatcherService
	appRepo    repositories.AppServiceRepository
}

func NewDeploymentListeners(dispatcher *DispatcherService, appRepo repositories.AppServiceRepository) *DeploymentListeners {
	return &DeploymentListeners{dispatcher: dispatcher, appRepo: appRepo}
}

func (l *DeploymentListeners) SendNotification(e DeploymentCompleted) {
	commit := e.CommitHash
	if len(commit) > 7 {
		commit = commit[:7]
	}
	msg := fmt.Sprintf("Deploy %s: %s (%s)", e.Status, e.ServiceID, commit)
	notifEvent := &models.NotificationEvent{
		ProjectID: e.ProjectID,
		Level:     e.Status,
		Title:     "Deployment " + e.Status,
		Message:   msg,
		URL:       fmt.Sprintf("https://codedock.local/projects/%s/services/%s", e.ProjectID, e.ServiceID),
	}
	l.dispatcher.Dispatch(notifEvent)
}

func (l *DeploymentListeners) UpdateAuditLog(e DeploymentCompleted) {
	slog.Info("deployment completed", "serviceID", e.ServiceID, "status", e.Status)
}

func (l *DeploymentListeners) TriggerWebhook(e DeploymentCompleted) {
	slog.Info("triggering webhook", "projectID", e.ProjectID)
	if l.appRepo == nil {
		return
	}
	ctx := context.Background()
	svc, err := l.appRepo.GetByID(ctx, e.ServiceID)
	if err != nil || svc == nil {
		slog.Error("service not found for webhook dispatch", "serviceID", e.ServiceID, "err", err)
		return
	}
	if svc.ProjectID != e.ProjectID {
		slog.Error("service project mismatch for webhook dispatch", "serviceID", e.ServiceID, "eventProjectID", e.ProjectID, "serviceProjectID", svc.ProjectID)
		return
	}

	webhooks, err := l.appRepo.ListWebhooksByService(ctx, e.ServiceID)
	if err != nil {
		slog.Error("failed to list webhooks", "serviceID", e.ServiceID, "err", err)
		return
	}
	if len(webhooks) == 0 {
		return
	}

	app, err := l.appRepo.GetByID(ctx, e.ServiceID)
	if err != nil || app == nil {
		slog.Error("failed to get app service", "serviceID", e.ServiceID, "err", err)
		return
	}

	isPREnvironment := false
	if e.Branch != "" && app.Branch != "" && e.Branch != app.Branch {
		isPREnvironment = true
	}

	eventType := "deployment." + strings.ToLower(e.Status)

	payloadMap := map[string]interface{}{
		"event":      eventType,
		"serviceId":  e.ServiceID,
		"projectId":  e.ProjectID,
		"status":     e.Status,
		"commitHash": e.CommitHash,
		"branch":     e.Branch,
	}
	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, hook := range webhooks {
		shouldSend := false
		for _, et := range hook.EventTypes {
			if et == "*" || et == eventType || et == "deployment.*" {
				shouldSend = true
				break
			}
		}
		if !shouldSend || (isPREnvironment && !hook.IncludePREnvironments) {
			continue
		}

		req, err := http.NewRequestWithContext(ctx, "POST", hook.URL, bytes.NewBuffer(payload))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			go func(r *http.Request) {
				resp, err := client.Do(r)
				if err == nil {
					resp.Body.Close()
				} else {
					slog.Warn("webhook dispatch failed", "url", hook.URL, "err", err)
				}
			}(req)
		}
	}
}

func (l *DeploymentListeners) Register() {
	On("deployment.completed", l.SendNotification)
	On("deployment.completed", l.UpdateAuditLog)
	On("deployment.completed", l.TriggerWebhook)
}
