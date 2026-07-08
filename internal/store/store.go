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
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return err
		}
	}
	log.Println("✅ Vessel SQLite schema initialized successfully (`data/vessel.db`)")
	return nil
}

// Close gracefully closes the SQLite database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
