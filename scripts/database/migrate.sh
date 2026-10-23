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

case "${command_name}" in
  status|up|down)
    ;;
  *)
    die "Unsupported database migration command: ${command_name}. Expected status, up, or down."
    ;;
esac

load_backend_runtime_env

[[ -n "${DATABASE_URL:-}" ]] || die "DATABASE_URL is required. Add the Aiven Service URI to the repository root .env file."

ensure_goose

migrations_dir="${BACKEND_DIR}/services/foundation/migrations"
[[ -d "${migrations_dir}" ]] || die "Foundation migrations directory not found: ${migrations_dir}"

export GOOSE_DRIVER="postgres"
export GOOSE_DBSTRING="${DATABASE_URL}"
export GOOSE_MIGRATION_DIR="${migrations_dir}"
export GOOSE_TABLE="foundation_schema_migrations"

log_info "Running Foundation database migration command: ${command_name}"
"${GOOSE_BIN}" -env=none "${command_name}"
