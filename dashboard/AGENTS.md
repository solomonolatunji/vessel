# Development

```sh
npm run dev
```

The dev server starts at `http://localhost:3000`.

## Build

```sh
npm run build
```

Output goes to `dist/`.

## Routing

Routes live in `src/routes/` following TanStack Router file conventions. Run `npm run generate-routes` after adding new route files.

## Components

Add shadcn/Radix UI components with:

```sh
npx shadcn@latest add button
```

Components go in `src/components/ui/`. Always prefer an existing Radix component over building from scratch.

## Documentation

- [TanStack Router](https://tanstack.com/router)
- [TanStack Query](https://tanstack.com/query)
- [Radix UI](https://www.radix-ui.com/)
- [Tailwind CSS v4](https://tailwindcss.com/)
- [Lucide Icons](https://lucide.dev/)
- [Hugeicons](https://hugeicons.com/)
