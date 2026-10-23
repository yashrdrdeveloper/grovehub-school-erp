#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

require_command go
ensure_runtime_dirs

OUTPUT="${RUNTIME_BIN_DIR}/groovehub-platform"

log_info "Resolving GroveHub backend Go modules"
(
  cd -- "${BACKEND_DIR}"
  go mod download
)

log_info "Building GroveHub backend"
(
  cd -- "${BACKEND_DIR}"
  go build -o "${OUTPUT}" ./cmd/platform
)

log_info "Backend binary ready: ${OUTPUT}"
