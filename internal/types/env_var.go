package types

import "time"

// EnvVar represents an encrypted environment variable record stored in SQLite.
type EnvVar struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"projectId"`
	Key            string    `json:"key"`
	EncryptedValue string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
