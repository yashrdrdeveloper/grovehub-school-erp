#!/usr/bin/env bash
# Shared lifecycle for the optional frontend-only Cloudflare Quick Tunnel.
# Sourced by the frontend start/status/stop scripts after common and paths.

FRONTEND_TUNNEL_PID_FILE="$RUNTIME_PID_DIR/frontend-tunnel.pid"
FRONTEND_TUNNEL_URL_FILE="$RUNTIME_PID_DIR/frontend-tunnel.url"
FRONTEND_TUNNEL_LOG_FILE="$RUNTIME_LOG_DIR/frontend-tunnel.log"

frontend_tunnel_running() {
  owned_pid_is_running "$FRONTEND_TUNNEL_PID_FILE" "cloudflared tunnel" &&
    command_matches_any "$(process_command "$(cat "$FRONTEND_TUNNEL_PID_FILE")")" "--url http://127.0.0.1:"
}

read_frontend_tunnel_url() {
  local url
  frontend_tunnel_running || return 1
  [[ -f "$FRONTEND_TUNNEL_URL_FILE" ]] || return 1
  url="$(cat "$FRONTEND_TUNNEL_URL_FILE")"
  [[ "$url" =~ ^https://[a-z0-9-]+\.trycloudflare\.com$ ]] || return 1
  printf '%s\n' "$url"
}

stop_frontend_tunnel() {
  if frontend_tunnel_running; then
    stop_owned_process "$FRONTEND_TUNNEL_PID_FILE" "cloudflared tunnel" "Frontend Cloudflare tunnel"
  fi
  rm -f -- "$FRONTEND_TUNNEL_PID_FILE" "$FRONTEND_TUNNEL_URL_FILE"
}

prepare_frontend_tunnel() {
  local origin="$1"
  local url attempt command_line
  if [[ "${FRONTEND_TUNNEL:-1}" == "0" ]]; then
    stop_frontend_tunnel
    return 0
  fi

  if url="$(read_frontend_tunnel_url)"; then
    command_line="$(process_command "$(cat "$FRONTEND_TUNNEL_PID_FILE")")"
    if [[ "$command_line" == *"--url $origin "* || "$command_line" == *"--url $origin" ]]; then
      return 0
    fi
  fi

  stop_frontend_tunnel
  if ! command -v cloudflared >/dev/null 2>&1; then
    log_warn "Cloudflare tunnel unavailable: install cloudflared (macOS: brew install cloudflared). Continuing locally."
    return 0
  fi

  log_warn "Starting a public development tunnel. Use test data only; anyone with the URL can reach the frontend."
  : >"$FRONTEND_TUNNEL_LOG_FILE"
  nohup cloudflared tunnel --no-autoupdate --url "$origin" </dev/null >>"$FRONTEND_TUNNEL_LOG_FILE" 2>&1 &
  printf '%s\n' "$!" >"$FRONTEND_TUNNEL_PID_FILE"

  for ((attempt = 1; attempt <= 30; attempt++)); do
    if ! pid_is_running "$(cat "$FRONTEND_TUNNEL_PID_FILE")"; then
      break
    fi
    url="$(grep -Eo 'https://[a-z0-9-]+\.trycloudflare\.com' "$FRONTEND_TUNNEL_LOG_FILE" | head -n 1 || true)"
    if [[ -n "$url" ]] && frontend_tunnel_running; then
      printf '%s\n' "$url" >"$FRONTEND_TUNNEL_URL_FILE"
      return 0
    fi
    sleep 1
  done
  log_warn "Cloudflare tunnel did not start; see $FRONTEND_TUNNEL_LOG_FILE. Continuing locally."
  stop_frontend_tunnel
}

report_frontend_tunnel() {
  local url
  if url="$(read_frontend_tunnel_url)"; then
    log_info "Cloudflare HTTPS URL: $url"
    log_info "Cloudflare tunnel log: $FRONTEND_TUNNEL_LOG_FILE"
  fi
}
