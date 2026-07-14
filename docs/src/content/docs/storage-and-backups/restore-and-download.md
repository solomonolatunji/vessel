---
title: Restore and Download
description: Download, delete, and restore database backups from disk or R2.
---

Successful backups can be downloaded, restored, or deleted from the database Backups tab.

## Download

When the local backup file exists, Vessl downloads it directly from disk.

When the backup exists only in R2, Vessl downloads the object from R2 into a temporary file, streams it to the browser, then cleans up the temporary file.

If neither a local file nor an accessible R2 object exists, the backup file cannot be downloaded.

## Restore

Restore starts or recreates the database container, waits for readiness, then loads the dump.

Supported restore paths:

- PostgreSQL and TimescaleDB: `pg_restore --clean --if-exists --no-owner`.
- MySQL: `mysql < backup.sql`.
- MongoDB: `mongorestore --archive --gzip --drop`.
- Redis: replace `/data/dump.rdb` and restart Redis.

ClickHouse restore is not available yet.

## Restore Safety

Restores overwrite target data. Treat restore as destructive:

- Create a fresh manual backup first when possible.
- Confirm you are on the correct database service.
- Expect client interruption during restore.
- Verify app variables after restoring into a migrated or recreated service.

## Delete

Deleting a backup removes the local disk file when present. If the backup has an R2 key and R2 is connected, Vessl also tries to delete the R2 object.

The database backup record is removed from Vessl after deletion.

## R2 Dependency

If the backup only exists in R2, restore and download require an active R2 connection with access to the same bucket and key.

For highest recovery confidence, use `disk+r2` for recent backups and a server-level offsite copy strategy for long-term retention.
