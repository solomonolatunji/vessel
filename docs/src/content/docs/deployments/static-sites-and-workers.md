---
title: Static Sites and Workers
description: Configure static output services and background worker deployments.
---

Static sites and workers are both source or Docker image services with a different runtime shape than a normal web app.

## Static Sites

A static site is a web service with a `Static output` directory. Vessl builds the source service, creates an image, copies the static output directory from that image into `DATA_DIR/static-sites/{serviceId}`, and serves the files through Traefik.

Set `Static output` to the directory created by your build, for example:

- `dist`
- `build`
- `.output/public`
- `apps/web/dist`

If the folder does not contain `index.html`, Vessl treats it as a failed static deployment. That usually means the app is server-rendered or the output directory is wrong.

## Static Output Detection

Leave `Static output` blank when the app should run as a server. Set it only when the deployment should serve files directly.

For TanStack Start SSR apps, leave `Static output` blank. Vessl detects TanStack Start source services and starts the server-rendered app instead of exporting `dist/client` as a static site.

For custom commands, Vessl expects the static output path to match the output inside the built image. For auto-detected builds, use the framework's output folder.

## Workers

Workers run background processes without public traffic. Use workers for queues, schedulers, job processors, event consumers, or any process that should stay alive but not listen behind Traefik.

Worker services:

- Do not need an internal port.
- Do not receive generated public web routes.
- Are considered healthy when the container stays running.
- Still have deployments, runtime logs, environment variables, and settings.

## Choosing the Right Mode

Use `web` when users, APIs, or webhooks need to reach the service.

Use `worker` when the process should run privately in the background.

Use `web` plus `Static output` when the result is static files, not a long-running server process.
