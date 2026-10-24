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
  make mail-provider-stop    Stop Haraka while preserving queue and DKIM volumes
  make mail-provider-status  Show Haraka container and health status
  make mail-domain-ready DOMAIN=school.example PUBLIC_IP=203.0.113.10
                             Register/reuse the domain, generate DKIM, configure Cloudflare DNS, and verify it

Database
  make db-migrate-status  Show all GroveHub service migration status using DATABASE_URL
  make db-migrate-up      Apply pending migrations for all services
  make db-migrate-down    Roll back the latest migration for each service in reverse order

Swagger
  make swagger-start    Start Swagger UI on http://localhost:8081
  make swagger-stop     Stop the managed Swagger UI container
  make swagger-status   Show Swagger UI status

Development
  make dev-start            Auto-start container runtime, Haraka, migrations, Air backend, and Swagger
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
