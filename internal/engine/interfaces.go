package engine

import (
	"time"

	"codedock.dev/codedock/internal/models"
)

type DeployerStore interface {
	ContainerManagerStore
	ListAppServicesByProject(projectID string) ([]*models.AppService, error)
	GetEnvVars(projectID string) (map[string]string, error)
	ListServiceVariables(serviceID string) ([]*models.Variable, error)
	GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error)
	UpdateAppService(app *models.AppService) error
	ListLogDrainsByService(serviceID string) ([]*models.LogDrain, error)
}

type DatabaseDeployerStore interface {
	GetServerSettings() (*models.ServerSettings, error)
	UpdateDatabaseStatus(id string, status models.DatabaseStatus, containerID string) error
	GetDatabase(id string) (*models.Database, error)
}

type CronManagerStore interface {
	ListScheduledTasks() ([]models.ScheduledTask, error)
	GetScheduledTask(id string) (*models.ScheduledTask, error)
	GetProject(id string) (*models.ProjectConfig, error)
	GetAppService(id string) (*models.AppService, error)
	UpdateScheduledTaskStatusAndOutput(id string, status models.ScheduledTaskStatus, lastRunAt *time.Time, output string) error
}

type BackupManagerStore interface {
	ListAllActiveBackupConfigs() ([]*models.BackupConfig, error)
	GetBackupConfig(id string) (*models.BackupConfig, error)
	CreateBackupRecord(rec *models.BackupRecord) error
	GetDatabase(id string) (*models.Database, error)
	UpdateBackupRecord(opts models.UpdateBackupRecordOpts) error
	GetS3Destination(id string) (*models.S3Destination, error)
	GetBackupRecord(id string) (*models.BackupRecord, error)
	ListBackupRecords(backupConfigID string) ([]*models.BackupRecord, error)
}

type ContainerManagerStore interface {
	GetServerSettings() (*models.ServerSettings, error)
}
