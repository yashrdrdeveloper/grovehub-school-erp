#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

if [[ ! -f "${REPO_ROOT}/frontend/projects/school/node_modules/.bin/vite" || ! -d "${REPO_ROOT}/frontend/projects/school/node_modules/@groovehub/base" ]]; then
  bash "${SCRIPT_DIR}/setup.sh"
fi
cd -- "${REPO_ROOT}/frontend/projects/school"
node node_modules/typescript/bin/tsc --noEmit
node node_modules/vite/bin/vite.js build
