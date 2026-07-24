---
title: CLI Reference
description: Command-line interface reference for Codedock — covering both the server daemon (codedockd) and the remote client (codedock).
---

Codedock ships two CLI tools with distinct responsibilities.

| Tool        | Runs on            | Connects to              |
| ----------- | ------------------ | ------------------------ |
| `codedockd` | Your VPS / server  | SQLite + Docker directly |
| `codedock`  | Your local machine | `codedockd` over HTTP    |

---

## `codedockd` — Server Daemon CLI

The `codedockd` binary is the Codedock server process. It runs on your VPS and exposes the HTTP API that the dashboard and the `codedock` remote CLI consume. It also doubles as a management CLI for direct server-side operations without needing the dashboard.

### Server Commands

#### `serve`

Starts the Codedock daemon. This is the default command when no subcommand is provided.

```sh
codedockd serve
```

#### `setup`

Runs the interactive setup wizard to initialise the database and create the initial admin account.

```sh
codedockd setup
```

#### `reset-password`

Resets the password for the admin account. Useful if you lose access to the dashboard.

```sh
codedockd reset-password
```

#### `config`

View or update global server configuration variables.

```sh
codedockd config
```

#### `restart`

Gracefully restarts the Codedock daemon via Docker Compose.

```sh
codedockd restart
```

#### `mcp`

Runs the Model Context Protocol (MCP) server over standard I/O for AI assistant integrations.

```sh
codedockd mcp
```

#### `version`

Prints the current daemon version.

```sh
codedockd version
```

### Deployment Commands

#### `deploy`

Deploy an application directly from the server terminal.

```sh
# From a Git repository
codedockd deploy https://github.com/your/repo.git

# From a template (e.g. go-fiber, nextjs)
codedockd deploy --template nextjs

# From a Docker image
codedockd deploy --image nginx:latest --port 80

# From a Docker Compose file
codedockd deploy --compose ./docker-compose.yml
```

### Resource Management Commands

All resource commands use a `<resource>:<action>` syntax and interact with the database directly — no HTTP, no auth required.

#### Projects

```sh
codedockd project:list                        # List all projects
codedockd project:show <id>                   # Show project details
codedockd project:create <name>               # Create a project
codedockd project:destroy <id>                # Delete a project
```

#### Applications

```sh
codedockd apps:list                           # List all apps across all projects
codedockd apps:show <id>                      # Show app details and env vars
codedockd apps:create <name> --project <id>   # Create an app
codedockd apps:destroy <id>                   # Delete an app
```

#### Databases

```sh
codedockd db:list                             # List all databases
codedockd db:show <id>                        # Show database details and connection string
codedockd db:create <name> <engine> --project <id>  # Create a database
codedockd db:destroy <id>                     # Delete a database
```

Supported engines: `postgres`, `mysql`, `mariadb`, `redis`, `mongodb`, `clickhouse`, `kafka`, `rabbitmq`, `nats`.

#### Environment Variables

```sh
codedockd env:list --project <id>             # List all env vars for a project
codedockd env:set KEY=VALUE --project <id>    # Set one or more env vars
codedockd env:unset KEY --project <id>        # Remove an env var
```

#### Deployments & Logs

```sh
codedockd deployment:list --service <id>      # List deployment history for a service
codedockd deployment:show <id>                # Show deployment details
codedockd deployment:logs <id>                # Print build logs for a deployment
```

#### Custom Domains

```sh
codedockd domain:list --project <id>          # List custom domains for a project
codedockd domain:add <hostname> --project <id> # Add a custom domain
codedockd domain:remove <id>                  # Remove a custom domain
```

---

## `codedock` — Remote CLI

The `codedock` binary runs on your **local machine** and communicates with your self-hosted `codedockd` server over HTTP. This is what you install and use day-to-day from your laptop.

### Installation

```sh
curl -fsSL https://get.codedock.run/cli | sh
```

Or if you have Go installed:

```sh
go install codedock.run/codedock/cmd/codedock@latest
```

Or download a pre-built binary from the [releases page](https://github.com/buildwithtechx/codedock/releases).

### Authentication

Before running any command, authenticate against your self-hosted server.

#### `login`

Prompts for your server URL, email, and password. Saves a token to `~/.codedock/config.json`.

```sh
codedock login
```

#### `logout`

Clears your saved credentials.

```sh
codedock logout
```

#### `me`

Shows the currently authenticated user.

```sh
codedock me
```

### Projects

```sh
codedock project list                         # List all projects
codedock project create <name>                # Create a project
codedock project destroy <id>                 # Delete a project
```

### Environments

```sh
codedock env list --project <id>             # List environments for a project
codedock env create <name> --project <id>    # Create an environment
codedock env destroy <id>                    # Delete an environment
```

### Applications

```sh
codedock apps list --environment <id>        # List apps in an environment
codedock apps create                         # Create an app (interactive flags)
codedock apps destroy <id>                   # Delete an app
```

#### Secrets (Environment Variables)

```sh
codedock apps secrets list --project <id>             # List env vars
codedock apps secrets set KEY=VALUE --project <id>    # Set one or more env vars
```

#### Custom Domains

```sh
codedock apps domains list --project <id>             # List custom domains
codedock apps domains add --domain <host> --project <id>  # Add a domain
codedock apps domains remove <id>                     # Remove a domain
```

#### Deployments & Logs

```sh
codedock apps deployments list --service <id>         # List deployment history
codedock apps logs <deployment-id>                    # View build logs
```

### Databases

```sh
codedock db list --project <id>              # List databases
codedock db create                           # Provision a database (interactive flags)
codedock db destroy <id>                     # Delete a database
```

#### Backups

```sh
codedock db backups list --project <id>      # List backup configurations
codedock db backups create                   # Create a backup config
codedock db backups trigger <id>             # Trigger a manual backup
codedock db backups history <id>             # View backup history
```

### Trigger a Deployment

```sh
codedock deploy <service-id>                 # Trigger a remote deployment for a service
```
