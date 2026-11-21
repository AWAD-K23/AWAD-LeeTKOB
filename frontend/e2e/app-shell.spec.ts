// Check the signed-out entry point and its accessible login fields.
import { expect, test } from "@playwright/test";

test("home page shows the LeeTKOB app shell", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "LeeTKOB" })).toBeVisible();
  await expect(
    page.getByRole("heading", { name: "Welcome back" }),
  ).toBeVisible();
  await expect(page.getByLabel("Email")).toBeVisible();
  await expect(page.getByLabel("Password")).toHaveAttribute("type", "password");
});
