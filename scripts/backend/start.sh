#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/process.sh"
source "${SCRIPT_DIR}/../config/services.sh"

ensure_runtime_dirs
require_command curl

PID_FILE="${RUNTIME_PID_DIR}/backend.pid"
LOG_FILE="${RUNTIME_LOG_DIR}/backend.log"
BINARY="${RUNTIME_BIN_DIR}/groovehub-platform"

if owned_pid_is_running "${PID_FILE}" "${BINARY}"; then
  log_info "Backend is already running (PID $(cat "${PID_FILE}"))"
  exit 0
fi

rm -f -- "${PID_FILE}"

for health_path in "${BACKEND_HEALTH_ENDPOINTS[@]}"; do
  if curl --fail --silent --max-time 1 "${BACKEND_BASE_URL}${health_path}" >/dev/null 2>&1; then
    die "Backend endpoint is already reachable but is not owned by this repository: ${BACKEND_BASE_URL}${health_path}"
  fi
done

bash "${SCRIPT_DIR}/build.sh"

log_info "Starting GroveHub backend"
nohup "${BINARY}" >>"${LOG_FILE}" 2>&1 &
backend_pid=$!
printf '%s\n' "${backend_pid}" >"${PID_FILE}"

for health_path in "${BACKEND_HEALTH_ENDPOINTS[@]}"; do
  if ! wait_for_http "${BACKEND_BASE_URL}${health_path}" "Backend health" 30; then
    log_error "Backend failed readiness check: ${health_path}"
    tail -n 50 "${LOG_FILE}" >&2 || true
    if pid_is_running "${backend_pid}"; then
      kill -TERM "${backend_pid}" >/dev/null 2>&1 || true
    fi
    rm -f -- "${PID_FILE}"
    exit 1
  fi
done

log_info "Backend started (PID ${backend_pid})"
log_info "Backend URL: ${BACKEND_BASE_URL}"
log_info "Backend log: ${LOG_FILE}"
