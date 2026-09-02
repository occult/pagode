# Community chat

Any visitor can browse rooms and read messages; signed-in users create rooms, and anyone with a nickname can join a room and talk in real time over a WebSocket. Room owners can ban participants and delete rooms.

## Sub-features

- `chat-list` shows every room with owner, lock icon for password rooms, and online count on `/chat`.
- `chat-create` creates a public room, optionally password protected, from the `Create Room` dialog (signed-in only, at most 5 rooms per user).
- `chat-join` opens `/chat/rooms/:id`, connects the WebSocket, and shows `Connected` plus the participant list.
- `chat-send` sends a text message with Enter or the send button; it appears for every participant and is stored.
- `chat-guest` prompts a signed-out visitor for a nickname before joining; the nickname is kept in the browser.
- `chat-password` prompts for the room password before connecting; a wrong password is refused.
- `chat-attach` sends an image, a camera photo, or a voice recording as a message.
- `chat-moderate` lets the owner ban/unban a participant and delete the room.

## How to get to it (user POV)

- Choose `Chat` in the sidebar, or open `/chat` (also works signed out, in the public layout).
- Choose a room card to enter it, or open `/chat/rooms/<id>` directly.
- Choose `Create Room` on `/chat` (signed-in only).

## Driving it with Playwright

Preconditions:

- Instance healthy; a signed-in user (use `registerViaUI` from `scenarios/lib/auth.mjs`).
- Room name from `unique()`; `chat_rooms.name` is unique.

- **Open the list.** `page.goto("/chat")`. `getByRole("heading", { name: "Community Chat" })` is visible; an empty database shows `No chat rooms yet`.
- **Create a room.** Click `getByRole("button", { name: "Create Room" })`. In `getByRole("dialog", { name: "Create Chat Room" })` fill `#room-name` (and `#room-pw` for a password room) and click its `getByRole("button", { name: "Create", exact: true })`. Toast `Room created successfully!` appears and the card with the room name is listed.
- **Persistence.** `db("select id, is_public, user_owned_chat_rooms from chat_rooms where name = '<name>'")` returns one row owned by the user's id.
- **Enter.** Click `getByText("<name>")`. URL becomes `/chat/rooms/<id>`; the composer `getByPlaceholder("Type a message...")` is visible and the header shows `Connected` and `Online (1)`.
- **Send.** Fill the composer and press `Enter`. The message bubble with the text appears; `db("select sender_name from chat_messages where body = '<text>'")` returns the user's name.
- **Guest path.** In a fresh `context.newPage()` with no session open `/chat/rooms/<id>`; the `Choose a nickname` dialog appears, fill `getByPlaceholder("Your nickname")` and submit. The guest's join appears in the first page as `<nick> joined` and the participant count reads `Online (2)`.
- **Password room.** Entering a locked room shows the password prompt (`getByPlaceholder("Room password")`); a wrong value keeps the prompt, the right one connects.
- **Proof.** Run `node scripts/pw.mjs scenarios/chat-room.mjs`; screenshot 06 shows the sent bubble and `run.log` shows the `chat_messages` row.

## Gotchas

- The room card is not a link or button; click the visible room name text.
- The `Create` button in the dialog is disabled until `#room-name` has non-blank text; `exact: true` avoids matching `Create Room` behind the dialog.
- Messages are stored by the server after broadcast; assert the DB row after the bubble appears, not before.
- Rate limit: 10 messages per 10 seconds per connection. A burst test beyond that gets dropped messages by design.
- Attachments are written to `static/chat-uploads/` in the repo, not the instance's uploads directory, and are served from `/files/chat-uploads/<name>`. Cleanup does not delete them.
- The nickname for guests lives in `localStorage` under `chat_nickname`; a new browser context is a new guest.
- Radix logs a `Missing Description or aria-describedby` console warning for the dialog; it is noise, not a failure.
