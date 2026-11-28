#!/usr/bin/env bash
set -euo pipefail

# Exercise the standalone adapter without contacting Docker or a registry.
script="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)/generate.sh"
fixture="$(mktemp -d)"
trap 'rm -rf -- "$fixture"' EXIT
mkdir -p "$fixture/bin" "$fixture/output"
printf 'openapi: 3.0.3\n' > "$fixture/openapi.yaml"
export CLIENTGEN_TEST_LOG="$fixture/docker.log"
export PATH="$fixture/bin:$PATH"

cat > "$fixture/bin/docker" <<'DOCKER'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$@" > "$CLIENTGEN_TEST_LOG"
stage=""
for arg in "$@"; do
  case "$arg" in
    *:/output) stage="${arg%:/output}" ;;
  esac
done
test -n "$stage"
case "${CLIENTGEN_TEST_MODE:-success}" in
  failure) exit 17 ;;
  empty) exit 0 ;;
  success) printf 'package generated\n' > "$stage/client.go" ;;
  *) exit 99 ;;
esac
DOCKER
chmod +x "$fixture/bin/docker"

fail_expected() {
  if bash "$script" "$@" > "$fixture/stdout" 2> "$fixture/stderr"; then
    echo "Expected generator failure" >&2
    exit 1
  fi
}

assert_no_staging() {
  if find "$fixture/output" -maxdepth 1 -name '.clientgen.*' -print -quit | grep -q .; then
    echo "Generator left a temporary directory" >&2
    exit 1
  fi
}

fail_expected
fail_expected "$fixture/missing.yaml" "$fixture/output/missing" valid
fail_expected "$fixture/openapi.yaml" "$fixture/output/invalid" Invalid
test ! -e "$CLIENTGEN_TEST_LOG"

bash "$script" "$fixture/openapi.yaml" "$fixture/output/generated" valid
test -f "$fixture/output/generated/client.go"
grep -Fxq -- '--network' "$CLIENTGEN_TEST_LOG"
grep -Fxq -- 'none' "$CLIENTGEN_TEST_LOG"
grep -Fxq -- 'openapitools/openapi-generator-cli:v7.15.0' "$CLIENTGEN_TEST_LOG"
grep -Fxq -- "$fixture/openapi.yaml:/input/openapi.yaml:ro" "$CLIENTGEN_TEST_LOG"
grep -Fxq -- 'packageName=valid,isGoSubmodule=true,hideGenerationTimestamp=true,withGoMod=false' "$CLIENTGEN_TEST_LOG"
grep -Fxq -- 'apiTests=false,modelTests=false,apiDocs=false,modelDocs=false' "$CLIENTGEN_TEST_LOG"
assert_no_staging

printf 'preserve me\n' > "$fixture/output/generated/client.go"
fail_expected "$fixture/openapi.yaml" "$fixture/output/generated" valid
grep -Fxq 'preserve me' "$fixture/output/generated/client.go"

export CLIENTGEN_TEST_MODE=failure
fail_expected "$fixture/openapi.yaml" "$fixture/output/failure" valid
test ! -e "$fixture/output/failure"
assert_no_staging

export CLIENTGEN_TEST_MODE=empty
fail_expected "$fixture/openapi.yaml" "$fixture/output/empty" valid
test ! -e "$fixture/output/empty"
assert_no_staging

echo "Client generator tests passed."
