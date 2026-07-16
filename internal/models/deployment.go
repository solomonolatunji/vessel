package models

import "time"

type Deployment struct {
	ID            string           `json:"id" db:"id"`
	ServiceID     string           `json:"serviceId" db:"service_id"`
	EnvironmentID string           `json:"environmentId" db:"environment_id"`
	ProjectID     string           `json:"projectId" db:"project_id"`
	Status        DeploymentStatus `json:"status" db:"status"`
	Branch        string           `json:"branch,omitempty" db:"branch"`
	CommitHash    string           `json:"commitHash,omitempty" db:"commit_hash"`
	CommitMessage string           `json:"commitMessage,omitempty" db:"commit_message"`
	Trigger       string           `json:"trigger,omitempty" db:"trigger"`
	BuildLogs     string           `json:"buildLogs,omitempty" db:"build_logs"`
	ContainerID   string           `json:"containerId,omitempty" db:"container_id"`
	CreatedAt     time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time        `json:"updatedAt" db:"updated_at"`
	FinishedAt    *time.Time       `json:"finishedAt,omitempty" db:"finished_at"`
}

type ServiceMetric struct {
	Timestamp  string  `json:"timestamp"`
	CPUPercent float64 `json:"cpuPercent"`
	MemoryMB   float64 `json:"memoryMB"`
	NetworkRx  float64 `json:"networkRxKB"`
	NetworkTx  float64 `json:"networkTxKB"`
}

type TriggerDeploymentRequest struct {
	Branch *string `json:"branch,omitempty"`
}

type RuntimeMode string

const (
	RuntimeModeWeb    RuntimeMode = "web"
	RuntimeModeWorker RuntimeMode = "worker"
)

type DeploymentStatus string

const (
	DeploymentStatusPending  DeploymentStatus = "pending"
	DeploymentStatusCloning  DeploymentStatus = "CLONING"
	DeploymentStatusPulling  DeploymentStatus = "PULLING"
	DeploymentStatusBuilding DeploymentStatus = "BUILDING"
	DeploymentStatusReady    DeploymentStatus = "READY"
	DeploymentStatusActive   DeploymentStatus = "ACTIVE"
	DeploymentStatusFailed   DeploymentStatus = "FAILED"
	DeploymentStatusRemoved  DeploymentStatus = "REMOVED"
	DeploymentStatusSlept    DeploymentStatus = "SLEPT"
)

type AppServiceStatus string

const (
	AppServiceStatusCreated  AppServiceStatus = "created"
	AppServiceStatusBuilding AppServiceStatus = "building"
	AppServiceStatusStopped  AppServiceStatus = "stopped"
	AppServiceStatusRunning  AppServiceStatus = "running"
	AppServiceStatusError    AppServiceStatus = "error"
)

type BuildEngine string

const (
	BuildEngineAuto       BuildEngine = "auto"
	BuildEngineDockerfile BuildEngine = "dockerfile"
	BuildEngineNixpacks   BuildEngine = "nixpacks"
	BuildEngineBuildpacks BuildEngine = "buildpacks"
	BuildEngineRailpack   BuildEngine = "railpack"
	BuildEngineServerless BuildEngine = "serverless"
)

type JobStatus string

const (
	JobStatusActive    JobStatus = "active"
	JobStatusInactive  JobStatus = "inactive"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusError     JobStatus = "error"
)

type BackupConfigStatus string

const (
	BackupConfigStatusActive   BackupConfigStatus = "active"
	BackupConfigStatusInactive BackupConfigStatus = "inactive"
)

type BackupRecordStatus string

const (
	BackupRecordStatusRunning   BackupRecordStatus = "running"
	BackupRecordStatusCompleted BackupRecordStatus = "completed"
	BackupRecordStatusFailed    BackupRecordStatus = "failed"
	BackupRecordStatusExpired   BackupRecordStatus = "expired"
)

type PRPreviewStatus string

const (
	PRPreviewStatusPending PRPreviewStatus = "PENDING"
	PRPreviewStatusReady   PRPreviewStatus = "READY"
	PRPreviewStatusFailed  PRPreviewStatus = "FAILED"
)

type AppService struct {
	ID              string           `json:"id" db:"id"`
	ProjectID       string           `json:"projectId" db:"project_id"`
	EnvironmentID   string           `json:"environmentId" db:"environment_id"`
	Name            string           `json:"name" db:"name"`
	RepositoryURL   string           `json:"repositoryUrl" db:"repository_url"`
	ImageRef        string           `json:"imageRef,omitempty" db:"image_ref"`
	Branch          string           `json:"branch" db:"branch"`
	RootDirectory   string           `json:"rootDirectory" db:"root_directory"`
	RuntimeMode     RuntimeMode      `json:"runtimeMode" db:"runtime_mode"`
	InstallCommand  string           `json:"installCommand" db:"install_command"`
	BuildCommand    string           `json:"buildCommand" db:"build_command"`
	StartCommand    string           `json:"startCommand" db:"start_command"`
	DockerfilePath  string           `json:"dockerfilePath" db:"dockerfile_path"`
	BuildEngine     BuildEngine      `json:"buildEngine" db:"build_engine"`
	InternalPort    int              `json:"internalPort" db:"internal_port"`
	Domain          string           `json:"domain" db:"domain"`
	StaticOutput    string           `json:"staticOutput" db:"static_output"`
	HealthCheckPath string           `json:"healthCheckPath" db:"health_check_path"`
	ContainerID     string           `json:"containerId" db:"container_id"`
	Status          AppServiceStatus `json:"status" db:"status"`
	CreatedAt       time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time        `json:"updatedAt" db:"updated_at"`
}

type CreateAppServiceRequest struct {
	ProjectID       string      `json:"projectId"`
	Name            string      `json:"name"`
	RepositoryURL   string      `json:"repositoryUrl"`
	Branch          string      `json:"branch"`
	RootDirectory   string      `json:"rootDirectory"`
	RuntimeMode     RuntimeMode `json:"runtimeMode"`
	InstallCommand  string      `json:"installCommand"`
	BuildCommand    string      `json:"buildCommand"`
	StartCommand    string      `json:"startCommand"`
	DockerfilePath  string      `json:"dockerfilePath"`
	BuildEngine     string      `json:"buildEngine"`
	InternalPort    int         `json:"internalPort"`
	Domain          string      `json:"domain"`
	StaticOutput    string      `json:"staticOutput"`
	HealthCheckPath string      `json:"healthCheckPath"`
}

type UpdateAppServiceRequest struct {
	Name            string      `json:"name"`
	RepositoryURL   string      `json:"repositoryUrl"`
	Branch          string      `json:"branch"`
	RootDirectory   string      `json:"rootDirectory"`
	RuntimeMode     RuntimeMode `json:"runtimeMode"`
	InstallCommand  string      `json:"installCommand"`
	BuildCommand    string      `json:"buildCommand"`
	StartCommand    string      `json:"startCommand"`
	DockerfilePath  string      `json:"dockerfilePath"`
	BuildEngine     string      `json:"buildEngine"`
	InternalPort    int         `json:"internalPort"`
	Domain          string      `json:"domain"`
	StaticOutput    string      `json:"staticOutput"`
	HealthCheckPath string      `json:"healthCheckPath"`
	ContainerID     string      `json:"containerId"`
	Status          string      `json:"status"`
}

type Variable struct {
	ID            string    `json:"id" db:"id"`
	ServiceID     string    `json:"serviceId" db:"service_id"`
	ProjectID     string    `json:"projectId" db:"project_id"`
	EnvironmentID string    `json:"environmentId" db:"environment_id"`
	Key           string    `json:"key" db:"key"`
	Value         string    `json:"value" db:"value"`
	IsSecret      bool      `json:"isSecret" db:"is_secret"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateServiceVarRequest struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
}

type UpdateServiceVarRequest struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
}

type Job struct {
	ID         string     `json:"id" db:"id"`
	ProjectID  string     `json:"projectId" db:"project_id"`
	Name       string     `json:"name" db:"name"`
	Schedule   string     `json:"schedule" db:"schedule"`
	Command    string     `json:"command" db:"command"`
	Status     JobStatus  `json:"status" db:"status"`
	LastRunAt  *time.Time `json:"lastRunAt" db:"last_run_at"`
	LastOutput string     `json:"lastOutput" db:"last_output"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
}

type CreateJobRequest struct {
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Schedule  string `json:"schedule"`
	Command   string `json:"command"`
}

type UpdateJobRequest struct {
	Name     *string `json:"name,omitempty"`
	Schedule *string `json:"schedule,omitempty"`
	Command  *string `json:"command,omitempty"`
	Status   *string `json:"status,omitempty"`
}

type BackupConfig struct {
	ID              string             `json:"id" db:"id"`
	ProjectID       string             `json:"projectId" db:"project_id"`
	DatabaseID      string             `json:"databaseId,omitempty" db:"database_id"`
	StorageID       string             `json:"storageId,omitempty" db:"storage_id"`
	S3DestinationID string             `json:"s3DestinationId,omitempty" db:"s3_destination_id"`
	Name            string             `json:"name" db:"name"`
	Schedule        string             `json:"schedule" db:"schedule"`
	RetentionDays   int                `json:"retentionDays" db:"retention_days"`
	Status          BackupConfigStatus `json:"status" db:"status"`
	CreatedAt       string             `json:"createdAt" db:"created_at"`
	UpdatedAt       string             `json:"updatedAt" db:"updated_at"`
}

type BackupRecord struct {
	ID             string             `json:"id" db:"id"`
	BackupConfigID string             `json:"backupConfigId" db:"backup_config_id"`
	ProjectID      string             `json:"projectId" db:"project_id"`
	DatabaseID     string             `json:"databaseId,omitempty" db:"database_id"`
	Status         BackupRecordStatus `json:"status" db:"status"`
	FilePath       string             `json:"filePath" db:"file_path"`
	FileSizeBytes  int64              `json:"fileSizeBytes" db:"file_size_bytes"`
	S3URL          string             `json:"s3Url,omitempty" db:"s3_url"`
	Logs           string             `json:"logs" db:"logs"`
	StartedAt      string             `json:"startedAt" db:"started_at"`
	CompletedAt    string             `json:"completedAt" db:"completed_at"`
}

type UpdateBackupRecordOpts struct {
	ID            string
	Status        BackupRecordStatus
	FilePath      string
	S3URL         string
	Logs          string
	FileSizeBytes int64
	CompletedAt   string
}

type S3Destination struct {
	ID              string `json:"id" db:"id"`
	ProjectID       string `json:"projectId" db:"project_id"`
	Name            string `json:"name" db:"name"`
	Endpoint        string `json:"endpoint" db:"endpoint"`
	Bucket          string `json:"bucket" db:"bucket"`
	Region          string `json:"region" db:"region"`
	AccessKeyID     string `json:"accessKeyId" db:"access_key_id"`
	SecretAccessKey string `json:"secretAccessKey" db:"secret_access_key"`
	CreatedAt       string `json:"createdAt" db:"created_at"`
}

type PRPreview struct {
	ID            string          `json:"id" db:"id"`
	ServiceID     string          `json:"serviceId" db:"service_id"`
	ProjectID     string          `json:"projectId" db:"project_id"`
	PRNumber      int             `json:"prNumber" db:"pr_number"`
	Branch        string          `json:"branch" db:"branch"`
	CommitHash    string          `json:"commitHash" db:"commit_hash"`
	Status        PRPreviewStatus `json:"status" db:"status"`
	PreviewDomain string          `json:"previewDomain" db:"preview_domain"`
	ContainerID   string          `json:"containerId" db:"container_id"`
	CreatedAt     time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time       `json:"updatedAt" db:"updated_at"`
}
