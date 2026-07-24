package http

import (
	"context"
	"time"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type engineAdapter struct {
	settingsRepo      repositories.SettingsRepository
	appRepo           repositories.AppServiceRepository
	envRepo           repositories.EnvRepository
	dbRepo            repositories.DatabaseRepository
	projectRepo       repositories.ProjectRepository
	scheduledTaskRepo repositories.ScheduledTaskRepository
	backupRepo        repositories.BackupRepository
	s3Repo            repositories.S3DestinationRepository
	serviceVarRepo    repositories.ServiceVarRepository
	serverlessRepo    repositories.ServerlessRepository
}

func newEngineAdapter(
	settingsRepo repositories.SettingsRepository,
	appRepo repositories.AppServiceRepository,
	envRepo repositories.EnvRepository,
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
	scheduledTaskRepo repositories.ScheduledTaskRepository,
	backupRepo repositories.BackupRepository,
	s3Repo repositories.S3DestinationRepository,
	serviceVarRepo repositories.ServiceVarRepository,
	serverlessRepo repositories.ServerlessRepository,
) *engineAdapter {
	return &engineAdapter{
		settingsRepo:      settingsRepo,
		appRepo:           appRepo,
		envRepo:           envRepo,
		dbRepo:            dbRepo,
		projectRepo:       projectRepo,
		scheduledTaskRepo: scheduledTaskRepo,
		backupRepo:        backupRepo,
		s3Repo:            s3Repo,
		serviceVarRepo:    serviceVarRepo,
		serverlessRepo:    serverlessRepo,
	}
}

func (a *engineAdapter) GetServerSettings() (*models.ServerSettings, error) {
	return a.settingsRepo.GetServerSettings(context.Background())
}

func (a *engineAdapter) ListAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return a.appRepo.ListByProject(context.Background(), projectID)
}

func (a *engineAdapter) GetEnvVars(projectID string) (map[string]string, error) {
	return a.envRepo.GetVars(context.Background(), projectID)
}

func (a *engineAdapter) ListServiceVariables(serviceID string) ([]*models.Variable, error) {
	return a.serviceVarRepo.ListByService(context.Background(), serviceID)
}

func (a *engineAdapter) GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error) {
	return a.serverlessRepo.GetCodeByServiceID(context.Background(), serviceID)
}

func (a *engineAdapter) UpdateDatabaseStatus(id string, status models.DatabaseStatus, containerID string) error {
	db, err := a.dbRepo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	db.Status = status
	db.ContainerID = containerID
	return a.dbRepo.Update(context.Background(), db)
}

func (a *engineAdapter) GetDatabase(id string) (*models.Database, error) {
	return a.dbRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) ListScheduledTasks() ([]models.ScheduledTask, error) {
	return a.scheduledTaskRepo.ListAll(context.Background())
}

func (a *engineAdapter) GetScheduledTask(id string) (*models.ScheduledTask, error) {
	return a.scheduledTaskRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) GetProject(id string) (*models.ProjectConfig, error) {
	return a.projectRepo.Get(context.Background(), id)
}

func (a *engineAdapter) GetAppService(id string) (*models.AppService, error) {
	return a.appRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) UpdateScheduledTaskStatusAndOutput(id string, status models.ScheduledTaskStatus, lastRunAt *time.Time, output string) error {
	return a.scheduledTaskRepo.UpdateStatus(context.Background(), id, status, lastRunAt, output)
}

func (a *engineAdapter) ListAllActiveBackupConfigs() ([]*models.BackupConfig, error) {
	return a.backupRepo.ListAllActiveConfigs(context.Background())
}

func (a *engineAdapter) GetBackupConfig(id string) (*models.BackupConfig, error) {
	return a.backupRepo.GetConfigByID(context.Background(), id)
}

func (a *engineAdapter) CreateBackupRecord(rec *models.BackupRecord) error {
	return a.backupRepo.CreateRecord(context.Background(), rec)
}

func (a *engineAdapter) UpdateBackupRecord(opts models.UpdateBackupRecordOpts) error {
	rec, err := a.backupRepo.GetRecordByID(context.Background(), opts.ID)
	if err != nil {
		return err
	}
	rec.Status = opts.Status
	rec.FilePath = opts.FilePath
	rec.S3URL = opts.S3URL
	rec.Logs = opts.Logs
	rec.FileSizeBytes = opts.FileSizeBytes
	rec.CompletedAt = opts.CompletedAt
	return a.backupRepo.UpdateRecord(context.Background(), rec)
}

func (a *engineAdapter) GetBackupRecord(id string) (*models.BackupRecord, error) {
	return a.backupRepo.GetRecordByID(context.Background(), id)
}

func (a *engineAdapter) GetS3Destination(id string) (*models.S3Destination, error) {
	return a.s3Repo.GetS3Destination(context.Background(), id)
}

func (a *engineAdapter) ListBackupRecords(backupConfigID string) ([]*models.BackupRecord, error) {
	return a.backupRepo.ListRecordsByConfig(context.Background(), backupConfigID)
}
