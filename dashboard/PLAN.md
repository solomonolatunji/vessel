# Vessl Dashboard In-Depth Plan (v2)

After a deep secondary analysis of the existing Go handlers and API routes, this plan has been expanded to cover the complete feature surface area of Vessl. It goes beyond a simple PaaS and includes Enterprise features, Super Admin capabilities, and AI integrations.

## 1. Core Layout & Navigation Architecture

The layout utilizes a contextual **Sidebar** and a global **Topbar**.

### Topbar

- **Workspace/Team Context Switcher:** Allows users to switch between Personal Account, various Teams, or Workspaces.
- **Global Command Menu (Cmd+K):** Deep-link search across all projects, environments, databases, and settings.
- **Notification Center:** Real-time updates via SSE for deployments, backup completions, or team invites.
- **Profile / Theme:** Dark/Light mode, link to Personal Settings.

### Contextual Sidebar States

#### A. Global / Workspace Sidebar

_Applies when no specific project is selected._

- **Overview:** Aggregated stats (running containers, total deployments, active jobs).
- **Projects:** Grid/List of projects.
- **Databases:** Global view of all DBs (Postgres, MySQL, Redis, etc.) across the workspace.
- **Storage:** S3-compatible MinIO instances.
- **Backups & S3 Destinations:** Global backup policies and connected S3 buckets.
- **Teams & Workspaces:**
  - Members & Roles
  - Trusted Domains (SSO config)
  - SSH Keys
  - Audit Logs (`/teams/:teamId/audit-logs`)
- **Git Integrations:** GitHub/GitLab/Bitbucket app connections.
- **AI & Notifications:**
  - AI Settings (OpenAI/Anthropic keys for the team)
  - Alert Channels (Discord, Slack, Email)

#### B. Project Sidebar

_Applies when navigating inside a specific project._

- **Project Overview:** Readme, status indicators, quick metrics.
- **Canvas View (`/environments/:id/canvas`):** A Railway-inspired interactive node graph. Visualize how Apps, Databases, and Storage connect within an environment.
- **Environments:** Switch between Production, Staging, and Previews.
- **Apps / Services:** The microservices running in this project.
- **Project Settings:**
  - Webhooks (Outgoing)
  - Project Tokens (For CLI/CI usage)
  - Members (Project-level RBAC)

#### C. Service / App Deep Dive

_Applies when managing a specific application/service._

- **Overview:** Live CPU/Memory metrics (`/services/:serviceId/metrics`).
- **Deployments:** List of builds. Includes **Rollback** button (`/deployments/:id/rollback`).
- **Logs:** Live streaming logs with timestamp filtering.
- **Variables:** Environment variables and secrets (`/services/:serviceId/variables`).
- **Domains:** Custom domains, routing prefixes, SSL status (`/projects/:id/domains`).
- **Serverless Editor:** In-browser code editor for edge functions (`/services/:serviceId/serverless/code`).
- **Terminal:** Web-based SSH into the running container (`/ws/services/:id/terminal`).
- **Settings:** Build commands, start commands, replicas, auto-deploy triggers.

#### D. Database View

_Applies when viewing a Database instance._

- **Connection Details:** URI strings, ports, credentials.
- **Controls:** Start/Stop container.
- **SQL Studio (`/databases/:id/query`):** A raw SQL execution playground built directly into the dashboard!
- **Snapshots:** Trigger manual backups or view history.

#### E. Super Admin Settings (Instance Management)

_Because Vessl is self-hosted, instance owners need an admin panel (`/settings`)._

- **License & Updates:** Activate licenses, check for Vessl updates, and trigger 1-click updates (`/settings/updates/deploy`).
- **Global Config:** Configure the Traefik wildcard IP, custom DNS resolvers, generic webhooks, and instance-wide SMTP/Resend configs.

---

## 2. Advanced Features & Flows

### 2.1 AI Diagnostics & MCP Integration

Vessl has built-in AI capabilities:

- **Build Failures:** When a deployment fails, an "AI Diagnose" button calls `POST /deployments/:id/diagnostics` to analyze the logs and suggest a fix.
- **MCP Chat:** The dashboard will feature a floating AI chat assistant. Because the backend implements an MCP server (`/mcp/sse`), the AI can actively query project states, restart services, or read logs on behalf of the user.

### 2.2 Authentication & Security Flows

- **Login/Signup:** Standard flows with OAuth provider support.
- **2FA:** Dedicated flows for `/auth/2fa/setup` and enforcement during login.
- **Audit Logs:** A dedicated table view for enterprise teams to track who did what and when.

### 2.3 PR Previews

- **Flow:** When a webhook from GitHub arrives (`/webhooks/github/services/:serviceId`), Vessl spins up an ephemeral environment.
- **UI:** The dashboard will have a "Previews" tab under the Project, showing temporary active deployments linked to Pull Requests.

### 2.4 Vercel Importer

- **Flow:** The Vessl dashboard will provide an integration (`/api/oauth/vercel/callback`, `/api/vercel/projects`) to authenticate with Vercel, pull down a user's Vercel projects and environment variables, and seamlessly migrate them to Vessl as serverless functions or container deployments.

---

## 3. Folder Structure (Feature-Driven Architecture)

Instead of putting all hooks, components, and API calls in global folders, we isolate them by business domain under `src/features/`.

```text
src/
├── features/                  # Domain-specific modules
│   ├── auth/                  # Login, Signup, 2FA, OAuth flows
│   ├── projects/              # Project listings, creation wizard, Vercel Importer
│   ├── canvas/                # Visual environment node graph (Railway style)
│   ├── services/              # Apps, Serverless Code editor, Live Logs, Terminal
│   ├── databases/             # DB instances, SQL Studio, Backups
│   ├── storage/               # Object storage buckets
│   ├── jobs/                  # Background Cron Tasks
│   ├── workspaces/            # Team settings, Members, Integrations
│   └── instance/              # Super Admin Vessl settings, Updates
│
├── routes/                    # TanStack file-based routing
│   ├── __root.tsx             # Root layout & providers
│   ├── _auth/                 # Auth layout (no sidebar)
│   │   ├── login.tsx
│   │   └── register.tsx
│   ├── _workspace/            # Global Workspace layout (Sidebar: Global Context)
│   │   ├── index.tsx          # Dashboard overview
│   │   ├── projects.tsx
│   │   ├── databases.tsx
│   │   └── settings.tsx
│   ├── _project/              # Project layout (Sidebar: Project Context)
│   │   ├── $projectId/
│   │   │   ├── index.tsx      # Project summary & Canvas
│   │   │   ├── deployments.tsx
│   │   │   ├── services/
│   │   │   │   └── $serviceId.tsx # Deep dive into a specific app
│   │   │   └── settings.tsx
│
├── components/                # Global UI primitives
│   ├── ui/                    # Shadcn/Radix components (Button, Input, etc)
│   ├── layout/                # Shell, Sidebar, Topbar
│   └── shared/                # Status badges, Error boundaries
│
├── hooks/                     # Global utilities (useTheme, useToast)
├── lib/                       # API clients (axios/fetch), auth tokens, formatting
└── types/                     # Global TypeScript interfaces
```

---

## 4. Implementation Phases

### Phase 1: Foundation & Auth

- [x] Configure Vite + TanStack Router + Tailwind CSS.
- [x] Install necessary packages (React Flow, Monaco, XTerm, Shadcn UI).
- [ ] Setup `src/lib/apiClient.ts` to handle auth tokens and global fetch interceptors.
- [ ] Create `_auth` layout (clean, centered card design).
- [ ] Implement `src/routes/auth/login.tsx` and `register.tsx`.

### Phase 2: Core Layouts

- [ ] Build the Topbar component (Workspace switcher, User profile).
- [ ] Build the dynamic Sidebar component.
- [ ] Create `_workspace` layout routing.
- [ ] Create `_project` layout routing (Sidebar items change based on `$projectId`).

### Phase 3: Workspace & Project Views

- [ ] Implement Workspace Dashboard (`/projects`, `/databases`).
- [ ] Implement Project Summary page (`/projects/$projectId`).
- [ ] Implement the **Interactive Canvas** using React Flow (`features/canvas`).
  - _Must allow dragging DB nodes to App nodes to link them._

### Phase 4: Aeroplane/Vercel Deep Dive Features

- [ ] Service overview and metrics.
- [ ] Implement live log streaming UI (`features/services/components/LiveLogs`).
- [ ] Integrate XTerm.js for the Web Terminal (`/ws/terminal/:id`).
- [ ] Integrate Monaco Editor for Serverless edge functions.
- [ ] Integrate Monaco Editor for the Database **SQL Studio**.

### Phase 5: Admin & Enterprise Polish

- [ ] Team management, RBAC, and Invites UI.
- [ ] Git webhook and repository connection flow.
- [ ] "Super Admin" settings panel for self-hosted Vessl instance updates.
- [ ] Dark Mode refinements and responsive design checks.
