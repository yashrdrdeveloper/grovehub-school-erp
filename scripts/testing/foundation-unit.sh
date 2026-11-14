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
# Go's cover instrumenter reads logical source paths directly and does not
# honor virtual -overlay replacements for generated production files.
# Materialize virtual files temporarily, then remove them even on failure.
# No generated file is committed or left behind in maintained packages.
manifest="$(mktemp)"
cleanup_generated() {
  if [[ -f "${manifest}" ]]; then
    while IFS= read -r logical; do
      [[ -n "${logical}" ]] && rm -f -- "${logical}"
    done < "${manifest}"
    rm -f -- "${manifest}"
  fi
}
trap cleanup_generated EXIT

python3 - "${BACKEND_DIR}/.generated/overlay.json" "${manifest}" <<'PY'
import json, pathlib, shutil, sys
overlay = pathlib.Path(sys.argv[1])
manifest = pathlib.Path(sys.argv[2])
replace = json.loads(overlay.read_text())["Replace"]
with manifest.open("w") as paths:
    for logical, generated in sorted(replace.items()):
        target = pathlib.Path(logical)
        if target.exists():
            raise SystemExit(f"generated file would overwrite maintained source: {target}")
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(generated, logical)
        paths.write(logical + "\n")
        paths.flush()
PY

coverage_file="${BACKEND_DIR}/.generated/foundation-unit-coverage.out"
log_info "Running isolated Foundation Ginkgo/Gomega unit tests with no real database or broker"
(
  cd "${BACKEND_DIR}"
  go test -vet=off -overlay="${BACKEND_DIR}/.generated/overlay.json" \
    -covermode=atomic -coverpkg=./services/foundation/internal/app/... \
    -coverprofile="${coverage_file}" \
    ./services/foundation/internal/app/... \
    ./services/foundation/internal/tests/unit-test/...
)

# Run every production package as an explicit target: this counts 0%-covered
# packages too, not just the transitive dependencies of existing specs.
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
