# Profile settings

A signed-in user edits their name and email, changes their password, switches appearance, and can delete their account after confirming with their password.

## Sub-features

- `profile-info` updates name and email from `/profile/info`.
- `profile-password` changes the password after checking the current one on `/profile/password`.
- `profile-appearance` switches light/dark/system on `/profile/appearance`.
- `profile-delete` deletes the account from the `Delete account` section on `/profile/info`.

## How to get to it (user POV)

- Sidebar user button, then `Settings` (opens `/profile/info`).
- Settings sub-navigation: `Profile`, `Password`, `Appearance`.
- Open `/profile/info`, `/profile/password`, `/profile/appearance` directly.

## Driving it with Playwright

Preconditions:

- Instance healthy; a signed-in user created with `registerViaUI` (keep its password for the password checks).

- **Open profile.** `page.goto("/profile/info")`. Heading `Profile information` is visible with `#name` and `#email` prefilled.
- **Update info.** Fill `#name` with a new value and click `getByRole("button", { name: "Save" })`. Toast `Your profile has been updated.` appears; `db("select name from users where id = <id>")` shows the new name. Submitting without changes shows `Nothing to update.`.
- **Change password.** On `/profile/password` fill `#current_password`, `#password`, `#password_confirmation` (min 8 chars) and click `Save password`. Toast `Your password has been updated successfully.` appears. Prove it by logging out and back in with the new password. A wrong current password shows `The current password you entered is incorrect.`.
- **Appearance.** On `/profile/appearance` choose `Dark`/`Light`/`System`; the `html` element's class toggles `dark` and the choice persists after reload.
- **Delete account.** On `/profile/info` click `getByRole("button", { name: "Delete account" })`, fill `#password` in the dialog, confirm. Toast `Your account has been deleted.` appears on `/`; `db("select count(*) from users where id = <id>")` is 0 and `/dashboard` answers 401.
- **Proof.** Screenshot after each save plus the `users` row before and after.

## Gotchas

- Validation errors redirect back with the warning toast `Please fix the errors in the form and try again.` and inline messages; the field values are not preserved.
- Changing the email does not re-verify it; `verified` stays as it was.
- Appearance is stored client-side (cookie and `localStorage`), not in the database; there is no DB proof for it.
- The `Delete account` dialog and the login form both use `#password`; scope to `getByRole("dialog")`.
