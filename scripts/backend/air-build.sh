#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/go.sh"

sync_backend_go_modules

OUTPUT="${RUNTIME_BIN_DIR}/groovehub-dev"

log_info "Building GroveHub backend for Air"
(
  cd -- "${BACKEND_DIR}"
  go build -o "${OUTPUT}" ./cmd/platform
)
