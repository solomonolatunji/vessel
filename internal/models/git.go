package models

import "time"

type GitProviderConfig struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"userId" db:"user_id"`
	Provider    string    `json:"provider" db:"provider"`
	AccessToken string    `json:"accessToken,omitempty" db:"encrypted_access_token"`
	AccountName string    `json:"accountName" db:"account_name"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type GitRepository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"fullName"`
	Private       bool   `json:"private"`
	CloneURL      string `json:"cloneUrl"`
	HTMLURL       string `json:"htmlUrl"`
	DefaultBranch string `json:"defaultBranch"`
}

type GitConnectRequest struct {
	Provider    string `json:"provider"`
	AccessToken string `json:"accessToken"`
	AccountName string `json:"accountName"`
}

type GithubApp struct {
	ID             string    `json:"id" db:"id"`
	WorkspaceID    string    `json:"workspaceId" db:"workspace_id"`
	Name           string    `json:"name" db:"name"`
	AppID          string    `json:"appId" db:"app_id"`
	InstallationID string    `json:"installationId" db:"installation_id"`
	ClientID       string    `json:"clientId" db:"client_id"`
	ClientSecret   string    `json:"clientSecret" db:"client_secret"`
	WebhookSecret  string    `json:"webhookSecret" db:"webhook_secret"`
	PrivateKey     string    `json:"privateKey" db:"private_key"`
	IsPublic       bool      `json:"isPublic" db:"is_public"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type GitlabApp struct {
	ID            string    `json:"id" db:"id"`
	WorkspaceID   string    `json:"workspaceId" db:"workspace_id"`
	Name          string    `json:"name" db:"name"`
	AppID         string    `json:"appId" db:"app_id"`
	AppSecret     string    `json:"appSecret" db:"app_secret"`
	WebhookSecret string    `json:"webhookSecret" db:"webhook_secret"`
	APIURL        string    `json:"apiUrl" db:"api_url"`
	IsPublic      bool      `json:"isPublic" db:"is_public"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type BitbucketApp struct {
	ID            string    `json:"id" db:"id"`
	WorkspaceID   string    `json:"workspaceId" db:"workspace_id"`
	Name          string    `json:"name" db:"name"`
	Workspace     string    `json:"workspace" db:"workspace"`
	ClientID      string    `json:"clientId" db:"client_id"`
	ClientSecret  string    `json:"clientSecret" db:"client_secret"`
	WebhookSecret string    `json:"webhookSecret" db:"webhook_secret"`
	IsPublic      bool      `json:"isPublic" db:"is_public"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}
