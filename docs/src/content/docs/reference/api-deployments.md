---
title: Deployments API
description: Deployment API endpoints, payloads, log responses, and stream behavior.
---

Deployment endpoints queue, inspect, abort, and stream deployment work.

All examples assume:

```bash
export VESSL_URL="https://pilot.example.com"
export VESSL_API_KEY="ap_..."
```

## Create Deployment

```txt
POST /api/services/:serviceId/deployments
```

Required access: `write`

Project scope: service project must be visible to the key.

Example:

```bash
curl -X POST "$VESSL_URL/api/services/svc_web/deployments" \
  -H "Authorization: Bearer $VESSL_API_KEY"
```

Response:

```json
{
  "deployment": {
    "id": "dep_123",
    "serviceId": "svc_web",
    "commitSha": null,
    "status": "queued",
    "trigger": "manual",
    "imageTag": null,
    "containerName": null,
    "startedAt": null,
    "finishedAt": null,
    "createdAt": "2026-06-10T08:45:00.000Z"
  }
}
```

## List Service Deployments

```txt
GET /api/services/:serviceId/deployments
```

Required access: `read`

Project scope: service project must be visible to the key.

Example:

```bash
curl "$VESSL_URL/api/services/svc_web/deployments" \
  -H "Authorization: Bearer $VESSL_API_KEY"
```

Response:

```json
{
  "deployments": [
    {
      "id": "dep_123",
      "serviceId": "svc_web",
      "commitSha": "abc1234",
      "status": "running",
      "trigger": "manual",
      "imageTag": "vessl-svc_web-dep_123",
      "containerName": "vessl-svc_web-stable",
      "startedAt": "2026-06-10T08:45:00.000Z",
      "finishedAt": "2026-06-10T08:46:00.000Z",
      "createdAt": "2026-06-10T08:45:00.000Z"
    }
  ]
}
```

## Abort Deployment

```txt
POST /api/deployments/:deploymentId/abort
```

Required access: `write`

Project scope: deployment service project must be visible to the key.

Example:

```bash
curl -X POST "$VESSL_URL/api/deployments/dep_123/abort" \
  -H "Authorization: Bearer $VESSL_API_KEY"
```

Response:

```json
{
  "accepted": true
}
```

If the deployment cannot be aborted, Vessl returns `409`.

## Deployment Logs

```txt
GET /api/deployments/:deploymentId/logs
```

Required access: `read`

Project scope: deployment service project must be visible to the key.

Example:

```bash
curl "$VESSL_URL/api/deployments/dep_123/logs" \
  -H "Authorization: Bearer $VESSL_API_KEY"
```

Response:

```json
{
  "logs": [
    {
      "id": 1,
      "deploymentId": "dep_123",
      "line": "Cloning repository acme/web",
      "stream": "stdout",
      "createdAt": "2026-06-10T08:45:03.000Z"
    },
    {
      "id": 2,
      "deploymentId": "dep_123",
      "line": "Build completed",
      "stream": "stdout",
      "createdAt": "2026-06-10T08:45:55.000Z"
    }
  ]
}
```

## Deployment Log Stream

```txt
GET /api/deployments/:deploymentId/stream
```

Required access: `read`

Project scope: deployment service project must be visible to the key.

This endpoint returns Server-Sent Events.

Example:

```bash
curl -N "$VESSL_URL/api/deployments/dep_123/stream" \
  -H "Authorization: Bearer $VESSL_API_KEY"
```

Events:

```txt
event: snapshot
data: [{"id":1,"deploymentId":"dep_123","line":"Cloning repository acme/web","stream":"stdout","createdAt":"2026-06-10T08:45:03.000Z"}]

event: log
data: {"id":2,"deploymentId":"dep_123","line":"Build completed","stream":"stdout","createdAt":"2026-06-10T08:45:55.000Z"}

event: ping
data: {"t":1781081155000}
```

## Runtime Log Stream

```txt
GET /api/services/:serviceId/runtime-logs/stream
```

Required access: `read`

Project scope: service project must be visible to the key.

This endpoint streams logs from the running container.

Example:

```bash
curl -N "$VESSL_URL/api/services/svc_web/runtime-logs/stream" \
  -H "Authorization: Bearer $VESSL_API_KEY"
```

Events:

```txt
event: snapshot
data: [{"id":1,"line":"Server listening on port 3000","stream":"stdout","createdAt":"2026-06-10T08:47:00.000Z"}]

event: log
data: {"id":2,"line":"GET /health 200","stream":"stdout","createdAt":"2026-06-10T08:47:10.000Z"}

event: status
data: {"ok":true,"closed":true}
```

## Explain Failure

```txt
POST /api/deployments/:deploymentId/explain-failure
```

Required access: `write`

Project scope: deployment service project must be visible to the key.

This endpoint uses the configured AI provider to explain a failed deployment.

Payload:

```json
{
  "providerId": "openai",
  "model": "gpt-5-mini"
}
```

Response:

```json
{
  "explanation": {
    "summary": "The deployment built successfully, but the service did not answer on the configured port.",
    "likelyCause": "The app is listening on port 8080 while the service is configured for port 3000.",
    "suggestedFix": "Update the service internal port to 8080 or change the app to listen on 3000.",
    "commands": [
      "npm start"
    ]
  }
}
```
