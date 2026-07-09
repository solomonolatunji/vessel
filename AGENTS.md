# Agent Instructions

## Code Style

- **One component per file.** Never cram thousands of lines into a single file. Break components down into individual files.
- **Name files in `kebab-case`** (e.g. `project-card.tsx`, `use-logs-stream.ts`).
- Use named exports over default exports.
- Prefer concise code. No inline `//` comments allowed. Only JSDoc/TSDoc (`/** */`) for TS/JS and GoDoc (`// PackageName`) for Go types/funcs, and only when the logic is non-obvious. Do not write GoDoc on types that are self-explanatory (e.g., `// User represents a user`) or on HTTP handlers (handler names already describe what they do).

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
| Backend              | Go (`cmd/vesseld`, `internal/`)                                              |
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

- **Feature/Domain Packaging (`internal/<domain>/`):** Modern Go code inside `internal/` must be organized by feature/domain (e.g., `internal/auth/`, `internal/projects/`, `internal/cron/`) rather than rigid horizontal layers (`api/`, `services/`, `store/`). Each domain package encapsulates its own `model.go`, `dto.go`, `handler.go`, `service.go`, `repository.go` (if needed), and database implementation (`sqlite.go`).
- **Consumer-Defined Interfaces:** Define interfaces where they are _consumed_, not where they are implemented (`Accept interfaces, return structs`). Services define narrow `Repository` interfaces (e.g., `type Repository interface { GetByID(ctx, id) (*Job, error) }`) specifying exact dependencies. Concrete database adapters (`*SQLiteRepository`) satisfy them implicitly without `implements` clauses.
- **Incremental Migration:** New features (e.g. `apikeys/`, `dns/`) must adopt feature-based packaging immediately. Existing legacy modules (`internal/api`, `internal/store`) should be migrated incrementally ("Boy Scout rule") when touched, avoiding massive all-at-once rewrites.
- **File naming:** lowercase snake_case (`container_health.go`).
- **Package naming:** short, lowercase, single word (`cron`, `auth`, `apikeys`).
- **Error handling:** always check errors; wrap with `fmt.Errorf("context: %w", err)`.
- **No global state.** Pass dependencies via struct fields — wire up in `cmd/vesseld/main.go`.
- **JSON tags** on every exported struct field.
- Use `modernc.org/sqlite` (CGO-free) for SQLite. No `database/sql` driver imports for `mattn/go-sqlite3`.
- Use official `github.com/docker/docker/client` for Docker SDK. Use `gorilla/websocket` for WebSocket upgrades.
- Avoid `init()` functions. Use explicit constructor functions instead.
- **Testing:** Integration and end-to-end tests go in `internal/tests/`, grouped by domain in subdirectories (e.g. `internal/tests/auth/`, `internal/tests/database/`). Unit tests (`Service` logic testing via consumer mocks) stay co-located with their source files (`service_test.go`).
