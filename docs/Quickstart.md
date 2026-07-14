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

Open the Admin Portal at `http://localhost:3001` and sign in with the default
bootstrapped credentials:

-   **Email:** `admin@example.com`
-   **Password:** `change-me`

(Change these before running in production by setting `BOOTSTRAP_ADMIN_EMAIL` /
`BOOTSTRAP_ADMIN_PASSWORD`.)

<details>
<summary>Prefer curl?</summary>

``` bash
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "organization": "org_default",
    "email": "admin@example.com",
    "password": "change-me"
  }'
```

Use the returned `token` as a `Bearer` token on the `/admin` requests below in
place of the Admin Portal steps.

</details>

## 3. Open your project

On the **Projects** page, click your project's name — this takes you into the
project and adds **Campaigns** and **API Keys** to the top nav.

## 4. Create a campaign

From the **Campaigns** tab, click **New campaign**, then:

-   **Name:** `10% off orders over €50`
-   **Status:** `Active`
-   **Rules (JSON):** rules are entered as a JSON block rather than individual
    fields — paste:

    ``` json
    {
      "id": "rule_10_percent_over_50",
      "name": "10% off orders over €50",
      "conditions": {
        "all": [
          { "fact": "cart.subtotal", "operator": "gte", "value": 50 }
        ]
      },
      "actions": [
        { "type": "percentage_discount", "params": { "value": 10 } }
      ]
    }
    ```

Click **Save**. A read-only summary of the rule renders below the JSON editor
so you can sanity-check it before saving.

<details>
<summary>Prefer curl?</summary>

Replace `proj_xxxxxxxxxxxxx` with your project's public id, and `<TOKEN>` with
the token from step 2.

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

</details>

## 5. Create an API key

From the **API Keys** tab, click **New API key**, give it a name and
description, then click **Create**. The key is shown exactly once — copy it
now, since the Admin Portal never displays it again. You'll use it to
authenticate `/v1/evaluate` requests below.

<details>
<summary>Prefer curl?</summary>

``` bash
curl -X POST http://localhost:8080/admin/projects/proj_xxxxxxxxxxxxx/api-keys \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "new-key"
  }'
```

</details>

## 6. Evaluate

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

Replace `api_xxxxx.yyyyy` with the API key you created in step 5.

If an active campaign matches the request, the evaluation response will include the
matching incentives. `cart.items[]` is optional — you can send `cart.subtotal`
directly instead if you don't need per-item pricing.

## 7. Try the demo store (optional)

`apps/demo-store` is a small Next.js storefront that shows how a third-party store
integrates with the `/v1/evaluate` API — see `apps/demo-store/README.md` for details.

It isn't started by `docker compose up` by default (it lives behind the `demo-store`
Compose profile). To run it:

``` bash
DEMO_STORE_API_KEY=api_xxxxx.yyyyy docker compose --profile demo-store up --build
```

Then open `http://localhost:3000`. Use the API key you created in step 5 — the
demo store's server-side code is the only place it's read; it's never sent to the
browser.

To run the demo store without Docker instead:

``` bash
cd apps/demo-store
cp .env.example .env.local
npm install
npm run dev
```
