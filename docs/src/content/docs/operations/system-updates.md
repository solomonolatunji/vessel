---
title: System Updates
description: Check, review, and apply Vessl updates for git and image installs.
---

System Updates compares the running Vessl install with the configured Vessl repository and branch.

## Git Installs

For installs running from a git checkout, Vessl checks the current commit against the remote branch.

Statuses:

- `current`: installed commit matches GitHub.
- `available`: remote branch has commits that can be fast-forwarded.
- `diverged`: local checkout and remote branch do not fast-forward cleanly.
- `unknown`: Vessl could not check status.

When an update is available and the checkout is clean, Vessl can:

1. Fetch the target commit.
2. Fast-forward merge.
3. Run `npm ci --include=dev`.
4. Build.
5. Prune dev dependencies.
6. Queue a restart when `VESSL_UPDATE_RESTART_CMD` is configured.

If the checkout has local changes, Vessl refuses to update automatically.

## Image Installs

For Docker image installs, Vessl compares the image commit metadata with GitHub. Images need `VESSL_COMMIT_SHA` to make that comparison useful.

If image self-update is configured, Vessl can pull the latest image and replace the running app container.

If one-click image update is not configured, Vessl shows the server command:

```bash
cd /opt/vessl && sudo docker compose pull vessl && sudo docker compose up -d vessl
```

## Pending Commits

When updates are available, the Updates tab lists the commits that will be applied. Review them before updating, especially on a production server.

## Restart Behavior

Updates that replace the running process may briefly interrupt the dashboard. The UI refreshes after a restart is queued.

If Vessl builds the update but cannot queue a restart, restart the service manually after reviewing logs.

## When Manual Update Is Required

Manual update is required when:

- The local checkout has real local changes.
- The checkout diverged from GitHub.
- The image commit is not an ancestor of the remote branch.
- The image lacks commit metadata.
- Image self-update command is not configured.
