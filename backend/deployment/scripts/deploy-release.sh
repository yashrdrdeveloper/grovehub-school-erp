#!/usr/bin/env bash
set -Eeuo pipefail
umask 077

# Usage: deploy-release.sh PROJECT COMMIT TRUSTED_SHA256 ARCHIVE [--staging|--prepare-only]
# Rollback: deploy-release.sh PROJECT PREVIOUS_COMMIT TRUSTED_SHA256 ORIGINAL_ARCHIVE [--staging|--prepare-only]
# The digest must come from a separately trusted CI release record.
[[ $# -eq 4 || $# -eq 5 ]] || { echo "Usage: $0 PROJECT COMMIT TRUSTED_SHA256 ARCHIVE [--staging|--prepare-only]" >&2; exit 2; }
[[ "$(id -u)" == 0 ]] || { echo "Root required" >&2; exit 2; }
project="$1"; commit="$2"; digest="$3"; archive="$4"; mode="${5:-}"
[[ "$mode" == "" || "$mode" == "--staging" || "$mode" == "--prepare-only" ]] || exit 2
[[ "$commit" =~ ^[0-9a-f]{40}$ && "$digest" =~ ^[0-9a-f]{64}$ ]] || exit 2
case "$project" in
  platform) port=8080; ready=/ready; live=/health; backend=1 ;;
  school-backend) port=8081; ready=/health/ready; live=/health/live; backend=1 ;;
  internal-tool-backend) port=8082; ready=/health/ready; live=/health/live; backend=1 ;;
  school-frontend) port=8081; ready=/health/ready; live=/health/live; backend=0 ;;
  internal-tool-frontend) port=8082; ready=/health/ready; live=/health/live; backend=0 ;;
  *) echo "Unknown project" >&2; exit 2 ;;
esac
[[ "$(uname -m)" == x86_64 && "$(uname -s)" == Linux ]] || { echo "Requires Linux amd64" >&2; exit 2; }
script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
root="/opt/groovehub/$project"
unit="groovehub-$project.service"
install -d -m 0755 /opt/groovehub "$root" "$root/releases"
# All installation directories must be root-controlled and inaccessible for writes
# by application service users. Do not run against a user-controlled mount.
for dir in /opt/groovehub "$root" "$root/releases"; do
  [[ -d "$dir" && ! -L "$dir" && "$(stat -c %u "$dir")" == 0 ]] || exit 2
  perms="$(stat -c %a "$dir")"
  (( (8#$perms & 022) == 0 )) || { echo "Unsafe release directory" >&2; exit 2; }
done
exec 9>"$root/.deploy.lock"
flock -x 9
work="$(mktemp -d "$root/.deploy.XXXXXX")"
trap 'rm -rf -- "$work"' EXIT
cp -- "$archive" "$work/release.tar.gz"
python3 "$script_dir/verify-release.py" "$work/release.tar.gz" "$work/verified" "$project" "$commit" "$digest"
target="$root/releases/$commit"
# Existing revisions are immutable. Verify every existing byte against the
# freshly verified artifact rather than trusting an old checksum manifest.
if [[ -e "$target" || -L "$target" ]]; then
  [[ -d "$target" && ! -L "$target" ]] || exit 2
  diff -qr --no-dereference "$work/verified" "$target" >/dev/null || { echo "Existing release differs" >&2; exit 2; }
  if ((backend)); then [[ -x "$target/payload/$project" ]] || exit 2; fi
else
  mv -T -- "$work/verified" "$target"
fi
if [[ "$mode" == "--prepare-only" ]]; then
  printf 'Verified release prepared: %s\\n' "$target"
  exit 0
fi
old=""
if [[ -L "$root/current" ]]; then
  old="$(readlink "$root/current")"
  [[ "$old" =~ ^releases/[0-9a-f]{40}$ && -d "$root/$old" && ! -L "$root/$old" ]] || { echo "Invalid current pointer" >&2; exit 2; }
elif [[ -e "$root/current" ]]; then
  echo "Current pointer is not a symlink" >&2; exit 2
fi
was_active=0
if ((backend)); then
  [[ "$(systemctl show -p LoadState --value "$unit")" == loaded ]] || { echo "Install the systemd unit first" >&2; exit 2; }
  if systemctl is-active --quiet "$unit"; then was_active=1; fi
else
  nginx -t
  systemctl is-active --quiet nginx.service
fi
probe() {
  local path="$1" status
  status="$(curl --noproxy '*' --max-time 2 --silent --output /dev/null --write-out '%{http_code}' "http://127.0.0.1:$port$path")" || return 1
  [[ "$status" =~ ^2[0-9]{2}$ ]]
}
healthy() {
  if [[ "$mode" == "--staging" ]]; then
    probe "$live" || return 1
    if [[ "$project" != platform ]]; then probe /health/database || return 1; fi
  else
    probe "$ready" || return 1
  fi
}
if [[ "$mode" == "--staging" ]]; then echo "STAGING ONLY: full readiness is not being asserted" >&2; fi
switch_pointer() {
  ln -s "$1" "$work/pointer"
  mv -Tf -- "$work/pointer" "$root/current"
}
restore() {
  if [[ -n "$old" ]]; then switch_pointer "$old"; else rm -f -- "$root/current"; fi
  if ((backend)); then
    if ((was_active)); then systemctl restart "$unit" || return 1
    else systemctl stop "$unit" || return 1; fi
  fi
}
# A frontend must not be activated against an unhealthy backend.
if ((!backend)); then healthy || { echo "Backend gate failed" >&2; exit 1; }; fi
switch_pointer "releases/$commit"
ok=0
if ((backend)); then
  if systemctl restart "$unit"; then
    for attempt in $(seq 1 30); do
      if systemctl is-active --quiet "$unit" && healthy; then ok=1; break; fi
      sleep 1
    done
  fi
else
  if [[ -f "$root/current/payload/index.html" ]] && healthy; then ok=1; fi
fi
if (( ! ok )); then
  echo "Activation failed; restoring previous pointer and service state" >&2
  restore || { echo "ROLLBACK FAILED: operator intervention required" >&2; exit 3; }
  echo "Previous pointer restored; verify its application health separately" >&2
  exit 1
fi
printf 'Activated %s at %s (%s)\n' "$project" "$commit" "${mode:-full-readiness}"
