---
title: Database Backups
description: Configure manual and automatic database backups to disk, R2, or both.
---

Database backups are configured per database service from the Backups tab.

## Supported Engines

| Engine | Backup tool | Format |
| --- | --- | --- |
| PostgreSQL | `pg_dump -Fc` | `pg_dump custom` |
| TimescaleDB | `pg_dump -Fc` | `pg_dump custom` |
| MySQL | `mysqldump` | SQL |
| MongoDB | `mongodump --archive --gzip` | gzip archive |
| Redis | `redis-cli --rdb` | RDB |

ClickHouse backups are not available yet.

## Backup Destinations

Codedock creates a local backup file first. Then it follows the selected destination:

- `disk`: keep the local file.
- `r2`: upload to R2, then remove the local file.
- `disk+r2`: keep the local file and upload a copy to R2.

If `disk+r2` is selected and the R2 upload fails after the local backup exists, Codedock marks the backup succeeded on disk and records the R2 error.

If `r2` is selected and upload fails, the backup fails.

## Local Disk Backups

Disk backups live under:

```txt
DATA_DIR/backups/{serviceId}
```

The filename includes timestamp, service slug, engine, and backup ID.

Local disk backups are fast and easy to restore, but they are still on the same server. Use R2 or `disk+r2` when you need off-server recovery.

## Automatic Schedules

Automatic backups are off by default for database services unless you enable schedules during onboarding. Daily, weekly, and monthly schedules can be toggled individually per database service from the Backups tab.

Codedock runs:

| Trigger | Interval | Retention |
| --- | ---: | ---: |
| Daily | 24 hours | 6 days |
| Weekly | 7 days | 31 days |
| Monthly | 30 days | 90 days |

The scheduler starts after the control plane has been running for about 60 seconds, then checks hourly.

## Manual Backups

Click `Create backup` from the database Backups tab to create a manual backup immediately. Manual backups use the selected destination unless you choose a specific destination for that run.

Use manual backups:

- Before imports.
- Before application migrations.
- Before database redeploys.
- Before deleting old services or volumes.

## Backup Records

Each backup records engine, status, trigger, storage target, format, local path, R2 key, size, checksum, error, and timestamps.

Checksums are used when restoring or validating bundle database dumps.
