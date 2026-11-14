#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

require_command go

# Fail closed on unimplemented annotations, especially security policies.
log_info "Validating annotation targets and rejecting metadata-only directives"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/check "${BACKEND_DIR}"
)

# Generate typed boilerplate directly into service-local ignored .generated/
# directories. The original Go packages stay entirely developer-maintained.
log_info "Generating annotated entity builders/accessors"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/entities "${BACKEND_DIR}"
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
        go run ./libs/core/annotation/generation/gormtable "${entity_dir}" "${output_dir}/zz_tables.gen.go"
      )
    fi
  done
done

log_info "Generating annotated HTTP route registrations"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/routes "${BACKEND_DIR}"
)

log_info "Generating annotated in-process client proxies"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/clients "${BACKEND_DIR}"
)

log_info "Generating annotated repository searches"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/search "${BACKEND_DIR}"
)

log_info "Generating @PatchRequestDto JSON methods"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/patchdto "${BACKEND_DIR}"
)

log_info "Generating strict DTO mappers"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/mappers "${BACKEND_DIR}"
)

log_info "Generating domain-owned event schema identifiers"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/events "${BACKEND_DIR}"
)

# Only generate test suites and test-owned mocks in test directories.
# These files are included by go test -overlay, not compiled into the production binary.
log_info "Generating annotated Ginkgo suites and isolated GoMock implementations"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/testing "${BACKEND_DIR}"
)

log_info "Generating consolidated Go source overlay"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/overlay "${BACKEND_DIR}"
)

log_info "Annotation source overlay ready: ${BACKEND_DIR}/.generated/overlay.json"
