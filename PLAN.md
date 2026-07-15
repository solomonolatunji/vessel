# Vessl Dashboard v2: Complete Feature & Architectural Plan

This document serves as the master specification for building out **Vessl Dashboard v2** (`dashboard/`). It aligns the React 19 / TanStack Router frontend with all **18+ Enterprise & Parity capabilities** recently implemented in the Go (`vessld`) backend.

---

## 1. Tech Stack & Architectural Conventions

- **Framework:** React 19 + Vite (`npm run dev`)
- **Routing:** TanStack Router (`src/routes/`, file-based, auto-generating `routeTree.gen.ts`)
- **State & Querying:** TanStack Query v5 + TanStack Store + Zod validation
- **Styling:** Tailwind CSS v4 (`@theme`), Radix UI Primitives, `tailwind-merge` + `clsx` + `class-variance-authority`
- **Icons:** Lucide React (`lucide-react`)
- **Editors & Terminal:** Monaco Editor (`@monaco-editor/react`), XTerm.js (`xterm`, `xterm-addon-fit`), React Flow (`@xyflow/react`)
- **Code Rules (`AGENTS.md`):** Max 350 lines per file, one component per file, named exports, `kebab-case` filenames (`project-card.tsx`), formatted strictly with Biome (`npm run format:fix`).

---

## 2. Global Shell & Navigation Architecture

The dashboard uses a **Contextual Shell** (`_shell.tsx`) with a global **Topbar** and a context-aware **Sidebar**.

### 2.1 Topbar Navigation

- **Contextual Breadcrumb/Title:** Dynamically updates based on global vs. project scope (e.g., "Dashboard", "Canvas", "Deployments").
- **Global Command Menu (`Cmd+K`):** Spotlight search across Projects, Environments, App Services, Databases, Domains, and System Settings.
- **Active Deployment / System Health Indicator:** Real-time badge showing CPU/Memory pressure, active background builds (`SSE /ws/events`), and available server updates.
- **Theme & Profile:** Dark/Light/System theme toggle and user logout (`POST /auth/logout`).

### 2.2 Contextual Sidebar States

#### A. Global Instance Context (`/projects`, `/databases`, `/settings/*`)

- **Overview:** Aggregated stats (running containers, total CPU/RAM usage, active jobs).
- **Projects (`/projects`):** Grid/List of all projects with quick-status badges.
- **Databases (`/databases`):** Global view of all relational, NoSQL, and broker instances across all projects.
- **Marketplace (`/marketplace`):** One-click templates and object storage integrations.
- **Imports:** Links to Railway/Vercel Importers.
- **System Settings (`/settings/*` - Admin Only):**
  - **DNS Providers:** Cloudflare, Namecheap, Spaceship API integrations (`/settings/dns`).
  - **Maintenance & Cleanup:** Garbage collection (`docker system prune`), disk usage gauges (`/settings/maintenance`).
  - **Updates:** Control plane version checks and auto-update configuration (`/settings/updates`).
  - **Migration Bundles:** Server-wide `.vessl` export and import (`/settings/migration`).
  - **Instance Users:** Global user management (`/settings/users`).

#### B. Project & Environment Context (`/projects/$projectId/*`)

- **Project Overview (`/projects/$projectId`):** Readme, active environment summary, quick metrics.
- **Interactive Canvas (`/projects/$projectId/canvas`):** Railway-inspired React Flow node graph visualizing App Services, Databases, and Storage volumes with drag-and-drop linking.
- **Services (`/projects/$projectId/services/$serviceId/*`):** Microservices and background workers inside the project.
- **Databases (`/projects/$projectId/databases/$dbId/*`):** Co-located database management, Data Browser, and SQL Studio.
- **Project Settings (`/projects/$projectId/settings`):** Project-level environment variables, custom domains, webhooks, and team RBAC.

---

## 3. Mapping Backend Features to Dashboard UI

Every single feature from the **Aeroplane vs. Vessl Feature Gap Analysis** maps directly to a dedicated UI experience:

|   #    | Feature                            | Target Route / Component                                                                      | UI Specification & User Flow                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| :----: | :--------------------------------- | :-------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **1**  | **Browser Onboarding Wizard**      | `/routes/onboarding.tsx`<br>`src/features/onboarding/*`                                       | **First-Run Interception:** If `GET /system/onboarding/status` indicates no owner account exists, redirect all traffic to `/onboarding`.<br>**Wizard Steps:**<br>1. **Owner Account:** Email and password setup.<br>2. **Control Plane Domain:** Configure `pilot.example.com` or fallback to `http://IP:8080`.<br>3. **GitHub App Setup:** Client ID, Secret, App ID, and Webhook (`/api/github/app/webhook`).<br>4. **Wildcard Root Domain:** Set `*.pilot.example.com` for auto-generated hostnames.<br>5. **Backup & R2 Storage:** Connect Cloudflare R2 bucket credentials.<br>6. **Bundle Import (Optional):** Dropzone to upload and decrypt `.vessl` server state.<br>_System Settings also includes a "Restart Onboarding" action._ |
| **2**  | **One-Line Installer**             | `/routes/getting-started.tsx`                                                                 | Show highlighted command (`curl -fsSL https://get.vessl.dev \| sh`) and system readiness checklist.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| **3**  | **Service Runtime Modes**          | `src/features/services/service-settings.tsx`                                                  | **Web vs. Worker Switcher:** Radio card selector when creating/editing an `AppService`.<br>- **Web:** Internal port input, public route generator, HTTP health checks (`/healthz`).<br>- **Background Worker:** No internal port, no public route, process uptime check badge (`runtimeMode === 'worker'`).                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| **4**  | **Static Site Deployments**        | `src/features/services/build-settings.tsx`                                                    | **Static Output Input:** Text field for `Static output directory` (e.g., `dist`, `build`, `.output/public`). When set, UI displays badge indicating the service runs inside an optimized `nginx:alpine` wrapper on internal port 80.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| **5**  | **Zero-Downtime Hot Swaps**        | `src/features/services/service-deployments.tsx`                                               | **Live Transition UI:** During `deploying` status, display both the active `running` container (`UUID-A`) and the probing `starting` container (`UUID-B`). Show real-time Traefik health check status before old container cleanup.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| **6**  | **Build Overrides**                | `src/features/services/build-settings.tsx`                                                    | **Command Override Inputs:** Expandable accordion under Railpack/Nixpacks settings:<br>- **Install Command:** (`--install-cmd`) e.g., `npm ci`<br>- **Build Command:** (`--build-cmd`) e.g., `npm run build`<br>- **Start Command:** (`--start-cmd`) e.g., `npm start`                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| **7**  | **Intelligent Env Var Linking**    | `src/features/services/service-variables.tsx`                                                 | **Smart Variable Drawer:** When editing service `.env` secrets, a side-drawer suggests auto-linked variables (`${postgres-db.POSTGRES_URL}`, `${timescaledb-db.TIMESCALE_URL}`). Includes `.env.example` parser pills to quickly autofill required keys.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| **8**  | **Database Data Imports**          | `src/features/databases/database-import-modal.tsx`                                            | **Import Data Modal (`POST /databases/:id/import`):**<br>- **URL Import:** Input for `postgres://` or `redis://` schemes with immediate syntax validation.<br>- **Railway Sync:** Auto-detects public DB URLs from imported Railway variables.<br>- **TimescaleDB Check:** Displays source/target extension compatibility warning badge.<br>- **History Table:** Live streaming progress of `pg_dump` / `redis-cli --rdb` jobs.                                                                                                                                                                                                                                                                                                              |
| **9**  | **Server Migration Bundles**       | `src/features/instance/migration-settings.tsx`<br>`/routes/_shell/settings/migration.tsx`     | **Bundle Manager (`/settings/migration`):**<br>- **Export Card:** Passphrase input + `Export Server Bundle` button (`GET /system/export`).<br>- **Import Card:** File upload dropzone (`.vessl`) + Passphrase + `Restore Server State` destructive action modal (`POST /system/import`).                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| **10** | **Database Provisioning Engines**  | `src/features/databases/create-database-modal.tsx`                                            | **Engine Selection Grid:** Organized categorizations with custom icon badges:<br>- **Relational:** PostgreSQL (`16-alpine`), **TimescaleDB (`latest`)**, MySQL (`8.0`), MariaDB (`11`), ClickHouse (`latest`)<br>- **NoSQL:** MongoDB (`7.0`), Redis (`7-alpine`), Dragonfly (`latest`), KeyDB (`latest`)<br>- **Brokers:** Kafka (`9092`), RabbitMQ (`5672`), NATS (`4222`)<br>- **One-Click:** NocoDB, Plausible, WordPress, Gitea                                                                                                                                                                                                                                                                                                         |
| **11** | **Data Browser & Row Editing**     | `src/features/databases/data-browser.tsx`<br>`/routes/_shell/databases/$dbId/data.tsx`        | **Relational Table Grid (`GET /databases/:id/data/:table`):**<br>- Table switcher (`GET /databases/:id/schemas`), filtering (`=`, `contains`, `>`), server-side pagination.<br>- **Inline Row Editing:** Double-click cells to edit or click `+ Add Row` (`POST /databases/:id/data/:table`).<br>**Redis Key Browser:** Specialized grid showing key names, types (`string`, `hash`, `list`, `set`), values, and interactive TTL editor.                                                                                                                                                                                                                                                                                                     |
| **12** | **Public Database Access & TLS**   | `src/features/databases/database-networking.tsx`                                              | **Public Hostname Controller (`PUT /databases/:id`):**<br>- Toggle for **Public Access (`ExternalDNS`)**.<br>- Displays generated TCP endpoint: `postgres-db.pilot.example.com:5432`.<br>- TLS Status badge (`Let's Encrypt TCP SNI enabled`) + one-click copy buttons for external clients.                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| **13** | **Logical Replication (CDC)**      | `src/features/databases/database-settings.tsx`                                                | **CDC Toggle Switch:** In Postgres & TimescaleDB configuration, toggle `Logical Replication (`wal*level=logical`)`. Displays warning badge: *"Enables max*replication_slots=10 and WAL retention for Change Data Capture tools."*                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| **14** | **Database Restore & Download**    | `src/features/databases/backup-manager.tsx`                                                   | **Backup Action Table (`/databases/:id/backups`):**<br>- **Download:** Button to stream `.sql` or `.rdb` directly from disk/R2 (`GET /backups/:id/download`).<br>- **Destructive Restore:** Red warning modal asking user to type database name before piping backup (`pg_restore --clean`, `mysql <`, `mongorestore`) into running container (`POST /backups/:id/restore`).                                                                                                                                                                                                                                                                                                                                                                 |
| **15** | **DNS Provider Automation**        | `src/features/instance/dns-settings.tsx`<br>`src/features/projects/project-domains.tsx`       | **Provider Credentials (`/settings/dns`):** Forms to save API keys for Cloudflare, Namecheap, and Spaceship.<br>**1-Click A-Record Sync:** On any service custom domain card (`/projects/:id/domains`), a `Sync A Record via DNS Provider` button that automatically writes the `1800` TTL record targeting the server IP.                                                                                                                                                                                                                                                                                                                                                                                                                   |
| **16** | **System Maintenance & Cleanup**   | `src/features/instance/maintenance-settings.tsx`<br>`/routes/_shell/settings/maintenance.tsx` | **Maintenance Dashboard (`/settings/maintenance`):**<br>- **Storage Gauges:** Root FS disk %, Docker storage reclaimable MBs, Backup volume size.<br>- **Garbage Collection:** `Run Docker Cleanup Now` button (`POST /system/maintenance/cleanup`) running `docker system prune -af --volumes`.<br>- **Cron Config:** Schedule selector for automated background pruning (`Docker Cleanup Cron`).                                                                                                                                                                                                                                                                                                                                           |
| **17** | **Railway Importer Specification** | `src/features/projects/railway-importer.tsx`<br>`/routes/_shell/import/railway.tsx`           | **Multi-Step Railway Import Wizard (`POST /import/railway`):**<br>1. **Token & Discovery:** Paste Railway Personal API Token (`Bearer <token>`), query GraphQL v2, select project.<br>2. **Service Classification Table:** Displays detected Git repos, Docker images, and database engines (Postgres, TimescaleDB, Redis, Mongo, ClickHouse).<br>3. **Configuration Checkboxes:** `Exclude RAILWAY_* variables` (default ON), `Recreate database engines` (creates local Vessl DBs), `Auto-deploy services`, and `Import database data` (runs automated `pg_dump`/`redis-cli` from public Railway URLs).                                                                                                                                    |
| **18** | **Control Plane Auto-Updates**     | `src/features/instance/update-settings.tsx`<br>`/routes/_shell/settings/updates.tsx`          | **Version Controller (`/settings/updates`):**<br>- **Version Card:** Shows `Current Version`, `Latest Version` (`GET /settings/updates/check`), and release notes.<br>- **Auto-Update Toggle:** Switch for `Auto Update Enabled` + `Update Check Cron` selector.<br>- **Manual Trigger:** `Deploy Update Now` button (`POST /settings/updates/deploy`) triggering `scripts/upgrade.sh` and graceful `vessld` container restart.                                                                                                                                                                                                                                                                                                              |

---

## 4. TanStack Router Route Tree Design

We will structure `src/routes/` with clear functional layout boundaries (`_shell`, `_auth`) and explicit URL paths:

```text
src/routes/
├── __root.tsx                               # Global QueryProvider, ThemeProvider, Toast, CommandMenu
├── _auth/                                   # Unauthenticated layout (centered card, no sidebar)
│   ├── login.tsx                            # GET /auth/login
│   └── register.tsx                         # GET /auth/register
│
├── onboarding.tsx                           # First-run browser onboarding wizard (/onboarding)
│
├── _shell/                                  # Authenticated layout (Topbar + Contextual Sidebar)
│   ├── index.tsx                            # Global Overview / Dashboard Home
│   ├── projects.tsx                         # Project List (`/projects`)
│   ├── databases.tsx                        # Global Database Inventory (`/databases`)
│   ├── marketplace.tsx                      # One-Click Apps & Storage Templates (`/marketplace`)
│   │
│   ├── import/                              # Migration Importers
│   │   ├── railway.tsx                      # Railway Project Importer (`/import/railway`)
│   │   └── vercel.tsx                       # Vercel Project Importer (`/import/vercel`)
│   │
│   ├── settings/                            # Super Admin & Instance Settings (`/settings/*`)
│   │   ├── index.tsx                        # General Instance Configuration
│   │   ├── dns.tsx                          # Cloudflare/Namecheap/Spaceship DNS Credentials
│   │   ├── maintenance.tsx                  # Garbage Collection & Disk Usage Alerts
│   │   ├── updates.tsx                      # Control Plane Version & Auto-Update Toggles
│   │   ├── migration.tsx                    # AES-256 `.vessl` Bundle Export/Import
│   │   └── users.tsx                        # Instance-Wide User Management
│   │
│   └── projects/                            # Project Context Routing (`/projects/$projectId/*`)
│       └── $projectId/
│           ├── index.tsx                    # Project Overview & Quick Stats
│           ├── canvas.tsx                   # Railway-Style React Flow Node Graph
│           ├── settings.tsx                 # Project RBAC, Webhooks, and Global Secrets
│           │
│           ├── services/
│           │   └── $serviceId/
│           │       ├── index.tsx            # Service Metrics & Overview
│           │       ├── deployments.tsx      # Build History, Logs & Rollback (`/deployments`)
│           │       ├── variables.tsx        # Secret Editor & Smart Env Var Linking
│           │       ├── domains.tsx          # Custom Domains, Wildcard & DNS Provider Sync
│           │       ├── build.tsx            # Railpack/Nixpacks Overrides & Static Output
│           │       ├── terminal.tsx         # XTerm.js Container Terminal (`/ws/terminal`)
│           │       └── serverless.tsx       # Monaco Editor for Edge Functions
│           │
│           └── databases/
│               └── $dbId/
│                   ├── index.tsx            # Connection Details, Credentials & Public Networking
│                   ├── data.tsx             # Table Browser, Row-Level Editor & Redis Explorer
│                   ├── query.tsx            # SQL Studio Playground (`/query`)
│                   └── backups.tsx          # Manual Snapshots, R2 Sync & Destructive Restore
```

---

## 5. Domain Feature Folder Organization (`src/features/`)

To keep files well below **350 lines** and strictly **one component per file**, every domain has its own modular folder:

```text
src/features/
├── auth/
│   ├── login-form.tsx
│   ├── register-form.tsx
│   └── use-auth.ts
│
├── onboarding/
│   ├── onboarding-wizard.tsx
│   ├── step-owner-account.tsx
│   ├── step-control-plane-domain.tsx
│   ├── step-github-app.tsx
│   ├── step-wildcard-domain.tsx
│   ├── step-backup-storage.tsx
│   └── step-bundle-import.tsx
│
├── projects/
│   ├── project-list.tsx
│   ├── project-card.tsx
│   ├── create-project-modal.tsx
│   ├── project-domains.tsx
│   ├── railway-importer.tsx
│   └── vercel-importer.tsx
│
├── canvas/
│   ├── environment-canvas.tsx
│   ├── app-service-node.tsx
│   ├── database-node.tsx
│   ├── storage-node.tsx
│   └── use-canvas-sync.ts
│
├── services/
│   ├── service-metrics.tsx
│   ├── service-deployments.tsx
│   ├── deployment-row.tsx
│   ├── live-logs-viewer.tsx
│   ├── service-variables.tsx
│   ├── smart-linker-drawer.tsx
│   ├── build-settings.tsx
│   ├── runtime-mode-card.tsx
│   └── web-terminal.tsx
│
├── databases/
│   ├── database-list.tsx
│   ├── create-database-modal.tsx
│   ├── database-connection-card.tsx
│   ├── database-networking.tsx
│   ├── data-browser.tsx
│   ├── table-data-grid.tsx
│   ├── row-editor-modal.tsx
│   ├── redis-key-browser.tsx
│   ├── sql-studio.tsx
│   ├── backup-manager.tsx
│   └── database-import-modal.tsx
│
└── instance/
    ├── dns-settings.tsx
    ├── maintenance-settings.tsx
    ├── update-settings.tsx
    └── migration-settings.tsx
```

---

## 6. API Client & SSE Interceptor Architecture (`src/lib/`)

### 6.1 Typed API Client (`src/lib/api-client.ts`)

- Uses `fetch` with automatic JSON parsing and Bearer token injection (`localStorage.getItem('vessl_token')`).
- Global 401 interceptor that automatically triggers `useAuthStore.getState().logout()` and redirects to `/login`.
- Standardized error extraction for clean toast notifications (`toast.error(err.message)`).

### 6.2 Real-Time Event Stream (`src/lib/use-event-stream.ts`)

- Connects to `/api/ws/events` (or `/api/events/sse`) via `EventSource`.
- Listens for backend deployment updates (`deployment:started`, `deployment:success`, `deployment:failed`), backup status changes, and database import progress.
- Automatically invalidates relevant TanStack Query queries (`queryClient.invalidateQueries({ queryKey: ['deployments', serviceId] })`) to update the UI instantly without manual polling.

---

## 7. Implementation Roadmap & Next Steps

1. **Phase 1: Router & API Client Verification**
   - Flesh out `src/lib/api-client.ts` and verify TanStack Query defaults (`staleTime`, retry behavior).
   - Register all expanded routes (`onboarding.tsx`, `settings/*`, `import/*`, `databases/*`) inside `src/routes/` and run `npm run generate-routes`.
2. **Phase 2: First-Run Onboarding Wizard (`/onboarding`)**
   - Build `src/features/onboarding/onboarding-wizard.tsx` and all step components to handle fresh server installations cleanly.
3. **Phase 3: Service Deep-Dive & Build Overrides**
   - Upgrade `AppService` components to support `RuntimeMode` (`web`/`worker`), `StaticOutput`, and command overrides (`--install-cmd`, `--build-cmd`, `--start-cmd`).
4. **Phase 4: Database Suite & SQL Studio**
   - Build the engine selection grid (`TimescaleDB`, `ClickHouse`, etc.), `database-networking.tsx` (`ExternalDNS` + CDC toggle), `data-browser.tsx` (row editing), and the Monaco `sql-studio.tsx`.
5. **Phase 5: Importers & Super Admin Settings**
   - Implement `railway-importer.tsx`, `migration-settings.tsx` (`.vessl` tar.gz export/import), `dns-settings.tsx`, and `maintenance-settings.tsx` (`docker system prune`).
6. **Phase 6: Biome Format & Verification**
   - Ensure every created or modified file passes `npm run format:fix` (`biome check --write .`) with zero errors.
