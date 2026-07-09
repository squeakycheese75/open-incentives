#!/usr/bin/env bash
set -euo pipefail

API=${API:-http://localhost:8080}

: "${TOKEN:?Please run 'source ./examples/admin/auth/login.sh' first.}"

curl -X POST "$API/admin/campaigns" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "VIP Campaign",
    "status": "active",
    "rule": {
      "rules": [
        {
          "field": "country",
          "operator": "eq",
          "value": "US"
        },
        {
          "field": "cart_total",
          "operator": "gte",
          "value": 100
        }
      ]
    }
  }'