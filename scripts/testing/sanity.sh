#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"

SANITY_SCOPE="${1:-./tests/sanity/...}"
REPORT_DIR="${RUNTIME_DIR}/reports/sanity"
REPORT_FILE="${REPORT_DIR}/index.html"
RAW_REPORT_FILE="${REPORT_DIR}/output.log"

open_report() {
  local file="$1"

  case "$(uname -s)" in
    Darwin)
      open "${file}" >/dev/null 2>&1 || log_warn "Unable to open sanity report automatically"
      ;;
    Linux)
      if command -v xdg-open >/dev/null 2>&1; then
        xdg-open "${file}" >/dev/null 2>&1 || log_warn "Unable to open sanity report automatically"
      elif command -v wslview >/dev/null 2>&1; then
        wslview "${file}" >/dev/null 2>&1 || log_warn "Unable to open sanity report automatically"
      else
        log_warn "No supported browser opener found; open ${file} manually"
      fi
      ;;
    MINGW*|MSYS*|CYGWIN*)
      if command -v cmd.exe >/dev/null 2>&1; then
        cmd.exe /C start "" "$(cygpath -w "${file}")" >/dev/null 2>&1 || log_warn "Unable to open sanity report automatically"
      else
        log_warn "No supported browser opener found; open ${file} manually"
      fi
      ;;
    *)
      log_warn "Unsupported host for automatic browser open; open ${file} manually"
      ;;
  esac
}

html_escape_file() {
  sed \
    -e 's/&/\&amp;/g' \
    -e 's/</\&lt;/g' \
    -e 's/>/\&gt;/g' \
    "$1"
}

mkdir -p -- "${REPORT_DIR}"
: >"${RAW_REPORT_FILE}"

require_command go

log_info "Ensuring NATS JetStream is available for sanity tests"
bash "${SCRIPT_DIR}/../infrastructure/nats/start.sh"

bash "${SCRIPT_DIR}/../backend/generate-annotations.sh"

log_info "Running backend sanity tests: ${SANITY_SCOPE}"

set +e
(
  cd "${BACKEND_DIR}"
  go test -vet=off -overlay="${BACKEND_DIR}/.generated/overlay.json" -count=1 -v "${SANITY_SCOPE}"
) 2>&1 | tee "${RAW_REPORT_FILE}"
test_status=${PIPESTATUS[0]}
set -e

if [[ "${test_status}" -eq 0 ]]; then
  result_label="PASS"
  result_class="pass"
else
  result_label="FAIL"
  result_class="fail"
fi

package_pass_count="$(grep -c '^ok[[:space:]]' "${RAW_REPORT_FILE}" || true)"
package_fail_count="$(grep -c '^FAIL[[:space:]]' "${RAW_REPORT_FILE}" || true)"
started_at="$(date '+%Y-%m-%d %H:%M:%S %Z')"

cat >"${REPORT_FILE}" <<EOF
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GroveHub Sanity Test Report</title>
  <style>
    :root { color-scheme: light dark; }
    body { margin: 0; font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #0b1020; color: #e8ecf5; }
    main { max-width: 1200px; margin: 0 auto; padding: 32px 20px 48px; }
    h1 { margin: 0 0 8px; font-size: 30px; }
    .muted { color: #9ca9c3; }
    .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 14px; margin: 24px 0; }
    .card { background: #131a2e; border: 1px solid #25304a; border-radius: 14px; padding: 18px; }
    .value { margin-top: 8px; font-size: 28px; font-weight: 700; }
    .pass { color: #58d68d; }
    .fail { color: #ff6b6b; }
    pre { margin: 0; white-space: pre-wrap; word-break: break-word; background: #090d18; border: 1px solid #25304a; border-radius: 14px; padding: 20px; overflow-x: auto; line-height: 1.45; }
    .section-title { margin: 30px 0 12px; font-size: 18px; }
    code { color: #c7d2fe; }
  </style>
</head>
<body>
<main>
  <h1>GroveHub Sanity Test Report</h1>
  <div class="muted">Generated ${started_at}</div>
  <div class="summary">
    <div class="card"><div class="muted">Result</div><div class="value ${result_class}">${result_label}</div></div>
    <div class="card"><div class="muted">Passed packages</div><div class="value">${package_pass_count}</div></div>
    <div class="card"><div class="muted">Failed packages</div><div class="value">${package_fail_count}</div></div>
    <div class="card"><div class="muted">Scope</div><div class="value" style="font-size:16px"><code>${SANITY_SCOPE}</code></div></div>
  </div>
  <div class="section-title">Complete test output</div>
  <pre>
EOF

html_escape_file "${RAW_REPORT_FILE}" >>"${REPORT_FILE}"

cat >>"${REPORT_FILE}" <<'EOF'
  </pre>
</main>
</body>
</html>
EOF

log_info "Sanity report generated: ${REPORT_FILE}"
open_report "${REPORT_FILE}"

exit "${test_status}"
