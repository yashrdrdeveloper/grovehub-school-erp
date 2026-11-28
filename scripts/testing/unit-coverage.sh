#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

require_command go
mode="${1:---strict}"
case "${mode}" in
  --report-only|--strict|--inventory) ;;
  *) die "usage: unit-coverage.sh [--report-only|--strict|--inventory]" ;;
esac

# Single service, multi-service, or all services: the Go runner is shared.
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" coverage \
    -core "${CORE_DIR}" \
    -generator-script "${REPO_ROOT}/scripts/backend/generate-annotations.sh" \
    -reports "${RUNTIME_DIR}/reports/unit-coverage" \
    -backend "${BACKEND_DIR}" \
    -mode "${mode}" \
    -package "${UNIT_PACKAGE:-}" \
    -packages "${UNIT_PACKAGES:-}"
)
