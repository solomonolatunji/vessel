# ☁️ Vessel Cloud Dashboard

**Cloud user-facing dashboard** — hosted at `cloud.vessel.dev`.

Built with **TanStack Router + React + Tailwind CSS v4**.

## Routes

| Path      | Description                       |
| --------- | --------------------------------- |
| `/signup` | Create a cloud account            |
| `/signin` | Log in to cloud                   |
| `/`       | Project dashboard (requires auth) |
| `/admin`  | Staff-only admin panel            |

## Dev

```sh
npm run dev      # http://localhost:3002
npm run build    # outputs to dist/
```

## How It Fits In

- **Self-hosted users** — download the daemon and use `dashboard/` (served by the Go binary at your own IP/domain).
- **Cloud users** — sign up at `cloud.vessel.dev` and use this dashboard (talks to `internal/cloud/` API).
- **Marketing** — `vessel.dev` (`web/`) routes users to either path.
- **Docs** — `docs.vessel.dev` (`docs/`).
