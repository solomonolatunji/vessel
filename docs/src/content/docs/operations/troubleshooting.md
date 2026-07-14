---
title: Troubleshooting
description: Common deployment, DNS, database, backup, Railway import, and update issues.
---

Start with the page that owns the failing workflow, then use server logs for the runtime view.

## Useful Server Logs

```bash
sudo journalctl -u vessl -f
cd /opt/vessl && sudo docker compose logs -f traefik buildkit
```

Deployment logs show Vessl's job-level view. Runtime logs show container output. Server logs show the underlying control plane, Traefik, and BuildKit behavior.

## BuildKit Is Unavailable

Vessl source builds need BuildKit at the configured `BUILDKIT_HOST`, usually `tcp://127.0.0.1:1234`.

Check:

- The `deploy-buildkit` service is running.
- The host port is not blocked locally.
- Docker is running.
- The BuildKit container is healthy.

## App Starts but Does Not Respond

If a web deployment starts but fails the port probe:

- Confirm the service internal port matches the app's listening port.
- Check runtime logs for a different port.
- Make sure the app binds `0.0.0.0`, not only `127.0.0.1`.
- Check that the service is not a worker or static site by mistake.

## Static Output Fails

Static site deployments require an output directory containing `index.html`.

Check:

- The build command creates the expected folder.
- The `Static output` setting points at the correct path inside the image.
- The app is actually static. Server-rendered apps should run as web services.

## DNS Is Pending

If a domain remains pending:

- Confirm the `A` record points to the server public IPv4.
- Confirm wildcard DNS for generated service hostnames.
- Wait for DNS propagation.
- Confirm Traefik is running and ports `80` and `443` are open.
- Click refresh and verify in the service Domains tab.

## Railway Data Import Fails

Common causes:

- Railway returned only a `.railway.internal` database URL.
- Public networking is not enabled on the Railway database.
- The target Vessl database is not deployed.
- TimescaleDB extension versions differ.
- The target PostgreSQL major version is older than the source.

Enable a public Railway database URL or use the direct URL import option.

## Backups Fail

For database backups:

- Deploy the database service first.
- Confirm the engine supports backups.
- Confirm R2 is connected before choosing `r2` or `disk+r2`.
- Check disk space for the local backup file.
- Check database credentials in service variables.

For `disk+r2`, a local backup can still succeed even when R2 upload fails. The backup record will show the R2 error.

## Updates Fail

System Updates will not apply when the checkout is dirty, diverged, or cannot fast-forward.

For image installs, one-click update requires commit metadata and an image update command. Otherwise, run the shown Docker Compose command from the server.
