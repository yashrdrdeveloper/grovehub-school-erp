#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

bash "${SCRIPT_DIR}/../backend/status.sh"
bash "${SCRIPT_DIR}/../mailProvider/status.sh"
bash "${SCRIPT_DIR}/../swagger/status.sh"
