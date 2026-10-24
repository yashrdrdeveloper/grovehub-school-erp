#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/container.sh"

COMPOSE_FILE="${REPO_ROOT}/deployments/docker/docker-compose/mail-provider.yaml"

prepare_existing_docker_cli
if ! configure_docker_access || ! docker info >/dev/null 2>&1; then
  log_info "Container runtime is not running; Mail Provider transport is already unavailable"
  exit 0
fi

log_info "Stopping GroveHub Mail Provider transport"
docker_compose -f "${COMPOSE_FILE}" down
