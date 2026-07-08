package types

import "time"

type StorageConfig struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	APIPort       int       `json:"apiPort"`
	ConsolePort   int       `json:"consolePort"`
	AccessKey     string    `json:"accessKey"`
	SecretKey     string    `json:"secretKey"`
	BucketName    string    `json:"bucketName"`
	VolumePath    string    `json:"volumePath"`
	ContainerID   string    `json:"containerId"`
	Status        string    `json:"status"`
	InternalDNS   string    `json:"internalDns"`
	ExternalDNS   string    `json:"externalDns"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
