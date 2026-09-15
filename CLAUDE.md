# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Pagode** is a full-stack web application starter kit combining Go backend with React + InertiaJS frontend, styled with Tailwind CSS v4 and shadcn/ui components. It uses SQLite for data storage and provides a complete foundation for building modern web applications.

## Key Technologies

- **Backend**: Go 1.24, Echo v4 web framework, Ent ORM, SQLite database
- **Frontend**: React 19, TypeScript, InertiaJS, Tailwind CSS v4, shadcn/ui
- **Build Tools**: Vite (frontend), Air (Go hot reload)

## Common Development Commands

### Essential Commands
- `make run` - Start the application (default: http://localhost:8000)
- `make watch` - Start with hot reload (requires `make air-install` first)
- `make test` - Run all tests
- `go build -o /dev/null ./cmd/web` - Test compilation without generating binaries
- `make admin email=user@example.com` - Create admin user account

### Database & ORM
- `make ent-gen` - Generate Ent ORM code (run after schema changes)
- `make ent-new name=Entity` - Create new Ent entity
- `make ent-install` - Install Ent code generation tools

### Frontend Development
Frontend assets are built automatically when the app starts. Vite handles hot module replacement for React components.

## Architecture Overview

### Service Container Pattern
All application dependencies are managed through a service container (`pkg/services/container.go`) that houses:
- Database client, Cache, Authentication, Mail, Tasks, Files, Configuration
- Use `services.NewContainer()` to initialize, `container.Shutdown()` to cleanup
- Available in handler `Init()` methods for dependency injection

### Directory Structure
- `cmd/` - Application entry points (`web/` for main app, `admin/` for CLI tools)
- `pkg/handlers/` - HTTP request handlers with route registration
- `pkg/middleware/` - Custom middleware
- `pkg/services/` - Service container and business logic
- `ent/` - Ent ORM generated code and schemas (`ent/schema/` for entity definitions)
- `resources/js/` - React/TypeScript frontend code
- `resources/js/Pages/` - InertiaJS page components
- `config/` - Configuration management (YAML + environment overrides)

### Frontend Architecture
- **InertiaJS** bridges Go controllers with React components
- No separate API layer - controllers return JSON "page props"
- Server-side routing with client-side interactivity
- Pages in `resources/js/Pages/`, components in `resources/js/components/`

### Authentication & Authorization
- Session-based auth with JWT email verification
- `middleware.RequireAuthentication()` and `middleware.RequireAdmin()` for route protection
- `middleware.RequirePaidUser()` for premium content protection
- User entity has `Admin` boolean field for admin access
- Admin panel auto-generated for all Ent entities

### Payment Integration
- Comprehensive Stripe integration for subscriptions and one-time payments
- Provider abstraction pattern allows easy switching between payment processors
- Database entities: PaymentCustomer, PaymentMethod, Subscription, PaymentIntent
- PCI-compliant payment processing using Stripe Elements
- Premium content access control based on payment status

### Database
- SQLite with automatic migrations on startup
- Separate in-memory test database when `PAGODE_APP_ENVIRONMENT=test`
- WAL mode enabled for better concurrent access

## Development Patterns

### Adding New Routes
1. Create handler in `pkg/handlers/` 
2. Embed required services in handler struct
3. Use `Register(new(YourHandler))` in `init()` 
4. Implement `Init(c *services.Container)` for dependency injection
5. Add routes in `Routes(g *echo.Group)` method with named routes
6. Route names should be defined in `routenames` package

### Testing
- Use `config.SwitchEnvironment(config.EnvTest)` before creating container
- Test helpers in `pkg/handlers/router_test.go` for HTTP testing
- goquery available for HTML response testing
- **IMPORTANT**: Always run `make test` after making changes to ensure all tests pass

### Background Tasks
- Uses Backlite (SQLite-based task queues)
- Register queues in `pkg/tasks/register.go`
- Admin UI available for monitoring tasks
- Task dispatcher starts automatically with app

### Forms & Validation
- Embed `form.Submission` in form structs
- Use struct tags: `form:"field_name" validate:"required"`
- `form.Submit()` handles parsing and validation
- Inline validation errors automatically displayed

## Configuration

Configuration uses Viper with `config/config.yaml` as base. Environment variables override config using `PAGODE_` prefix (e.g., `PAGODE_HTTP_PORT` overrides `http.port`).

Key config paths:
- `Config.App.Environment` - Controls behavior (Local/Test/Production)
- `Config.Database.Connection` - SQLite database path
- `Config.App.EncryptionKey` - Must change for production (sessions/JWT)

## File Generation

When modifying Ent schemas:
1. Edit schema files in `ent/schema/`
2. Run `make ent-gen` to regenerate all Ent code
3. Admin panel code auto-generates for all entities

## Payment System

### Configuration
Payment settings are configured in `config/config.yaml`:
```yaml
payment:
  provider: "stripe"
  stripe:
    secretKey: "sk_test_your_stripe_secret_key_here"
    publishableKey: "pk_test_your_stripe_publishable_key_here"
    webhookSecret: "whsec_your_webhook_secret_here"
    currency: "usd"
```

### Payment Routes & Pages
- `/plans` - Subscription management (monthly/yearly billing)
- `/products` - One-time product purchases
- `/premium` - Protected premium content (requires payment)
- `/billing` - Subscription and payment method management

### Payment Entities (Ent ORM)
- **PaymentCustomer**: Links users to Stripe customers
- **PaymentMethod**: Stores payment method metadata (no sensitive card data)
- **Subscription**: Tracks subscription status and billing periods
- **PaymentIntent**: Handles one-time payments

### Security & Compliance
- No sensitive card data stored in database
- Stripe Elements for PCI-compliant card collection
- Server-side payment validation and processing
- Payment method IDs used instead of raw card data

### Premium Access Control
Use `middleware.RequirePaidUser(orm)` to protect routes requiring payment:
- Checks for active subscriptions (`subscription.StatusActive`)
- Checks for successful payment intents (`paymentintent.StatusSucceeded`)
- Redirects unauthorized users to purchase pages

### Development & Testing
- Use Stripe test cards: `4242 4242 4242 4242` (success)
- Test environment automatically uses Stripe test mode
- Payment status visible in admin panel for debugging

## Modular Setup CLI (Feature Opt-In/Opt-Out)

Pagode ships with optional features that can be removed at project init time via `make setup` (or `go run ./cmd/setup/`). The tool presents a TUI checklist and surgically removes deselected features — files, schemas, config sections, frontend pages, sidebar items, and dependencies.

Current optional features: Payment (Stripe), Chat (WebSocket), Mail (Resend), Background Tasks (Backlite), File Upload.

**IMPORTANT: Every new feature that is not part of the core (auth, dashboard, admin) MUST be built as an optional module so users can opt out via the setup CLI.**

### The Marker System

Feature-specific code inside shared files is wrapped with comment markers. The setup tool finds and removes everything between `[feature:X] start` / `[feature:X] end` (inclusive).

**Go / JS / TS files:**
```go
// [feature:myfeature] start
c.initMyFeature()
// [feature:myfeature] end
```

**YAML / Makefile:**
```yaml
# [feature:myfeature] start
myfeature:
  key: value
# [feature:myfeature] end
```

**TSX / JSX (inside JSX expressions):**
```tsx
{/* [feature:myfeature] start */}
{ title: "My Feature", href: "/my-feature", icon: SomeIcon },
{/* [feature:myfeature] end */}
```

Markers also wrap feature-specific imports so they are cleanly removed:
```go
import (
    "context"
    // [feature:myfeature] start
    "github.com/example/myfeature-dep"
    // [feature:myfeature] end
)
```

### How to Add a New Optional Feature

1. **Build the feature normally** — handlers, services, pages, schemas, routes, etc.
2. **Identify what is feature-owned vs shared:**
   - Feature-owned files/dirs: will be deleted entirely
   - Shared files: need `[feature:name]` markers around feature-specific blocks
3. **Add a `Module` entry** in `cmd/setup/modules.go`:
   ```go
   Module{
       Name:        "My Feature",
       FeatureName: "myfeature",  // used in [feature:myfeature] markers
       Description: "Short description for TUI",
       Files:       []string{"pkg/handlers/myfeature.go", "resources/js/Pages/MyFeature.tsx"},
       Dirs:        []string{"pkg/myfeature/"},
       EntSchemas:  []string{"ent/schema/myentity.go"},
   }
   ```
4. **Wrap ALL shared references** with `[feature:myfeature]` markers. Common places:
   - `pkg/services/container.go` — field, init call, init function, shutdown block
   - `pkg/handlers/router.go` — route registration
   - `ent/schema/user.go` — edges pointing to feature entities
   - `resources/js/components/AppSidebar.tsx` — nav items and icon imports
   - `config/config.yaml` — feature config section
   - `pkg/routenames/names.go` — route name constants
   - `pkg/handlers/admin.go` — admin routes/handlers if applicable
   - `pkg/middleware/auth.go` — feature-specific middleware
   - `Makefile` — feature-specific targets
5. **Test removal in isolation:** reset repo, run `go run ./cmd/setup/`, remove only the new feature, confirm Go build (`go build -o /dev/null ./cmd/web`) and TS build (`npx tsc --noEmit`) pass.

### Critical Rules for Module Authors

- **Never leave unmarked references** in shared files. The #1 cause of setup failures is a missing marker around a single import, route name, menu item, or init call.
- **Feature modules must not import each other.** Payment code must never reference Chat, etc. Cross-feature coupling in shared handlers (e.g., Auth uses Mail) must be marker-guarded.
- **If the feature owns Ent schemas**, list them in `EntSchemas`. Also marker-guard any edges from `ent/schema/user.go` (or other shared schemas) to feature entities.
- **Prefer deleting whole files** over marking. If a file is purely feature-owned, add it to `Files`/`Dirs` in the manifest rather than marking its contents.
- **Do not modify `cmd/setup/` core files** (`main.go`, `remover.go`, `cleanup.go`) when adding a feature. Only add a new `Module` entry in `modules.go` and add markers to your code.
- Full design details are in `docs/plans/2026-03-27-modular-setup-cli-design.md`. Operational reference is in `SETUP.md`.

## Important Notes

- Admin panel is dynamically generated - entity constraints must be defined in Ent schemas
- CSRF protection enabled by default for non-GET requests
- Static files served from `static/` directory at `/files/` URL prefix
- Email client is skeleton - must implement `MailClient.send()` method
- Cache is in-memory (otter library) - tags supported but come with performance cost
- Payment credentials should be stored as environment variables in production