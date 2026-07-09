# Vessel Website

**Public marketing landing page** — hosted at `vessel.dev`.

Built with Astro 7 and Tailwind CSS v4. Routes visitors to either self-hosted (`dashboard/`) or cloud (`cloud.vessel.dev`).

```sh
npm install
npm run dev       # http://localhost:4321
npm run build     # outputs to dist/
```

## Project Structure

```text
src/
├── pages/          # Astro page routes
├── components/     # Reusable Astro components
├── layouts/        # Page layout templates
└── assets/         # Static assets (images, icons)
```

## Commands

| Command           | Action                               |
| ----------------- | ------------------------------------ |
| `npm run dev`     | Start dev server at `localhost:4321` |
| `npm run build`   | Build for production to `dist/`      |
| `npm run preview` | Preview production build locally     |

## Learn More

- [Astro docs](https://docs.astro.build)
