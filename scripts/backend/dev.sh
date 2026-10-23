#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/toolchain.sh"
source "${SCRIPT_DIR}/../lib/go.sh"
source "${SCRIPT_DIR}/../config/services.sh"

ensure_runtime_dirs
ensure_air
sync_backend_go_modules

log_info "Running GroveHub backend with Air hot reload"
log_info "Backend URL: ${BACKEND_BASE_URL}"
log_info "Air version: ${AIR_VERSION}"

cd -- "${BACKEND_DIR}"
exec "${AIR_BIN}" -c .air.toml
