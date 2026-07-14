---
title: System Maintenance
description: Watch disk pressure, Docker growth, build artifacts, logs, and cleanup actions.
---

System Maintenance is the server health view inside Vessl. It helps you catch disk pressure before deployments, backups, or database volumes run the server out of space.

## What Vessl Measures

Vessl records:

- Root filesystem usage from `df`.
- Docker storage and reclaimable data from Docker system data.
- Vessl data directory usage.
- Build artifact usage.
- Database backup usage.
- APT cache usage.
- System journal usage.

Maintenance history keeps the last 48 points for disk usage, Docker reclaimable bytes, and build artifact usage.

## Alerts

Vessl raises maintenance alerts when:

- Disk usage crosses warning or critical thresholds.
- Docker reclaimable data grows large.
- Build artifacts grow large.

The exact underlying numbers are visible in the maintenance panels so you can decide whether cleanup is enough or the server needs more capacity.

## Safe Cleanup

Safe cleanup can run these cleanup targets:

- Stopped Docker containers.
- Unused Docker images.
- Docker build cache.
- APT cache.
- System journals down to about `100M`.
- Vessl build artifacts older than 24 hours.

Safe cleanup avoids Docker volume pruning because old database data can live in unattached volumes.

## Volume Cleanup

`Clean volumes` is separate and requires confirmation. It runs Docker volume prune and can delete old data from removed containers.

Do not use volume cleanup as a casual disk fix. Create backups first and make sure you do not need any detached database volumes.

## When to Use Maintenance

- Before a large Railway import.
- Before exporting an Vessl migration bundle.
- Before enabling `disk+r2` backups on many databases.
- After repeated failed builds.
- When deploys fail with Docker or disk space errors.

Maintenance does not replace server monitoring, but it gives the control plane a first-class view of the resources it consumes.
