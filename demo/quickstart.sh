#!/usr/bin/env bash
# Bootstraps a campaign + API key (idempotent) then evaluates a cart,
# so the quickstart GIF only ever shows the interesting part.
set -euo pipefail

API="http://localhost:8080"

TOKEN=$(curl -s -X POST "$API/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"organization":"org_default","email":"admin@example.com","password":"change-me"}' \
  | jq -r .token)

PROJECT=$(curl -s "$API/admin/projects" -H "Authorization: Bearer $TOKEN" \
  | jq -r '.projects[0].publicId')

EXISTING_CAMPAIGN=$(curl -s "$API/admin/projects/$PROJECT/campaigns" -H "Authorization: Bearer $TOKEN" \
  | jq -r '.campaigns[] | select(.name == "10% off orders over €50") | .publicId' | head -n1)

if [ -z "$EXISTING_CAMPAIGN" ]; then
  curl -s -X POST "$API/admin/projects/$PROJECT/campaigns" \
    -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d '{
      "name": "10% off orders over €50",
      "status": "active",
      "rules": [
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
      ]
    }' >/dev/null
fi

APIKEY=$(curl -s -X POST "$API/admin/projects/$PROJECT/api-keys" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name": "quickstart-demo"}' | jq -r '.apiKey')

echo '$ curl -X POST http://localhost:8080/v1/evaluate \'
echo '    -H "Authorization: ApiKey <your-key>" \'
echo '    -d '"'"'{ "customer": {"id": "user_123", "tier": "gold"}, "cart": {"subtotal": 120, "currency": "EUR"} }'"'"''
echo

curl -s -X POST "$API/v1/evaluate" \
  -H "Authorization: ApiKey $APIKEY" \
  -H "Content-Type: application/json" \
  -d '{
    "customer": { "id": "user_123", "country": "DE", "tier": "gold" },
    "cart": { "subtotal": 120, "currency": "EUR" }
  }' | jq .
