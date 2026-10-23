#!/usr/bin/env bash
set -Eeuo pipefail

pid_is_running() {
  local pid="$1"
  kill -0 "${pid}" >/dev/null 2>&1
}

process_command() {
  local pid="$1"
  ps -p "${pid}" -o command= 2>/dev/null || true
}

command_matches_any() {
  local command_line="$1"
  shift

  local fragment
  for fragment in "$@"; do
    if [[ "${command_line}" == *"${fragment}"* ]]; then
      return 0
    fi
  done

  return 1
}

owned_pid_is_running() {
  local pid_file="$1"
  local command_fragment="$2"

  [[ -f "${pid_file}" ]] || return 1

  local pid
  pid="$(cat "${pid_file}")"
  [[ "${pid}" =~ ^[0-9]+$ ]] || return 1
  pid_is_running "${pid}" || return 1

  local command_line
  command_line="$(process_command "${pid}")"
  [[ "${command_line}" == *"${command_fragment}"* ]]
}

owned_pid_is_running_any() {
  local pid_file="$1"
  shift

  [[ -f "${pid_file}" ]] || return 1

  local pid
  pid="$(cat "${pid_file}")"
  [[ "${pid}" =~ ^[0-9]+$ ]] || return 1
  pid_is_running "${pid}" || return 1

  command_matches_any "$(process_command "${pid}")" "$@"
}

stop_owned_process() {
  local pid_file="$1"
  local command_fragment="$2"
  local label="$3"

  stop_owned_process_any "${pid_file}" "${label}" "${command_fragment}"
}

stop_owned_process_any() {
  local pid_file="$1"
  local label="$2"
  shift 2

  if [[ ! -f "${pid_file}" ]]; then
    log_info "${label} is already stopped"
    return 0
  fi

  local pid
  pid="$(cat "${pid_file}")"

  if [[ ! "${pid}" =~ ^[0-9]+$ ]] || ! pid_is_running "${pid}"; then
    rm -f -- "${pid_file}"
    log_info "${label} is already stopped"
    return 0
  fi

  local command_line
  command_line="$(process_command "${pid}")"
  if ! command_matches_any "${command_line}" "$@"; then
    die "Refusing to stop PID ${pid}; it is not the managed ${label} process"
  fi

  log_info "Stopping ${label} (PID ${pid})"
  kill -TERM "${pid}"

  local attempt
  for ((attempt = 1; attempt <= 10; attempt++)); do
    if ! pid_is_running "${pid}"; then
      rm -f -- "${pid_file}"
      log_info "${label} stopped"
      return 0
    fi
    sleep 1
  done

  log_warn "${label} did not stop after 10 seconds; forcing termination"
  kill -KILL "${pid}"
  rm -f -- "${pid_file}"
}
