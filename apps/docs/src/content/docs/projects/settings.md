---
title: Project Settings
description: Configure project tokens, webhooks, and members.
---

### Webhooks

Configure outgoing webhooks for project-level events:

1. Go to **Project Settings → Webhooks**.
2. Click **Add Webhook**.
3. Enter the URL and select events to listen for.
4. Optionally, add a secret for HMAC verification.

### Project Tokens

Generate API tokens scoped to a specific project:

1. Go to **Project Settings → Tokens**.
2. Click **Generate Token**.
3. Select permissions (deploy, read logs, manage variables, etc.).
4. Copy the token — it won't be shown again.

### Members

Invite team members to collaborate on a project:

1. Go to **Project Settings → Members**.
2. Click **Add Member**.
3. Enter their email and select a role.
4. They'll receive an invitation.

## Deleting a Project

1. Open the project's settings.
2. Scroll to the bottom and click **Delete Project**.
3. Confirm the deletion.

This removes all services, databases, volumes, and configurations for the project.
