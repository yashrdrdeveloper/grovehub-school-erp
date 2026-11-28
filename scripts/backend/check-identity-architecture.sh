#!/usr/bin/env bash
set -Eeuo pipefail

service="backend/projects/platform/services/identity"
app="$service/internal/app"

for forbidden in \
  "$service/cmd" \
  "$app/bootstrap/impl/bootstrap.go" \
  "$app/security" \
  "$app/controller/account/public" \
  "$service/client/account/impl"; do
  if [[ -e "$forbidden" ]]; then
    printf '[ERROR] Forbidden Identity architecture path: %s\n' "$forbidden" >&2
    exit 1
  fi
done

for required in \
  "$app/bootstrap/impl/build.go" \
  "$app/bootstrap/impl/runtime.go" \
  "$app/bootstrap/impl/client.go" \
  "$app/services/security/security.go" \
  "$app/validators/account/helpers/search.go" \
  "$app/controller/activation/public/activation.go" \
  "$app/client/account/impl/account.go" \
  "backend/projects/platform/cmd/internal/src/strategies/bootstrapadmin/bootstrapadmin.go"; do
  if [[ ! -f "$required" ]]; then
    printf '[ERROR] Required Identity architecture path is missing: %s\n' "$required" >&2
    exit 1
  fi
done

if grep -R --include='*.go' -nE '@Proxy|IDENTITY_SYSTEM_MAIL_TENANT_ID|IDENTITY_SYSTEM_MAIL_CONNECTION_ID' "$service/client/account" "$app/config"; then
  printf '[ERROR] Legacy Identity account proxy or fixed mail environment selection is forbidden\n' >&2
  exit 1
fi

if grep -E 'func \(module \*runtime\) (Get|Search)\(' "$app/bootstrap/impl/runtime.go"; then
  printf '[ERROR] Identity runtime must not delegate account service operations\n' >&2
  exit 1
fi

printf '[INFO] Identity architecture consistency check passed\n'
