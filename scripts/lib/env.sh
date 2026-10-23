#!/usr/bin/env bash
set -Eeuo pipefail

load_backend_runtime_env() {
  local env_file="${REPO_ROOT}/.env"
  [[ -f "${env_file}" ]] || return 0

  set +u
  # shellcheck disable=SC1090
  source "${env_file}"
  set -u

  local name
  for name in \
    DATABASE_URL \
    GROVEHUB_LOG_LEVELS \
    GROVEHUB_TRACE_SAMPLE_RATIO \
    GROVEHUB_METRIC_EXPORT_INTERVAL \
    OTEL_SERVICE_VERSION \
    OTEL_SERVICE_INSTANCE_ID \
    OTEL_DEPLOYMENT_ENVIRONMENT; do
    if [[ -n "${!name+x}" ]]; then
      export "${name}"
    fi
  done
}
