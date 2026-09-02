#!/usr/bin/env bash
# Create a verified admin user in the running instance's database via the repo's own admin CLI.
# Usage: scripts/seed-admin.sh [email]   (default verify-admin@example.com)
# Writes ADMIN_EMAIL / ADMIN_PASSWORD to $RUN/admin-credentials.env and prints them.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

EMAIL="${1:-verify-admin@example.com}"
[ -x "$ADMIN_BIN" ] || { echo "verify-pagode: $ADMIN_BIN missing, run scripts/launch.sh first" >&2; exit 1; }

cd "$REPO"
out="$(app_env "$ADMIN_BIN" --email="$EMAIL")"
pw="$(printf '%s\n' "$out" | awk -F': ' '/^Password:/{print $2}')"
if [ -z "$pw" ]; then
  echo "verify-pagode: admin CLI did not print a password:" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
printf 'ADMIN_EMAIL=%s\nADMIN_PASSWORD=%s\n' "$EMAIL" "$pw" >"$CREDS_FILE"
chmod 600 "$CREDS_FILE"
cat "$CREDS_FILE"
