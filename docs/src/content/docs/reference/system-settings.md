---
title: System Settings
description: Reference for Vessl system-wide settings and what they affect.
---

System Settings controls server-wide behavior. Service settings control one service.

## Domains

`Control plane hostname` serves the Vessl dashboard through Traefik.

`Root domain` enables generated service hostnames and generated database public hostnames.

Both can be set during onboarding or later from System Settings.

## DNS

DNS provider settings store credentials for Cloudflare, Namecheap, and Spaceship. Connected providers appear as actions on service custom domains.

Provider automation writes IPv4 `A` records only.

## GitHub

GitHub settings support either:

- GitHub App credentials.
- Host-level `GITHUB_ACCESS_TOKEN`.

The GitHub App path enables repository discovery and push webhooks.

## API Access

API Access creates scoped API keys for programmatic requests to Vessl.

Keys can be read-only or read/write, scoped to all projects or selected projects, and configured to expire after `7`, `30`, or `90` days, or never expire.

API keys use bearer authentication against the same `/api/*` endpoints as the dashboard. See [API Access](/docs/reference/api-access/) for examples and endpoint behavior.

## Storage

R2 settings store Cloudflare account ID, bucket, endpoint, access key suffix, encrypted secret access key, and timestamps.

R2 must be connected before database services can select `r2` or `disk+r2` backup destinations.

## Migration

Migration settings export encrypted `.vessl` bundles. The export requires a passphrase with at least 8 characters.

Bundle import is available during onboarding.

## Maintenance

Maintenance settings show disk, Docker, Vessl data paths, build artifacts, database backup storage, APT cache, system logs, cleanup actions, and recent history.

Safe cleanup avoids Docker volumes. Volume cleanup is separate and destructive for detached volumes.

## Deployments

`Concurrent deployments` controls the global number of deployment jobs that can run at once.

Allowed values:

```txt
1 through 10
```

Default:

```txt
3
```

Vessl also enforces one active deployment per service.

## Updates

Updates settings compare the running install against the configured repository and branch.

Git installs can fast-forward, install dependencies, build, prune, and queue a restart.

Image installs can compare commit metadata and run an image update command when configured.
