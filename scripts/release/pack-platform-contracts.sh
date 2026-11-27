#!/usr/bin/env bash
set -euo pipefail
if [[ $# -ne 2 ]]; then
  echo "Usage: pack-platform-contracts.sh DESTINATION COMMIT" >&2
  exit 2
fi
destination="$1"
revision="$2"
[[ "$revision" =~ ^[0-9a-f]{40}$ ]] || exit 2
root="$(git rev-parse --show-toplevel)"
actual="$(git -C "$root" rev-parse HEAD)"
[[ "$actual" == "$revision" ]] || { echo "Contract revision does not match checkout" >&2; exit 2; }
mkdir -p "$destination"
stage="$(mktemp -d)"
trap 'rm -rf -- "$stage"' EXIT
for service in audit foundation identity integration; do
  relative="$service/internal/resources/api/openapi/v1/openapi.yaml"
  mkdir -p "$stage/services/$(dirname "$relative")"
  git -C "$root" show "$revision:backend/projects/platform/services/$relative" > "$stage/services/$relative"
done
mkdir -p "$stage/school"
git -C "$root" show "$revision:backend/projects/school/services/registry/internal/resources/api/openapi/v1/openapi.yaml" > "$stage/school/registry.yaml"
mkdir -p "$stage/tools"
git -C "$root" show "$revision:backend/libs/core/tools/clientgen/generate.sh" > "$stage/tools/generate.sh"
printf '{"commit":"%s","generator":"openapitools/openapi-generator-cli:v7.15.0"}\n' "$revision" > "$stage/manifest.json"
(
  cd "$stage"
  find services tools school -type f -print0 | LC_ALL=C sort -z | xargs -0 sha256sum > SHA256SUMS
  sha256sum manifest.json >> SHA256SUMS
  sha256sum -c SHA256SUMS
)
tar -C "$stage" -czf "$destination/platform-contracts.tar.gz" manifest.json SHA256SUMS services tools school
(cd "$destination" && sha256sum platform-contracts.tar.gz > platform-contracts.tar.gz.sha256)
