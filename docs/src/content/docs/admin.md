---
title: Administration
description: Instance management, system updates, license management, and global settings.
---

Vessl administration covers instance-wide configuration available to instance admins.

## CLI Admin Tool

After installation, `vesslctl` is available at `/usr/local/bin/vesslctl` for managing Vessl from the terminal.

```sh
vesslctl status           # Show daemon health + running containers
vesslctl setup            # Interactive admin account wizard
vesslctl reset-password   # Reset admin password
vesslctl config           # View current configuration
vesslctl config <key>=<value>  # Update a setting (site-name, registration, telemetry)
vesslctl logs -f          # Tail daemon logs
vesslctl update           # Upgrade to the latest version
vesslctl downgrade <ver>  # Downgrade to a specific version (with backup + confirmation)
vesslctl backup           # Create a manual database backup
vesslctl restart          # Restart the Vessl daemon

# App management
vesslctl deploy <git-url>           # Deploy an app from a Git URL
vesslctl deploy --image nginx:latest --port 80  # Deploy from a Docker image
vesslctl apps:list                  # List all apps
vesslctl apps:show <id>             # Show app details
vesslctl apps:create <name>         # Create an app
vesslctl apps:destroy <id>         # Delete an app

# Database management
vesslctl db:list                    # List all databases
vesslctl db:show <id>              # Show database details
vesslctl db:create <name> <engine> # Create a database (postgres, mysql, redis, etc.)
vesslctl db:destroy <id>           # Delete a database
```

For development or standalone (non-Docker) mode, the `vessld` binary supports the same subcommands:

```sh
vessld serve              # Start the daemon (default)
vessld setup              # Setup wizard
vessld reset-password     # Reset admin password
vessld config             # View/update configuration
vessld deploy <url>       # Deploy from Git URL
vessld apps:list          # List apps
vessld db:list            # List databases
vessld mcp                # Run MCP stdio server
vessld version            # Show version
```

### Update with `vesslctl update`

1. Shows current and latest available version.
2. Creates a pre-upgrade database backup automatically.
3. Pulls the new Docker image and recreates the container.
4. Your apps and databases experience zero downtime.

### Downgrade with `vesslctl downgrade`

1. Requires you to type `downgrade` to confirm (safety gate).
2. Creates a pre-downgrade database backup automatically.
3. Pulls the specified version and recreates the container.
4. If something breaks, restore from backup: `cp /vessl/data/backups/vessl-pre-downgrade-*.db /vessl/data/vessl.db`

## Instance Settings

Access from **Settings** in the dashboard. Only the first registered user (instance admin) can modify these.

### Global Configuration

- **Wildcard Domain**: Set the base domain for all services across all workspaces.
- **SMTP Configuration**: Instance-wide SMTP settings for transactional emails.
- **DNS Resolvers**: Custom DNS resolvers for container networking.
- **Port Ranges**: Configure the port pool for service allocation.

### Traefik Configuration

- **SSL Provider**: Let's Encrypt configuration.
- **HTTP Redirect**: Force HTTPS across all services.
- **Dashboard**: Enable or disable the Traefik dashboard.

## System Updates

### Automatic Update Checks

Vessl periodically checks GitHub releases for new versions. The dashboard displays a notification when an update is available.

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
sudo vesslctl rollback
```

Restores the previous binary and a database backup taken before the upgrade.

## License Management

### Activating a License

1. Go to **Settings → License**.
2. Enter your license key.
3. Click **Activate**.

The license is validated against the licensing server and applied immediately.

### License Features

Plans may include:

- Seat limits (number of users)
- Workspace limits
- Premium features (audit logs, SSO, advanced RBAC)

## Telemetry

By default, Vessl collects anonymized usage data to improve the product. This can be disabled in **Settings → Privacy**.

### What's Collected

- Instance version and uptime
- Count of projects, services, and users (no names or content)
- Deployment success/failure rates
- Enabled features and integrations

## Maintenance

### Data Directory

All persistent data is stored in the configured `VESSL_DATA_DIR` (default: `data/`):

```text
data/
├── vessl.db          # SQLite database
├── .vault_key        # Encryption key (keep safe)
├── databases/        # Database volumes
├── storage/          # MinIO storage volumes
└── backups/          # Backup archives
```

### Backup

Back up your Vessl instance:

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
