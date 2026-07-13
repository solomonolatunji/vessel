# Agent Instructions

## Code Style

- **Max 350 lines per file.** If a file exceeds this, extract helpers, adapters, or types into separate files.
- **One component per file.** Never cram thousands of lines into a single file. Break components down into individual files.
- **Name files in `kebab-case`** (e.g. `project-card.tsx`, `use-logs-stream.ts`).
- Use named exports over default exports.
- **No comments or GoDoc/JSDoc allowed.** Code should be self-explanatory. If logic is truly non-obvious, refactor to make it clear rather than adding a comment.

## Workflow

- **Do not run build or test commands after every change.** Just make the code change. If something breaks, the user will say so.
- Run `npm run format:fix` (`biome check --write .`) for TS/JS/JSON and `go fmt ./...` for Go before committing or finishing a session. NEVER run `prettier` (`npx prettier`) — Biome is our strict formatter.
- Prefer `read`/`grep`/`glob` tools over `bash` for file exploration.
- When making edits, read the file first, then use `edit` for targeted changes.

## Stack

| Layer                | Tech                                                                         |
| -------------------- | ---------------------------------------------------------------------------- |
| Frontend (dashboard) | React 19, TanStack Router, TanStack Query, Radix UI, Tailwind CSS v4, Vite   |
| Marketing (web)      | Astro 7, Tailwind CSS v4                                                     |
| Docs                 | Astro 7, Starlight                                                           |
| Backend              | Go (`cmd`, `internal/`)                                                      |
| State (dashboard)    | TanStack Store, TanStack Query, Zod validation                               |
| Styling (dashboard)  | `tailwind-merge` + `clsx` + `class-variance-authority` for class composition |
| Monorepo             | npm workspaces (`dashboard/`, `web/`, `docs/`)                               |

## Conventions

- **Dashboard routes** live in `dashboard/src/routes/` following TanStack Router file conventions. Generated route tree is in `routeTree.gen.ts` — do not edit by hand.
- **Dashboard components** go in `dashboard/src/components/`, grouped by domain (e.g. `projects/`, `databases/`, `ui/`).
- **Hooks** go in `dashboard/src/hooks/`.
- **Lib/utils** go in `dashboard/src/lib/`.
- **Marketing pages** live in `web/src/pages/`, components in `web/src/components/`.
- Use Tailwind CSS v4 `@theme` directives for design tokens; avoid custom CSS where Tailwind utilities suffice.
- **Format strictly with Biome** (`npm run format:fix` / `biome check --write .`) and `go fmt ./...`. NEVER use Prettier (`npx prettier`).

## Go Conventions & Architecture

- **Layered Monolith Architecture (`internal/`):** All Go code inside `internal/` must be organized by clean functional layers:
  - `internal/models/` — ALL domain structs, DTOs, and database entities (no circular imports).
  - `internal/repositories/` — ALL database persistence, SQL interfaces, and SQLite implementations (`project.go`, `user.go`, `auth.go`).
  - `internal/services/` — ALL business logic and external integrations (`auth.go`, `git.go`, `deploy.go`).
  - `internal/handlers/` — ALL HTTP controllers and Echo route handlers (`auth.go`, `project.go`, `oauth.go`).
  - `internal/http/` — HTTP server setup, routes, CORS, and auth middleware wiring.
  - `internal/engine/` — Container engine, Docker deployer, runtime management, cron, and backup workers.
- **Consumer-Defined Interfaces:** Define narrow interfaces where consumed (`Accept interfaces, return structs`).
- **Max 350 lines per file:** If a file exceeds 350 lines, split into smaller focused files.
- **No comments or GoDoc allowed:** Code must be self-explanatory without comments.
- **File naming:** lowercase snake_case (`container_health.go`).
- **Package naming:** short, lowercase, single word (`cron`, `auth`, `apikeys`).
- **Error handling:** always check errors; wrap with `fmt.Errorf("context: %w", err)`.
- **No global state.** Pass dependencies via struct fields — wire up in `cmd/vessld/main.go`.
- **JSON tags** on every exported struct field.
- Use `modernc.org/sqlite` (CGO-free) for SQLite. No `database/sql` driver imports for `mattn/go-sqlite3`.
- Use official `github.com/docker/docker/client` for Docker SDK. Use `gorilla/websocket` for WebSocket upgrades.
- Avoid `init()` functions. Use explicit constructor functions instead.
- **Testing:** Unit tests stay co-located with their source files (`service_test.go`). New domain packages must include tests for their service/handler logic.

## Dashboard

```sh
npm run dev        # starts at http://localhost:3000
npm run build      # output → dashboard/dist/
npm run generate-routes  # regenerate route tree after adding routes
```

**Conventions:**

- Routes live in `dashboard/src/routes/` — TanStack Router file conventions.
- Add shadcn/Radix UI components via `npx shadcn@latest add button`.
- Components go in `src/components/ui/`. Prefer existing Radix over building from scratch.

## Web (Marketing)

```sh
astro dev --background   # use background mode
astro dev stop|status|logs
```

Pure Astro + Tailwind CSS v4 — no React/Vue/Svelte islands.

## Docs

```sh
npm run dev   # starts at http://localhost:4322
npm run build # output → docs/dist/
```

Add pages via `.md` files in `src/content/docs/`. Starlight uses file-based routing.
