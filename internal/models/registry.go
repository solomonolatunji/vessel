package models

import "time"

type Registry struct {
	ID            string    `json:"id" db:"id"`
	ProjectID     string    `json:"projectId" db:"project_id"`
	Name          string    `json:"name" db:"name"`
	RegistryURL   string    `json:"registryUrl" db:"registry_url"`
	Username      string    `json:"username" db:"username"`
	PasswordToken string    `json:"passwordToken,omitempty" db:"password_token"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}
