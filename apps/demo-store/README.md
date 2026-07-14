# Open Incentives Demo Store

A small demo e-commerce storefront showing how a third-party store integrates with the
Open Incentives evaluation API. It is a product demonstration, not a production
e-commerce platform — see `demo-store-mvp-spec.md` for the full spec.

## Local setup

```bash
cp .env.example .env.local
npm install
npm run dev
```

The app runs at http://localhost:3000.

## Required configuration

```env
INCENTIVES_API_URL=http://localhost:8080
INCENTIVES_API_KEY=api_example.secret
```

`INCENTIVES_API_KEY` is read only on the server and is never sent to the browser.

## Docker setup

The `demo-store` service is behind the `demo-store` Compose profile, so it does not
start with a plain `docker compose up`. From the repository root:

```bash
DEMO_STORE_API_KEY=api_xxxxx.yyyyy docker compose --profile demo-store up --build
```

This starts both the Open Incentives API (`api`) and this store (`demo-store`), with
`demo-store` talking to `api` over the Docker network at `http://api:8080`.
`DEMO_STORE_API_KEY` must be an API key created for a project on that API instance
(see the repo-root `docs/Quickstart.md`).

## Integration overview

```text
Browser → Demo Store backend → Open Incentives API
```

- The browser only ever talks to this app's own `POST /api/evaluate` route.
- That route runs on the server, reads `INCENTIVES_API_URL` / `INCENTIVES_API_KEY`,
  and calls `POST {INCENTIVES_API_URL}/v1/evaluate` with
  `Authorization: ApiKey <INCENTIVES_API_KEY>`.
- The API key is only ever used by the demo-store's server-side code
  (`lib/incentives/client.ts`) — it is never included in any browser JavaScript,
  frontend environment variable, or client-side log.

## Tests

```bash
npm test
```
