#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

mode="${1:---strict}"
case "${mode}" in
  --report-only|--strict) ;;
  *) die "usage: foundation-unit.sh [--strict|--report-only]" ;;
esac

require_command go
# Build the exact same generated sources as production (plus tests and mocks).
bash "${SCRIPT_DIR}/../backend/generate-annotations.sh"

mkdir -p "${BACKEND_DIR}/.generated"
coverage_file="${BACKEND_DIR}/.generated/foundation-unit-coverage.out"
log_info "Running isolated Foundation Ginkgo/Gomega unit tests with no real database or broker"
(
  cd "${BACKEND_DIR}"
  go test -vet=off -overlay="${BACKEND_DIR}/.generated/overlay.json" \
    -covermode=atomic -coverpkg=./services/foundation/internal/app/... \
    -coverprofile="${coverage_file}" \
    ./services/foundation/internal/tests/unit-test/...
)

# Calculate actual instrumented statement coverage directly from the profile.
# Do not call go tool cover -func: its virtual overlay filenames are not present
# on disk, even though the generated Go code was compiled and executed.
read -r total covered < <(awk 'NR > 1 { total += $2; if ($3 > 0) covered += $2 } END { printf "%d %d\n",total,covered }' "${coverage_file}")
((total > 0)) || die "No instrumented Foundation statements found"
percentage="$(awk -v covered="${covered}" -v total="${total}" 'BEGIN { printf "%.2f",100*covered/total }')"
log_info "Foundation statement coverage: ${percentage}% (${covered}/${total} statements)"
if [[ "${mode}" == "--strict" && "${covered}" != "${total}" ]]; then
  die "Foundation unit coverage is not 100%: ${percentage}%. Every uncovered statement must be covered before passing the strict gate."
fi
