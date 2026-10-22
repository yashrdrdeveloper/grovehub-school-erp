#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/container.sh"
source "${SCRIPT_DIR}/../config/services.sh"

require_docker

if container_is_running "${SWAGGER_CONTAINER_NAME}"; then
  log_info "Swagger UI: UP (${SWAGGER_BASE_URL})"
else
  log_warn "Swagger UI: DOWN"
fi
