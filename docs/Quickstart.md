# Quickstart

Get Open Incentives running in a few minutes.

## 1. Start Open Incentives

``` bash
docker compose up --build
```

The application will:

-   Create the SQLite database (if it doesn't exist)
-   Run database migrations
-   Bootstrap a default organization and project
-   Start the API on `http://localhost:8080`

## 2. Log in

The default bootstrapped login credentials are:

-   **Email:** `admin@example.com`
-   **Password:** `change-me`

(Change these before running in production by setting `BOOTSTRAP_ADMIN_EMAIL` /
`BOOTSTRAP_ADMIN_PASSWORD`.)

### Option A: Admin Portal (recommended)

Open `http://localhost:3001` and sign in with the credentials above. From there
you can create projects, campaigns, and API keys without touching curl — skip
straight to [step 6](#6-try-the-demo-store-optional) once you've created a
campaign and API key through the UI.

### Option B: curl

``` bash
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "organization": "org_default",
    "email": "admin@example.com",
    "password": "change-me"
  }'
```

You'll need the returned `token` to authenticate the rest of the `/admin` requests
in the steps below.

## 3. Create a campaign

If you logged in via the Admin Portal (Option A above), you can create a campaign
there directly and skip to [step 5](#5-evaluate). The curl steps below continue
the API-only path (Option B).

Replace `proj_xxxxxxxxxxxxx` with your project's public id, and `<TOKEN>` with the
token from step 2.

``` bash
curl -X POST http://localhost:8080/admin/projects/proj_xxxxxxxxxxxxx/campaigns \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
  "name": "10% off orders over €50",
  "status": "active",
  "rules": [
    {
      "id": "rule_10_percent_over_50",
      "name": "10% off orders over €50",
      "conditions": {
        "all": [
          {
            "fact": "cart.subtotal",
            "operator": "gte",
            "value": 50
          }
        ]
      },
      "actions": [
        {
          "type": "percentage_discount",
          "params": {
            "value": 10
          }
        }
      ]
    }
  ]
}'
```

## 4. Create an API key

Copy the generated key — you'll use it to authenticate `/v1/evaluate` requests.

``` bash
curl -X POST http://localhost:8080/admin/projects/proj_xxxxxxxxxxxxx/api-keys \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "new-key"
  }'
```

## 5. Evaluate

``` bash
curl -X POST http://localhost:8080/v1/evaluate \
  -H "Authorization: ApiKey api_xxxxx.yyyyy" \
  -H "Content-Type: application/json" \
  -d '{
    "customer": {
      "id": "user_123",
      "country": "DE",
      "tier": "gold"
    },
    "cart": {
      "currency": "EUR",
      "items": [
        { "productId": "prod_coffee", "quantity": 2, "unitPrice": 18 },
        { "productId": "prod_mug", "quantity": 1, "unitPrice": 14 }
      ]
    }
  }'
```

Replace `api_xxxxx.yyyyy` with the API key you created in step 4.

If an active campaign matches the request, the evaluation response will include the
matching incentives. `cart.items[]` is optional — you can send `cart.subtotal`
directly instead if you don't need per-item pricing.

## 6. Try the demo store (optional)

`apps/demo-store` is a small Next.js storefront that shows how a third-party store
integrates with the `/v1/evaluate` API — see `apps/demo-store/README.md` for details.

It isn't started by `docker compose up` by default (it lives behind the `demo-store`
Compose profile). To run it:

``` bash
DEMO_STORE_API_KEY=api_xxxxx.yyyyy docker compose --profile demo-store up --build
```

Then open `http://localhost:3000`. Use the API key you created in step 4 — the
demo store's server-side code is the only place it's read; it's never sent to the
browser.

To run the demo store without Docker instead:

``` bash
cd apps/demo-store
cp .env.example .env.local
npm install
npm run dev
```
