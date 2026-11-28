import { expect, test } from "@playwright/test";

// FAQ-02 keeps only each student's latest submission in the queue; FAQ-10 lets
// the TA open every submission. Earlier ones are listed read-only, not lost.
test("TA review shows earlier submissions read-only, outside the queue", async ({
  page,
}) => {
  await page.route("**/api/labs/lab-a/publication", (route) =>
    route.fulfill({
      json: {
        reviewRevision: 1,
        required: 1,
        reviewed: 0,
        pending: 1,
        incompleteSubmissions: 0,
        published: false,
      },
    }),
  );
  await page.route("**/api/me", (route) =>
    route.fulfill({ json: { id: "ta-a", role: "ta" } }),
  );
  await page.route("**/api/labs/lab-a/review-queue", (route) =>
    route.fulfill({
      json: [
        {
          submissionId: "sub-latest",
          studentId: "student-a",
          createdAt: "2026-11-24 14:00:00+00",
          judgeVerdict: "Accepted",
          sourceCode: "int main() {\n  bubble();\n}",
          constraintId: "c2",
          constraint: "Use merge sort",
          deduction: 4,
          verdict: "violation",
          reason: "",
          evidence: [{ line: 2, text: "  bubble();" }],
          source: "model",
          decision: "",
          reviewEvidence: [],
          appliedDeduction: 0,
        },
      ],
    }),
  );
  await page.route("**/api/labs/lab-a/earlier-submissions", (route) =>
    route.fulfill({
      json: [
        {
          submissionId: "sub-earlier",
          studentId: "student-a",
          createdAt: "2026-11-24 13:00:00+00",
          judgeVerdict: "Accepted",
          sourceCode: "int main() {\n  mergeSort();\n}",
          findings: [
            {
              constraint: "Use merge sort",
              verdict: "pass",
              source: "model",
              reason: "",
              evidence: [{ line: 2, text: "  mergeSort();" }],
            },
          ],
        },
      ],
    }),
  );

  await page.goto("/labs/lab-a/review");
  await expect(
    page.getByRole("heading", { name: "Submission sub-latest" }),
  ).toBeVisible();

  const earlier = page.getByRole("region", {
    name: "Earlier submissions (not counted)",
  });
  await expect(earlier).toBeVisible();
  await earlier.getByText("Submission sub-earlier").click();
  await expect(earlier.getByText("mergeSort();").first()).toBeVisible();
  await expect(earlier).toContainText("Use merge sort");
  await expect(earlier).toContainText("pass (model)");
  // Read-only: no decision controls outside the queue.
  await expect(earlier.getByRole("button")).toHaveCount(0);
});
