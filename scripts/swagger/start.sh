#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/container.sh"
source "${SCRIPT_DIR}/../config/services.sh"

require_docker
require_command curl

if container_is_running "${SWAGGER_CONTAINER_NAME}"; then
  log_info "Swagger UI is already running: ${SWAGGER_BASE_URL}"
  exit 0
fi

if container_exists "${SWAGGER_CONTAINER_NAME}"; then
  docker rm "${SWAGGER_CONTAINER_NAME}" >/dev/null
fi

if curl --fail --silent --max-time 1 "${SWAGGER_BASE_URL}" >/dev/null 2>&1; then
  die "Swagger port ${SWAGGER_PORT} is already in use by another service"
fi

swagger_urls="["
separator=""
for spec in "${SWAGGER_SPECS[@]}"; do
  IFS='|' read -r spec_name spec_path <<<"${spec}"
  [[ -n "${spec_name}" && -n "${spec_path}" ]] || die "Invalid Swagger spec registration: ${spec}"
  [[ -f "${REPO_ROOT}/backend/services/${spec_path}" ]] || die "Swagger spec not found: backend/services/${spec_path}"
  swagger_urls+="${separator}{url:\"/specs/${spec_path}\",name:\"${spec_name}\"}"
  separator=","
done
swagger_urls+="]"

log_info "Starting Swagger UI (${SWAGGER_IMAGE})"
docker run --detach \
  --name "${SWAGGER_CONTAINER_NAME}" \
  --publish "127.0.0.1:${SWAGGER_PORT}:8080" \
  --env "URLS=${swagger_urls}" \
  --env "DEEP_LINKING=true" \
  --volume "${REPO_ROOT}/backend/services:/usr/share/nginx/html/specs:ro" \
  "${SWAGGER_IMAGE}" >/dev/null

if ! wait_for_http "${SWAGGER_BASE_URL}" "Swagger UI" 30; then
  log_error "Swagger UI failed to become ready"
  docker logs --tail 50 "${SWAGGER_CONTAINER_NAME}" >&2 || true
  docker rm --force "${SWAGGER_CONTAINER_NAME}" >/dev/null 2>&1 || true
  exit 1
fi

log_info "Swagger UI: ${SWAGGER_BASE_URL}"
