---
title: CI/CD & Git Integration
description: Set up automatic deployments on push.
---

## Git Providers

Connect GitHub, GitLab, or Bitbucket to enable automatic deployments:

1. Go to **Settings → Git Apps**.
2. Install the Codedock GitHub App, or configure a GitLab / Bitbucket App.
3. Grant repository access to the repos you want to deploy.

## Automatic Deployments

Once connected, every push to the configured branch triggers a new deployment:

1. Codedock receives the webhook.
2. Clones the latest commit.
3. Builds a new container image.
4. Runs a health check on the new container.
5. Swaps traffic to the new container (zero-downtime).
6. Cleans up the old container.

## Manual Deploy

Trigger a deployment from the dashboard or CLI:

```sh
curl -X POST /api/projects/:id/deploy \
  -H "Authorization: Bearer vpt_xxx"
```

## Webhooks

Configure outgoing webhooks for deployment events:

- `deployment.started`
- `deployment.completed`
- `deployment.failed`

Webhooks can notify external services, chat platforms, or your own automation.

## PR Previews

When a pull request is opened against your connected repository, Codedock can spin up an ephemeral preview environment. The webhook triggers a new deployment on the PR branch, and a comment is added to the PR with the preview URL.
