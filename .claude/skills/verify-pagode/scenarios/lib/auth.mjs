// Shared user-path helpers for scenarios. Every helper drives the real UI.

// Registers a fresh user through /user/register and leaves the page on /dashboard.
export async function registerViaUI({ page, step, unique }) {
  const user = {
    name: "Verify User",
    email: `${unique("verify")}@example.com`,
    password: "Verify-pass-123",
  };
  await step("register a new user", async () => {
    await page.goto("/user/register");
    await page.getByRole("heading", { name: "Create an account" }).waitFor();
    await page.fill("#name", user.name);
    await page.fill("#email", user.email);
    await page.fill("#password", user.password);
    await page.fill("#password-confirm", user.password);
    await page.getByRole("button", { name: "Create Account" }).click();
    await page.waitForURL("**/dashboard");
    await page.getByText("Your account has been created. You are now logged in.").waitFor();
  });
  return user;
}

// Logs in through /user/login and leaves the page on /dashboard.
export async function loginViaUI({ page, step }, { email, password, name }) {
  await step("log in", async () => {
    await page.goto("/user/login");
    await page.getByRole("heading", { name: "Log in to your account" }).waitFor();
    await page.fill("#email", email);
    await page.fill("#password", password);
    await page.getByRole("button", { name: "Log in" }).click();
    await page.waitForURL("**/dashboard");
    if (name) await page.getByText(`Welcome back, ${name}. You are now logged in.`).waitFor();
  });
}
