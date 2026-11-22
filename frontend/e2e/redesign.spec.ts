import { expect, test } from "@playwright/test";

test("page navigation animates and honors reduced motion", async ({ page }) => {
  await page.emulateMedia({ reducedMotion: "no-preference" });
  await page.route("**/api/me", (route) =>
    route.fulfill({ json: { id: "admin", role: "admin" } }),
  );
  await page.route("**/api/classes", (route) => route.fulfill({ json: [] }));
  await page.route("**/api/people?*", (route) => route.fulfill({ json: [] }));
  await page.goto("/classes");
  await page
    .getByRole("navigation", { name: "Main navigation" })
    .getByRole("link", { name: "People", exact: true })
    .click();
  await expect(
    page.getByRole("heading", { name: "People", exact: true }),
  ).toBeVisible();

  // The destination retains an entrance fallback even without native snapshots.
  await expect(page.locator("main#content")).toHaveCSS(
    "animation-name",
    "page-reveal",
  );
  await expect(page.locator("main#content")).toHaveCSS(
    "animation-duration",
    "0.42s",
  );

  await page.emulateMedia({ reducedMotion: "reduce" });
  await page
    .getByRole("navigation", { name: "Main navigation" })
    .getByRole("link", { name: "Classes", exact: true })
    .click();
  await expect(
    page.getByRole("heading", { name: "Classes", exact: true }),
  ).toBeVisible();
  await expect(page.locator("main#content")).toHaveCSS(
    "animation-name",
    "none",
  );
});

test("sign in resolves the session and sign out clears the workspace", async ({
  page,
}) => {
  let signedIn = false;
  await page.route("**/api/me", (route) =>
    signedIn
      ? route.fulfill({ json: { id: "admin", role: "admin" } })
      : route.fulfill({ status: 401, json: { detail: "Sign in first." } }),
  );
  await page.route("**/api/auth/login", (route) => {
    expect(route.request().postDataJSON()).toEqual({
      email: "admin@test.example",
      password: "test-password",
    });
    signedIn = true;
    return route.fulfill({ status: 204 });
  });
  await page.route("**/api/auth/logout", (route) => {
    signedIn = false;
    return route.fulfill({ status: 204 });
  });
  await page.route("**/api/classes", (route) => route.fulfill({ json: [] }));
  await page.goto("/");
  await page.getByLabel("Email", { exact: true }).fill("admin@test.example");
  await page.getByLabel("Password", { exact: true }).fill("test-password");
  await page.getByRole("button", { name: "Show password" }).click();
  await expect(page.getByLabel("Password", { exact: true })).toHaveAttribute(
    "type",
    "text",
  );
  await page.getByRole("button", { name: "Sign in", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "Classes", exact: true }),
  ).toBeVisible();
  await page.getByRole("button", { name: "Sign out", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "Welcome back" }),
  ).toBeVisible();
});

test("admin selects people by email and opens the saved class", async ({
  page,
}) => {
  await page.route("**/api/me", (route) =>
    route.fulfill({ json: { id: "admin", role: "admin" } }),
  );
  await page.route("**/api/classes", async (route) => {
    if (route.request().method() === "POST") {
      expect(route.request().postDataJSON()).toEqual({
        name: "Sorting A",
        taId: "ta",
        studentIds: ["student"],
      });
      return route.fulfill({ status: 201, json: { id: "class" } });
    }

    return route.fulfill({ json: [] });
  });
  const members = [
    { id: "ta", email: "ta@test.example", role: "ta" },
    { id: "student", email: "student@test.example", role: "student" },
  ];
  await page.route("**/api/people?*", (route) =>
    route.fulfill({ json: members }),
  );
  await page.route("**/api/classes/class", (route) =>
    route.fulfill({ json: { id: "class", name: "Sorting A", members } }),
  );
  await page.goto("/classes");
  await page.getByRole("link", { name: "Create class" }).click();
  await page.getByLabel("Class name").fill("Sorting A");
  await page.getByRole("radio", { name: "ta@test.example" }).check();
  await page.getByRole("checkbox", { name: "student@test.example" }).check();
  await page.getByRole("button", { name: "Create class" }).click();
  await expect(page.getByRole("heading", { name: "Sorting A" })).toBeVisible();
  await expect(
    page.getByText("Ready for your teaching assistant"),
  ).toBeVisible();
  await page.reload();
  await expect(page.getByRole("heading", { name: "Sorting A" })).toBeVisible();
});

test("connection failure is visible and retry restores the session", async ({
  page,
}) => {
  let unavailable = true;
  await page.route("**/api/me", (route) =>
    unavailable
      ? route.fulfill({
          status: 502,
          body: "<html>Unavailable</html>",
          contentType: "text/html",
        })
      : route.fulfill({ status: 401, json: { detail: "Sign in first." } }),
  );
  await page.goto("/");
  await expect(page.getByRole("alert")).toContainText("502");
  unavailable = false;
  await page.getByRole("button", { name: "Retry" }).click();
  await expect(
    page.getByRole("heading", { name: "Welcome back" }),
  ).toBeVisible();
});

test("mobile workspace preserves a draft across navigation with reduced motion", async ({
  page,
}) => {
  await page.setViewportSize({ width: 375, height: 812 });
  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.route("**/api/me", (route) =>
    route.fulfill({ json: { id: "student", role: "student" } }),
  );
  await page.route("**/api/classes", (route) => route.fulfill({ json: [] }));
  await page.route("**/api/labs", (route) =>
    route.fulfill({
      json: [{ id: "lab", title: "Sorting", statement: "Sort ascending." }],
    }),
  );
  await page.goto("/labs/lab");
  await page.getByRole("tab", { name: "code", exact: true }).click();
  await page
    .getByLabel("C++ source", { exact: true })
    .fill("int main() { return 0; }");
  await page.reload();
  await page.getByRole("tab", { name: "code", exact: true }).click();
  await expect(page.getByLabel("C++ source", { exact: true })).toHaveValue(
    "int main() { return 0; }",
  );
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= window.innerWidth,
    ),
  ).toBeTruthy();
});
