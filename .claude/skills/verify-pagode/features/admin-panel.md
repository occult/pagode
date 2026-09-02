# Admin panel

An admin manages users from `/admin/users` (list, add, edit, delete), edits any entity through the generated pages under `/admin/users/entity/<entity>` (lowercase Go type name: `user`, `chatroom`, `chatmessage`, `chatban`, `passwordtoken`, `paymentcustomer`, `paymentmethod`, `subscription`, `paymentintent`), and watches background tasks under `/admin/tasks`.

## Sub-features

- `admin-users-list` paginates all users, 10 per page.
- `admin-user-add` creates a user from the `Add User` dialog.
- `admin-user-edit` changes name, email, admin flag of a user.
- `admin-user-delete` removes a user.
- `admin-entity` lists/adds/edits/deletes any Ent entity through its generated pages.
- `admin-tasks` shows running, succeeded, failed, and upcoming Backlite tasks.

## How to get to it (user POV)

- Sidebar `Admin Panel` link (shown only when the signed-in user is an admin), which opens `/admin/users`.
- Open `/admin/users/entity/<entity>` or `/admin/tasks` directly.

## Driving it with Playwright

Preconditions:

- Instance healthy; `scripts/seed-admin.sh` has run and `tmp/verify-pagode/run/18000/admin-credentials.env` exists.
- Log in with those credentials through `loginViaUI`.

- **Open.** Click sidebar `Admin Panel` (or `page.goto("/admin/users")`). The user table lists the seeded admin.
- **Add.** Click `getByRole("button", { name: "Add User" })`; in the `Add User` dialog fill `#name`, `#email`, set the password and admin fields, submit. Toast `User successfully created.` appears and the new row is listed; `db("select admin from users where email = '<email>'")` matches the flag chosen.
- **Edit.** Click the row's edit button, change `#name` in the `Edit User` dialog, submit. The row shows the new name and the DB agrees.
- **Delete.** Click the row's delete button and confirm. The row disappears and `db` count for that email is 0.
- **Entity pages.** `page.goto("/admin/users/entity/chatroom")` lists rooms with add/edit/delete links; create one there and confirm it appears on `/chat`.
- **Tasks.** `page.goto("/admin/tasks")` renders the Backlite UI; submit the demo task from `/task` first to have rows under `Upcoming` or `Succeeded`.
- **Proof.** Screenshots of the list after each mutation plus the `users` rows.

## Gotchas

- As of this map, `/admin/users`, its add/edit/delete routes, and every `/admin/users/entity/<entity>` page have no authentication or admin middleware in `pkg/handlers/admin.go`: signed-out `GET /admin/users` and `GET /admin/users/entity/user` returned 200 with user data (including emails) during the skill's first live run, while `/admin/tasks` returned 401. Treat unauthenticated access as a bug to report, and do not use it as a shortcut for proofs.
- The entity path under `/admin/users/entity/...` is the lowercased Ent type name (`user`, `chatroom`); `User` or `ChatRoom` is a 404.
- Deleting a user cascades to owned chat rooms via `ON DELETE SET NULL` on the room owner; rooms survive without an owner.
- Backlite pages are server-rendered HTML, not Inertia; assert on page text, not toasts.
