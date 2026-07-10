// @ts-check

import starlight from '@astrojs/starlight';
import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://vessel.dev',
  base: '/docs',
  integrations: [
    starlight({
      title: 'Vessel Docs',
      customCss: ['./src/styles/theme.css'],
      sidebar: [
        { label: 'Getting Started', slug: 'getting-started' },
        { label: 'Deployment', slug: 'deployment' },
        { label: 'Databases', slug: 'databases' },
        { label: 'Configuration', slug: 'configuration' },
      ],
    }),
  ],
});
