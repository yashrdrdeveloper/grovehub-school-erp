#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

require_command go

# Explicit roots let the same processors serve another application.
BACKEND_DIR="${1:-${BACKEND_DIR}}"
CORE_DIR="${2:-${CORE_DIR}}"
BACKEND_DIR="$(cd -- "$BACKEND_DIR" && pwd)"
CORE_DIR="$(cd -- "$CORE_DIR" && pwd)"

# A clean checkout needs module checksums for generated mocks and unit-only
# dependencies (GoMock, sqlmock, etc.) before running annotation generators
# or instrumentation. Go does not discover all generated source imports during
# normal go mod tidy because those files live in a virtual overlay.
log_info "Preparing pinned Go modules for generated source and unit tests"
(
  cd -- "${BACKEND_DIR}"
  # Go's lazy module loading may omit checksum entries for code-generation
  # and isolated-test dependencies even when plain 'go mod download' succeeds.
  go mod download
)
(
  cd -- "${CORE_DIR}"
  go mod download
  go mod download go.uber.org/mock github.com/DATA-DOG/go-sqlmock
)

# Generation is a fresh snapshot, not an incremental accumulation. Removing
# ignored service-local output prevents orphaned virtual .go files from
# referring to packages deleted by a refactor (e.g. client/person/dto/).
# The overlay generator still fails closed on genuine output/source conflicts.
log_info "Cleaning stale generated annotation outputs"
for owner in "${BACKEND_DIR}"/services/* "${BACKEND_DIR}"/launcher "${CORE_DIR}"; do
  [[ -d "${owner}" ]] || continue
  generated="${owner}/.generated"
  [[ ! -L "${generated}" ]] || die "Refusing to clean symlinked generated directory: ${generated}"
  if [[ -d "${generated}" ]]; then
    rm -rf -- "${generated}"
  fi
done
rm -f -- "${BACKEND_DIR}/.generated/overlay.json"

# Fail closed on unimplemented annotations, especially security policies.
log_info "Validating annotation targets and rejecting metadata-only directives"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" check "${BACKEND_DIR}"
)

# Generate typed boilerplate directly into service-local ignored .generated/
# directories. The original Go packages stay entirely developer-maintained.
log_info "Generating annotated entity builders/accessors"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" entities "${BACKEND_DIR}"
)

# GORM TableName() uses the same original @Entity declarations.
for service_dir in "${BACKEND_DIR}"/services/*; do
  [[ -d "${service_dir}" ]] || continue
  for entity_dir in "${service_dir}/entity" "${service_dir}/internal/app/entity"; do
    [[ -d "${entity_dir}" ]] || continue
    # The Go overlay mirrors the owning package relative to the service root.
    if grep -Eq '^[[:space:]]*//[[:space:]]*@Entity[[:space:]]+' "${entity_dir}"/*.go; then
      relative_dir="${entity_dir#${service_dir}/}"
      output_dir="${service_dir}/.generated/${relative_dir}"
      (
        cd -- "${BACKEND_DIR}"
        go run "${CORE_MODULE}/annotation/cmd" gormtable "${entity_dir}" "${output_dir}/zz_tables.gen.go"
      )
    fi
  done
done

log_info "Generating annotated HTTP route registrations"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" routes "${BACKEND_DIR}"
)

log_info "Generating annotated in-process client proxies"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" clients "${BACKEND_DIR}"
)

log_info "Generating annotated repository searches"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" search "${BACKEND_DIR}"
)

log_info "Generating @PatchRequestDto JSON methods"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" patchdto "${BACKEND_DIR}"
)

log_info "Generating strict DTO mappers"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" mappers "${BACKEND_DIR}"
)

log_info "Generating domain-owned event schema identifiers"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" events "${BACKEND_DIR}"
)

# Only generate test suites and test-owned mocks in test directories.
# These files are included by go test -overlay, not compiled into the production binary.
log_info "Generating annotated Ginkgo suites and isolated GoMock implementations"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" testing "${BACKEND_DIR}" "${CORE_DIR}"
)

log_info "Generating consolidated Go source overlay"
(
  cd -- "${BACKEND_DIR}"
  go run "${CORE_MODULE}/annotation/cmd" overlay "${BACKEND_DIR}" "${CORE_DIR}"
)

log_info "Annotation source overlay ready: ${BACKEND_DIR}/.generated/overlay.json"
