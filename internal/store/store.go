package store

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Store encapsulates the embedded SQLite 3 database and AES secret vault.
type Store struct {
	db    *sql.DB
	vault *Vault
}

// NewStore initializes modernc pure-Go SQLite inside the specified data directory and runs schema migrations.
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
			name TEXT UNIQUE NOT NULL,
			repository_url TEXT,
			branch TEXT,
			build_command TEXT,
			start_command TEXT,
			dockerfile_path TEXT,
			internal_port INTEGER DEFAULT 3000,
			domain TEXT,
			auto_deploy_webhook BOOLEAN DEFAULT 1,
			cpu_request REAL DEFAULT 0.5,
			memory_limit_mb INTEGER DEFAULT 512,
			health_check_path TEXT DEFAULT '/',
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
	log.Println("✅ Vessel SQLite schema initialized successfully (`data/vessel.db`)")
	return nil
}

// Close gracefully closes the SQLite database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
