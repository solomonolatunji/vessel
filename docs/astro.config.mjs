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
        { label: 'Getting Started', slug: 'getting-started' },
        { label: 'Deploy Your First App', slug: 'tutorial' },
        {
          label: 'Deployment',
          collapsed: false,
          items: [
            { label: 'Overview', slug: 'deployment' },
            { label: 'Projects & Environments', slug: 'projects' },
            { label: 'Serverless Functions', slug: 'serverless' },
          ],
        },
        { label: 'Databases', slug: 'databases' },
        { label: 'Storage', slug: 'storage' },
        { label: 'Workspaces & Teams', slug: 'workspaces' },
        { label: 'Integrations', slug: 'integrations' },
        { label: 'Configuration', slug: 'configuration' },
        { label: 'Administration', slug: 'admin' },
        { label: 'API Reference', slug: 'api' },
      ],
    }),
  ],
});
