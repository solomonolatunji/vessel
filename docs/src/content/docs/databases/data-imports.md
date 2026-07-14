---
title: Data Imports
description: Import PostgreSQL-compatible or Redis data from URLs or Railway import sources.
---

Vessl data imports move data into an existing database service. Deploy the target database before importing data.

## Supported Imports

| Target engine | Source options | Notes |
| --- | --- | --- |
| PostgreSQL | Postgres URL, Railway | Uses `pg_dump` custom format and `pg_restore`. |
| TimescaleDB | Postgres URL, Railway | Requires compatible PostgreSQL and TimescaleDB versions. |
| Redis | Redis URL, Railway | Uses `redis-cli --rdb` and loads the RDB into Redis. |

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

Vessl redacts passwords from error messages when a dump command fails.

## Import from Railway

Railway data import uses the service import source metadata saved during Railway project import. Vessl fetches Railway service variables again with the Railway API token and searches for public database URLs.

For PostgreSQL, Vessl checks keys such as `POSTGRES_PUBLIC_URL`, `DATABASE_PUBLIC_URL`, `POSTGRES_URL`, and `DATABASE_URL`, then falls back to `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, and `PGPASSWORD`.

For Redis, Vessl checks keys such as `REDIS_PUBLIC_URL`, `DATABASE_PUBLIC_URL`, `REDIS_TLS_URL`, `REDIS_URL`, and `DATABASE_URL`, then falls back to Redis host, port, user, and password parts.

## Railway Internal URLs

Railway internal URLs ending in `.railway.internal` cannot be reached from your Vessl server. Enable public networking on the Railway database or use the URL import option with a public URL.

## TimescaleDB Compatibility

TimescaleDB imports need source and target extension compatibility. Vessl tries to read the source PostgreSQL major version and TimescaleDB extension version. If the target is older or the TimescaleDB extension version differs, import fails with a corrective message.

## Import History

Vessl records data import jobs with status, source label, source variable key, dump size, checksum, error, and timestamps. Use the history to confirm whether an automated Railway import actually copied data or only recreated the service.
