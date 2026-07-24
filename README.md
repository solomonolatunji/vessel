# 🛰️ Codedock

**Self-hosted PaaS. Turn any VPS into your own Vercel or Railway in 60 seconds.**

---

Codedock is a lightweight, open-source Platform-as-a-Service (PaaS) designed to simplify deployments. Whether you're deploying a static site, a full-stack monorepo, or a complex microservice architecture, Codedock provides a frictionless developer experience without the vendor lock-in.

## 🚀 Quick Start

Install Codedock on any fresh Linux server (Ubuntu/Debian recommended):

```bash
curl -fsSL https://get.codedock.run | sh
```

Once installed, your dashboard will be available at `http://your-server-ip:8080`.

## ✨ Features

Codedock is built to be simple but powerful, giving you everything you need to run production workloads out of the box.

- **Deploy Anything:** Native support for Dockerfiles, Railpack, Nixpacks, and standard Buildpacks.
- **Managed Databases:** Provision PostgreSQL, MySQL, Redis, MongoDB, and more with a single click.
- **Smart Environment:** Database credentials (`DATABASE_URL`, `REDIS_URL`) are automatically injected into your linked applications.
- **Zero-Downtime Deploys:** Seamless container swaps with built-in health checks and instant rollbacks.
- **Custom Domains & SSL:** Automatic Let's Encrypt certificates managed via Traefik v3.
- **GitOps Ready:** Connect to GitHub/GitLab for automatic deployments on push and PR preview environments.
- **Marketplace Templates:** Instantly deploy popular frameworks (Node.js, Go, Python, Ruby, PHP) from our built-in marketplace.
- **No Lock-in:** Codedock orchestrates standard Docker containers. If you ever remove Codedock, your apps keep running.

## 💻 CLI

Codedock ships two CLI tools. See their individual READMEs for full command references.

| Tool                            | Purpose                                                                | Docs                                |
| ------------------------------- | ---------------------------------------------------------------------- | ----------------------------------- |
| [`codedockd`](./cmd/codedockd/) | Server daemon — runs on your VPS, manages Docker & SQLite directly     | [README](./cmd/codedockd/README.md) |
| [`codedock`](./cmd/codedock/)   | Remote client — runs on your laptop, connects to `codedockd` over HTTP | [README](./cmd/codedock/README.md)  |

**Quick example:**

```sh
# On your server
codedockd serve
codedockd deploy --template nextjs
codedockd deploy https://github.com/user/repo.git

# On your local machine
codedock login
codedock project list
codedock deploy <service-id>
```

## 🛠️ Local Development

Want to contribute or hack on Codedock locally?

````bash
# 1. Clone the repository
git clone https://github.com/buildwithtechx/codedock.git
cd codedock

# 2. Setup your environment
cp .env.example .env

# 3. Run locally
You can run the Go daemon and the frontend dashboard concurrently using the provided Makefiles.

To run the daemon normally + frontend:
```bash
make dev
````

To run the daemon in dry-run mode (skips Docker actions) + frontend:

```bash
make dev-dryrun
```

Alternatively, you can manually build the dashboard and run the daemon:

```bash
cd dashboard && npm install && npm run build && cd ..
go run ./cmd/codedockd
```

**Requirements:** Go 1.22+, Node.js 22+, and Docker.

## 📚 Documentation

For complete guides, API references, and advanced configuration, please visit our documentation at **[docs.codedock.run](https://docs.codedock.run)**.
