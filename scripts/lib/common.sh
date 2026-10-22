#!/usr/bin/env bash
set -Eeuo pipefail

log_info() {
  printf '[INFO] %s\n' "$*"
}

log_warn() {
  printf '[WARN] %s\n' "$*" >&2
}

log_error() {
  printf '[ERROR] %s\n' "$*" >&2
}

die() {
  log_error "$*"
  exit 1
}

require_command() {
  local command_name="$1"
  command -v "${command_name}" >/dev/null 2>&1 || die "Required command not found: ${command_name}"
}

wait_for_http() {
  local url="$1"
  local label="$2"
  local attempts="${3:-30}"
  local attempt

  require_command curl

  for ((attempt = 1; attempt <= attempts; attempt++)); do
    if curl --fail --silent --show-error --max-time 2 "${url}" >/dev/null 2>&1; then
      log_info "${label} is ready: ${url}"
      return 0
    fi
    sleep 1
  done

  return 1
}
