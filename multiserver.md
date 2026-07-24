# Multi-Server & Codedock Cloud Roadmap

This document outlines the architectural changes, features, and steps required to transform Codedock from a single-node manager into a **Distributed Fleet Manager** (Multi-Server) and **Codedock Cloud SaaS**.

## 1. Architectural Strategy: ✅ Agent-based Worker Daemon (DECIDED)

**Decision: Worker Daemon** — chosen because Codedock Cloud is the long-term goal.

### Why Worker Daemon wins for SaaS

| Concern                    | Agentless SSH                                                   | Worker Daemon ✅                                                       |
| -------------------------- | --------------------------------------------------------------- | ---------------------------------------------------------------------- |
| **Security**               | Control plane dials INTO user's server (requires open SSH port) | Worker dials OUT to `api.codedock.dev` (outbound only, firewall-friendly) |
| **Deployment reliability** | Deployment dies if connection drops                             | Deployment continues locally, reports back on reconnect                |
| **Real-time metrics**      | SSH polling — expensive at scale                                | Persistent WebSocket push — cheap at scale                             |
| **Monetisation gate**      | No natural gate                                                 | License key validated on worker registration                           |
| **Enterprise acceptance**  | Blocked by most firewalls                                       | Accepted (same model as GitHub Actions runners)                        |

### Architecture

```text
User's VPS                          Codedock Cloud (api.codedock.dev)
┌─────────────────────┐             ┌──────────────────────────┐
│  codedock-worker       │──WebSocket──▶  Control Plane (Go API)  │
│  - Runs deployments │◀─Commands───│  - Dashboard UI          │
│  - Streams logs     │──Metrics───▶│  - Billing (Stripe)      │
│  - Reports health   │             │  - License validation    │
└─────────────────────┘             └──────────────────────────┘
```

### How it works

1. User creates an account on `app.codedock.dev` and gets a **license/registration key**.
2. User runs a one-liner on their VPS: `curl -sL get.codedock.dev | bash -s -- --key <LICENSE_KEY>`
3. The script installs `codedock-worker` as a systemd service.
4. `codedock-worker` dials out to `api.codedock.dev` via WebSocket and registers itself.
5. The user's server appears in their dashboard. All deployments, logs, and metrics flow through the persistent WebSocket tunnel.

### Worker Binary (`cmd/codedock-worker/`)

- Written in Go, single static binary (~15 MB stripped).
- Connects to control plane via `gorilla/websocket`.
- Executes deployment commands received from the control plane using the local Docker socket.
- Streams container logs and metrics back via the same WebSocket.
- Reconnects automatically with exponential backoff if the connection drops.

---

## 2. Database Schema Changes

We need to make Codedock aware of physical servers. The **project** is the server boundary — all apps and databases inside a project inherit the server automatically. Users never pick a server when creating individual apps or databases.

- **`servers` table (New):**
  - `id` (UUID)
  - `user_id` — owner (for Codedock Cloud multi-tenancy)
  - `name` (e.g., "EU Production")
  - `ip_address` (e.g., "198.51.100.1")
  - `status` (`online`, `offline`, `provisioning`)
  - `worker_token` — the secret the `codedock-worker` binary uses to authenticate
  - `last_seen_at` — heartbeat timestamp
  - `metrics` (JSON — latest CPU/RAM/Disk snapshot pushed by the worker)

- **`projects` table (Update):**
  - Add `server_id` FK → `servers.id`
  - `NULL` means local (for single-node self-hosted installs — backward compatible)

- **`app_services` & `databases` tables — NO CHANGE ✅**
  - They resolve their server by joining through their parent project.
  - No `server_id` column needed here. The user never selects a server per-app or per-database.

---

## 3. Core Engine Updates (`internal/engine/`)

The Deployer engine currently assumes `localhost`. It needs to become **Project-Server-Aware**.

- **Command Routing via Worker WebSocket:**
  Instead of executing Docker commands directly, the deployer checks the project's `server_id`. If a server is attached, it serialises the deployment command as a JSON message and sends it down the worker's persistent WebSocket connection. The worker executes it locally and streams results back.

  ```text
  Deployer.Deploy(app)
    → look up app.Project.ServerID
    → if nil: run locally (existing behaviour, no change)
    → if set:  serialize command → send to WorkerHub → worker executes → stream back
  ```

- **WorkerHub (`internal/engine/worker_hub.go` — New):**
  Maintains a registry of live `server_id → WebSocket connection` pairs. When a worker connects, it registers here. When the deployer needs a remote server, it looks up the live connection from the hub.

- **No SFTP / SSH needed:** The worker has direct access to its own Docker socket. Build context (Dockerfile, source) is sent as a binary payload over the WebSocket, not via SFTP.

---

## 4. Networking & Traefik Routing

Currently, the single Traefik container routes all traffic. In a multi-server setup:

- **Distributed Proxies:** Every Worker Node must run its own instance of Traefik.
- **DNS Resolution:** When an app is deployed to `Server B`, the dashboard must show the user the IP address of `Server B` so they can point their Custom Domain's A-Record to the correct worker node.
- **Wildcard Domains:** If the user has a wildcard domain (e.g., `*.apps.mycodedock.com`), the DNS A-Record for the wildcard must point to the specific server hosting those apps, or we need a central load balancer.

---

## 5. Frontend & UI Integrations

- **Servers Dashboard (`/dashboard/servers` — New):**
  - List all connected worker servers with live status (Online/Offline), last seen, CPU/RAM.
  - "Add Server" page — shows the one-liner install command pre-filled with a fresh `worker_token`.
  - Per-server detail page: resource graphs, connected projects, events log.

- **Project Creation (Update — minor):**
  - Add an optional **"Deploy to Server"** dropdown. Defaults to "Local" for self-hosted installs.
  - Once set on the project, all apps and databases inside automatically deploy to that server. No per-app or per-database server selection ever shown. ✅

- **App/Database Creation — NO CHANGE ✅**
  - Server is fully inherited from the project. Users never see a server dropdown here.

- **Metrics UI:**
  - Per-server resource graphs on the Servers page.
  - Per-project resource graphs continue to work as before (metrics aggregated from the project's server).

---

## 6. Codedock Cloud (SaaS) Considerations

Once Multi-Server is built, launching Codedock Cloud is trivial:

1. Host the **Codedock Control Plane** centrally at `app.codedock.dev`.
2. Users create accounts (user auth already exists in the codebase).
3. Users spin up any VPS (DigitalOcean, Hetzner, AWS EC2, bare-metal), run the one-liner install command from their Servers dashboard, and the server appears live in seconds.
4. Users create a Project, select the server to deploy it to, and add their apps/databases as normal.
5. **Billing:** Stripe integration gates plan limits (e.g., Free = 1 server, Pro = 5 servers, Enterprise = unlimited).

---

## TODO Checklist

### Backend — Models & Repositories

- [x] Create `servers` model (`internal/models/server.go`) — id, user_id, name, ip_address, status, worker_token, last_seen_at, metrics JSON.
- [x] Create `ServerRepository` (`internal/repositories/server.go`) — CRUD + `GetByWorkerToken`.
- [x] Add `server_id` (nullable FK) to `projects` table only — NOT to app_services or databases.

### Backend — Services & Handlers

- [x] Create `ServerService` (`internal/services/server.go`) — create server, generate worker token, list, delete.
- [x] Create `ServerHandler` (`internal/handlers/server.go`) — REST endpoints for server management.
- [x] Add server routes to `internal/http/routes.go`.

### Backend — Worker Engine

- [x] Create `WorkerHub` (`internal/engine/worker_hub.go`) — registry of `server_id → live WebSocket conn`.
- [x] Create Worker WebSocket endpoint (`/ws/worker`) — workers dial in, authenticate with `worker_token`, register in the hub.
- [x] Update `Deployer` — if `project.ServerID != nil`, route deployment command through `WorkerHub` instead of local Docker socket.
- [ ] WorkerHub handles heartbeats and updates `servers.last_seen_at` + `servers.status`.

### Worker Binary (New)

- [x] Scaffold `cmd/codedock-worker/` — new Go entrypoint.
- [x] Worker connects to control plane via WebSocket using its `worker_token`.
- [x] Worker receives deployment commands (JSON), executes them on the local Docker socket.
- [ ] Worker streams container logs and CPU/RAM/disk metrics back over the WebSocket.
- [ ] Worker reconnects with exponential backoff on disconnect.
- [ ] Worker installs Traefik on first boot if not already running.
- [ ] Build & release `codedock-worker` as a separate binary in CI.

### Frontend

- [ ] Create Servers dashboard page (`/servers`) — list servers, status, metrics, last seen.
- [ ] Create Add Server page — generates and displays the one-liner install command with the worker token pre-filled.
- [ ] Update Project creation form — add optional "Deploy to Server" dropdown (defaults to "Local").
- [ ] App and Database creation forms — NO changes needed. ✅
- [ ] Per-server resource graphs (CPU/RAM/Disk over time).

### Codedock Cloud SaaS (Later)

- [ ] Add user registration + email verification (currently only admin account exists).
- [ ] Integrate Stripe for subscription billing.
- [ ] Add plan-based limits (server count, project count) validated server-side.
- [ ] Set up hosted deployment of control plane (`app.codedock.dev`).

## 7. Competitive Analysis (Codedock vs Dokploy)

### Features Dokploy has that Codedock doesn't (Yet)

- **PR Previews:** Dokploy can automatically spin up ephemeral environments when a Pull Request is opened on GitHub, and tear it down when closed. We don't have this yet.
- **More Git Providers:** We currently only support GitHub. Dokploy supports GitHub, GitLab, Bitbucket, and Gitea.
- **Private Docker Registries:** They allow users to link AWS ECR, Google GCR, or private DockerHub accounts to pull private images.
- **Organizations & Teams:** They have a full RBAC (Role-Based Access Control) system where users can create Organizations, invite members, and assign permissions.
- **Volume Backups:** They can back up persistent Docker volumes to S3. We currently only back up Databases natively.

### Features we both have, but we did MUCH better (Where Dokploy went wrong)

- **Multi-Server Clustering (Docker Swarm vs Worker Daemon):**
  - *Their mistake:* They used Docker Swarm, forcing users to open SSH ports and deal with extremely fragile Swarm overlay networking (UDP 4789).
  - *Our fix:* Our WebSocket Worker Daemon (what we are building now) dials outward over standard HTTPS (port 443), bypassing firewalls completely. It's infinitely more stable and secure.
- **Background Tasks & Schedules:**
  - *Their mistake:* Because Node.js is single-threaded, Dokploy had to build completely separate microservices (apps/monitoring and apps/schedules) just to run cron jobs without freezing the dashboard.
  - *Our fix:* Go has goroutines. We run our cron scheduler and metrics monitors inside the exact same binary seamlessly, drastically reducing RAM usage and deployment complexity.
- **Database Models:**
  - *Their mistake:* They created separate database tables and APIs for every single database type (postgres.ts, mysql.ts, mongo.ts, mariadb.ts, redis.ts). If they want to add a feature to databases, they have to update 5 different files.
  - *Our fix:* We have one unified database.go model with an engine_type enum. It’s vastly cleaner to maintain.

### A Feature we both had, but WE did totally wrong (The NATS / Type Safety issue)

- **Type Safety between Backend and Workers:**
  - *Dokploy's win:* Because Dokploy is 100% TypeScript (Next.js frontend, Node.js backend), they use tRPC. If they rename a variable in the backend, the frontend instantly throws a compiler error. Perfect type safety.
  - *Our massive flaw (historically):* As noted earlier, our Go backend was sending raw JSON payloads over NATS to the workers without a single source of truth. If a Go struct changed, the worker wouldn't know until it crashed in production!
  - *The fix:* We are moving to the new Worker Architecture using shared schemas (like Protobuf or central types) so our Go backend and Go workers share the exact same structs natively.
