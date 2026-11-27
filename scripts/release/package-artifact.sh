#!/usr/bin/env bash
set -euo pipefail

# Package only explicitly selected build outputs; never archive the checkout.
# This script runs on Linux CI, not on the target host.
if [[ $# -ne 4 ]]; then
  echo "Usage: package-artifact.sh PROJECT SOURCE DESTINATION COMMIT" >&2
  exit 2
fi
project="$1"
source_path="$2"
destination="$3"
revision="$4"
case "$project" in
  platform|school-backend|internal-tool-backend|school-frontend|internal-tool-frontend) ;;
  *) echo "Unknown project" >&2; exit 2 ;;
esac
[[ "$revision" =~ ^[0-9a-f]{40}$ ]] || { echo "Invalid commit SHA" >&2; exit 2; }
[[ -e "$source_path" && ! -L "$source_path" ]] || { echo "Missing build output" >&2; exit 2; }
mkdir -p "$destination"
stage="$(mktemp -d)"
trap 'rm -rf -- "$stage"' EXIT
mkdir -p "$stage/payload"
case "$project" in
  *-frontend)
    [[ -d "$source_path" && -f "$source_path/index.html" ]] || exit 2
    # Reject symlinks and local secrets rather than copying arbitrary targets.
    if [[ -n "$(find "$source_path" \( -type l -o -name '.env*' -o -name '*.pem' -o -name '*.key' \) -print -quit)" ]]; then
      echo "Unsafe frontend output" >&2; exit 2
    fi
    cp -R "$source_path/." "$stage/payload/"
    ;;
  *)
    [[ -f "$source_path" && -x "$source_path" ]] || exit 2
    install -m 0755 "$source_path" "$stage/payload/$project"
    ;;
esac
printf '{"project":"%s","commit":"%s","platform":"linux-amd64"}\n' "$project" "$revision" > "$stage/manifest.json"
(
  cd "$stage"
  find payload -type f -print0 | LC_ALL=C sort -z | xargs -0 sha256sum > SHA256SUMS
  sha256sum manifest.json >> SHA256SUMS
  sha256sum -c SHA256SUMS
)
tar -C "$stage" -czf "$destination/$project.tar.gz" manifest.json SHA256SUMS payload
(
  cd "$destination"
  sha256sum "$project.tar.gz" > "$project.tar.gz.sha256"
)
