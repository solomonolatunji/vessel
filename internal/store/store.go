package store

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

type Store struct {
	mu    sync.RWMutex
	db    *sql.DB
	vault *Vault
}

func NewStore(dataDir string) (*Store, error) {
	vault, err := NewVault(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize secrets vault: %w", err)
	}

	dbPath := filepath.Join(dataDir, "vessel.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	s := &Store{db: db, vault: vault}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("failed schema migration: %w", err)
	}

	return s, nil
}

func (s *Store) migrate() error {
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
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return err
		}
	}

	_, _ = s.db.Exec("ALTER TABLE databases ADD COLUMN project_id TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE storage ADD COLUMN project_id TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE databases ADD COLUMN environment_id TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE storage ADD COLUMN environment_id TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE projects ADD COLUMN team_id TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE projects ADD COLUMN description TEXT DEFAULT '';")

	if err := s.initEnvironmentTable(); err != nil {
		return fmt.Errorf("failed to initialize environments table: %w", err)
	}
	if err := s.initAppServiceTable(); err != nil {
		return fmt.Errorf("failed to initialize app_services table: %w", err)
	}
	if err := s.initDeploymentsTable(); err != nil {
		return fmt.Errorf("failed to initialize deployments table: %w", err)
	}
	if err := s.initServiceVarsTable(); err != nil {
		return fmt.Errorf("failed to initialize service_vars table: %w", err)
	}
	if err := s.initProjectWebhooksTable(); err != nil {
		return fmt.Errorf("failed to initialize project_webhooks table: %w", err)
	}
	if err := s.initProjectTokensTable(); err != nil {
		return fmt.Errorf("failed to initialize project_tokens table: %w", err)
	}
	if err := s.initProjectMembersTable(); err != nil {
		return fmt.Errorf("failed to initialize project_members table: %w", err)
	}
	if err := s.initBackupTables(); err != nil {
		return fmt.Errorf("failed to initialize backup tables: %w", err)
	}
	if err := s.initTeamTables(); err != nil {
		return fmt.Errorf("failed to initialize team tables: %w", err)
	}
	if err := s.initWorkspaceTables(); err != nil {
		return fmt.Errorf("failed to initialize workspace tables: %w", err)
	}
	if err := s.initSettingsTables(); err != nil {
		return fmt.Errorf("failed to initialize settings tables: %w", err)
	}

	log.Println("✅ Vessel SQLite schema initialized successfully (`data/vessel.db`)")
	return nil
}

func (s *Store) initSettingsTables() error {
	queries := []string{
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
			notification_alerts BOOLEAN DEFAULT TRUE,
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
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) initTeamTables() error {
	queries := []string{
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
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	_, _ = s.db.Exec("ALTER TABLE teams ADD COLUMN avatar_url TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE teams ADD COLUMN preferred_deployment_region TEXT DEFAULT 'local';")
	return nil
}

func (s *Store) initWorkspaceTables() error {
	queries := []string{
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
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	_, _ = s.db.Exec("ALTER TABLE projects ADD COLUMN workspace_id TEXT DEFAULT '';")
	return nil
}

func (s *Store) initBackupTables() error {
	queries := []string{
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
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

// Close gracefully closes the SQLite database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
