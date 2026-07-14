# Open Incentives Admin Portal

The admin portal for managing Projects, Campaigns, and API Keys. Built with Vite, React, TypeScript, Tailwind CSS, and TanStack Query.

It talks directly to the Go admin API from the browser using the session token returned at login — there is no server-side proxy.

## Running via Docker Compose

This service is part of the default stack (not behind a `--profile` flag), so it starts automatically:

```bash
docker compose up
```

Once running: http://localhost:3001

## Local development

Requires the API running locally (`go run ./apps/api` or `make up` for just the backend), reachable at `http://localhost:8080` by default.

```bash
cp .env.example .env
npm install
npm run dev
```

The dev server runs at http://localhost:5173.

### Environment variables

| Variable | Purpose | Default |
|---|---|---|
| `VITE_API_BASE_URL` | Base URL of the admin API, used at build/dev time | `http://localhost:8080` |
| `VITE_DEFAULT_ORG_ID` | Organization public ID used for login (single-tenant today — see below) | `org_default` |

## Why login only asks for email/password

The backend's login endpoint takes `{email, password, organization}`. Since this is currently a single-tenant deployment, `organization` is sourced from `VITE_DEFAULT_ORG_ID` rather than shown as a form field. Adding a real org switcher later means: adding an org ID to the URL structure alongside `:projectId`, an `useOrgs()`/`OrgSwitcher` mirroring `ProjectSwitcher`, and sourcing `organization` from that selection instead of the env default — an additive change, not a rewrite.

## Docker runtime environment injection

Vite bakes `VITE_*` variables into the JS bundle at **build time**. But the API URL a deployed image should use depends on how the operator exposes the `api` service — baking it in once would make the image non-portable across environments.

To avoid that, the Docker build (`.env.production`) builds with a placeholder value, `VITE_API_BASE_URL=__RUNTIME_API_BASE_URL__`. At container **start**, `docker-entrypoint.sh` (auto-run by the official nginx image) does a literal string substitution over the built JS assets, replacing the placeholder with the real `API_BASE_URL` environment variable read at that point. This is a different variable from `VITE_API_BASE_URL` on purpose — `API_BASE_URL` (no `VITE_` prefix) is consumed by the entrypoint script at runtime, not by Vite at build time. Don't "fix" this into a single build-time-only variable; that's what reintroduces the portability problem.

`API_BASE_URL` must be a **browser-resolvable** address (e.g. `http://localhost:8080`), never the Docker-internal service DNS name (`http://api:8080`) — the browser, not the container, makes the actual API call.

## Adding a new resource

Each feature area under `src/features/<resource>/` owns an `api.ts` with all its TanStack Query hooks, following one convention throughout:

```ts
export function useThings(projectId: string) {
  return useQuery({
    queryKey: ['things', projectId],
    queryFn: () => apiFetch<{ things: Thing[] }>(`/admin/projects/${projectId}/things`),
  });
}
```

Mutations invalidate the relevant query key on success. Follow this pattern (see `features/campaigns/api.ts`, `features/projects/api.ts`, `features/apikeys/api.ts`) when adding new resources like Coupons or Events.

## Prior art

`docs/` in this folder holds the original product spec documents written before this app was built, kept for historical context — not the operational source of truth (this README is).
