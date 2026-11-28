#!/usr/bin/env bash
set -Eeuo pipefail
set +x

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/env.sh"

if [[ "$#" -lt 4 || "$#" -gt 5 ]]; then
  printf 'Usage: bash scripts/identity/bootstrap-admin.sh EMAIL FIRST_NAME LAST_NAME TENANT_ID [PERSON_GID]\n' >&2
  exit 2
fi

load_backend_runtime_env
if [[ -z "${DATABASE_URL:-}" ]]; then
  printf 'DATABASE_URL must be configured in the repository .env\n' >&2
  exit 1
fi
if [[ ! -t 0 ]]; then
  printf 'Bootstrap requires an interactive terminal; refusing piped invocation\n' >&2
  exit 1
fi

printf 'Bootstrap the FIRST Identity account only. Email: %s | Tenant: %s\n' "$1" "$4"
printf 'New password (12-256 characters, hidden): '
IFS= read -r -s PASSWORD
printf '\n'
printf 'Confirm password (hidden): '
IFS= read -r -s CONFIRM_PASSWORD
printf '\n'
if [[ "${PASSWORD}" != "${CONFIRM_PASSWORD}" ]]; then
  printf 'Passwords do not match\n' >&2
  exit 1
fi

bash "${SCRIPT_DIR}/../backend/generate-annotations.sh"

args=(--confirm-bootstrap --password-stdin --email "$1" --first-name "$2" --last-name "$3" --tenant-id "$4")
if [[ "$#" -eq 5 ]]; then args+=(--person-gid "$5"); fi
(
  cd -- "${BACKEND_DIR}"
  printf '%s\n' "${PASSWORD}" | go run -overlay="${BACKEND_DIR}/.generated/overlay.json" ./cmd identity-bootstrap-admin "${args[@]}"
)
unset PASSWORD CONFIRM_PASSWORD
