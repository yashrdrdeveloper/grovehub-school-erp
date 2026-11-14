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
  go run ./libs/core/annotation/generation/testing/coverage \
    -backend "${BACKEND_DIR}" \
    -mode "${mode}" \
    -package "${UNIT_PACKAGE:-}" \
    -packages "${UNIT_PACKAGES:-}"
)
