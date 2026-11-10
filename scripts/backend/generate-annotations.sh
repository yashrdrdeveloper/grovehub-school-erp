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
    return
  fi
  log_info "Installing ${name}@${version}"
  GOBIN="${TOOLS_BIN_DIR}" go install "${package}@${version}"
  printf '%s\n' "${version}" >"${version_file}"
}

# Foundation entity Gombok is preserved; all services may additionally generate routes and client adapters.
# The generated methods never reside in services/foundation/entity.
entity_dir="${BACKEND_DIR}/services/foundation/entity"
generated_dir="${BACKEND_DIR}/services/foundation/.generated/entity"
state_file="${STATE_DIR}/foundation-entity-codegen.sha256"
declare -a source_files=()
declare -a generated_files=()
fingerprint=""

for file in "${entity_dir}"/*.go; do
  [[ -f "${file}" ]] || continue
  case "${file}" in *_gombok.go|*/zz_grove_tables.gen.go|*/zz_tables.gen.go) continue ;; esac
  source_files+=("${file}")
  fingerprint+="$(shasum -a 256 "${file}")"$'\n'
  if grep -Eq '^[[:space:]]*//[[:space:]]*@(Entity|Builder|Getter|Setter|AllArgsConstructor|NoArgsConstructor|RequiredArgsConstructor|ToString|Equals)([.[:space:]]|$)' "${file}"; then
    generated_files+=("$(basename "${file}" .go)_gombok.go")
  fi
done

fingerprint+="gombok=${GOMBOK_VERSION};goimports=${GOIMPORTS_VERSION};entity=v1;table=v3;overlay=v1"
digest="$(printf '%s' "${fingerprint}" | shasum -a 256 | awk '{print $1}')"
ready=1
[[ -f "${state_file}" && "$(cat "${state_file}")" == "${digest}" ]] || ready=0
for generated in "${generated_files[@]}" "zz_tables.gen.go"; do
  [[ -s "${generated_dir}/${generated}" ]] || ready=0
done

if [[ "${ready}" != 1 ]]; then
  install_generator gombok github.com/YangTaeyoung/gombok "${GOMBOK_VERSION}"
  install_generator goimports golang.org/x/tools/cmd/goimports "${GOIMPORTS_VERSION}"
  export PATH="${TOOLS_BIN_DIR}:${PATH}"

  stage="$(mktemp -d "${STATE_DIR}/entity-stage.XXXXXXXX")"
  trap 'rm -rf -- "${stage}"' EXIT
  for file in "${source_files[@]}"; do
    cp -- "${file}" "${stage}/"
  done

  log_info "Expanding @Entity annotations in isolated staging directory"
  (
    cd -- "${BACKEND_DIR}"
    go run ./libs/core/annotation/generation/expand "${stage}"
  )
  log_info "Generating Foundation Gombok implementations in isolated staging directory"
  (
    cd -- "${stage}"
    "${TOOLS_BIN_DIR}/gombok"
  )
  for generated in "${generated_files[@]}"; do
    [[ -s "${stage}/${generated}" ]] || die "Gombok did not generate ${generated}"
  done

  mkdir -p -- "${generated_dir}"
  # Clean old generated files so removed annotations cannot leave stale methods.
  find "${generated_dir}" -maxdepth 1 -type f \( -name '*_gombok.go' -o -name 'zz_tables.gen.go' \) -delete
  for generated in "${generated_files[@]}"; do
    cp -- "${stage}/${generated}" "${generated_dir}/"
  done
  (
    cd -- "${BACKEND_DIR}"
    go run ./libs/core/annotation/generation/gormtable "${entity_dir}" "${generated_dir}/zz_tables.gen.go"
  )
  [[ -s "${generated_dir}/zz_tables.gen.go" ]] || die "GORM table generator did not produce an output"

  printf '%s\n' "${digest}" >"${state_file}"
fi

# Generate typed HTTP registrations from @Routes and @GET/@POST/etc
# in controller contracts. The router writes only to service-local .generated/.
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/routes "${BACKEND_DIR}"
)

# Generate in-process @Client/@Proxy forwarding methods for supported clients.
# Source client/impl constructors remain explicit and unchanged.
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/clients "${BACKEND_DIR}"
)

# The overlay manifest is recreated even when the source fingerprint is unchanged.
# This ensures normal Go builds see virtual files under the original entity package.
(
  cd -- "${BACKEND_DIR}"
  go run ./libs/core/annotation/generation/overlay "${BACKEND_DIR}"
)
log_info "Annotation overlay ready (Foundation entities, HTTP routes, client proxies): ${BACKEND_DIR}/.generated/overlay.json"
