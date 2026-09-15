# Modular Setup CLI — Test Plan

**Date:** 2026-03-27
**Related:** [Design Document](2026-03-27-modular-setup-cli-design.md)
**Goal:** Verify that the setup CLI correctly removes optional features — individually, in combination, and in edge cases — without breaking the Go or frontend build.

---

## Prerequisites (run once)

```bash
# 1. Commit all current changes so git is clean
git add -A && git commit -m "feat: add modular setup CLI with feature markers"

# 2. Install npm dependencies (needed for frontend verification)
npm install
```

---

## Test 1: No changes (keep all features)

**Goal:** Selecting all features exits cleanly with no modifications.

```bash
# Run setup, keep all checkboxes checked, submit
go run cmd/setup/main.go
```

**Expected:** Prints `No changes needed.` and exits. No files modified.

**Verify:**
```bash
git status   # should show nothing changed
```

---

## Test 2: Remove single feature — Chat

**Goal:** Removing only Chat deletes all chat files, patches shared files, rebuilds Ent.

```bash
go run cmd/setup/main.go
# Uncheck ONLY "Chat (WebSocket)", submit
```

**Verify after completion:**
```bash
# Files deleted
test ! -f pkg/handlers/chat.go
test ! -f resources/js/hooks/useChat.ts
test ! -f resources/js/types/chat.d.ts
test ! -f resources/js/hooks/useAudioRecorder.ts
test ! -f ent/schema/chatroom.go
test ! -f ent/schema/chatmessage.go
test ! -f ent/schema/chatban.go

# Markers removed from shared files
grep -c "feature:chat" pkg/services/container.go                # should be 0
grep -c "feature:chat" pkg/handlers/router.go                   # should be 0
grep -c "feature:chat" resources/js/components/AppSidebar.tsx   # should be 0
grep -c "feature:chat" config/config.yaml                       # should be 0
grep -c "feature:chat" Makefile                                 # should be 0
grep -c "feature:chat" ent/schema/user.go                       # should be 0

# Chat nav item gone from sidebar
grep "Chat" resources/js/components/AppSidebar.tsx       # should not find "Chat" nav item

# Chat config section gone
grep "chat:" config/config.yaml                          # should not find chat section

# chat-clear target gone from Makefile
grep "chat-clear" Makefile                               # should not find it

# Build passes
go build -o /dev/null ./cmd/web

# Other features still work
grep "feature:payment" pkg/services/container.go         # should still exist
grep "feature:mail" pkg/handlers/auth.go                 # should still exist
```

**Recover:** `git checkout .`

---

## Test 3: Remove single feature — Payment

**Goal:** Removing Payment deletes Stripe files, patches middleware, rebuilds Ent.

```bash
go run cmd/setup/main.go
# Uncheck ONLY "Payment (Stripe)", submit
```

**Verify:**
```bash
# Files deleted
test ! -f pkg/handlers/billing.go
test ! -f pkg/handlers/plans.go
test ! -f pkg/services/payment.go
test ! -f pkg/services/payment_stripe.go
test ! -f resources/js/Pages/Billing.tsx
test ! -f resources/js/Pages/Plans.tsx
test ! -f resources/js/components/PaymentForm.tsx
test ! -f ent/schema/subscription.go
test ! -f ent/schema/paymentcustomer.go

# RequirePaidUser middleware removed
grep "RequirePaidUser" pkg/middleware/auth.go             # should not find it

# Payment imports removed from middleware
grep "paymentcustomer" pkg/middleware/auth.go             # should not find it

# Payment config gone
grep "payment:" config/config.yaml                       # should not find it

# Build passes
go build -o /dev/null ./cmd/web
```

**Recover:** `git checkout .`

---

## Test 4: Remove single feature — Mail

**Goal:** Removing Mail deletes mail service, patches auth handler (email verification + forgot password).

```bash
go run cmd/setup/main.go
# Uncheck ONLY "Mail (Resend)", submit
```

**Verify:**
```bash
# Files deleted
test ! -f pkg/handlers/contact.go
test ! -f pkg/services/mail.go
test ! -f pkg/services/mail_test.go
test ! -f pkg/ui/emails/auth.go

# Auth handler: mail field removed
grep "mail.*MailClient" pkg/handlers/auth.go             # should not find it

# Auth handler: sendVerificationEmail function removed
grep "sendVerificationEmail" pkg/handlers/auth.go        # should not find it

# Auth handler: ForgotPasswordSubmit still exists but without mail sending
grep "ForgotPasswordSubmit" pkg/handlers/auth.go         # should find the function
grep "h.mail" pkg/handlers/auth.go                       # should NOT find it

# Contact route names removed
grep "Contact" pkg/routenames/names.go                   # should not find Contact/ContactSubmit

# Build passes
go build -o /dev/null ./cmd/web
```

**Recover:** `git checkout .`

---

## Test 5: Remove single feature — Background Tasks

**Goal:** Removing Tasks patches admin handler and container.

```bash
go run cmd/setup/main.go
# Uncheck ONLY "Background Tasks", submit
```

**Verify:**
```bash
# Files deleted
test ! -f pkg/handlers/task.go
test ! -d pkg/tasks

# Admin handler: backlite references removed
grep "backlite" pkg/handlers/admin.go                    # should not find it
grep "AdminTasks" pkg/routenames/names.go                # should not find it

# Container: tasks shutdown removed
grep "initTasks" pkg/services/container.go               # should not find it

# Build passes
go build -o /dev/null ./cmd/web
```

**Recover:** `git checkout .`

---

## Test 6: Remove single feature — File Upload

**Goal:** Removing Files patches container and sidebar.

```bash
go run cmd/setup/main.go
# Uncheck ONLY "File Upload", submit
```

**Verify:**
```bash
# Files deleted
test ! -f pkg/handlers/files.go
test ! -f resources/js/Pages/UploadFile.tsx

# Sidebar: Upload Files nav item gone
grep "Upload Files" resources/js/components/AppSidebar.tsx   # should not find it
grep "UploadCloud" resources/js/components/AppSidebar.tsx    # should not find it

# Container: initFiles removed
grep "initFiles" pkg/services/container.go               # should not find it
grep "afero" pkg/services/container.go                   # should not find it

# Build passes
go build -o /dev/null ./cmd/web
```

**Recover:** `git checkout .`

---

## Test 7: Remove multiple features — Chat + Payment

**Goal:** Two features with Ent schemas removed simultaneously.

```bash
go run cmd/setup/main.go
# Uncheck "Chat (WebSocket)" AND "Payment (Stripe)", submit
```

**Verify:**
```bash
# All chat + payment files gone
test ! -f pkg/handlers/chat.go
test ! -f pkg/handlers/billing.go
test ! -f ent/schema/chatroom.go
test ! -f ent/schema/subscription.go

# User schema: both edge groups removed
grep "chat_rooms\|payment_customer" ent/schema/user.go   # should not find either

# Ent regenerated successfully (no stale files)
ls ent/*.go | head -5                                    # should exist (regenerated)
test ! -d ent/chatroom                                   # chat entity dir gone
test ! -d ent/subscription                               # subscription entity dir gone

# Build passes
go build -o /dev/null ./cmd/web
```

**Recover:** `git checkout .`

---

## Test 8: Remove ALL features

**Goal:** Stripping everything leaves a core auth + dashboard skeleton.

```bash
go run cmd/setup/main.go
# Uncheck ALL 5 features, submit
```

**Verify:**
```bash
# Core files still exist
test -f pkg/handlers/auth.go
test -f pkg/handlers/router.go
test -f pkg/services/container.go
test -f ent/schema/user.go

# No feature markers remain in any file
grep -r "feature:" --include='*.go' --include='*.tsx' --include='*.yaml' . \
  | grep -v cmd/setup | grep -v node_modules
# Should return nothing

# Config only has core sections
cat config/config.yaml | grep -E "^[a-z]+:"
# Should only show: http, app, cache, database (no files, tasks, mail, payment, chat)

# Sidebar only has Dashboard
grep "title:" resources/js/components/AppSidebar.tsx
# Should show Dashboard, Repository, Documentation, Admin Panel only

# Build passes
go build -o /dev/null ./cmd/web

# Setup tool removed itself
test ! -d cmd/setup
grep "setup" Makefile                                    # should not find make setup target
```

**Recover:** `git checkout .`

---

## Test 9: Edge case — Dirty git working tree

**Goal:** Tool refuses to run with uncommitted changes.

```bash
echo "dirty" >> README.md
go run cmd/setup/main.go
```

**Expected:** `Error: Please commit or stash changes before running setup.`

**Recover:** `git checkout .`

---

## Test 10: Edge case — Idempotent run

**Goal:** Running the tool twice with the same selection doesn't fail.

```bash
# First run: remove Chat
go run cmd/setup/main.go   # uncheck Chat, submit

# Commit the result
git add -A && git commit -m "test: removed chat"

# Second run: remove Chat again (already gone)
go run cmd/setup/main.go   # uncheck Chat again, submit
```

**Expected:** Second run completes successfully (skips already-deleted files, no-ops on missing markers). Build still passes.

**Recover:** `git checkout .` then `git reset HEAD~1`

---

## Recommended Test Order

| Priority | Test | What it validates |
|----------|------|-------------------|
| P0 | Test 1 | No-op path works |
| P0 | Test 2 | Single feature (Chat) — Ent regen, WebSocket removal |
| P0 | Test 4 | Single feature (Mail) — Auth handler graceful degradation |
| P0 | Test 8 | All features removed — nuclear option |
| P1 | Test 3 | Single feature (Payment) — middleware removal |
| P1 | Test 5 | Single feature (Tasks) — admin handler patching |
| P1 | Test 6 | Single feature (Files) — container + sidebar |
| P1 | Test 7 | Multi-feature — Ent regen with multiple schema removals |
| P2 | Test 9 | Dirty git guard |
| P2 | Test 10 | Idempotency |
