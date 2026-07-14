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
		case "clickhouse":
			envMap["CLICKHOUSE_URL"] = fmt.Sprintf("clickhouse://%s:%s@%s:8123/%s", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			envMap["CLICKHOUSE_HOST"] = db.InternalDNS
			envMap["CLICKHOUSE_PORT"] = "8123"
			envMap["CLICKHOUSE_USER"] = db.Username
			envMap["CLICKHOUSE_PASSWORD"] = db.Password
			envMap["CLICKHOUSE_DB"] = db.DatabaseName
		case "timescaledb":
			connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			envMap["DATABASE_URL"] = connStr
			envMap["TIMESCALE_URL"] = connStr
			envMap["PGHOST"] = db.InternalDNS
			envMap["PGPORT"] = "5432"
			envMap["PGUSER"] = db.Username
			envMap["PGPASSWORD"] = db.Password
			envMap["PGDATABASE"] = db.DatabaseName
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

func (sl *ServiceLinker) GetNamespacedVariables(ctx context.Context, projectID string) (map[string]map[string]string, error) {
	registry := make(map[string]map[string]string)
	if projectID == "" {
		return registry, nil
	}

	databases, err := sl.databases.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list databases for interpolation: %w", err)
	}

	for _, db := range databases {
		vars := make(map[string]string)
		switch db.Engine {
		case "postgres":
			connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			vars["DATABASE_URL"] = connStr
			vars["POSTGRES_URL"] = connStr
			vars["POSTGRES_HOST"] = db.InternalDNS
			vars["POSTGRES_PORT"] = "5432"
			vars["POSTGRES_USER"] = db.Username
			vars["POSTGRES_PASSWORD"] = db.Password
			vars["POSTGRES_DB"] = db.DatabaseName
		case "mysql", "mariadb":
			connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			vars["DATABASE_URL"] = connStr
			vars["MYSQL_URL"] = connStr
			vars["MYSQL_HOST"] = db.InternalDNS
			vars["MYSQL_PORT"] = "3306"
			vars["MYSQL_USER"] = db.Username
			vars["MYSQL_PASSWORD"] = db.Password
			vars["MYSQL_DATABASE"] = db.DatabaseName
		case "redis":
			connStr := fmt.Sprintf("redis://:%s@%s:6379", db.Password, db.InternalDNS)
			vars["REDIS_URL"] = connStr
			vars["REDIS_HOST"] = db.InternalDNS
			vars["REDIS_PORT"] = "6379"
			vars["REDIS_PASSWORD"] = db.Password
		case "mongodb":
			connStr := fmt.Sprintf("mongodb://%s:%s@%s:27017/%s?authSource=admin", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			vars["DATABASE_URL"] = connStr
			vars["MONGO_URL"] = connStr
			vars["MONGO_HOST"] = db.InternalDNS
			vars["MONGO_PORT"] = "27017"
			vars["MONGO_USER"] = db.Username
			vars["MONGO_PASSWORD"] = db.Password
			vars["MONGO_DB"] = db.DatabaseName
		case "clickhouse":
			vars["CLICKHOUSE_URL"] = fmt.Sprintf("clickhouse://%s:%s@%s:8123/%s", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			vars["CLICKHOUSE_HOST"] = db.InternalDNS
			vars["CLICKHOUSE_PORT"] = "8123"
			vars["CLICKHOUSE_USER"] = db.Username
			vars["CLICKHOUSE_PASSWORD"] = db.Password
			vars["CLICKHOUSE_DB"] = db.DatabaseName
		case "timescaledb":
			connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", db.Username, db.Password, db.InternalDNS, db.DatabaseName)
			vars["DATABASE_URL"] = connStr
			vars["TIMESCALE_URL"] = connStr
			vars["PGHOST"] = db.InternalDNS
			vars["PGPORT"] = "5432"
			vars["PGUSER"] = db.Username
			vars["PGPASSWORD"] = db.Password
			vars["PGDATABASE"] = db.DatabaseName
		}
		if len(vars) > 0 {
			registry[db.Name] = vars
		}
	}

	storages, err := sl.storages.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list storages for interpolation: %w", err)
	}
	for _, st := range storages {
		if st.Type == "minio" {
			registry[st.Name] = map[string]string{
				"S3_ENDPOINT":       fmt.Sprintf("http://%s:9000", st.InternalDNS),
				"S3_ACCESS_KEY":     st.AccessKey,
				"S3_SECRET_KEY":     st.SecretKey,
				"S3_BUCKET":         st.BucketName,
				"MINIO_URL":         fmt.Sprintf("http://%s:9000", st.InternalDNS),
				"MINIO_CONSOLE_URL": fmt.Sprintf("http://%s:9001", st.InternalDNS),
			}
		}
	}

	return registry, nil
}
