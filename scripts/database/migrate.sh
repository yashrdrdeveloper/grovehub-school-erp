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

command_name="${1:-status}"
service_name="${2:-all}"

case "${command_name}" in
  status|up|down)
    ;;
  *)
    die "Unsupported database migration command: ${command_name}. Expected status, up, or down."
    ;;
esac

case "${service_name}" in
  all|foundation|identity)
    ;;
  *)
    die "Unsupported migration service: ${service_name}. Expected all, foundation, or identity."
    ;;
esac

load_backend_runtime_env

[[ -n "${DATABASE_URL:-}" ]] || die "DATABASE_URL is required. Add the Aiven Service URI to the repository root .env file."

ensure_goose

run_service_migration() {
  local service="$1"
  local migrations_dir="${BACKEND_DIR}/services/${service}/migrations"

  [[ -d "${migrations_dir}" ]] || die "Migrations directory not found for ${service}: ${migrations_dir}"

  log_info "Running ${service} database migration command: ${command_name}"

  GOOSE_DRIVER="postgres" \
  GOOSE_DBSTRING="${DATABASE_URL}" \
  GOOSE_MIGRATION_DIR="${migrations_dir}" \
  GOOSE_TABLE="${service}_schema_migrations" \
    "${GOOSE_BIN}" -env=none "${command_name}"
}

if [[ "${service_name}" != "all" ]]; then
  run_service_migration "${service_name}"
  exit 0
fi

if [[ "${command_name}" == "down" ]]; then
  run_service_migration "identity"
  run_service_migration "foundation"
  exit 0
fi

run_service_migration "foundation"
run_service_migration "identity"
