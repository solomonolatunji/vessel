package types

import "time"

type DatabaseConfig struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	Engine        string    `json:"engine"`
	Version       string    `json:"version"`
	Port          int       `json:"port"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	DatabaseName  string    `json:"databaseName"`
	VolumePath    string    `json:"volumePath"`
	ContainerID   string    `json:"containerId"`
	Status        string    `json:"status"`
	InternalDNS   string    `json:"internalDns"`
	ExternalDNS   string    `json:"externalDns"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
