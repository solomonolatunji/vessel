---
title: Public Access and TLS
description: Configure database public hostnames, Postgres logical replication, and TLS-related behavior.
---

Database services can be private to the server or exposed through generated public hostnames.

## Public Hostnames

When a root domain is configured, Vessl generates database hostnames from the service name. The service settings show the hostname and connection target.

Example:

```txt
postgres-db.pilot.example.com:5432
```

Use the generated host for tools or clients that need to connect from outside the Docker runtime network.

## DNS

Public database hostnames depend on the wildcard root domain resolving to the server. If the wildcard DNS record is missing or still propagating, the hostname will not work from outside the server.

For custom public domains, add the hostname to the service and point an `A` record at the server IP.

## PostgreSQL TLS Assets

For PostgreSQL-family services, Vessl prepares TLS assets used by public database access. Those assets are stored under the Vessl data directory and are included in Vessl migration bundles when present.

After an Vessl bundle import, check restored database hostnames and redeploy services if you need to refresh runtime containers with the restored assets.

## Logical Replication

PostgreSQL-family database services can enable logical replication. When enabled, Vessl deploys Postgres with:

```txt
wal_level=logical
max_replication_slots=10
max_wal_senders=10
```

Use this when downstream tools need logical replication or change data capture. Leave it disabled if you do not need it.

## Security Notes

- Prefer private project-local database URLs for app services on the same server.
- Use public database hostnames only when external access is needed.
- Keep generated credentials in service variables.
- Rotate credentials after sharing public access with temporary tools.
- Confirm firewall and DNS behavior before assuming a public hostname is reachable.
