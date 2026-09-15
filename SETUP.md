# Setup Guide

This document explains the interactive `setup` tool in Pagode:

1. what it is
2. how to use it safely
3. how it works internally
4. how to add a new optional feature correctly

The guide is written for both users of the template and maintainers of the template.

## What Setup Is

`setup` is an interactive project-pruning tool.

Pagode ships with multiple optional features enabled in the template repo. The setup CLI lets a developer keep only the features they want in their own project. Instead of starting from a minimal skeleton and adding capabilities one by one, the template starts full-featured and `setup` removes the parts you do not want.

Today, optional features are defined in [`cmd/setup/modules.go`](./cmd/setup/modules.go):

- `Payment (Stripe)`
- `Chat (WebSocket)`
- `Mail (Resend)`
- `Background Tasks`
- `File Upload`

The tool is intentionally destructive:

- it deletes feature-specific files
- it removes feature-specific blocks from shared files
- it regenerates Ent if schemas were removed
- it runs hard validation gates
- on success, it removes itself from the generated app

This means `setup` is meant to be run early in a project, usually once.

## How To Use It

From the repository root:

```bash
make setup
```

This is a convenience wrapper for:

```bash
go run ./cmd/setup/
```

The UI shows a checklist of features to include.

- Checked = keep the feature
- Unchecked = remove the feature

When the selection is confirmed, the tool executes the full removal pipeline.

### Preconditions

`setup` requires a clean git worktree before it runs.

If there are uncommitted changes, it exits with:

```text
Please commit or stash changes before running setup.
```

This is deliberate. The tool can delete many files, so rollback must stay trivial.

### Safe Workflow

Recommended usage:

```bash
git clone <repo>
cd pagode
npm install
go run ./cmd/setup/
```

If setup fails mid-run during development of the template, recover with:

```bash
git restore --worktree --source=HEAD -- .
```

If you are testing the setup tool itself, use the same pattern between tests so every run starts from a clean template state.

## What Setup Actually Does

The orchestration lives in [`cmd/setup/main.go`](./cmd/setup/main.go).

After the user chooses which features to keep, setup removes the unchecked modules in 8 steps:

1. Delete feature files
2. Patch shared files by removing `[feature:<name>]` blocks
3. Regenerate Ent ORM if schemas were removed
4. Clean stale Ent generated files
5. Run `go mod tidy`
6. Verify `go build -o /dev/null ./cmd/web`
7. Verify `npx tsc --noEmit`
8. Self-cleanup: remove `cmd/setup`, remove the `setup` Make target, and clean the `huh` dependency

If any hard gate fails, setup stops and does **not** self-delete.

That is intentional: the user should still have the tool available while fixing the template.

## How Features Are Defined

Each optional feature is declared as a `Module` in [`cmd/setup/modules.go`](./cmd/setup/modules.go).

Each module has:

- `Name`: user-facing label in the checklist
- `FeatureName`: marker name used in shared files
- `Description`: short description shown in the UI
- `Files`: specific files to delete
- `Dirs`: directories to remove recursively
- `EntSchemas`: Ent schema files to remove, which also triggers Ent regeneration

Example shape:

```go
type Module struct {
    Name        string
    FeatureName string
    Description string
    Files       []string
    Dirs        []string
    EntSchemas  []string
}
```

## The Marker System

Feature-specific code inside shared files is removed with comment markers.

The remover is implemented in [`cmd/setup/remover.go`](./cmd/setup/remover.go).

Supported marker styles:

- Go / JS / TS:
```go
// [feature:tasks] start
...
// [feature:tasks] end
```

- YAML / Makefile:
```yaml
# [feature:mail] start
...
# [feature:mail] end
```

- JSX / TSX:
```tsx
{/* [feature:files] start */}
...
{/* [feature:files] end */}
```

The remover scans the repository for those blocks and removes them for every feature being removed.

### Important Scanner Rules

The scanner:

- scans normal project code, including `cmd/web`
- skips `.git`, `node_modules`, and `vendor`
- skips `cmd/setup` intentionally so the tool does not rewrite itself while running

This behavior is covered by [`cmd/setup/remover_test.go`](./cmd/setup/remover_test.go).

## Ent Regeneration Behavior

If a feature owns Ent schemas, setup must update generated ORM code safely.

That logic lives in [`cmd/setup/cleanup.go`](./cmd/setup/cleanup.go).

The current sequence is:

1. remove stale generated admin files in `ent/admin`
2. run Ent codegen directly from `ent/` with:
```bash
go run -mod=mod entc.go
```
3. remove orphaned generated Ent artifacts whose schema files no longer exist

This order matters.

Earlier versions broke because:

- stale generated `ent/admin/*.go` files could import deleted entity packages
- deleting too much under `ent/` before codegen could break schema compilation
- `go generate ./ent` was less reliable here than running `entc.go` directly from `ent/`

Those behaviors are covered by [`cmd/setup/cleanup_test.go`](./cmd/setup/cleanup_test.go).

## Validation Gates

Setup deliberately ends with hard validation:

- `go mod tidy`
- `go build -o /dev/null ./cmd/web`
- `npx tsc --noEmit` if `node_modules` exists

Why this matters:

- marker coverage mistakes surface immediately
- dead imports are caught immediately
- route/menu leftovers are caught immediately
- frontend references to removed routes/pages are caught immediately

If `node_modules` is missing, the frontend check is skipped with a visible warning.

## Self-Cleanup

On success, setup removes itself from the generated application:

- deletes [`cmd/setup`](./cmd/setup)
- removes the `setup` target from [`Makefile`](./Makefile)
- removes `github.com/charmbracelet/huh` from `go.mod` / `go.sum`

This is intentional: the generated app should not carry template-only setup machinery after pruning is complete.

## How To Add A New Optional Feature

This is the maintainer checklist for building a new opt-in/opt-out feature correctly.

### 1. Decide What Is Truly Feature-Owned

Start by identifying:

- standalone files that can be deleted
- directories that can be deleted recursively
- shared files that need conditional blocks
- Ent schemas owned by the feature
- startup wiring in `cmd/web`
- service/container wiring
- routes and route names
- admin links, sidebar links, navigation links
- frontend pages/components/types
- tests owned by the feature

If a file is purely feature-owned, prefer deleting the whole file.
If a file is shared, use markers.

### 2. Add The Module Manifest Entry

Add a new `Module` in [`cmd/setup/modules.go`](./cmd/setup/modules.go).

Populate:

- `FeatureName`
- `Files`
- `Dirs`
- `EntSchemas`

Be explicit. If setup should remove it, list it or mark it.

### 3. Mark Shared Code With `[feature:<name>]`

Every shared reference to the feature must be wrapped.

Common places that are easy to miss:

- imports used only by the feature
- container fields and initialization
- shutdown hooks
- route name constants
- menu/sidebar links
- `cmd/web` startup code
- admin handlers and admin routes
- Ent schema edges pointing to feature entities
- tests that depend on the feature

The most common setup failures come from missing markers around tiny leftovers such as:

- one route name
- one menu item
- one import
- one startup hook

### 4. Handle Ent Carefully

If the feature owns Ent schemas:

- list them in `EntSchemas`
- ensure any references from remaining schemas are feature-marked
- ensure no generated admin file or shared Ent support file is deleted too early

If the feature does **not** own schemas, leave `EntSchemas` empty.

### 5. Consider Frontend And Build-Time References

Do not stop at file deletions.

Also check:

- sidebar links
- routes referenced in components
- page imports
- TS-only types
- motion/React/utility imports that become unused after marker removal

### 6. Test The Feature In Isolation

For every new optional feature:

1. reset the repo to a clean state
2. run `go run ./cmd/setup/`
3. remove only that feature
4. confirm the full pipeline reaches self-cleanup

At minimum, verify:

- `go mod tidy`
- Go build
- TypeScript build

### 7. Fix Marker Coverage, Not Just The Symptom

When setup fails after removing a feature, the right fix is usually:

- add a missing marker
- move a feature-only import inside a marker
- include a missing file in the module manifest
- expand the scanner/test coverage

Do not paper over failures by relaxing validation gates.

## Patterns For Shared Code

Use these patterns consistently.

### Feature-only import

```go
// [feature:mail] start
importedThing "example.com/mail"
// [feature:mail] end
```

### Feature-only struct field

```go
// [feature:tasks] start
Tasks *backlite.Client
// [feature:tasks] end
```

### Feature-only startup code

```go
// [feature:tasks] start
tasks.Register(c)
c.Tasks.Start(context.Background())
// [feature:tasks] end
```

### Feature-only menu link

```go
// [feature:files] start
MenuLink(r, "Files", routenames.Files),
// [feature:files] end
```

## What Not To Do

- Do not skip the validation gates to make setup “pass”
- Do not put feature markers inside `cmd/setup` expecting them to be applied during the same run
- Do not assume deleting feature-owned files is enough; shared references are the real source of breakage
- Do not delete all generated Ent files before regeneration
- Do not use `go generate ./ent` here unless you have verified it works for the removal path; the current implementation intentionally runs `entc.go` directly

## Recommended Maintainer Workflow

When adding or modifying an optional feature:

1. implement the feature normally
2. decide which parts are optional
3. add the module manifest entry
4. wrap all shared references with `[feature:<name>]`
5. run setup removing only that feature
6. fix every failure by tightening marker coverage or manifest coverage
7. add tests if the failure exposed a setup-system bug rather than just missing markers
8. repeat until the full setup pipeline succeeds

## Summary

`setup` is not just a UI checklist. It is a repository transformation pipeline with hard guarantees.

To work correctly, an optional feature must satisfy all three:

1. feature-owned files are declared
2. shared references are marker-guarded
3. the remaining project still passes Go and frontend verification after removal

If you treat feature removal as a first-class architecture concern while adding the feature, the setup system stays simple and reliable.
