#!/usr/bin/env bash

API=${API:-http://localhost:8081}

response=$(
curl -s -X POST "$API/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "organization": "default",
    "email": "admin@example.com",
    "password": "change-me"
  }'
)

export TOKEN=$(echo "$response" | jq -r '.token')

echo "Logged in."

#  source ./examples/admin/auth/login.sh 