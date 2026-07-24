---
title: Administration
description: Instance management, system updates, license management, and global settings.
---

Codedock administration covers instance-wide configuration available to instance admins.

## CLI Administration

Codedock provides three CLI tools for different environments and use cases:

### Server CLI (`codedockd`)

After installation, `codedockd` is available at `/usr/local/bin/codedockd`. This is a shell wrapper that manages the Codedock daemon by executing commands **inside the Docker container**. Use it for day-to-day server administration:

```sh
codedockd status           # Show daemon health + running containers
codedockd setup            # Interactive admin account wizard
codedockd reset-password   # Reset admin password
codedockd config           # View current configuration
codedockd config <key>=<value>  # Update a setting (site-name, registration, telemetry)
codedockd logs -f          # Tail daemon logs
codedockd update           # Upgrade to the latest version
codedockd downgrade <ver>  # Downgrade to a specific version (with backup + confirmation)
codedockd backup           # Create a manual database backup
codedockd restart          # Restart the Codedock daemon

# App management
codedockd deploy <git-url>           # Deploy an app from a Git URL
codedockd deploy --template nextjs   # Deploy a template from buildwithtechx/codedock-examples
codedockd deploy --image nginx:latest --port 80  # Deploy from a Docker image
codedockd apps:list                  # List all apps
codedockd apps:show <id>             # Show app details
codedockd apps:create <name>         # Create an app
codedockd apps:destroy <id>         # Delete an app

# Database management
codedockd db:list                    # List all databases
codedockd db:show <id>              # Show database details
codedockd db:create <name> <engine> # Create a database (postgres, mysql, redis, etc.)
codedockd db:destroy <id>           # Delete a database
```

### Daemon CLI (`codedockd`)

```sh
codedockd serve              # Start the daemon (default)
codedockd setup              # Setup wizard
codedockd reset-password     # Reset admin password
codedockd config             # View/update configuration
codedockd deploy <url>       # Deploy from Git URL
codedockd deploy --template  # Deploy a template
codedockd apps:list          # List apps
codedockd db:list            # List databases
codedockd mcp                # Run MCP stdio server
codedockd version            # Show version
```

### Remote CLI (`codedock`)

For remote management from your local machine, install the `codedock` client:

```sh
curl -fsSL https://get.codedock.run/cli | sh
codedock login    # Connect to your server
```

For a full list of remote commands, see the [CLI Reference](/cli/).

### Update with `codedockd update`

1. Shows current and latest available version.
2. Creates a pre-upgrade database backup automatically.
3. Pulls the new Docker image and recreates the container.
4. Your apps and databases experience zero downtime.

### Downgrade with `codedockd downgrade`

1. Requires you to type `downgrade` to confirm (safety gate).
2. Creates a pre-downgrade database backup automatically.
3. Pulls the specified version and recreates the container.
4. If something breaks, restore from backup: `cp /codedock/data/backups/codedock-pre-downgrade-*.db /codedock/data/codedock.db`

## Instance Settings

Access from **Settings** in the dashboard. Only the first registered user (instance admin) can modify these.

### Global Configuration

- **Wildcard Domain**: Set the base domain for all services across all projects.
- **SMTP Configuration**: Instance-wide SMTP settings for transactional emails.
- **DNS Resolvers**: Custom DNS resolvers for container networking.
- **Port Ranges**: Configure the port pool for service allocation.

### Traefik Configuration

- **SSL Provider**: Let's Encrypt configuration.
- **HTTP Redirect**: Force HTTPS across all services.
- **Dashboard**: Enable or disable the Traefik dashboard.

### AI Configuration

Codedock supports multi-provider AI integrations for features like **AI Log Diagnosis**.

- **Supported Providers**: OpenAI, Groq, Mistral, DeepSeek, xAI, Moonshot.
- **Provider Settings**: Configure your default provider, preferred models, and API keys.

## System Updates

### Automatic Update Checks

Codedock periodically checks GitHub releases for new versions. The dashboard displays a notification when an update is available.

### Manual Check

```sh
curl -X POST /api/settings/updates/check \
  -H "Authorization: Bearer vpt_xxx"
```

### Deploying an Update

1. Go to **Settings → Updates**.
2. Click **Check for Updates**.
3. If an update is available, click **Deploy Update**.
4. The system downloads and applies the update.
5. The dashboard displays real-time update progress.

### Update Process

1. New binary is downloaded from GitHub releases.
2. Database migrations are applied (backward-compatible).
3. Services are restarted gracefully.
4. Old binary is kept for rollback.

### Rollback

```sh
sudo codedockd rollback
```

Restores the previous binary and a database backup taken before the upgrade.

## Telemetry

By default, Codedock collects anonymized usage data to improve the product. This can be disabled in **Settings → Privacy**.

### What's Collected

- Instance version and uptime
- Count of projects, services, and users (no names or content)
- Deployment success/failure rates
- Enabled features and integrations

## Maintenance

### Data Directory

All persistent data is stored in the configured `CODEDOCK_DATA_DIR` (default: `data/`):

```text
data/
├── codedock.db          # SQLite database
├── .vault_key        # Encryption key (keep safe)
├── databases/        # Database volumes
├── storage/          # MinIO storage volumes
└── backups/          # Backup archives
```

### Backup

Back up your Codedock instance:

```sh
# Stop the daemon
# Copy the data directory
cp -r data/ data-backup-$(date +%Y%m%d)
# Restart the daemon
```

### Restore

1. Stop the daemon.
2. Restore the data directory from your backup.
3. Restart the daemon.
4. Database migrations run automatically on startup.
