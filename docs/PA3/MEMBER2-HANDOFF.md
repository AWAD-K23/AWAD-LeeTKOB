# Member 2 handoff

## Delivered parts

| Part | Result | Evidence |
|---|---|---|
| Real model | Gemini checks pending semantic constraints after Accepted | `eval/runs/real-worker/result.json` |
| Guardrails | Invalid schema/evidence and timeouts become unsure | Fake tests in `backend/internal/semantic/` |
| Usage and cost | Provider IDs, tokens, duration, prompt, response, cost estimate | `eval/runs/real-dev/calls.jsonl` |
| Lab budget | Persistent database reservations with concurrent-worker protection | `worker/model.go`, migration 00005, integration tests |
| Request pacing | Five seconds between request starts per process | `semantic/rate.go`, real usage summary |
| Checker | Generic command and standalone offline image | [Evaluation commands](../../eval/README.md) |
| Data | Unchanged official snapshot plus five separate Accepted cases | Official hashes and `eval/extra/` |

## Real evaluation

| Metric | Result |
|---|---:|
| Pairs | 102 |
| Wrongful deductions | 0 |
| Violations caught | 21/37, or 56.8% |
| Unsure | 31/102, or 30.4% |
| Original evaluator exit | 1: valid predictions, score thresholds missed |

- Files: [predictions](../../eval/runs/real-dev/predictions.jsonl),
  [original report](../../eval/runs/real-dev/report.txt),
  [metrics](../../eval/runs/real-dev/summary.json), and
  [usage summary](../../eval/runs/real-dev/usage-summary.json).
- The run made 34 requests. There were 29 successful responses and five HTTP 503
  responses. Errors were not retried. Unknown usage retained its full reservation;
  budget blocks left later constraints unsure.
- Strict evidence validation rejected several responses whose line/text did not
  match the sanitized code. A matching source line can still be weak proof:
  one predicted violation cited real lines outside the labelled evidence lines
  and was correctly not counted as caught by the original evaluator.
- Peak request rate was 12 per rolling minute, below the user's 15 RPM quota.
- PA3 requires honest reporting. These results do not meet the PA4 thresholds.
- No labels, evaluator rules, or tests were loosened to improve the score.

## Measured cost

- Model: `gemini-3.5-flash-lite`.
- On 29 calls with known usage: 22,841 input tokens, 2,536 output tokens,
  and zero reported thinking tokens.
- Standard rates: US$0.30 input and US$2.50 output per million tokens.
- Sum of rounded per-call estimates: **US$0.013204**.
- For 80 similar calls, assuming one call per submission:
  `0.013204 / 29 × 80 = US$0.036425`.
- This excludes unknown usage on the five HTTP 503 calls. Their retained
  reservations total US$2.392065; that is a conservative budget charge,
  not a measured provider bill. Do not report it as actual token usage.
- The separate checker smoke and worker smoke cost an estimated US$0.000220
  and US$0.000444 respectively. Across all 36 real requests, peak RPM remained 12.

## Worker proof

- Submission: `e7ad52ed-686a-4d50-abc7-b426fb8d1931` in the local manual-test lab.
- A labelled SQL fixture queued one submission for the enrolled test student;
  this was a backend smoke test, not a claim of a new frontend test.
- `go -C backend run ./cmd/worker --once` ran the real sandbox tests and exited
  after one job. The result is Accepted.
- The sort ban is `pass` / `hard_rule`; the comparison requirement is
  `pass` / `model`, with the exact comparison at line 8.
- Known identity text in both a comment and a literal was removed from the
  provider prompt. The original evidence line is retained for TA review.
- [Saved result](../../eval/runs/real-worker/result.json) includes the provider
  request ID and audit record. [Read-only SQL](../../eval/runs/real-worker/query.sql)
  reproduces the database inspection while that local fixture exists.
- No marks were published and no deductions were applied.

## Interface for Member 3

```json
{
  "constraintId": "constraint-uuid",
  "verdict": "pass",
  "evidence": [{"line": 8, "text": "    if (a > b) { int temp = a; a = b; b = temp; }"}],
  "source": "model"
}
```

| Value | Meaning |
|---|---|
| `hard_rule` | Syntax-based ban decision, or unsure after analysis failure |
| `semantic_pending` | Semantic checking disabled |
| `model` | Validated model finding, or unsure after a failed/blocked attempt |

- Read findings from `constraint_verdicts` through a TA/class-protected endpoint.
- Keep judge verdicts in `judge_results` separate from constraint findings.
- TA-selected deductions come from `lab_constraints`, never from model output.
- Student status endpoints still exclude findings and deductions.
- Member 3 owns the review endpoint/UI and publication controls.
- Provider logs in `model_calls` are protected data; do not expose them to students.
- Member 3 acceptance and public deployment verification remain their tasks.

## Limits to carry forward

- Improve evidence accuracy and semantic coverage before PA4; keep this real run.
- Unknown personal names cannot be inferred reliably from the current account
  schema. Comments and known UUID/email identifiers are removed.
- Rate limiting is shared within one process. Multiple workers or other tools
  using the same key need coordinated provider-level pacing.
- Budgets use reviewed standard rates and conservative reservations; update both
  rates and reservation bounds if the provider changes them.
- Verify the final public project/key spending cap before public testing.
- The worker, API, and checker need restarting/rebuilding after code changes.
- Offline image, fake guardrails, real checker smoke, and real worker smoke have
  local evidence. Remote CI and teammate acceptance are not claimed here.
