#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

# Auth/Identity backend is the only required development service.
BACKEND_HOT_RELOAD=1 bash "${SCRIPT_DIR}/../backend/start.sh"

# Swagger is optional until an Auth OpenAPI spec is available.
# To start it explicitly: DEV_SWAGGER=1 make dev-start
if [[ "${DEV_SWAGGER:-0}" == "1" ]]; then
  bash "${SCRIPT_DIR}/../swagger/start.sh"
fi
