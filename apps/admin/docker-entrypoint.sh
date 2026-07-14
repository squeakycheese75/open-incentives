#!/bin/sh
# The frontend is built with a placeholder API base URL baked into the JS
# bundle (Vite env vars are compile-time only). This substitutes the real,
# browser-resolvable API URL into the built assets at container start, so
# the same image works regardless of where/how the api service is exposed.
set -eu

: "${API_BASE_URL:=http://localhost:8080}"

grep -rl '__RUNTIME_API_BASE_URL__' /usr/share/nginx/html/assets/*.js 2>/dev/null | while read -r file; do
  sed -i "s|__RUNTIME_API_BASE_URL__|${API_BASE_URL}|g" "$file"
done
