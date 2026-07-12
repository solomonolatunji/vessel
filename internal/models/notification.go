package models

import (
	"encoding/json"
	"time"
)

type WorkspaceNotificationChannel struct {
	ID          string          `json:"id"`
	WorkspaceID string          `json:"workspaceId"`
	Provider    string          `json:"provider"` // e.g., "discord", "slack", "smtp"
	Config      json.RawMessage `json:"config"`   // Generic JSON config tailored to the provider
	Events      json.RawMessage `json:"events"`   // Array of strings e.g. ["deploy.success", "deploy.failure"]
	IsEnabled   bool            `json:"isEnabled"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type NotificationEvent struct {
	Title       string `json:"title"`
	Message     string `json:"message"`
	Level       string `json:"level"`
	EventType   string `json:"eventType"`
	WorkspaceID string `json:"workspaceId"`
	ProjectID   string `json:"projectId,omitempty"`
	URL         string `json:"url,omitempty"`
}
