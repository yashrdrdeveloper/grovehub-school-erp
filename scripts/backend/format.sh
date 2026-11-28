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
  core)
    path="${CORE_DIR}"
    ;;
  integration)
    path="${BACKEND_DIR}/services/integration"
    ;;
  foundation)
    path="${BACKEND_DIR}/services/foundation"
    ;;
  foundation-repository)
    path="${BACKEND_DIR}/services/foundation/internal/app/repository"
    ;;
  *)
    die "Unsupported format scope: ${scope}. Expected all, core, foundation, foundation-repository, or integration."
    ;;
esac

mode="${2:-}"
case "${mode}" in
  --check)
    go run "${SCRIPT_DIR}/groovefmt/main.go" -check "${path}"
    if [[ "$scope" == all ]]; then go run "${SCRIPT_DIR}/groovefmt/main.go" -check "${CORE_DIR}"; fi
    log_info "GroveHub formatting lint passed under ${path}"
    ;;
  "")
    go run "${SCRIPT_DIR}/groovefmt/main.go" "${path}"
    if [[ "$scope" == all ]]; then go run "${SCRIPT_DIR}/groovefmt/main.go" "${CORE_DIR}"; fi
    log_info "Formatted Go code with GroveHub style under ${path}"
    ;;
  *)
    die "Unsupported formatting mode: ${mode}"
    ;;
esac
