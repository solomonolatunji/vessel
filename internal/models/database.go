package models

import "time"

type Database struct {
	ID                 string    `json:"id" db:"id"`
	ProjectID          string    `json:"projectId" db:"project_id"`
	EnvironmentID      string    `json:"environmentId" db:"environment_id"`
	Name               string    `json:"name" db:"name"`
	Engine             string    `json:"engine" db:"engine"`
	Version            string    `json:"version" db:"version"`
	Port               int       `json:"port" db:"port"`
	Username           string    `json:"username" db:"username"`
	Password           string    `json:"password" db:"-"`
	EncryptedPassword  string    `json:"-" db:"encrypted_password"`
	DatabaseName       string    `json:"databaseName" db:"database_name"`
	VolumePath         string    `json:"volumePath" db:"volume_path"`
	ContainerID        string    `json:"containerId" db:"container_id"`
	Status             string    `json:"status" db:"status"`
	InternalDNS        string    `json:"internalDns" db:"internal_dns"`
	ExternalDNS        string    `json:"externalDns" db:"external_dns"`
	CustomArgs         string    `json:"customArgs" db:"custom_args"`
	LogicalReplication bool      `json:"logicalReplication" db:"logical_replication"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateDatabaseRequest struct {
	ProjectID          string `json:"projectId"`
	EnvironmentID      string `json:"environmentId"`
	Name               string `json:"name"`
	Engine             string `json:"engine"`
	Version            string `json:"version"`
	Port               int    `json:"port"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	DatabaseName       string `json:"databaseName"`
	VolumePath         string `json:"volumePath"`
	CustomArgs         string `json:"customArgs"`
	LogicalReplication bool   `json:"logicalReplication"`
}

type UpdateDatabaseRequest struct {
	ExternalDNS        string `json:"externalDns"`
	CustomArgs         string `json:"customArgs"`
	LogicalReplication bool   `json:"logicalReplication"`
}

type ImportDatabaseRequest struct {
	SourceURL string `json:"sourceUrl"`
}

type Storage struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	APIPort       int       `json:"apiPort"`
	ConsolePort   int       `json:"consolePort"`
	AccessKey     string    `json:"accessKey"`
	SecretKey     string    `json:"secretKey,omitempty"`
	BucketName    string    `json:"bucketName"`
	VolumePath    string    `json:"volumePath"`
	ContainerID   string    `json:"containerId"`
	Status        string    `json:"status"`
	InternalDNS   string    `json:"internalDns"`
	ExternalDNS   string    `json:"externalDns"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type DatabaseQueryRequest struct {
	Query string `json:"query"`
}

type DatabaseQueryResponse struct {
	Columns []string         `json:"columns,omitempty"`
	Rows    []map[string]any `json:"rows,omitempty"`
	Result  any              `json:"result,omitempty"`
}

type TableSchema struct {
	Name    string         `json:"name"`
	Columns []ColumnSchema `json:"columns"`
}

type ColumnSchema struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsNullable bool   `json:"isNullable"`
	IsPrimary  bool   `json:"isPrimary"`
}
