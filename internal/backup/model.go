package backup

import "vessel.dev/vessel/internal/models"

type BackupConfig struct {
	ID              string `json:"id"`
	ProjectID       string `json:"projectId"`
	DatabaseID      string `json:"databaseId,omitempty"`
	StorageID       string `json:"storageId,omitempty"`
	S3DestinationID string `json:"s3DestinationId,omitempty"`
	Name            string `json:"name"`
	Schedule        string `json:"schedule"`
	RetentionDays   int    `json:"retentionDays"`
	Status          string `json:"status"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type BackupRecord struct {
	ID             string `json:"id"`
	BackupConfigID string `json:"backupConfigId"`
	ProjectID      string `json:"projectId"`
	DatabaseID     string `json:"databaseId,omitempty"`
	Status         string `json:"status"`
	FilePath       string `json:"filePath"`
	FileSizeBytes  int64  `json:"fileSizeBytes"`
	S3URL          string `json:"s3Url,omitempty"`
	Logs           string `json:"logs"`
	StartedAt      string `json:"startedAt"`
	CompletedAt    string `json:"completedAt"`
}

type S3Destination struct {
	ID              string `json:"id"`
	ProjectID       string `json:"projectId"`
	Name            string `json:"name"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	CreatedAt       string `json:"createdAt"`
}

func toModelBackupConfig(c *BackupConfig) *models.BackupConfig {
	if c == nil {
		return nil
	}
	return &models.BackupConfig{
		ID:              c.ID,
		ProjectID:       c.ProjectID,
		DatabaseID:      c.DatabaseID,
		StorageID:       c.StorageID,
		S3DestinationID: c.S3DestinationID,
		Name:            c.Name,
		Schedule:        c.Schedule,
		RetentionDays:   c.RetentionDays,
		Status:          c.Status,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}
