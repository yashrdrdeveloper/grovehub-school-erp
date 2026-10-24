#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/container.sh"

HARAKA_CONTAINER_NAME="grovehub-haraka"

prepare_existing_docker_cli
if ! configure_docker_access || ! docker info >/dev/null 2>&1; then
  log_info "Mail Provider transport: STOPPED (container runtime unavailable)"
  exit 0
fi

if ! container_exists "${HARAKA_CONTAINER_NAME}"; then
  log_info "Mail Provider transport: STOPPED"
  exit 0
fi

status="$(docker inspect --format '{{.State.Status}}' "${HARAKA_CONTAINER_NAME}")"
health="$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}n/a{{end}}' "${HARAKA_CONTAINER_NAME}")"
log_info "Mail Provider transport: ${status} (health=${health})"
