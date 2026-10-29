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
load_mail_provider_state_env

export MAIL_PROVIDER_SMTP_PORT="${MAIL_PROVIDER_SMTP_PORT:-2525}"
export HARAKA_MEMORY_LIMIT="${HARAKA_MEMORY_LIMIT:-512m}"
export MAIL_PROVIDER_HARAKA_HOSTNAME="${MAIL_PROVIDER_HARAKA_HOSTNAME:-mail.local.grovehub}"
export MAIL_PROVIDER_DKIM_STATE_DIR="${MAIL_PROVIDER_DKIM_STATE_DIR:-${STATE_DIR}/mailProvider/dkim}"

[[ -f "${COMPOSE_FILE}" ]] || die "Mail Provider compose file not found: ${COMPOSE_FILE}"

mkdir -p -- "${MAIL_PROVIDER_DKIM_STATE_DIR}"
chmod 700 "${MAIL_PROVIDER_DKIM_STATE_DIR}" 2>/dev/null || true

has_persisted_dkim_key() {
  find "${MAIL_PROVIDER_DKIM_STATE_DIR}" -type f -name private -size +0c -print -quit 2>/dev/null | grep -q .
}

migrate_existing_dkim_key() {
  if has_persisted_dkim_key; then
    log_info "Reusing persistent Haraka DKIM key from ${MAIL_PROVIDER_DKIM_STATE_DIR}"
    return 0
  fi

  if ! docker inspect "${HARAKA_CONTAINER_NAME}" >/dev/null 2>&1; then
    return 0
  fi

  log_info "Migrating existing Haraka DKIM key into persistent host storage"
  if docker cp "${HARAKA_CONTAINER_NAME}:/opt/haraka/config/dkim/." "${MAIL_PROVIDER_DKIM_STATE_DIR}/" >/dev/null 2>&1 && has_persisted_dkim_key; then
    find "${MAIL_PROVIDER_DKIM_STATE_DIR}" -type f -name private -exec chmod 600 {} + 2>/dev/null || true
    log_info "Existing Haraka DKIM key migrated to ${MAIL_PROVIDER_DKIM_STATE_DIR}"
    return 0
  fi

  log_info "No existing Haraka DKIM key found; one will be generated when the mail domain is prepared"
}

migrate_existing_dkim_key

log_info "Starting GroveHub Mail Provider transport (Haraka)"
docker_compose -f "${COMPOSE_FILE}" up --detach --build

for _ in {1..30}; do
  health="$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "${HARAKA_CONTAINER_NAME}" 2>/dev/null || true)"
  if [[ "${health}" == "healthy" ]]; then
    log_info "Haraka is ready on 127.0.0.1:${MAIL_PROVIDER_SMTP_PORT}"
    log_info "Haraka public hostname: ${MAIL_PROVIDER_HARAKA_HOSTNAME}"
    log_info "Haraka DKIM state: ${MAIL_PROVIDER_DKIM_STATE_DIR}"
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
