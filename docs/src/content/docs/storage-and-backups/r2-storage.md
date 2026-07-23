---
title: R2 Storage
description: Connect Cloudflare R2 so database backups can upload off the server.
---

Codedock can upload database backups to Cloudflare R2. R2 is configured once in System Settings, then selected per database service.

## Required Values

Open System Settings, choose Storage, and enter:

- Cloudflare account ID.
- Bucket name.
- R2 access key ID.
- R2 secret access key.

The default bucket name shown in the UI is `codedock-backups`.

## Create or Verify Bucket

The R2 form includes `Create or verify bucket`. When enabled, Codedock checks whether the bucket exists and creates it if R2 returns not found.

R2 credential errors are surfaced with specific guidance for account ID, access key ID, secret access key, bucket access, and token permissions.

## What Codedock Stores

Database backup uploads use keys like:

```txt
database-backups/{projectId}/{serviceSlug}/{filename}
```

Codedock stores only the public R2 status in the UI: account ID, bucket, endpoint, access key suffix, and timestamps. The secret access key is encrypted in system settings.

## Backup Destinations

Once R2 is connected, database backup settings can use:

- `disk`: keep the backup only on the server.
- `r2`: upload to R2 and remove the local backup file.
- `disk+r2`: keep the local file and upload a copy to R2.

When R2 is not connected, disk is the only available destination.

## Disconnecting R2

Disconnecting R2 removes the stored R2 connection and disables future R2 uploads. Existing backup records stay in Codedock.

If a backup exists only in R2 and you disconnect R2, Codedock will not be able to download or restore that backup until R2 is reconnected with access to the same bucket and object.
