#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/common.sh"

ensure_nats_runtime

if nats_container_running; then
  log_info "NATS JetStream is already running: ${NATS_CONTAINER_NAME}"
  exit 0
fi

if nats_container_exists; then
  docker rm --force "${NATS_CONTAINER_NAME}" >/dev/null
fi

log_info "Starting NATS JetStream"
docker run --detach \
  --name "${NATS_CONTAINER_NAME}" \
  --pull missing \
  --restart unless-stopped \
  --publish "127.0.0.1:${NATS_CLIENT_PORT}:4222" \
  --publish "127.0.0.1:${NATS_MONITOR_PORT}:8222" \
  --volume "${NATS_DATA_VOLUME}:/data" \
  "${NATS_IMAGE}" \
  --jetstream \
  --store_dir=/data \
  --http_port=8222 >/dev/null

if ! wait_for_http "${NATS_HEALTH_URL}" "NATS JetStream" 30; then
  docker logs "${NATS_CONTAINER_NAME}" >&2 || true
  docker rm --force "${NATS_CONTAINER_NAME}" >/dev/null 2>&1 || true
  die "NATS JetStream failed to start"
fi

log_info "NATS client endpoint: nats://127.0.0.1:${NATS_CLIENT_PORT}"
log_info "NATS monitoring endpoint: http://127.0.0.1:${NATS_MONITOR_PORT}"
log_info "NATS persistent volume: ${NATS_DATA_VOLUME}"
