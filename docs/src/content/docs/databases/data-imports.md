---
title: Data Imports
description: Import PostgreSQL-compatible or Redis data from URLs.
---

Codedock data imports move data into an existing database service. Deploy the target database before importing data.

## Supported Imports

| Target engine | Source options | Notes |
| --- | --- | --- |
| PostgreSQL | Postgres URL | Uses `pg_dump` custom format and `pg_restore`. |
| TimescaleDB | Postgres URL | Requires compatible PostgreSQL and TimescaleDB versions. |
| Redis | Redis URL | Uses `redis-cli --rdb` and loads the RDB into Redis. |

MySQL, MongoDB, and ClickHouse data imports are not available yet.

## Import from URL

Use a direct source URL when you already have a public connection string.

Accepted Postgres schemes:

```txt
postgres://
postgresql://
```

Accepted Redis schemes:

```txt
redis://
rediss://
```

Codedock redacts passwords from error messages when a dump command fails.

For PostgreSQL, Codedock checks keys such as `POSTGRES_PUBLIC_URL`, `DATABASE_PUBLIC_URL`, `POSTGRES_URL`, and `DATABASE_URL`, then falls back to `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, and `PGPASSWORD`.

For Redis, Codedock checks keys such as `REDIS_PUBLIC_URL`, `DATABASE_PUBLIC_URL`, `REDIS_TLS_URL`, `REDIS_URL`, and `DATABASE_URL`, then falls back to Redis host, port, user, and password parts.

## TimescaleDB Compatibility

TimescaleDB imports need source and target extension compatibility. Codedock tries to read the source PostgreSQL major version and TimescaleDB extension version. If the target is older or the TimescaleDB extension version differs, import fails with a corrective message.

## Import History

Codedock records data import jobs with status, source label, source variable key, dump size, checksum, error, and timestamps.
