# `pkg` — Codedock SDK

This directory contains the public Go packages that power the `codedock` remote CLI. You can import them directly to build your own tooling on top of a self-hosted Codedock server — CI/CD scripts, GitHub Actions, Terraform providers, custom dashboards, etc.

## Packages

### `pkg/http` — API Client

A typed HTTP client for the Codedock API. Handles authentication, request building, and response decoding.

```go
import codedockhttp "codedock.run/codedock/pkg/http"

client := codedockhttp.NewClient("https://your-server.com", "your-jwt-token")
```

#### Auth

```go
// Login and get a token
resp, err := client.Login("admin@example.com", "password")
// resp.Token, resp.User

// Get current user profile
user, err := client.Me()

// Logout (server-side session clear)
err := client.Logout()
```

#### Projects

```go
projects, err := client.ListProjects()
project, err  := client.GetProject("project-id")
project, err  := client.CreateProject(&models.ProjectConfig{Name: "my-app"})
err           := client.DeleteProject("project-id")
```

#### Environments

```go
envs, err := client.ListEnvironments("project-id")
env, err  := client.CreateEnvironment("project-id", &models.EnvironmentConfig{Name: "production"})
err       := client.DeleteEnvironment("env-id")
```

#### Applications

```go
services, err := client.ListServices("env-id")
service, err  := client.GetService("service-id")
service, err  := client.CreateService(&models.AppService{...})
err           := client.DeleteService("service-id")
```

#### Secrets

```go
vars, err := client.GetSecrets("project-id")
err       := client.SetSecrets("project-id", models.SetEnvVarsRequest{...})
```

#### Custom Domains

```go
domains, err := client.ListDomains("project-id")
domain, err  := client.AddDomain("project-id", &models.DomainConfig{DomainName: "app.example.com"})
err          := client.RemoveDomain("domain-id")
```

#### Databases

```go
dbs, err := client.ListDatabases("project-id")
db, err  := client.GetDatabase("db-id")
db, err  := client.CreateDatabase(&models.CreateDatabaseRequest{...})
err      := client.DeleteDatabase("db-id")
```

#### Backups

```go
configs, err := client.ListBackups("project-id")
config, err  := client.CreateBackup(&models.BackupConfig{...})
err          := client.TriggerBackup("backup-id")
records, err := client.ListBackupRecords("backup-id")
```

#### Deployments

```go
deployments, err := client.ListDeployments("service-id")
deployment, err  := client.TriggerDeployment("service-id")
logs, err        := client.GetDeploymentLogs("deployment-id")
metrics, err     := client.GetServiceMetrics("service-id")
```

---

### `pkg/config` — Credential Store

Manages saved server credentials on the local filesystem at `~/.codedock/config.json`.

```go
import "codedock.run/codedock/pkg/config"
```

#### Load saved config

```go
cfg, err := config.Load()
// cfg.ServerURL, cfg.Token, cfg.Email
```

#### Save config

```go
err := config.Save(&config.Config{
    ServerURL: "https://your-server.com",
    Token:     "jwt-token",
    Email:     "you@example.com",
})
```

#### Get config file path

```go
path, err := config.GetConfigPath()
// e.g. /home/user/.codedock/config.json
```

## Full Example

```go
package main

import (
    "fmt"

    "codedock.run/codedock/pkg/config"
    codedockhttp "codedock.run/codedock/pkg/http"
)

func main() {
    cfg, err := config.Load()
    if err != nil || cfg.Token == "" {
        panic("run 'codedock login' first")
    }

    client := codedockhttp.NewClient(cfg.ServerURL, cfg.Token)

    projects, err := client.ListProjects()
    if err != nil {
        panic(err)
    }

    for _, p := range projects {
        fmt.Printf("%s  %s\n", p.ID[:8], p.Name)
    }
}
```
