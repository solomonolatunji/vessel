# Vessel Dashboard

**Self-hosted control panel** — served by the `vesseld` daemon binary on your own server.

Built with React 19, TanStack Router, TanStack Query, Radix UI, and Tailwind CSS v4.

```sh
npm install
npm run dev       # http://localhost:3000
npm run build     # → served by daemon at /dashboard
```

## Project Structure

```text
src/
├── routes/          # TanStack Router file-based routes
├── components/      # UI components grouped by domain
├── hooks/           # WebSocket streaming hooks, TanStack Query mutations
├── lib/             # Utilities and helpers
└── styles.css       # Tailwind v4 global styles
```

## Commands

| Command              | Action                               |
| -------------------- | ------------------------------------ |
| `npm run dev`        | Start dev server at `localhost:3000` |
| `npm run build`      | Build for production to `dist/`      |
| `npm run test`       | Run Vitest tests                     |
| `npm run format`     | Check formatting with Biome          |
| `npm run format:fix` | Fix formatting with Biome            |

## Domain

This dashboard ships with the daemon binary. Access it at your instance's IP/domain after running `vesseld`. Not used by cloud users — they use `cloud.vessel.dev` instead.

## Learn More

- [TanStack Router docs](https://tanstack.com/router)
- [TanStack Query docs](https://tanstack.com/query)
- [Radix UI docs](https://www.radix-ui.com/)
- [Tailwind CSS v4 docs](https://tailwindcss.com/)
