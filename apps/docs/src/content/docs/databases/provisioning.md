---
title: Database Provisioning
description: Managed database engines with one-click provisioning and connection string injection.
---

Spin up managed databases directly from the Codedock dashboard. Each database runs in its own Docker container with persistent volumes, automatic health checks, and daily backups.

## Supported Engines

### Relational

| Engine      | Version   | Default Port |
| ----------- | --------- | ------------ |
| PostgreSQL  | 16-alpine | 5432         |
| TimescaleDB | latest    | 5432         |
| MySQL       | 8.0       | 3306         |
| MariaDB     | 11        | 3306         |
| ClickHouse  | latest    | 9000         |

### NoSQL

| Engine    | Version  | Default Port |
| --------- | -------- | ------------ |
| MongoDB   | 7.0      | 27017        |
| Redis     | 7-alpine | 6379         |
| Dragonfly | latest   | 6379         |
| KeyDB     | latest   | 6379         |

### Message Brokers

| Engine   | Version  | Default Port |
| -------- | -------- | ------------ |
| Kafka    | latest   | 9092         |
| RabbitMQ | 4-alpine | 5672         |
| NATS     | 2-alpine | 4222         |

### One-Click Deployers

| Service   | Purpose                          | Port |
| --------- | -------------------------------- | ---- |
| NocoDB    | Open-source Airtable alternative | 8080 |
| Plausible | Web analytics                    | 8000 |
| WordPress | CMS                              | 80   |
| Gitea     | Self-hosted Git service          | 3000 |

## Creating a Database

1. Navigate to **Databases** in the sidebar.
2. Click **New Database**.
3. Select an engine from the list.
4. Optionally set a custom name and port.
5. Click **Create**.

Codedock provisions the container, creates a default database and user, and mounts a persistent volume at `/var/lib/data`.

## Connection Strings

Once created, the connection string is automatically injected into every service in the same project:

```
DATABASE_URL=postgresql://codedock:<password>@<service-name>:5432/codedock
TIMESCALE_URL=postgresql://codedock:<password>@<service-name>:5432/codedock
REDIS_URL=redis://<service-name>:6379
MONGO_URL=mongodb://codedock:<password>@<service-name>:27017/codedock
```

You can also find the connection details on the database's detail page in the dashboard.

## Managing Databases

### Start / Stop

Databases can be started and stopped from the dashboard. Stopping a database frees resources while preserving the volume data.

### Configuration

Each database has sensible defaults:

- **Port**: Assigned from the engine's default port range
- **Username**: `codedock`
- **Database name**: `codedock`
- **Data volume**: Persisted at `<data-dir>/databases/<id>/`
