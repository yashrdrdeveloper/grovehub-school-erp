#!/usr/bin/env bash

sync_backend_go_modules() {
  require_command go

  log_info "Syncing GroveHub backend Go modules"
  (
    cd -- "${CORE_DIR}"
    GOWORK=off go mod tidy
  )
  (
    cd -- "${BACKEND_DIR}"
    go mod tidy
  )
}
