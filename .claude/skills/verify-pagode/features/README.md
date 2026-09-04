# Pagode verification map

This directory is the maintained source for verifying the user-facing behavior of Pagode.
Read this index before driving the app, then use the matching feature file as the recipe.

## Baseline preconditions

- Launch with `scripts/launch.sh` from `.claude/skills/verify-pagode/`. The instance answers at
  `http://localhost:18000` with its own database and uploads directory under `tmp/verify-pagode/run/18000/`.
- `scripts/doctor.sh` exits 0.
- The database starts empty. Features that need a session register a fresh user through the UI
  (`scenarios/lib/auth.mjs`); admin-panel proofs seed one admin with `scripts/seed-admin.sh`.
- Never drive an instance this run did not launch.

## Driving conventions

- Run browser actions through scenarios executed by `node scripts/pw.mjs <scenario>`.
- Prefer ids and ARIA roles named in the feature files over CSS classes or positions.
- Assert flash messages by their exact server text; they show as toasts at the top center.
- Use `unique()` for every name and email; `users.email` and `chat_rooms.name` are unique columns.
- Check side effects with `db()` in a scenario or `scripts/db.sh` by hand.
- Do not remove proof artifacts during cleanup.

## Proof and skip reporting

- Capture the user action and the resulting state, not only the final screen.
- UI proof is the step screenshots plus `trace.zip`; database proof is the `db()` line in `run.log`.
- Record the feature ID and entry point used with every artifact.
- Report an unreachable path with the attempted command and the unmet precondition.
- Do not report a skipped entry point as verified through a different path.

## Feature entry contract

Each feature file starts with an H1 title and one paragraph describing the user-visible behavior.
It then uses exactly four H2 sections in this order.

1. `Sub-features` lists short IDs with one line for each behavior.
2. `How to get to it (user POV)` lists every user entry point.
3. `Driving it with Playwright` starts with `Preconditions:` and uses labeled bullets that pair each user action with an exact handle and observable result.
4. `Gotchas` lists traps that can waste or invalidate a verification run.

## Features

- [Authentication](./auth.md) covers register, log in, log out, forgot password, and the 401 wall on protected pages. Scenario: `scenarios/register-login.mjs`.
- [Community chat](./chat.md) covers rooms, the WebSocket conversation, guest nicknames, password rooms, and attachments. Scenario: `scenarios/chat-room.mjs`.
- [Profile settings](./profile-settings.md) covers name/email update, password change, appearance, and account deletion.
- [File upload](./file-upload.md) covers uploading through `/files` and the on-disk result.
- [Admin panel](./admin-panel.md) covers the user list, add/edit/delete, the per-entity CRUD pages, and the task monitor.

Unmapped, pending Stripe test keys in `config/config.yaml`: `/plans`, `/products`, `/premium`, `/billing`.
