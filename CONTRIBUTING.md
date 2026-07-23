# Contributing to Codedock 🛰️

> "First, thanks for considering contributing. It really means a lot!"

Ask for guidance on our [Discord server](https://discord.gg/codedock) in `#contribute`.

---

## Table of Contents

1. [Setup Development Environment](#1-setup-development-environment)
2. [Verify Installation](#2-verify-installation)
3. [Fork and Clone](#3-fork-and-clone)
4. [Environment Variables](#4-environment-variables)
5. [Start Codedock](#5-start-codedock)
6. [Start Developing](#6-start-developing)
7. [Pull Requests](#7-pull-requests)
8. [Development Notes](#8-development-notes)
9. [Reset Dev Environment](#9-reset-dev-environment)

---

## 1. Setup Development Environment

### Prerequisites

- **Go** `v1.22+`
- **Node.js** `v22.12+`
- **Docker** `v24+`

### Linux

```bash
# Go
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Node.js
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo bash -
sudo apt install -y nodejs

# Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

### macOS

```bash
brew install go node docker
```

Or use [Orbstack](https://docs.orbstack.dev/) instead of Docker Desktop (faster, lighter).

### Windows

1. Install [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install) with Ubuntu.
2. Inside WSL2, follow the Linux steps above.
3. Or use [Docker Desktop](https://docs.docker.com/desktop/install/windows-install/) with WSL2 backend.

---

## 2. Verify Installation

```bash
go version    # v1.22+
node --version # v22.12+
docker --version # v24+
npm --version
```

---

## 3. Fork and Clone

1. Fork [codedock](https://github.com/buildwithtechx/codedock) on GitHub.
2. Clone your fork:

```bash
git clone https://github.com/<your-username>/codedock.git
cd codedock
git remote add upstream https://github.com/buildwithtechx/codedock.git
```

---

## 4. Environment Variables

```bash
cp .env.example .env
```

| Variable           | Default          | Description                    |
| ------------------ | ---------------- | ------------------------------ |
| `PORT`             | `8080`           | Daemon HTTP port               |
| `CODEDOCK_DATA_DIR`   | `data`           | SQLite DB + vault storage      |
| `CODEDOCK_STATIC_DIR` | `dashboard/dist` | Built dashboard files          |
| `CODEDOCK_TLS_EMAIL`  | —                | Let's Encrypt email (optional) |

---

## 5. Start Codedock

Two terminals needed — daemon + dashboard.

### Terminal 1: Go Daemon

```bash
go run ./cmd
# API at http://localhost:8080
```

### Terminal 2: Dashboard

```bash
cd dashboard
npm install
npm run dev
# UI at http://localhost:3000 with HMR
```

### Optional: Marketing Site

```bash
cd web
npm install
npm run dev
# http://localhost:4321
```

---

## 6. Start Developing

| Tool      | URL                            | Purpose               |
| --------- | ------------------------------ | --------------------- |
| Dashboard | `http://localhost:3000`        | Main UI (HMR enabled) |
| API       | `http://localhost:8080`        | Go REST API           |
| Health    | `http://localhost:8080/health` | API health check      |
| Marketing | `http://localhost:4321`        | Public site           |

---

## 7. Pull Requests

### Branch

```bash
git checkout -b feat/my-thing next
git add .
git commit -m "feat(daemon): add container CPU usage websocket stream"
git push origin feat/my-thing
```

Open a PR on GitHub with `base: next` and `compare: feat/my-thing`.

### Guidelines

- **Target `next`**, never `main`. PRs against `main` are closed.
- **One PR = one thing**. No bundled unrelated changes.
- **Title**: `type(scope): description` — e.g. `fix(db): prevent nil pointer on empty env`.
- **Description**: What, why, how to test, screenshots if UI, linked issue.
- **Draft**: Use for WIP. Convert when ready. Stale drafts (>7d) may be closed.
- **AI code**: You must understand every line. AI submissions you can't explain will be rejected.

### Review Process

- Maintainers review promptly. Complex PRs take longer.
- Address all feedback. Unresolved comments block merge.
- Merge requires:
  - `go build ./cmd/... ./internal/...` passes
  - `npm run format:fix` + `tsc --noEmit` pass
  - Code review approval
- PRs closed for inactivity (>7d), guideline violations, or being superseded.

### Formatting

```bash
# Go
go fmt ./...
go vet ./...

# TypeScript/JS
npm run format:fix   # Biome, not Prettier
tsc --noEmit
```

---

## 8. Development Notes

### Go Conventions

- `internal/` uses horizontal layers: `models/` → `repositories/` → `services/` → `handlers/` → `engine/`.
- **No comments allowed.** Code must be self-explanatory.
- JSON tags on every exported struct field.
- Snake_case for file names, single-word package names.
- SQLite via `modernc.org/sqlite` (CGO-free). Migrations auto-run on startup in `cmd/codedockd/main.go`.

### Dashboard Conventions

- Routes in `src/routes/` — TanStack Router file conventions. Run `npm run generate-routes` after adding routes.
- shadcn components: `npx shadcn@latest add <name>` → `src/components/ui/`.
- Biome is the formatter — never use Prettier.

### Adding a New API Endpoint

1. Add handler function in `internal/handlers/<name>.go`.
2. Register route in `internal/http/routes.go`.
3. Wire dependencies in `internal/http/server.go`.

### Troubleshooting

| Problem                   | Fix                                                      |
| ------------------------- | -------------------------------------------------------- |
| Daemon won't start        | Delete `data/codedock.db` and restart (schema auto-creates) |
| Dashboard can't reach API | Ensure daemon runs on `:8080`                            |
| Port conflict             | Change `PORT` in `.env`                                  |
| Build errors after pull   | `go mod tidy` + `npm install`                            |

---

## 9. Reset Dev Environment

```bash
# Stop daemon (Ctrl+C), then:
rm -f data/codedock.db data/.vault_key
go run ./cmd
```

Schema is recreated automatically.

---

## Reporting Bugs

Include:

- `uname -r` (kernel version)
- `docker --version`
- `codedockd` logs (run with `--debug` for verbose)
- Steps to reproduce
