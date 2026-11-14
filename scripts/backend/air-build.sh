#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
# Module dependencies are synchronized once by dev.sh; rebuilding on every
# source edit must not rewrite go.mod or go.sum and re-trigger Air.

require_command go
ensure_runtime_dirs
bash "${SCRIPT_DIR}/generate-annotations.sh"

OUTPUT="${RUNTIME_BIN_DIR}/groovehub-dev"

log_info "Building GroveHub backend for Air"
(
  cd -- "${BACKEND_DIR}"
  go build -overlay="${BACKEND_DIR}/.generated/overlay.json" -o "${OUTPUT}" ./cmd/platform
)
