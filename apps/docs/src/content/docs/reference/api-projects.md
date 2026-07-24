---
title: Projects API
description: Project API endpoints, request payloads, and response examples.
---

Project endpoints manage Codedock project groups. Services live inside projects.

All examples assume:

```bash
export CODEDOCK_URL="https://pilot.example.com"
export CODEDOCK_API_KEY="ap_..."
```

## List Projects

```txt
GET /api/projects
```

Required access: `read`

Project scope: returns all visible projects for the key.

Example:

```bash
curl "$CODEDOCK_URL/api/projects" \
  -H "Authorization: Bearer $CODEDOCK_API_KEY"
```

Response:

```json
{
  "projects": [
    {
      "id": "project_123",
      "name": "Acme",
      "slug": "acme",
      "description": "Production services",
      "status": "active",
      "serviceCount": 2,
      "lastUpdatedAt": "2026-06-10T08:40:00.000Z",
      "services": [
        {
          "id": "svc_web",
          "projectId": "project_123",
          "name": "Web",
          "slug": "web",
          "repoFullName": "acme/web",
          "repoUrl": "https://github.com/acme/web",
          "dockerImage": null,
          "branch": "main",
          "rootDir": null,
          "hasGithubToken": false,
          "installCommand": null,
          "buildCommand": null,
          "startCommand": null,
          "staticOutput": null,
          "runtimeMode": "web",
          "internalPort": 3000,
          "hostPort": 41001,
          "databasePublicEnabled": false,
          "databasePublicHostname": null,
          "postgresLogicalReplicationEnabled": false,
          "status": "active",
          "reachable": false,
          "localUrl": "",
          "primaryUrl": "",
          "preferredDomain": null,
          "framework": null,
          "lastDeployedAt": "2026-06-10T08:35:00.000Z",
          "createdAt": "2026-06-10T08:30:00.000Z",
          "updatedAt": "2026-06-10T08:35:00.000Z"
        }
      ]
    }
  ]
}
```

## Create Project

```txt
POST /api/projects
```

Required access: `write`

Project scope: requires an all-project key. Selected-project keys cannot create projects.

Payload:

```json
{
  "name": "Acme",
  "description": "Production services"
}
```

Example:

```bash
curl -X POST "$CODEDOCK_URL/api/projects" \
  -H "Authorization: Bearer $CODEDOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme","description":"Production services"}'
```

Response:

```json
{
  "project": {
    "id": "project_123",
    "name": "Acme",
    "slug": "acme",
    "description": "Production services",
    "status": "idle",
    "serviceCount": 0,
    "lastUpdatedAt": "2026-06-10T08:40:00.000Z",
    "services": []
  }
}
```

## Get Project

```txt
GET /api/projects/:projectSlug
```

Required access: `read`

Project scope: project must be visible to the key.

Example:

```bash
curl "$CODEDOCK_URL/api/projects/acme" \
  -H "Authorization: Bearer $CODEDOCK_API_KEY"
```

Response:

```json
{
  "project": {
    "id": "project_123",
    "name": "Acme",
    "slug": "acme",
    "description": "Production services",
    "status": "active",
    "serviceCount": 2,
    "lastUpdatedAt": "2026-06-10T08:40:00.000Z",
    "services": []
  }
}
```

## Update Project

```txt
PATCH /api/projects/:projectId
```

Required access: `write`

Project scope: project must be visible to the key.

Payload:

```json
{
  "name": "Acme Production",
  "description": "Customer-facing services"
}
```

Example:

```bash
curl -X PATCH "$CODEDOCK_URL/api/projects/project_123" \
  -H "Authorization: Bearer $CODEDOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme Production","description":"Customer-facing services"}'
```

Response:

```json
{
  "project": {
    "id": "project_123",
    "name": "Acme Production",
    "slug": "acme",
    "description": "Customer-facing services",
    "status": "active",
    "serviceCount": 2,
    "lastUpdatedAt": "2026-06-10T08:45:00.000Z",
    "services": []
  }
}
```

## Delete Project

```txt
DELETE /api/projects/:projectId
```

Required access: `write`

Project scope: project must be visible to the key.

Deleting a project removes its services, domains, environment variables, deployments, and deployment logs.

Example:

```bash
curl -X DELETE "$CODEDOCK_URL/api/projects/project_123" \
  -H "Authorization: Bearer $CODEDOCK_API_KEY"
```

Response:

```json
{
  "ok": true,
  "traefik": {
    "ok": true,
    "detail": "Traefik reloaded"
  }
}
```

## Database Variable Suggestions

```txt
GET /api/projects/:projectId/database-variable-suggestions
```

Required access: `read`

Project scope: project must be visible to the key.

Example:

```bash
curl "$CODEDOCK_URL/api/projects/project_123/database-variable-suggestions" \
  -H "Authorization: Bearer $CODEDOCK_API_KEY"
```

Response:

```json
{
  "suggestions": [
    {
      "key": "POSTGRES_URL",
      "value": "${postgres-db.POSTGRES_URL}",
      "label": "PostgreSQL private URL"
    }
  ]
}
```
