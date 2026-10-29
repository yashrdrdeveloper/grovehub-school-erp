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
MAIL_PROVIDER_STATE_DIR="${STATE_DIR}/mailProvider"
MAIL_PROVIDER_RUNTIME_ENV="${MAIL_PROVIDER_STATE_DIR}/runtime.env"
BACKEND_NO_PROXY="localhost,127.0.0.1,::1"

load_backend_runtime_env

DKIM_SELECTOR="${MAIL_PROVIDER_DKIM_SELECTOR:-groovehub1}"

require_value() {
  local name="$1"
  local value="$2"
  [[ -n "${value}" ]] || die "Missing ${name}. Add it to the repository root .env file."
}

backend_request() {
  local method="$1"
  local path="$2"
  local payload="${3:-}"
  local response
  local args=(
    --silent
    --show-error
    --connect-timeout 5
    --max-time 15
    --noproxy "${BACKEND_NO_PROXY}"
    --request "${method}"
    --url "${BACKEND_BASE_URL}${path}"
  )

  if [[ -n "${payload}" ]]; then
    args+=(
      --header "Content-Type: application/json"
      --data "${payload}"
    )
  fi

  if response="$(curl "${args[@]}")"; then
    printf '%s\n' "${response}"
    return 0
  fi

  local exit_code=$?
  log_warn "Backend request failed (${method} ${path}, curl exit ${exit_code})"
  return "${exit_code}"
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

discover_zone() {
  local response
  local status

  response="$(cloudflare_request GET "/zones/${CLOUDFLARE_ZONE_ID}")"
  DOMAIN="$(jq -r '.result.name // empty' <<<"${response}")"
  status="$(jq -r '.result.status // empty' <<<"${response}")"

  [[ -n "${DOMAIN}" ]] || die "Cloudflare zone ${CLOUDFLARE_ZONE_ID} did not return a zone name"
  [[ "${status}" == "active" ]] || die "Cloudflare zone ${DOMAIN} is not active (status: ${status:-unknown})"

  MAIL_HOST="mail.${DOMAIN}"
}

discover_public_ip() {
  if [[ -n "${MAIL_PROVIDER_PUBLIC_IP:-}" ]]; then
    PUBLIC_IP="${MAIL_PROVIDER_PUBLIC_IP}"
    log_info "Using MAIL_PROVIDER_PUBLIC_IP override: ${PUBLIC_IP}"
    return 0
  fi

  local trace
  trace="$(curl -4 --fail --silent --show-error --max-time 10 \
    "https://www.cloudflare.com/cdn-cgi/trace" 2>/dev/null || true)"
  PUBLIC_IP="$(awk -F= '$1 == "ip" { print $2; exit }' <<<"${trace}")"

  if [[ -z "${PUBLIC_IP}" ]]; then
    trace="$(curl -6 --fail --silent --show-error --max-time 10 \
      "https://www.cloudflare.com/cdn-cgi/trace" 2>/dev/null || true)"
    PUBLIC_IP="$(awk -F= '$1 == "ip" { print $2; exit }' <<<"${trace}")"
  fi

  [[ -n "${PUBLIC_IP}" ]] || die "Unable to discover this machine's public IP. Set MAIL_PROVIDER_PUBLIC_IP in .env to override discovery."
  log_info "Discovered machine public IP: ${PUBLIC_IP}"
}

persist_mail_runtime() {
  mkdir -p -- "${MAIL_PROVIDER_STATE_DIR}"
  cat >"${MAIL_PROVIDER_RUNTIME_ENV}" <<EOF
MAIL_PROVIDER_PUBLIC_IP=${PUBLIC_IP}
MAIL_PROVIDER_HARAKA_HOSTNAME=${MAIL_HOST}
EOF
  chmod 600 "${MAIL_PROVIDER_RUNTIME_ENV}" 2>/dev/null || true
}

upsert_address_record() {
  local zone_id="$1"
  local type="$2"
  local name="$3"
  local content="$4"
  local records
  local count
  local record_id
  local payload

  records="$(cloudflare_request GET "/zones/${zone_id}/dns_records?type=${type}&name=${name}&per_page=100")"
  count="$(jq '.result | length' <<<"${records}")"
  if (( count > 1 )); then
    die "Multiple ${type} records already exist for ${name}; refusing to guess which one GroveHub should manage"
  fi

  record_id="$(jq -r '.result[0].id // empty' <<<"${records}")"
  payload="$(jq -nc \
    --arg type "${type}" \
    --arg name "${name}" \
    --arg content "${content}" \
    '{type:$type,name:$name,content:$content,ttl:1,proxied:false}')"

  if [[ -n "${record_id}" ]]; then
    cloudflare_request PUT "/zones/${zone_id}/dns_records/${record_id}" "${payload}" >/dev/null
    log_info "Updated ${type} ${name} -> ${content} (DNS only)"
  else
    cloudflare_request POST "/zones/${zone_id}/dns_records" "${payload}" >/dev/null
    log_info "Created ${type} ${name} -> ${content} (DNS only)"
  fi
}

ensure_mail_host_record() {
  local zone_id="$1"
  local conflict_records
  local conflict
  local type
  local other_type

  conflict_records="$(cloudflare_request GET "/zones/${zone_id}/dns_records?name=${MAIL_HOST}&per_page=100")"
  conflict="$(jq -r '.result[] | select(.type == "CNAME") | .content' <<<"${conflict_records}" | head -n 1)"
  [[ -z "${conflict}" ]] || die "${MAIL_HOST} is already a CNAME to ${conflict}; remove or rename it before using this hostname for GroveHub Mail"

  if [[ "${PUBLIC_IP}" == *:* ]]; then
    type="AAAA"
    other_type="A"
  else
    [[ "${PUBLIC_IP}" =~ ^[0-9]{1,3}(\.[0-9]{1,3}){3}$ ]] || die "Discovered public IP is not a valid IPv4/IPv6 address: ${PUBLIC_IP}"
    type="A"
    other_type="AAAA"
  fi

  upsert_address_record "${zone_id}" "${type}" "${MAIL_HOST}" "${PUBLIC_IP}"

  if jq -e --arg type "${other_type}" '.result[] | select(.type == $type)' <<<"${conflict_records}" >/dev/null 2>&1; then
    log_warn "${MAIL_HOST} also has an existing ${other_type} record. GroveHub left it untouched; ensure that address also belongs to this mail server or remove it."
  fi
}

upsert_managed_txt() {
  local zone_id="$1"
  local name="$2"
  local content="$3"
  local prefix="$4"
  local records
  local record_ids
  local record_id
  local stale_id
  local payload

  records="$(cloudflare_request GET "/zones/${zone_id}/dns_records?type=TXT&name=${name}&per_page=100")"
  record_ids="$(jq -r --arg prefix "${prefix}" \
    '.result[] | select((.content // "") | startswith($prefix)) | .id' <<<"${records}")"
  record_id="$(head -n 1 <<<"${record_ids}")"
  payload="$(jq -nc \
    --arg name "${name}" \
    --arg content "${content}" \
    '{type:"TXT",name:$name,content:$content,ttl:1}')"

  if [[ -n "${record_id}" ]]; then
    cloudflare_request PUT "/zones/${zone_id}/dns_records/${record_id}" "${payload}" >/dev/null
    log_info "Updated TXT ${name}"
  else
    cloudflare_request POST "/zones/${zone_id}/dns_records" "${payload}" >/dev/null
    log_info "Created TXT ${name}"
  fi

  while IFS= read -r stale_id; do
    [[ -n "${stale_id}" ]] || continue
    [[ "${stale_id}" == "${record_id}" ]] && continue
    cloudflare_request DELETE "/zones/${zone_id}/dns_records/${stale_id}" >/dev/null
    log_info "Removed stale managed TXT ${name}"
  done <<<"${record_ids}"
}

ensure_spf() {
  local zone_id="$1"
  local mechanism="$2"
  local response
  local count
  local record_id
  local existing
  local updated
  local last
  local payload

  response="$(cloudflare_request GET "/zones/${zone_id}/dns_records?type=TXT&name=${DOMAIN}&per_page=100")"
  count="$(jq '[.result[] | select((.content // "") | startswith("v=spf1"))] | length' <<<"${response}")"
  if (( count > 1 )); then
    die "Multiple SPF records exist for ${DOMAIN}. SPF permits only one record; merge them before continuing."
  fi

  record_id="$(jq -r '.result[] | select((.content // "") | startswith("v=spf1")) | .id' <<<"${response}" | head -n 1)"
  existing="$(jq -r '.result[] | select((.content // "") | startswith("v=spf1")) | .content' <<<"${response}" | head -n 1)"

  if [[ -z "${existing}" ]]; then
    upsert_managed_txt "${zone_id}" "${DOMAIN}" "v=spf1 ${mechanism} ~all" "v=spf1"
    return 0
  fi

  if [[ "${existing}" == *"${mechanism}"* ]]; then
    log_info "SPF already authorizes ${mechanism}"
    return 0
  fi

  last="${existing##* }"
  case "${last}" in
    +all|-all|~all|\?all)
      updated="${existing% *} ${mechanism} ${last}"
      ;;
    *)
      updated="${existing} ${mechanism}"
      ;;
  esac

  payload="$(jq -nc \
    --arg name "${DOMAIN}" \
    --arg content "${updated}" \
    '{type:"TXT",name:$name,content:$content,ttl:1}')"
  cloudflare_request PUT "/zones/${zone_id}/dns_records/${record_id}" "${payload}" >/dev/null
  log_info "Updated existing SPF to authorize ${mechanism}"
}

ensure_dmarc() {
  local zone_id="$1"
  local name="_dmarc.${DOMAIN}"
  local response
  local existing

  response="$(cloudflare_request GET "/zones/${zone_id}/dns_records?type=TXT&name=${name}&per_page=100")"
  existing="$(jq -r '.result[] | select((.content // "") | startswith("v=DMARC1")) | .content' <<<"${response}" | head -n 1)"
  if [[ -n "${existing}" ]]; then
    log_info "Existing DMARC policy preserved: ${existing}"
    return 0
  fi

  upsert_managed_txt "${zone_id}" "${name}" "v=DMARC1; p=none; adkim=r; aspf=r" "v=DMARC1"
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
  local payload

  if ! response="$(backend_request GET "/mail-provider/domains/by-name/${DOMAIN}")"; then
    log_error "Unable to reach the Mail Provider backend while looking up ${DOMAIN}"
    return 1
  fi

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

  payload="$(jq -nc --arg name "${DOMAIN}" '{name:$name}')"
  if ! response="$(backend_request POST "/mail-provider/domains" "${payload}")"; then
    log_error "Unable to reach the Mail Provider backend while registering ${DOMAIN}"
    return 1
  fi

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
    if ! response="$(backend_request POST "/mail-provider/domains/${domain_id}/verify")"; then
      log_warn "Mail Provider verification request could not reach the backend; retrying (${attempt}/36)"
      sleep 5
      continue
    fi

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

check_reverse_dns() {
  if [[ "${PUBLIC_IP}" == *:* ]]; then
    log_warn "Automatic PTR verification currently supports IPv4 only. Configure reverse DNS for ${PUBLIC_IP} -> ${MAIL_HOST} with the IP/VPS provider."
    return 0
  fi

  local reverse_name
  local ptr_response
  local ptr
  IFS='.' read -r ip1 ip2 ip3 ip4 <<<"${PUBLIC_IP}"
  reverse_name="${ip4}.${ip3}.${ip2}.${ip1}.in-addr.arpa"

  ptr_response="$(curl --silent --show-error --fail --max-time 10 \
    --get "https://cloudflare-dns.com/dns-query" \
    --header 'accept: application/dns-json' \
    --data-urlencode "name=${reverse_name}" \
    --data-urlencode 'type=PTR' 2>/dev/null || true)"
  ptr="$(jq -r '.Answer[]?.data // empty' <<<"${ptr_response:-{}}" 2>/dev/null | head -n 1 | sed 's/\.$//')"

  if [[ "${ptr}" == "${MAIL_HOST}" ]]; then
    log_info "PTR/rDNS is aligned: ${PUBLIC_IP} -> ${ptr}"
  elif [[ -n "${ptr}" ]]; then
    log_warn "PTR/rDNS mismatch: ${PUBLIC_IP} -> ${ptr}; ask the VPS/cloud provider to change it to ${MAIL_HOST}."
  else
    log_warn "No PTR/rDNS was found for ${PUBLIC_IP}; ask the VPS/cloud provider to set it to ${MAIL_HOST}. Cloudflare DNS cannot set PTR for an IP owned by another provider."
  fi
}

check_outbound_smtp() {
  if ! command -v nc >/dev/null 2>&1; then
    log_warn "Skipped outbound TCP/25 check because nc is not installed. Haraka requires outbound TCP/25 for direct MX delivery."
    return 0
  fi

  if nc -z -w 5 gmail-smtp-in.l.google.com 25 >/dev/null 2>&1; then
    log_info "Outbound TCP/25 is reachable from this machine"
  else
    log_warn "Outbound TCP/25 appears blocked. Ask the hosting/VPS provider to allow outbound SMTP before production delivery."
  fi
}

require_value CLOUDFLARE_API_TOKEN "${CLOUDFLARE_API_TOKEN:-}"
require_value CLOUDFLARE_ZONE_ID "${CLOUDFLARE_ZONE_ID:-}"

require_command curl
ensure_jq
require_docker
ensure_runtime_dirs

log_info "Checking Cloudflare API token"
cloudflare_request GET "/user/tokens/verify" >/dev/null

discover_zone
discover_public_ip
persist_mail_runtime

log_info "Cloudflare zone: ${DOMAIN} (${CLOUDFLARE_ZONE_ID})"
log_info "GroveHub public mail hostname: ${MAIL_HOST}"

bash "${SCRIPT_DIR}/start.sh"

if ! curl --noproxy "${BACKEND_NO_PROXY}" --fail --silent --max-time 1 "${BACKEND_BASE_URL}/health" >/dev/null 2>&1; then
  log_info "Backend is not running; applying Mail Provider migrations and starting it"
  bash "${SCRIPT_DIR}/../database/migrate.sh" up mailProvider
  BACKEND_HOT_RELOAD=1 bash "${SCRIPT_DIR}/../backend/start.sh"
fi

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

log_info "Publishing outbound-mail DNS records in Cloudflare"
ensure_mail_host_record "${CLOUDFLARE_ZONE_ID}"
upsert_managed_txt "${CLOUDFLARE_ZONE_ID}" "${VERIFICATION_HOST}" "${VERIFICATION_VALUE}" "grovehub-verification="
upsert_managed_txt "${CLOUDFLARE_ZONE_ID}" "${DKIM_HOST}" "${DKIM_VALUE}" "v=DKIM1"

if [[ "${PUBLIC_IP}" == *:* ]]; then
  SPF_MECHANISM="ip6:${PUBLIC_IP}"
else
  SPF_MECHANISM="ip4:${PUBLIC_IP}"
fi
ensure_spf "${CLOUDFLARE_ZONE_ID}" "${SPF_MECHANISM}"
ensure_dmarc "${CLOUDFLARE_ZONE_ID}"

log_info "Reloading Haraka with the domain DKIM key and public hostname"
docker restart "${HARAKA_CONTAINER_NAME}" >/dev/null
if ! wait_for_haraka; then
  docker logs "${HARAKA_CONTAINER_NAME}" >&2 || true
  die "Haraka did not become healthy after loading DKIM"
fi

VERIFIED_RESPONSE="${DOMAIN_RESPONSE}"
if [[ "${DOMAIN_STATUS}" == "VERIFIED" ]]; then
  log_info "Mail Provider domain is already verified; skipping ownership verification wait"
else
  log_info "Waiting for the GroveHub ownership TXT record to become resolvable"
  if ! VERIFIED_RESPONSE="$(verify_domain "${DOMAIN_ID}")"; then
    log_warn "Cloudflare records were created, but DNS verification is still pending"
    log_warn "mail-domain-ready will not fail for DNS propagation; rerun it later to confirm verification"
    VERIFIED_RESPONSE="${DOMAIN_RESPONSE}"
  fi
fi

check_reverse_dns
check_outbound_smtp

if [[ "$(jq -r '.data.status // empty' <<<"${VERIFIED_RESPONSE}")" == "VERIFIED" ]]; then
  log_info "Mail Provider domain is verified: $(jq -r '.data.name' <<<"${VERIFIED_RESPONSE}")"
else
  log_warn "Mail Provider domain DNS is prepared but verification is still pending: ${DOMAIN}"
fi
log_info "Domain ID: ${DOMAIN_ID}"
log_info "Mail hostname: ${MAIL_HOST} -> ${PUBLIC_IP}"
log_info "DKIM: ${DKIM_HOST}"
log_info "Cloudflare outbound-mail DNS is configured"
log_info "No MX record was changed; GroveHub Mail Provider is outbound-only and will not take over inbound mail."
