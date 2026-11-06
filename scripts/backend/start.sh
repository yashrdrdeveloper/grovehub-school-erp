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
source "${SCRIPT_DIR}/../lib/env.sh"
source "${SCRIPT_DIR}/../config/services.sh"

ensure_runtime_dirs
require_command curl
load_backend_runtime_env

export MAIL_INFRASTRUCTURE_DKIM_STATE_DIR="${MAIL_INFRASTRUCTURE_DKIM_STATE_DIR:-${STATE_DIR}/infrastructure/dkim}"

PID_FILE="${RUNTIME_PID_DIR}/backend.pid"
LOG_FILE="${RUNTIME_LOG_DIR}/backend.log"
BINARY="${RUNTIME_BIN_DIR}/groovehub-platform"
HOT_RELOAD="${BACKEND_HOT_RELOAD:-0}"

if owned_pid_is_running "${PID_FILE}" "${AIR_BIN}"; then
  log_info "Backend is already running with Air hot reload (PID $(cat "${PID_FILE}"))"
  exit 0
fi

if owned_pid_is_running "${PID_FILE}" "${BINARY}"; then
  if [[ "${HOT_RELOAD}" == "1" ]]; then
    log_info "Switching GroveHub backend from one-time binary to Air hot reload"
    stop_owned_process "${PID_FILE}" "${BINARY}" "GroveHub backend"
  else
    log_info "Backend is already running (PID $(cat "${PID_FILE}"))"
    exit 0
  fi
fi

# A stale/missing PID file can happen after a manual 'make backend-dev'
# or a previous shell exit. Recover only an Air process whose executable and
# working directory prove it belongs to this checkout.
if [[ "${HOT_RELOAD}" == "1" ]]; then
  listener_pids="$(listening_tcp_pids "${BACKEND_PORT}")"
  if [[ -n "${listener_pids}" ]]; then
    for listener_pid in ${listener_pids}; do
      supervisor_pid="$(repo_air_supervisor_for_listener "${listener_pid}" "${AIR_BIN}" "${BACKEND_DIR}" || true)"
      if [[ -n "${supervisor_pid}" ]]; then
        printf '%s\n' "${supervisor_pid}" >"${PID_FILE}"
        log_info "Recovered this repository's existing Air backend (PID ${supervisor_pid}); skipping duplicate startup"
        exit 0
      fi
    done
  fi
fi

rm -f -- "${PID_FILE}"

for health_path in "${BACKEND_HEALTH_ENDPOINTS[@]}"; do
  if curl --fail --silent --max-time 1 "${BACKEND_BASE_URL}${health_path}" >/dev/null 2>&1; then
    log_error "Backend endpoint is reachable but no managed GroveHub process could be verified: ${BACKEND_BASE_URL}${health_path}"
    if command -v lsof >/dev/null 2>&1; then
      log_error "Port ${BACKEND_PORT} listeners (inspect before stopping anything):"
      lsof -nP -iTCP:"${BACKEND_PORT}" -sTCP:LISTEN >&2 || true
    fi
    die "Port ${BACKEND_PORT} is occupied. Inspect the listener, stop only the confirmed old process, and rerun make dev-start."
  fi
done

if [[ "${HOT_RELOAD}" == "1" ]]; then
  bash "${SCRIPT_DIR}/setup.sh"

  log_info "Starting GroveHub backend with Air hot reload"
  (
    cd -- "${BACKEND_DIR}"
    nohup "${AIR_BIN}" -c .air.toml >>"${LOG_FILE}" 2>&1 &
    printf '%s\n' "$!" >"${PID_FILE}"
  )
else
  bash "${SCRIPT_DIR}/build.sh"

  log_info "Starting GroveHub backend"
  nohup "${BINARY}" >>"${LOG_FILE}" 2>&1 &
  printf '%s\n' "$!" >"${PID_FILE}"
fi

backend_pid="$(cat "${PID_FILE}")"

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

if [[ "${HOT_RELOAD}" == "1" ]]; then
  log_info "Backend started with Air hot reload (PID ${backend_pid})"
else
  log_info "Backend started (PID ${backend_pid})"
fi
log_info "Backend URL: ${BACKEND_BASE_URL}"
log_info "Backend log: ${LOG_FILE}"
