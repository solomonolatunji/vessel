# 🛰️ Vessel Development Roadmap & Tasks (`TODO.md`)

---

## 💻 Phase 3: Control Panel Dashboard (`dashboard/`)

- [ ] **GitHub One-Click Connect (UI)**: Add a flow in the dashboard to trigger the GitHub App Manifest creation.
- [ ] **Vercel Project Imports (UI)**: One-click UI flow to authenticate and select existing projects from Vercel.
- [ ] **AI-Powered Deployment Diagnostics (UI)**: Explain deployment failures in plain English within the deployment logs tab.
- [ ] **Serverless Functions Editor**: Built-in GUI editor for serverless functions with an embedded AI assistant.
- [ ] **Shared Confirmation Dialogs**: Robust shared components to prevent accidental deletions of services, databases, domains, and env vars.
- [ ] **Guard Active Deployments**: UI logic to disable system updates or conflicting actions while an active deployment is running.

- [ ] **Navigation & Shell Layout**:
  - Responsive dark-mode glassmorphism sidebar
    (`Dashboard`, `Projects`, `Databases & Storage`, `Jobs & Backups`, `Teams`, `Settings`).
  - System health indicator header (`CPU %`, `RAM %`, `Docker Status`, `Upgrade Available banner`).
- [ ] **Project Management & Deployment Pages**:
  - "New Project" Wizard
    (`Connect GitHub/GitLab OAuth/PAT`, `Select from authenticated public/private repositories` OR `paste public Git URL` -> `Select Branch` -> `Configure Build Port`).
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
- [ ] **Server Settings & Domain Management (BYOK Architecture)**:
  - **Global Wildcard Domain Configurator**: Input box in Server Settings for `DefaultWildcardDomain` (e.g. `cloud.yourdomain.com`) to enable instant Caddy v2 Let's Encrypt wildcard subdomains for all apps.
  - **OAuth Provider Manager (`/settings/auth`)**: Table of OAuth providers with enable/disable toggles, client ID/secret fields, redirect URI display, and provider-specific fields (base URL, tenant).
  - **2FA Setup Page**: TOTP QR code scanner, verification code input, recovery codes display.
  - **Notification Integrations Page (`/settings/notifications`)**: Per-provider configuration cards (Bring Your Own Key) for App-level notifications — SMTP, Resend, Slack, Discord, Telegram, Pushover, Generic Webhooks. Protects platform email reputation by keeping app-level outbound emails under user-controlled API keys.
  - **AI Settings Page (`/settings/ai`)**: "Bring Your Own Key" configuration for OpenAI/Anthropic to unlock unlimited AI Deployment Diagnostics and MCP features (prevents runaway API costs while empowering users).
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

## 🤖 Phase 7: AI Agent Protocol (MCP) & API Ecosystem (OSS & Cloud)

> The MCP server and API Ecosystem is a core feature built into the `vesseld` Go daemon directly so self-hosters can use it for free, but it is also securely exposed and proxied via the Vessel Cloud control plane for managed users.

- [x] **REST API to MCP Bridge**:
  - Expose Vessel's REST API as an MCP server (`@modelcontextprotocol/sdk`) so AI agents (Claude Code, Cursor, etc.) can deploy apps, manage databases, and query logs programmatically.
  - Implement Local stdio transport for the CLI daemon and SSE/WebSocket transport for the Cloud.
- [ ] **SDKs**:
  - Publish an official Vessel API client SDK for Node.js and Go.
