---
title: Services API
description: Service API endpoints, request payloads, and response examples.
---

Services are deployable units inside projects. A service can be a Git source service, Docker image service, database service, worker, or static site.

Database service creation has its own page: [Databases API](/docs/reference/api-databases/).

All examples assume:

```bash
export CODEDOCK_URL="https://pilot.example.com"
export CODEDOCK_API_KEY="ap_..."
```

## Create Source Service

```txt
POST /api/projects/:projectId/services
```

Required access: `write`

Project scope: project must be visible to the key.

Payload:

```json
{
  "name": "web",
  "repoFullName": "acme/web",
  "branch": "main",
  "rootDir": null,
  "installCommand": "npm install",
  "buildCommand": "npm run build",
  "startCommand": "npm start",
  "buildMethod": "auto",
  "dockerfilePath": null,
  "runtimeMode": "web",
  "internalPort": 3000,
  "env": [
    { "key": "NODE_ENV", "value": "production" }
  ]
}
```

`buildMethod` controls how the image is built:

- `auto` (default) — if the repository contains a `Dockerfile` at its root, the deployment builds it with `docker build`; otherwise Railpack analyzes the project.
- `dockerfile` — always build with the repository's Dockerfile. The deployment fails if it is missing.
- `railpack` — always build with Railpack, even when a Dockerfile is present.

`dockerfilePath` optionally points at a Dockerfile in a non-standard location, relative to the service root directory (for example `docker/Dockerfile.web`). The `CODEDOCK_DOCKERFILE_PATH` service environment variable takes precedence when set. Custom install, build, and start commands do not apply to Dockerfile builds.

Use `repoFullName` for GitHub repositories. Use `repoUrl` instead for a direct Git URL:

```json
{
  "name": "web",
  "repoFullName": null,
  "repoUrl": "https://github.com/acme/web.git",
  "branch": "main",
  "runtimeMode": "web",
  "internalPort": 3000
}
```

Example:

```bash
curl -X POST "$CODEDOCK_URL/api/projects/project_123/services" \
  -H "Authorization: Bearer $CODEDOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "web",
    "repoFullName": "acme/web",
    "branch": "main",
    "runtimeMode": "web",
    "internalPort": 3000,
    "env": [
      { "key": "NODE_ENV", "value": "production" }
    ]
  }'
```

Response:

```json
{
  "service": {
    "id": "svc_web",
    "projectId": "project_123",
    "name": "web",
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
    "buildMethod": "auto",
    "dockerfilePath": null,
    "detectedBuildMethod": null,
    "runtimeMode": "web",
    "internalPort": 3000,
    "hostPort": 41001,
    "databasePublicEnabled": false,
    "databasePublicHostname": null,
    "postgresLogicalReplicationEnabled": false,
    "status": "idle",
    "reachable": false,
    "localUrl": "http://127.0.0.1:41001",
    "primaryUrl": "http://127.0.0.1:41001",
    "preferredDomain": null,
    "framework": null,
    "lastDeployedAt": null,
    "createdAt": "2026-06-10T08:40:00.000Z",
    "updatedAt": "2026-06-10T08:40:00.000Z"
  }
}
```

## Create Docker Image Service

```txt
POST /api/projects/:projectId/services
```

Required access: `write`

Payload:

```json
{
  "name": "api",
  "repoUrl": "docker-image",
  "dockerImage": "ghcr.io/acme/api:latest",
  "branch": "main",
  "runtimeMode": "web",
  "internalPort": 8080,
  "env": [
    { "key": "NODE_ENV", "value": "production" }
  ]
}
```

Response:

```json
{
  "service": {
    "id": "svc_api",
    "projectId": "project_123",
    "name": "api",
    "slug": "api",
    "repoFullName": "image:ghcr.io/acme/api:latest",
    "repoUrl": "docker-image",
    "dockerImage": "ghcr.io/acme/api:latest",
    "branch": "main",
    "rootDir": null,
    "runtimeMode": "web",
    "internalPort": 8080,
    "hostPort": 41002,
    "status": "idle",
    "createdAt": "2026-06-10T08:42:00.000Z",
    "updatedAt": "2026-06-10T08:42:00.000Z"
  }
}
```

## Service Overview

```txt
GET /api/services/:serviceId/overview
```

Required access: `read`

Project scope: service project must be visible to the key.

Example:

```bash
curl "$CODEDOCK_URL/api/services/svc_web/overview" \
  -H "Authorization: Bearer $CODEDOCK_API_KEY"
```

Response:

```json
{
  "service": {
    "id": "svc_web",
    "projectId": "project_123",
    "name": "web",
    "status": "active",
    "runtimeMode": "web",
    "internalPort": 3000,
    "hostPort": 41001,
    "primaryUrl": "https://web.example.com"
  },
  "deployments": [
    {
      "id": "dep_123",
      "serviceId": "svc_web",
      "commitSha": "abc1234",
      "status": "running",
      "trigger": "manual",
      "imageTag": "codedock-svc_web-dep_123",
      "containerName": "codedock-svc_web-stable",
      "startedAt": "2026-06-10T08:45:00.000Z",
      "finishedAt": "2026-06-10T08:46:00.000Z",
      "createdAt": "2026-06-10T08:45:00.000Z"
    }
  ],
  "env": [
    {
      "id": "env_123",
      "key": "NODE_ENV",
      "hasValue": true,
      "value": "production",
      "resolvedValue": "production",
      "createdAt": "2026-06-10T08:40:00.000Z",
      "updatedAt": "2026-06-10T08:40:00.000Z"
    }
  ],
  "domains": [
    {
      "id": "domain_123",
      "serviceId": "svc_web",
      "hostname": "web.example.com",
      "status": "active",
      "createdAt": "2026-06-10T08:41:00.000Z",
      "updatedAt": "2026-06-10T08:41:00.000Z"
    }
  ],
  "publicIp": "203.0.113.10"
}
```

## Update Service

```txt
PATCH /api/services/:serviceId
```

Required access: `write`

Project scope: service project must be visible to the key.

Payload:

```json
{
  "branch": "main",
  "rootDir": "apps/web",
  "buildCommand": "npm run build",
  "startCommand": "npm start",
  "buildMethod": "auto",
  "dockerfilePath": null,
  "runtimeMode": "web",
  "internalPort": 3000
}
```

`detectedBuildMethod` in service responses reports which builder the most recent deployment actually used (`dockerfile` or `railpack`); it is `null` until a deployment has run.

Response:

```json
{
  "service": {
    "id": "svc_web",
    "projectId": "project_123",
    "name": "web",
    "rootDir": "apps/web",
    "buildCommand": "npm run build",
    "startCommand": "npm start",
    "runtimeMode": "web",
    "internalPort": 3000,
    "updatedAt": "2026-06-10T08:50:00.000Z"
  }
}
```

## Transfer Service

```txt
POST /api/services/:serviceId/transfer
```

Required access: `write`

Project scope: key must access both the source project and the target project.

Payload:

```json
{
  "targetProjectId": "project_456"
}
```

Response:

```json
{
  "service": {
    "id": "svc_web",
    "projectId": "project_456",
    "name": "web",
    "slug": "web"
  },
  "project": {
    "id": "project_456",
    "name": "Internal",
    "slug": "internal",
    "serviceCount": 1,
    "services": []
  },
  "traefik": {
    "ok": true,
    "detail": "Traefik reloaded"
  }
}
```

## Delete Service

```txt
DELETE /api/services/:serviceId
```

Required access: `write`

Project scope: service project must be visible to the key.

Example:

```bash
curl -X DELETE "$CODEDOCK_URL/api/services/svc_web" \
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

## Variable Suggestions

```txt
GET /api/services/:serviceId/suggestion-keys
```

Required access: `read`

Response:

```json
{
  "suggestions": [
    {
      "key": "hostPort",
      "label": "Local service hostPort"
    },
    {
      "key": "postgres-db.POSTGRES_URL",
      "label": "Service postgres-db variable"
    }
  ],
  "databaseVariables": [
    {
      "key": "POSTGRES_URL",
      "value": "${postgres-db.POSTGRES_URL}",
      "label": "PostgreSQL private URL"
    }
  ]
}
```
