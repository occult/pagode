# Sourced by every verify-pagode script. Not executable on its own.
# shellcheck shell=bash

VERIFY_SKILL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO="$(cd "$VERIFY_SKILL_DIR/../../.." && pwd)"

PORT="${VERIFY_PORT:-18000}"
BASE_URL="http://localhost:$PORT"

STATE_ROOT="$REPO/tmp/verify-pagode"          # tmp/ is gitignored
RUN="$STATE_ROOT/run/$PORT"                    # per-instance scratch, removed by cleanup.sh
EVIDENCE_ROOT="$STATE_ROOT/evidence"           # proof artifacts, never removed by cleanup.sh

DB_FILE="$RUN/verify.db"
UPLOADS_DIR="$RUN/uploads"
WEB_BIN="$RUN/pagode-web"
ADMIN_BIN="$RUN/pagode-admin"
PID_FILE="$RUN/pid"
SERVER_LOG="$RUN/server.log"
CREDS_FILE="$RUN/admin-credentials.env"

# Run a command with the environment that isolates this instance from dbs/main.db,
# uploads/, and port 8000. Viper maps PAGODE_A_B to config key a.b.
app_env() {
  env PAGODE_HTTP_PORT="$PORT" \
      PAGODE_APP_HOST="$BASE_URL" \
      PAGODE_DATABASE_CONNECTION="$DB_FILE?_journal=WAL&_timeout=5000&_fk=true" \
      PAGODE_FILES_DIRECTORY="$UPLOADS_DIR" \
      "$@"
}

instance_pid() {
  [ -f "$PID_FILE" ] && cat "$PID_FILE"
}

# True only when the pid is alive and its argv[0] is the binary launch.sh built.
pid_is_ours() {
  local pid="$1"
  [ -n "$pid" ] || return 1
  local argv0
  argv0="$(ps -p "$pid" -o args= 2>/dev/null | awk '{print $1}')"
  [ "$argv0" = "$WEB_BIN" ]
}

http_code() {
  curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$1" || true
}
