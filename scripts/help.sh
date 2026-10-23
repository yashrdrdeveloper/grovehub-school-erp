#!/usr/bin/env bash
set -Eeuo pipefail

cat <<'EOF'
GroveHub developer commands

Backend
  make backend-setup    Install/pin local backend development tools
  make backend-build    Build the platform binary into .runtime/bin
  make backend-start    Start the backend in the background
  make backend-stop     Stop the managed backend process
  make backend-status   Show backend process and health status
  make backend-dev      Run backend with Air hot reload

Swagger
  make swagger-start    Start Swagger UI on http://localhost:8081
  make swagger-stop     Stop the managed Swagger UI container
  make swagger-status   Show Swagger UI status

Development
  make dev-start            Start backend + Swagger UI
  make dev-start-with-logs  Start backend + Swagger UI + local OTEL log collector
  make dev-stop             Stop Swagger UI + backend + optional OTEL collector
  make dev-status           Show backend + Swagger UI status

Overrides
  make swagger-start SWAGGER_PORT=8082
EOF
