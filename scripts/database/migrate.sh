#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/env.sh"
source "${SCRIPT_DIR}/../lib/toolchain.sh"

command_name="${1:-init}"
service_name="${2:-all}"

case "${command_name}" in
  init|reset|status)
    ;;
  *)
    die "Unsupported database schema command: ${command_name}. Expected init, reset, or status."
    ;;
esac

case "${service_name}" in
  all|foundation|identity|integration)
    ;;
  *)
    die "Unsupported database service: ${service_name}. Expected all, foundation, identity, or integration."
    ;;
esac

load_backend_runtime_env

[[ -n "${DATABASE_URL:-}" ]] || die "DATABASE_URL is required. Add the Aiven Service URI to the repository root .env file."

ensure_goose

schema_state_table() {
  case "$1" in
    foundation) printf '%s\n' "foundation_schema_state" ;;
    identity) printf '%s\n' "identity_schema_state" ;;
    integration) printf '%s\n' "integration_schema_state" ;;
    infrastructure) printf '%s\n' "infrastructure_schema_state" ;;
  esac
}

run_service_schema() {
  local service="$1"
  local action="$2"
  local schema_dir="${BACKEND_DIR}/services/${service}/migrations"
  if [[ "${service}" == "foundation" ]]; then
    schema_dir="${BACKEND_DIR}/services/foundation/internal/resources/migrations"
  fi
  local state_table
  state_table="$(schema_state_table "${service}")"

  if [[ ! -d "${schema_dir}" ]]; then
    log_info "Skipping ${service} database schema command (${action}): no migrations directory"
    return 0
  fi

  local -a migration_files=()
  shopt -s nullglob
  migration_files=("${schema_dir}"/*.sql)
  shopt -u nullglob

  if (( ${#migration_files[@]} == 0 )); then
    log_info "Skipping ${service} database schema command (${action}): no SQL migrations"
    return 0
  fi

  log_info "Running ${service} database schema command: ${action}"

  GOOSE_DRIVER="postgres" \
  GOOSE_DBSTRING="${DATABASE_URL}" \
  GOOSE_MIGRATION_DIR="${schema_dir}" \
  GOOSE_TABLE="${state_table}" \
    "${GOOSE_BIN}" -env=none "${action}"
}

init_service() {
  run_service_schema "$1" up
}

reset_service() {
  local service="$1"
  run_service_schema "${service}" reset
  init_service "${service}"
}

status_service() {
  run_service_schema "$1" status
}

if [[ "${service_name}" != "all" ]]; then
  case "${command_name}" in
    init) init_service "${service_name}" ;;
    reset) reset_service "${service_name}" ;;
    status) status_service "${service_name}" ;;
  esac
  exit 0
fi

case "${command_name}" in
  init)
    init_service foundation
    init_service identity
    init_service integration
    ;;
  reset)
    run_service_schema integration reset
    run_service_schema identity reset
    run_service_schema foundation reset
    init_service foundation
    init_service identity
    init_service integration
    ;;
  status)
    status_service foundation
    status_service identity
    status_service integration
    ;;
esac
