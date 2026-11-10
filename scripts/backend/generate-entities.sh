#!/usr/bin/env bash
set -Eeuo pipefail
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
# Backward-compatible entity generation command; the unified pipeline owns all generators.
bash "${SCRIPT_DIR}/generate-annotations.sh" "$@"
