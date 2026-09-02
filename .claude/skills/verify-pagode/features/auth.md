# Authentication

A visitor creates an account, is signed in immediately, can sign out from the sidebar user menu, sign back in, and request a password reset. Signed-out visitors get a 401 page on protected routes.

## Sub-features

- `auth-register` creates a user from `/user/register` and lands on `/dashboard` signed in.
- `auth-login` signs an existing user in from `/user/login`.
- `auth-logout` ends the session from the sidebar user menu and returns to `/`.
- `auth-guard` answers 401 on `/dashboard`, `/profile/*`, `/files`, `/plans`, `/products`, `/premium`, `/billing`, `/admin/tasks` when signed out.
- `auth-forgot` accepts an email on `/user/password` and always confirms with the same message.
- `auth-verify` marks the account verified when the emailed `/email/verify/:token` link is opened.

## How to get to it (user POV)

- Open `/user/register` or choose `Sign up` on the login page.
- Open `/user/login` or choose `Log in` on the register page.
- Choose the sidebar user button (shows the user's name) then `Log out`.
- Choose `Forgot password?` on the login page.
- Open the verification link from the registration email.

## Driving it with Playwright

Preconditions:

- Instance healthy at `http://localhost:18000` (`scripts/doctor.sh` exits 0).
- Use a `unique()` email; `users.email` is unique.

- **Open register.** `page.goto("/user/register")`. `getByRole("heading", { name: "Create an account" })` is visible.
- **Fill the form.** Fill `#name`, `#email`, `#password`, `#password-confirm` (same password). Both password fields are required and must match.
- **Submit.** Click `getByRole("button", { name: "Create Account" })`. URL becomes `/dashboard`; toast `Your account has been created. You are now logged in.` appears.
- **Persistence.** `db("select name, email, verified, admin from users where email = '<email>'")` returns one row with `admin = 0` and `verified = 0`.
- **Log out.** Click `getByRole("button", { name: /<name>/ }).first()` in the sidebar, then `getByRole("menuitem", { name: "Log out" })`. URL becomes `/`; toast `You have been logged out successfully.` appears.
- **Guard.** `page.goto("/dashboard")` responds with status 401 and renders the error page.
- **Log in.** `page.goto("/user/login")`, fill `#email` and `#password`, click `getByRole("button", { name: "Log in" })`. URL becomes `/dashboard`; toast `Welcome back, <name>. You are now logged in.` appears.
- **Wrong password.** Same form with a bad password stays on `/user/login` with toast `Invalid credentials. Please try again.`.
- **Forgot password.** On `/user/password` fill `#email`, submit. Toast `An email containing a link to reset your password will be sent to this address if it exists in our system.` appears whether or not the email exists; `db("select count(*) from password_tokens")` grows by one only for an existing user.
- **Proof.** Run `node scripts/pw.mjs scenarios/register-login.mjs`; the evidence directory holds screenshots 01-08 and `run.log` with the `users` row.

## Gotchas

- `getByText("Create an account")` matches two elements (a screen-reader-only logo link and the h1). Use the heading role.
- The register handler logs the user in before sending the verification email; the mail client is a skeleton, so the server log shows a send error and no email is ever delivered. `auth-verify` can only be driven by reading the token from the log or DB, which is not a user path; report it as unreachable rather than faking it.
- Validation failures on register redirect to `/user/login` with a danger toast, not back to the form.
- Login lowercases the email before lookup; register stores it as typed.
- Protected pages return 401 rather than redirecting, so `waitForURL("**/user/login")` after a guard check will hang.
