// Check the signed-out entry point and its accessible login fields.
import { expect, test } from "@playwright/test";

test("home page shows the LeeTKOB app shell", async ({ page }) => {
  await page.route("**/api/me", (route) =>
    route.fulfill({
      status: 401,
      json: { detail: "Sign in first." },
    }),
  );
  await page.goto("/");
  await expect(page.getByRole("link", { name: "LeeTKOB home" })).toBeVisible();
  await expect(
    page.getByRole("heading", { name: "Welcome back" }),
  ).toBeVisible();
  await expect(page.getByLabel("Email")).toBeVisible();
  await expect(page.getByLabel("Password", { exact: true })).toHaveAttribute(
    "type",
    "password",
  );
});
