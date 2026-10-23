#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/process.sh"
source "${SCRIPT_DIR}/../lib/toolchain.sh"
source "${SCRIPT_DIR}/../config/services.sh"

PID_FILE="${RUNTIME_PID_DIR}/backend.pid"
BINARY="${RUNTIME_BIN_DIR}/groovehub-platform"

if owned_pid_is_running "${PID_FILE}" "${AIR_BIN}"; then
  log_info "Backend process: UP with Air hot reload (PID $(cat "${PID_FILE}"))"
elif owned_pid_is_running "${PID_FILE}" "${BINARY}"; then
  log_info "Backend process: UP (PID $(cat "${PID_FILE}"))"
else
  log_warn "Backend process: DOWN"
fi

for health_path in "${BACKEND_HEALTH_ENDPOINTS[@]}"; do
  if curl --fail --silent --max-time 2 "${BACKEND_BASE_URL}${health_path}" >/dev/null 2>&1; then
    log_info "Health UP: ${BACKEND_BASE_URL}${health_path}"
  else
    log_warn "Health DOWN: ${BACKEND_BASE_URL}${health_path}"
  fi
done
