#!/usr/bin/env bash
# Read-only health check: is the instance on $PORT ours and worth driving?
# Usage: scripts/doctor.sh   (exit 0 = drive it, exit 1 = fix or relaunch first)
set -uo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

fail=0
ok()   { echo "OK   $*"; }
bad()  { echo "FAIL $*"; fail=1; }
warn() { echo "WARN $*"; }

pid="$(instance_pid || true)"
if [ -z "$pid" ]; then
  bad "no pid file at $PID_FILE (run scripts/launch.sh)"
elif ! kill -0 "$pid" 2>/dev/null; then
  bad "pid $pid from $PID_FILE is not running (run scripts/cleanup.sh then scripts/launch.sh)"
elif ! pid_is_ours "$pid"; then
  bad "pid $pid is alive but is not $WEB_BIN: $(ps -p "$pid" -o args=)"
else
  ok "process $pid is $WEB_BIN"
fi

owner="$(lsof -nP -iTCP:"$PORT" -sTCP:LISTEN -t 2>/dev/null | head -1)"
if [ -z "$owner" ]; then
  bad "nothing is listening on port $PORT"
elif [ "$owner" != "$pid" ]; then
  bad "port $PORT is owned by pid $owner, not our pid ${pid:-?}"
else
  ok "port $PORT is owned by pid $pid"
fi

code="$(http_code "$BASE_URL/user/login")"
if [ "$code" = "200" ]; then
  if curl -s --max-time 5 "$BASE_URL/user/login" | grep -q 'data-page='; then
    ok "GET $BASE_URL/user/login -> 200 with an Inertia page"
  else
    bad "GET $BASE_URL/user/login -> 200 but no Inertia data-page attribute in the HTML"
  fi
else
  bad "GET $BASE_URL/user/login -> $code"
fi

if [ -f "$DB_FILE" ]; then
  users="$(sqlite3 "$DB_FILE" 'select count(*) from users;' 2>&1)"
  if [ $? -eq 0 ]; then ok "database $DB_FILE has $users user(s)"; else bad "sqlite3 cannot read $DB_FILE: $users"; fi
else
  bad "database file $DB_FILE missing"
fi

[ -f "$CREDS_FILE" ] && ok "admin credentials seeded at $CREDS_FILE" || warn "no seeded admin (scripts/seed-admin.sh) - only needed for admin-panel proofs"

if [ -f "$REPO/public/hot" ]; then
  hot="$(cat "$REPO/public/hot")"
  if [ "$(http_code "$hot/@vite/client")" = "200" ]; then
    warn "public/hot points at $hot (answering): pages load JS from the dev server, not the embedded build"
  else
    bad "public/hot points at $hot but it does not answer; pages will render without JS. Delete public/hot or start npx vite"
  fi
else
  ok "no public/hot: pages use the embedded Vite build"
fi

if node "$VERIFY_SKILL_DIR/scripts/pw.mjs" --doctor; then
  ok "playwright + chromium usable"
else
  bad "playwright/chromium not usable (see message above)"
fi

exit $fail
