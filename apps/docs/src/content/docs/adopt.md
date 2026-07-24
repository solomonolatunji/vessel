---
title: No Lock-In
description: Codedock apps are standard Docker containers — they survive removal of the Codedock daemon.
---

Codedock is designed so you never lose access to your applications. Every app and database runs as a standard Docker container with persistent volumes. Removing Codedock leaves your containers running.

## How It Works

- **App containers** — deployed with `--restart unless-stopped` via standard Docker.
- **Database containers** — run on named volumes that persist independently.
- **Traefik reverse proxy** — is managed by Codedock, but your containers keep running if Codedock stops.

## Uninstall Codedock Without Losing Apps

```sh
codedockd uninstall
```

This command:
1. Stops the Codedock daemon container.
2. Removes the systemd service.
3. Removes `codedockd` from PATH.
4. **Leaves all your app and database containers running.**

After uninstall, your apps continue serving traffic if you set up your own reverse proxy. The Traefik routing will stop, but your containers are still running on the Codedock Docker network with their assigned ports.

## Adopt Your Containers (After Uninstall)

To take manual control of your containers:

```sh
# List all running containers
docker ps --filter network=codedock-network

# Inspect an app container
docker inspect <container-name>

# View logs
docker logs <container-name>

# Set up your own reverse proxy (nginx example):
# docker run -d --name my-proxy -p 80:80 -p 443:443 ...
# Point it to your app containers on the codedock-network
```

### Databases

Database containers have persistent volumes:

```sh
# List volumes
docker volume ls | grep codedock-db

# Backup a database volume
docker run --rm -v codedock-db-data-<id>:/data -v $(pwd):/backup alpine tar czf /backup/db-backup.tar.gz -C /data .
```

## Migration to Another Platform

Since everything is standard Docker, migrating is straightforward:

1. List all running containers: `docker ps --filter network=codedock-network`
2. For each container, note the image, env vars, and volume mounts.
3. Recreate them on your new platform with the same configuration.

```sh
# Example: recreate a database container manually
docker run -d \
  --name my-postgres \
  --network codedock-network \
  -e POSTGRES_USER=codedock \
  -e POSTGRES_PASSWORD=<password> \
  -e POSTGRES_DB=codedock \
  -v codedock-db-data-<id>:/var/lib/postgresql/data \
  postgres:16-alpine
```

## Backup Before Changes

Before any major operation:

```sh
codedockd backup
```

This creates a timestamped copy of the Codedock database at `/codedock/data/backups/`.
