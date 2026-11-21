import { expect, test } from "@playwright/test";

test("TA creates a lab with tests and chosen deductions", async ({ page }) => {
  await page.route("**/me", (route) => route.fulfill({ json: { role: "ta" } }));
  await page.route("**/classes", (route) =>
    route.fulfill({ json: [{ id: "class-a", name: "Sorting A" }] }),
  );
  await page.route("**/labs", (route) => route.fulfill({ json: [] }));
  let saved: Record<string, unknown> | undefined;
  await page.route("**/classes/class-a/labs", async (route) => {
    saved = route.request().postDataJSON();
    await route.fulfill({ status: 201, json: { id: "lab-a" } });
  });
  await page.goto("/");
  await page.getByLabel("Class", { exact: true }).selectOption("class-a");
  await page.getByLabel("Title", { exact: true }).fill("Merge sort");
  await page
    .getByLabel("Statement", { exact: true })
    .fill("Sort ascending using merge sort.");
  await page.getByLabel("Input", { exact: true }).fill("3 1");
  await page.getByLabel("Expected output", { exact: true }).fill("1 3");
  await page
    .getByLabel("Rule", { exact: true })
    .fill("Do not call `std::sort`");
  await page.getByLabel("Deduction", { exact: true }).fill("3");
  await page.getByRole("button", { name: "Create lab", exact: true }).click();
  await expect(page.getByRole("alert")).toContainText("Lab created");
  expect(saved).toMatchObject({
    title: "Merge sort",
    tests: [{ input: "3 1", expectedOutput: "1 3" }],
    constraints: [{ description: "Do not call `std::sort`", deduction: 3 }],
  });
});

test("student submits code and sees the polled judge result", async ({
  page,
}) => {
  await page.route("**/me", (route) =>
    route.fulfill({ json: { role: "student" } }),
  );
  await page.route("**/classes", (route) => route.fulfill({ json: [] }));
  await page.route("**/labs", (route) =>
    route.fulfill({
      json: [{ id: "lab-a", title: "Sorting", statement: "Sort two numbers." }],
    }),
  );
  await page.route("**/labs/lab-a/submissions", async (route) => {
    expect(route.request().headers()["idempotency-key"]).toBeTruthy();
    expect(route.request().postDataJSON()).toEqual({ source: "int main(){}" });
    await route.fulfill({ status: 201, json: { id: "submission-a" } });
  });
  await page.route("**/submissions/submission-a", (route) =>
    route.fulfill({
      json: {
        id: "submission-a",
        status: "completed",
        verdict: "Accepted",
        detail: "All test cases passed.",
      },
    }),
  );
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Create a lab" })).toHaveCount(
    0,
  );
  await page.getByLabel("Lab", { exact: true }).selectOption("lab-a");
  await page.getByLabel("C++ source").fill("int main(){}");
  await page.getByRole("button", { name: "Submit", exact: true }).click();
  await expect(page.getByRole("status")).toContainText("Accepted");
  await expect(page.getByRole("status")).not.toContainText("deduction");
});
