CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS domains (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			domain_name TEXT UNIQUE NOT NULL,
			redirect_to TEXT,
			ssl_cert_status TEXT DEFAULT 'pending',
			path_prefix TEXT DEFAULT '/',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);

CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			name TEXT DEFAULT '',
			password_hash TEXT NOT NULL,
			role TEXT DEFAULT 'developer',
			totp_enabled BOOLEAN DEFAULT FALSE,
			totp_secret TEXT DEFAULT '',
			recovery_codes TEXT DEFAULT '',
			oauth_provider TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			role TEXT DEFAULT 'developer',
			token TEXT UNIQUE NOT NULL,
			invited_by TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			accepted_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

CREATE TABLE IF NOT EXISTS env_vars (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			key TEXT NOT NULL,
			encrypted_value TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			UNIQUE(project_id, key)
		);

CREATE TABLE IF NOT EXISTS databases (
			id TEXT PRIMARY KEY,
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
			custom_args TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		,
			environment_id TEXT DEFAULT '',
            logical_replication INTEGER DEFAULT 0,
            project_id TEXT DEFAULT '');

CREATE TABLE IF NOT EXISTS storage (
			id TEXT PRIMARY KEY,
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
		,
			environment_id TEXT DEFAULT '',
            project_id TEXT DEFAULT '');

CREATE TABLE IF NOT EXISTS jobs (
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
		);

CREATE TABLE IF NOT EXISTS user_git_providers (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			provider TEXT NOT NULL,
			encrypted_access_token TEXT NOT NULL,
			account_name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, provider)
		);

CREATE TABLE IF NOT EXISTS user_vercel_accounts (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			encrypted_access_token TEXT NOT NULL,
			vercel_team_id TEXT, -- Vercel team ID if they authenticated a team, or NULL for personal account
			account_name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, vercel_team_id)
		);

CREATE TABLE IF NOT EXISTS environments (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			is_default BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE(project_id, name)
		);

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
			build_engine TEXT DEFAULT 'railpack',
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
			git_repo_full_name TEXT DEFAULT '',
			wait_for_ci BOOLEAN DEFAULT 1,
			auto_deploy_branch BOOLEAN DEFAULT 1,
			public_networking_domain TEXT DEFAULT '',
			private_networking_internal TEXT DEFAULT '',
			enable_outbound_ipv6 BOOLEAN DEFAULT 0,
			image_ref TEXT NOT NULL DEFAULT '',
			runtime_mode TEXT DEFAULT 'web',
			install_command TEXT DEFAULT '',
			static_output TEXT DEFAULT '',
			UNIQUE(environment_id, name)
		);

CREATE TABLE IF NOT EXISTS deployments (
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
		);

CREATE TABLE IF NOT EXISTS pr_previews (
			id TEXT PRIMARY KEY,
			service_id TEXT,
			project_id TEXT,
			pr_number INTEGER,
			branch TEXT,
			commit_hash TEXT,
			status TEXT,
			preview_domain TEXT,
			container_id TEXT,
			created_at DATETIME,
			updated_at DATETIME
		);

CREATE TABLE IF NOT EXISTS service_vars (
			id TEXT PRIMARY KEY,
			service_id TEXT NOT NULL,
			environment_id TEXT,
			key TEXT NOT NULL,
			value TEXT NOT NULL,
			is_secret INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			UNIQUE(service_id, key)
		);

CREATE TABLE IF NOT EXISTS project_webhooks (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			provider TEXT DEFAULT 'github',
			webhook_secret TEXT,
			webhook_url TEXT,
			auto_deploy BOOLEAN DEFAULT 1,
			branch TEXT DEFAULT 'main',
			created_at TEXT,
			updated_at TEXT,
			url TEXT DEFAULT '',
			event_types TEXT DEFAULT '',
			include_pr_environments BOOLEAN DEFAULT FALSE
		);

CREATE TABLE IF NOT EXISTS serverless_functions_code (
			id TEXT PRIMARY KEY,
			service_id TEXT NOT NULL,
			runtime TEXT NOT NULL,
			code_content TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			UNIQUE(service_id)
		);

CREATE TABLE IF NOT EXISTS project_tokens (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			environment_id TEXT NOT NULL,
			name TEXT NOT NULL,
			token_prefix TEXT NOT NULL,
			token_hash TEXT NOT NULL,
			scopes TEXT DEFAULT '',
			ip_allowlist TEXT DEFAULT '',
			expires_at TEXT,
			last_used_at TEXT,
			created_at TEXT
		);

CREATE TABLE IF NOT EXISTS project_members (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			user_email TEXT DEFAULT '',
			role TEXT NOT NULL,
			joined_at TEXT,
			email TEXT DEFAULT '',
			permission TEXT DEFAULT 'Can Edit',
			status TEXT DEFAULT 'pending',
			invited_at TEXT,
			accepted_at TEXT,
			UNIQUE(project_id, user_id)
		);

CREATE TABLE IF NOT EXISTS backup_configs (
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
		);

CREATE TABLE IF NOT EXISTS backup_records (
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
		);

CREATE TABLE IF NOT EXISTS s3_destinations (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			provider TEXT DEFAULT 's3',
			endpoint TEXT NOT NULL,
			bucket TEXT NOT NULL,
			region TEXT NOT NULL,
			access_key_id TEXT NOT NULL,
			secret_access_key TEXT NOT NULL,
			created_at TEXT
		);





CREATE TABLE IF NOT EXISTS server_settings (
			id TEXT PRIMARY KEY,
			traefik_wildcard_ip TEXT DEFAULT '127.0.0.1',
			registration_enabled BOOLEAN DEFAULT TRUE,
			registration_domain_allowlist TEXT DEFAULT '',
			custom_dns_resolvers TEXT DEFAULT '',
			dns_validation_enabled BOOLEAN DEFAULT TRUE,
			ip_allowlist TEXT DEFAULT '',
			mcp_server_enabled BOOLEAN DEFAULT TRUE,
			default_wildcard_domain TEXT DEFAULT '',
			update_check_cron TEXT DEFAULT '0 * * * *',
			auto_update_enabled BOOLEAN DEFAULT FALSE,
			current_version TEXT DEFAULT '0.1.0',
			latest_version TEXT DEFAULT '0.1.0',
			last_update_check TEXT DEFAULT '',
			updated_at TEXT,
			site_name TEXT DEFAULT '',
			public_ipv4 TEXT DEFAULT '',
			public_ipv6 TEXT DEFAULT '',
			show_sponsorship_popup BOOLEAN DEFAULT 1,
			disable_two_step_confirmation BOOLEAN DEFAULT 0,
			panel_domain TEXT NOT NULL DEFAULT '',
			concurrent_builds INTEGER NOT NULL DEFAULT 2,
			deployment_timeout INTEGER NOT NULL DEFAULT 3600,
			server_timezone TEXT NOT NULL DEFAULT 'UTC',
			docker_cleanup_cron TEXT NOT NULL DEFAULT '0 0 * * *',
			disk_usage_threshold INTEGER NOT NULL DEFAULT 80,
			disk_usage_cron TEXT NOT NULL DEFAULT '0 23 * * *',
			cloudflare_api_token TEXT DEFAULT '',
			namecheap_api_user TEXT DEFAULT '',
			namecheap_api_key TEXT DEFAULT '',
			namecheap_client_ip TEXT DEFAULT '',
			spaceship_api_key TEXT DEFAULT ''
		);

CREATE TABLE IF NOT EXISTS personal_access_tokens (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			token_hash TEXT NOT NULL,
			prefix TEXT NOT NULL,
			access_level TEXT DEFAULT 'read_write',
			project_scope TEXT DEFAULT 'all',
			allowed_projects TEXT,
			expires_at TEXT,
			created_at TEXT
		);



CREATE TABLE IF NOT EXISTS github_apps (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			app_id TEXT NOT NULL,
			installation_id TEXT,
			client_id TEXT NOT NULL,
			client_secret TEXT NOT NULL,
			webhook_secret TEXT NOT NULL,
			private_key TEXT NOT NULL,
			is_public BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

CREATE TABLE IF NOT EXISTS gitlab_apps (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			app_id TEXT NOT NULL,
			app_secret TEXT NOT NULL,
			webhook_secret TEXT NOT NULL,
			api_url TEXT NOT NULL DEFAULT 'https://gitlab.com',
			is_public BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

CREATE TABLE IF NOT EXISTS bitbucket_apps (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			owner TEXT NOT NULL,
			client_id TEXT NOT NULL,
			client_secret TEXT NOT NULL,
			webhook_secret TEXT NOT NULL,
			is_public BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

CREATE TABLE IF NOT EXISTS oauth_providers (
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
		);
