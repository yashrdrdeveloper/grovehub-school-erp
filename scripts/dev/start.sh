#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

BACKEND_HOT_RELOAD=1 bash "${SCRIPT_DIR}/../backend/start.sh"

# Swagger is optional for Auth development while its OpenAPI spec is missing.
if [[ "${DEV_SWAGGER:-0}" == "1" ]]; then
  bash "${SCRIPT_DIR}/../swagger/start.sh"
fi
