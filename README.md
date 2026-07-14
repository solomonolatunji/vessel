# 🛰️ Vessl

**Self-hosted PaaS. Turn any VPS into your own Vercel or Railway in 60 seconds.**

---

Vessl is a lightweight, open-source Platform-as-a-Service (PaaS) designed to simplify deployments. Whether you're deploying a static site, a full-stack monorepo, or a complex microservice architecture, Vessl provides a frictionless developer experience without the vendor lock-in.

## 🚀 Quick Start

Install Vessl on any fresh Linux server (Ubuntu/Debian recommended):

```bash
curl -fsSL https://get.vessl.dev | sh
```

Once installed, your dashboard will be available at `http://your-server-ip:8080`.

## ✨ Features

Vessl is built to be simple but powerful, giving you everything you need to run production workloads out of the box.

- **Deploy Anything:** Native support for Dockerfiles, Railpack, Nixpacks, and standard Buildpacks.
- **Managed Databases:** Provision PostgreSQL, MySQL, Redis, MongoDB, and more with a single click.
- **Smart Environment:** Database credentials (`DATABASE_URL`, `REDIS_URL`) are automatically injected into your linked applications.
- **Zero-Downtime Deploys:** Seamless container swaps with built-in health checks and instant rollbacks.
- **Custom Domains & SSL:** Automatic Let's Encrypt certificates managed via Traefik v3.
- **GitOps Ready:** Connect to GitHub/GitLab for automatic deployments on push and PR preview environments.
- **Marketplace Templates:** Instantly deploy popular frameworks (Node.js, Go, Python, Ruby, PHP) from our built-in marketplace.
- **No Lock-in:** Vessl orchestrates standard Docker containers. If you ever remove Vessl, your apps keep running.

## 💻 CLI Operations

Manage your entire infrastructure directly from the terminal.

```bash
# General Management
vessld status                # View cluster health and running containers
vessld setup                 # Run the initial admin configuration wizard
vessld reset-password        # Reset the dashboard admin password
vessld logs -f               # Tail global system logs

# Deployments
vessld deploy --template node-express            # Deploy an official starter template
vessld deploy https://github.com/user/repo.git   # Deploy any Git repository
vessld deploy --image nginx:latest               # Deploy a pre-built Docker image
vessld deploy --compose docker-compose.yml       # Deploy a full Docker Compose stack

# App & Resource Management
vessld apps:list                                 # List all running applications
vessld db:create my-db postgres --project <id>   # Provision a managed database
```

## 🛠️ Local Development

Want to contribute or hack on Vessl locally?

```bash
# 1. Clone the repository
git clone https://github.com/vesslhq/vessl.git
cd vessl

# 2. Setup your environment
cp .env.example .env

# 3. Run the Go daemon locally (Starts on port :8080)
go run ./cmd
```

**Requirements:** Go 1.25+, Node.js 22+, and Docker.

## 📚 Documentation

For complete guides, API references, and advanced configuration, please visit our documentation at **[docs.vessl.dev](https://docs.vessl.dev)**.
