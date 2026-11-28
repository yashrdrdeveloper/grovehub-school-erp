#!/usr/bin/env bash
set -Eeuo pipefail

module_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)"
export GOWORK=off
cd -- "$module_dir"

profile="${1:-$module_dir/.generated/coverage.out}"
mkdir -p -- "$(dirname -- "$profile")"
# Include every production package and use all standalone test trees.
go test -race -count=1 -covermode=atomic -coverpkg=./session/... -coverprofile="$profile" ./...
go tool cover -func="$profile"
awk '
  NR == 1 { next }
  { statements[$1] = $2; if ($3 > 0) covered[$1] = 1 }
  END {
    for (block in statements) {
      total += statements[block]
      if (covered[block]) hit += statements[block]
      else missing += statements[block]
    }
    printf "Platform library: %d/%d covered statements, %d uncovered\n", hit, total, missing
    if (total == 0 || missing != 0) exit 1
  }
' "$profile"
