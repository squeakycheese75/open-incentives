#!/usr/bin/env bash
set -euo pipefail

API=${API:-http://localhost:8080}

: "${TOKEN:?Please run 'source ./examples/admin/auth/login.sh' first.}"

curl -X POST "$API/admin/projects/proj_5043b130-2174-45c6-9165-e9d736d754e0/campaigns" \
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