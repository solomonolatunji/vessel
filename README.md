# 🛰️ Vessl

Self-hosted PaaS. Turn any VPS into your own Vercel or Railway in 60 seconds.

```bash
curl -fsSL https://get.vessl.dev | sh
```

Dashboard at `http://your-server-ip:8080`.

## Features

- **Deploy apps** — Docker, Railpack, Nixpacks, Buildpacks, or Serverless
- **15 databases** — PostgreSQL, MySQL, MariaDB, MongoDB, Redis, Dragonfly, KeyDB, ClickHouse, Kafka, RabbitMQ, NATS, NocoDB, Plausible, WordPress, Gitea
- **S3 storage** — One-click MinIO buckets, auto-injected credentials
- **Auto connection strings** — `DATABASE_URL`, `REDIS_URL`, etc. injected automatically
- **Custom domains** — Let's Encrypt SSL via Traefik v3
- **Git auto-deploy** — Connect GitHub/GitLab, push to deploy, PR previews
- **Workspaces & teams** — RBAC, audit logs, SSH keys, trusted domain SSO
- **Serverless functions** — Node.js, Python, Go in-browser editor
- **Cron jobs** — Scheduled tasks inside containers
- **Notifications** — Discord, Slack, Telegram, Pushover, SMTP, Resend, webhooks
- **AI diagnostics** — Analyze failed builds via OpenAI/Anthropic
- **Zero-downtime deploys** — Health-checked container swaps with rollback
- **No lock-in** — Standard Docker containers. Remove Vessl, apps keep running
- **CLI management** — Full admin from terminal without the dashboard

## CLI

```bash
vesslctl status                # Health + containers
vesslctl setup                 # Admin wizard
vesslctl reset-password        # Reset admin password
vesslctl config                # View config
vesslctl config site-name=Prod # Update setting
vesslctl logs -f               # Tail logs
vesslctl restart               # Restart daemon (applies domain/config changes)
vesslctl deploy <git-url>           # Deploy app from Git
vesslctl deploy --image nginx:latest  # Deploy from Docker image
vesslctl apps:list             # List apps
vesslctl db:create my-db postgres --project <id>  # Create DB
vesslctl backup                # Backup database
vesslctl update                # Upgrade Vessl
vesslctl downgrade v0.1.0      # Downgrade
vesslctl uninstall             # Remove Vessl, keep apps
```

## Local Dev

```bash
cp .env.example .env
go run ./cmd           # Daemon on :8080
```

Requires Go 1.25+, Node.js 22+, Docker.

## Docs

[docs.vessl.dev](https://docs.vessl.dev)
