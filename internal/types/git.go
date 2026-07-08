package types

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
