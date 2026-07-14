package models

import "time"

type Deployment struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"serviceId"`
	EnvironmentID string    `json:"environmentId"`
	ProjectID     string    `json:"projectId"`
	Status        string    `json:"status"`
	Branch        string    `json:"branch,omitempty"`
	CommitHash    string    `json:"commitHash,omitempty"`
	CommitMessage string    `json:"commitMessage,omitempty"`
	Trigger       string    `json:"trigger,omitempty"`
	BuildLogs     string    `json:"buildLogs,omitempty"`
	ContainerID   string    `json:"containerId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	FinishedAt    time.Time `json:"finishedAt,omitempty"`
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

type AppService struct {
	ID              string    `json:"id"`
	ProjectID       string    `json:"projectId"`
	EnvironmentID   string    `json:"environmentId"`
	Name            string    `json:"name"`
	RepositoryURL   string    `json:"repositoryUrl"`
	ImageRef        string    `json:"imageRef,omitempty"`
	Branch          string    `json:"branch"`
	RootDirectory   string    `json:"rootDirectory"`
	BuildCommand    string    `json:"buildCommand"`
	StartCommand    string    `json:"startCommand"`
	DockerfilePath  string    `json:"dockerfilePath"`
	BuildEngine     string    `json:"buildEngine"`
	InternalPort    int       `json:"internalPort"`
	Domain          string    `json:"domain"`
	HealthCheckPath string    `json:"healthCheckPath"`
	ContainerID     string    `json:"containerId"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type CreateAppServiceRequest struct {
	ProjectID       string `json:"projectId"`
	Name            string `json:"name"`
	RepositoryURL   string `json:"repositoryUrl"`
	Branch          string `json:"branch"`
	RootDirectory   string `json:"rootDirectory"`
	BuildCommand    string `json:"buildCommand"`
	StartCommand    string `json:"startCommand"`
	DockerfilePath  string `json:"dockerfilePath"`
	BuildEngine     string `json:"buildEngine"`
	InternalPort    int    `json:"internalPort"`
	Domain          string `json:"domain"`
	HealthCheckPath string `json:"healthCheckPath"`
}

type UpdateAppServiceRequest struct {
	Name            string `json:"name"`
	RepositoryURL   string `json:"repositoryUrl"`
	Branch          string `json:"branch"`
	RootDirectory   string `json:"rootDirectory"`
	BuildCommand    string `json:"buildCommand"`
	StartCommand    string `json:"startCommand"`
	DockerfilePath  string `json:"dockerfilePath"`
	BuildEngine     string `json:"buildEngine"`
	InternalPort    int    `json:"internalPort"`
	Domain          string `json:"domain"`
	HealthCheckPath string `json:"healthCheckPath"`
	ContainerID     string `json:"containerId"`
	Status          string `json:"status"`
}

type Variable struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"serviceId"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	IsSecret      bool      `json:"isSecret"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
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
	ID         string     `json:"id"`
	ProjectID  string     `json:"projectId"`
	Name       string     `json:"name"`
	Schedule   string     `json:"schedule"`
	Command    string     `json:"command"`
	Status     string     `json:"status"`
	LastRunAt  *time.Time `json:"lastRunAt"`
	LastOutput string     `json:"lastOutput"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
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

type PRPreview struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"serviceId"`
	ProjectID     string    `json:"projectId"`
	PRNumber      int       `json:"prNumber"`
	Branch        string    `json:"branch"`
	CommitHash    string    `json:"commitHash"`
	Status        string    `json:"status"`
	PreviewDomain string    `json:"previewDomain"`
	ContainerID   string    `json:"containerId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
