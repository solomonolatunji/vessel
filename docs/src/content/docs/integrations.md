---
title: Integrations
description: Connect Codedock with Git providers, OAuth, Analytics, and external services.
---

Codedock integrates with popular development tools and services.

## Git Providers

Connect GitHub or GitLab for automatic deployments from repository pushes.

### GitHub App

1. Go to **Settings → Git Apps → GitHub**.
2. Click **Create GitHub App** or **Configure**.
3. Follow the GitHub App manifest flow to install the app in your organization.
4. Select the repositories to grant access to.

### GitLab App

1. Go to **Settings → Git Apps → GitLab**.
2. Enter your GitLab instance URL and application credentials.
3. Configure the webhook URL pointing to your Codedock instance.

### Bitbucket App

1. Go to **Settings → Git Apps → Bitbucket**.
2. Follow the OAuth setup flow.
3. Grant repository access permissions.

### Repository Status

Check connected providers and their status:

1. Go to **Sources** in the sidebar.
2. View all connected Git providers and their sync status.
3. Disconnect providers from the same view.

## OAuth Authentication

Configure OAuth providers for login:

1. Go to **Settings → OAuth Providers**.
2. Click **Add Provider**.
3. Select the provider type (GitHub, Google, GitLab, custom).
4. Enter the **Client ID** and **Client Secret** from the provider.
5. Set the **Redirect URI** to your Codedock OAuth callback URL.

### Enabling Providers

After adding a provider, enable it from the same page. The login screen will show the provider's button.

### Custom OpenID Connect

For enterprise SSO, configure a custom OpenID Connect provider:

1. Select **OpenID Connect** as the provider type.
2. Enter the issuer URL, client ID, and client secret.
3. Map the user claims (name, email, avatar).

## Analytics & Telemetry

Integrate analytics to monitor platform usage and performance:

### PostHog

Codedock supports PostHog for product analytics and telemetry.

1. Configure your PostHog integration in `dashboard/src/integrations/`.
2. The platform automatically tracks page views, user sessions, and key interactions using `@posthog/react` and `posthog-js`.

## Outgoing Webhooks

Codedock can send webhook notifications to external services when events occur:

### Event Types

- `deployment.started` — A deployment begins
- `deployment.completed` — A deployment succeeds
- `deployment.failed` — A deployment fails
- `backup.completed` — A backup finishes
- `backup.failed` — A backup fails

### Configuring Webhooks

1. Go to **Project Settings → Webhooks**.
2. Click **Add Webhook**.
3. Enter the payload URL.
4. Select the events to trigger the webhook.
5. Optionally set a secret for HMAC-SHA256 signature verification.

### Payload Format

```json
{
  "event": "deployment.completed",
  "project": { "id": "...", "name": "my-app" },
  "deployment": {
    "id": "...",
    "status": "success",
    "service": "web",
    "commit": "abc123"
  },
  "timestamp": "2026-07-13T12:00:00Z"
}
```
