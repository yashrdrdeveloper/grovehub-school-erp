#!/bin/sh
set -eu

if [ -n "${HARAKA_HOSTNAME:-}" ]; then
  printf '%s\n' "${HARAKA_HOSTNAME}" > /opt/haraka/config/me
fi

exec haraka -c /opt/haraka
