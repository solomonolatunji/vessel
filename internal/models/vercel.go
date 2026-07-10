package models

import "time"

type UserVercelAccount struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	AccessToken string    `json:"accessToken,omitempty"`
	TeamID      *string   `json:"teamId,omitempty"`
	AccountName string    `json:"accountName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type VercelProject struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Framework   string `json:"framework"`
	NodeVersion string `json:"nodeVersion"`
	AccountID   string `json:"accountId"`
}

type VercelEnvVar struct {
	Type   string   `json:"type"` // "system" or "secret" or "plain"
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	Target []string `json:"target"` // "production", "preview", "development"
}
