# 🛰️ Vessel: The Ultra-Sleek, Lightweight Self-Hosted PaaS

> **Tagline**: _Turn any bare-metal VPS into your own private Vercel & Railway in 60 seconds._  
> **Project Name**: **Vessel** (`vessel.dev` / `github.com/vessel-run/vessel`)  
> **Mission**: Build an open-source, blazing-fast, developer-first self-hosted PaaS with a next-generation GUI, automated Docker container builds, zero-downtime deployments, Caddy edge SSL routing, and one-click self-updating/upgrading.

---

## 1. System Anatomy & Terminology

To ensure crystal clarity across the repository, Vessel separates the **Public Landing Page** from the **Self-Hosted Control Panel (GUI)** and **Orchestration Server**:

1. **`website/` (Public Marketing Landing Page)**:
   - **What it is**: The public website hosted at `vessel.dev` (built with Astro / Vite / Tailwind).
   - **Who sees it**: Developers globally discovering the project, reading features, documentation, and copying the one-line install command (`curl -fsSL https://get.vessel.dev | sh`).
2. **`src/client/` (Self-Hosted Web GUI Dashboard)**:
   - **What it is**: The interactive management dashboard built with **Vite + TanStack Router + React + Tailwind + `@xterm/xterm`**.
   - **Who sees it**: Anyone who installs Vessel on their VPS. When they visit `https://their-server-ip:3000` or `https://app.their-domain.com`, this is the GUI control panel they log into to deploy apps, manage databases, view live logs, and configure `.env` variables.
3. **`src/server/` (Orchestrator Backend & Daemon — `vesseld`)**:
   - **What it is**: High-performance backend orchestrator built in **Go (Golang)**.
   - **What it does**: Talks directly to the Docker socket (`/var/run/docker.sock`), manages Caddy SSL/reverse-proxy rules, executes git webhooks, streams WebSocket terminal logs (`gorilla/websocket`), manages SQLite state, and handles self-upgrade/update commands (`scripts/upgrade.sh`).
4. **`get-vessel/` (Installation & Script Delivery Engine)**:
   - **What it is**: Lightweight server / script host (`get.vessel.dev`) serving `install.sh`, `upgrade.sh`, and system bootstrap scripts.
5. **`scripts/` (Self-Update, Upgrade & Diagnostics Scripts)**:
   - **What it is**: Core shell automation allowing the user (and the GUI dashboard via "Check for Updates" / "Upgrade Now" button) to self-update Vessel in place without losing containers or data.

---

## 2. Comprehensive Repository Structure

Modeled after modern, clean open-source PaaS repositories (`aeroplane`), Vessel uses the following unified monorepo structure:

```
vessel/
├── src/
│   ├── client/                # 💻 SELF-HOSTED WEB GUI (What the VPS installer sees & uses)
│   │   ├── src/
│   │   │   ├── routes/        # TanStack Router type-safe routes (Dashboard, Projects, Databases, Settings)
│   │   │   ├── components/    # Glassmorphism UI cards, xterm.js live terminal, `.env` vault editor
│   │   │   ├── hooks/         # WebSocket streaming hooks, TanStack Query mutations
│   │   │   └── lib/           # API client, auth token manager, theme toggler
│   │   ├── index.html         # Main entry point for local GUI
│   │   ├── package.json       # Frontend dependencies (@tanstack/react-router, @xterm/xterm, lucide-react)
│   │   └── vite.config.ts     # Vite build configuration
│   │
│   ├── server/                # ⚙️ ORCHESTRATOR API & DAEMON (`vesseld` in Go)
│   │   ├── cmd/
│   │   │   └── vesseld/       # Main entrypoint (`main.go`)
│   │   ├── internal/
│   │   │   ├── api/           # REST & WebSocket API endpoints for `src/client/` GUI
│   │   │   ├── docker/        # Native Docker SDK orchestration, volume & container management
│   │   │   ├── proxy/         # Caddy v2 Caddyfile generator & dynamic reload (`caddy reload`)
│   │   │   ├── store/         # Embedded SQLite (`CGO_ENABLED=0` sqlite) + AES-256 `.env` vault
│   │   │   ├── updater/       # Self-upgrade & update manager (triggers `scripts/upgrade.sh`)
│   │   │   └── ssh/           # Remote node terminal execution via gRPC/SSH
│   │   ├── go.mod             # Go module definition
│   │   └── go.sum             # Go module checksums
│   │
│   └── shared/                # 🔗 SHARED TYPES, DTOs & CONSTANTS
│       ├── types.ts           # Shared TypeScript interfaces (ContainerHealth, ProjectConfig, DeployEvent)
│       └── constants.ts       # API routes, WebSocket event names (`deploy:log`, `stats:cpu`)
│
├── website/                   # 🌐 PUBLIC MARKETING LANDING PAGE (vessel.dev)
│   ├── src/
│   │   ├── pages/             # Landing page, Features overview, Documentation, FAQ
│   │   └── components/        # Hero section, one-click copy install banner, interactive demo screenshots
│   ├── public/                # Logos, OG graphics, favicons
│   ├── package.json           # Astro / Vite landing page dependencies
│   └── astro.config.mjs       # Astro site configuration
│
├── get-vessel/                # 📦 INSTALLATION & UPGRADE SCRIPT HOST (get.vessel.dev)
│   ├── install.sh             # One-click curl installation (`curl -fsSL https://get.vessel.dev | sh`)
│   ├── upgrade.sh             # Automated upgrade script to pull latest Vessel release safely
│   ├── server.js              # Lightweight delivery server with telemetry & version checking
│   └── Dockerfile             # Container definition for the installation host
│
├── scripts/                   # 🛠️ SYSTEM AUTOMATION, UPGRADE & BOOTSTRAP SCRIPTS
│   ├── upgrade.sh             # In-place self-upgrade script executed by GUI/backend during updates
│   ├── backup.sh              # SQLite state & Caddy configuration automated backup (`/data/backups`)
│   ├── restore.sh             # One-click disaster recovery restore script
│   ├── bootstrap-host.sh      # Linux OS check, Docker checking, systemd unit provisioning
│   └── generate-ssl.sh        # Fallback local self-signed SSL / Let's Encrypt renewal helper
│
├── data/                      # 💾 PERSISTENT DATA DIRECTORY (Mounted in Docker/Systemd)
│   ├── vessel.db              # SQLite primary database storing projects, users, and encrypted envs
│   ├── backups/               # Automated daily backups (`.tar.gz`)
│   └── caddy/                 # Caddy certificates and dynamic `Caddyfile`
│
├── docker-compose.yml         # All-in-one local dev & production container deployment
├── Dockerfile                 # Multi-stage production build uniting `src/client` and `src/server`
├── Makefile                   # Developer commands (`make dev`, `make build`, `make upgrade-test`)
└── README.md                  # Main open-source repository documentation
```

---

## 3. How Self-Updates & Upgrades Work (`Upgrade/Update Engine`)

One of Vessel's killer open-source features is **1-Click Self-Updating**:

1. **Check for Updates**: The GUI (`src/client`) pings `GET /api/system/version` on the orchestrator (`src/server`). The orchestrator checks GitHub Releases (`github.com/vessel-run/vessel/releases/latest`) or `https://get.vessel.dev/version`.
2. **One-Click Upgrade**: If a new version (`v1.2.0`) is found, the GUI displays a banner with an **"Upgrade Now"** button.
3. **Safe In-Place Execution**:
   - When clicked, `src/server` triggers `scripts/upgrade.sh` in the background.
   - `upgrade.sh` automatically creates a snapshot (`backup.sh -> /data/backups/vessel-before-upgrade.db`).
   - Pulls the latest Docker image (`docker pull ghcr.io/vessel-run/vessel:latest` or downloads precompiled Go binary + built `web/dist`).
   - Gracefully restarts `vesseld` systemd service or Docker container (`docker-compose up -d --force-recreate vessel`).
   - Running user app containers (`Vite`, `NestJS`, `Postgres`, `Redis`) are **completely untouched and experience ZERO downtime during the Vessel upgrade**.

---

## 4. Why This Structure Wins

- **Separation of Concerns**: The public marketing site (`website/`) never bloats the local VPS deployment. When `install.sh` runs on a VPS, it only downloads and mounts `src/server` (`vesseld`) and `src/client` (`vessel-ui dist`).
- **Clean Shared Contracts**: `src/shared/types.ts` ensures that the frontend GUI (`Vite + TanStack Router`) and backend orchestrator speak the exact same data structures.
- **Easy Open-Source Contributions**: Any developer from GitHub can fork the repo, run `make dev`, and instantly work on the GUI (`src/client`), Go backend (`src/server`), or landing page (`website`) independently.
