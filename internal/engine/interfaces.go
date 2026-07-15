package engine

import (
	"time"

	"vessl.dev/vessl/internal/models"
)

type DeployerStore interface {
	ContainerManagerStore
	ListAppServicesByProject(projectID string) ([]*models.AppService, error)
	GetEnvVars(projectID string) (map[string]string, error)
	ListServiceVariables(serviceID string) ([]*models.Variable, error)
	GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error)
}

type DatabaseDeployerStore interface {
	GetServerSettings() (*models.ServerSettings, error)
	UpdateDatabaseStatus(id string, status models.DatabaseStatus, containerID string) error
	GetDatabase(id string) (*models.Database, error)
}

type StorageDeployerStore interface {
	GetServerSettings() (*models.ServerSettings, error)
	UpdateStorageStatus(id string, status models.StorageStatus, containerID string) error
	GetStorage(id string) (*models.Storage, error)
}

type CronManagerStore interface {
	ListJobs() ([]models.Job, error)
	GetJob(id string) (*models.Job, error)
	GetProject(id string) (*models.ProjectConfig, error)
	UpdateJobStatusAndOutput(id string, status models.JobStatus, lastRunAt *time.Time, output string) error
}

type BackupManagerStore interface {
	ListAllActiveBackupConfigs() ([]*models.BackupConfig, error)
	GetBackupConfig(id string) (*models.BackupConfig, error)
	CreateBackupRecord(rec *models.BackupRecord) error
	GetDatabase(id string) (*models.Database, error)
	UpdateBackupRecord(id string, status models.BackupRecordStatus, filePath, s3URL, logs string, fileSizeBytes int64, completedAt string) error
	GetS3Destination(id string) (*models.S3Destination, error)
	GetBackupRecord(id string) (*models.BackupRecord, error)
	ListBackupRecords(backupConfigID string) ([]*models.BackupRecord, error)
}

type ContainerManagerStore interface {
	GetServerSettings() (*models.ServerSettings, error)
}
