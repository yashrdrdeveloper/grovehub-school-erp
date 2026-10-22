#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

"${SCRIPT_DIR}/../backend/start.sh"
"${SCRIPT_DIR}/../swagger/start.sh"
