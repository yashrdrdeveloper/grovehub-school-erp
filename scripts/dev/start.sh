#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

bash "${SCRIPT_DIR}/../mailProvider/start.sh"
bash "${SCRIPT_DIR}/../database/migrate.sh" init all
BACKEND_HOT_RELOAD=1 bash "${SCRIPT_DIR}/../backend/start.sh"
bash "${SCRIPT_DIR}/../swagger/start.sh"
