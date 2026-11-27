#!/usr/bin/env bash
set -euo pipefail

project="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
generator="$project/../../libs/core/tools/clientgen/generate.sh"
spec="$project/internal/clientgen/specs/audit.yaml"
parent="$project/internal/generated"
output="$parent/audit"
mkdir -p "$parent"
if [[ -f "$output/api_default.go" && -f "$output/.source.yaml" ]] && cmp -s "$spec" "$output/.source.yaml"; then
  exit 0
fi
stage="$(mktemp -d "$parent/.generate.XXXXXX")"
trap 'rm -rf -- "$stage"' EXIT
bash "$generator" "$spec" "$stage/audit" auditclient
cp "$spec" "$stage/audit/.source.yaml"
# Do not remove the prior generated client until generation succeeds.
rm -rf -- "$output"
mv "$stage/audit" "$output"
