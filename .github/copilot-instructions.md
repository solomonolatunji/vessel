# Vessel — Agent Instructions

## What This Is

This is the Vessel monorepo — an open-source, self-hosted PaaS that turns any bare-metal VPS into a private Vercel & Railway. It's written in Go (backend daemon) and TypeScript (dashboard UI, marketing site, docs).

- **Language (backend)**: Go — `cmd/vesseld`, `internal/`
- **Language (frontend)**: TypeScript/TSX — React 19, Astro 7
- **Runtime (dashboard)**: Vite + TanStack Start
- **Runtime (web/docs)**: Astro 7 + Starlight
- **Monorepo**: npm workspaces (`dashboard/`, `web/`, `docs/`)
- **Database**: embedded SQLite (`modernc.org/sqlite`, CGO-free)
- **Container runtime**: Docker SDK (`github.com/docker/docker/client`)

## Architecture

### Backend (`cmd/`, `internal/`)

1. **Entrypoint** — `cmd/vesseld/main.go`: HTTP server daemon, wires dependencies
2. **API Handlers** — `internal/api/`: REST + WebSocket endpoints (projects, databases, env vars, git, terminal)
3. **Orchestrator** — `internal/orchestrator/`: Multi-language build engine (Dockerfile, Railpack, Nixpacks), container lifecycle, zero-downtime deploys
4. **Proxy** — `internal/proxy/`: Dynamic Caddy v2 reverse proxy config generation and hot-reload
5. **Store** — `internal/store/`: SQLite repositories (projects, domains, env vars, users, invites) + AES-256-GCM `.env` vault
6. **Middleware** — `internal/middleware/`: JWT auth guards, CORS

### Dashboard (`dashboard/`)

React 19 + TanStack Router + TanStack Query + Radix UI + Tailwind CSS v4. The self-hosted control panel where users deploy apps, manage databases, view logs, and configure settings.

### Marketing Site (`web/`)

Astro 7 + Tailwind CSS v4. Public landing page at `vessel.dev` — hero section, feature comparisons, install command.

### Docs (`docs/`)

Astro 7 + Starlight. Documentation site with full-text search, sidebar navigation.

## Code Style & Conventions

### Go & Architecture

- **Feature/Domain Packaging (`internal/<domain>/`)**: Organize packages by domain/feature (`internal/auth/`, `internal/projects/`, `internal/cron/`) encapsulating `model.go`, `dto.go`, `handler.go`, `service.go`, `repository.go`, and `sqlite.go` in each domain. Avoid horizontal layers (`api/`, `store/`).
- **Consumer-Defined Interfaces**: Define narrow `Repository` interfaces where they are *consumed* (`Accept interfaces, return structs`) rather than where implemented (`sqlite.go` satisfies implicitly without `implements` clauses).
- **Files**: `snake_case.go` — `container_manager.go`
- **Packages**: short, lowercase, single word — `cron`, `auth`, `apikeys`
- **No inline `//` comments** — only GoDoc when logic is non-obvious
- **No GoDoc** on self-explanatory types (`// User represents a user`) or HTTP handlers
- **No global state** — pass deps via struct fields, wire in `main.go`
- **Always check errors** — wrap with `fmt.Errorf("context: %w", err)`
- **JSON tags** on every exported struct field
- **No `init()`** — use explicit constructors

### TypeScript (Dashboard)

- **Files**: `kebab-case.tsx` — `project-card.tsx`
- **Named exports** over default exports
- **One component per file**
- Routes in `dashboard/src/routes/` (TanStack Router file conventions)
- Do **not** edit `routeTree.gen.ts` by hand
- Use `tailwind-merge` + `clsx` + `class-variance-authority` for class composition

### npm Workspace Scripts

| Command                 | Action                                   |
| ----------------------- | ---------------------------------------- |
| `npm run dev:dashboard` | Start dashboard at `localhost:3000`      |
| `npm run dev:web`       | Start marketing site at `localhost:4321` |
| `npm run dev:docs`      | Start docs at `localhost:4322`           |
| `npm run build:all`     | Build all workspaces                     |
| `npm run format:fix`    | Format all files with Prettier           |

## Navigation Tips

- **Find an API handler**: `internal/api/<domain>_handler.go`
- **Find a store repository**: `internal/store/<entity>_store.go`
- **Find an orchestrator component**: `internal/orchestrator/<component>.go`
- **Find a dashboard route**: `dashboard/src/routes/` (TanStack Router file-based)
- **Find a dashboard component**: `dashboard/src/components/<domain>/`
- **Find a web page**: `web/src/pages/`
- **Find docs content**: `docs/src/content/docs/`

## Constraints

- Do NOT run build or test commands after every change — the user will say if something breaks
- Do NOT commit unless explicitly asked
- Do NOT edit `routeTree.gen.ts` by hand
- Do NOT add `init()` in Go — use explicit constructors
- Do NOT use `mattn/go-sqlite3` — use `modernc.org/sqlite`
- Do NOT add unnecessary dependencies
- Run `npm run format:fix` before finishing a session
