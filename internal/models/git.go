package models

import "time"

type GitProviderConfig struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Provider    string    `json:"provider"`
	AccessToken string    `json:"accessToken,omitempty"`
	AccountName string    `json:"accountName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
	ID             string    `json:"id"`
	WorkspaceID    string    `json:"workspaceId"`
	Name           string    `json:"name"`
	AppID          string    `json:"appId"`
	InstallationID string    `json:"installationId"`
	ClientID       string    `json:"clientId"`
	ClientSecret   string    `json:"clientSecret"`
	WebhookSecret  string    `json:"webhookSecret"`
	PrivateKey     string    `json:"privateKey"`
	IsPublic       bool      `json:"isPublic"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type GitlabApp struct {
	ID            string    `json:"id"`
	WorkspaceID   string    `json:"workspaceId"`
	Name          string    `json:"name"`
	AppID         string    `json:"appId"`
	AppSecret     string    `json:"appSecret"`
	WebhookSecret string    `json:"webhookSecret"`
	APIURL        string    `json:"apiUrl"`
	IsPublic      bool      `json:"isPublic"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type BitbucketApp struct {
	ID            string    `json:"id"`
	WorkspaceID   string    `json:"workspaceId"`
	Name          string    `json:"name"`
	Workspace     string    `json:"workspace"`
	ClientID      string    `json:"clientId"`
	ClientSecret  string    `json:"clientSecret"`
	WebhookSecret string    `json:"webhookSecret"`
	IsPublic      bool      `json:"isPublic"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
