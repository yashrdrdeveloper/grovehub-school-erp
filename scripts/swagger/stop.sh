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

if ! container_exists "${SWAGGER_CONTAINER_NAME}"; then
  log_info "Swagger UI is already stopped"
  exit 0
fi

log_info "Stopping Swagger UI"
docker rm --force "${SWAGGER_CONTAINER_NAME}" >/dev/null
log_info "Swagger UI stopped"
