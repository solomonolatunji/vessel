---
title: Projects Overview
description: Organize your services into projects with isolated environments.
---

Projects are the top-level organizational unit in Codedock. Each project contains services, databases, storage, and environment configuration.

## Creating a Project

1. Click **New Project** in the dashboard.
2. Enter a name and optional description.
3. Optionally, select a Git repository to connect.
4. Click **Create**.

A default **production** environment is created automatically.

## Environments

Environments provide isolation between development stages. Each environment has its own services, databases, and variables.

### Default Environments

Every new project gets a **production** environment. You can add more:

- **staging** — pre-production testing
- **dev** — development work
- **preview** — ephemeral PR previews

### Creating Environments

1. Open your project.
2. Go to **Environments**.
3. Click **New Environment**.
4. Enter a name (e.g. `staging`, `dev`).
5. Click **Create**.

### Environment Variables

Set environment variables per environment:

1. Open your project.
2. Go to **Environment Variables**.
3. Add key-value pairs.
4. These are injected into every service and job in the environment.

## Canvas View

The **Canvas** provides a visual overview of your environment — a node graph showing how apps, databases, and storage connect:

- **Service Nodes**: Each running service with status indicators.
- **Database Nodes**: Connected database instances.
- **Storage Nodes**: S3-compatible buckets.
- **Edges**: Connection lines showing which services link to which databases.

Access the Canvas from **Project → Canvas** or the environment detail page.
