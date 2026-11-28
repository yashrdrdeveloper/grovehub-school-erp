#!/usr/bin/env bash
set -euo pipefail

export GOWORK=off
module="$(go list -m)"
root="github.com/yashrdrdeveloper/grovehub-school-erp"
case "$module" in
  "$root/backend/projects/school"|"$root/backend/projects/internal-tool") ;;
  *) echo "Unexpected consumer module: $module" >&2; exit 1 ;;
esac
core="grovehub/libs/core"

# Include test imports and transitives; generated SDKs belong to the consumer.
deps="$(go list -mod=mod -deps -test -f '{{if .Module}}{{.Module.Path}}{{end}}' ./... | sort -u)"
while IFS= read -r dep; do
  case "$dep" in
    ""|"$module"|"$core") ;;
    grovehub|grovehub/*|"$root"|"$root/"*) echo "Forbidden project dependency: $dep" >&2; exit 1 ;;
  esac
done <<< "$deps"
