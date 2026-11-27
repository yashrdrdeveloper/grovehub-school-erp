#!/usr/bin/env bash
set -euo pipefail

# Generate a consumer-local Go client from an explicitly supplied, versioned spec.
# Requires Docker; generation is offline after the pinned image is available.
if [[ $# -ne 3 ]]; then
  echo "usage: $0 SPEC OUTPUT_DIR PACKAGE_NAME" >&2
  exit 2
fi
spec="$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
output_parent="$(cd "$(dirname "$2")" && pwd)"
output="$output_parent/$(basename "$2")"
package="$3"
[[ -f "$spec" && "$package" =~ ^[a-z][a-z0-9_]*$ ]] || exit 2
if [[ -e "$output" ]]; then
  echo "Output already exists; generate into a fresh directory and review the diff." >&2
  exit 2
fi
stage="$(mktemp -d "$output_parent/.clientgen.XXXXXX")"
trap 'rm -rf -- "$stage"' EXIT
image="openapitools/openapi-generator-cli:v7.15.0"
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$spec:/input/openapi.yaml:ro" \
  -v "$stage:/output" \
  "$image" generate -i /input/openapi.yaml -g go -o /output \
  --additional-properties "packageName=$package,isGoSubmodule=true,hideGenerationTimestamp=true,withGoMod=false" \
  --global-property "apiTests=false,modelTests=false,apiDocs=false,modelDocs=false"
find "$stage" -maxdepth 1 -name '*.go' -print -quit | grep -q .
mv -- "$stage" "$output"
trap - EXIT
