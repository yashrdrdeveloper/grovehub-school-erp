#!/usr/bin/env bash
set -Eeuo pipefail

cat <<'EOF'
GroveHub developer commands

Backend
  make backend-setup    Install/pin local backend development tools
  make backend-build    Build the platform binary into .runtime/bin
  make backend-start    Start the backend binary in the background (no hot reload)
  make backend-stop     Stop the managed backend process (binary or Air)
  make backend-status   Show backend process, hot-reload mode, and health status
  make backend-dev      Run backend with Air hot reload in the foreground

Database
  make db-schema-status                           Show current GroveHub schema baseline state
  make db-schema-init                             Initialize all service schemas from the current baseline
  make db-schema-reset                            Destructively reset all service schemas to the current baseline
  make db-schema-reset DB_SERVICE=foundation      Reset only the Foundation schema
  make db-schema-reset DB_SERVICE=infrastructure  Reset only the Infrastructure schema

  DB_SERVICE can be: all, foundation, identity, or infrastructure.
  Schema reset is development-only and deletes data in the selected service schema.

Messaging / NATS JetStream
  make nats-start     Start local NATS JetStream with persistent storage
  make nats-stop      Stop the NATS container and preserve JetStream data
  make nats-status    Show NATS health, endpoints, and persistent volume

Swagger
  make swagger-start    Start Swagger UI on http://localhost:8081
  make swagger-stop     Stop the managed Swagger UI container
  make swagger-status   Show Swagger UI status

Development
  make dev-start            Start NATS, initialize schemas, start Air backend, and start Swagger
  make dev-start-with-logs  Same as dev-start plus the local OTEL collector
  make dev-stop             Stop Swagger UI, backend, NATS, and optional OTEL collector
  make dev-status           Show NATS, backend, and Swagger UI status

Container runtime used by dev-start
  macOS    Colima (installed/started automatically; Docker CLI/Compose via Homebrew)
  Windows  Docker Desktop (installed through winget/Chocolatey and started automatically)
  Linux    Docker Engine (installed through the official Docker installer and started automatically)

Overrides
  make swagger-start SWAGGER_PORT=8082
  make nats-start NATS_CLIENT_PORT=4223 NATS_MONITOR_PORT=8223
EOF
