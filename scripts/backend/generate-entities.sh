#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

require_command go
ensure_tool_dirs

# Pin tools; do not depend on globally installed executables.
GOMBOK_VERSION="v1.2.0"
GOIMPORTS_VERSION="v0.51.0"
install_generator() {
  local name="$1" package="$2" version="$3" version_file="${STATE_DIR}/${name}.version"
  if [[ -x "${TOOLS_BIN_DIR}/${name}" && -f "${version_file}" && "$(cat "${version_file}")" == "${version}" ]]; then
    return
  fi
  log_info "Installing ${name}@${version} for GroveHub entity generation"
  GOBIN="${TOOLS_BIN_DIR}" go install "${package}@${version}"
  printf '%s\n' "${version}" >"${version_file}"
}

# Intentionally limited to entity packages that have opted into @Builder/@GroveTable.
# Never walk the full backend: Gombok scans recursively and will overwrite generated files.
ENTITY_DIRS=("${BACKEND_DIR}/services/foundation/entity")
for entity_dir in "${ENTITY_DIRS[@]}"; do
  [[ -d "${entity_dir}" ]] || die "Missing entity package: ${entity_dir}"

  # Do not rewrite generated Go files on each Air build: it triggers another rebuild.
  fingerprint=""
  for file in "${entity_dir}"/*.go; do
    [[ -f "${file}" ]] || continue
    case "${file}" in *_gombok.go|*/zz_grove_tables.gen.go) continue ;; esac
    fingerprint+="$(shasum -a 256 "${file}")"$'\n'
  done
  fingerprint+="Gombok=${GOMBOK_VERSION};GormTable=v1"
  digest="$(printf '%s' "${fingerprint}" | shasum -a 256 | awk '{print $1}')"
  state_file="${STATE_DIR}/$(basename "$(dirname "${entity_dir}")")-entity-codegen.sha256"
  if [[ -f "${state_file}" && "$(cat "${state_file}")" == "${digest}" && -f "${entity_dir}/contact_gombok.go" && -f "${entity_dir}/person_gombok.go" && -f "${entity_dir}/zz_grove_tables.gen.go" ]]; then
    continue
  fi

  install_generator gombok github.com/YangTaeyoung/gombok "${GOMBOK_VERSION}"
  install_generator goimports golang.org/x/tools/cmd/goimports "${GOIMPORTS_VERSION}"
  export PATH="${TOOLS_BIN_DIR}:${PATH}"
  log_info "Generating Gombok methods under ${entity_dir}"
  (
    cd -- "${entity_dir}"
    "${TOOLS_BIN_DIR}/gombok"
  )
  (cd -- "${BACKEND_DIR}" && go run ./libs/core/entity/cmd/gormtable "${entity_dir}")
  for expected in contact_gombok.go person_gombok.go zz_grove_tables.gen.go; do
    [[ -s "${entity_dir}/${expected}" ]] || die "Entity code generation did not produce ${expected}"
  done
  printf '%s\n' "${digest}" >"${state_file}"
done
