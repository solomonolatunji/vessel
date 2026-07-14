// @ts-check

import starlight from '@astrojs/starlight';
import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://docs.vessl.dev',
  integrations: [
    starlight({
      title: 'Vessl Docs',
      customCss: ['./src/styles/theme.css'],
      sidebar: [
        {
          label: "Start here",
          items: [
            { label: 'Vessl Docs', slug: 'index' },
            { label: 'Installation', slug: 'getting-started' },
            { label: 'Browser Onboarding', slug: 'getting-started/onboarding' },
            { label: 'Deploy Your First App', slug: 'tutorial' },
          ],
        },
        {
          label: "Core concepts",
          items: [
            { label: 'Architecture', slug: 'core-concepts/architecture' },
            { label: 'Projects & Environments', slug: 'projects' },
          ],
        },
        {
          label: "Deployments",
          items: [
            { label: 'Overview', slug: 'deployment' },
            { label: 'Source Services', slug: 'deployments/source-services' },
            { label: 'Docker Image Services', slug: 'deployments/docker-image-services' },
            { label: 'Serverless Functions', slug: 'serverless' },
            { label: 'Static Sites & Workers', slug: 'deployments/static-sites-and-workers' },
            { label: 'Environment Variables', slug: 'deployments/environment-variables' },
            { label: 'Deployment Lifecycle', slug: 'deployments/deployment-lifecycle' },
          ],
        },
        {
          label: "Databases",
          items: [
            { label: 'Overview', slug: 'databases' },
            { label: 'Data Browser', slug: 'databases/data-browser' },
            { label: 'Data Imports', slug: 'databases/data-imports' },
            { label: 'Public Access & TLS', slug: 'databases/public-access-and-tls' },
          ],
        },
        {
          label: "Storage & Backups",
          items: [
            { label: 'Storage (S3)', slug: 'storage' },
            { label: 'R2 Storage Integration', slug: 'storage-and-backups/r2-storage' },
            { label: 'Database Backups', slug: 'storage-and-backups/database-backups' },
            { label: 'Restore & Download', slug: 'storage-and-backups/restore-and-download' },
          ],
        },
        {
          label: "Migration",
          items: [
            { label: 'Vessl Bundles', slug: 'migration/vessl-bundles' },
            { label: 'Railway Import', slug: 'migration/railway-import' },
          ],
        },
        {
          label: "Operations",
          items: [
            { label: 'Administration', slug: 'admin' },
            { label: 'Integrations', slug: 'integrations' },
            { label: 'Domains', slug: 'operations/domains' },
            { label: 'System Maintenance', slug: 'operations/system-maintenance' },
            { label: 'System Updates', slug: 'operations/system-updates' },
            { label: 'Backups & Updates', slug: 'operations/backups-and-updates' },
            { label: 'Troubleshooting', slug: 'operations/troubleshooting' },
          ],
        },
        {
          label: "Reference",
          items: [
            { label: 'Configuration', slug: 'configuration' },
            { label: 'CLI Reference', slug: 'cli' },
            { label: 'API Reference (Full)', slug: 'api' },
            { label: 'System Settings', slug: 'reference/system-settings' },
            { label: 'API Access', slug: 'reference/api-access' },
            { label: 'API Projects', slug: 'reference/api-projects' },
            { label: 'API Services', slug: 'reference/api-services' },
            { label: 'API Deployments', slug: 'reference/api-deployments' },
            { label: 'API Databases', slug: 'reference/api-databases' },
            { label: 'API Env & Domains', slug: 'reference/api-environment-and-domains' },
            { label: 'No Lock-In', slug: 'adopt' },
          ],
        },
      ],
      components: {
        SiteTitle: "./src/components/docs-site-title.astro",
        ThemeSelect: "./src/components/docs-theme-select.astro",
      },
    }),
  ],
});
