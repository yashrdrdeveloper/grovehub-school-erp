#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/process.sh"

PID_FILE="${RUNTIME_PID_DIR}/backend.pid"
BINARY="${RUNTIME_BIN_DIR}/groovehub-platform"

stop_owned_process "${PID_FILE}" "${BINARY}" "GroveHub backend"
