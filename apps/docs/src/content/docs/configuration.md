---
title: Configuration
description: Server settings, environment variables, custom domains, notifications, and email configuration.
---

Codedock offers extensive configuration options through the dashboard and environment variables.

## Server Settings

Access server settings from **Settings → Server Settings** in the dashboard. Only instance admins can modify these.

### Traefik Reverse Proxy

- **Wildcard Domain**: Set the base domain for all services (e.g. `codedock.example.com`).
- **Let's Encrypt Email**: Email for SSL certificate notifications.
- **HTTP to HTTPS Redirect**: Auto-redirect all HTTP traffic to HTTPS.
- **Custom Port Bindings**: Map service ports to custom host ports.

### DNS

- **Custom DNS Resolvers**: Override system DNS for container networking.
- **Network**: Codedock containers run on a dedicated Docker network.

### System

- **Data Directory**: Location for SQLite database, vault keys, and container volumes.
- **Static Directory**: Path to the built dashboard frontend.
- **Port**: The HTTP port for the Codedock daemon (default: `8080`).

## Environment Variables

Set at the server level via `.env`:

```sh
PORT=8080
CODEDOCK_DATA_DIR=data
CODEDOCK_STATIC_DIR=dashboard/dist
CODEDOCK_TLS_EMAIL=admin@example.com
CODEDOCK_MAGIC_DOMAIN=traefik.me # Magic DNS domain (options: sslip.io, traefik.me, nip.io)
```

## Notifications

Codedock sends deployment, backup, and system notifications through configurable channels.

### Supported Channels

| Channel             | Configuration                          |
| ------------------- | -------------------------------------- |
| **Discord**         | Webhook URL                            |
| **Slack**           | Webhook URL                            |
| **Telegram**        | Bot token + Chat ID                    |
| **Pushover**        | User key + App token                   |
| **Email (SMTP)**    | SMTP server, credentials, from address |
| **Resend**          | Resend API key                         |
| **Generic Webhook** | Custom URL + headers                   |

### Setting Up a Channel

1. Go to **Settings → Notifications**.
2. Click **Add Channel**.
3. Select the channel type and enter the required credentials.
4. Test the channel with a sample notification.

### Notification Events

Channels can be configured for specific events:

- Deployment started, completed, failed
- Backup completed, failed
- Backup upload to S3 completed
- System update available
- Backup retention cleanup

### Email Configuration

You can configure SMTP or Resend settings for transactional emails (invites, password resets):

1. Go to **Server Settings → Email**.
2. Configure SMTP server or Resend API key.
3. Verify the configuration with a test email.

## OAuth Providers

Configure external authentication providers:

1. Go to **Settings → OAuth Providers**.
2. Add a provider (GitHub, Google, GitLab, or custom).
3. Enter the Client ID and Client Secret from the provider's OAuth app.
4. Enable the provider for login.

Supported providers:

- GitHub
- Google
- GitLab
- Bitbucket
- Custom OpenID Connect

## Two-Factor Authentication

Enable 2FA on your account for additional security:

1. Go to **Profile → Security**.
2. Click **Setup 2FA**.
3. Scan the QR code with your authenticator app (Authy, Google Authenticator, 1Password).
4. Enter the verification code to confirm.

## AI Integration

Configure an AI provider for deployment diagnostics:

1. Go to **Server Settings → AI**.
2. Select a provider (OpenAI or Anthropic).
3. Enter your API key.
4. When a deployment fails, click **AI Diagnose** to analyze logs and get fix suggestions.

## License Management

For self-hosted instances, manage your license from **Settings → License**:

- Activate a license key
- View current plan, seat limits, and expiration
- Instance admins can update license details

## System Updates

Codedock checks for updates automatically. To manage updates:

1. Go to **Settings → Updates**.
2. Check the current version and available updates.
3. Click **Deploy Update** to upgrade.
4. The dashboard will display update progress.

Updates can also be triggered manually via the CLI:

```sh
curl -fsSL https://get.codedock.run | sh
```
