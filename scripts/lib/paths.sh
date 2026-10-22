#!/usr/bin/env bash
set -Eeuo pipefail

PATHS_SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

REPO_ROOT="$(
  cd -- "${PATHS_SCRIPT_DIR}/../.." >/dev/null 2>&1
  pwd
)"

BACKEND_DIR="${REPO_ROOT}/backend"
RUNTIME_DIR="${REPO_ROOT}/.runtime"
RUNTIME_BIN_DIR="${RUNTIME_DIR}/bin"
RUNTIME_LOG_DIR="${RUNTIME_DIR}/logs"
RUNTIME_PID_DIR="${RUNTIME_DIR}/pids"

ensure_runtime_dirs() {
  mkdir -p -- "${RUNTIME_BIN_DIR}" "${RUNTIME_LOG_DIR}" "${RUNTIME_PID_DIR}"
}
