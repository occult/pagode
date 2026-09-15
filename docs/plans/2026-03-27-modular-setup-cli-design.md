# Modular Setup CLI — Design Document

**Date:** 2026-03-27
**Status:** Approved
**Goal:** Allow developers to opt-out of features they don't need when bootstrapping a Pagode project, via a one-time interactive CLI tool that physically removes unused code.

---

## Overview

Pagode ships with 5 optional features: Payment (Stripe), Chat (WebSocket), Mail (Resend), Background Tasks (Backlite), and File Upload. Today, removing any of them is a manual, error-prone process. This design introduces a setup CLI (`cmd/setup/main.go`) that presents a TUI multi-select checklist and surgically removes deselected features — files, database schemas, config sections, frontend pages, sidebar nav items, and dependencies.

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Removal strategy | Code deletion (not config-driven disable) | Clean slate, no dead code |
| When to run | Project init time (one-shot after clone) | Simplest mental model |
| UX | TUI multi-select via `charmbracelet/huh` | Best Go TUI library, polished UX |
| Post-removal cleanup | Fully automated (ent-gen, go mod tidy, npm cleanup, build verify) | Developer shouldn't need manual steps |
| Extensibility | Declarative module manifests + comment markers | Adding future features requires no changes to setup tool core |

## Architecture

### Execution Flow

```
go run cmd/setup/main.go
        │
        ▼
┌─────────────────────┐
│  1. TUI Checklist   │  All features pre-checked. Dev unchecks what to remove.
│     (huh library)   │
└────────┬────────────┘
         ▼
┌─────────────────────┐
│  2. Delete Files    │  Remove standalone files/dirs listed in module manifest.
└────────┬────────────┘
         ▼
┌─────────────────────┐
│  3. Patch Shared    │  Find [feature:X] markers in shared files,
│     Files           │  remove delimited blocks.
└────────┬────────────┘
         ▼
┌─────────────────────┐
│  4. Clean Stale Ent │  Delete generated files for removed schemas before regen.
│     Generated Files │  See "Ent Regeneration Strategy" below.
└────────┬────────────┘
         ▼
┌─────────────────────┐
│  5. Regenerate Ent  │  `go generate ./ent` — rebuilds ORM for remaining schemas.
└────────┬────────────┘
         ▼
┌─────────────────────┐
│  6. Clean Go Deps   │  `go mod tidy` — removes unused Go modules.
└────────┬────────────┘
         ▼
┌─────────────────────┐
│  7. Verify Go Build │  `go build -o /dev/null ./cmd/web` — HARD GATE.
│                     │  If this fails, report error and stop.
└────────┬────────────┘
         ▼
┌─────────────────────┐
│  8. Verify Frontend │  `npx tsc --noEmit` — HARD GATE.
│     Build           │  Catches broken imports/types after page removal.
└────────┬────────────┘
         ▼
┌─────────────────────┐
│  9. Self-Cleanup    │  Remove cmd/setup/, huh dependency, go mod tidy again.
└─────────────────────┘
```

**Note:** No npm cleanup step is needed. No feature currently uses feature-specific npm packages (Stripe is loaded via CDN). All npm dependencies are shared by the core app.

### Module Manifest Structure

Each optional feature is described declaratively as a `Module` struct:

```go
type Module struct {
    Name        string   // Display name in TUI (e.g., "Payment (Stripe)")
    Description string   // One-liner description
    Files       []string // Standalone files to delete
    Dirs        []string // Directories to delete recursively
    EntSchemas  []string // Ent schema files to remove (triggers ent-gen)
    UserEdges   []string // Edge names to remove from ent/schema/user.go
    NavItems    []string // Sidebar nav item titles to remove from AppSidebar.tsx
    RouteNames  []string // Constants to remove from routenames/names.go
    ConfigKeys  []string // Top-level keys to remove from config.yaml
    GoDeps      []string // Go module paths to expect removal via go mod tidy
    MakeTargets []string // Makefile targets to remove
}
```

**Note on npm dependencies:** No feature currently requires npm package removal. Stripe is loaded via CDN (`window.Stripe`), and all other features use only packages shared by the core app (React, Inertia, lucide-react, etc.). If a future feature introduces a feature-specific npm package, add an `NpmDeps []string` field back to the struct.

All modules are registered in `cmd/setup/modules.go`. Adding a new optional feature = defining a new `Module` + wrapping code in markers.

### Comment Markers for Shared Files

Shared files use delimited markers to identify feature-specific blocks. The setup tool finds `[feature:X] start` / `[feature:X] end` and removes everything between them (inclusive).

**Go files (container.go, user.go, router.go, middleware/auth.go, admin.go):**
```go
// [feature:payment] start
c.initPayment()
// [feature:payment] end
```

**TSX files (AppSidebar.tsx):**

Icon imports must be split to separate lines so each feature's icons can be independently marked:
```tsx
import {
  BookOpen,
  Folder,
  LayoutGrid,
  // [feature:files] start
  UploadCloud,
  // [feature:files] end
  // [feature:payment] start
  CreditCard,
  Receipt,
  ShoppingBag,
  Crown,
  // [feature:payment] end
  // [feature:chat] start
  MessageCircle,
  // [feature:chat] end
} from "lucide-react";
```

Nav item blocks:
```tsx
{/* [feature:chat] start */}
{ title: "Chat", href: "/chat", icon: MessageCircle },
{/* [feature:chat] end */}
```

**YAML files (config.yaml):**
```yaml
# [feature:payment] start
payment:
  provider: "stripe"
  stripe:
    secretKey: "sk_test_..."
    publishableKey: "pk_test_..."
    webhookSecret: "whsec_..."
    currency: "usd"
# [feature:payment] end
```

**Makefile:**
```makefile
# [feature:chat] start
.PHONY: chat-clear
chat-clear: ## Clear all chat data
	...
# [feature:chat] end
```

**Go imports (handled specially):** Feature-specific imports in shared files are also wrapped:
```go
import (
    "context"
    // [feature:tasks] start
    "github.com/mikestefanello/backlite"
    // [feature:tasks] end
)
```

### Marker Cleanup

After removing marked blocks, the tool cleans up artifacts:
- Consecutive blank lines collapsed to one
- Trailing commas in Go/TSX after removed blocks
- Empty import groups `import ()` removed
- Dangling commas before closing `}` in TSX/JS import statements

### Ent Regeneration Strategy

`go generate ./ent` reads schemas from `ent/schema/` and generates code, but it does **not** delete files for schemas that no longer exist. If `ent/schema/chatroom.go` is removed but `ent/chatroom.go`, `ent/chatroom_create.go`, etc. remain, the build will fail because those generated files import the now-deleted `ent/chatroom/` package.

**Solution:** Before running `go generate ./ent`, the tool deletes all generated Ent files and directories, preserving only the non-generated scaffolding:

```
Preserved (never deleted):
  ent/schema/          — user-authored schema definitions
  ent/generate.go      — go:generate directive
  ent/entc.go          — Ent codegen configuration

Deleted (regenerated by go generate):
  ent/*.go             — all .go files in ent/ root except generate.go and entc.go
  ent/<entity>/        — all entity subdirectories (constants, predicates)
  ent/enttest/
  ent/hook/
  ent/migrate/
  ent/predicate/
  ent/runtime/
  ent/admin/           — auto-generated admin panel (regenerated from remaining schemas)
```

This is safe because `go generate ./ent` fully recreates all deleted files from the remaining schemas. The admin panel in `ent/admin/` is also regenerated automatically, so it will only reference entities that still exist.

---

## Feature Module Manifests

### 1. Payment (Stripe)

**Files to delete:**
- `pkg/handlers/billing.go`
- `pkg/handlers/plans.go`
- `pkg/handlers/products.go`
- `pkg/handlers/premium.go`
- `pkg/services/payment.go`
- `pkg/services/payment_stripe.go`
- `resources/js/Pages/Billing.tsx`
- `resources/js/Pages/Plans.tsx`
- `resources/js/Pages/Products.tsx`
- `resources/js/Pages/Premium.tsx`
- `resources/js/components/PaymentForm.tsx`
- `resources/js/types/stripe.d.ts`

**Ent schemas to delete:**
- `ent/schema/paymentcustomer.go`
- `ent/schema/paymentmethod.go`
- `ent/schema/paymentintent.go`
- `ent/schema/subscription.go`

**User edges to remove:** `payment_customer`

**Nav items to remove:** Plans, Products, Premium, Billing

**Route names to remove:** Plans, PlansSubscribe, Products, ProductsPurchase, Premium, Billing, BillingCancel

**Config keys to remove:** `payment`

**Go deps removed (via tidy):** `github.com/stripe/stripe-go/v82`

**Marked blocks in shared files:**
- `pkg/services/container.go`: Payment field, initPayment() call, initPayment() function
- `pkg/middleware/auth.go`: RequirePaidUser() function + payment-specific imports (`paymentcustomer`, `paymentintent`, `subscription`)
- `ent/schema/user.go`: payment_customer edge
- `resources/js/components/AppSidebar.tsx`: Plans, Products, Premium, Billing nav items + icon imports (`CreditCard`, `Receipt`, `ShoppingBag`, `Crown`)
- `config/config.yaml`: payment section
- `pkg/routenames/names.go`: payment route constants

### 2. Chat (WebSocket)

**Files to delete:**
- `pkg/handlers/chat.go`
- `resources/js/hooks/useChat.ts`
- `resources/js/types/chat.d.ts`
- `resources/js/hooks/useAudioRecorder.ts`

**Dirs to delete (all contents removed recursively):**
- `static/chat-uploads/` (created at runtime — may not exist; skip if missing)
- `resources/js/Pages/Chat/`
- `resources/js/components/chat/`
- `pkg/chat/`

**Ent schemas to delete:**
- `ent/schema/chatroom.go`
- `ent/schema/chatmessage.go`
- `ent/schema/chatban.go`

**User edges to remove:** `owned_chat_rooms`, `chat_messages`, `chat_bans`, `chat_bans_issued`

**Nav items to remove:** Chat

**Route names to remove:** ChatRooms, ChatRoomCreate, ChatRoom, ChatWebSocket, ChatBanUser, ChatUnbanUser, ChatDeleteRoom

**Config keys to remove:** `chat`

**Make targets to remove:** `chat-clear`

**Marked blocks in shared files:**
- `pkg/services/container.go`: Chat field + import, WebSocketGroup field, initChat() call, Chat shutdown block, initChat() function
- `pkg/handlers/router.go`: WebSocketHandler interface, WebSocket group setup, WebSocket route registration
- `ent/schema/user.go`: 4 chat edges
- `resources/js/components/AppSidebar.tsx`: Chat nav item + icon import (`MessageCircle`)
- `config/config.yaml`: chat section
- `pkg/routenames/names.go`: chat route constants
- `Makefile`: chat-clear target

### 3. Mail (Resend)

**Files to delete:**
- `pkg/handlers/contact.go`
- `pkg/services/mail.go`
- `pkg/services/mail_test.go`
- `pkg/ui/forms/contact.go`
- `pkg/ui/pages/contact.go`
- `pkg/ui/emails/auth.go`

**Ent schemas to delete:** none

**Nav items to remove:** none (mail has no nav presence)

**Config keys to remove:** `mail`

**Go deps removed (via tidy):** `github.com/resend/resend-go/v2`

**Marked blocks in shared files:**
- `pkg/services/container.go`: Mail field, initMail() call, initMail() function
- `pkg/handlers/auth.go` (5 marker zones):
  1. `mail *services.MailClient` field in Auth struct
  2. `h.mail = c.Mail` initialization in Init()
  3. `h.sendVerificationEmail(ctx, u)` call + error handling block in RegisterSubmit()
  4. `sendVerificationEmail()` entire function definition
  5. User lookup + token generation + email sending block in ForgotPasswordSubmit() (from `u, err := h.orm.User.Query()...` through the mail send error handling, lines 404-450). Must include the user query to avoid unused variable compile errors (`u`, `token`, `pt` are only used within the mail-sending code). After removal, the function validates the form and returns `succeed()` immediately — preserving anti-enumeration without DB queries or token creation.
- `config/config.yaml`: mail section
- `pkg/routenames/names.go`: Contact, ContactSubmit route constants

**Important note:** Removing mail does NOT break authentication. Login and register flows still work. The auth handler gracefully degrades:
- **Register:** User is created and logged in. The `sendVerificationEmail()` call is removed, so no verification email is sent, but the account works immediately.
- **Forgot password:** `ForgotPasswordSubmit` validates the form and returns the generic success message without querying the DB or creating tokens. This preserves anti-enumeration (same response regardless of email existence). Password reset is effectively disabled — users cannot receive reset links. The `ResetPassword` page/route still exists but is unreachable without a valid token URL.
- **Verify email:** Route still exists and works if a user somehow has a valid token URL, but no new verification emails are sent.
- The email template file (`pkg/ui/emails/auth.go`) is deleted alongside the mail service since it has no callers after the auth handler markers are removed.

### 4. Background Tasks (Backlite)

**Files to delete:**
- `pkg/handlers/task.go`
- `pkg/tasks/example.go`
- `pkg/tasks/register.go`
- `pkg/ui/forms/task.go`
- `pkg/ui/pages/task.go`

**Dirs to delete:**
- `pkg/tasks/`

**Ent schemas to delete:** none (Backlite manages its own tables)

**Route names to remove:** Task, TaskSubmit, AdminTasks

**Nav items to remove:** none (tasks have no sidebar nav)

**Config keys to remove:** `tasks`

**Go deps removed (via tidy):** `github.com/mikestefanello/backlite`

**Marked blocks in shared files:**
- `pkg/services/container.go`: Tasks field + import, initTasks() call, Tasks shutdown block, initTasks() function
- `pkg/handlers/admin.go`: backlite field, backlite import, backlite handler init, task routes in Routes(), Backlite() helper method
- `config/config.yaml`: tasks section
- `pkg/routenames/names.go`: Task, TaskSubmit, AdminTasks constants

### 5. File Upload

**Files to delete:**
- `pkg/handlers/files.go`
- `resources/js/Pages/UploadFile.tsx`
- `pkg/ui/forms/file.go`
- `pkg/ui/pages/file.go`
- `pkg/ui/models/file.go`

**Ent schemas to delete:** none

**Nav items to remove:** Upload Files

**Config keys to remove:** `files`

**Go deps removed (via tidy):** `github.com/spf13/afero`

**Note:** The `pkg/ui/` files (`forms/file.go`, `pages/file.go`, `models/file.go`) are legacy gomponents UI files that are no longer imported by the handler (which now uses InertiaJS). They are included in the manifest for dead code cleanup.

**Marked blocks in shared files:**
- `pkg/services/container.go`: Files field + import, initFiles() call, initFiles() function
- `resources/js/components/AppSidebar.tsx`: Upload Files nav item + icon import (`UploadCloud`)
- `config/config.yaml`: files section
- `pkg/routenames/names.go`: Files, FilesSubmit constants

---

## Setup Tool File Structure

```
cmd/setup/
├── main.go          # Entry point: TUI form, orchestration, build verification
├── modules.go       # Module manifest definitions (all 5 features)
├── remover.go       # Core logic: file deletion, marker-based patching, cleanup
└── cleanup.go       # Post-removal: ent cleanup, ent-gen, go mod tidy, build verify, self-removal
```

## TUI Behavior

```
$ go run cmd/setup/main.go

  Pagode Setup

  Select the features to INCLUDE in your project.
  Uncheck any features you don't need.

  > [x] Payment (Stripe)      Subscription billing, one-time payments, premium access
    [x] Chat (WebSocket)       Real-time community chat with rooms, voice messages
    [x] Mail (Resend)          Transactional emails (verification, password reset, contact)
    [x] Background Tasks       Async job queue with admin monitoring UI
    [x] File Upload            File upload handling with filesystem abstraction

  [ Submit ]

  Removing: Chat, Background Tasks
  [1/8] Deleting feature files...          done (23 files, 4 directories)
  [2/8] Patching shared files...           done (8 files patched)
  [3/8] Cleaning stale Ent files...        done
  [4/8] Regenerating Ent ORM...            done
  [5/8] Cleaning Go dependencies...        done
  [6/8] Verifying Go build...              done
  [7/8] Verifying frontend build...        done
  [8/8] Cleaning up setup tool...          done

  Setup complete! Run `make run` to start your project.
```

All features are pre-checked. The developer unchecks what they don't want. If nothing is unchecked, the tool exits with "No changes needed."

## Preconditions

- **Clean git working tree required.** The tool checks `git status --porcelain` on startup and refuses to run if there are uncommitted changes. If the build fails after removal, the developer can recover with `git checkout .` to restore all files.

## Safety Guarantees

1. **Go build verification** — After all removals, `go build -o /dev/null ./cmd/web` must pass. If it fails, the tool reports the error and exits without self-cleanup so the developer can debug.
2. **Frontend build verification** — After all removals, `npx tsc --noEmit` must pass. This catches broken TypeScript imports and type errors from removed pages/components.
3. **Marker-based patching** — No fragile line-number or regex matching. Markers are explicit, grep-able, and won't drift with code changes.
4. **Idempotent** — Running the tool twice with the same selections is safe (already-deleted files are skipped, already-removed markers are no-ops).
5. **Feature isolation** — Feature modules never import each other (Payment code never references Chat, etc.). Features do integrate through shared handlers (Auth uses Mail, Admin uses Tasks, middleware uses Payment), but all such coupling points are wrapped in `[feature:X]` markers and surgically removed. Removing one feature never affects another.

## Extensibility: Adding a New Feature Module

To make a new feature optional:

1. **Wrap feature-specific blocks** in shared files with `[feature:name] start/end` markers
2. **Define a Module struct** in `cmd/setup/modules.go`
3. **Append it** to the `Modules` slice

No changes needed to `main.go`, `remover.go`, or `cleanup.go`.

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| Dirty git working tree | Tool refuses to start: "Please commit or stash changes before running setup." |
| All features kept | Tool exits: "No changes needed." |
| All features removed | Works — leaves core auth + dashboard skeleton |
| Run after already modifying code | Works if markers are intact; warns if markers not found |
| Run twice | Idempotent — skips already-removed files/markers |
| Go build fails after removal | Tool stops, reports error, does NOT self-cleanup. Use `git checkout .` to recover. |
| Frontend build fails after removal | Same as above — stops before self-cleanup. |
| `static/chat-uploads/` doesn't exist | Skipped silently (directory is created at runtime, may not exist in fresh clone) |
| Mail removed but user expects email verification | Auth still works, just doesn't send emails. README/comments should note this. |

## Implementation Order

1. Add `[feature:X]` comment markers to all shared files (non-breaking, markers are just comments)
2. Build `cmd/setup/` tool: remover logic, module manifests, TUI
3. Test each feature removal independently (verify build passes)
4. Test removing all features at once
5. Test removing various combinations
6. Add `make setup` target to Makefile
7. Update project README with setup instructions
