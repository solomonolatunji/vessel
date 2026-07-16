package http

import (
	"context"
	"time"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type engineAdapter struct {
	settingsRepo   repositories.SettingsRepository
	appServiceRepo repositories.AppServiceRepository
	envRepo        repositories.EnvRepository
	databaseRepo   repositories.DatabaseRepository
	storageRepo    repositories.StorageRepository
	projectRepo    repositories.ProjectRepository
	jobRepo        repositories.JobRepository
	backupRepo     repositories.BackupRepository
	s3Repo         repositories.S3DestinationRepository
	serviceVarRepo repositories.ServiceVarRepository
	serverlessRepo repositories.ServerlessRepository
}

func newEngineAdapter(
	settingsRepo repositories.SettingsRepository,
	appServiceRepo repositories.AppServiceRepository,
	envRepo repositories.EnvRepository,
	databaseRepo repositories.DatabaseRepository,
	storageRepo repositories.StorageRepository,
	projectRepo repositories.ProjectRepository,
	jobRepo repositories.JobRepository,
	backupRepo repositories.BackupRepository,
	s3Repo repositories.S3DestinationRepository,
	serviceVarRepo repositories.ServiceVarRepository,
	serverlessRepo repositories.ServerlessRepository,
) *engineAdapter {
	return &engineAdapter{
		settingsRepo:   settingsRepo,
		appServiceRepo: appServiceRepo,
		envRepo:        envRepo,
		databaseRepo:   databaseRepo,
		storageRepo:    storageRepo,
		projectRepo:    projectRepo,
		jobRepo:        jobRepo,
		backupRepo:     backupRepo,
		s3Repo:         s3Repo,
		serviceVarRepo: serviceVarRepo,
		serverlessRepo: serverlessRepo,
	}
}

func (a *engineAdapter) GetServerSettings() (*models.ServerSettings, error) {
	return a.settingsRepo.GetServerSettings(context.Background())
}

func (a *engineAdapter) ListAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return a.appServiceRepo.ListByProject(context.Background(), projectID)
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
	db, err := a.databaseRepo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	db.Status = status
	db.ContainerID = containerID
	return a.databaseRepo.Update(context.Background(), db)
}

func (a *engineAdapter) GetDatabase(id string) (*models.Database, error) {
	return a.databaseRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) UpdateStorageStatus(id string, status models.StorageStatus, containerID string) error {
	st, err := a.storageRepo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	st.Status = status
	st.ContainerID = containerID
	return a.storageRepo.Update(context.Background(), st)
}

func (a *engineAdapter) GetStorage(id string) (*models.Storage, error) {
	return a.storageRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) ListJobs() ([]models.Job, error) {
	return a.jobRepo.ListAll(context.Background())
}

func (a *engineAdapter) GetJob(id string) (*models.Job, error) {
	return a.jobRepo.GetByID(context.Background(), id)
}

func (a *engineAdapter) GetProject(id string) (*models.ProjectConfig, error) {
	return a.projectRepo.Get(context.Background(), id)
}

func (a *engineAdapter) UpdateJobStatusAndOutput(id string, status models.JobStatus, lastRunAt *time.Time, output string) error {
	return a.jobRepo.UpdateStatus(context.Background(), id, status, lastRunAt, output)
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
