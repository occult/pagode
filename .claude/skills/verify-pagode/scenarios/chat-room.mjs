// Feature: chat (features/chat.md). Create a room, enter it, send a message over the WebSocket.
import { registerViaUI } from "./lib/auth.mjs";

export default async function (ctx) {
  const { page, step, db, expect, unique, log } = ctx;
  const user = await registerViaUI(ctx);
  const roomName = unique("verify-room");
  const body = `hello from verify ${unique("msg")}`;

  await step("open chat index", async () => {
    await page.goto("/chat");
    await page.getByRole("heading", { name: "Community Chat" }).waitFor();
  });

  await step("create a public room", async () => {
    await page.getByRole("button", { name: "Create Room" }).click();
    const dialog = page.getByRole("dialog", { name: "Create Chat Room" });
    await dialog.locator("#room-name").fill(roomName);
    await dialog.getByRole("button", { name: "Create", exact: true }).click();
    await page.getByText("Room created successfully!").waitFor();
    await page.getByText(roomName).waitFor();
  });

  const room = await step("room row exists in the database", async () => {
    const rows = db(`select r.id, r.name, r.is_public, u.email as owner from chat_rooms r join users u on u.id = r.user_owned_chat_rooms where r.name = '${roomName}'`);
    expect(rows.length === 1, `one chat_rooms row named ${roomName}, got ${rows.length}`);
    expect(rows[0].owner === user.email, `room owner is ${rows[0].owner}`);
    return rows[0];
  });

  await step("enter the room", async () => {
    await page.getByText(roomName).click();
    await page.waitForURL(`**/chat/rooms/${room.id}`);
    await page.getByPlaceholder("Type a message...").waitFor();
  });

  await step("send a message over the websocket", async () => {
    const input = page.getByPlaceholder("Type a message...");
    await input.fill(body);
    await input.press("Enter");
    await page.getByText(body).waitFor();
  });

  await step("message persisted with the sender name", async () => {
    const rows = db(`select body, sender_name from chat_messages where body = '${body}'`);
    expect(rows.length === 1, `one chat_messages row for the body, got ${rows.length}`);
    expect(rows[0].sender_name === user.name, `sender_name is ${rows[0].sender_name}`);
    log(JSON.stringify(rows[0]));
  });
}
