#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

require_command go
ensure_tool_dirs

GOMBOK_VERSION="v1.2.0"
GOIMPORTS_VERSION="v0.51.0"

install_generator() {
  local name="$1" package="$2" version="$3" version_file="${STATE_DIR}/${name}.version"
  if [[ -x "${TOOLS_BIN_DIR}/${name}" && -f "${version_file}" && "$(cat "${version_file}")" == "${version}" ]]; then
    return 0
  fi
  log_info "Installing ${name}@${version} for GroveHub entity generation"
  GOBIN="${TOOLS_BIN_DIR}" go install "${package}@${version}"
  printf '%s\n' "${version}" >"${version_file}"
}

# Scan service entity packages individually: Gombok walks recursively from its CWD.
# Non-annotated packages remain untouched, and tool installations are on demand.
for entity_dir in "${BACKEND_DIR}"/services/*/entity; do
  [[ -d "${entity_dir}" ]] || continue
  declared=()
  has_table=0
  fingerprint=""
  for source in "${entity_dir}"/*.go; do
    [[ -f "${source}" ]] || continue
    case "${source}" in *_gombok.go|*/zz_grove_tables.gen.go) continue ;; esac
    fingerprint+="$(shasum -a 256 "${source}")"$'\n'
    if grep -Eq '^[[:space:]]*//[[:space:]]*@(Builder|Getter|Setter|AllArgsConstructor|NoArgsConstructor|RequiredArgsConstructor|ToString|Equals)([.[:space:]]|$)' "${source}"; then
      declared+=("${entity_dir}/$(basename "${source}" .go)_gombok.go")
    fi
    if grep -Eq '^[[:space:]]*//[[:space:]]*@GroveTable[[:space:]]+' "${source}"; then
      has_table=1
    fi
  done
  if [[ ${#declared[@]} -eq 0 && "${has_table}" == 0 ]]; then
    continue
  fi

  fingerprint+="Gombok=${GOMBOK_VERSION};goimports=${GOIMPORTS_VERSION};GroveTable=v1"
  digest="$(printf '%s' "${fingerprint}" | shasum -a 256 | awk '{print $1}')"
  state_file="${STATE_DIR}/$(basename "$(dirname "${entity_dir}")")-entity-codegen.sha256"
  ready=1
  if [[ ! -f "${state_file}" || "$(cat "${state_file}")" != "${digest}" ]]; then ready=0; fi
  for generated in "${declared[@]}"; do [[ -s "${generated}" ]] || ready=0; done
  if [[ "${has_table}" == 1 && ! -s "${entity_dir}/zz_grove_tables.gen.go" ]]; then ready=0; fi
  if [[ "${ready}" == 1 ]]; then continue; fi

  if [[ ${#declared[@]} -gt 0 ]]; then
    install_generator gombok github.com/YangTaeyoung/gombok "${GOMBOK_VERSION}"
    install_generator goimports golang.org/x/tools/cmd/goimports "${GOIMPORTS_VERSION}"
    export PATH="${TOOLS_BIN_DIR}:${PATH}"
    log_info "Running Gombok in ${entity_dir}"
    (cd -- "${entity_dir}" && "${TOOLS_BIN_DIR}/gombok")
    for generated in "${declared[@]}"; do
      [[ -s "${generated}" ]] || die "Gombok did not generate ${generated}"
    done
  fi
  if [[ "${has_table}" == 1 ]]; then
    log_info "Generating GORM table methods in ${entity_dir}"
    (cd -- "${BACKEND_DIR}" && go run ./libs/core/entity/cmd/gormtable "${entity_dir}")
    [[ -s "${entity_dir}/zz_grove_tables.gen.go" ]] || die "GORM TableName generation failed for ${entity_dir}"
  fi
  printf '%s\n' "${digest}" >"${state_file}"
done
