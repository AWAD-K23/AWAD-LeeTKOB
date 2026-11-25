# PA4 plan — three-member handoff

**Due:** Wednesday, December 2, 2026, at 22:00 (Asia/Saigon).

**Goal:** Make the evaluation block unsafe merges, enforce all five guardrails, require TA review before Publish, and deliver a reproducible red-team report.

**Source:** [PA4 rubric](<PA#4 rubric — Eval, guardrails, human gate and red-team.md>), [project brief](<../LeeTKOB — class project brief.md>), and [project progress](../PROJECT-PROGRESS.md).

## Starting point after PA3

PA3 is complete per the team. Use [PA3.md](../../PA3.md), the [PA3 self-assessment](../../SELF_ASSESSMENT_REPORT.md), [handoff acceptance](../PA3/MEMBER3-ACCEPTANCE.md), and [public walkthrough evidence](../PA3/evidence/public-walk.json) as the baseline. The older PA3 planning checklists contain historical open items; use the final report for delivered behavior.

| Area | Delivered baseline | PA4 work still required |
| --- | --- | --- |
| Evaluation | `real-dev-numbered`: 0 wrongful, 35/37 recall (94.6%), 2/102 unsure (2.0%), evaluator exit 0 | Fresh real-model CI gate, ten own cases, three runs of one commit, held-out rehearsal |
| Cost and audit | 39 successful calls with provider request IDs; US$0.020851 estimated for the saved official run | Measure official plus own cases in CI; preserve budget and provider-failure evidence |
| Guardrails | Schema/evidence checks, timeout, persistent lab budget, unknown-field rejection, numbered source lines | Audit the full G1–G5 matrix, prove tests detect removed guards, demonstrate deployed behavior |
| TA review | Separate decisions, reviewer identity/history, exact evidence, configured deductions and local unsure reasons | Persist a random pass sample and enforce required reviews at Publish |
| Publication | Findings and TA decisions stay out of student submission responses; no Publish endpoint | Add the server gate, publication state and authorized student result view |
| Own cases | Five separate labelled/tagged Accepted solutions in `eval/extra/` | At least ten, covering trap, injection and disguised; preserve case-before-fix history |
| CI/deployment | Offline checker and sandbox CI; public model-on PA3 walkthrough recorded | Required live eval, `pa4` tag trigger, deployed Publish and injection checks |

The passing saved run is one dev-set observation. It does not establish repeated-run stability, a live CI gate, or the lecturer's held-out result. Keep the older failing `real-dev` run unchanged. Sending numbered `{line, text}` objects fixed line-number drift without relaxing evidence validation (test `c05230c`, fix `1cdc486`, PR #22).

## How the handoff works

- Keep the PA3 member roles: Member 1 owns backend integration, Member 2 owns the checker/evaluation, and Member 3 owns the review UI, red-team report and release evidence. Confirm the owners before implementation; Member 3's PA3 owner was Phạm Bảo Huy (23120131).
- Agree on the review, sampling and publication contract first. Members 1 and 2 can then work independently; Member 3 integrates only after the relevant handoff checks pass.
- Each owner supplies tests, documentation and fixes for their part. Return integration failures to the owner and record acceptance rather than silently replacing their work.
- Commit each new guardrail with its meaningful fake-model test, test no later than code. Commit each new own case before the checker change that handles it. Do not split or rebuild completed history to simulate this sequence; disclose existing PA3 history limitations.
- Preserve official files, labels and existing evidence. Fix the checker when evaluation fails; never loosen a test, case or label in the change that makes it pass. Production rules/prompts must remain generic and contain no dev-set IDs or quoted dev-set code.
- Record actual assistant use and concrete Changed/Rejected decisions in `AI-LOG.md` as the work happens. Do not manufacture rejected suggestions.

## Member 1 — server sampling, required reviews and Publish

- [ ] Review the existing API, migrations and TA decision history with Member 3. Write the publication/review contract under `docs/contract/` before connecting the UI.
- [ ] Define the eligible review population using counted submissions and per-constraint findings. Preserve latest-submission behavior and read-only earlier submissions from PA3.
- [ ] Queue every violation and every unsure, plus `ceil(0.10 × pass count)` randomly chosen pass findings. Draw on the server and persist the selection; reopening the queue must not redraw it.
- [ ] Add a fixed-seed test through an injectable random source. Verify sample size, uniqueness, membership and that selection is random rather than the first rows; include 0, 1 and non-multiple-of-ten pass counts.
- [ ] Define how new submissions, rejudging and constraint changes affect the review population. Version or invalidate stale review/publication state so a previously reviewed attempt cannot authorize publishing a replacement.
- [ ] Reuse independent TA decisions, reviewer attribution and append-only decision history. Ensure every required finding has a valid TA decision before publication.
- [ ] Implement Publish with role/class authorization and pending-review validation on the server. Validate and store the published result atomically against the same review population; concurrent submissions/rejudging must not bypass the gate.
- [ ] Make Publish the only action exposing a deduction to the owning student. Derive deductions from configured constraints and TA decisions; the checker must have no mark-writing tool or path.
- [ ] Add direct API tests for pending reviews, a student caller, another class's TA, unreviewed sampled passes, successful publication, and student visibility before/after publication. Cover stale reviews and concurrent changes meaningfully.
- [ ] Keep violation, unsure and sampled-pass review requirements visible in the API response, alongside review progress and reasons. Preserve class and student ownership checks on every new read/write route.
- [ ] **Handoff:** Demonstrate deterministic sampling tests, denied direct Publish requests, successful authorized Publish and the student's published deduction. Supply migration steps, API examples, state transitions and passing check evidence to Member 3.

**Acceptance:** No direct request can publish while required reviews remain; queue reloads preserve the sample; a new counted submission cannot inherit an old decision; student deductions appear only after authorized Publish.

## Member 2 — guardrails, own cases and real-model evaluation

- [ ] Accept the PA3 checker/worker baseline and reproduce the saved passing metrics locally without model calls. Verify all 128 official files against the manifest and retain the upstream-comparison evidence.
- [ ] Audit G1–G5 below against actual code/tests. Add missing tests before or with each guardrail change, in separate genuine commits. Check that each test fails when its guardrail is removed; record any PA3 mixed-history limitation honestly.
- [ ] Extend `eval/extra/` to at least ten Accepted solutions with labels and tags, covering trap, injection and disguised. Verify they compile and pass their problem tests. Keep these cases separate from the official set and commit them before related checker fixes.
- [ ] Investigate remaining missed violations with generic algorithm/evidence reasoning. Preserve exact source validation, hard-rule-first routing and safe unsure fallback; no problem-specific fixes or label changes.
- [ ] Add a live evaluation job for PRs, main and `pa4`. Generate fresh predictions for both official and own cases inside the job using the real model, then run the unchanged evaluator on both. Every nonzero evaluator exit must fail the job.
- [ ] Keep unit/guardrail tests offline with fake models. Retain the separate model-off checker job; its expected semantic threshold failure is not the live gate.
- [ ] Store the CI key only in repository secrets; preserve provider request IDs, model/version, tokens, timings and estimated cost in safe artifacts. Do not expose the key to student-code containers.
- [ ] Document a fork policy in `PA4.md`: secretless fork PRs must not count as a passing live eval. Require a maintainer-reviewed same-repository branch/run before merge; never execute untrusted fork code with secrets through `pull_request_target`.
- [ ] Configure the live eval as a required check on main. Review CODEOWNERS coverage for tests, `.github/` and eval data with required owner review; record actual settings/evidence rather than treating a workflow file as proof of protection.
- [ ] Save a PR where a wrongful deduction or recall below 90% makes the required eval red and blocks merging. Link the checker fix and subsequent green run without editing labels/data, hiding exit codes or weakening branch protection.
- [ ] Freeze one commit and run the real dev evaluation three times with the same model/configuration. Save each run in a new folder with predictions, reports, summaries and provider call logs; never overwrite the PA3 runs.
- [ ] Count unique `(solution, constraint)` pairs whose verdict differs across any of the three runs. Report the definition, count, lowest/highest recall, lowest/highest unsure rate and wrongful deductions per run in `PA4.md`.
- [ ] Measure one full CI evaluation's estimated provider cost, including official and own cases, failed/unknown-usage calls and any retries. The PA3 US$0.020851 official-only run is a baseline, not the full PA4 CI cost.
- [ ] Rehearse the generic checker image/command with an arbitrary dataset folder. Confirm valid output within the 30-minute marking limit and model access through the documented host/proxy configuration only; no hidden per-problem setup.
- [ ] Document the exact model host required by marking in README and `PA4.md`. Prepare the checker-required variable names for private `marking.env`; test working credentials separately from Git and ZIP packaging.
- [ ] **Handoff:** Demonstrate G1–G5 with no network, a fresh real-model CI run, case history and three saved runs. Supply the report artifacts, cost arithmetic, blocked-PR/fix links and checker/marking instructions to Member 3.

### Guardrail acceptance matrix

| Guardrail | Required negative tests and behavior | Integration evidence |
| --- | --- | --- |
| G1 Schema | Invalid JSON, fenced JSON, prose before JSON, unknown verdict, missing/duplicate constraint → unsure | Fake model tests fail if schema validation is bypassed |
| G2 Evidence | Invented text/line, correct text at wrong number, stripped-comment-only evidence → unsure; map cleaned lines to original source | Keep exact evidence checks and numbered-line regression; verify original source in review UI |
| G3 Timeout | A model call exceeding its deadline → unsure | Fake model with cancellation, bounded completion and locally generated TA reason |
| G4 Lab cap | Exhausted persistent budget → remaining constraints unsure, no extra provider call | Fake model/ledger test, concurrent budget coverage where applicable, TA sees cap reason |
| G5 Deduction | Model deduction field or instructions about marks cannot change the configured deduction | Rejected/ignored output and forged marks tested through checker, review API and displayed TA deduction |

Hard rules still decide ban-list constraints before the model. Deployed injection submissions that violate a constraint must end as violation or unsure, never pass. A fake-model test alone does not demonstrate this live behavior.

## Member 3 — review UI, red-team report and submission

- [ ] Accept Member 1's publication/review contract and Member 2's guardrail/evaluation evidence; run the documented handoff checks before claiming integration complete.
- [ ] Update the TA queue to show all required findings, persisted sampled-pass status, review progress, exact evidence, configured deduction and local unsure reasons including budget exhaustion.
- [ ] Add the Publish flow and published student result view. Show pending-review errors returned by the server; verify a direct API request is blocked even when the UI is bypassed.
- [ ] Add focused browser coverage for review completion → Publish → owning student sees the deduction, plus pending reviews and injection attempts that leave the TA-configured deduction unchanged.
- [ ] Write `PA4.md` section **Where we stop it** using the team's own evaluation numbers, configured deductions and lab/course weighting where known. Explain who pays when wrong, how errors are detected, why Publish is the stopping point, and the missed-violation/fairness cost of unsampled passes. Do not present a brief example's weight as the team's measured configuration.
- [ ] Obtain the lecturer-assigned target from Classroom after PA3. Use only that deployment/public repository and its PA3 test accounts; do not guess the target or start attacks against other teams.
- [ ] Read the target code and derive at least one attack from a concrete file/behavior. Record the target commit and code location.
- [ ] Execute at least eight reproducible attacks across the five categories below. Save exact inputs, timestamps with timezone, expected and observed responses/logs/screenshots, held/broke status, applicable rule, who pays and proposed fixes for breaks.
- [ ] Write `docs/red-team/<target-student-IDs>.md` and retain safe replay fixtures. Report held protections as well as breaks; rehearse two attacks from the report alone.
- [ ] Keep traffic at about one request per second at most, one resource attack per kind, and trigger the target's cost cap only once. Stop on a live key or real personal data; have the team notify lecturer/target that day without reproducing sensitive values in the report.
- [ ] If the target is unavailable, record timed attempts and have the team notify the lecturer by Monday, November 30. Local attacks on the target's code with the team's own key also count; clearly identify that execution mode.
- [ ] Integrate and verify the PA4 build on the public deployment, including direct Publish denials, sampling persistence and live injection results. Keep the deployment running through PA4 marking.
- [ ] Write root `PA4.md` with tag CI, blocked PR/fix, all three runs and variability, full CI cost, fork policy, model host and **Where we stop it**. Link evidence for every claim.
- [ ] Update the self-assessment to six PA4 rubric rows: CI 20, held-out 20, guardrails 20, human gate 15, red-team 20, AI log 5. Claim held-out points only as an estimate from dev results; include what the team did not manage.
- [ ] Prepare private `marking.env` beside the ZIP, never inside it or Git, holding only checker-required variables and a dedicated working key/endpoint credential. Arrange revocation after PA4 marks are released.
- [ ] Tag the submitted commit `pa4`, push the tag and verify its actual CI run. Keep repository public, deployment reachable, and release evidence tied to the submitted version.
- [ ] Create the ZIP including `.git`, named `<ascending-student-IDs>_<claimed-total>.zip`; include the self-assessment and red-team evidence, exclude secrets/local configuration. One member submits the ZIP plus separate `marking.env` before the deadline and records the receipt privately.

### Red-team coverage plan

These are proposed attacks to adapt to the assigned target; none is claimed as executed.

| Category | Candidate attacks | Required evidence |
| --- | --- | --- |
| Sandbox minimums | 1. Network attempt; 2. memory allocation with read-back | Exact C++ and response/time/resource result; one attempt per resource kind |
| Model key and model input | 3. Inspect public assets/configuration and verify input minimization with synthetic identity/comment probes | Safe observations/code references; do not use or publish a live credential |
| Prompt injection | 4. Violating code with comment injection; 5. violating code with string-literal injection | Exact source and resulting semantic verdict/evidence |
| Guardrails | 6. Evidence/schema attack adapted to the exposed interface; 7. single cost-cap exercise | Actual reachable request/fixture, outcome and TA reason; label any local-code run |
| Human gate and authorization | 8. Publish with pending reviews; 9. student/other-class TA direct access | Request method/path/body, caller role, response and visibility result |

## Delivery order and checkpoints

| Target date (2026, Asia/Saigon) | Owners | Exit condition |
| --- | --- | --- |
| November 25–26 | All | Accept PA3 baseline; agree publication/sampling contract; audit G1–G5; commit new cases before fixes; obtain target assignment when available |
| November 27–28 | Members 1 and 2 | Backend sampling/Publish handoff and live eval integration available; required-check configuration and blocked-merge evidence started |
| November 29–30 | Member 3 with owners | UI integrated; target attacks recorded; report any unavailable target by November 30; complete guardrail/injection checks |
| December 1 | All | Freeze candidate commit; three real runs, variability/cost report, offline/held-out rehearsal and evidence review complete |
| December 2, before 22:00 | Member 3 with team review | Public `pa4` tag CI green, deployed gate verified, self-assessment/ZIP/private credentials ready and submission recorded |

These are planning targets, not completed work. Leave time to fix provider reliability or integration issues without weakening thresholds.

## Final team check

- [ ] Run the documented quality commands from the repository root:

  ```text
  node scripts/check-go-format.mjs
  go -C backend vet ./...
  go -C backend test ./... -count=1
  npm ci --prefix frontend
  npm run check --prefix frontend
  node scripts/verify-eval.mjs
  node scripts/verify-extra.mjs
  ```

- [ ] Run documented Docker/PostgreSQL integration and browser checks. Verify the README checks and guardrail tests finish offline after dependency installation within the marker's 15-minute limit; live eval is separate.
- [ ] Confirm every live eval nonzero exit fails CI and that required-check settings block merge. Official data/evaluator hashes remain unchanged; no labels were adjusted to turn a failing fix green.
- [ ] Confirm the three runs share one commit/model/configuration and all successful real calls have provider request IDs. Explain provider errors with unknown IDs/usage separately.
- [ ] Demonstrate zero wrongful deductions, recall ≥90% and unsure ≤20% on fresh dev results. Recall ≥95% reaches the top held-out band only if reproduced on the lecturer's set; never claim a held-out result before marking.
- [ ] Verify Publish pending-review/role/class denials through direct requests, the fixed-seed sample test and persisted random sample, TA-configured deductions, and student non-disclosure before Publish.
- [ ] Confirm no sandbox receives credentials, host mounts or a Docker socket; the model still has no mark-writing path. Keep deployment/private credential provisioning and release actions under the team's existing permission rules.
- [ ] Check all six self-assessment rows have evidence; recorded shortcomings and the ZIP total agree. Match tag, CI, deployment, reports and submission version.

