#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/container.sh"
source "${SCRIPT_DIR}/../lib/env.sh"

COMPOSE_FILE="${REPO_ROOT}/deployments/docker/docker-compose/mail-provider.yaml"
HARAKA_CONTAINER_NAME="grovehub-haraka"

require_docker
load_backend_runtime_env

export MAIL_PROVIDER_SMTP_PORT="${MAIL_PROVIDER_SMTP_PORT:-2525}"
export HARAKA_MEMORY_LIMIT="${HARAKA_MEMORY_LIMIT:-512m}"

[[ -f "${COMPOSE_FILE}" ]] || die "Mail Provider compose file not found: ${COMPOSE_FILE}"

log_info "Starting GroveHub Mail Provider transport (Haraka)"
docker_compose -f "${COMPOSE_FILE}" up --detach --build

for _ in {1..30}; do
  health="$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "${HARAKA_CONTAINER_NAME}" 2>/dev/null || true)"
  if [[ "${health}" == "healthy" ]]; then
    log_info "Haraka is ready on 127.0.0.1:${MAIL_PROVIDER_SMTP_PORT}"
    exit 0
  fi
  if [[ "${health}" == "unhealthy" || "${health}" == "exited" ]]; then
    docker logs "${HARAKA_CONTAINER_NAME}" >&2 || true
    die "Haraka failed to start"
  fi
  sleep 1
done

docker logs "${HARAKA_CONTAINER_NAME}" >&2 || true
die "Haraka did not become healthy"
