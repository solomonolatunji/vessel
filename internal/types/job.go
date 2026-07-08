package types

import "time"

type JobConfig struct {
	ID         string     `json:"id"`
	ProjectID  string     `json:"projectId"`
	Name       string     `json:"name"`
	Schedule   string     `json:"schedule"`
	Command    string     `json:"command"`
	Status     string     `json:"status"`
	LastRunAt  *time.Time `json:"lastRunAt"`
	LastOutput string     `json:"lastOutput"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}
