#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

bash "${SCRIPT_DIR}/../infrastructure/nats/start.sh"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/env.sh"

# Export Gmail/Outlook provider variables before starting Air.
# Provider connections stay disabled when the optional OAuth credentials
# are not configured, rather than failing the entire developer startup.
load_backend_runtime_env
# Development Vite proxy has an explicit loopback-only origin policy.
export IDENTITY_ALLOWED_ORIGIN="${IDENTITY_ALLOWED_ORIGIN:-http://127.0.0.1:4173}"

# init all already includes the Integration schema.
bash "${SCRIPT_DIR}/../database/migrate.sh" init all
BACKEND_HOT_RELOAD=1 bash "${SCRIPT_DIR}/../backend/start.sh"
bash "${SCRIPT_DIR}/../swagger/start.sh"
bash "${SCRIPT_DIR}/../frontend/start.sh"
