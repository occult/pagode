#!/usr/bin/env bash
# Query the running instance's SQLite database (side-effect checks).
# Usage: scripts/db.sh "select id, email from users;"
#        scripts/db.sh -json "select * from chat_messages;"
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"
[ -f "$DB_FILE" ] || { echo "verify-pagode: no database at $DB_FILE" >&2; exit 1; }
exec sqlite3 -header -column "$DB_FILE" "$@"
