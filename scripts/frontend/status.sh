#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/process.sh"

ENTRY="${REPO_ROOT}/frontend/node_modules/vite/bin/vite.js"
if owned_pid_is_running "${RUNTIME_PID_DIR}/frontend.pid" "${ENTRY}"; then
  log_info "Frontend running (PID $(cat "${RUNTIME_PID_DIR}/frontend.pid")) — http://127.0.0.1:${FRONTEND_PORT:-4173}"
else
  log_info "Frontend stopped"
fi
