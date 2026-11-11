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
for entity_dir in "${BACKEND_DIR}"/services/*/entity; do
  [[ -d "${entity_dir}" ]] || continue
  # Avoid creating empty generated folders for services without @Entity.
  if grep -Eq '^[[:space:]]*//[[:space:]]*@Entity[[:space:]]+' "${entity_dir}"/*.go; then
    output_dir="$(dirname "${entity_dir}")/.generated/entity"
    (
      cd -- "${BACKEND_DIR}"
      go run ./libs/core/annotation/generation/gormtable "${entity_dir}" "${output_dir}/zz_tables.gen.go"
    )
  fi
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

log_info "Generating consolidated Go source overlay"
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/overlay "${BACKEND_DIR}"
)

log_info "Annotation source overlay ready: ${BACKEND_DIR}/.generated/overlay.json"
