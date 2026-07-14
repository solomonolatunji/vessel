---
title: Domains
description: Configure generated service hostnames and custom domains for Vessl services.
sidebar:
  order: 2
---

Vessl uses Traefik to route traffic and manage certificates once DNS resolves to the server.

## Control Plane Domain

The control plane domain serves Vessl itself:

```txt
A     pilot.example.com     YOUR_SERVER_IPV4
AAAA  pilot.example.com     YOUR_SERVER_IPV6
```

Set or update this hostname during onboarding or in system settings.

## Wildcard Root Domain

Set a wildcard root domain to generate service hostnames automatically.

DNS record:

```txt
A     *.pilot.example.com     YOUR_SERVER_IPV4
AAAA  *.pilot.example.com     YOUR_SERVER_IPV6
```

With this configured, services can receive hostnames like:

```txt
api.pilot.example.com
web.pilot.example.com
worker.pilot.example.com
```

Database services can also receive generated public hostnames when database public access is enabled.

## Custom Domains

For a custom domain, point DNS at the same server:

```txt
A     app.example.com     YOUR_SERVER_IPV4
AAAA  app.example.com     YOUR_SERVER_IPV6
```

Then add the domain to the service in Vessl. Traefik handles routing and certificates after DNS resolves.

## Service Domain Tab

The service Domains tab shows:

- The hostname.
- Whether the domain is active or pending.
- The target `A` record.
- The server IP to point at.
- DNS provider actions when a provider is connected.
- Refresh and verify action for propagation checks.

Local loopback domains such as `localhost` do not need public DNS records.

## Traefik Behavior

Traefik serves the Vessl dashboard, app services, static sites, and custom domains. Vessl rewrites and reloads Traefik configuration when domain settings change.

If Traefik reload fails, Vessl surfaces the reload detail in the deployment or settings flow. Check Traefik logs from the server when DNS is correct but routing still fails.

## DNS Checklist

- The hostname resolves to the server.
- Ports `80` and `443` are reachable.
- No other process is bound to those public ports.
- Traefik is running through the Vessl Docker Compose stack.
- The service is deployed and active.
- The service internal port matches the app container's listening port.

## Related Pages

- [DNS Providers](/docs/operations/dns-providers/) for Cloudflare, Namecheap, and Spaceship automation.
- [Public Access and TLS](/docs/databases/public-access-and-tls/) for database hostnames.
