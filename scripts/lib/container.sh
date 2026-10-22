#!/usr/bin/env bash
set -Eeuo pipefail

require_docker() {
  require_command docker
  docker info >/dev/null 2>&1 || die "Docker daemon is not available"
}

container_exists() {
  local container_name="$1"
  docker inspect "${container_name}" >/dev/null 2>&1
}

container_is_running() {
  local container_name="$1"
  [[ "$(docker inspect --format '{{.State.Running}}' "${container_name}" 2>/dev/null || true)" == "true" ]]
}
