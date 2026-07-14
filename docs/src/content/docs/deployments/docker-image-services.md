---
title: Docker Image Services
description: Run prebuilt Docker images as web services or workers.
---

Docker image services skip the source build. Vessl pulls the configured image and runs it directly on the server.

## When to Use Them

Use Docker image services when:

- CI already builds and publishes images.
- You want to deploy a service from GHCR, Docker Hub, or another registry.
- You are running third-party software.
- You need full control over the image build outside Vessl.

Private images use the host Docker daemon's registry login. Sign in on the server with Docker before deploying private image references.

## Required Settings

Each Docker image service stores:

- Service name.
- Image reference, for example `ghcr.io/org/app:latest`.
- Runtime mode, either `web` or `worker`.
- Internal port for `web` services.
- Environment variables.

Vessl validates the image reference before creating the service.

## Web Images

For a web image, set the internal port the container listens on. Vessl starts a temporary container, probes the port, and only switches traffic after the container is reachable.

If the container starts but the port never responds, check the image logs. Vessl attempts to detect port mismatch hints from container output.

## Worker Images

For a worker image, choose `worker` runtime mode. Vessl starts the container without a public port and checks that the process stays running.

Worker image deployments reload Traefik only to keep routing state in sync with the rest of the project.

## Hot Swaps

For web Docker image services, Vessl uses a temporary container and an active port swap:

1. Pull the image.
2. Start a temporary container on an available host port.
3. Probe the configured internal port.
4. Update the service active port.
5. Reload Traefik.
6. Remove the previous stable container.
7. Rename the temporary container to the stable name.

This keeps traffic on the old container until the new one is ready.
