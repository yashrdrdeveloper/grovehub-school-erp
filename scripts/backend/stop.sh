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

PID_FILE="${RUNTIME_PID_DIR}/backend.pid"
BINARY="${RUNTIME_BIN_DIR}/groovehub-platform"

stop_owned_process_any \
  "${PID_FILE}" \
  "GroveHub backend" \
  "${BINARY}" \
  "${AIR_BIN}"
