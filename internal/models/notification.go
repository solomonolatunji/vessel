package models

import (
	"encoding/json"
	"time"
)

type WorkspaceNotificationChannel struct {
	ID          string          `json:"id" db:"id"`
	WorkspaceID string          `json:"workspaceId" db:"workspace_id"`
	Provider    string          `json:"provider" db:"provider"` // e.g., "discord", "slack", "smtp"
	Config      json.RawMessage `json:"config" db:"-"`          // Generic JSON config tailored to the provider
	Events      json.RawMessage `json:"events" db:"-"`          // Array of strings e.g. ["deploy.success", "deploy.failure"]
	IsEnabled   bool            `json:"isEnabled" db:"is_enabled"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" db:"updated_at"`
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
