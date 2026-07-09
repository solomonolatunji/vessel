package types

import "time"

type ProjectConfig struct {
	ID          string    `json:"id"`
	TeamID      string    `json:"teamId,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
