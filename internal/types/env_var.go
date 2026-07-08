package types

import "time"

type EnvVar struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"projectId"`
	Key            string    `json:"key"`
	EncryptedValue string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
