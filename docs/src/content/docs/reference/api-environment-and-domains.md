---
title: Environment and Domains API
description: Environment variable and service domain endpoints with payload and response examples.
---

Environment variables and domains belong to services.

All examples assume:

```bash
export VESSL_URL="https://pilot.example.com"
export VESSL_API_KEY="ap_..."
```

## Set Environment Variable

```txt
POST /api/services/:serviceId/env
```

Required access: `write`

Project scope: service project must be visible to the key.

Payload:

```json
{
  "key": "NODE_ENV",
  "value": "production"
}
```

Example:

```bash
curl -X POST "$VESSL_URL/api/services/svc_web/env" \
  -H "Authorization: Bearer $VESSL_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"key":"NODE_ENV","value":"production"}'
```

Response:

```json
{
  "ok": true
}
```

The same endpoint creates or updates a variable. Environment variable keys must match:

```txt
^[A-Z_][A-Z0-9_]*$
```

## Delete Environment Variable

```txt
DELETE /api/services/:serviceId/env/:envId
```

Required access: `write`

Project scope: service project must be visible to the key.

Example:

```bash
curl -X DELETE "$VESSL_URL/api/services/svc_web/env/env_123" \
  -H "Authorization: Bearer $VESSL_API_KEY"
```

Response:

```json
{
  "ok": true
}
```

## Add Domain

```txt
POST /api/services/:serviceId/domains
```

Required access: `write`

Project scope: service project must be visible to the key.

Workers do not accept custom domains.

Payload:

```json
{
  "hostname": "app.example.com"
}
```

Example:

```bash
curl -X POST "$VESSL_URL/api/services/svc_web/domains" \
  -H "Authorization: Bearer $VESSL_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"hostname":"app.example.com"}'
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

The created domain appears in `GET /api/services/:serviceId/overview`.

## Update Domain

```txt
PATCH /api/services/:serviceId/domains/:domainId
```

Required access: `write`

Payload:

```json
{
  "hostname": "www.example.com"
}
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

## Delete Domain

```txt
DELETE /api/services/:serviceId/domains/:domainId
```

Required access: `write`

Example:

```bash
curl -X DELETE "$VESSL_URL/api/services/svc_web/domains/domain_123" \
  -H "Authorization: Bearer $VESSL_API_KEY"
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

## Apply DNS Provider Record

```txt
POST /api/services/:serviceId/domains/:domainId/dns-records
```

Required access: `write`

Project scope: service project must be visible to the key.

This endpoint uses a DNS provider that has already been connected in System Settings.

Payload:

```json
{
  "providerId": "cloudflare"
}
```

Supported provider IDs:

- `cloudflare`
- `namecheap`
- `spaceship`

Example:

```bash
curl -X POST "$VESSL_URL/api/services/svc_web/domains/domain_123/dns-records" \
  -H "Authorization: Bearer $VESSL_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"providerId":"cloudflare"}'
```

Response:

```json
{
  "ok": true,
  "result": {
    "provider": "cloudflare",
    "providerName": "Cloudflare",
    "action": "created",
    "hostname": "app.example.com",
    "recordType": "A",
    "host": "app",
    "zone": "example.com",
    "targetIp": "203.0.113.10"
  },
  "domain": {
    "id": "domain_123",
    "serviceId": "svc_web",
    "hostname": "app.example.com",
    "status": "active",
    "createdAt": "2026-06-10T08:41:00.000Z",
    "updatedAt": "2026-06-10T08:42:00.000Z"
  }
}
```

## Read Environment and Domains

Environment variables and domains are returned by service overview:

```txt
GET /api/services/:serviceId/overview
```

Response excerpt:

```json
{
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
      "hostname": "app.example.com",
      "status": "active",
      "createdAt": "2026-06-10T08:41:00.000Z",
      "updatedAt": "2026-06-10T08:42:00.000Z"
    }
  ]
}
```
