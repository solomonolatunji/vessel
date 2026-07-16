package models

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
