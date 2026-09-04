---
name: verify-pagode
description: Launch an isolated Pagode instance (Go + Echo + Inertia/React web app) and drive it through a real headless Chromium the way a user would, capturing screenshots, traces, and SQLite side effects as proof. Use when a change must be shown working in the running app (register, login, chat, profile, uploads, admin panel), not just in `make test`.
---

# verify-pagode

Pagode's user surface is the web UI on port 8000: Inertia pages rendered by Go handlers under
`pkg/handlers/` and React pages under `resources/js/Pages/`. Secondary surfaces are the WebSocket
chat (`/ws/chat/:id`) and the admin CLI (`cmd/admin`, used here only to seed an admin user).
This skill never touches the developer's own instance: it builds its own binary, runs it on
port 18000 with its own SQLite file and uploads directory, and drives that.

All paths below are relative to this skill directory (`.claude/skills/verify-pagode/`).
Every script accepts `VERIFY_PORT=<port>` to run a second instance side by side; the default is 18000.

## Launch

```bash
scripts/launch.sh
```

What it does, in order: `npm ci` if `node_modules/` is missing, `npm run build` (the Go binary
embeds `public/build/assets/*` at compile time, so a stale build ships stale JS), `go build` of
`cmd/web` and `cmd/admin` into the run directory, then starts the server with:

| Setting | Value |
|---|---|
| Port | `PAGODE_HTTP_PORT=18000` |
| Database | `tmp/verify-pagode/run/18000/verify.db` (WAL, foreign keys on) |
| Uploads | `tmp/verify-pagode/run/18000/uploads/` |
| Log | `tmp/verify-pagode/run/18000/server.log` |

Ready means `GET http://localhost:18000/user/login` answers 200, which the script polls for up to
30 s. It prints the instance facts and writes them to `tmp/verify-pagode/run/18000/env`.
`tmp/` is gitignored. Set `VERIFY_SKIP_FRONTEND=1` to skip the Vite build when you know the
assets are fresh.

Launch refuses to start if the port is already listening and the owner is not a process this
skill started. Never drive a Pagode on a port you did not launch: that is the developer's session
and its `dbs/main.db`.

Seed an admin (only needed for admin-panel proofs; everything else registers through the UI):

```bash
scripts/seed-admin.sh              # verify-admin@example.com, prints and stores the password
scripts/seed-admin.sh me@example.com
```

Credentials land in `tmp/verify-pagode/run/18000/admin-credentials.env` (`ADMIN_EMAIL`,
`ADMIN_PASSWORD`). The user is created verified and admin via the repo's own `cmd/admin`.

## Doctor

```bash
scripts/doctor.sh      # read-only; exit 0 = drive it
```

Checks: the pid file names a live process whose argv[0] is our built binary; that pid owns port
18000; `/user/login` returns 200 with an Inertia `data-page` attribute; the SQLite file is readable;
whether an admin is seeded; whether `public/hot` exists (if it does, pages load JS from the Vite dev
server named in that file instead of the embedded build); Playwright resolves and Chromium launches.
Run it first whenever anything looks off, and after launch before the first drive.

## Drive

Harness: Playwright (Node library) driving headless Chromium, wrapped by `scripts/pw.mjs`.
Playwright is resolved from `<repo>/node_modules`, then `~/.npm/_npx/*/node_modules`, then the
global npm root. If none has it:

```bash
npm install --no-save playwright && npx playwright install chromium   # leaves package.json untouched
```

Run a scenario:

```bash
node scripts/pw.mjs scenarios/register-login.mjs           # auth feature
node scripts/pw.mjs scenarios/chat-room.mjs                # chat feature
node scripts/pw.mjs scenarios/<mine>.mjs --name my-label --headed   # watch it
```

A scenario is an ES module whose default export receives:

| Arg | Use |
|---|---|
| `page`, `context` | Playwright page/context, `baseURL` preset to the instance |
| `step(name, fn)` | Runs `fn`, screenshots after it (or `-FAILED.png` on throw), logs timing |
| `db(sql)` | Runs `sqlite3 -json` against the instance DB, returns rows |
| `expect(cond, msg)` | Throws on false; use for every assertion |
| `unique(prefix)` | Collision-free names/emails so reruns never hit the unique constraints |
| `log(line)` | Appends to `run.log` and stdout |

Stable handles in this app (from the React sources; prefer these over CSS or coordinates):

- Auth forms use ids: `#name`, `#email`, `#password`, `#password-confirm` (register),
  `#email` / `#password` (login). Buttons by role: `Create Account`, `Log in`.
- Page identity: `getByRole("heading", { name: "Create an account" })`,
  `"Log in to your account"`, `"Community Chat"`. Do not use `getByText` for these; the auth layout
  repeats the title in a screen-reader-only link.
- Flash messages render as sonner toasts; assert the exact server text with `page.getByText(...)`.
  The texts live in `msg.Success/Danger/...` calls inside `pkg/handlers/*.go`.
- Sidebar user menu: `getByRole("button", { name: /<user name>/ })` then
  `getByRole("menuitem", { name: "Log out" })`.
- Chat: button `Create Room`, dialog `Create Chat Room` with `#room-name` and `#room-pw`, submit
  button `Create`; inside a room the composer is `getByPlaceholder("Type a message...")` and Enter sends.
- Protected routes answer 401 (an Inertia error page), not a redirect to login, when logged out.

`scenarios/lib/auth.mjs` exports `registerViaUI` and `loginViaUI` for scenarios that need a session;
both go through the real forms.

The feature map in `features/` is the source of what to drive and what proves it. A proof that
takes one convenient entry point is incomplete when the map lists others.

## Evidence

Each run writes `tmp/verify-pagode/evidence/<UTC timestamp>-<label>/`:

- `NN-<step>.png` after every step; `NN-<step>-FAILED.png` where it threw.
- `trace.zip`, open with `npx playwright show-trace <path>` for DOM snapshots and network.
- `run.log`: steps, every `db()` query and result, browser console errors, HTTP >= 400 responses.
- `summary.json`: status, steps with timings, scenario path, base URL, Playwright version.

Proof standards for this app:

- Exercise the real user path: forms, buttons, menus. Never call an internal setter, seed rows
  for the thing under test, or hit a handler with curl and call it verified.
- Capture the action and the resulting state, not just the final screen. `step()` does this when
  each step is one user action.
- Verify side effects with `db()` (tables are `users`, `chat_rooms`, `chat_messages`, `chat_bans`,
  `password_tokens`, payment tables) and on disk (`tmp/verify-pagode/run/18000/uploads/` for
  `/files`, `static/chat-uploads/` for chat attachments), alongside what is visible.
- Mocks only at existing production boundaries: mail is a skeleton client that logs instead of
  sending, and Stripe needs real test keys in `config/config.yaml`. Payment features are therefore
  unmapped until keys exist; do not stub them in.
- Ad hoc checks outside a scenario: `scripts/db.sh "select id, email, admin from users;"` and
  `curl -s -o /dev/null -w '%{http_code}' http://localhost:18000/<path>`. Paste the command and
  output into the report; they are not captured automatically.

## Cleanup

```bash
scripts/cleanup.sh
```

Kills only the pid recorded by launch, and only after confirming its argv[0] is our built
binary; then copies the server log to `tmp/verify-pagode/evidence/server-18000-<timestamp>.log`
and deletes the run directory (binaries, DB, uploads, credentials). Evidence directories are
never removed. Run cleanup after every failed attempt too, so ports and processes are not stranded.

Cleanup does not touch `static/chat-uploads/` (chat attachments are written there by any
instance, including the developer's) or `node_modules/` and `public/build/` (shared build outputs).

## Helpers

| Script | Invocation | Purpose |
|---|---|---|
| `scripts/launch.sh` | `scripts/launch.sh` | Build + start the isolated instance, wait for ready |
| `scripts/doctor.sh` | `scripts/doctor.sh` | Read-only health check, exit 1 on any FAIL |
| `scripts/seed-admin.sh` | `scripts/seed-admin.sh [email]` | Create a verified admin via `cmd/admin` |
| `scripts/db.sh` | `scripts/db.sh "<sql>"` | Query the instance DB (`-json` accepted) |
| `scripts/pw.mjs` | `node scripts/pw.mjs <scenario> [--name L] [--headed]` / `--doctor` | Run a scenario, collect evidence |
| `scripts/cleanup.sh` | `scripts/cleanup.sh` | Stop our instance, keep evidence |
| `scripts/common.sh` | sourced by the others | Paths, ports, env, pid ownership check |

Full loop:

```bash
cd .claude/skills/verify-pagode
scripts/launch.sh && scripts/doctor.sh
node scripts/pw.mjs scenarios/register-login.mjs
scripts/cleanup.sh
ls ../../../tmp/verify-pagode/evidence/        # tmp/verify-pagode/evidence/ from the repo root
```

Keep the map honest as the app changes with `/maintain-verification-skill`.
