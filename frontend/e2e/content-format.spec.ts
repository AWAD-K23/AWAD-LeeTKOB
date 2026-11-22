import { expect, test } from "@playwright/test";

test("Markdown statements and colored C++ preserve editable source", async ({
  page,
}) => {
  const statement = [
    "## Input format",
    "",
    "Read **two integers**.",
    "",
    "| Input | Output |",
    "| --- | --- |",
    "| 3 1 | 1 3 |",
    "",
    "```cpp",
    "int main() { return 0; }",
    "```",
    "",
    "<script>window.markdownExecuted = true</script>",
  ].join("\n");
  await page.route("**/api/me", (route) =>
    route.fulfill({ json: { id: "student", role: "student" } }),
  );
  await page.route("**/api/labs", (route) =>
    route.fulfill({ json: [{ id: "lab", title: "Numbers", statement }] }),
  );
  await page.goto("/labs/lab");
  await expect(
    page.getByRole("heading", { name: "Input format" }),
  ).toBeVisible();
  await expect(page.locator(".markdown-content strong")).toHaveText(
    "two integers",
  );
  await expect(page.locator(".markdown-content table")).toBeVisible();
  await expect(page.locator(".markdown-content script")).toHaveCount(0);
  await expect(page.locator(".markdown-content .hljs-type").first()).toHaveText(
    "int",
  );

  const source =
    '#include <iostream>\nint main() {\n    // Keep this comment.\n    std::cout << "<script>";\n}\n';
  await page.getByLabel("C++ source", { exact: true }).fill(source);
  await expect(page.getByLabel("C++ source", { exact: true })).toHaveValue(
    source,
  );
  await expect(page.locator(".code-colors .hljs-comment")).toContainText(
    "Keep this comment.",
  );
  await expect(page.locator(".code-colors .hljs-string")).toHaveText([
    "<iostream>",
    '"<script>"',
  ]);
  await expect(page.locator(".code-colors script")).toHaveCount(0);
});
