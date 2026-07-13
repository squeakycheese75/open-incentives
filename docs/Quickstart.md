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
-  ~~ Start the API and Admin UI~~

<!-- ## 2. Open the Admin UI (WIP)

``` text
http://localhost:8080/admin
``` -->
## 2. Login

``` bash
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{

    "organization": "org_default",
    "email": "admin@example.com",
    "password": "change-me"
  }'
```

You'll need the returned Token to authenticate the rest of the /admin requests.  

## 3. Create a campaign

Replace `oi_api_key_xxx` with the Token you generated.

Create your first campaign and mark it as **Active**.  You'll need to lookup the default project public id.

``` bash
curl -X POST http://localhost:8080//admin/projects/proj_xxxxxxxxxxxxx/campaigns\
  -H "Authorization: Bearer oi_live_xxx" \
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

Create your first API key.

Copy the generated key---you'll use it to authenticate requests.

``` bash
curl -X POST http://localhost:8080//admin/projects/proj_xxxxxxxxxxxx/api-keys \
  -H "Authorization: Bearer oi_live_xxx" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"new-key",
    "description":"used for me"
  }'
```


## 5. Evaluate

``` bash
curl -X POST http://localhost:8080/v1/evaluate \
  -H "Authorization: Bearer oi_api_key_xxx" \
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

Replace `oi_api_key_xxx` with the API key you created.

If an active campaign matches the request, the evaluation response will
include the matching incentives.
