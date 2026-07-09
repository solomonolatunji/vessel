# 🛰️ Vessel Development Roadmap & Tasks (`TODO.md`)

## 📌 Phase 1: Foundation & Core Layout (COMPLETED ✅)

- [x] Establish standard Go Cloud Monorepo architecture (`cmd/`, `internal/`, `dashboard/`, `web/`, `bootstrap/`, `scripts/`, `data/`)
- [x] Initialize Go daemon skeleton (`cmd/vesseld/main.go` + `internal/types/types.go`)
- [x] Initialize Dashboard GUI (`dashboard/`) with TanStack Start/Router + Vite + Tailwind CSS v4 + Radix UI + Lucide
- [x] Initialize Marketing Landing Page (`web/`) with Astro 7 + Tailwind CSS
- [x] Create automated zero-downtime self-upgrade and backup scripts (`scripts/upgrade.sh`, `scripts/backup.sh`, `scripts/restore.sh`)
- [x] Create open-source governance files (`LICENSE`, `CONTRIBUTING.md`, `SECURITY.md`, `Makefile`)

---

## ⚙️ Phase 2: Go Backend Engine (`cmd/vesseld` & `internal/`)

- [x] **Docker Engine Client & Multi-Language Builders (`internal/engine/`)**:

  The core bridge where `vesseld` talks to the host Docker daemon (`/var/run/docker.sock`) using the official Docker Go SDK (`github.com/docker/docker/client`). Supports any language/framework out of the box (Node.js, Python, Go, Rust, PHP, Ruby, Java, etc.) without requiring users to write their own Dockerfile — just like Railway, Coolify, and Vercel.

  ```text
  internal/engine/
  ├── builder.go              # Builder Interface & Strategy Dispatcher
  ├── dockerfile_builder.go   # Standard Docker build strategy (when Dockerfile exists)
  ├── railpack_builder.go     # Railpack / Nixpacks / Buildpacks auto-detect OCI builder
  ├── container_manager.go    # Docker SDK client (Start/Stop/Inspect/Logs)
  ├── deployer.go             # End-to-End Zero-Downtime Deployment Coordinator
  └── stats_monitor.go        # Live CPU/RAM container statistics polling engine
  ```

  **How Railpack, Nixpacks & Buildpacks work in `railpack_builder.go`:**
  - [x] **Auto-Detection:** `builder.go` inspects the cloned repo root for `package.json`, `requirements.txt`, `go.mod`, `Cargo.toml`, `composer.json`, `Gemfile`, etc.
  - [x] **Build Plan:** `railpack_builder.go` invokes the Railpack/Nixpacks/Buildpacks engine to determine system packages, build phase commands, and start command.
  - [x] **OCI Image & Log Streaming:** Compiles an optimized multi-stage Docker image (`vessel-app-[project-id]:latest`) on the local Docker socket, streaming every build log line over WebSocket (`gorilla/websocket`) to the dashboard's `@xterm/xterm` live terminal.

  **How `deployer.go`, `container_manager.go` & `stats_monitor.go` complete the deployment:**
  - [x] **Secret Injection:** Decrypts `.env` secrets from `internal/store` and injects them as container environment variables (`deployer.go`).
  - [x] **Container Launch:** Launches the container with resource constraints (`CPURequest`, `MemoryLimitMB`) and port mappings (`container_manager.go`).
  - [x] **Health Verification:** Hits the `HealthCheckPath`; on success, instructs `internal/proxy` to update the Caddyfile and execute `caddy reload` — completing a zero-downtime deployment.
  - [x] Implement live CPU/RAM stats polling generator for active containers (`stats_monitor.go`).

- [x] **Caddy v2 Dynamic Proxy Manager (`internal/proxy/`)**:
  - [x] Auto-generate `/data/caddy/Caddyfile` rules when containers are deployed or custom domains are attached (`caddyfile_generator.go`).
  - [x] Execute `caddy reload` cleanly inside Docker network when configurations change (`proxy_manager.go`).
- [x] **Embedded SQLite Store & `.env` Vault (`internal/store/`)**:
  - [x] Initialize `CGO_ENABLED=0` modernc.org/sqlite database instance in `data/vessel.db` (`store.go`).
  - [x] Create modular repositories following strict `snake_case` and one-component-per-file (`project_store.go`, `domain_store.go`, `user_store.go`, `invite_store.go`, `env_var_store.go`).
  - [x] Implement AES-256-GCM encryption/decryption for `.env` secrets at rest (`vault.go`).
  - [x] Verify complete CRUD capabilities and schema initialization via unit tests (`store_test.go`).
- [x] **REST & WebSocket API Handlers (`internal/api/`)**:
  - [x] `GET /api/projects` / `POST /api/projects` (with auto `sslip.io` wildcard domain generation) / `DELETE /api/projects/:id` (`project_handler.go`)
  - [x] `POST /api/projects/:id/deploy` (Triggers multi-language OCI build & zero-downtime container rollout)
  - [x] `GET /api/projects/:id/domains` / `POST /api/projects/:id/domains` / `DELETE /api/domains/:id` (`domain_handler.go`)
  - [x] `GET /api/projects/:id/env` / `PUT /api/projects/:id/env` (Encrypted `.env` vault read/write via `env_handler.go`)
  - [x] `GET /ws/terminal/:id` (Interactive xterm shell via Docker `exec -it bash` over WebSocket using `gorilla/websocket`)
- [ ] **Enterprise PaaS Core Services & Advanced API (`internal/api/` & `internal/engine/`)**:
  - [x] **Auth API & RBAC Guards (`auth_handler.go`)**: `POST /api/auth/register`, `POST /api/auth/login` (JWT token issuance), `GET /api/auth/me`, `POST /api/auth/logout`, and route-level auth middleware (`RequireAuth`, `RequireRole`).
  - [x] **One-Click Managed Databases & S3 (`database_handler.go`, `database_deployer.go`, `storage_handler.go`, `storage_deployer.go`)**: Instant provisioning of `PostgreSQL`, `MySQL`/`MariaDB`, `Redis`, `MongoDB`, and `MinIO (S3 Object Storage)` with persistent volumes, encrypted credentials (`vault.go`), and internal DNS strings (`GET/POST /api/databases`, `POST /api/databases/:id/start|stop`, `GET/POST /api/storage`, `POST /api/storage/:id/start|stop`).
  - [x] **Live Database Data Browser & Web SQL Console (`/api/databases/:id/query`)**: Execute SQL queries (`SELECT`, `INSERT`, `UPDATE`, `EXPLAIN`) against managed PostgreSQL/MySQL instances or inspect/run commands against Redis (`KEYS`, `GET`, `SET`) directly over REST/WebSocket without external database GUI clients.
  - [x] **Git Providers Authentication, Auto-Clone & Webhooks (`git_handler.go`, `git_service.go`, `git_store.go`)**: Authenticate users with GitHub & GitLab (`POST /api/git/connect`, `GET /api/git/repos`), support public git URLs or private OAuth/PAT repos, auto-clone `repositoryUrl` (`git pull/clone`) before build in `handleDeployProject`, and handle push webhook triggers (`POST /api/webhooks/git/{projectId}`).
  - [x] **Pull Request (PR) Preview Environments (`pr_preview_deployer.go`)**: Handle `pull_request.opened` and `synchronize` webhook events by launching isolated, ephemeral container environments (`pr-123.my-project.sslip.io`) with cloned variables, post commit status check updates (`Deploying -> Deployed`) to GitHub/GitLab PR check APIs, and automatically destroy containers when PR is closed/merged.
  - [x] **True Railway Canvas Architecture: Multi-App & Multi-Environment Support (`environment_handler.go`, `environment_store.go`, `app_service_handler.go`, `app_service_store.go`)**: Transform `ProjectConfig` into an overarching Railway/Coolify workspace canvas with independent environments (`production`, `staging`, `preview`). Enable multiple Git application containers (`AppServiceConfig`) plus managed databases and S3 storage buckets per environment (`GET/POST /api/projects/:id/environments`, `GET/POST /api/projects/:id/environments/:envId/services`).
  - [x] **Service-Level Configuration & Deployment Management (`deployment_store.go`, `service_var_store.go`, `deployment_handler.go`, `service_var_handler.go`)**: Implement full Service Settings (`rootDirectory`, `buildCommand`, `startCommand`, `dockerfilePath`, `internalPort`, `domain`, `replicas`, `restartPolicy`, `teardownTimeout`, `healthCheckPath`, `serverless/autoSleep`), Service Deployments Tab (`history`, `logs`, `metrics`), and Service Variables Tab (`raw/editor`, variable references `${{ Postgres.DATABASE_URL }}`).
  - [x] **Multi-Engine Build Strategy & Layer Caching (`railpack_builder.go`)**: Allow per-service build engine selection (`dockerfile`, `nixpacks`, `buildpacks`, or `railpack`), supporting BuildKit layer caching (`--cache-from` / `--cache-to`) for `<5 second` incremental deployments.
  - [x] **Health-Checked Zero-Downtime Swaps & Auto-Abort (`deployer.go`)**: Enhance container rollout engine with `HealthCheckPath` active polling and configurable `WarmupSeconds`. If the green container fails health checks within `teardownTimeout`, automatically abort the swap and retain 100% of production traffic on the blue container.
  - [x] **Project Settings Architecture (`project_settings_handler.go`, `webhook_store.go`, `token_store.go`, `member_store.go`)**: Implement complete Project Settings (`General`, `Usage/Billing` moved to commercial side, `Environments`, `Shared Variables`, `Webhooks` for deployment/alert triggers, `Members & Workspace roles`, `API Tokens`, and `Danger zone`).
  - [x] **Scoped API Keys & Granular Permissions (`token_handler.go`, `token_store.go`)**: Issue API tokens with granular RBAC scopes (`deploy:write`, `logs:read`, `env:read`, `database:manage`), IP allowlists, and expiration dates for secure CI/CD and Terraform integration.
  - [x] **Workspaces Architecture & Trusted Domains (`workspace_handler.go`, `workspace_store.go`)**: Implement Workspaces where projects belong (`General People/Teams` separated from `Project Members`), Trusted Domains (`/api/workspaces/:id/domains`), and SSH Keys (`/api/workspaces/:id/ssh-keys`).
  - [x] **Automated DB & Volume Backups (`backup_handler.go`, `backup_manager.go`)**: Automated backup scheduling, `pg_dump`/`mysqldump`/`sqlite3 .dump` execution, and automated S3/MinIO offsite uploads (`GET/POST /api/backups`, `POST /api/backups/trigger`, `GET/POST /api/s3-destinations`).
  - [x] **Teams, Organizations & Invitations (`team_handler.go`, `team_store.go`)**: Multi-tenant collaboration with `Owner`, `Admin`, and `Member` roles (`GET/POST /api/teams`, `POST /api/teams/:id/invite`, `DELETE /api/teams/:id/members/:userId`).
  - [x] **Server Settings & Profile (`settings_handler.go`)**: System configurations (Docker system prune, Caddy wildcard IP, Notification Webhooks: Discord/Slack/Telegram/Email) and Personal Access Tokens (PATs) for CLI.
  - [x] **Wildcard Root Domain & Automatic Subdomain Provisioning (`domain_handler.go`, `caddyfile_generator.go`)**: Support a global `DefaultWildcardDomain` (e.g. `apps.yourdomain.com`). When a new service or project is launched, automatically provision `my-app.apps.yourdomain.com` with Caddy v2 Let's Encrypt wildcard certificates alongside `sslip.io` fallback.
  - [x] **Notification Integrations (`internal/notifier/`, `internal/api/`)**:
  - **Email (SMTP / Resend)**: SMTP host, port, user, password; or Resend API key. Send templated emails for team/project member invitations (with copy-link fallback) and deployment success/failure alerts. Commercial cloud will use AWS SES instead.
  - **Slack**: Enabled toggle, webhook URL, send test notification.
  - **Discord**: Enabled toggle, webhook URL, ping enable/disable, send test notification.
  - **Telegram**: Enabled toggle, bot API token, chat ID, send test notification.
  - **Pushover**: Enabled toggle, user key, API token, send test notification.
  - **Generic Webhook**: Enabled toggle, POST webhook URL, send test notification.
  - **Notification Preferences**: Per-project toggle to enable/disable each channel independently for deployment events.
- [x] **OAuth 2.0 Authentication Providers & 2FA (`internal/api/auth_handler.go`, `internal/services/oauth/`)**:
  - [x] Admin-configurable OAuth providers with enable/disable toggle, client ID, client secret, redirect URI, and provider-specific fields (base URL, tenant).
  - [x] Supported providers: GitHub, GitLab, Google, Azure AD, Discord, Authentik, Bitbucket, Clerk, Infomaniak, Zitadel.
  - [x] Member login via OAuth — users authenticate with their chosen provider to access the dashboard.
  - [x] Two-factor authentication (2FA/TOTP) for local accounts with recovery codes and optional enforced policy per workspace.
- [x] **Advanced Server Settings (`internal/api/settings_handler.go`)**:
  - [x] **Registration Control**: Toggle to enable/disable new user signups.
  - [x] **DNS Configuration**: Custom DNS resolver addresses for container networking (e.g. `1.1.1.1`), DNS validation toggle.
  - [x] **API Access Control**: Restrict API access to specific IPs/CIDR ranges (e.g. `192.168.1.100`, `10.0.0.0/8`). Empty = all IPs allowed.
  - [x] **MCP Server Toggle**: Enable/disable the MCP server endpoint for AI agent integrations (`/api/mcp` JSON-RPC).
- [x] **Agent Mode (`cmd/vesseld --agent`)**:
  - [x] Add `--agent --token=<auth_token>` and `--server=<wss://...>` flags to `vesseld`.
  - [x] Implement secure outbound WebSocket / mTLS tunnel for remote control.
  - [x] Allow remote execution of Docker commands over the tunnel without exposing public ports.
  - [x] Ships in the open-source daemon; the cloud-side connection acceptor is in Phase 6.
- [x] **Update Management (`internal/api/settings_handler.go`, `internal/updater/`)**:
  - [x] **Update Check Frequency**: Configurable cron expression for automatic update checks (e.g. `0 * * * *`).
  - [x] **Manual Check Button**: Trigger an immediate update check from the API (`POST /api/settings/updates/check`).
  - [x] **Auto-Update & Deploy**: Toggle to enable/disable automatic updates and execute binary rollout (`POST /api/settings/updates/deploy`).

---

## 💻 Phase 3: Control Panel Dashboard (`dashboard/`)

- [ ] **Navigation & Shell Layout**:
  - Responsive dark-mode glassmorphism sidebar (`Dashboard`, `Projects`, `Databases & Storage`, `Jobs & Backups`, `Teams`, `Settings`).
  - System health indicator header (`CPU %`, `RAM %`, `Docker Status`, `Upgrade Available banner`).
- [ ] **Project Management & Deployment Pages**:
  - "New Project" Wizard (`Connect GitHub/GitLab OAuth/PAT`, `Select from authenticated public/private repositories` OR `paste public Git URL` -> `Select Branch` -> `Configure Build Port`).
  - Project Details View with Tab Navigation (`Overview`, `Live Logs`, `Environment Variables`, `Settings`).
  - **Pull Request (PR) Previews Tab (`/projects/:id/previews`)**: Monitor active ephemeral PR preview environments, inspect Git commit links, view isolated logs, and trigger 1-click manual teardown.
  - **Service Build Strategy Configurator**: Dropdown in Service Settings to select build engine (`Dockerfile`, `Nixpacks`, `Buildpacks`, or `Railpack`) and toggle BuildKit layer caching (`--cache-from`).
- [ ] **Managed Databases & S3 Storage Pages**:
  - One-click spin-up modal (`PostgreSQL 16`, `MySQL 8`, `Redis 7`, `MongoDB 7`, `MinIO S3`).
  - Database Details Dashboard (`ConnectionString copy`, `One-Click Terminal console`, `Resource metrics`).
  - **Live Web SQL Console & Redis Table Explorer Tab (`/databases/:id/query`)**: Interactive query editor with syntax highlighting, visual schema table inspector, and instant data grid execution (`SELECT`, `INSERT`, `UPDATE` or Redis key browser).
- [ ] **Scheduled Jobs & Automated Backups UI**:
  - Cron Job Manager with visual cron builder and execution history table.
  - Backup Schedule Configurator with S3 destination dropdown and 1-click restore/download.
- [ ] **Teams & Organization Management Pages**:
  - Team member invitation modal (`Email` + `Role selection`) and project assignment table.
- [ ] **Live Interactive Terminal (`@xterm/xterm`)**:
  - Mount `@xterm/xterm` inside custom React component with dark theme & automatic resizing (`@xterm/addon-fit`).
  - Connect directly to `ws://host:8080/ws/terminal/:id` for live container bash access and live `docker build` streams.
- [ ] **`.env` Secret Vault Editor & Profile Settings**:
  - Multi-line secure `.env` key-value editor with instant encryption and 1-click rolling container restart.
  - User profile settings (`Update Name/Email/Password`, `Manage CLI Personal Access Tokens`).
  - **Scoped API Keys Manager (`/settings/tokens`)**: Modal to create API keys with granular RBAC checkboxes (`deploy:write`, `logs:read`, `env:read`, `db:manage`), IP allowlist restriction, and expiry dates.
- [ ] **Server Settings & Domain Management**:
  - **Global Wildcard Domain Configurator**: Input box in Server Settings for `DefaultWildcardDomain` (e.g. `cloud.yourdomain.com`) to enable instant Caddy v2 Let's Encrypt wildcard subdomains for all apps.
  - **OAuth Provider Manager (`/settings/auth`)**: Table of OAuth providers with enable/disable toggles, client ID/secret fields, redirect URI display, and provider-specific fields (base URL, tenant).
  - **2FA Setup Page**: TOTP QR code scanner, verification code input, recovery codes display.
  - **Notification Integrations Page (`/settings/notifications`)**: Per-provider configuration cards with enable toggle and fields — SMTP (host, port, user, password), Resend (API key), Slack (webhook URL), Discord (webhook URL, ping toggle), Telegram (bot token, chat ID), Pushover (user key, API token), Generic Webhook (POST URL). Each card has a "Send Test" button.
  - **Advanced Settings Page (`/settings/advanced`)**: Registration allowed toggle, custom DNS resolvers, API IP allowlist (CIDR input), MCP server toggle.
  - **Update Settings Page (`/settings/updates`)**: Auto-update toggle, update check frequency cron input, manual "Check for Updates" button with status display.

---

## 🌐 Phase 4: Public Marketing Site (`web/`)

- [ ] **Hero Section & Quick-Install Banner**:
  - High-conversion hero banner with one-click copyable install command: `curl -fsSL https://get.vessel.dev | sh`.
  - Interactive terminal mockup tabs showing instant container rollouts & CPU/RAM usage.
- [ ] **Comparison Tables vs. Existing Solutions**:
  - Comparison grid highlighting `<30MB RAM` Go daemon vs. Coolify/CapRover/Dokku/Railway/Vercel.
- [ ] **Documentation & FAQ pages**:
  - Setup guides for DigitalOcean, Hetzner, AWS EC2, and local bare-metal servers.

---

## 📦 Phase 5: Production Docker & Script Verification

- [ ] Build unified multi-stage `Dockerfile` packaging `dashboard/dist/` inside the `vesseld` Go binary.
- [ ] Test end-to-end `deploy/install.sh` and `scripts/upgrade.sh` inside isolated test containers.

---

## ☁️ Phase 6: Commercial SaaS — Vessel Cloud (`cloud.vessel.dev` & BYOS Agent Mode)

> The cloud control plane lives in a separate private repository (`vessel-cloud`). The open-source repo (`vessel`) contains the self-hosted daemon, dashboard, marketing site, and agent — the cloud backend is proprietary. The cloud reuses the open-source dashboard (users connect to cloud-hosted instances with the same UI) and adds a lightweight staff admin panel.

```text
vessel-cloud/
├── cmd/apid/                  # Cloud API server entrypoint
├── internal/
│   ├── api/                   # REST handlers (teams, billing, instances)
│   ├── store/                 # PostgreSQL repositories
│   ├── agent/                 # Agent tunnel manager (accepts WebSocket)
│   ├── billing/               # Stripe, invoicing, metering
│   ├── mailer/                # SES transactional email
│   ├── auth/                  # OAuth, SSO/SAML, JWT
│   ├── featureflags/          # GrowthBook evaluation
│   ├── audit/                 # Immutable audit log
│   └── types/                 # Domain structs
├── migrations/                # PostgreSQL schema migrations
├── dashboard/                 # Staff admin dashboard (React)
├── scripts/                   # DB migrations, seed, deploy
├── Dockerfile
├── docker-compose.yml         # API + Postgres + Redis
└── Makefile
```

- [ ] **Private Cloud Repository (`vessel-cloud`)**:
  - Initialize private repo with Go API server, multi-tenant PostgreSQL database.
  - Staff admin dashboard for managing customers, instances, billing.
  - Stripe integration for checkout, billing, invoicing, and seat management.
  - GrowthBook feature flag evaluation for per-tier feature gating.
- [ ] **Go Daemon Agent Mode (`cmd/vesseld --agent`)**:
  - The `--agent` flag ships in the open-source daemon (see Phase 2). Phase 6 adds the cloud-side connection acceptor.
- [ ] **Vessel Cloud Control Plane**:
  - Agent connection acceptor (`internal/agent/`) — accepts inbound WebSocket tunnels from `vesseld --agent` instances and routes commands.
  - "Connect Server" 1-Click Wizard generating unique install tokens: `curl -fsSL https://get.vessel.dev/agent | sh -s -- --token=vsl_live_xyz`.
  - Multi-server fleet deployment dashboard allowing 1-click deployments to multiple geographic VPS regions.
- [ ] **Billing & Subscription Integration**:
  - Integrate Stripe checkout and subscription management (`Hobby / Pro / Team` tiers).
  - Implement automated BYOS seat limits and deployment rate limiting.
  - Usage metering — track deployments, container hours, bandwidth per account for billing.
  - Integrate GrowthBook for feature flag gating of Pro/Team features and gradual rollouts.
- [ ] **AI Agent Protocol (MCP) & API Ecosystem**:
  - Expose Vessel's REST API as an MCP server (`@modelcontextprotocol/sdk`) so AI agents (Claude Code, Cursor, etc.) can deploy apps, manage databases, and query logs programmatically.
  - Publish an official Vessel API client SDK for Node.js and Go.
- [ ] **Enterprise Features**:
  - **SSO / SAML**: Enterprise single sign-on with Okta, Azure AD, Google Workspace.
  - **Audit Logs**: Immutable event log of all actions (deployments, config changes, member management) with export API.
  - **Email Delivery**: Use AWS SES for all transactional emails (invites, alerts) — no SMTP config needed on cloud.
  - **Custom Branding**: White-label dashboard with custom domain, logo, and colors for enterprise tenants.
- [ ] **Self-Hosted License & Telemetry**:
  - Optional telemetry ping (version, anonymous usage stats) for upgrade notifications and feature analytics.
  - License key system for self-hosted enterprise tier with offline activation.
