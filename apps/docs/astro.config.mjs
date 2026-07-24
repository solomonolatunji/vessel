// @ts-check

import starlight from '@astrojs/starlight';
import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://docs.codedock.run',
  integrations: [
    starlight({
      title: 'Codedock Docs',
      favicon: '/favicon.svg',
      customCss: ['./src/styles/theme.css'],
      sidebar: [
        {
          label: 'Start here',
          items: [
            { label: 'Codedock Docs', slug: 'index' },
            { label: 'Installation', slug: 'getting-started/installation' },
            { label: 'Browser Onboarding', slug: 'getting-started/onboarding' },
            { label: 'Quick Start', slug: 'getting-started/quick-start' },
            { label: 'Deploy Your First App', slug: 'tutorial' },
          ],
        },
        {
          label: 'Core concepts',
          items: [
            { label: 'Architecture', slug: 'core-concepts/architecture' },
            { label: 'Projects Overview', slug: 'projects/overview' },
            { label: 'Project Settings', slug: 'projects/settings' },
          ],
        },
        {
          label: 'Deployments',
          items: [
            { label: 'Build Strategies', slug: 'deployments/build-strategies' },
            { label: 'Templates & Examples', slug: 'deployments/templates' },
            { label: 'CI/CD & Git', slug: 'deployments/ci-cd' },
            { label: 'Service Types', slug: 'deployments/service-types' },
            { label: 'Cron Jobs', slug: 'deployments/jobs' },
            {
              label: 'Environment Variables',
              slug: 'deployments/environment-variables',
            },
            {
              label: 'Deployment Lifecycle',
              slug: 'deployments/deployment-lifecycle',
            },
          ],
        },
        {
          label: 'Databases',
          items: [
            { label: 'Database Provisioning', slug: 'databases/provisioning' },
            { label: 'SQL Studio', slug: 'databases/sql-studio' },
            { label: 'Data Browser', slug: 'databases/data-browser' },
            { label: 'Data Imports', slug: 'databases/data-imports' },
            {
              label: 'Public Access & TLS',
              slug: 'databases/public-access-and-tls',
            },
          ],
        },
        {
          label: 'Storage & Backups',
          items: [
            {
              label: 'Storage (MinIO)',
              slug: 'storage-and-backups/minio-storage',
            },
            {
              label: 'R2 Storage Integration',
              slug: 'storage-and-backups/r2-storage',
            },
            {
              label: 'Database Backups',
              slug: 'storage-and-backups/database-backups',
            },
            {
              label: 'Restore & Download',
              slug: 'storage-and-backups/restore-and-download',
            },
          ],
        },
        {
          label: 'Migration',
          items: [{ label: 'Codedock Bundles', slug: 'migration/codedock-bundles' }],
        },
        {
          label: 'Operations',
          items: [
            { label: 'Teams & Collaboration', slug: 'operations/teams' },
            { label: 'Administration', slug: 'admin' },
            { label: 'Integrations', slug: 'integrations' },
            { label: 'Observability & Logs', slug: 'operations/observability' },
            { label: 'Domains & DNS', slug: 'operations/domains-and-dns' },
            { label: 'Account Security', slug: 'operations/account-security' },
            {
              label: 'Maintenance & Updates',
              slug: 'operations/maintenance-and-updates',
            },
            { label: 'Troubleshooting', slug: 'operations/troubleshooting' },
          ],
        },
        {
          label: 'Reference',
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
            {
              label: 'API Env & Domains',
              slug: 'reference/api-environment-and-domains',
            },
            { label: 'No Lock-In', slug: 'adopt' },
          ],
        },
      ],
      components: {
        SiteTitle: './src/components/docs-site-title.astro',
        ThemeSelect: './src/components/docs-theme-select.astro',
      },
      social: [
        {
          icon: 'github',
          label: 'GitHub',
          href: 'https://github.com/buildwithtechx/codedock',
        },
      ],
    }),
  ],
});
