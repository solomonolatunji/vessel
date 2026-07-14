---
title: Backups and Updates
description: Keep database backups, host maintenance, and Vessl updates visible.
sidebar:
  order: 3
---

This page is an operations overview. For detailed behavior, use the dedicated backup, R2, maintenance, and update pages.

## Database Backups

Use [Database Backups](/docs/storage-and-backups/database-backups/) to configure manual and automatic backups per database service.

Backups can target disk, R2, or both. PostgreSQL-family, MySQL, MongoDB, and Redis services have backup and restore support. ClickHouse backups are not available yet.

Use [Restore and Download](/docs/storage-and-backups/restore-and-download/) before you rely on a recovery path.

## R2 Storage

Use [R2 Storage](/docs/storage-and-backups/r2-storage/) to connect Cloudflare R2. Once connected, database services can upload backups off the server.

## Data Access

Database browser panels help with day-to-day inspection. See [Data Browser](/docs/databases/data-browser/) for engine support and [Data Imports](/docs/databases/data-imports/) for Postgres, TimescaleDB, Redis, and Railway import data flows.

## System Maintenance

Use [System Maintenance](/docs/operations/system-maintenance/) to watch disk, Docker usage, build artifacts, database backup size, APT cache, and system logs.

## System Updates

Use [System Updates](/docs/operations/system-updates/) to review pending commits, update git installs, and understand image install update behavior.

## Logs to Keep Handy

When you need the server-side view, these commands are the fastest places to start:

```bash
sudo journalctl -u vessl -f
cd /opt/vessl && sudo docker compose logs -f traefik buildkit
```

Treat the dashboard and the host logs as a pair: the dashboard shows what Vessl believes is happening, and the host logs show what the server is doing underneath.
