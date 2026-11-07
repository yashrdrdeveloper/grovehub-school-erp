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

log_info "Swagger configuration: ${SWAGGER_BASE_URL} -> container port 8080 (Docker localhost publish)"

if container_is_running "${SWAGGER_CONTAINER_NAME}"; then
  if curl --noproxy '*' --fail --silent --max-time 3 "${SWAGGER_BASE_URL}/index.html" >/dev/null 2>&1; then
    log_info "Swagger UI is already running: ${SWAGGER_BASE_URL}"
    exit 0
  fi
  log_warn "Swagger container is running but the host HTTP endpoint is unreachable"
  docker port "${SWAGGER_CONTAINER_NAME}" >&2 || true
  die "Existing Swagger container failed readiness; inspect it before replacing it (docker logs ${SWAGGER_CONTAINER_NAME})"
fi

if container_exists "${SWAGGER_CONTAINER_NAME}"; then
  docker rm "${SWAGGER_CONTAINER_NAME}" >/dev/null
fi

if command -v lsof >/dev/null 2>&1 && lsof -nP -t -iTCP:"${SWAGGER_PORT}" -sTCP:LISTEN 2>/dev/null | grep -q .; then
  log_error "Another process is listening on Swagger port ${SWAGGER_PORT}:"
  lsof -nP -iTCP:"${SWAGGER_PORT}" -sTCP:LISTEN >&2 || true
  die "Swagger port ${SWAGGER_PORT} is occupied; choose an unused port or stop the confirmed process"
fi
if curl --noproxy '*' --silent --max-time 2 "${SWAGGER_BASE_URL}/index.html" >/dev/null 2>&1; then
  die "Swagger port ${SWAGGER_PORT} already responds to HTTP; refusing to replace it"
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
  --env "PORT=8080" \
  --env "URLS=${swagger_urls}" \
  --env "DEEP_LINKING=true" \
  --volume "${REPO_ROOT}/backend/services:/usr/share/nginx/html/specs:ro" \
  "${SWAGGER_IMAGE}" >/dev/null

swagger_ready=0
for ((attempt=1; attempt<=60; attempt++)); do
  if curl --noproxy '*' --fail --silent --max-time 2 "${SWAGGER_BASE_URL}/index.html" >/dev/null 2>&1; then
    swagger_ready=1
    break
  fi
  if ! container_is_running "${SWAGGER_CONTAINER_NAME}"; then
    log_error "Swagger container exited during startup"
    break
  fi
  sleep 1
done

if [[ "${swagger_ready}" != "1" ]]; then
  log_error "Swagger UI cannot be reached on the host at ${SWAGGER_BASE_URL}/index.html"
  log_error "Docker context: $(docker context show)"
  log_error "Docker published ports:"
  docker port "${SWAGGER_CONTAINER_NAME}" >&2 || true
  log_error "Docker container state:"
  docker inspect --format 'status={{.State.Status}} exit={{.State.ExitCode}} network={{.HostConfig.NetworkMode}} ports={{json .NetworkSettings.Ports}}' "${SWAGGER_CONTAINER_NAME}" >&2 || true
  log_error "HTTP diagnostic (direct connection without proxy):"
  curl --noproxy '*' --verbose --location --max-time 5 "${SWAGGER_BASE_URL}/index.html" -o /dev/null >&2 || true
  log_error "Nginx listen directive:"
  docker exec "${SWAGGER_CONTAINER_NAME}" sh -c 'grep -n "listen " /etc/nginx/conf.d/default.conf' >&2 || true
  log_error "HTTP diagnostic from inside container:"
  docker exec "${SWAGGER_CONTAINER_NAME}" sh -c 'wget -S -O /dev/null http://127.0.0.1:8080/index.html 2>&1' >&2 || true
  log_error "Container log tail:"
  docker logs --tail 35 "${SWAGGER_CONTAINER_NAME}" >&2 || true
  die "Swagger readiness failed. Container retained as ${SWAGGER_CONTAINER_NAME} for diagnosis; use make swagger-stop after inspecting."
fi
# Confirm each registered OpenAPI document is served, not just the static UI.
for spec in "${SWAGGER_SPECS[@]}"; do
  IFS='|' read -r spec_name spec_path <<<"${spec}"
  if ! curl --noproxy '*' --fail --silent --max-time 5 "${SWAGGER_BASE_URL}/specs/${spec_path}" >/dev/null; then
    log_error "Swagger UI is reachable, but the OpenAPI spec is not: ${spec_name} (${spec_path})"
    log_error "Verify the bind mount under ${REPO_ROOT}/backend/services and the container file path"
    die "Swagger OpenAPI specification readiness failed; container retained for diagnosis."
  fi
done

log_info "Swagger UI is ready: ${SWAGGER_BASE_URL}"
log_info "Swagger UI: ${SWAGGER_BASE_URL}"
