#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/go.sh"

require_command go
ensure_runtime_dirs

OUTPUT="${RUNTIME_BIN_DIR}/groovehub-platform"

sync_backend_go_modules
bash "${SCRIPT_DIR}/generate-entities.sh"

log_info "Building GroveHub backend"
(
  cd -- "${BACKEND_DIR}"
  go build -overlay="${BACKEND_DIR}/.generated/overlay.json" -o "${OUTPUT}" ./cmd/platform
)

log_info "Backend binary ready: ${OUTPUT}"
