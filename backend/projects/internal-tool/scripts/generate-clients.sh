#!/usr/bin/env bash
set -euo pipefail

project="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
generator="${GROOVEHUB_CLIENT_GENERATOR:-$project/../../libs/core/tools/clientgen/generate.sh}"
providers="${GROOVEHUB_CONTRACT_ROOT:-$project/../platform/services}"
[[ -f "$generator" && -d "$providers" ]] || { echo "Missing pinned client generator or provider contracts" >&2; exit 1; }
parent="$project/internal/generated"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.generate.XXXXXX")"
trap 'rm -rf -- "$stage"' EXIT

# Providers own contracts. Do not maintain consumer-local copies of the specs.
for service in integration audit registry; do
  spec="$providers/$service/internal/resources/api/openapi/v1/openapi.yaml"
  if [[ "$service" == registry ]]; then
    spec="${GROOVEHUB_SCHOOL_CONTRACT:-$project/../school/services/registry/internal/resources/api/openapi/v1/openapi.yaml}"
  fi
  output="$parent/$service"
  if [[ -f "$output/client.go" && -f "$output/.source.yaml" && -f "$output/.generator.sh" ]] &&
    cmp -s "$spec" "$output/.source.yaml" &&
    cmp -s "$generator" "$output/.generator.sh"; then
    continue
  fi
  bash "$generator" "$spec" "$stage/$service" "${service}client"
  cp "$spec" "$stage/$service/.source.yaml"
  cp "$generator" "$stage/$service/.generator.sh"
done

# Generation of every changed SDK must succeed before replacing any existing SDK.
for service in integration audit registry; do
  if [[ -d "$stage/$service" ]]; then
    rm -rf -- "$parent/$service"
    mv -- "$stage/$service" "$parent/$service"
  fi
done
