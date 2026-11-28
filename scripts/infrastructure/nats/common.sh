#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/container.sh"

NATS_CONTAINER_NAME="${NATS_CONTAINER_NAME:-grovehub-nats}"
NATS_IMAGE="${NATS_IMAGE:-nats:2.11.9-alpine}"
NATS_CLIENT_PORT="${NATS_CLIENT_PORT:-4222}"
NATS_MONITOR_PORT="${NATS_MONITOR_PORT:-8222}"
NATS_DATA_VOLUME="${NATS_DATA_VOLUME:-grovehub-nats-data}"
NATS_HEALTH_URL="http://127.0.0.1:${NATS_MONITOR_PORT}/healthz"

ensure_nats_runtime() {
  require_docker
}

nats_container_exists() {
  container_exists "${NATS_CONTAINER_NAME}"
}

nats_container_running() {
  container_is_running "${NATS_CONTAINER_NAME}"
}
