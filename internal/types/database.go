package types

import "time"

// DatabaseConfig represents a managed stateful database engine instance provisioned by Vessel.
type DatabaseConfig struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"` // e.g. production/staging environment ID
	Name          string    `json:"name"`
	Engine        string    `json:"engine"` // postgres, mysql, redis, mongodb
	Version       string    `json:"version"`
	Port          int       `json:"port"`
	Username      string    `json:"username"`
	Password      string    `json:"password"` // returned plaintext or masked depending on API context
	DatabaseName  string    `json:"databaseName"`
	VolumePath    string    `json:"volumePath"`
	ContainerID   string    `json:"containerId"`
	Status        string    `json:"status"` // running, stopped, failed
	InternalDNS   string    `json:"internalDns"`
	ExternalDNS   string    `json:"externalDns"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
