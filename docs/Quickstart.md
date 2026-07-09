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
-   Start the API and Admin UI

## 2. Open the Admin UI

``` text
http://localhost:8080/admin
```

## 3. Create an API key

From the Admin UI, create your first API key.

Copy the generated key---you'll use it to authenticate requests.

## 4. Create a campaign

Create your first campaign and mark it as **Active**.

## 5. Evaluate

``` bash
curl -X POST http://localhost:8080/v1/evaluate \
  -H "Authorization: Bearer oi_live_xxx" \
  -H "Content-Type: application/json" \
  -d '{
    "customer": {
      "id": "user_123",
      "country": "DE",
      "tier": "gold"
    },
    "cart": {
      "total": 120,
      "currency": "EUR"
    }
  }'
```

Replace `oi_live_xxx` with the API key you created in the Admin UI.

If an active campaign matches the request, the evaluation response will
include the matching incentives.
