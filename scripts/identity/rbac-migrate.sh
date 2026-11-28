#!/usr/bin/env bash
set -Eeuo pipefail
set +x
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/env.sh"
load_backend_runtime_env
# Use absolute --input paths; go -C runs from the backend directory.
exec go -C "${BACKEND_DIR}" run ./cmd/identity-rbac-migrate "$@"
