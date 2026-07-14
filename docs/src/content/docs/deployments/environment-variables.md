---
title: Environment Variables
description: Add variables, use database suggestions, and understand how Vessl keeps secrets scoped to services.
---

Environment variables belong to a service. Add them during service creation or from the Variables tab after the service exists.

## Adding Variables

You can add variables one at a time or paste multiple `.env` lines. Vessl understands lines like:

```txt
DATABASE_URL=postgres://user:password@host:5432/app
REDIS_URL=redis://default:password@host:6379
```

Lines starting with `#` are ignored. Duplicate keys are replaced by the latest pasted value.

## Database Suggestions

When a project has database services, Vessl suggests variables for app services in that project. Suggestions reference the database service instead of freezing one old URL.

Example:

```txt
DATABASE_URL=${postgres-db.POSTGRES_URL}
```

That makes app configuration easier to keep correct after Railway imports, database recreation, or Vessl migration bundle restores.

## `.env.example` Suggestions

For GitHub repository services, Vessl can inspect `.env.example` variables from the selected branch and root directory. These suggestions help you remember required app variables without copying real secrets into source control.

## Railway Imports

Railway imports can copy service variables. The default import option excludes `RAILWAY_*` variables because they usually describe Railway's runtime, not the app's portable configuration.

After a Railway import, Vessl syncs database variables for recreated database services so app services point at the Vessl-managed database URLs.

## Secret Handling

Values are masked in the UI. Treat Vessl service variables as runtime secrets and avoid using them as documentation. Keep source-controlled examples in `.env.example` without secret values.

## Good Patterns

- Keep variables scoped to the service that needs them.
- Use database suggestions instead of hand-copying generated URLs when possible.
- Keep host-level settings, such as GitHub App credentials or update commands, in system settings or server env.
- Recheck variables after migrations and imports before deploying public traffic.
