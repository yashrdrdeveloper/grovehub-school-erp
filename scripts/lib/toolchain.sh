#!/usr/bin/env bash
set -Eeuo pipefail

TOOLCHAIN_LIB_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${TOOLCHAIN_LIB_DIR}/../config/toolchain.sh"

AIR_BIN="${TOOLS_BIN_DIR}/air"
GOOSE_BIN="${TOOLS_BIN_DIR}/goose"

ensure_air() {
  require_command go
  ensure_tool_dirs

  local version_file="${STATE_DIR}/air.version"
  local installed_version=""

  if [[ -f "${version_file}" ]]; then
    installed_version="$(cat "${version_file}")"
  fi

  if [[ -x "${AIR_BIN}" && "${installed_version}" == "${AIR_VERSION}" ]]; then
    return 0
  fi

  log_info "Installing Air ${AIR_VERSION} into ${TOOLS_BIN_DIR}"
  GOBIN="${TOOLS_BIN_DIR}" go install "${AIR_PACKAGE}@${AIR_VERSION}"
  printf '%s\n' "${AIR_VERSION}" >"${version_file}"
}

ensure_goose() {
  require_command go
  ensure_tool_dirs

  local version_file="${STATE_DIR}/goose.version"
  local installed_version=""

  if [[ -f "${version_file}" ]]; then
    installed_version="$(cat "${version_file}")"
  fi

  if [[ -x "${GOOSE_BIN}" && "${installed_version}" == "${GOOSE_VERSION}" ]]; then
    return 0
  fi

  log_info "Installing Goose ${GOOSE_VERSION} into ${TOOLS_BIN_DIR}"
  GOBIN="${TOOLS_BIN_DIR}" go install "${GOOSE_PACKAGE}@${GOOSE_VERSION}"
  printf '%s\n' "${GOOSE_VERSION}" >"${version_file}"
}
