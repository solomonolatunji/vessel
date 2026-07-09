package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// initAppServiceTable creates the app_services table and runs safe migrations if needed.
func (s *Store) initAppServiceTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS app_services (
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
	);`
	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	// Safe backward compatibility migrations
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN icon TEXT DEFAULT 'git';")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN root_directory TEXT DEFAULT '/';")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN replicas INTEGER DEFAULT 1;")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN restart_policy TEXT DEFAULT 'on_failure';")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN teardown_timeout INTEGER DEFAULT 30;")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN serverless BOOLEAN DEFAULT 0;")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN cron_schedule TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN git_repo_full_name TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN wait_for_ci BOOLEAN DEFAULT 1;")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN auto_deploy_branch BOOLEAN DEFAULT 1;")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN public_networking_domain TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN private_networking_internal TEXT DEFAULT '';")
	_, _ = s.db.Exec("ALTER TABLE app_services ADD COLUMN enable_outbound_ipv6 BOOLEAN DEFAULT 0;")
	return nil
}

// CreateAppService adds a new Git deployable service (e.g., recovery, wallet-bot) inside an environment.
func (s *Store) CreateAppService(service *types.AppServiceConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if service.ID == "" {
		service.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	service.CreatedAt = now
	service.UpdatedAt = now
	if service.Status == "" {
		service.Status = "building"
	}
	if service.Icon == "" {
		service.Icon = "git"
	}
	if service.RootDirectory == "" {
		service.RootDirectory = "/"
	}
	if service.Replicas <= 0 {
		service.Replicas = 1
	}
	if service.RestartPolicy == "" {
		service.RestartPolicy = "on_failure"
	}
	if service.TeardownTimeout <= 0 {
		service.TeardownTimeout = 30
	}

	query := `INSERT INTO app_services (
		id, project_id, environment_id, name, icon, repository_url, branch, root_directory, build_command,
		start_command, dockerfile_path, internal_port, domain, env_vars_count,
		auto_deploy_webhook, git_repo_full_name, wait_for_ci, auto_deploy_branch, public_networking_domain,
		private_networking_internal, enable_outbound_ipv6, cpu_request, memory_limit_mb, replicas, restart_policy, teardown_timeout,
		serverless, cron_schedule, health_check_path, status, container_id, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		service.ID, service.ProjectID, service.EnvironmentID, service.Name, service.Icon, service.RepositoryURL, service.Branch, service.RootDirectory,
		service.BuildCommand, service.StartCommand, service.DockerfilePath, service.InternalPort, service.Domain,
		service.EnvVarsCount, service.AutoDeployWebhook, service.GitRepoFullName, service.WaitForCI, service.AutoDeployBranch, service.PublicNetworkingDomain,
		service.PrivateNetworkingInternal, service.EnableOutboundIPv6, service.CPURequest, service.MemoryLimitMB, service.Replicas, service.RestartPolicy, service.TeardownTimeout,
		service.Serverless, service.CronSchedule, service.HealthCheckPath, service.Status, service.ContainerID, service.CreatedAt, service.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create app service: %w", err)
	}
	return nil
}

// GetAppService retrieves an application service by its ID.
func (s *Store) GetAppService(id string) (*types.AppServiceConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, environment_id, name, icon, repository_url, branch, root_directory, build_command,
		start_command, dockerfile_path, internal_port, domain, env_vars_count,
		auto_deploy_webhook, COALESCE(git_repo_full_name, ''), COALESCE(wait_for_ci, 1), COALESCE(auto_deploy_branch, 1),
		COALESCE(public_networking_domain, ''), COALESCE(private_networking_internal, ''), COALESCE(enable_outbound_ipv6, 0),
		cpu_request, memory_limit_mb, replicas, restart_policy, teardown_timeout,
		serverless, cron_schedule, health_check_path, status, container_id, created_at, updated_at FROM app_services WHERE id = ?`

	row := s.db.QueryRow(query, id)
	var app types.AppServiceConfig
	var autoDeploy, waitForCI, autoDeployBranch, enableIPv6, serverless int
	err := row.Scan(
		&app.ID, &app.ProjectID, &app.EnvironmentID, &app.Name, &app.Icon, &app.RepositoryURL, &app.Branch, &app.RootDirectory,
		&app.BuildCommand, &app.StartCommand, &app.DockerfilePath, &app.InternalPort, &app.Domain,
		&app.EnvVarsCount, &autoDeploy, &app.GitRepoFullName, &waitForCI, &autoDeployBranch,
		&app.PublicNetworkingDomain, &app.PrivateNetworkingInternal, &enableIPv6,
		&app.CPURequest, &app.MemoryLimitMB, &app.Replicas, &app.RestartPolicy, &app.TeardownTimeout,
		&serverless, &app.CronSchedule, &app.HealthCheckPath, &app.Status, &app.ContainerID, &app.CreatedAt, &app.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("app service not found: %s", id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get app service: %w", err)
	}
	app.AutoDeployWebhook = autoDeploy == 1
	app.WaitForCI = waitForCI == 1
	app.AutoDeployBranch = autoDeployBranch == 1
	app.EnableOutboundIPv6 = enableIPv6 == 1
	app.Serverless = serverless == 1
	return &app, nil
}

// ListAppServicesByEnvironment returns all application container services running inside a specific environment.
func (s *Store) ListAppServicesByEnvironment(environmentID string) ([]*types.AppServiceConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, environment_id, name, icon, repository_url, branch, root_directory, build_command,
		start_command, dockerfile_path, internal_port, domain, env_vars_count,
		auto_deploy_webhook, COALESCE(git_repo_full_name, ''), COALESCE(wait_for_ci, 1), COALESCE(auto_deploy_branch, 1),
		COALESCE(public_networking_domain, ''), COALESCE(private_networking_internal, ''), COALESCE(enable_outbound_ipv6, 0),
		cpu_request, memory_limit_mb, replicas, restart_policy, teardown_timeout,
		serverless, cron_schedule, health_check_path, status, container_id, created_at, updated_at FROM app_services WHERE environment_id = ? ORDER BY created_at ASC`

	rows, err := s.db.Query(query, environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by environment: %w", err)
	}
	defer rows.Close()

	var apps []*types.AppServiceConfig
	for rows.Next() {
		var app types.AppServiceConfig
		var autoDeploy, waitForCI, autoDeployBranch, enableIPv6, serverless int
		if err := rows.Scan(
			&app.ID, &app.ProjectID, &app.EnvironmentID, &app.Name, &app.Icon, &app.RepositoryURL, &app.Branch, &app.RootDirectory,
			&app.BuildCommand, &app.StartCommand, &app.DockerfilePath, &app.InternalPort, &app.Domain,
			&app.EnvVarsCount, &autoDeploy, &app.GitRepoFullName, &waitForCI, &autoDeployBranch,
			&app.PublicNetworkingDomain, &app.PrivateNetworkingInternal, &enableIPv6,
			&app.CPURequest, &app.MemoryLimitMB, &app.Replicas, &app.RestartPolicy, &app.TeardownTimeout,
			&serverless, &app.CronSchedule, &app.HealthCheckPath, &app.Status, &app.ContainerID, &app.CreatedAt, &app.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan app service row: %w", err)
		}
		app.AutoDeployWebhook = autoDeploy == 1
		app.WaitForCI = waitForCI == 1
		app.AutoDeployBranch = autoDeployBranch == 1
		app.EnableOutboundIPv6 = enableIPv6 == 1
		app.Serverless = serverless == 1
		apps = append(apps, &app)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return apps, nil
}

// ListAppServicesByProject returns all application container services across all environments in a project.
func (s *Store) ListAppServicesByProject(projectID string) ([]*types.AppServiceConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, environment_id, name, icon, repository_url, branch, root_directory, build_command,
		start_command, dockerfile_path, internal_port, domain, env_vars_count,
		auto_deploy_webhook, COALESCE(git_repo_full_name, ''), COALESCE(wait_for_ci, 1), COALESCE(auto_deploy_branch, 1),
		COALESCE(public_networking_domain, ''), COALESCE(private_networking_internal, ''), COALESCE(enable_outbound_ipv6, 0),
		cpu_request, memory_limit_mb, replicas, restart_policy, teardown_timeout,
		serverless, cron_schedule, health_check_path, status, container_id, created_at, updated_at FROM app_services WHERE project_id = ? ORDER BY created_at ASC`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by project: %w", err)
	}
	defer rows.Close()

	var apps []*types.AppServiceConfig
	for rows.Next() {
		var app types.AppServiceConfig
		var autoDeploy, waitForCI, autoDeployBranch, enableIPv6, serverless int
		if err := rows.Scan(
			&app.ID, &app.ProjectID, &app.EnvironmentID, &app.Name, &app.Icon, &app.RepositoryURL, &app.Branch, &app.RootDirectory,
			&app.BuildCommand, &app.StartCommand, &app.DockerfilePath, &app.InternalPort, &app.Domain,
			&app.EnvVarsCount, &autoDeploy, &app.GitRepoFullName, &waitForCI, &autoDeployBranch,
			&app.PublicNetworkingDomain, &app.PrivateNetworkingInternal, &enableIPv6,
			&app.CPURequest, &app.MemoryLimitMB, &app.Replicas, &app.RestartPolicy, &app.TeardownTimeout,
			&serverless, &app.CronSchedule, &app.HealthCheckPath, &app.Status, &app.ContainerID, &app.CreatedAt, &app.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan app service row: %w", err)
		}
		app.AutoDeployWebhook = autoDeploy == 1
		app.WaitForCI = waitForCI == 1
		app.AutoDeployBranch = autoDeployBranch == 1
		app.EnableOutboundIPv6 = enableIPv6 == 1
		app.Serverless = serverless == 1
		apps = append(apps, &app)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return apps, nil
}

// UpdateAppService modifies existing service settings (`Settings` tab in UI).
func (s *Store) UpdateAppService(service *types.AppServiceConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	service.UpdatedAt = time.Now().UTC()
	query := `UPDATE app_services SET
		name = ?, icon = ?, repository_url = ?, branch = ?, root_directory = ?, build_command = ?,
		start_command = ?, dockerfile_path = ?, internal_port = ?, domain = ?,
		auto_deploy_webhook = ?, git_repo_full_name = ?, wait_for_ci = ?, auto_deploy_branch = ?,
		public_networking_domain = ?, private_networking_internal = ?, enable_outbound_ipv6 = ?,
		cpu_request = ?, memory_limit_mb = ?, replicas = ?, restart_policy = ?,
		teardown_timeout = ?, serverless = ?, cron_schedule = ?, health_check_path = ?, updated_at = ?
		WHERE id = ?`

	_, err := s.db.Exec(query,
		service.Name, service.Icon, service.RepositoryURL, service.Branch, service.RootDirectory, service.BuildCommand,
		service.StartCommand, service.DockerfilePath, service.InternalPort, service.Domain,
		service.AutoDeployWebhook, service.GitRepoFullName, service.WaitForCI, service.AutoDeployBranch,
		service.PublicNetworkingDomain, service.PrivateNetworkingInternal, service.EnableOutboundIPv6,
		service.CPURequest, service.MemoryLimitMB, service.Replicas, service.RestartPolicy,
		service.TeardownTimeout, service.Serverless, service.CronSchedule, service.HealthCheckPath, service.UpdatedAt,
		service.ID,
	)
	return err
}

// UpdateAppServiceStatus updates the container status and container ID of an application service.
func (s *Store) UpdateAppServiceStatus(id, status, containerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `UPDATE app_services SET status = ?, container_id = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, status, containerID, time.Now().UTC(), id)
	return err
}

// ListAllAppServices returns all application container services across the entire platform.
func (s *Store) ListAllAppServices() ([]*types.AppServiceConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, environment_id, name, icon, repository_url, branch, root_directory, build_command,
		start_command, dockerfile_path, internal_port, domain, env_vars_count,
		auto_deploy_webhook, COALESCE(git_repo_full_name, ''), COALESCE(wait_for_ci, 1), COALESCE(auto_deploy_branch, 1),
		COALESCE(public_networking_domain, ''), COALESCE(private_networking_internal, ''), COALESCE(enable_outbound_ipv6, 0),
		cpu_request, memory_limit_mb, replicas, restart_policy, teardown_timeout,
		serverless, cron_schedule, health_check_path, status, container_id, created_at, updated_at FROM app_services ORDER BY created_at ASC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all app services: %w", err)
	}
	defer rows.Close()

	var apps []*types.AppServiceConfig
	for rows.Next() {
		var app types.AppServiceConfig
		var autoDeploy, waitForCI, autoDeployBranch, enableIPv6, serverless int
		if err := rows.Scan(
			&app.ID, &app.ProjectID, &app.EnvironmentID, &app.Name, &app.Icon, &app.RepositoryURL, &app.Branch, &app.RootDirectory,
			&app.BuildCommand, &app.StartCommand, &app.DockerfilePath, &app.InternalPort, &app.Domain,
			&app.EnvVarsCount, &autoDeploy, &app.GitRepoFullName, &waitForCI, &autoDeployBranch,
			&app.PublicNetworkingDomain, &app.PrivateNetworkingInternal, &enableIPv6,
			&app.CPURequest, &app.MemoryLimitMB, &app.Replicas, &app.RestartPolicy, &app.TeardownTimeout,
			&serverless, &app.CronSchedule, &app.HealthCheckPath, &app.Status, &app.ContainerID, &app.CreatedAt, &app.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan app service row: %w", err)
		}
		app.AutoDeployWebhook = autoDeploy == 1
		app.WaitForCI = waitForCI == 1
		app.AutoDeployBranch = autoDeployBranch == 1
		app.EnableOutboundIPv6 = enableIPv6 == 1
		app.Serverless = serverless == 1
		apps = append(apps, &app)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return apps, nil
}

// DeleteAppService deletes an application service configuration.
func (s *Store) DeleteAppService(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM app_services WHERE id = ?`, id)
	return err
}
