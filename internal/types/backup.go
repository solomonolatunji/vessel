package types

// BackupConfig defines an automated database or storage volume backup schedule.
type BackupConfig struct {
	ID              string `json:"id"`
	ProjectID       string `json:"projectId"`
	DatabaseID      string `json:"databaseId,omitempty"` // Target database container ID (if backing up DB)
	StorageID       string `json:"storageId,omitempty"`  // Target storage container ID (if backing up S3/volume)
	S3DestinationID string `json:"s3DestinationId,omitempty"` // Optional external S3 offsite upload target
	Name            string `json:"name"`
	Schedule        string `json:"schedule"` // Cron expression, e.g. "0 2 * * *" for daily at 2am
	RetentionDays   int    `json:"retentionDays"`
	Status          string `json:"status"` // active, paused
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// BackupRecord represents an execution history entry for a completed or failed backup run.
type BackupRecord struct {
	ID             string `json:"id"`
	BackupConfigID string `json:"backupConfigId"`
	ProjectID      string `json:"projectId"`
	DatabaseID     string `json:"databaseId,omitempty"`
	Status         string `json:"status"` // running, completed, failed
	FilePath       string `json:"filePath"`
	FileSizeBytes  int64  `json:"fileSizeBytes"`
	S3URL          string `json:"s3Url,omitempty"`
	Logs           string `json:"logs"`
	StartedAt      string `json:"startedAt"`
	CompletedAt    string `json:"completedAt"`
}

// S3Destination represents an external or internal S3/MinIO bucket credentials for offsite backups.
type S3Destination struct {
	ID              string `json:"id"`
	ProjectID       string `json:"projectId"`
	Name            string `json:"name"`
	Endpoint        string `json:"endpoint"` // e.g., s3.amazonaws.com or minio.internal:9000
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	CreatedAt       string `json:"createdAt"`
}
