#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

require_command go

scope="${1:-all}"
case "${scope}" in
  all)
    path="${BACKEND_DIR}"
    ;;
  integration)
    path="${BACKEND_DIR}/services/integration"
    ;;
  foundation)
    path="${BACKEND_DIR}/services/foundation"
    ;;
  foundation-repository)
    path="${BACKEND_DIR}/services/foundation/repository"
    ;;
  *)
    die "Unsupported format scope: ${scope}. Expected all, foundation, foundation-repository, or integration."
    ;;
esac

mode="${2:-}"
case "${mode}" in
  --check)
    go run "${SCRIPT_DIR}/groovefmt/main.go" -check "${path}"
    log_info "GroveHub formatting lint passed under ${path}"
    ;;
  "")
    go run "${SCRIPT_DIR}/groovefmt/main.go" "${path}"
    log_info "Formatted Go code with GroveHub style under ${path}"
    ;;
  *)
    die "Unsupported formatting mode: ${mode}"
    ;;
esac
