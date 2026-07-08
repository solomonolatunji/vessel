# Agent Instructions

## Code Style

- **One component per file.** Never cram thousands of lines into a single file. Break components down into individual files.
- **Name files in `kebab-case`** (e.g. `project-card.tsx`, `use-logs-stream.ts`).
- Use named exports over default exports.
- Prefer concise code. No JSDoc/TSDoc comments unless the logic is non-obvious.

## Workflow

- **Do not run build or test commands after every change.** Just make the code change. If something breaks, the user will say so.
- Run `npm run format:fix` before finishing a session to keep formatting clean.
- Prefer `read`/`grep`/`glob` tools over `bash` for file exploration.
- When making edits, read the file first, then use `edit` for targeted changes.

## Stack

| Layer                | Tech                                                                         |
| -------------------- | ---------------------------------------------------------------------------- |
| Frontend (dashboard) | React 19, TanStack Router, TanStack Query, Radix UI, Tailwind CSS v4, Vite   |
| Marketing (website)  | Astro 7, Tailwind CSS v4                                                     |
| Backend              | Go (`cmd/vesseld`, `internal/`)                                              |
| State (dashboard)    | TanStack Store, TanStack Query, Zod validation                               |
| Styling (dashboard)  | `tailwind-merge` + `clsx` + `class-variance-authority` for class composition |
| Monorepo             | npm workspaces (`dashboard/`, `website/`)                                    |

## Conventions

- **Dashboard routes** live in `dashboard/src/routes/` following TanStack Router file conventions. Generated route tree is in `routeTree.gen.ts` — do not edit by hand.
- **Dashboard components** go in `dashboard/src/components/`, grouped by domain (e.g. `projects/`, `databases/`, `ui/`).
- **Hooks** go in `dashboard/src/hooks/`.
- **Lib/utils** go in `dashboard/src/lib/`.
- **Marketing pages** live in `website/src/pages/`, components in `website/src/components/`.
- Use Tailwind CSS v4 `@theme` directives for design tokens; avoid custom CSS where Tailwind utilities suffice.
- Format with Prettier (root `.prettierrc` + `.prettierignore`).

## Go Conventions

- **File naming:** lowercase snake_case (`container_health.go`).
- **Package naming:** short, lowercase, single word (`store`, `api`, `orchestrator`).
- **Error handling:** always check errors; wrap with `fmt.Errorf("context: %w", err)`.
- **No global state.** Pass dependencies via struct fields — wire up in `cmd/vesseld/main.go`.
- **JSON tags** on every exported struct field.
- Use `modernc.org/sqlite` (CGO-free) for SQLite. No `database/sql` driver imports for `mattn/go-sqlite3`.
- Use official `github.com/docker/docker/client` for Docker SDK. Use `gorilla/websocket` for WebSocket upgrades.
- Hashicorp-style Go layout: `internal/` packages are private, `cmd/` binaries are thin entrypoints.
- Avoid `init()` functions. Use explicit constructor functions instead.
