#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../config/services.sh"

require_command go

log_info "Running GroveHub backend in foreground"
log_info "Backend URL: ${BACKEND_BASE_URL}"
cd -- "${BACKEND_DIR}"
exec go run ./cmd/platform
