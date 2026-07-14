---
title: Railway Import
description: Import Railway projects into Vessl, including services, variables, databases, custom domains, and optional data imports.
---

Railway import translates a Railway project into Vessl projects and services. It is for moving workloads from Railway into your own server.

## What You Need

- A Railway Personal API token.
- A running Vessl server.
- Enough disk for cloned apps, Docker images, database volumes, and optional database data imports.
- Public networking enabled for Railway databases if you want Vessl to import data automatically.

Railway internal hosts ending in `.railway.internal` are only reachable inside Railway. Vessl rejects those for Postgres and Redis data import and asks for a public URL instead.

## Import Flow

The import modal walks through four steps:

1. Enter Railway API token.
2. Select a Railway project.
3. Choose environment and migration options.
4. Run the import.

Vessl reads Railway environments, service instances, service variables, source metadata, repo triggers, command overrides, root directory, custom domains, and detected service images when Railway exposes them.

## Supported Railway Service Shapes

Vessl can recreate:

- Git services.
- Docker image services.
- Database services.

Database detection is based on Railway image and service metadata. Vessl recognizes PostgreSQL, TimescaleDB, MySQL/MariaDB, Redis, MongoDB, and ClickHouse shapes.

Services that are not Git services, Docker image services, or recognized databases are marked unsupported instead of being guessed.

## Migration Options

`Exclude RAILWAY_* variables` is enabled by default. Keep it on unless your app intentionally uses those variables outside Railway.

`Recreate database engines` creates Vessl database services for Railway databases. Vessl creates fresh local credentials and new Vessl-managed connection variables.

`Auto-deploy services` queues deployments after the import. Database services deploy first, then app services deploy after database work settles.

`Import database data` is available only when database recreation and auto-deploy are enabled. Vessl needs deployed target databases before it can import data.

## Database Data Import

Railway data import currently runs for:

- PostgreSQL-compatible services, including TimescaleDB.
- Redis services.

For PostgreSQL, Vessl finds a public Postgres URL from Railway variables, creates a `pg_dump` custom dump, then restores it into the Vessl database.

For Redis, Vessl finds a public Redis URL, creates an RDB dump with `redis-cli --rdb`, then loads it into the Vessl Redis container.

MySQL, MongoDB, and ClickHouse services can be recreated by the Railway import, but automatic Railway data import is not available for them yet.

## TimescaleDB Notes

For TimescaleDB imports, Vessl tries to detect the source PostgreSQL major version and TimescaleDB extension version. The target must be compatible before restore.

If extension versions differ, Vessl explains the mismatch and may suggest resetting the local target volume, redeploying with the pinned source image, then retrying the import.

## Custom Domains

Vessl imports Railway custom domains on a best-effort basis. It uses Railway domain metadata to find target ports when available and falls back to `PORT` variables when it must.

After import, verify DNS. Domains that were pointed at Railway need to be pointed at your Vessl server.

## After Import

Check the imported project before sending traffic:

- Confirm all services exist with the expected type.
- Review variables and remove Railway-only values.
- Deploy databases first if auto-deploy was not enabled.
- Deploy app services.
- Verify generated and custom domains.
- Run manual database backups after the first successful import.
