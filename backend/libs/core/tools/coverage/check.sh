#!/usr/bin/env bash
# Check exact covered statements; never rely on rounded percentages.
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/../.."
export GOWORK=off
mode="${1:-}"
out="${2:-$PWD/.coverage}"
mkdir -p "$out"
out="$(cd "$out" && pwd)"
check() {
  local profile="$1" label="$2"
  test -s "$profile" || { echo "$label: missing coverage profile" >&2; return 1; }
  go tool cover -func="$profile"
  awk -v label="$label" '
    NR == 1 { if ($0 != "mode: atomic" && $0 != "mode: count" && $0 != "mode: set") bad=1; next }
    NF != 3 || $2 !~ /^[0-9]+$/ || $3 !~ /^[0-9]+$/ { bad=1; next }
    { statements[$1]=$2; counts[$1]+=$3 }
    END {
      for (block in statements) {
        total+=statements[block]
        if (counts[block]>0) covered+=statements[block]
        else if (statements[block]>0) {
          print "UNCOVERED " block " (" statements[block] " statements)"
          missing+=statements[block]
        }
      }
      if (total == 0) print label ": no executable statements"
      else printf "%s: %d/%d statements covered (%.4f%%)\n",label,covered,total,100*covered/total
      if (bad || missing) exit 1
    }
  ' "$profile"
}
run_scope() {
  local label="$1" packages="$2"
  shift 2
  local profile="$out/$label.out"
  # Include the instrumented source packages as well as the owning tests,
  # so entirely unimported source packages cannot disappear from the report.
  rm -f "$profile"
  go test -count=1 -covermode=atomic -coverpkg="$packages" -coverprofile="$profile" "$@" || return 1
  check "$profile" "$label"
}
case "$mode" in
  list)
    for dir in */internal/src; do
      cap="${dir%%/*}"
      test -f "$cap/$cap.go" || { echo "Missing capability facade: $cap" >&2; exit 1; }
      printf '%s\n' "$cap"
    done
    printf '%s\n' tools test-support
    ;;
  tools)
    status=0
    run_scope modulepack-internal ./tools/modulepack/internal/src/... ./tools/modulepack/internal/src/... ./tools/modulepack/internal/resources/tests/unit-test/... || status=1
    for command in annotation modulepack; do
      case "$command" in
        annotation) target=./annotation/cmd ;;
        modulepack) target=./tools/modulepack ;;
      esac
      cov="$out/$command-cli-data"
      rm -rf "$cov"; mkdir -p "$cov"
      bin="$out/$command"
      go build -cover -covermode=atomic -coverpkg="$target" -o "$bin" "$target" || { status=1; continue; }
      if [[ "$command" == annotation ]]; then
        GOCOVERDIR="$cov" "$bin" --help || status=1
        if GOCOVERDIR="$cov" "$bin"; then status=1; else test "$?" -eq 2 || status=1; fi
      else
        src="$(mktemp -d)"; proxy="$(mktemp -d)"
        printf 'module example.com/corefixture\n\ngo 1.27.0\n' > "$src/go.mod"
        printf 'package corefixture\n' > "$src/core.go"
        GOCOVERDIR="$cov" "$bin" "$src" "$proxy" v0.0.1 2026-11-28T00:00:00Z || status=1
        if GOCOVERDIR="$cov" "$bin"; then status=1; else test "$?" -eq 1 || status=1; fi
        rm -rf "$src" "$proxy"
      fi
      go tool covdata textfmt -i="$cov" -o="$out/$command-cli.out" && check "$out/$command-cli.out" "$command-cli" || status=1
    done
    bash tools/clientgen/internal/resources/tests/unit-test/generate_test.sh || status=1
    exit "$status"
    ;;
  test-support)
    run_scope test-support './tests/helpers/...,./tests/mocks/...' ./tests/helpers/... ./tests/mocks/... ./tests/unit-test/test-support/...
    ;;
  *)
    [[ "$mode" =~ ^[a-z][a-z0-9]*$ && -f "$mode/$mode.go" && -d "$mode/internal/src" && -d "tests/unit-test/$mode" ]] || { echo "Unknown capability: $mode" >&2; exit 2; }
    status=0
    run_scope "$mode-internal" "./$mode/internal/src/..." "./$mode/internal/src/..." "./$mode/internal/resources/tests/unit-test/..." || status=1
    run_scope "$mode-public" "./$mode" "./$mode" "./tests/unit-test/$mode/..." || status=1
    exit "$status"
    ;;
esac
