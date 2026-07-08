package types

import "time"

// StorageConfig represents a managed MinIO or S3-compatible object storage container provisioned by Vessel.
type StorageConfig struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // minio
	APIPort     int       `json:"apiPort"`
	ConsolePort int       `json:"consolePort"`
	AccessKey   string    `json:"accessKey"`
	SecretKey   string    `json:"secretKey"`
	BucketName  string    `json:"bucketName"`
	VolumePath  string    `json:"volumePath"`
	ContainerID string    `json:"containerId"`
	Status      string    `json:"status"` // running, stopped, failed
	InternalDNS string    `json:"internalDns"`
	ExternalDNS string    `json:"externalDns"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
