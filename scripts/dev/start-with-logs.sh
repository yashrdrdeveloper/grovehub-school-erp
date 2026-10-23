#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/container.sh"

ENV_FILE="${REPO_ROOT}/.env"
COLLECTOR_CONFIG="${REPO_ROOT}/scripts/observability/otel-collector.local.yaml"
COLLECTOR_NAME="grovehub-otel-collector"
COLLECTOR_IMAGE="otel/opentelemetry-collector-contrib:0.147.0"

require_env() {
  local name="$1"
  [[ -n "${!name:-}" ]] || die "Missing required value in .env: ${name}"
}

ensure_runtime_dirs
require_docker

[[ -f "${ENV_FILE}" ]] || die "Missing ${ENV_FILE}. Copy .env.example to .env and add your Grafana Cloud credentials."
[[ -f "${COLLECTOR_CONFIG}" ]] || die "Missing collector config: ${COLLECTOR_CONFIG}"

set -a
# shellcheck disable=SC1090
source "${ENV_FILE}"
set +a

require_env GRAFANA_CLOUD_OTLP_ENDPOINT
require_env GRAFANA_CLOUD_INSTANCE_ID
require_env GRAFANA_CLOUD_API_KEY
require_env OTEL_SERVICE_NAME
require_env OTEL_SERVICE_NAMESPACE
require_env OTEL_DEPLOYMENT_ENVIRONMENT

export OTEL_SERVICE_VERSION="${OTEL_SERVICE_VERSION:-dev}"
export OTEL_EXPORTER_OTLP_ENDPOINT="${OTEL_EXPORTER_OTLP_ENDPOINT:-http://127.0.0.1:4318}"
export OTEL_EXPORTER_OTLP_PROTOCOL="${OTEL_EXPORTER_OTLP_PROTOCOL:-http/protobuf}"
export OTEL_RESOURCE_ATTRIBUTES="${OTEL_RESOURCE_ATTRIBUTES:-service.namespace=${OTEL_SERVICE_NAMESPACE},service.version=${OTEL_SERVICE_VERSION},deployment.environment.name=${OTEL_DEPLOYMENT_ENVIRONMENT}}"

if container_exists "${COLLECTOR_NAME}"; then
  docker rm --force "${COLLECTOR_NAME}" >/dev/null
fi

touch "${RUNTIME_LOG_DIR}/backend.log"

log_info "Starting local OpenTelemetry collector"
docker run --detach \
  --name "${COLLECTOR_NAME}" \
  --pull missing \
  --publish 127.0.0.1:4317:4317 \
  --publish 127.0.0.1:4318:4318 \
  --volume "${COLLECTOR_CONFIG}:/etc/otelcol-contrib/config.yaml:ro" \
  --volume "${RUNTIME_LOG_DIR}:/var/log/grovehub:ro" \
  --env GRAFANA_CLOUD_OTLP_ENDPOINT \
  --env GRAFANA_CLOUD_INSTANCE_ID \
  --env GRAFANA_CLOUD_API_KEY \
  --env OTEL_SERVICE_NAME \
  --env OTEL_SERVICE_NAMESPACE \
  --env OTEL_SERVICE_VERSION \
  --env OTEL_DEPLOYMENT_ENVIRONMENT \
  "${COLLECTOR_IMAGE}" \
  --config=/etc/otelcol-contrib/config.yaml >/dev/null

sleep 1
if ! container_is_running "${COLLECTOR_NAME}"; then
  docker logs "${COLLECTOR_NAME}" >&2 || true
  docker rm --force "${COLLECTOR_NAME}" >/dev/null 2>&1 || true
  die "OpenTelemetry collector failed to start"
fi

unset GRAFANA_CLOUD_API_KEY
unset GRAFANA_CLOUD_INSTANCE_ID
unset GRAFANA_CLOUD_OTLP_ENDPOINT

if ! bash "${SCRIPT_DIR}/start.sh"; then
  docker rm --force "${COLLECTOR_NAME}" >/dev/null 2>&1 || true
  exit 1
fi

log_info "Grafana Cloud observability export enabled"
log_info "OTLP endpoint for GroveHub: ${OTEL_EXPORTER_OTLP_ENDPOINT}"
log_info "Collector container: ${COLLECTOR_NAME}"
