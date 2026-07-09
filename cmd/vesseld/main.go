package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/client"
	_ "modernc.org/sqlite"
	"vessel.dev/vessel/internal/api"
	"vessel.dev/vessel/internal/engine"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/repositories"
	"vessel.dev/vessel/internal/vault"
)

const vesselVersion = "0.1.0-alpha"

type dbProjectLister struct{ db *sql.DB }

func (a *dbProjectLister) ListProjects() ([]models.ProjectConfig, error) {
	rows, err := a.db.Query(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ProjectConfig
	for rows.Next() {
		var p models.ProjectConfig
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

type dbServiceLister struct{ db *sql.DB }

func (a *dbServiceLister) ListServices() ([]models.AppService, error) {
	rows, err := a.db.Query(`SELECT id, project_id, environment_id, name, repository_url, branch, internal_port, domain, container_id, status, created_at, updated_at FROM app_services ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.AppService
	for rows.Next() {
		var svc models.AppService
		if err := rows.Scan(&svc.ID, &svc.ProjectID, &svc.EnvironmentID, &svc.Name, &svc.RepositoryURL, &svc.Branch, &svc.InternalPort, &svc.Domain, &svc.ContainerID, &svc.Status, &svc.CreatedAt, &svc.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, svc)
	}
	return list, rows.Err()
}

type dbDomainLister struct{ db *sql.DB }

func (a *dbDomainLister) ListAllDomains() ([]models.DomainConfig, error) {
	rows, err := a.db.Query(`SELECT id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at FROM domains ORDER BY domain_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.DomainConfig
	for rows.Next() {
		var d models.DomainConfig
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.DomainName, &d.RedirectTo, &d.SSLCertStatus, &d.PathPrefix, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}



type dbDeployerStore struct {
	db    *sql.DB
	vault *vault.Vault
}

func (a *dbDeployerStore) GetServerSettings() (*models.ServerSettings, error) {
	return repositories.NewSettingsSQLiteRepository(a.db).GetServerSettings(context.Background())
}

func (a *dbDeployerStore) ListAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return repositories.NewAppServiceSQLiteRepository(a.db).ListByProject(context.Background(), projectID)
}

func (a *dbDeployerStore) GetEnvVars(projectID string) (map[string]string, error) {
	return repositories.NewEnvSQLiteRepository(a.db, a.vault).GetVars(context.Background(), projectID)
}

func (a *dbDeployerStore) ListServiceVariables(serviceID string) ([]*models.Variable, error) {
	svVarRepo := repositories.NewServiceVarSQLiteRepository(a.db)
	return svVarRepo.ListByService(context.Background(), serviceID)
}

func main() {
	log.Printf(" Booting Vessel Daemon (`vesseld`) v%s [%s/%s]...", vesselVersion, runtime.GOOS, runtime.GOARCH)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dataDir := os.Getenv("VESSEL_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf(" Failed to create data directory: %v", err)
	}

	vlt, err := vault.NewVault(dataDir)
	if err != nil {
		log.Fatalf(" Failed to initialize secrets vault: %v", err)
	}

	dbPath := filepath.Join(dataDir, "vessel.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		log.Fatalf(" Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatalf(" Failed schema migration: %v", err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf(" Docker daemon connection warning: %v (container deployment features disabled)", err)
	}

	proxyCfg := proxy.NewCaddyConfig(dataDir, os.Getenv("VESSEL_TLS_EMAIL"))
	proxyMgr := proxy.NewProxyManager(proxyCfg, &dbProjectLister{db: db}, &dbServiceLister{db: db}, &dbDomainLister{db: db}, dockerClient)
	_ = proxyMgr.Reload(context.Background())

	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})

	apiServer := api.NewServer(db, vlt, deployer, proxyMgr, dockerClient)

	log.Printf(" Vessel control plane listening on :%s", port)
	if err := http.ListenAndServe(":"+port, apiServer.Handler()); err != nil {
		log.Fatalf(" Server crashed: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			team_id TEXT DEFAULT '',
			name TEXT UNIQUE NOT NULL,
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS domains (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			domain_name TEXT UNIQUE NOT NULL,
			redirect_to TEXT,
			ssl_cert_status TEXT DEFAULT 'pending',
			path_prefix TEXT DEFAULT '/',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT DEFAULT 'developer',
			totp_enabled BOOLEAN DEFAULT FALSE,
			totp_secret TEXT DEFAULT '',
			recovery_codes TEXT DEFAULT '',
			oauth_provider TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			role TEXT DEFAULT 'developer',
			token TEXT UNIQUE NOT NULL,
			invited_by TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			accepted_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS env_vars (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			key TEXT NOT NULL,
			encrypted_value TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			UNIQUE(project_id, key)
		);`,
		`CREATE TABLE IF NOT EXISTS databases (
			id TEXT PRIMARY KEY,
			project_id TEXT DEFAULT '',
			name TEXT UNIQUE NOT NULL,
			engine TEXT NOT NULL,
			version TEXT NOT NULL,
			port INTEGER NOT NULL,
			username TEXT NOT NULL,
			encrypted_password TEXT NOT NULL,
			database_name TEXT NOT NULL,
			volume_path TEXT NOT NULL,
			container_id TEXT,
			status TEXT DEFAULT 'stopped',
			internal_dns TEXT,
			external_dns TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS storage (
			id TEXT PRIMARY KEY,
			project_id TEXT DEFAULT '',
			name TEXT UNIQUE NOT NULL,
			type TEXT DEFAULT 'minio',
			api_port INTEGER DEFAULT 9000,
			console_port INTEGER DEFAULT 9001,
			access_key TEXT NOT NULL,
			encrypted_secret_key TEXT NOT NULL,
			bucket_name TEXT NOT NULL,
			volume_path TEXT NOT NULL,
			container_id TEXT,
			status TEXT DEFAULT 'stopped',
			internal_dns TEXT,
			external_dns TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			schedule TEXT NOT NULL,
			command TEXT NOT NULL,
			status TEXT DEFAULT 'active',
			last_run_at DATETIME,
			last_output TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS user_git_providers (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			provider TEXT NOT NULL,
			encrypted_access_token TEXT NOT NULL,
			account_name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, provider)
		);`,
		`CREATE TABLE IF NOT EXISTS environments (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			is_default BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE(project_id, name)
		);`,
		`CREATE TABLE IF NOT EXISTS app_services (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			environment_id TEXT NOT NULL,
			name TEXT NOT NULL,
			icon TEXT DEFAULT 'git',
			repository_url TEXT DEFAULT '',
			branch TEXT DEFAULT 'main',
			root_directory TEXT DEFAULT '/',
			build_command TEXT DEFAULT '',
			start_command TEXT DEFAULT '',
			dockerfile_path TEXT DEFAULT '',
			internal_port INTEGER DEFAULT 3000,
			domain TEXT DEFAULT '',
			env_vars_count INTEGER DEFAULT 0,
			auto_deploy_webhook BOOLEAN DEFAULT 1,
			cpu_request REAL DEFAULT 0.5,
			memory_limit_mb INTEGER DEFAULT 512,
			replicas INTEGER DEFAULT 1,
			restart_policy TEXT DEFAULT 'on_failure',
			teardown_timeout INTEGER DEFAULT 30,
			serverless BOOLEAN DEFAULT 0,
			cron_schedule TEXT DEFAULT '',
			health_check_path TEXT DEFAULT '/',
			status TEXT DEFAULT 'building',
			container_id TEXT DEFAULT '',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE(environment_id, name)
		);`,
		`CREATE TABLE IF NOT EXISTS deployments (
			id TEXT PRIMARY KEY,
			service_id TEXT NOT NULL,
			environment_id TEXT NOT NULL,
			project_id TEXT NOT NULL,
			status TEXT NOT NULL,
			commit_hash TEXT DEFAULT '',
			commit_message TEXT DEFAULT '',
			branch TEXT DEFAULT '',
			trigger TEXT DEFAULT 'Manual',
			build_logs TEXT DEFAULT '',
			container_id TEXT DEFAULT '',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			finished_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS service_vars (
			id TEXT PRIMARY KEY,
			service_id TEXT NOT NULL,
			environment_id TEXT,
			key TEXT NOT NULL,
			value TEXT NOT NULL,
			is_secret INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			UNIQUE(service_id, key)
		);`,
		`CREATE TABLE IF NOT EXISTS project_webhooks (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			provider TEXT DEFAULT 'github',
			webhook_secret TEXT,
			webhook_url TEXT,
			auto_deploy BOOLEAN DEFAULT 1,
			branch TEXT DEFAULT 'main',
			created_at TEXT,
			updated_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS project_tokens (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			token TEXT NOT NULL,
			scopes TEXT DEFAULT '',
			expires_at TEXT,
			last_used_at TEXT,
			created_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS project_members (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			user_email TEXT DEFAULT '',
			role TEXT NOT NULL,
			joined_at TEXT,
			UNIQUE(project_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS backup_configs (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			database_id TEXT,
			storage_id TEXT,
			s3_destination_id TEXT,
			name TEXT NOT NULL,
			schedule TEXT NOT NULL,
			retention_days INTEGER DEFAULT 7,
			status TEXT DEFAULT 'active',
			created_at TEXT,
			updated_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS backup_records (
			id TEXT PRIMARY KEY,
			backup_config_id TEXT NOT NULL,
			project_id TEXT NOT NULL,
			database_id TEXT,
			status TEXT DEFAULT 'running',
			file_path TEXT,
			file_size_bytes INTEGER DEFAULT 0,
			s3_url TEXT,
			logs TEXT,
			started_at TEXT,
			completed_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS s3_destinations (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			endpoint TEXT NOT NULL,
			bucket TEXT NOT NULL,
			region TEXT NOT NULL,
			access_key_id TEXT NOT NULL,
			secret_access_key TEXT NOT NULL,
			created_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			avatar_url TEXT DEFAULT '',
			preferred_deployment_region TEXT DEFAULT 'local',
			owner_id TEXT NOT NULL,
			created_at TEXT,
			updated_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS team_members (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			user_email TEXT,
			role TEXT NOT NULL,
			joined_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS team_invites (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			email TEXT NOT NULL,
			role TEXT NOT NULL,
			token TEXT UNIQUE NOT NULL,
			invited_by TEXT NOT NULL,
			expires_at TEXT,
			created_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS workspaces (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			avatar_url TEXT DEFAULT '',
			preferred_region TEXT DEFAULT 'local',
			owner_id TEXT NOT NULL,
			created_at TEXT,
			updated_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS workspace_trusted_domains (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			domain TEXT NOT NULL,
			role TEXT DEFAULT 'developer',
			created_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS workspace_ssh_keys (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			name TEXT NOT NULL,
			public_key TEXT NOT NULL,
			created_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS workspace_audit_logs (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			project_id TEXT,
			environment_id TEXT,
			action TEXT NOT NULL,
			actor TEXT NOT NULL,
			created_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS server_settings (
			id TEXT PRIMARY KEY,
			caddy_wildcard_ip TEXT DEFAULT '127.0.0.1',
			discord_webhook_url TEXT,
			slack_webhook_url TEXT,
			telegram_bot_token TEXT,
			telegram_chat_id TEXT,
			smtp_host TEXT,
			smtp_port INTEGER DEFAULT 587,
			smtp_user TEXT,
			smtp_password TEXT,
			smtp_from_name TEXT DEFAULT '',
			smtp_from_address TEXT DEFAULT '',
			notification_alerts BOOLEAN DEFAULT TRUE,
			registration_enabled BOOLEAN DEFAULT TRUE,
			custom_dns_resolvers TEXT DEFAULT '',
			dns_validation_enabled BOOLEAN DEFAULT TRUE,
			ip_allowlist TEXT DEFAULT '',
			mcp_server_enabled BOOLEAN DEFAULT TRUE,
			update_check_cron TEXT DEFAULT '0 * * * *',
			auto_update_enabled BOOLEAN DEFAULT FALSE,
			current_version TEXT DEFAULT '0.1.0',
			latest_version TEXT DEFAULT '0.1.0',
			last_update_check TEXT DEFAULT '',
			updated_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS personal_access_tokens (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			token_hash TEXT NOT NULL,
			prefix TEXT NOT NULL,
			expires_at TEXT,
			created_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS notification_integrations (
			id TEXT PRIMARY KEY,
			smtp_enabled BOOLEAN DEFAULT FALSE,
			smtp_host TEXT,
			smtp_port INTEGER DEFAULT 587,
			smtp_user TEXT,
			smtp_password TEXT,
			smtp_from_name TEXT,
			smtp_from_address TEXT,
			resend_enabled BOOLEAN DEFAULT FALSE,
			resend_api_key TEXT,
			slack_enabled BOOLEAN DEFAULT FALSE,
			slack_webhook_url TEXT,
			discord_enabled BOOLEAN DEFAULT FALSE,
			discord_webhook_url TEXT,
			discord_ping_enabled BOOLEAN DEFAULT FALSE,
			telegram_enabled BOOLEAN DEFAULT FALSE,
			telegram_bot_token TEXT,
			telegram_chat_id TEXT,
			pushover_enabled BOOLEAN DEFAULT FALSE,
			pushover_user_key TEXT,
			pushover_api_token TEXT,
			webhook_enabled BOOLEAN DEFAULT FALSE,
			webhook_url TEXT,
			updated_at TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS project_notification_prefs (
			project_id TEXT PRIMARY KEY,
			email_enabled BOOLEAN DEFAULT TRUE,
			slack_enabled BOOLEAN DEFAULT TRUE,
			discord_enabled BOOLEAN DEFAULT TRUE,
			telegram_enabled BOOLEAN DEFAULT TRUE,
			pushover_enabled BOOLEAN DEFAULT TRUE,
			webhook_enabled BOOLEAN DEFAULT TRUE,
			events TEXT DEFAULT 'deploy.success,deploy.failure,invite',
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS oauth_providers (
			id TEXT PRIMARY KEY,
			provider_name TEXT UNIQUE NOT NULL,
			enabled BOOLEAN DEFAULT FALSE,
			client_id TEXT DEFAULT '',
			client_secret TEXT DEFAULT '',
			redirect_uri TEXT DEFAULT '',
			base_url TEXT DEFAULT '',
			tenant TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	alterQueries := []string{
		"ALTER TABLE databases ADD COLUMN project_id TEXT DEFAULT '';",
		"ALTER TABLE storage ADD COLUMN project_id TEXT DEFAULT '';",
		"ALTER TABLE databases ADD COLUMN environment_id TEXT DEFAULT '';",
		"ALTER TABLE storage ADD COLUMN environment_id TEXT DEFAULT '';",
		"ALTER TABLE projects ADD COLUMN team_id TEXT DEFAULT '';",
		"ALTER TABLE projects ADD COLUMN description TEXT DEFAULT '';",
		"ALTER TABLE projects ADD COLUMN workspace_id TEXT DEFAULT '';",
		"ALTER TABLE users ADD COLUMN totp_enabled BOOLEAN DEFAULT FALSE;",
		"ALTER TABLE users ADD COLUMN totp_secret TEXT DEFAULT '';",
		"ALTER TABLE users ADD COLUMN recovery_codes TEXT DEFAULT '';",
		"ALTER TABLE users ADD COLUMN oauth_provider TEXT DEFAULT '';",
		"ALTER TABLE server_settings ADD COLUMN registration_enabled BOOLEAN DEFAULT TRUE;",
		"ALTER TABLE server_settings ADD COLUMN custom_dns_resolvers TEXT DEFAULT '';",
		"ALTER TABLE server_settings ADD COLUMN dns_validation_enabled BOOLEAN DEFAULT TRUE;",
		"ALTER TABLE server_settings ADD COLUMN ip_allowlist TEXT DEFAULT '';",
		"ALTER TABLE server_settings ADD COLUMN mcp_server_enabled BOOLEAN DEFAULT TRUE;",
		"ALTER TABLE server_settings ADD COLUMN update_check_cron TEXT DEFAULT '0 * * * *';",
		"ALTER TABLE server_settings ADD COLUMN auto_update_enabled BOOLEAN DEFAULT FALSE;",
		"ALTER TABLE server_settings ADD COLUMN current_version TEXT DEFAULT '0.1.0';",
		"ALTER TABLE server_settings ADD COLUMN latest_version TEXT DEFAULT '0.1.0';",
		"ALTER TABLE server_settings ADD COLUMN last_update_check TEXT DEFAULT '';",
		"ALTER TABLE server_settings ADD COLUMN smtp_from_name TEXT DEFAULT '';",
		"ALTER TABLE server_settings ADD COLUMN smtp_from_address TEXT DEFAULT '';",
		"ALTER TABLE app_services ADD COLUMN icon TEXT DEFAULT 'git';",
		"ALTER TABLE app_services ADD COLUMN root_directory TEXT DEFAULT '/';",
		"ALTER TABLE app_services ADD COLUMN replicas INTEGER DEFAULT 1;",
		"ALTER TABLE app_services ADD COLUMN restart_policy TEXT DEFAULT 'on_failure';",
		"ALTER TABLE app_services ADD COLUMN teardown_timeout INTEGER DEFAULT 30;",
		"ALTER TABLE app_services ADD COLUMN serverless BOOLEAN DEFAULT 0;",
		"ALTER TABLE app_services ADD COLUMN cron_schedule TEXT DEFAULT '';",
		"ALTER TABLE app_services ADD COLUMN git_repo_full_name TEXT DEFAULT '';",
		"ALTER TABLE app_services ADD COLUMN wait_for_ci BOOLEAN DEFAULT 1;",
		"ALTER TABLE app_services ADD COLUMN auto_deploy_branch BOOLEAN DEFAULT 1;",
		"ALTER TABLE app_services ADD COLUMN public_networking_domain TEXT DEFAULT '';",
		"ALTER TABLE app_services ADD COLUMN private_networking_internal TEXT DEFAULT '';",
		"ALTER TABLE app_services ADD COLUMN enable_outbound_ipv6 BOOLEAN DEFAULT 0;",
		"ALTER TABLE notification_integrations ADD COLUMN smtp_from_name TEXT DEFAULT '';",
		"ALTER TABLE notification_integrations ADD COLUMN smtp_from_address TEXT DEFAULT '';",
		"ALTER TABLE teams ADD COLUMN avatar_url TEXT DEFAULT '';",
		"ALTER TABLE teams ADD COLUMN preferred_deployment_region TEXT DEFAULT 'local';",
	}

	for _, q := range alterQueries {
		_, _ = db.Exec(q)
	}

	log.Println(" Vessel SQLite schema initialized successfully (`data/vessel.db`)")
	return nil
}
