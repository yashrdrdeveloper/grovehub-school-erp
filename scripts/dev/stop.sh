#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/container.sh"

bash "${SCRIPT_DIR}/../swagger/stop.sh"
bash "${SCRIPT_DIR}/../backend/stop.sh"
bash "${SCRIPT_DIR}/../mailProvider/stop.sh"

COLLECTOR_NAME="grovehub-otel-collector"
prepare_existing_docker_cli
if configure_docker_access && docker info >/dev/null 2>&1; then
  if container_exists "${COLLECTOR_NAME}"; then
    docker rm --force "${COLLECTOR_NAME}" >/dev/null
  fi
fi
