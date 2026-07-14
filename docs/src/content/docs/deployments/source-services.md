---
title: Source Services
description: Deploy apps from GitHub repositories or direct Git URLs.
---

Source services are services Vessl builds from code. Use them for most web apps, APIs, workers, and static sites.

## Source Options

Vessl supports two source paths:

- `GitHub repository`: browse repositories available to the connected GitHub App, choose a branch, and select a root directory.
- `Git URL`: enter an HTTPS or SSH Git URL manually, then provide the branch and root directory.

Use GitHub repositories when you want repository discovery and push-triggered deployments. Use direct Git URLs when the repository is public, uses SSH, or lives outside the GitHub App flow.

## Repository Settings

Each source service stores:

- Service name.
- Repository or Git URL.
- Branch.
- Root directory.
- Runtime mode.
- Internal port for web services.
- Optional static output directory.
- Optional install, build, and start command overrides.

The root directory matters for monorepos. Choose the folder that contains the app you want Vessl to build. For direct Git URLs, type the root directory manually.

## Build Detection and Overrides

After cloning the repository, Vessl checks for a `Dockerfile` at the service root. When one exists, the deployment builds it with `docker build` and skips Railpack entirely — the Dockerfile controls how the image is built and started. Otherwise Vessl uses Railpack through BuildKit to detect and build the app. Leave install, build, and start commands blank when auto detection is correct.

TanStack Start apps are detected automatically when the service root depends on `@tanstack/react-start` or `@tanstack/start`. If the app has no start command, Vessl supplies the correct production start behavior for Nitro output or for the `dist/client` plus `dist/server` fetch-handler output used by TanStack Start server builds.

The build method can be pinned per service in Settings:

- **Auto** (default) — use the Dockerfile when present, otherwise Railpack.
- **Dockerfile** — always build with the Dockerfile; the deployment fails if it is missing.
- **Railpack** — always build with Railpack, even when a Dockerfile is present.

A Dockerfile in a non-standard location can be selected with the Dockerfile path setting, or with the `VESSL_DOCKERFILE_PATH` service environment variable (the variable wins when both are set). Paths are relative to the service root directory. Install, build, and start command overrides apply only to Railpack builds.

Use command overrides when:

- The repo has a custom install command.
- The build script is not the default script.
- The service needs a custom start command.
- A monorepo package must be launched from a specific path.

The static output setting turns a source service into a static site deployment. Vessl copies the output directory out of the built image and serves it through Traefik instead of running the app server.

## Runtime Mode

`web` services must listen on the configured internal port. Vessl starts the container, checks the port, reloads Traefik, and routes traffic.

`worker` services run without a published port. Vessl checks that the container process stays running.

## GitHub Push Deploys

When GitHub App webhooks are configured, pushes can enqueue deployments for connected services. The webhook URL is:

```txt
https://YOUR_PUBLIC_HOST/api/github/app/webhook
```

Vessl also supports manual deployments from the service Deployments tab. Manual deployments are useful for first deploys, retries, and direct Git URL services.

## Common Failures

- BuildKit is not reachable at the configured `BUILDKIT_HOST`.
- The selected root directory does not contain the app.
- The app listens on a different port than the service internal port.
- The static output folder does not contain `index.html`.
- The GitHub App cannot read the repository.
- A direct SSH Git URL needs host-level SSH credentials.
