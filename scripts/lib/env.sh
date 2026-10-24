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
    MAIL_PROVIDER_SMTP_HOST \
    MAIL_PROVIDER_SMTP_PORT \
    HARAKA_MEMORY_LIMIT \
    MAIL_PROVIDER_DKIM_SELECTOR \
    MAIL_PROVIDER_PUBLIC_IP \
    MAIL_PROVIDER_HARAKA_HOSTNAME \
    CLOUDFLARE_API_TOKEN \
    CLOUDFLARE_ZONE_ID \
    OTEL_SERVICE_VERSION \
    OTEL_SERVICE_INSTANCE_ID \
    OTEL_DEPLOYMENT_ENVIRONMENT; do
    if [[ -n "${!name+x}" ]]; then
      export "${name}"
    fi
  done
}

load_mail_provider_state_env() {
  local env_file="${STATE_DIR}/mailProvider/runtime.env"
  [[ -f "${env_file}" ]] || return 0

  set +u
  # shellcheck disable=SC1090
  source "${env_file}"
  set -u

  local name
  for name in MAIL_PROVIDER_PUBLIC_IP MAIL_PROVIDER_HARAKA_HOSTNAME; do
    if [[ -n "${!name+x}" ]]; then
      export "${name}"
    fi
  done
}
