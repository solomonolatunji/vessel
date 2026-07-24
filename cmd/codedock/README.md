# `codedock` — Remote CLI

The `codedock` binary is the remote client for your self-hosted Codedock server. It runs on your **local machine** and communicates with `codedockd` over HTTP using a saved token.

## Installation

```sh
curl -fsSL https://get.codedock.run/cli | sh
```

Or if you have Go installed:

```sh
go install codedock.run/codedock/cmd/codedock@latest
```

After installing, authenticate against your server:

```sh
codedock login
```

This prompts for your server URL, email, and password and saves credentials to `~/.codedock/config.json`.

## Commands

### Auth

```sh
codedock login                            # Authenticate to your server
codedock logout                           # Clear saved credentials
codedock me                               # Show current logged-in user
```

### Projects

```sh
codedock project list                     # List all projects
codedock project create <name>            # Create a project
codedock project destroy <id>             # Delete a project
```

### Environments

```sh
codedock env list --project <id>          # List environments
codedock env create <name> --project <id> # Create an environment
codedock env destroy <id>                 # Delete an environment
```

### Applications

```sh
codedock apps list --environment <id>     # List apps
codedock apps create                      # Create an app
codedock apps destroy <id>                # Delete an app
```

#### Secrets (Env Vars)

```sh
codedock apps secrets list --project <id>
codedock apps secrets set KEY=VALUE --project <id>
```

#### Custom Domains

```sh
codedock apps domains list --project <id>
codedock apps domains add --domain <host> --project <id>
codedock apps domains remove <id>
```

#### Deployments & Logs

```sh
codedock apps deployments list --service <id>
codedock apps logs <deployment-id>
```

### Databases

```sh
codedock db list --project <id>           # List databases
codedock db create                        # Provision a database
codedock db destroy <id>                  # Delete a database
```

#### Backups

```sh
codedock db backups list --project <id>
codedock db backups create
codedock db backups trigger <id>
codedock db backups history <id>
```

### Deployments

```sh
codedock deploy <service-id>              # Trigger a remote deployment
```

## Config

Credentials are stored at `~/.codedock/config.json`:

```json
{
  "serverUrl": "https://your-server.com",
  "token": "<jwt>",
  "email": "you@example.com"
}
```

Run `codedock logout` to clear this file.
