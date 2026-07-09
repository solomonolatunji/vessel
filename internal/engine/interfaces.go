package orchestrator

import (
	"time"

	"vessel.dev/vessel/internal/backup"
	"vessel.dev/vessel/internal/database"
	"vessel.dev/vessel/internal/job"
	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/service"
	"vessel.dev/vessel/internal/service_var"
	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/storage"
)

type DeployerStore interface {
	ContainerManagerStore
	ListAppServicesByProject(projectID string) ([]*service.AppService, error)
	GetEnvVars(projectID string) (map[string]string, error)
	ListServiceVariables(serviceID string) ([]*service_var.Variable, error)
}

type DatabaseDeployerStore interface {
	GetServerSettings() (*settings.ServerSettings, error)
	UpdateDatabaseStatus(id string, status string, containerID string) error
	GetDatabase(id string) (*database.Database, error)
}

type StorageDeployerStore interface {
	GetServerSettings() (*settings.ServerSettings, error)
	UpdateStorageStatus(id string, status string, containerID string) error
	GetStorage(id string) (*storage.Storage, error)
}

type CronManagerStore interface {
	ListJobs() ([]job.Job, error)
	GetJob(id string) (*job.Job, error)
	GetProject(id string) (*project.ProjectConfig, error)
	UpdateJobStatusAndOutput(id string, status string, lastRunAt *time.Time, output string) error
}

type BackupManagerStore interface {
	ListAllActiveBackupConfigs() ([]*backup.BackupConfig, error)
	GetBackupConfig(id string) (*backup.BackupConfig, error)
	CreateBackupRecord(rec *backup.BackupRecord) error
	GetDatabase(id string) (*database.Database, error)
	UpdateBackupRecord(id, status, filePath, s3URL, logs string, fileSizeBytes int64, completedAt string) error
	GetS3Destination(id string) (*backup.S3Destination, error)
	ListBackupRecords(backupConfigID string) ([]*backup.BackupRecord, error)
}

type ContainerManagerStore interface {
	GetServerSettings() (*settings.ServerSettings, error)
}
