#!/bin/sh
set -eu

CONFIG_DIR="/opt/haraka/config"
PLUGINS_FILE="${CONFIG_DIR}/plugins"
SMTP_FORWARD_CONFIG="${CONFIG_DIR}/smtp_forward.ini"

if [ -n "${HARAKA_HOSTNAME:-}" ]; then
  printf '%s\n' "${HARAKA_HOSTNAME}" > "${CONFIG_DIR}/me"
fi

cat > "${PLUGINS_FILE}" <<'EOF'
relay
dkim
EOF

relay_host="${MAIL_PROVIDER_RELAY_HOST:-}"

if [ -n "${relay_host}" ]; then
  relay_port="${MAIL_PROVIDER_RELAY_PORT:-587}"
  relay_username="${MAIL_PROVIDER_RELAY_USERNAME:-}"
  relay_password="${MAIL_PROVIDER_RELAY_PASSWORD:-}"
  relay_tls="${MAIL_PROVIDER_RELAY_TLS:-true}"

  case "${relay_port}" in
    ''|*[!0-9]*)
      printf '%s\n' "MAIL_PROVIDER_RELAY_PORT must be a numeric TCP port" >&2
      exit 1
      ;;
  esac

  if [ "${relay_port}" -lt 1 ] || [ "${relay_port}" -gt 65535 ]; then
    printf '%s\n' "MAIL_PROVIDER_RELAY_PORT must be between 1 and 65535" >&2
    exit 1
  fi

  case "${relay_tls}" in
    true|TRUE|1|yes|YES)
      relay_tls="true"
      ;;
    false|FALSE|0|no|NO)
      relay_tls="false"
      ;;
    *)
      printf '%s\n' "MAIL_PROVIDER_RELAY_TLS must be true or false" >&2
      exit 1
      ;;
  esac

  if { [ -n "${relay_username}" ] && [ -z "${relay_password}" ]; } || \
     { [ -z "${relay_username}" ] && [ -n "${relay_password}" ]; }; then
    printf '%s\n' "MAIL_PROVIDER_RELAY_USERNAME and MAIL_PROVIDER_RELAY_PASSWORD must be configured together" >&2
    exit 1
  fi

  {
    printf 'host=%s\n' "${relay_host}"
    printf 'port=%s\n' "${relay_port}"
    printf 'enable_outbound=true\n'
    printf 'enable_tls=%s\n' "${relay_tls}"
    printf 'connect_timeout=15\n'

    if [ -n "${relay_username}" ]; then
      printf 'auth_type=login\n'
      printf 'auth_user=%s\n' "${relay_username}"
      printf 'auth_pass=%s\n' "${relay_password}"
    fi
  } > "${SMTP_FORWARD_CONFIG}"

  printf '%s\n' 'queue/smtp_forward' >> "${PLUGINS_FILE}"

  if [ -n "${relay_username}" ]; then
    relay_auth="enabled"
  else
    relay_auth="disabled"
  fi

  printf '%s\n' "Haraka outbound relay enabled: ${relay_host}:${relay_port} TLS=${relay_tls} AUTH=${relay_auth}"
else
  rm -f "${SMTP_FORWARD_CONFIG}"
  printf '%s\n' "Haraka outbound relay disabled; using direct MX delivery"
fi

exec haraka -c /opt/haraka
