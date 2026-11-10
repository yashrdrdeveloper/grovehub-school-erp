#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/go.sh"

require_command go
sync_backend_go_modules
bash "${SCRIPT_DIR}/generate-entities.sh"

# go vet cannot read virtual source paths that are present only in a Go overlay.
# Run compiled generated-code tests with -vet=off, and vet ordinary packages separately.
log_info "Running backend tests with the generated annotation overlay"
(
  cd -- "${BACKEND_DIR}"
  go test -vet=off -overlay="${BACKEND_DIR}/.generated/overlay.json" "$@"
)
