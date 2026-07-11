package services

import (
	"context"
	"fmt"

	"vessl.dev/vessl/internal/repositories"
)

type ServiceLinker struct {
	databases repositories.DatabaseRepository
	storages  repositories.StorageRepository
}

func NewServiceLinker(dbRepo repositories.DatabaseRepository, stRepo repositories.StorageRepository) *ServiceLinker {
	return &ServiceLinker{databases: dbRepo, storages: stRepo}
}

func (sl *ServiceLinker) GetLinkedEnvironmentVariables(ctx context.Context, projectID string) (map[string]string, error) {
	envMap := make(map[string]string)
	if projectID == "" {
		return envMap, nil
	}
	databases, err := sl.databases.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list linked databases for project %s: %w", projectID, err)
	}
	for _, db := range databases {
		switch db.Engine {
		case "postgres":
			connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			envMap["DATABASE_URL"] = connStr
			envMap["POSTGRES_URL"] = connStr
			envMap["POSTGRES_HOST"] = db.InternalDNS
			envMap["POSTGRES_PORT"] = "5432"
			envMap["POSTGRES_USER"] = db.Username
			envMap["POSTGRES_PASSWORD"] = db.Password
			envMap["POSTGRES_DB"] = db.DatabaseName
		case "mysql":
			connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			envMap["DATABASE_URL"] = connStr
			envMap["MYSQL_URL"] = connStr
			envMap["MYSQL_HOST"] = db.InternalDNS
			envMap["MYSQL_PORT"] = "3306"
			envMap["MYSQL_USER"] = db.Username
			envMap["MYSQL_PASSWORD"] = db.Password
			envMap["MYSQL_DATABASE"] = db.DatabaseName
		case "redis":
			connStr := fmt.Sprintf("redis://:%s@%s:6379", db.Password, db.InternalDNS)
			envMap["REDIS_URL"] = connStr
			envMap["REDIS_HOST"] = db.InternalDNS
			envMap["REDIS_PORT"] = "6379"
			envMap["REDIS_PASSWORD"] = db.Password
		case "mongodb":
			connStr := fmt.Sprintf("mongodb://%s:%s@%s:27017/%s?authSource=admin", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			envMap["DATABASE_URL"] = connStr
			envMap["MONGO_URL"] = connStr
			envMap["MONGO_HOST"] = db.InternalDNS
			envMap["MONGO_PORT"] = "27017"
			envMap["MONGO_USER"] = db.Username
			envMap["MONGO_PASSWORD"] = db.Password
			envMap["MONGO_DB"] = db.DatabaseName
		}
	}
	storages, err := sl.storages.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list linked storage for project %s: %w", projectID, err)
	}
	for _, st := range storages {
		if st.Type == "minio" {
			envMap["S3_ENDPOINT"] = fmt.Sprintf("http://%s:9000", st.InternalDNS)
			envMap["S3_ACCESS_KEY"] = st.AccessKey
			envMap["S3_SECRET_KEY"] = st.SecretKey
			envMap["S3_BUCKET"] = st.BucketName
			envMap["MINIO_URL"] = fmt.Sprintf("http://%s:9000", st.InternalDNS)
			envMap["MINIO_CONSOLE_URL"] = fmt.Sprintf("http://%s:9001", st.InternalDNS)
		}
	}
	return envMap, nil
}
