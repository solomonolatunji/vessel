# 🛰️ Vessl

**The Ultra-Lightweight, Self-Hosted PaaS for Developers.**

Turn any bare-metal Linux VPS into your own private Vercel, Railway, or Heroku in 60 seconds with zero-downtime deployments, automated SSL routing, and an ultra-responsive web control panel.

---

## ✨ Features

- **⚡ Blazing-Fast Go Daemon (`vessld`)**: Uses native Go concurrency and official Docker SDK with `< 30MB RAM` idle overhead.
- **💻 Self-Hosted Dashboard**: Built with **Vite + TanStack Router + React + Tailwind CSS**. Features live terminal logs, real-time CPU/RAM stats, and dark-mode glassmorphism.
- **🔒 Automated Edge Routing**: Zero-config Let's Encrypt SSL/TLS certificates via Traefik v3.
- **🔐 Encrypted `.env` Vault**: AES-256 encrypted environment variables inside SQLite.
- **🗄️ 15 Managed Database Engines**: PostgreSQL, MySQL, MariaDB, MongoDB, Redis, Dragonfly, KeyDB, ClickHouse, Kafka, RabbitMQ, NATS, and one-click deployers for NocoDB, Plausible, WordPress, Gitea.
- **📦 S3-Compatible Storage**: MinIO-based object storage with one-click provisioning.
- **🔗 Auto-Injected Connection Strings**: Services get `DATABASE_URL`, `REDIS_URL`, `MONGO_URL`, `S3_ENDPOINT` automatically.
- **👥 Multi-Tenant Workspaces**: Teams, RBAC, audit logs, SSH keys, trusted domain SSO.
- **🔌 Git Integration**: Connect GitHub/GitLab, auto-deploy on push, PR previews.
- **🧩 5 Build Strategies**: Dockerfile, Railpack, Nixpacks, Buildpacks, Serverless.
- **📬 Notification Channels**: Discord, Slack, Telegram, Pushover, SMTP, Resend, generic webhooks.
- **🤖 AI Diagnostics**: Analyze failed deployments via OpenAI/Anthropic.
- **🔁 Zero-Downtime Deployments**: Health-checked container swaps with automatic rollback.

---

## 🚀 Quick Install (On any Linux VPS)

```bash
curl -fsSL https://get.vessl.dev | sh
```

Access your dashboard at `http://your-server-ip:8080`.

During installation, you can optionally set up an admin account via the terminal:

```bash
vesslctl setup    # Interactive wizard: create admin, set SSL email
```

### CLI Admin Tool (`vesslctl`)

After installation, manage Vessl entirely from the command line:

```bash
vesslctl status                  # Show daemon health + running containers
vesslctl setup                   # Create admin account (wizard)
vesslctl reset-password          # Reset admin password
vesslctl config                  # View configuration
vesslctl config site-name=MyVessl   # Update a setting
vesslctl logs -f                 # Tail daemon logs
vesslctl update                  # Upgrade to latest version
vesslctl downgrade v0.1.0        # Downgrade to a specific version
vesslctl backup                  # Manual database backup
vesslctl restart                 # Restart the daemon

# Deploy & manage apps
vesslctl deploy https://github.com/user/app.git   # Deploy from Git URL
vesslctl apps:list                                 # List all apps
vesslctl apps:show <id>                            # App details
vesslctl apps:create my-app --project <id>         # Create an app
vesslctl apps:destroy <id>                         # Delete an app

# Manage databases
vesslctl db:list                                   # List databases
vesslctl db:show <id>                              # Database details
vesslctl db:create my-db postgres --project <id>   # Create a database
vesslctl db:destroy <id>                           # Delete a database
```

---

## 📖 Documentation

Full docs at [docs.vessl.dev](https://docs.vessl.dev), including:
- Getting started guide
- Deployment (build strategies, domains, env vars, CI/CD)
- Database management (all 15 engines, backups)
- Storage (S3-compatible MinIO)
- Workspaces & teams
- Serverless functions
- Integrations (Git, OAuth, Vercel import, webhooks)
- Configuration (notifications, AI, 2FA, email)
- Administration (instance settings, updates)
- API reference (REST API, personal access tokens)

---

## 🏗️ Repository Layout

```text
vessl/
├── cmd/                  # Go daemon + CLI subcommands
│   ├── main.go           # Entrypoint, server startup
│   ├── cli.go            # CLI routing
│   ├── setup.go          # Setup wizard
│   ├── reset_password.go # Password reset
│   └── config.go         # Config management
├── internal/             # Core Go packages
│   ├── engine/           # Docker container orchestration
│   ├── handlers/         # HTTP handlers
│   ├── http/             # Server setup, routes, middleware
│   ├── models/           # Domain models
│   ├── services/         # Business logic
│   ├── repositories/     # SQLite persistence
│   ├── notifications/    # Channel integrations
│   ├── core/             # Event bus, dispatchers
│   └── utils/            # Helpers (vault, network, units)
├── dashboard/            # React/Vite frontend
├── web/                  # Marketing site
├── docs/                 # Starlight documentation
├── bootstrap/
│   ├── install.sh        # One-line install script
│   └── vesslctl          # CLI admin wrapper (installed to /usr/local/bin)
└── scripts/              # Helper scripts (upgrade, backup, restore, downgrade)
```

---

## ⚡ Local Development

1. **Prerequisites**: Go 1.25+, Node.js 22+, Docker
2. **Setup**:
   ```bash
   cp .env.example .env
   # Edit .env — set VESSL_JWT_SECRET to a random 32-char string
   ```
3. **Run** (daemon only):
   ```bash
   go run ./cmd           # vessld serve on :8080
   ```
   Or for the full stack (daemon + dashboard + docs):
   ```bash
   make dev
   ```

### CLI Subcommands (dev mode)

The `vessld` binary supports CLI commands even outside Docker:

```bash
go run ./cmd setup           # Run setup wizard
go run ./cmd reset-password  # Reset admin password
go run ./cmd config           # View config
go run ./cmd version          # Show version
```

### Environment Variables (.env)

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | No | `8080` | Daemon HTTP port |
| `VESSL_DATA_DIR` | No | `data` | Data directory (DB, vault, builds) |
| `VESSL_JWT_SECRET` | **Yes** | — | JWT signing key (generate a 32-char random string) |
| `VESSL_DASHBOARD_URL` | No | `http://localhost:8080` | Public dashboard URL for notification links |
| `DOCKER_SOCKET_PATH` | No | `/var/run/docker.sock` | Docker socket path (rootless/podman/colima) |
| `VESSL_TLS_EMAIL` | No | — | Email for Let's Encrypt SSL certificate notifications |
| `VESSL_RUNTIME_NETWORK` | No | `vessl-network` | Docker network name for services |
| `DEPLOY_HOST_PORT_START` | No | `4100` | Start of host port range for deployments |
| `DEPLOY_HOST_PORT_END` | No | `4999` | End of host port range |
| `VESSL_TRAEFIK_IMAGE` | No | `traefik:v3.0` | Traefik Docker image tag |
| `VESSL_MINIO_IMAGE` | No | `minio/minio:latest` | MinIO Docker image tag |
| `VESSL_DEFAULT_MEMORY_MB` | No | `512` | Default memory limit for app containers (MB) |
| `VESSL_DEFAULT_CPU` | No | `0.5` | Default CPU request for app containers |
| `VESSL_DEFAULT_DB_MEMORY_MB` | No | `1024` | Default memory limit for database containers (MB) |
| `VESSL_DEFAULT_DB_CPU` | No | `1.0` | Default CPU request for database containers |

---

## 📄 License

Vessl Source-Available License. You are free to view, use, and modify the code for personal or internal business use. However, redistribution, reselling, or using Vessl to provide a competing commercial managed PaaS is strictly prohibited without explicit written permission. See `LICENSE` for details.
