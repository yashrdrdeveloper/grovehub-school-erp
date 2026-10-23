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

ensure_air
sync_backend_go_modules

log_info "Backend development tools are ready"
log_info "Air: ${AIR_VERSION} (${AIR_BIN})"
