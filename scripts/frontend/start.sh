#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/process.sh"

FRONTEND_DIR="${REPO_ROOT}/frontend"
PID_FILE="${RUNTIME_PID_DIR}/frontend.pid"
LOG_FILE="${RUNTIME_LOG_DIR}/frontend.log"
FRONTEND_PORT="${FRONTEND_PORT:-4173}"
FRONTEND_URL="http://127.0.0.1:${FRONTEND_PORT}"
VITE_ENTRY="${FRONTEND_DIR}/node_modules/vite/bin/vite.js"

ensure_runtime_dirs
require_command node
require_command curl

if owned_pid_is_running "${PID_FILE}" "${VITE_ENTRY}"; then
  log_info "GrooveHub frontend is already running (PID $(cat "${PID_FILE}"))"
  exit 0
fi

rm -f -- "${PID_FILE}"

if curl --silent --fail --max-time 1 "${FRONTEND_URL}" >/dev/null 2>&1; then
  die "Port ${FRONTEND_PORT} is occupied by an unmanaged HTTP process; refusing to adopt it"
fi

if [[ ! -f "${VITE_ENTRY}" ]]; then
  bash "${SCRIPT_DIR}/setup.sh"
fi

log_info "Starting GrooveHub frontend with Vite"
(
  cd -- "${FRONTEND_DIR}"
  nohup node "${VITE_ENTRY}" --host 127.0.0.1 --port "${FRONTEND_PORT}" --strictPort >>"${LOG_FILE}" 2>&1 &
  printf '%s\n' "$!" >"${PID_FILE}"
)

frontend_pid="$(cat "${PID_FILE}")"
if ! wait_for_http "${FRONTEND_URL}" "Frontend" 30; then
  log_error "Frontend failed readiness: ${LOG_FILE}"
  tail -n 35 "${LOG_FILE}" >&2 || true
  if pid_is_running "${frontend_pid}"; then kill -TERM "${frontend_pid}" || true; fi
  rm -f -- "${PID_FILE}"
  exit 1
fi
log_info "Frontend URL: ${FRONTEND_URL}"
log_info "Frontend log: ${LOG_FILE}"
