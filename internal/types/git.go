package types

import "time"

// GitProviderConfig represents a user's authenticated OAuth or Personal Access Token connection to a Git platform.
type GitProviderConfig struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Provider    string    `json:"provider"`
	AccessToken string    `json:"accessToken,omitempty"`
	AccountName string    `json:"accountName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GitRepository represents repository metadata fetched from GitHub or GitLab APIs.
type GitRepository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"fullName"`
	Private       bool   `json:"private"`
	CloneURL      string `json:"cloneUrl"`
	HTMLURL       string `json:"htmlUrl"`
	DefaultBranch string `json:"defaultBranch"`
}

// GitConnectRequest represents the payload sent by the user when linking a Git provider token.
type GitConnectRequest struct {
	Provider    string `json:"provider"`
	AccessToken string `json:"accessToken"`
	AccountName string `json:"accountName"`
}
