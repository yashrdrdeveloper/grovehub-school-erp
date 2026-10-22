#!/usr/bin/env bash
set -Eeuo pipefail

cat <<'EOF'
GroveHub developer commands

Backend
  make backend-build    Build the platform binary into .runtime/bin
  make backend-start    Start the backend in the background
  make backend-stop     Stop the managed backend process
  make backend-status   Show backend process and health status
  make backend-dev      Run the backend in the foreground

Swagger
  make swagger-start    Start Swagger UI on http://localhost:8081
  make swagger-stop     Stop the managed Swagger UI container
  make swagger-status   Show Swagger UI status

Development
  make dev-start        Start backend + Swagger UI
  make dev-stop         Stop Swagger UI + backend
  make dev-status       Show backend + Swagger UI status

Overrides
  make swagger-start SWAGGER_PORT=8082
EOF
