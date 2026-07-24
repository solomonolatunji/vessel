---
title: API & CLI
description: Programmatic access to Codedock via Personal Access Tokens and the REST API.
---

Codedock exposes a REST API for programmatic access, CI/CD integration, and automation.

## Personal Access Tokens

Generate tokens to authenticate API requests without using your password.

### Creating a Token

1. Go to **Profile → Tokens**.
2. Click **Generate Token**.
3. Enter a name for the token.
4. Select scopes:

| Scope           | Access                                  |
| --------------- | --------------------------------------- |
| `deploy:write`  | Trigger deployments and rollbacks       |
| `logs:read`     | Stream deployment and service logs      |
| `env:read`      | Read environment variables              |
| `env:write`     | Create and update environment variables |
| `db:manage`     | Create, modify, and delete databases    |
| `project:read`  | Read project information                |
| `project:write` | Create and modify projects              |
| `admin`         | Full administrative access              |

1. Click **Generate**.

Copy the token immediately — it starts with `vpt_` and won't be shown again.

### Using a Token

```sh
curl -H "Authorization: Bearer vpt_xxx" \
  https://your-codedock-instance/api/projects
```

### Revoking a Token

1. Go to **Profile → Tokens**.
2. Click the delete button next to the token.

## API Endpoints

### Authentication

| Method | Path                        | Description            |
| ------ | --------------------------- | ---------------------- |
| POST   | `/api/auth/signup`          | Register a new account |
| POST   | `/api/auth/signin`          | Log in                 |
| POST   | `/api/auth/logout`          | Log out                |
| POST   | `/api/auth/forgot-password` | Request password reset |
| POST   | `/api/auth/reset-password`  | Reset password         |
| POST   | `/api/auth/2fa/setup`       | Setup 2FA              |
| POST   | `/api/auth/2fa/verify`      | Verify 2FA code        |
| POST   | `/api/auth/2fa/disable`     | Disable 2FA            |

### Projects

| Method | Path                       | Description         |
| ------ | -------------------------- | ------------------- |
| GET    | `/api/projects`            | List all projects   |
| POST   | `/api/projects`            | Create a project    |
| GET    | `/api/projects/:id`        | Get project details |
| DELETE | `/api/projects/:id`        | Delete a project    |
| POST   | `/api/projects/:id/deploy` | Trigger deployment  |

### Services

| Method | Path                         | Description                  |
| ------ | ---------------------------- | ---------------------------- |
| GET    | `/api/environments/:id/apps` | List services in environment |
| POST   | `/api/environments/:id/apps` | Create a service             |
| GET    | `/api/apps/:id`              | Get service details          |
| PUT    | `/api/apps/:id`              | Update a service             |
| DELETE | `/api/apps/:id`              | Delete a service             |
| POST   | `/api/apps/:id/stop`         | Stop service                 |
| POST   | `/api/apps/:id/restart`      | Restart service              |
| POST   | `/api/apps/:id/redeploy`     | Redeploy service             |
| GET    | `/api/services/:serviceId/serverless/code` | Get serverless function code |
| POST   | `/api/services/:serviceId/serverless/code` | Save serverless function code |

### Deployments

| Method | Path                                   | Description         |
| ------ | -------------------------------------- | ------------------- |
| GET    | `/api/services/:serviceId/deployments` | List deployments    |
| POST   | `/api/services/:serviceId/deploy`      | Trigger deployment  |
| POST   | `/api/deployments/:id/rollback`        | Rollback deployment |
| GET    | `/api/deployments/:id/logs`            | Get deployment logs |
| GET    | `/api/services/:serviceId/metrics`     | Get service metrics |

### Databases

| Method | Path                       | Description          |
| ------ | -------------------------- | -------------------- |
| GET    | `/api/databases`           | List databases       |
| POST   | `/api/databases`           | Create a database    |
| GET    | `/api/databases/:id`       | Get database details |
| DELETE | `/api/databases/:id`       | Delete a database    |
| POST   | `/api/databases/:id/start` | Start database       |
| POST   | `/api/databases/:id/stop`  | Stop database        |
| POST   | `/api/databases/:id/query` | Execute SQL query    |

### Storage

| Method | Path                     | Description            |
| ------ | ------------------------ | ---------------------- |
| GET    | `/api/storage`           | List storage instances |
| POST   | `/api/storage`           | Create storage         |
| GET    | `/api/storage/:id`       | Get storage details    |
| DELETE | `/api/storage/:id`       | Delete storage         |
| POST   | `/api/storage/:id/start` | Start storage          |
| POST   | `/api/storage/:id/stop`  | Stop storage           |

### Variables

| Method | Path                                     | Description            |
| ------ | ---------------------------------------- | ---------------------- |
| GET    | `/api/projects/:id/env`                  | Get project env vars   |
| PUT    | `/api/projects/:id/env`                  | Set project env vars   |
| GET    | `/api/services/:serviceId/variables`     | List service variables |
| POST   | `/api/services/:serviceId/variables`     | Add service variable   |
| PUT    | `/api/services/:serviceId/variables/:id` | Update variable        |
| DELETE | `/api/services/:serviceId/variables/:id` | Delete variable        |

### Domains

| Method | Path                        | Description     |
| ------ | --------------------------- | --------------- |
| GET    | `/api/projects/:id/domains` | List domains    |
| POST   | `/api/projects/:id/domains` | Add a domain    |
| DELETE | `/api/domains/:id`          | Remove a domain |

### Webhooks & Volumes

| Method | Path                                   | Description               |
| ------ | -------------------------------------- | ------------------------- |
| GET    | `/api/apps/:id/webhooks`               | List app webhooks         |
| POST   | `/api/apps/:id/webhooks`               | Create an app webhook     |
| DELETE | `/api/apps/:id/webhooks/:webhookId`    | Delete a webhook          |
| GET    | `/api/apps/:id/volumes`                | List app storage volumes  |
| POST   | `/api/apps/:id/volumes`                | Create an app volume      |
| DELETE | `/api/apps/:id/volumes/:volumeId`      | Delete an app volume      |
| GET    | `/api/services/:id/log-drains`         | List log drains           |
| POST   | `/api/services/:id/log-drains`         | Create a log drain        |
| DELETE | `/api/services/:id/log-drains/:drainId`| Delete a log drain        |



### Backups

| Method | Path                       | Description          |
| ------ | -------------------------- | -------------------- |
| GET    | `/api/backups`             | List backup configs  |
| POST   | `/api/backups`             | Create backup config |
| GET    | `/api/backups/:id`         | Get backup config    |
| DELETE | `/api/backups/:id`         | Delete backup config |
| POST   | `/api/backups/:id/trigger` | Trigger backup       |
| GET    | `/api/backups/:id/records` | List backup records  |
| GET    | `/api/s3-destinations`     | List S3 destinations |
| POST   | `/api/s3-destinations`     | Create S3 destination|
| DELETE | `/api/s3-destinations/:id` | Delete S3 destination|

### Miscellaneous & Integrations

| Method | Path                            | Description                    |
| ------ | ------------------------------- | ------------------------------ |
| POST   | `/api/compose/deploy`           | Deploy Docker Compose stack    |
| POST   | `/api/deploy/archive`           | Deploy from ZIP archive        |
| POST   | `/api/one-click/deploy`         | Deploy one-click template      |
| GET    | `/api/dns`                      | List DNS records               |
| POST   | `/api/dns`                      | Create DNS record              |
| GET    | `/api/scheduled-tasks`          | List scheduled tasks           |
| POST   | `/api/scheduled-tasks`          | Create scheduled task (cron)   |
| POST   | `/api/git/connect`              | Connect a Git provider         |
| GET    | `/api/canvas/projects`          | Get architecture canvas data   |
| GET    | `/api/audit-logs`               | List instance audit logs       |

### Settings

| Method | Path                    | Description            |
| ------ | ----------------------- | ---------------------- |
| GET    | `/api/settings`         | Get server settings    |
| PUT    | `/api/settings`         | Update server settings |
| GET    | `/api/system/public`    | Get public settings    |
| GET    | `/api/ai`               | Get AI config settings |
| PUT    | `/api/ai`               | Update AI settings     |
| POST   | `/api/ai/diagnose`      | Run AI Log Diagnosis   |
| POST   | `/api/settings/license` | Activate license       |

### Real-Time

| Type      | Path                            | Description        |
| --------- | ------------------------------- | ------------------ |
| WebSocket | `/api/ws/terminal/:id`          | Container terminal |
| WebSocket | `/api/ws/services/:id/terminal` | Service terminal   |
| SSE       | `/api/mcp/sse`                  | MCP SSE endpoint   |

## Rate Limiting

Authentication endpoints are rate-limited to prevent abuse. Limits are applied per IP address:

- Sign-up: 5 requests per minute
- Login: 10 requests per minute
- Password reset: 3 requests per minute

## Errors

The API returns standard HTTP status codes:

| Code | Meaning                                 |
| ---- | --------------------------------------- |
| 200  | Success                                 |
| 201  | Created                                 |
| 400  | Bad request (validation error)          |
| 401  | Unauthorized (missing or invalid token) |
| 403  | Forbidden (insufficient permissions)    |
| 404  | Not found                               |
| 429  | Rate limited                            |
| 500  | Internal server error                   |

Error responses include a JSON body with details:

```json
{
  "error": "validation_error",
  "message": "Name is required",
  "details": {
    "field": "name"
  }
}
```
