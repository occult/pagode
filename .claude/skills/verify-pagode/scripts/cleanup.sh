#!/usr/bin/env bash
# Stop the instance launch.sh started and remove its scratch state.
# Keeps every proof under tmp/verify-pagode/evidence/ and copies the server log there.
# Usage: scripts/cleanup.sh   (VERIFY_PORT selects the instance)
set -uo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

pid="$(instance_pid || true)"
if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
  if pid_is_ours "$pid"; then
    kill "$pid"
    for _ in $(seq 1 20); do
      kill -0 "$pid" 2>/dev/null || break
      sleep 0.5
    done
    if kill -0 "$pid" 2>/dev/null; then
      echo "verify-pagode: pid $pid ignored SIGTERM, sending SIGKILL" >&2
      kill -9 "$pid"
    fi
    echo "verify-pagode: stopped pid $pid"
  else
    echo "verify-pagode: pid $pid is not $WEB_BIN, leaving it alone: $(ps -p "$pid" -o args=)" >&2
  fi
fi

if [ -f "$SERVER_LOG" ]; then
  mkdir -p "$EVIDENCE_ROOT"
  dest="$EVIDENCE_ROOT/server-$PORT-$(date +%Y%m%d-%H%M%S).log"
  cp "$SERVER_LOG" "$dest"
  echo "verify-pagode: server log kept at $dest"
fi

if [ -d "$RUN" ]; then
  rm -rf "$RUN"
  echo "verify-pagode: removed $RUN"
fi
echo "verify-pagode: evidence preserved under $EVIDENCE_ROOT"
