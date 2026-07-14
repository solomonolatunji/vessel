---
title: API Access
description: Authenticate scripts and CI jobs with scoped Vessl API keys.
---

API keys let scripts, CI jobs, and local tools call Vessl without a browser session.

API keys use the same `/api/*` endpoints as the dashboard. Browser requests authenticate with the session cookie. Programmatic requests authenticate with a bearer token.

## Create a Key

Open `System Settings` -> `API Access`, then click `Create API key`.

Choose:

- `Name`: a label for the key, such as `CI deploys`.
- `Access`: `Read` or `Read and write`.
- `Projects`: all projects or selected projects.
- `Expiration`: `7 days`, `30 days`, `90 days`, or `No expiration`.

After creation, Vessl shows the full key once. Copy it before closing the dialog. Later, the API Access page only shows metadata such as prefix, scope, expiration, and last-used time.

## Authentication

Use bearer authentication:

```bash
export VESSL_URL="https://pilot.example.com"
export VESSL_API_KEY="ap_..."

curl "$VESSL_URL/api/projects" \
  -H "Authorization: Bearer $VESSL_API_KEY"
```

Vessl also accepts `X-API-Key`:

```bash
curl "$VESSL_URL/api/projects" \
  -H "X-API-Key: $VESSL_API_KEY"
```

## Access Levels

`Read` keys can call `GET` endpoints.

`Read and write` keys can call `GET`, `POST`, `PATCH`, and `DELETE` endpoints.

Read-only keys receive `403` for write actions.

## Project Scopes

`All projects` keys can access every project according to their access level.

`Specific projects` keys can only access the selected projects. For service routes, Vessl resolves the service to its project before allowing the request. For deployment routes, Vessl resolves the deployment to its service, then to the service project.

Selected-project keys cannot create new projects because new projects are outside their scope. Use an all-project write key for project creation.

## Session-Only Endpoints

Some endpoints always require a browser session and cannot be called with an API key:

- `/api/system/*`
- API key creation and revocation.
- GitHub repository browsing endpoints such as `/api/github/repos`, `/api/github/branches`, and `/api/github/directories`.
- Railway import endpoints under `/api/integrations/*`.
- Onboarding, login, logout, and setup flows.

GitHub webhooks and health checks keep their own public or signature-based behavior.

## Endpoint Pages

- [Projects API](/docs/reference/api-projects/)
- [Services API](/docs/reference/api-services/)
- [Deployments API](/docs/reference/api-deployments/)
- [Databases API](/docs/reference/api-databases/)
- [Environment and Domains API](/docs/reference/api-environment-and-domains/)

## Errors

Common responses:

- `401 Authentication required`: no session cookie or API key was provided.
- `401 Invalid or expired API key`: the key is wrong, revoked, or expired.
- `403 This API key is read-only`: the key tried to call a write endpoint.
- `403 This API key cannot access that project`: the key is scoped to different projects.
- `403 This endpoint requires a browser session`: the endpoint does not allow API keys.

Error response:

```json
{
  "error": "This API key is read-only"
}
```

## Revocation

Revoke keys from `System Settings` -> `API Access`.

Revocation takes effect immediately. Existing scripts using the key will receive `401 Invalid or expired API key`.
