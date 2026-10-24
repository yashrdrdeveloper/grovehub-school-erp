#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

source "${SCRIPT_DIR}/../lib/paths.sh"
source "${SCRIPT_DIR}/../lib/common.sh"
source "${SCRIPT_DIR}/../lib/json.sh"
source "${SCRIPT_DIR}/../lib/env.sh"
source "${SCRIPT_DIR}/../config/services.sh"

HARAKA_CONTAINER_NAME="grovehub-haraka"
CLOUDFLARE_API_BASE="https://api.cloudflare.com/client/v4"

load_backend_runtime_env

DOMAIN="${DOMAIN:-${1:-}}"
PUBLIC_IP="${PUBLIC_IP:-${2:-}}"
DKIM_SELECTOR="${MAIL_PROVIDER_DKIM_SELECTOR:-groovehub1}"

require_value() {
  local name="$1"
  local value="$2"
  [[ -n "${value}" ]] || die "Missing ${name}. Example: make mail-domain-ready DOMAIN=school.example PUBLIC_IP=203.0.113.10"
}

cloudflare_request() {
  local method="$1"
  local path="$2"
  local payload="${3:-}"
  local response
  local args=(
    --silent
    --show-error
    --request "${method}"
    --url "${CLOUDFLARE_API_BASE}${path}"
    --header "Authorization: Bearer ${CLOUDFLARE_API_TOKEN}"
    --header "Content-Type: application/json"
  )

  if [[ -n "${payload}" ]]; then
    args+=(--data "${payload}")
  fi

  response="$(curl "${args[@]}")"
  if [[ "$(jq -r '.success // false' <<<"${response}")" != "true" ]]; then
    log_error "Cloudflare API request failed: ${method} ${path}"
    jq -c '.errors // .' <<<"${response}" >&2 || true
    return 1
  fi

  printf '%s\n' "${response}"
}

find_zone_id() {
  if [[ -n "${CLOUDFLARE_ZONE_ID:-}" ]]; then
    printf '%s\n' "${CLOUDFLARE_ZONE_ID}"
    return 0
  fi

  local candidate="${DOMAIN}"
  local response
  local zone_id

  while [[ "${candidate}" == *.* ]]; do
    response="$(cloudflare_request GET "/zones?name=${candidate}&status=active&per_page=1")"
    zone_id="$(jq -r '.result[0].id // empty' <<<"${response}")"
    if [[ -n "${zone_id}" ]]; then
      printf '%s\n' "${zone_id}"
      return 0
    fi
    candidate="${candidate#*.}"
  done

  return 1
}

upsert_single_txt() {
  local zone_id="$1"
  local name="$2"
  local content="$3"
  local records
  local record_id
  local payload

  records="$(cloudflare_request GET "/zones/${zone_id}/dns_records?type=TXT&name=${name}&per_page=100")"
  record_id="$(jq -r '.result[0].id // empty' <<<"${records}")"
  payload="$(jq -nc \
    --arg name "${name}" \
    --arg content "${content}" \
    '{type:"TXT",name:$name,content:$content,ttl:1}')"

  if [[ -n "${record_id}" ]]; then
    cloudflare_request PUT "/zones/${zone_id}/dns_records/${record_id}" "${payload}" >/dev/null
    log_info "Updated TXT ${name}"
    return 0
  fi

  cloudflare_request POST "/zones/${zone_id}/dns_records" "${payload}" >/dev/null
  log_info "Created TXT ${name}"
}

create_txt() {
  local zone_id="$1"
  local name="$2"
  local content="$3"
  local payload

  payload="$(jq -nc \
    --arg name "${name}" \
    --arg content "${content}" \
    '{type:"TXT",name:$name,content:$content,ttl:1}')"
  cloudflare_request POST "/zones/${zone_id}/dns_records" "${payload}" >/dev/null
  log_info "Created TXT ${name}"
}

ensure_spf() {
  local zone_id="$1"
  local mechanism="$2"
  local response
  local existing
  local desired

  response="$(cloudflare_request GET "/zones/${zone_id}/dns_records?type=TXT&name=${DOMAIN}&per_page=100")"
  existing="$(jq -r '.result[] | select((.content // "") | startswith("v=spf1")) | .content' <<<"${response}" | head -n 1)"
  desired="v=spf1 ${mechanism} ~all"

  if [[ -z "${existing}" ]]; then
    create_txt "${zone_id}" "${DOMAIN}" "${desired}"
    return 0
  fi

  if [[ "${existing}" == *"${mechanism}"* ]]; then
    log_info "SPF already authorizes ${mechanism}"
    return 0
  fi

  log_warn "Existing SPF was preserved: ${existing}"
  log_warn "Add ${mechanism} to that existing SPF record before its final all mechanism. Do not create a second SPF record."
}

ensure_dmarc() {
  local zone_id="$1"
  local name="_dmarc.${DOMAIN}"
  local response
  local existing

  response="$(cloudflare_request GET "/zones/${zone_id}/dns_records?type=TXT&name=${name}&per_page=100")"
  existing="$(jq -r '.result[] | select((.content // "") | startswith("v=DMARC1")) | .content' <<<"${response}" | head -n 1)"
  if [[ -n "${existing}" ]]; then
    log_info "Existing DMARC policy preserved"
    return 0
  fi

  create_txt "${zone_id}" "${name}" "v=DMARC1; p=none; adkim=r; aspf=r"
}

ensure_haraka_dkim_key() {
  docker exec \
    --env "GROVEHUB_DOMAIN=${DOMAIN}" \
    --env "GROVEHUB_SELECTOR=${DKIM_SELECTOR}" \
    "${HARAKA_CONTAINER_NAME}" \
    sh -eu -c '
      key_dir="/opt/haraka/config/dkim/${GROVEHUB_DOMAIN}"
      mkdir -p "${key_dir}"
      if [ ! -s "${key_dir}/private" ]; then
        openssl genrsa -out "${key_dir}/private" 2048 >/dev/null 2>&1
      fi
      chmod 600 "${key_dir}/private"
      printf "%s\n" "${GROVEHUB_SELECTOR}" >"${key_dir}/selector"
      openssl rsa -in "${key_dir}/private" -pubout -out "${key_dir}/public" >/dev/null 2>&1
    '

  docker exec \
    --env "GROVEHUB_DOMAIN=${DOMAIN}" \
    "${HARAKA_CONTAINER_NAME}" \
    sh -eu -c '
      sed -e "/-----BEGIN PUBLIC KEY-----/d" -e "/-----END PUBLIC KEY-----/d" \
        "/opt/haraka/config/dkim/${GROVEHUB_DOMAIN}/public" | tr -d "\r\n"
    '
}

wait_for_haraka() {
  local attempt
  local health

  for ((attempt = 1; attempt <= 30; attempt++)); do
    health="$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "${HARAKA_CONTAINER_NAME}" 2>/dev/null || true)"
    if [[ "${health}" == "healthy" ]]; then
      return 0
    fi
    if [[ "${health}" == "unhealthy" || "${health}" == "exited" ]]; then
      return 1
    fi
    sleep 1
  done
  return 1
}

lookup_or_register_domain() {
  local response
  local error_code

  response="$(curl --silent --show-error "${BACKEND_BASE_URL}/mail-provider/domains/by-name/${DOMAIN}")"
  if [[ "$(jq -r '.success // false' <<<"${response}")" == "true" ]]; then
    printf '%s\n' "${response}"
    return 0
  fi

  error_code="$(jq -r '.errorCode // empty' <<<"${response}")"
  if [[ "${error_code}" != "MAIL_PROVIDER_DOMAIN_NOT_FOUND" ]]; then
    log_error "Unable to lookup mail-provider domain"
    jq -c '.' <<<"${response}" >&2 || true
    return 1
  fi

  response="$(curl --silent --show-error \
    --request POST \
    --url "${BACKEND_BASE_URL}/mail-provider/domains" \
    --header 'Content-Type: application/json' \
    --data "$(jq -nc --arg name "${DOMAIN}" '{name:$name}')")"
  if [[ "$(jq -r '.success // false' <<<"${response}")" != "true" ]]; then
    log_error "Unable to register mail-provider domain"
    jq -c '.' <<<"${response}" >&2 || true
    return 1
  fi

  printf '%s\n' "${response}"
}

verify_domain() {
  local domain_id="$1"
  local attempt
  local response
  local error_code

  for ((attempt = 1; attempt <= 36; attempt++)); do
    response="$(curl --silent --show-error \
      --request POST \
      --url "${BACKEND_BASE_URL}/mail-provider/domains/${domain_id}/verify")"
    if [[ "$(jq -r '.success // false' <<<"${response}")" == "true" ]]; then
      printf '%s\n' "${response}"
      return 0
    fi

    error_code="$(jq -r '.errorCode // empty' <<<"${response}")"
    if [[ "${error_code}" != "MAIL_PROVIDER_DOMAIN_NOT_VERIFIED" ]]; then
      log_error "Mail Provider rejected domain verification"
      jq -c '.' <<<"${response}" >&2 || true
      return 1
    fi
    sleep 5
  done

  return 1
}

require_value DOMAIN "${DOMAIN}"
require_value PUBLIC_IP "${PUBLIC_IP}"
require_value CLOUDFLARE_API_TOKEN "${CLOUDFLARE_API_TOKEN:-}"

require_command curl
ensure_jq
require_docker
ensure_runtime_dirs

bash "${SCRIPT_DIR}/start.sh"

if ! curl --fail --silent --max-time 1 "${BACKEND_BASE_URL}/health" >/dev/null 2>&1; then
  log_info "Backend is not running; applying Mail Provider migrations and starting it"
  bash "${SCRIPT_DIR}/../database/migrate.sh" up mailProvider
  BACKEND_HOT_RELOAD=1 bash "${SCRIPT_DIR}/../backend/start.sh"
fi

log_info "Checking Cloudflare API token"
cloudflare_request GET "/user/tokens/verify" >/dev/null

ZONE_ID="$(find_zone_id || true)"
[[ -n "${ZONE_ID}" ]] || die "No active Cloudflare zone was found for ${DOMAIN}. Set CLOUDFLARE_ZONE_ID explicitly if the zone is managed under a parent domain."
log_info "Cloudflare zone resolved: ${ZONE_ID}"

DOMAIN_RESPONSE="$(lookup_or_register_domain)"
DOMAIN_ID="$(jq -r '.data.domainId' <<<"${DOMAIN_RESPONSE}")"
DOMAIN_STATUS="$(jq -r '.data.status' <<<"${DOMAIN_RESPONSE}")"
VERIFICATION_HOST="$(jq -r '.data.verificationHost' <<<"${DOMAIN_RESPONSE}")"
VERIFICATION_VALUE="$(jq -r '.data.verificationValue' <<<"${DOMAIN_RESPONSE}")"

if [[ "${DOMAIN_STATUS}" == "DEACTIVATED" ]]; then
  die "Mail Provider domain ${DOMAIN} is deactivated and cannot be prepared"
fi

log_info "Preparing DKIM key for ${DOMAIN}"
DKIM_PUBLIC_KEY="$(ensure_haraka_dkim_key)"
[[ -n "${DKIM_PUBLIC_KEY}" ]] || die "Unable to read generated DKIM public key"
DKIM_HOST="${DKIM_SELECTOR}._domainkey.${DOMAIN}"
DKIM_VALUE="v=DKIM1; k=rsa; p=${DKIM_PUBLIC_KEY}"

log_info "Publishing GroveHub domain verification and DKIM records"
upsert_single_txt "${ZONE_ID}" "${VERIFICATION_HOST}" "${VERIFICATION_VALUE}"
upsert_single_txt "${ZONE_ID}" "${DKIM_HOST}" "${DKIM_VALUE}"

if [[ "${PUBLIC_IP}" == *:* ]]; then
  SPF_MECHANISM="ip6:${PUBLIC_IP}"
else
  SPF_MECHANISM="ip4:${PUBLIC_IP}"
fi
ensure_spf "${ZONE_ID}" "${SPF_MECHANISM}"
ensure_dmarc "${ZONE_ID}"

log_info "Reloading Haraka with the domain DKIM key"
docker restart "${HARAKA_CONTAINER_NAME}" >/dev/null
if ! wait_for_haraka; then
  docker logs "${HARAKA_CONTAINER_NAME}" >&2 || true
  die "Haraka did not become healthy after loading DKIM"
fi

log_info "Waiting for the GroveHub ownership TXT record to become resolvable"
if ! VERIFIED_RESPONSE="$(verify_domain "${DOMAIN_ID}")"; then
  log_error "Cloudflare records were created, but DNS verification has not propagated yet"
  log_error "Retry: curl --request POST --url ${BACKEND_BASE_URL}/mail-provider/domains/${DOMAIN_ID}/verify"
  exit 1
fi

log_info "Mail Provider domain is verified: $(jq -r '.data.name' <<<"${VERIFIED_RESPONSE}")"
log_info "Domain ID: ${DOMAIN_ID}"
log_info "DKIM: ${DKIM_HOST}"
log_info "Cloudflare DNS is configured for GroveHub Mail Provider"
log_warn "Cloudflare cannot configure reverse DNS/PTR for your server IP. Before production delivery, set the PTR/rDNS for ${PUBLIC_IP} with the VPS/cloud provider and make the Haraka public hostname match it."
