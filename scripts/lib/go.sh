#!/usr/bin/env bash

sync_backend_go_modules() {
  require_command go

  log_info "Syncing GroveHub backend Go modules"
  (
    cd -- "${BACKEND_DIR}"
    go mod tidy
  )
}
