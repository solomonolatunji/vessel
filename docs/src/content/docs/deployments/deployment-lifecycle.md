---
title: Deployment Lifecycle
description: Understand queueing, concurrency, aborts, hot swaps, status changes, and deployment logs.
---

Vessl treats deployments as jobs. A service can have many deployments over time, but only one deployment per service can build at once.

## Deployment Triggers

Deployments can be created by:

- Clicking `Deploy` in the service Deployments tab.
- GitHub push webhooks for GitHub-connected services.
- Vessl migration import follow-up deployment queueing for restored active app services.

Every deployment records its trigger.

## Queue and Concurrency

System Settings includes `Concurrent deployments`. The default is `3`; the allowed range is `1` through `10`.

Vessl uses that global number as the deployment slot limit. It also prevents two deployments for the same service from running at the same time.

If all slots are full, new deployments stay `queued` until a slot opens.

## Statuses

- `queued`: waiting for an available global slot.
- `building`: work has started.
- `running`: the deployment is the current live deployment for the service.
- `superseded`: a newer deployment became live.
- `aborted`: a queued or building deployment was stopped.
- `failed`: the deploy could not complete.

Older `running` deployments are marked `superseded` when a newer deployment succeeds.

## Aborting Deployments

You can abort queued or building deployments from the Deployments tab.

For queued deployments, Vessl marks the job aborted before it starts.

For building deployments, Vessl marks the job aborted, asks the active process to stop, then escalates to a force kill if needed. Any temporary containers are cleaned up by the deployment flow.

## Web Hot Swaps

For source-built web services and Docker image web services, Vessl keeps the old container live while the new container proves it can answer on the configured port.

The hot swap flow is:

1. Build or pull the image.
2. Start a temporary container.
3. Probe the configured internal port.
4. Store the new active port.
5. Reload Traefik.
6. Remove the old stable container.
7. Mark the new deployment running.

If the new container does not become reachable, Vessl keeps the previous service state instead of switching traffic.

## Workers and Static Sites

Worker deployments check that the process stays running. There is no HTTP port probe.

Static site deployments export files from the built image into the static site directory. Vessl verifies that the output contains `index.html` before marking the deployment running.

## Database Deployments

Database deployments are different from app hot swaps. Vessl runs database containers with persistent Docker volumes. Deploying a database service ensures the image, volume, runtime network, public hostname route, Postgres TLS assets, and logical replication settings are in place.

Database container replacement can interrupt connections, so schedule database deploys with more care than app deploys.

## Logs

Deployment logs show clone, build, image pull, runtime, Traefik, and health check output. Runtime logs show output from the running container after deployment.

When a deployment fails, start with the deployment log. If the service started but did not answer on the configured port, compare the app logs with the service internal port setting.
