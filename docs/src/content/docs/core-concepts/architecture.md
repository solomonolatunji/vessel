---
title: Architecture
description: How Vessl coordinates the control plane, Docker, Traefik, BuildKit, services, databases, backups, and migrations.
---

Vessl is a control plane that runs on your server and manages a Docker-based deployment runtime.

## Main Pieces

- The Vessl app stores projects, services, deployments, domains, environment variables, backup metadata, import sources, and users.
- Docker runs application containers, worker containers, database containers, and short-lived helper containers for builds, backups, imports, restores, and updates.
- BuildKit builds source services through Railpack.
- Traefik routes HTTP traffic, handles certificates, serves the control plane hostname, serves generated service hostnames, and routes custom domains.
- The runtime network lets app services and database services talk privately inside Docker.
- `DATA_DIR` stores Vessl state such as static sites, database backups, Postgres TLS assets, system settings, update history, maintenance history, and generated Traefik config.

## Request Flow

For a web service, public traffic reaches Traefik first. Traefik routes the request to the service's active container port. Vessl updates that active port during successful deployments, then reloads Traefik.

For static output, Traefik serves files from Vessl's static site directory.

For workers, there is no public route. Vessl starts the container and checks that the process stays running.

## Build Flow

Git services are cloned from GitHub or a direct Git URL. Vessl passes the selected branch, root directory, install command, build command, start command, runtime mode, internal port, and static output setting into the build and deployment flow.

Source services use BuildKit and Railpack. Docker image services skip the source build and pull the configured image directly.

## Database Flow

Database services run as Docker containers with persistent Docker volumes. Vessl creates connection variables, generated public hostnames when a root domain exists, and Traefik routes for public database access where supported.

PostgreSQL-family services can enable logical replication. Vessl also creates Postgres TLS assets for public database hostnames.

## Backup Flow

Database backups are created locally first. Depending on the selected storage target, Vessl then keeps the local disk file, uploads it to Cloudflare R2, or does both.

Enabled automatic schedules run in the background for daily, weekly, or monthly backups. Manual backups run from the database Backups tab.

## Migration Flow

Vessl supports two migration paths:

- Vessl migration bundles, which export a whole Vessl instance into an encrypted `.vessl` file and restore it into another server.

Vessl bundles move a Vessl instance between servers.
