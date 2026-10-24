#!/usr/bin/env bash
set -Eeuo pipefail

JSON_LIB_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

# shellcheck source=common.sh
source "${JSON_LIB_DIR}/common.sh"
# shellcheck source=container.sh
source "${JSON_LIB_DIR}/container.sh"

refresh_windows_winget_path() {
  if [[ -z "${LOCALAPPDATA:-}" ]] || ! command -v cygpath >/dev/null 2>&1; then
    return 0
  fi

  local winget_links
  winget_links="$(cygpath -u "${LOCALAPPDATA}\\Microsoft\\WinGet\\Links")"
  if [[ -d "${winget_links}" ]]; then
    export PATH="${PATH}:${winget_links}"
  fi
}

ensure_jq() {
  if command -v jq >/dev/null 2>&1; then
    return 0
  fi

  case "$(uname -s)" in
    Darwin)
      ensure_homebrew
      log_info "Installing jq"
      brew install jq
      ;;
    MINGW*|MSYS*|CYGWIN*)
      if command -v winget.exe >/dev/null 2>&1; then
        log_info "Installing jq with winget"
        winget.exe install --id jqlang.jq -e \
          --accept-package-agreements --accept-source-agreements
        refresh_windows_winget_path
      elif command -v choco.exe >/dev/null 2>&1; then
        log_info "Installing jq with Chocolatey"
        choco.exe install jq -y
      else
        die "jq is required and neither winget nor Chocolatey is available"
      fi
      ;;
    Linux)
      log_info "Installing jq"
      if command -v apt-get >/dev/null 2>&1; then
        run_linux_privileged apt-get update
        run_linux_privileged apt-get install -y jq
      elif command -v dnf >/dev/null 2>&1; then
        run_linux_privileged dnf install -y jq
      elif command -v yum >/dev/null 2>&1; then
        run_linux_privileged yum install -y jq
      elif command -v pacman >/dev/null 2>&1; then
        run_linux_privileged pacman -Sy --noconfirm jq
      else
        die "Unable to install jq automatically on this Linux distribution"
      fi
      ;;
    *)
      die "Unsupported development operating system: $(uname -s)"
      ;;
  esac

  require_command jq
}
