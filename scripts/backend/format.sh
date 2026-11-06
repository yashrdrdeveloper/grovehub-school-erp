#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

require_command gofmt

scope="${1:-all}"
case "${scope}" in
  all)
    path="${BACKEND_DIR}"
    ;;
  integration)
    path="${BACKEND_DIR}/services/integration"
    ;;
  *)
    die "Unsupported format scope: ${scope}. Expected all or integration."
    ;;
esac

# Use a null-delimited find pipeline to handle unusual file names safely.
find "${path}" -type f -name '*.go' -print0 | xargs -0 gofmt -w
log_info "Formatted Go code under ${path}"
