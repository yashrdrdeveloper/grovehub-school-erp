#!/usr/bin/env bash
set -Eeuo pipefail

# The launcher owns platform health for the single GroveHub backend process.
BACKEND_HOST="127.0.0.1"
BACKEND_PORT="8080"
BACKEND_BASE_URL="http://${BACKEND_HOST}:${BACKEND_PORT}"
BACKEND_HEALTH_ENDPOINTS=(
  "/health"
)

# Add one "Display Name|path-under-backend/services" entry for each OpenAPI spec.
# Swagger UI is intentionally pinned; upgrade this version deliberately.
SWAGGER_IMAGE="swaggerapi/swagger-ui:v5.33.3"
SWAGGER_CONTAINER_NAME="groovehub-swagger"
SWAGGER_PORT="${SWAGGER_PORT:-8081}"
SWAGGER_BASE_URL="http://127.0.0.1:${SWAGGER_PORT}"
SWAGGER_SPECS=(
  "Foundation API v1|foundation/api/openapi/v1/openapi.yaml"
  "Identity API v1|identity/api/openapi/v1/openapi.yaml"
)
