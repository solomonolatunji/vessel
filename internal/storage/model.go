package storage

import (
	"time"

	"vessel.dev/vessel/internal/models"
)

type Storage struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	APIPort       int       `json:"apiPort"`
	ConsolePort   int       `json:"consolePort"`
	AccessKey     string    `json:"accessKey"`
	SecretKey     string    `json:"secretKey,omitempty"`
	BucketName    string    `json:"bucketName"`
	VolumePath    string    `json:"volumePath"`
	ContainerID   string    `json:"containerId"`
	Status        string    `json:"status"`
	InternalDNS   string    `json:"internalDns"`
	ExternalDNS   string    `json:"externalDns"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func toModelStorage(s *Storage) *models.Storage {
	if s == nil {
		return nil
	}
	return &models.Storage{
		ID:            s.ID,
		ProjectID:     s.ProjectID,
		EnvironmentID: s.EnvironmentID,
		Name:          s.Name,
		Type:          s.Type,
		APIPort:       s.APIPort,
		ConsolePort:   s.ConsolePort,
		AccessKey:     s.AccessKey,
		SecretKey:     s.SecretKey,
		BucketName:    s.BucketName,
		VolumePath:    s.VolumePath,
		ContainerID:   s.ContainerID,
		Status:        s.Status,
		InternalDNS:   s.InternalDNS,
		ExternalDNS:   s.ExternalDNS,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
