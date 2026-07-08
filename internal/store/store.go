package store

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
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

// CreateProject inserts a new ProjectConfig record into SQLite.
func (s *Store) CreateProject(p *types.ProjectConfig) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	_, err := s.db.Exec(`INSERT INTO projects (
		id, name, repository_url, branch, build_command, start_command, dockerfile_path,
		internal_port, domain, auto_deploy_webhook, cpu_request, memory_limit_mb, health_check_path,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.RepositoryURL, p.Branch, p.BuildCommand, p.StartCommand, p.DockerfilePath,
		p.InternalPort, p.Domain, p.AutoDeployWebhook, p.CPURequest, p.MemoryLimitMB, p.HealthCheckPath,
		p.CreatedAt, p.UpdatedAt,
	)
	return err
}

// GetProject retrieves a ProjectConfig record by its ID.
func (s *Store) GetProject(id string) (*types.ProjectConfig, error) {
	row := s.db.QueryRow(`SELECT id, name, repository_url, branch, build_command, start_command, dockerfile_path,
		internal_port, domain, auto_deploy_webhook, cpu_request, memory_limit_mb, health_check_path, created_at, updated_at
		FROM projects WHERE id = ?`, id)

	var p types.ProjectConfig
	err := row.Scan(&p.ID, &p.Name, &p.RepositoryURL, &p.Branch, &p.BuildCommand, &p.StartCommand, &p.DockerfilePath,
		&p.InternalPort, &p.Domain, &p.AutoDeployWebhook, &p.CPURequest, &p.MemoryLimitMB, &p.HealthCheckPath,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// AddDomain registers a new custom domain routing rule for a project.
func (s *Store) AddDomain(d *types.DomainConfig) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now

	_, err := s.db.Exec(`INSERT INTO domains (id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ProjectID, d.DomainName, d.RedirectTo, d.SSLCertStatus, d.PathPrefix, d.CreatedAt, d.UpdatedAt)
	return err
}

// ListDomains returns all custom domain configurations attached to the specified project ID.
func (s *Store) ListDomains(projectID string) ([]types.DomainConfig, error) {
	rows, err := s.db.Query(`SELECT id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at
		FROM domains WHERE project_id = ? ORDER BY domain_name ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []types.DomainConfig
	for rows.Next() {
		var d types.DomainConfig
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.DomainName, &d.RedirectTo, &d.SSLCertStatus, &d.PathPrefix, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, nil
}

// SetEnvVar encrypts a plaintext environment variable value and stores it in SQLite.
func (s *Store) SetEnvVar(projectID, key, plaintextValue string) error {
	encrypted, err := s.vault.Encrypt(plaintextValue)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = s.db.Exec(`INSERT INTO env_vars (id, project_id, key, encrypted_value, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(project_id, key) DO UPDATE SET encrypted_value = excluded.encrypted_value, updated_at = excluded.updated_at`,
		uuid.NewString(), projectID, key, encrypted, now, now)
	return err
}

// GetEnvVars retrieves and decrypts all environment variables for a given project ID.
func (s *Store) GetEnvVars(projectID string) (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, encrypted_value FROM env_vars WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	envs := make(map[string]string)
	for rows.Next() {
		var key, encrypted string
		if err := rows.Scan(&key, &encrypted); err != nil {
			return nil, err
		}
		plaintext, err := s.vault.Decrypt(encrypted)
		if err != nil {
			continue
		}
		envs[key] = plaintext
	}
	return envs, nil
}

// CreateUser registers a new authenticated user in SQLite.
func (s *Store) CreateUser(u *types.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err := s.db.Exec(`INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`, u.ID, u.Email, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt)
	return err
}

// CreateInvite issues a new workspace invitation token with a 7-day expiration.
func (s *Store) CreateInvite(inv *types.Invite) error {
	if inv.ID == "" {
		inv.ID = uuid.NewString()
	}
	if inv.Token == "" {
		inv.Token = uuid.NewString()
	}
	inv.CreatedAt = time.Now()
	if inv.ExpiresAt.IsZero() {
		inv.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	}

	_, err := s.db.Exec(`INSERT INTO invites (id, email, role, token, invited_by, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, inv.ID, inv.Email, inv.Role, inv.Token, inv.InvitedBy, inv.ExpiresAt, inv.CreatedAt)
	return err
}
