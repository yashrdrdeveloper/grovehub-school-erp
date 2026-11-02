#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/common.sh"

prepare_existing_docker_cli
if ! configure_docker_access || ! docker info >/dev/null 2>&1; then
  log_warn "Container runtime is not available; NATS stop skipped"
  exit 0
fi

if ! nats_container_exists; then
  log_info "NATS JetStream is not running"
  exit 0
fi

log_info "Stopping NATS JetStream"
docker rm --force "${NATS_CONTAINER_NAME}" >/dev/null
log_info "NATS persistent volume preserved: ${NATS_DATA_VOLUME}"
