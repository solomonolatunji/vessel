package models

import "time"

type GithubWebhookPayload struct {
	Action      string            `json:"action"`
	Number      int               `json:"number"`
	PullRequest GithubPullRequest `json:"pull_request"`
	Ref         string            `json:"ref"`
	After       string            `json:"after"`
}

type GithubPullRequest struct {
	Head GithubPullRequestHead `json:"head"`
}

type GithubPullRequestHead struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

type Webhook struct {
	ID                    string    `json:"id" db:"id"`
	ServiceID             string    `json:"serviceId" db:"service_id"`
	URL                   string    `json:"url" db:"url"`
	EventTypes            []string  `json:"eventTypes" db:"-"`
	IncludePREnvironments bool      `json:"includePrEnvironments" db:"include_pr_environments"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateWebhookRequest struct {
	URL                   string   `json:"url"`
	EventTypes            []string `json:"eventTypes"`
	IncludePREnvironments bool     `json:"includePrEnvironments"`
}
