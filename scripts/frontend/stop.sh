#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/process.sh"

stop_owned_process_any "${RUNTIME_PID_DIR}/frontend.pid" "GrooveHub frontend" "${REPO_ROOT}/frontend/node_modules/vite/bin/vite.js"
