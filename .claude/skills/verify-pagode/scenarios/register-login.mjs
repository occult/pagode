// Feature: auth (features/auth.md). Register -> dashboard -> log out -> log in.
export default async function ({ page, step, db, expect, unique, log }) {
  const name = "Verify User";
  const email = `${unique("verify")}@example.com`;
  const password = "Verify-pass-123";

  await step("open register page", async () => {
    await page.goto("/user/register");
    await page.getByRole("heading", { name: "Create an account" }).waitFor();
  });

  await step("fill registration form", async () => {
    await page.fill("#name", name);
    await page.fill("#email", email);
    await page.fill("#password", password);
    await page.fill("#password-confirm", password);
  });

  await step("submit registration and land on dashboard", async () => {
    await page.getByRole("button", { name: "Create Account" }).click();
    await page.waitForURL("**/dashboard");
    await page.getByText("Your account has been created. You are now logged in.").waitFor();
  });

  await step("user row exists in the database", async () => {
    const rows = db(`select id, name, email, verified, admin from users where email = '${email}'`);
    expect(rows.length === 1, `one users row for ${email}, got ${rows.length}`);
    expect(rows[0].name === name, `stored name is ${rows[0].name}`);
    expect(rows[0].admin === 0, "self-registered user is not admin");
    log(`user id ${rows[0].id}`);
  });

  await step("log out from the user menu", async () => {
    await page.getByRole("button", { name: new RegExp(name) }).first().click();
    await page.getByRole("menuitem", { name: "Log out" }).click();
    await page.waitForURL((u) => u.pathname === "/");
    await page.getByText("You have been logged out successfully.").waitFor();
  });

  await step("dashboard is closed to a logged-out visitor", async () => {
    const res = await page.goto("/dashboard");
    expect(res.status() === 401, `GET /dashboard while logged out returns ${res.status()}`);
  });

  await step("log in with the new credentials", async () => {
    await page.goto("/user/login");
    await page.fill("#email", email);
    await page.fill("#password", password);
    await page.getByRole("button", { name: "Log in" }).click();
    await page.waitForURL("**/dashboard");
    await page.getByText(`Welcome back, ${name}. You are now logged in.`).waitFor();
  });

  await step("sidebar shows the signed-in user", async () => {
    await page.getByRole("button", { name: new RegExp(name) }).first().waitFor();
  });
}
