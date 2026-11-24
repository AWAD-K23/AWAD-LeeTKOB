import { expect, test } from "@playwright/test";

// Finding 2 of the PA3 Member 3 acceptance: a TA who had already loaded the lab
// list saw "Lab not available" right after creating a lab, because the new lab's
// page first rendered the cached list. The follow-up list request is slowed here
// so the stale render, if any, lasts long enough to be observed every time.
test("a TA who loaded the lab list never sees 'Lab not available' after creating a lab", async ({
  page,
}) => {
  await page.addInitScript(() => {
    const seen = () => {
      if (document.body?.innerText.includes("Lab not available")) {
        (window as unknown as { sawUnavailable: boolean }).sawUnavailable =
          true;
      }
    };
    new MutationObserver(seen).observe(document, {
      childList: true,
      subtree: true,
      characterData: true,
    });
  });
  await page.route("**/api/me", (route) =>
    route.fulfill({ json: { id: "ta-a", role: "ta" } }),
  );
  await page.route("**/api/classes", (route) =>
    route.fulfill({ json: [{ id: "class-a", name: "Sorting A" }] }),
  );
  let created = false;
  await page.route("**/api/labs", async (route) => {
    if (!created) {
      await route.fulfill({ json: [] });
      return;
    }
    await new Promise((resolve) => setTimeout(resolve, 800));
    await route.fulfill({
      json: [
        {
          id: "lab-a",
          classId: "class-a",
          title: "Merge sort",
          statement: "Sort ascending using merge sort.",
        },
      ],
    });
  });
  await page.route("**/api/classes/class-a/labs", async (route) => {
    created = true;
    await route.fulfill({ status: 201, json: { id: "lab-a" } });
  });
  await page.route("**/api/labs/lab-a/submissions", (route) =>
    route.fulfill({ json: [] }),
  );

  // The list is loaded and cached before the lab exists; the form is opened
  // with a client-side link so the cache survives.
  const listLoaded = page.waitForResponse("**/api/labs");
  await page.goto("/labs");
  await listLoaded;
  await page.getByRole("link", { name: "Create lab" }).first().click();
  await expect(page).toHaveURL(/\/labs\/new$/);
  await page.getByLabel("Class", { exact: true }).selectOption("class-a");
  await page.getByLabel("Title", { exact: true }).fill("Merge sort");
  await page
    .getByLabel("Statement", { exact: true })
    .fill("Sort ascending using merge sort.");
  await page.getByLabel("Input", { exact: true }).fill("3 1");
  await page.getByLabel("Expected output", { exact: true }).fill("1 3");
  await page.getByLabel("Rule", { exact: true }).fill("Use merge sort");
  await page.getByLabel("Deduction", { exact: true }).fill("2");
  await page.getByRole("button", { name: "Create lab", exact: true }).click();

  await expect(page).toHaveURL(/\/labs\/lab-a$/);
  await expect(page.getByRole("alert")).toContainText("Lab created");
  expect(
    await page.evaluate(
      () => (window as unknown as { sawUnavailable?: boolean }).sawUnavailable,
    ),
  ).toBeFalsy();
});
