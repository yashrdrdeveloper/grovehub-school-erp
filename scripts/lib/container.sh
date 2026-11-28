#!/usr/bin/env bash
set -Eeuo pipefail

CONTAINER_LIB_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

# shellcheck source=common.sh
source "${CONTAINER_LIB_DIR}/common.sh"

DOCKER_BINARY=""

refresh_homebrew_path() {
  if [[ -x /opt/homebrew/bin/brew ]]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"
  elif [[ -x /usr/local/bin/brew ]]; then
    eval "$(/usr/local/bin/brew shellenv)"
  fi
}

prepare_existing_docker_cli() {
  case "$(uname -s)" in
    Darwin)
      refresh_homebrew_path
      ;;
    MINGW*|MSYS*|CYGWIN*)
      local docker_desktop_bin="/c/Program Files/Docker/Docker/resources/bin"
      if [[ -d "${docker_desktop_bin}" ]]; then
        export PATH="${PATH}:${docker_desktop_bin}"
      fi
      ;;
  esac
}

ensure_homebrew() {
  refresh_homebrew_path
  if command -v brew >/dev/null 2>&1; then
    return 0
  fi

  require_command curl
  log_info "Homebrew is not installed; installing it for the GroveHub development toolchain"
  NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  refresh_homebrew_path
  require_command brew
}

ensure_macos_container_runtime() {
  ensure_homebrew

  local formula
  for formula in colima docker docker-compose; do
    if ! brew list --formula "${formula}" >/dev/null 2>&1; then
      log_info "Installing ${formula}"
      brew install "${formula}"
    fi
  done

  if ! colima status >/dev/null 2>&1; then
    log_info "Starting Colima"
    colima start \
      --cpu "${COLIMA_CPU:-2}" \
      --memory "${COLIMA_MEMORY_GB:-4}" \
      --disk "${COLIMA_DISK_GB:-20}"
  fi
}

ensure_windows_container_runtime() {
  prepare_existing_docker_cli

  if ! command -v docker >/dev/null 2>&1; then
    if command -v winget.exe >/dev/null 2>&1; then
      log_info "Docker Desktop is not installed; installing it with winget"
      winget.exe install --id Docker.DockerDesktop -e \
        --accept-package-agreements --accept-source-agreements
    elif command -v choco.exe >/dev/null 2>&1; then
      log_info "Docker Desktop is not installed; installing it with Chocolatey"
      choco.exe install docker-desktop -y
    else
      die "Docker Desktop is required on Windows and neither winget nor Chocolatey is available"
    fi
    prepare_existing_docker_cli
  fi

  require_command docker
  if docker info >/dev/null 2>&1; then
    return 0
  fi

  require_command powershell.exe
  log_info "Starting Docker Desktop"
  powershell.exe -NoProfile -Command \
    "Start-Process 'C:\Program Files\Docker\Docker\Docker Desktop.exe'" >/dev/null 2>&1 || true
}

run_linux_privileged() {
  if [[ "$(id -u)" == "0" ]]; then
    "$@"
    return
  fi
  require_command sudo
  sudo "$@"
}

ensure_linux_curl() {
  if command -v curl >/dev/null 2>&1; then
    return 0
  fi

  log_info "curl is not installed; installing it"
  if command -v apt-get >/dev/null 2>&1; then
    run_linux_privileged apt-get update
    run_linux_privileged apt-get install -y curl ca-certificates
  elif command -v dnf >/dev/null 2>&1; then
    run_linux_privileged dnf install -y curl ca-certificates
  elif command -v yum >/dev/null 2>&1; then
    run_linux_privileged yum install -y curl ca-certificates
  elif command -v pacman >/dev/null 2>&1; then
    run_linux_privileged pacman -Sy --noconfirm curl ca-certificates
  else
    die "Unable to install curl automatically on this Linux distribution"
  fi
}

ensure_linux_container_runtime() {
  if ! command -v docker >/dev/null 2>&1; then
    ensure_linux_curl
    log_info "Docker Engine is not installed; installing it"
    local installer
    installer="$(mktemp)"
    curl -fsSL https://get.docker.com -o "${installer}"
    run_linux_privileged sh "${installer}"
    rm -f -- "${installer}"
  fi

  if docker info >/dev/null 2>&1; then
    return 0
  fi

  log_info "Starting Docker Engine"
  if command -v systemctl >/dev/null 2>&1; then
    run_linux_privileged systemctl start docker
  elif command -v service >/dev/null 2>&1; then
    run_linux_privileged service docker start
  fi
}

configure_docker_access() {
  DOCKER_BINARY="$(type -P docker || true)"
  [[ -n "${DOCKER_BINARY}" ]] || return 1

  if "${DOCKER_BINARY}" info >/dev/null 2>&1; then
    return 0
  fi

  if [[ "$(uname -s)" == "Linux" ]] && command -v sudo >/dev/null 2>&1 \
    && sudo "${DOCKER_BINARY}" info >/dev/null 2>&1; then
    docker() {
      sudo "${DOCKER_BINARY}" "$@"
    }
    return 0
  fi

  return 1
}

wait_for_docker() {
  local attempt
  for ((attempt = 1; attempt <= 120; attempt++)); do
    if configure_docker_access && docker info >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
  done
  return 1
}

require_docker() {
  prepare_existing_docker_cli

  case "$(uname -s)" in
    Darwin)
      ensure_macos_container_runtime
      ;;
    MINGW*|MSYS*|CYGWIN*)
      ensure_windows_container_runtime
      ;;
    Linux)
      ensure_linux_container_runtime
      ;;
    *)
      die "Unsupported development operating system: $(uname -s)"
      ;;
  esac

  if ! wait_for_docker; then
    die "Container runtime could not be started"
  fi
}

docker_compose() {
  if docker compose version >/dev/null 2>&1; then
    docker compose "$@"
    return
  fi
  if command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
    return
  fi
  die "Docker Compose is not available"
}

container_exists() {
  local container_name="$1"
  docker inspect "${container_name}" >/dev/null 2>&1
}

container_is_running() {
  local container_name="$1"
  [[ "$(docker inspect --format '{{.State.Running}}' "${container_name}" 2>/dev/null || true)" == "true" ]]
}
