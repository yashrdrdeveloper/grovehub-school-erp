#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/common.sh"

prepare_existing_docker_cli
if ! configure_docker_access || ! docker info >/dev/null 2>&1; then
  log_warn "NATS JetStream: container runtime unavailable"
  exit 0
fi

if ! nats_container_running; then
  log_warn "NATS JetStream: DOWN"
  log_info "NATS persistent volume: ${NATS_DATA_VOLUME}"
  exit 0
fi

log_info "NATS JetStream: UP"
log_info "NATS container: ${NATS_CONTAINER_NAME}"
log_info "NATS client endpoint: nats://127.0.0.1:${NATS_CLIENT_PORT}"
log_info "NATS monitoring endpoint: http://127.0.0.1:${NATS_MONITOR_PORT}"
log_info "NATS persistent volume: ${NATS_DATA_VOLUME}"

if ! wait_for_http "${NATS_HEALTH_URL}" "NATS JetStream health" 1; then
  log_warn "NATS JetStream health: DOWN"
fi
