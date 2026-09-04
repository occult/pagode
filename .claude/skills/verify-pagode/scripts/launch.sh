#!/usr/bin/env bash
# Build the frontend and Go binaries, then start an isolated Pagode instance.
# Usage: scripts/launch.sh          (port 18000, override with VERIFY_PORT=...)
# Env:   VERIFY_SKIP_FRONTEND=1 skips `npm run build` (only when assets are known fresh).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

if pid_is_ours "$(instance_pid || true)"; then
  echo "verify-pagode: instance already running on port $PORT (pid $(instance_pid)). Run scripts/cleanup.sh first." >&2
  exit 1
fi
if lsof -nP -iTCP:"$PORT" -sTCP:LISTEN >/dev/null 2>&1; then
  echo "verify-pagode: port $PORT is owned by a process this skill did not start. Refusing to double-drive it:" >&2
  lsof -nP -iTCP:"$PORT" -sTCP:LISTEN >&2
  echo "Pick another port with VERIFY_PORT=<port>." >&2
  exit 1
fi

mkdir -p "$RUN" "$EVIDENCE_ROOT"
cd "$REPO"

if [ ! -d node_modules ]; then
  echo "verify-pagode: node_modules missing, running npm ci" >&2
  npm ci >"$RUN/npm-ci.log" 2>&1 || { cat "$RUN/npm-ci.log" >&2; exit 1; }
fi

# assets.go embeds public/build/assets/* at compile time, so the frontend must be
# built before the Go binary or the binary ships stale JS.
if [ "${VERIFY_SKIP_FRONTEND:-0}" != "1" ]; then
  echo "verify-pagode: building frontend (npm run build)" >&2
  npm run build >"$RUN/vite-build.log" 2>&1 || { cat "$RUN/vite-build.log" >&2; exit 1; }
fi
if [ ! -f public/build/.vite/manifest.json ] && [ ! -f public/build/manifest.json ]; then
  echo "verify-pagode: no Vite manifest under public/build; run npm run build" >&2
  exit 1
fi
if [ -f public/hot ]; then
  echo "verify-pagode: WARNING public/hot exists, so pages will load JS from the Vite dev server named in that file instead of the embedded build. Stop 'npx vite' (which removes public/hot) for a self-contained run." >&2
fi

echo "verify-pagode: building Go binaries" >&2
go build -o "$WEB_BIN" ./cmd/web
go build -o "$ADMIN_BIN" ./cmd/admin

# cwd must stay at the repo root: root.html, the Vite manifest, and static/ are
# resolved relative to the directory holding go.mod.
# Plain env-prefixed command (not the app_env function): backgrounding a function
# forks a subshell and $! would be that subshell, not the server.
PAGODE_HTTP_PORT="$PORT" \
PAGODE_APP_HOST="$BASE_URL" \
PAGODE_DATABASE_CONNECTION="$DB_FILE?_journal=WAL&_timeout=5000&_fk=true" \
PAGODE_FILES_DIRECTORY="$UPLOADS_DIR" \
nohup "$WEB_BIN" >"$SERVER_LOG" 2>&1 &
echo $! >"$PID_FILE"

for _ in $(seq 1 60); do
  if [ "$(http_code "$BASE_URL/user/login")" = "200" ]; then
    break
  fi
  if ! kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
    echo "verify-pagode: server exited during startup:" >&2
    cat "$SERVER_LOG" >&2
    rm -f "$PID_FILE"
    exit 1
  fi
  sleep 0.5
done
if [ "$(http_code "$BASE_URL/user/login")" != "200" ]; then
  echo "verify-pagode: server did not answer GET /user/login with 200 within 30s:" >&2
  cat "$SERVER_LOG" >&2
  exit 1
fi

cat >"$RUN/env" <<ENV
VERIFY_PORT=$PORT
BASE_URL=$BASE_URL
RUN_DIR=$RUN
DB_FILE=$DB_FILE
UPLOADS_DIR=$UPLOADS_DIR
SERVER_LOG=$SERVER_LOG
EVIDENCE_ROOT=$EVIDENCE_ROOT
PID=$(cat "$PID_FILE")
ENV

echo "verify-pagode: ready"
cat "$RUN/env"
