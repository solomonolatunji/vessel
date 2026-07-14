---
title: Browser Onboarding
description: Finish the first-run setup after Vessl is installed on your server.
sidebar:
  order: 2
---

Onboarding connects the running control plane to the domains and source control it needs for daily use.

## Choose the Owner

Create the first owner account in the browser. This account controls system settings, deployment configuration, GitHub connection details, and maintenance panels.

Keep the owner account separate from any app users. It can change server-wide routing, backup, migration, and update settings.

## Set the Control Plane Domain

You can serve Vessl itself through Traefik by choosing a dashboard domain such as:

```txt
pilot.example.com
```

Point the hostname at the server:

```txt
A     pilot.example.com     YOUR_SERVER_IPV4
AAAA  pilot.example.com     YOUR_SERVER_IPV6
```

Vessl keeps the raw `http://IP:8080` URL available as a fallback.

## Configure GitHub

Vessl works best with a GitHub App. Create one with these repository permissions:

- `Contents: Read`
- `Metadata: Read`

Subscribe it to the `Push` webhook event and point the webhook URL at:

```txt
https://YOUR_PUBLIC_HOST/api/github/app/webhook
```

Enter the app details during onboarding or later in system settings.

If you do not want to use a GitHub App, Vessl can also read a `GITHUB_ACCESS_TOKEN` from the server environment. The GitHub App path is better for repository selection and push webhooks.

## Add the Wildcard Root Domain

Set a wildcard root domain when you want Vessl to generate service hostnames automatically:

```txt
*.pilot.example.com
```

With that root, a service named `api` can receive a generated hostname like:

```txt
api.pilot.example.com
```

You can still add public custom domains to individual services. The wildcard root just gives every eligible web service and database service a predictable generated hostname.

## Connect Backup Storage

Automatic backups are off by default. In the Backups step, enable daily, weekly, or monthly schedules individually when new database services should create backups automatically. You can change these schedules later per database service from its Backups tab.

If you have a Cloudflare R2 bucket or token ready, connect R2 in the same step. Vessl can then save database backups to disk, R2, or both.

The default backup destination is:

- `disk` when R2 is not connected.
- `disk+r2` when R2 is connected.

You can change the destination per database service later.

## Import an Vessl Bundle

If you are moving from another Vessl server, use the migration import option during onboarding. Choose the `.vessl` file and enter the passphrase from the source server.

The import restores projects, services, users, domains, environment variables, static output, backup records, backup files, database dumps, Traefik config, and system settings. Vessl clears auth sessions after import, then sends you through the success path so the restored owner account can sign in cleanly.

## Restart Onboarding

After the first run, System Settings includes a `Restart onboarding` action. Use it when you need to walk through domain, GitHub, R2, or migration setup again without reinstalling the server.
