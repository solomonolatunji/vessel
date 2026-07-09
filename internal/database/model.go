package database

import (
	"time"

	"vessel.dev/vessel/internal/models"
)

type Database struct {
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

func toModelDatabase(d *Database) *models.Database {
	if d == nil {
		return nil
	}
	return &models.Database{
		ID:            d.ID,
		ProjectID:     d.ProjectID,
		EnvironmentID: d.EnvironmentID,
		Name:          d.Name,
		Engine:        d.Engine,
		Version:       d.Version,
		Port:          d.Port,
		Username:      d.Username,
		Password:      d.Password,
		DatabaseName:  d.DatabaseName,
		VolumePath:    d.VolumePath,
		ContainerID:   d.ContainerID,
		Status:        d.Status,
		InternalDNS:   d.InternalDNS,
		ExternalDNS:   d.ExternalDNS,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
