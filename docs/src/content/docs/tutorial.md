---
title: Deploy Your First App
description: Step-by-step tutorial — deploy a sample Node.js app on Vessl in under 5 minutes.
---

This tutorial walks you through deploying a real application on Vessl. You'll deploy a ready-to-run Node.js app, connect it to a database, and see it live.

## Prerequisites

- A running Vessl instance ([install guide](/getting-started/))
- An admin account (create via dashboard or `vesslctl setup`)

## Step 1: Create a Project

1. Log in to your Vessl dashboard at `http://your-server-ip:8080`.
2. Click **New Project** in the sidebar.
3. Enter `my-first-app` as the project name.
4. Click **Create**.

Vessl creates the project with a default **production** environment.

## Step 2: Deploy a Sample App

Since you don't have a Git repository connected yet, you'll deploy from a public Git URL.

1. In your project, click **New Service**.
2. Select **Deploy from Git URL**.
3. Paste this URL:
   ```
   https://github.com/expressjs/express
   ```
4. Vessl auto-detects the build strategy — it finds the `package.json` and uses Railpack.
5. Click **Deploy**.

Vessl clones the repository, detects Node.js, builds the image, and runs health checks.

**Check progress:** Click on the service to view live deployment logs.

Once the status shows **running**, your app is live at the URL shown in the service details.

## Step 3: Add a Database

1. Navigate to **Databases** in the sidebar.
2. Click **New Database**.
3. Select **PostgreSQL 16**.
4. Click **Create**.

Vessl provisions a PostgreSQL container with persistent storage.

After creation, Vessl automatically injects the connection string into your app:

```
DATABASE_URL=postgresql://vessl:<password>@<container-name>:5432/vessl
```

## Step 4: Set Environment Variables

1. Go to your service's **Variables** tab.
2. Click **Add Variable**.
3. Set `NODE_ENV` to `production`.
4. Save.

The service restarts automatically with the new environment.

## Step 5: Attach a Custom Domain (Optional)

1. Go to **Project → Domains**.
2. Click **Add Domain**.
3. Enter your domain (e.g. `app.example.com`).
4. Add the `A` record at your DNS provider.
5. SSL is provisioned automatically via Let's Encrypt.

## Step 6: Monitor

View real-time metrics from the service detail page:

- **CPU**: Current and historical usage
- **Memory**: RAM consumption
- **Status**: Running, deploying, or failed
- **Logs**: Live streaming output

## Clean Up

To remove everything:

```sh
# Via CLI
vesslctl backup   # Backup your database first

# In dashboard: Delete the project
```

## Next Steps

- [Configure a custom domain](/deployment/#custom-domains)
- [Add environment variables](/deployment/#environment-variables)
- [Set up automatic Git deployments](/integrations/#git-providers)
- [Create a serverless function](/serverless/)
