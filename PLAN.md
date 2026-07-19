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
- **Templates (`/templates`):** One-click templates and object storage integrations.
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

### 1. Browser Onboarding Wizard

- **Target Route / Component:** `src/routes/_auth/setup.tsx`, `src/features/auth/setup-form.tsx`
- **UI Specification & User Flow:**
  - **First-Run Interception:** If `GET /system/setup-status` indicates no owner account exists, redirect all traffic to `/setup`.
  - **Wizard Steps:**
    1. **Owner Account:** Email and password setup.
    2. **Runtime Environment:** Configure `DATA_DIR`, `PORT`, `PUBLIC_URL` which writes to `.env.local`.
    3. **Domains & Hostnames:** Freedom to configure Control Plane Domain (`pilot.example.com` or IP) and Wildcard Root Domain for generated apps.
    4. **Backup & Object Storage:** Connect S3/R2 bucket credentials.
    5. **Restore Instance (Optional):** Dropzone to upload `.vessl` server state bundle.
  - _Note: GitHub is NOT configured here (it's a simple 1-click install button in the dashboard)._

### 2. One-Line Installer

- **Target Route / Component:** `/routes/getting-started.tsx`
- **UI Specification & User Flow:**
  - Show highlighted command (`curl -fsSL https://get.vessl.dev | sh`) and system readiness checklist.

### 3. Service Runtime Modes

- **Target Route / Component:** `src/features/services/service-settings.tsx`
- **UI Specification & User Flow:**
  - **Web vs. Worker Switcher:** Radio card selector when creating/editing an `AppService`.
  - **Web:** Internal port input, public route generator, HTTP health checks (`/healthz`).
  - **Background Worker:** No internal port, no public route, process uptime check badge (`runtimeMode === 'worker'`).

### 4. Static Site Deployments

- **Target Route / Component:** `src/features/services/build-settings.tsx`
- **UI Specification & User Flow:**
  - **Static Output Input:** Text field for `Static output directory` (e.g., `dist`, `build`, `.output/public`). When set, UI displays badge indicating the service runs inside an optimized `nginx:alpine` wrapper on internal port 80.

### 5. Zero-Downtime Hot Swaps

- **Target Route / Component:** `src/features/services/service-deployments.tsx`
- **UI Specification & User Flow:**
  - **Live Transition UI:** During `deploying` status, display both the active `running` container (`UUID-A`) and the probing `starting` container (`UUID-B`). Show real-time Traefik health check status before old container cleanup.

### 6. Build Overrides

- **Target Route / Component:** `src/features/services/build-settings.tsx`
- **UI Specification & User Flow:**
  - **Command Override Inputs:** Expandable accordion under Railpack/Nixpacks settings:
    - **Install Command:** (`--install-cmd`) e.g., `npm ci`
    - **Build Command:** (`--build-cmd`) e.g., `npm run build`
    - **Start Command:** (`--start-cmd`) e.g., `npm start`

### 7. Intelligent Env Var Linking

- **Target Route / Component:** `src/features/services/service-variables.tsx`
- **UI Specification & User Flow:**
  - **Smart Variable Drawer:** When editing service `.env` secrets, a side-drawer suggests auto-linked variables (`${postgres-db.POSTGRES_URL}`, `${timescaledb-db.TIMESCALE_URL}`). Includes `.env.example` parser pills to quickly autofill required keys.

### 8. Database Data Imports

- **Target Route / Component:** `src/features/databases/database-import-modal.tsx`
- **UI Specification & User Flow:**
  - **Import Data Modal (`POST /databases/:id/import`):**
    - **URL Import:** Input for `postgres://` or `redis://` schemes with immediate syntax validation.
    - **Railway Sync:** Auto-detects public DB URLs from imported Railway variables.
    - **TimescaleDB Check:** Displays source/target extension compatibility warning badge.
    - **History Table:** Live streaming progress of `pg_dump` / `redis-cli --rdb` jobs.

### 9. Server Migration Bundles

- **Target Route / Component:** `src/features/instance/migration-settings.tsx`, `/routes/_shell/settings/migration.tsx`
- **UI Specification & User Flow:**
  - **Bundle Manager (`/settings/migration`):**
    - **Export Card:** Passphrase input + `Export Server Bundle` button (`GET /system/export`).
    - **Import Card:** File upload dropzone (`.vessl`) + Passphrase + `Restore Server State` destructive action modal (`POST /system/import`).

### 10. Database Provisioning Engines

- **Target Route / Component:** `src/features/databases/create-database-modal.tsx`
- **UI Specification & User Flow:**
  - **Engine Selection Grid:** Organized categorizations with custom icon badges:
    - **Relational:** PostgreSQL (`16-alpine`), **TimescaleDB (`latest`)**, MySQL (`8.0`), MariaDB (`11`), ClickHouse (`latest`)
    - **NoSQL:** MongoDB (`7.0`), Redis (`7-alpine`), Dragonfly (`latest`), KeyDB (`latest`)
    - **Brokers:** Kafka (`9092`), RabbitMQ (`5672`), NATS (`4222`)
    - **One-Click:** NocoDB, Plausible, WordPress, Gitea

### 11. Data Browser & Row Editing

- **Target Route / Component:** `src/features/databases/data-browser.tsx`, `/routes/_shell/databases/$dbId/data.tsx`
- **UI Specification & User Flow:**
  - **Relational Table Grid (`GET /databases/:id/data/:table`):**
    - Table switcher (`GET /databases/:id/schemas`), filtering (`=`, `contains`, `>`), server-side pagination.
    - **Inline Row Editing:** Double-click cells to edit or click `+ Add Row` (`POST /databases/:id/data/:table`).
  - **Redis Key Browser:** Specialized grid showing key names, types (`string`, `hash`, `list`, `set`), values, and interactive TTL editor.

### 12. Public Database Access & TLS

- **Target Route / Component:** `src/features/databases/database-networking.tsx`
- **UI Specification & User Flow:**
  - **Public Hostname Controller (`PUT /databases/:id`):**
    - Toggle for **Public Access (`ExternalDNS`)**.
    - Displays generated TCP endpoint: `postgres-db.pilot.example.com:5432`.
    - TLS Status badge (`Let's Encrypt TCP SNI enabled`) + one-click copy buttons for external clients.

### 13. Logical Replication (CDC)

- **Target Route / Component:** `src/features/databases/database-settings.tsx`
- **UI Specification & User Flow:**
  - **CDC Toggle Switch:** In Postgres & TimescaleDB configuration, toggle `Logical Replication (wal_level=logical)`. Displays warning badge: _"Enables max_replication_slots=10 and WAL retention for Change Data Capture tools."_

### 14. Database Restore & Download

- **Target Route / Component:** `src/features/databases/backup-manager.tsx`
- **UI Specification & User Flow:**
  - **Backup Action Table (`/databases/:id/backups`):**
    - **Download:** Button to stream `.sql` or `.rdb` directly from disk/R2 (`GET /backups/:id/download`).
    - **Destructive Restore:** Red warning modal asking user to type database name before piping backup (`pg_restore --clean`, `mysql <`, `mongorestore`) into running container (`POST /backups/:id/restore`).

### 15. DNS Provider Automation

- **Target Route / Component:** `src/features/instance/dns-settings.tsx`, `src/features/projects/project-domains.tsx`
- **UI Specification & User Flow:**
  - **Provider Credentials (`/settings/dns`):** Forms to save API keys for Cloudflare, Namecheap, and Spaceship.
  - **1-Click A-Record Sync:** On any service custom domain card (`/projects/:id/domains`), a `Sync A Record via DNS Provider` button that automatically writes the `1800` TTL record targeting the server IP.

### 16. System Maintenance & Cleanup

- **Target Route / Component:** `src/features/instance/maintenance-settings.tsx`, `/routes/_shell/settings/maintenance.tsx`
- **UI Specification & User Flow:**
  - **Maintenance Dashboard (`/settings/maintenance`):**
    - **Storage Gauges:** Root FS disk %, Docker storage reclaimable MBs, Backup volume size.
    - **Garbage Collection:** `Run Docker Cleanup Now` button (`POST /system/maintenance/cleanup`) running `docker system prune -af --volumes`.
    - **Cron Config:** Schedule selector for automated background pruning (`Docker Cleanup Cron`).

### 17. Railway Importer Specification

- **Target Route / Component:** `src/features/projects/railway-importer.tsx`, `/routes/_shell/import/railway.tsx`
- **UI Specification & User Flow:**
  - **Multi-Step Railway Import Wizard (`POST /import/railway`):**
    1. **Token & Discovery:** Paste Railway Personal API Token (`Bearer <token>`), query GraphQL v2, select project.
    2. **Service Classification Table:** Displays detected Git repos, Docker images, and database engines (Postgres, TimescaleDB, Redis, Mongo, ClickHouse).
    3. **Configuration Checkboxes:** `Exclude RAILWAY_* variables` (default ON), `Recreate database engines` (creates local Vessl DBs), `Auto-deploy services`, and `Import database data` (runs automated `pg_dump`/`redis-cli` from public Railway URLs).

### 18. Control Plane Auto-Updates

- **Target Route / Component:** `src/features/instance/update-settings.tsx`, `/routes/_shell/settings/updates.tsx`
- **UI Specification & User Flow:**
  - **Version Controller (`/settings/updates`):**
    - **Version Card:** Shows `Current Version`, `Latest Version` (`GET /settings/updates/check`), and release notes.
    - **Auto-Update Toggle:** Switch for `Auto Update Enabled` + `Update Check Cron` selector.
    - **Manual Trigger:** `Deploy Update Now` button (`POST /settings/updates/deploy`) triggering `scripts/upgrade.sh` and graceful `vessld` container restart.

### 19. Global Domain Management

- **Target Route / Component:** `src/features/domains/domain-list.tsx`, `/routes/_shell/domains.tsx`
- **UI Specification & User Flow:**
  - **Domains (`/domains`):** Centralized view of all custom domains mapped across services (`GET /domains`).

### 20. Storage & S3 Buckets

- **Target Route / Component:** `src/features/storage/storage-list.tsx`, `/routes/_shell/storage.tsx`
- **UI Specification & User Flow:**
  - **Storage (`/storage`):** List all active storage buckets, view status, and manage their lifecycles (`GET /storage`, `POST /storage/:id/start`, `POST /storage/:id/stop`).

### 21. Global Deployments View

- **Target Route / Component:** `src/features/deployments/deployments-list.tsx`, `/routes/_shell/deployments.tsx`
- **UI Specification & User Flow:**
  - **Global Deployments (`/deployments`):** A centralized view of all active and historical builds across all projects and services.

### 22. Background Jobs

- **Target Route / Component:** `src/features/jobs/jobs-list.tsx`, `/routes/_shell/jobs.tsx`
- **UI Specification & User Flow:**
  - **Jobs (`/jobs`):** View all background jobs and cron tasks across the platform (`GET /jobs`, `POST /jobs/:id/trigger`).

### 23. Git Sources Integration

- **Target Route / Component:** `src/features/instance/git-apps.tsx`, `/routes/_shell/sources.tsx`
- **UI Specification & User Flow:**
  - **Sources (`/sources`):** Connect and manage Git providers (GitHub, GitLab, Bitbucket), OAuth, and webhooks (`/git/connect`, `/git/repos`, `/webhooks/*`).

### 24. API Access (Tokens)

- **Target Route / Component:** `src/features/profile/access-tokens.tsx`, `/routes/_shell/settings/api.tsx`
- **UI Specification & User Flow:**
  - **Tokens (`/settings/api`):** Manage Personal Access Tokens (`/profile/tokens`) and Project-level Tokens (`/projects/:id/tokens`).

### 25. AI Assistant (MCP)

- **Target Route / Component:** `src/features/ai/ai-chat.tsx`, `/routes/_shell/ai.tsx`
- **UI Specification & User Flow:**
  - **AI (`/ai`):** Copilot interface powered by MCP for querying logs, managing infrastructure, and generating configurations (`/mcp/sse`, `/mcp/messages`).

### 26. Templates & One-Click Apps

- **Target Route / Component:** `src/features/templates/template-list.tsx`, `/routes/_shell/templates.tsx`
- **UI Specification & User Flow:**
  - **Templates (`/templates`):** Discover and deploy one-click apps from predefined templates (`/one-click`, `/one-click/deploy`).

### 27. Serverless Code Editor

- **Target Route / Component:** `src/features/services/serverless-editor.tsx`
- **UI Specification & User Flow:**
  - **Serverless (`/projects/:id/services/:serviceId/serverless`):** Integrated code editor for serverless functions (`GET/POST /services/:serviceId/serverless/code`).

---

## 4. TanStack Router Route Tree Design

We will structure `src/routes/` with clear functional layout boundaries (`_shell`, `_auth`) and explicit URL paths:

```text
src/routes/
├── __root.tsx                               # Global QueryProvider, ThemeProvider, Toast, CommandMenu
├── _auth.signin.tsx                         # POST /auth/login
├── _auth.signup.tsx                         # POST /auth/register
├── _auth.forgot-password.tsx                # POST /auth/forgot-password
├── _auth.reset-password.tsx                 # POST /auth/reset-password
├── _auth.setup.tsx                          # First-run browser onboarding wizard (/setup)
│
├── _dashboard/                                  # Authenticated layout (Topbar + Contextual Sidebar)
│   ├── index.tsx                            # Global Overview / Dashboard Home
│   ├── projects.tsx                         # Project List (`/projects`)
│   ├── databases.tsx                        # Global Database Inventory (`/databases`)
│   ├── templates.tsx                      # One-Click Apps & Storage Templates (`/templates`)
│   │
│   ├── imports/                             # Migration Importers
│   │   ├── railway.tsx                      # Railway Project Importer (`/import/railway`)
│   │   └── vercel.tsx                       # Vercel Project Importer (`/import/vercel`)
│   │
│   ├── deployments.tsx                      # Global Deployments View (`/deployments`)
│   ├── domains.tsx                          # Global Domains List (`/domains`)
│   ├── storage.tsx                          # Global Storage Buckets (`/storage`)
│   ├── ai.tsx                               # AI Copilot Interface (`/ai`)
│   ├── templates.tsx                        # One-Click Apps & Storage Templates (`/templates`)
│   ├── sources.tsx                          # Git Providers Integration (`/sources`)
│   │
│   ├── settings/                            # Super Admin & Instance Settings (`/settings/*`)
│   │   ├── index.tsx                        # General Instance Configuration
│   │   ├── dns.tsx                          # Cloudflare/Namecheap/Spaceship DNS Credentials
│   │   ├── maintenance.tsx                  # Garbage Collection & Disk Usage Alerts
│   │   ├── updates.tsx                      # Control Plane Version & Auto-Update Toggles
│   │   ├── migration.tsx                    # AES-256 `.vessl` Bundle Export/Import
│   │   ├── users.tsx                        # Instance-Wide User Management
│   │   ├── api.tsx                          # API Access & Access Tokens Management
│   │   ├── oauth.tsx                        # Global OAuth Providers (GitHub, Google)
│   │   └── backups.tsx                      # S3 Backup Destinations
│   │
│   ├── profile/                             # Current User Profile (`/profile`)
│   │   └── index.tsx                        # 2FA, Personal Access Tokens, Change Password
│   │
│   ├── projects/                            # Project Overviews
│   │   ├── $projectId.index.tsx             # Project Overview & Quick Stats
│   │   ├── $projectId.canvas.tsx            # Railway-Style React Flow Node Graph
│   │   ├── $projectId.settings.tsx          # Project RBAC, Webhooks, and Global Secrets
│   │   ├── $projectId.jobs.tsx              # Background Jobs & Cron Tasks
│   │   └── $projectId.compose.tsx           # Docker Compose Deployments
│   │
│   ├── services/                            # Decoupled Services
│   │   ├── $serviceId.index.tsx             # Service Metrics & Overview
│   │   ├── $serviceId.deployments.tsx       # Build History, Logs & Rollback (`/services/$serviceId/deployments`)
│   │   ├── $serviceId.variables.tsx         # Secret Editor & Smart Env Var Linking
│   │   ├── $serviceId.domains.tsx           # Custom Domains, Wildcard & DNS Provider Sync
│   │   ├── $serviceId.build.tsx             # Railpack/Nixpacks Overrides & Static Output
│   │   ├── $serviceId.terminal.tsx          # XTerm.js Container Terminal (`/ws/terminal`)
│   │   └── $serviceId.serverless.tsx        # Monaco Editor for Edge Functions
│   │
│   └── databases/                           # Decoupled Databases
│       ├── $dbId.index.tsx                  # Connection Details, Credentials & Public Networking
│       ├── $dbId.data.tsx                   # Table Browser, Row-Level Editor & Redis Explorer
│       ├── $dbId.query.tsx                  # SQL Studio Playground (`/query`)
│       └── $dbId.backups.tsx                # Manual Snapshots, R2 Sync & Destructive Restore
```

---

## 5. Domain Feature Folder Organization (`src/features/`)

To keep files well below **350 lines** and strictly **one component per file**, every domain has its own modular folder:

```text
src/features/
├── auth/
│   ├── login-form.tsx
│   ├── register-form.tsx
│   ├── forgot-password-form.tsx
│   ├── reset-password-form.tsx
│   ├── o-auth-buttons.tsx
│   ├── setup-form.tsx
│   └── use-auth.ts
│
├── profile/
│   ├── user-profile-form.tsx
│   ├── security-2fa-setup.tsx
│   └── access-tokens-list.tsx
│
├── projects/
│   ├── project-list.tsx
│   ├── project-card.tsx
│   ├── create-project-modal.tsx
│   ├── project-domains.tsx
│   ├── railway-importer.tsx
│   ├── vercel-importer.tsx
│   ├── compose-deploy-form.tsx
│   └── jobs-list.tsx
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
    ├── migration-settings.tsx
    ├── oauth-providers-list.tsx
    ├── git-apps-manager.tsx
    └── s3-destinations-list.tsx
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

## 7. Implementation Roadmap & Build Phases

We will build the dashboard in distinct phases to ensure stability, proper data-binding, and excellent UI/UX consistency across all domains.

**Phase 1: Core Foundation & Shell (Routing & Layout)**

- [x] Scaffold the new `AppSidebar` and `Topbar` adhering to the new group structures (Overview, Resources, System & Settings).
- [x] Establish the `TanStack Router` configuration (`src/routes/`) to reflect the new layout (e.g. `/domains`, `/ai`, `/settings/dns`).
- [x] Setup the core `api-client.ts` to seamlessly intercept 401s and standardize JSON parsing.
- [x] Build the `Auth` & `Setup Wizard` views (`/login`, `/setup`).

**Phase 2: Project & Service Management (The Bread & Butter)**

- [ ] Implement the `Projects` overview and environment grid (`/projects`).
- [ ] Build the `Services` domain (`/services/$serviceId/*`): metrics, build overrides, variables (`.env`), and deployment history.
- [ ] Map the new `/deployments` global view to track all system-wide builds in real-time using `EventSource` (SSE).

**Phase 3: Database & Storage Provisioning**

- [ ] Create the Database Engine Selection UI (`Postgres`, `Redis`, `Mongo`, etc.).
- [ ] Build out the DB dashboards (`/databases/$dbId/*`): connection strings, public networking (`ExternalDNS`), and data/table browser.
- [ ] Implement `Storage` bucket management (`/storage`), fetching and controlling bucket lifecycle.

**Phase 4: Integrations & AI**

- [ ] Map out the `GitHub` integration screens (`/settings/github`), handling OAuth flows and repository syncing (`/git/repos`).
- [x] Develop the `Domains` & `DNS` system (`/domains`, `/settings/dns`), syncing custom domains automatically.
- [x] Build the AI Assistant interface (`/ai`) connected to Vessl's backend knowledge base for auto-generating Compose files and config.

**Phase 5: Super Admin & Maintenance**

- [x] Create `Users` management (`/settings/users`) and API Tokens (`/settings/api`).
- [x] Build the `Maintenance` and `Updates` dashboards (`docker system prune`, `vessld` auto-updates).
- [x] Finalize the `Migration` bundle logic (`.vessl` export/import functionality).
- [x] Build Instance Settings (`/settings/general`), `Notifications` (`/settings/notifications`), and `OAuth` (`/settings/oauth`).
- [ ] Implement `Backups` configuration.

**Phase 6: Polish & Verification**

- [ ] Strict Biome formatting (`npm run format:fix`) and type-checking across all components.
- [ ] Audit all inputs and buttons to ensure minimalist "plain" design (no weird autofill backgrounds, sharp padding/margins).
- [ ] Ensure all components are under 350 lines and correctly modularized in `src/features/`.
      export TEST_TOKEN="<generate-via-script-or-env>"
