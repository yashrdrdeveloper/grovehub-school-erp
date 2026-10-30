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

Mail Provider
  make mail-provider-start   Ensure the container runtime and start Haraka
  make mail-provider-stop    Stop Haraka while preserving queue and DKIM state
  make mail-provider-status  Show Haraka container and health status
  make mail-domain-ready     Read Cloudflare token/zone from .env, discover domain/public IP,
                             prepare durable DKIM state, configure outbound-mail DNS, and verify the domain

Database
  make db-schema-status                           Show current GroveHub schema baseline state
  make db-schema-init                             Initialize all service schemas from the current baseline
  make db-schema-reset                            Destructively reset all service schemas to the current baseline
  make db-schema-reset DB_SERVICE=foundation      Reset only the Foundation schema
  make db-schema-reset DB_SERVICE=infrastructure  Reset only the Infrastructure schema

  DB_SERVICE can be: all, foundation, identity, infrastructure, or mailProvider.
  Schema reset is development-only and deletes data in the selected service schema.

Swagger
  make swagger-start    Start Swagger UI on http://localhost:8081
  make swagger-stop     Stop the managed Swagger UI container
  make swagger-status   Show Swagger UI status

Development
  make dev-start            Auto-start container runtime, Haraka, schema initialization, Air backend, and Swagger
  make dev-start-with-logs  Same as dev-start plus the local OTEL collector
  make dev-stop             Stop Swagger UI, backend, Haraka, and optional OTEL collector
  make dev-status           Show backend, Mail Provider, and Swagger UI status

Container runtime used by dev-start
  macOS    Colima (installed/started automatically; Docker CLI/Compose via Homebrew)
  Windows  Docker Desktop (installed through winget/Chocolatey and started automatically)
  Linux    Docker Engine (installed through the official Docker installer and started automatically)

Overrides
  make swagger-start SWAGGER_PORT=8082
EOF
