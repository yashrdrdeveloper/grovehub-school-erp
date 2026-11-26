#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/process.sh"
source "${SCRIPT_DIR}/tunnel.sh"

FRONTEND_DIR="${REPO_ROOT}/frontend"
PID_FILE="${RUNTIME_PID_DIR}/frontend.pid"
LOG_FILE="${RUNTIME_LOG_DIR}/frontend.log"
FRONTEND_PORT="${FRONTEND_PORT:-4173}"
FRONTEND_URL="http://127.0.0.1:${FRONTEND_PORT}"
VITE_ENTRY="${FRONTEND_DIR}/node_modules/vite/bin/vite.js"

ensure_runtime_dirs
require_command node
require_command curl

if ! owned_pid_is_running "${PID_FILE}" "${VITE_ENTRY}" && curl --silent --fail --max-time 1 "${FRONTEND_URL}" >/dev/null 2>&1; then
  die "Port ${FRONTEND_PORT} is occupied by an unmanaged HTTP process; refusing to adopt it"
fi

if [[ ! -f "${VITE_ENTRY}" ]]; then
  bash "${SCRIPT_DIR}/setup.sh"
fi

prepare_frontend_tunnel "${FRONTEND_URL}"
TUNNEL_URL="$(read_frontend_tunnel_url || true)"
VITE_TUNNEL_FILE="${RUNTIME_PID_DIR}/frontend-tunnel-vite.url"

if owned_pid_is_running "${PID_FILE}" "${VITE_ENTRY}"; then
  previous_url="$(cat "${VITE_TUNNEL_FILE}" 2>/dev/null || true)"
  if [[ "${previous_url}" == "${TUNNEL_URL}" ]]; then
    log_info "GrooveHub frontend is already running (PID $(cat "${PID_FILE}"))"
    log_info "Frontend URL: ${FRONTEND_URL}"
    report_frontend_tunnel
    exit 0
  fi
  log_info "Restarting managed frontend to apply the changed tunnel hostname"
  stop_owned_process "${PID_FILE}" "${VITE_ENTRY}" "GrooveHub frontend"
fi

rm -f -- "${PID_FILE}"

if curl --silent --fail --max-time 1 "${FRONTEND_URL}" >/dev/null 2>&1; then
  die "Port ${FRONTEND_PORT} is occupied by an unmanaged HTTP process; refusing to adopt it"
fi

log_info "Starting GrooveHub frontend with Vite"
(
  cd -- "${FRONTEND_DIR}"
  export GROVEHUB_FRONTEND_TUNNEL_URL="${TUNNEL_URL}"
  export GROVEHUB_FRONTEND_LOCAL_ORIGIN="${FRONTEND_URL}"
  nohup node "${VITE_ENTRY}" --host 127.0.0.1 --port "${FRONTEND_PORT}" --strictPort >>"${LOG_FILE}" 2>&1 &
  printf '%s\n' "$!" >"${PID_FILE}"
)

frontend_pid="$(cat "${PID_FILE}")"
if ! wait_for_http "${FRONTEND_URL}" "Frontend" 30; then
  log_error "Frontend failed readiness: ${LOG_FILE}"
  tail -n 35 "${LOG_FILE}" >&2 || true
  if pid_is_running "${frontend_pid}"; then kill -TERM "${frontend_pid}" || true; fi
  rm -f -- "${PID_FILE}" "${VITE_TUNNEL_FILE}"
  stop_frontend_tunnel
  exit 1
fi
log_info "Frontend URL: ${FRONTEND_URL}"
log_info "Frontend log: ${LOG_FILE}"
printf '%s\n' "${TUNNEL_URL}" >"${VITE_TUNNEL_FILE}"
report_frontend_tunnel
