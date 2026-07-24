# `codedockd` — Server Daemon

`codedockd` is the Codedock server process. It runs on your VPS, manages the SQLite database, orchestrates Docker containers, and exposes the HTTP API consumed by the dashboard and the `codedock` remote CLI.

## Running the Server

```sh
codedockd serve          # Start the daemon (default when no subcommand is given)
```

By default it listens on `:8080`. Configure with environment variables:

| Variable          | Default | Description                               |
| ----------------- | ------- | ----------------------------------------- |
| `PORT`            | `8080`  | HTTP port to listen on                    |
| `HOST`            | ``      | Bind address                              |
| `CODEDOCK_DATA_DIR`  | `data/` | Directory for SQLite DB and secrets vault |
| `CODEDOCK_TLS_EMAIL` | ``      | Email for Let's Encrypt (Traefik)         |

## Setup & Maintenance

```sh
codedockd setup            # Interactive first-time setup wizard
codedockd reset-password   # Reset the admin account password
codedockd config           # View or update server configuration
codedockd restart          # Gracefully restart the daemon
codedockd version          # Print the current version
```

## Deployment Commands

Deploy applications directly from the server terminal without using the dashboard.

```sh
# From a Git repository
codedockd deploy https://github.com/your/repo.git

# From a template (e.g. go-fiber, nextjs)
codedockd deploy --template go-fiber

# From a Docker image
codedockd deploy --image nginx:latest --port 80

# From a Docker Compose file
codedockd deploy --compose ./docker-compose.yml
```

## Resource Management

All commands below operate directly on the database — no HTTP, no auth token required. Useful for admin recovery or scripting.

### Projects

```sh
codedockd project:list
codedockd project:show <id>
codedockd project:create <name>
codedockd project:destroy <id>
```

### Applications

```sh
codedockd apps:list
codedockd apps:show <id>
codedockd apps:create <name> --project <id>
codedockd apps:destroy <id>
```

### Databases

```sh
codedockd db:list
codedockd db:show <id>
codedockd db:create <name> <engine> --project <id>
codedockd db:destroy <id>
```

Supported engines: `postgres`, `mysql`, `mariadb`, `redis`, `mongodb`, `clickhouse`, `kafka`, `rabbitmq`, `nats`.

### Environment Variables

```sh
codedockd env:list --project <id>
codedockd env:set KEY=VALUE --project <id>
codedockd env:unset KEY --project <id>
```

### Deployments & Logs

```sh
codedockd deployment:list --service <id>
codedockd deployment:show <id>
codedockd deployment:logs <id>
```

### Custom Domains

```sh
codedockd domain:list --project <id>
codedockd domain:add <hostname> --project <id>
codedockd domain:remove <id>
```

## Advanced

```sh
codedockd mcp              # Run the MCP stdio server for AI integrations
```
