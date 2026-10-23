#!/usr/bin/env bash
set -Eeuo pipefail

load_backend_runtime_env() {
  local env_file="${REPO_ROOT}/.env"
  [[ -f "${env_file}" ]] || return 0

  set +u
  # shellcheck disable=SC1090
  source "${env_file}"
  set -u

  if [[ -n "${GROVEHUB_LOG_LEVELS+x}" ]]; then
    export GROVEHUB_LOG_LEVELS
  fi
}
