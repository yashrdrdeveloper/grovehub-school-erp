#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

FRONTEND_DIR="${REPO_ROOT}/frontend"
export COREPACK_HOME="${TOOLS_DIR}/corepack"

require_command node
node_major="$(node --version | sed -E 's/^v([0-9]+).*/\1/')"
[[ "${node_major}" =~ ^[0-9]+$ && "${node_major}" -ge 20 ]] || die "Frontend requires Node.js 20 or newer"

run_pnpm() {
  if command -v corepack >/dev/null 2>&1; then
    corepack pnpm "$@"
  elif command -v pnpm >/dev/null 2>&1; then
    pnpm "$@"
  else
    die "pnpm not available. Install Corepack or pnpm using the version in frontend/package.json."
  fi
}

ensure_tool_dirs
cd -- "${FRONTEND_DIR}"
log_info "Installing frontend dependencies with pinned pnpm"
run_pnpm install --no-frozen-lockfile
log_info "Frontend dependencies installed"
