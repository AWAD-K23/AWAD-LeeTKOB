# AI log

## 2026-11-27 - Follow-up CI failure and independent confirmation

- **Tool / task:** Codex; diagnose reported job 140631667610 without new paid calls. Existing run 46036066453 failed: official 1 wrongful, 36/37 recall, 2/102 unsure; own cases passed. All 77 calls have IDs/known usage, US$0.043797 estimated. Downloaded raw artifacts locally and retained a [compact audit](docs/PA4/evidence/confirmation-repeat-failure.json).
- **Changed / verified:** Both responses repeated input-storage evidence rather than a prohibited operation. Removed the earlier finding from provider confirmation input/audit to reduce anchoring while keeping full sanitized source, constraints and strict validation. New fake-transport wire/audit regression failed before/pass after. No official data, cases, labels, thresholds or checks changed.
- **Limits:** The new candidate is not live-validated and cannot guarantee independent model errors. The prior green run is historical, not current CI status. No rerun, paid call or push was made for this diagnosis.

## 2026-11-27 - PA4 Member 1: review sampling and Publish

- **Tool / task:** Codex; implement Member 1 with a publication contract and tested handoff.
- **Changed:** Migration 00008, persisted random pass samples, mandatory reviews, revision invalidation, atomic authorized Publish and an owning-student published-result endpoint. Only TA decisions and configured deductions determine marks. Shared lab-row locks and database triggers serialize Publish with workers, new submissions, rejudging and constraint edits. Population changes invalidate reviews/publication; decision edits hide publication until Publish is repeated.
- **Verified:** Go format/vet and full PostgreSQL suite passed; publication tests passed three additional times. Coverage includes sample persistence/counts, ownership/roles, pending reviews, forged marks, student visibility, stale revisions and concurrent mutation. Goose applied through version 8. Frontend lint/types/build, all 128 official hashes and original own cases passed. Initial CSS/fixture CRLF blockers were later resolved by LF normalization/Git export without changing cases or assertions.
- **Commits / limits:** [Contract](docs/contract/publication.md) `6b40ca2`; implementation/tests `f287fec`; [handoff](docs/PA4/MEMBER1-HANDOFF.md) `0f7270e`. Work was grouped after implementation; no reconstructed test-first history is claimed. Member 3 UI/deployment acceptance remains separate.

## 2026-11-27 - PA4 Member 2: cases, guardrails and CI

- **Tool / task:** Codex; extend own cases, audit G1-G5 and automate live evaluation.
- **Changed:** Added five Accepted cases before checker fixes (`f1510fa`), bringing `eval/extra` to ten solutions / twenty labelled pairs. New cases cover comment/raw-string injection, macro loops, unused recursion and accumulator recursion. All ten compiled and passed their problem tests in restricted Docker; original cases/labels and 128 official files remain unchanged.
- **Guardrail commits:** G1 strict schema `ed609fe`; G2 exact evidence `4f85bdf`; G3 late-success timeout `3305f2a`; G4 persisted budget refusal with zero provider calls `a4f7d10`; G5 forged deductions/marks and publication API `55b7c85`. G4 adds an unexported fake-provider seam; production admission remains unchanged. These are focused fake-model regressions, not reconstructed origins of the existing PA3 guards.
- **CI / tests:** [Live job](.github/workflows/gates.yml) generates fresh official/own predictions, runs the unchanged evaluator, preserves artifacts and propagates failures; secretless forks fail explicitly. Trusted checker holds the key; isolated Docker scanners have no credentials/network/mounts. All seven final mutation baselines passed and every bypass failed; see [mutation evidence](docs/PA4/evidence/provider-schema-guard-mutations.json). Go/PostgreSQL/Docker, frontend and five Node tooling checks passed. A runner metadata-mismatch bug was fixed instead of weakening its test. Model-off image rehearsals produced valid pairs but correctly failed semantic thresholds; they are not live/proxy proof.
- **Human authorization:** User selected the team repository, verified the provider cap and explicitly approved key transfer, branch push/draft PR and exact ruleset changes. Existing checks/restrictions were preserved; live-eval, one approval and CODEOWNERS review were added. Initial automatic review rejected ruleset/secret operations for insufficient explicit authorization; no workaround was used, and retry followed exact approval. Key was neither printed nor committed; spending limits were unchanged.
- **Limits / evidence:** Dedicated private marking credentials, actual marking-proxy verification and Member 3 deployed UI/injection acceptance remain pending. See [Member 2 handoff](docs/PA4/MEMBER2-HANDOFF.md).

## 2026-11-27 - PA4 real failures, checker corrections and passing CI

- **Tool / task:** Codex; retain failures and validate generic fixes within authorized runs. All successful calls below have provider IDs and known usage, with no automatic retries or unknown-usage reservations. Costs are estimates, not provider bills.

| Frozen code / run | Actual result | Estimated US$ | Evidence |
| --- | --- | ---: | --- |
| `77ef5b0`, three dev runs | Each had 1 wrongful finding; recall 89.2-94.6%, unsure 0-2.9%; 7 changed pairs, 117 calls | 0.062328 | [Comparison](docs/PA4/evidence/baseline-variability.json) |
| Baseline CI 45968971182 | Official passed: 0 wrongful, 35/37 recall, 2/102 unsure. Own failed: 8/9 recall; 49 calls | 0.023860 | [Artifacts](eval/runs/pa4-ci-baseline/) |
| `61d9c89`, three first-fix dev runs | Two passed; third had 1 wrongful finding. All had 36/37 recall; 2 changed pairs, 117 calls | 0.068578 | [Comparison](docs/PA4/evidence/fixed-variability.json) |
| Confirmation CI 45987925317 | Official passed: 0 wrongful, 35/37 recall, 5/102 unsure. Own failed: 8/9 recall, 1/20 unsure; 75 calls | 0.041354 | [Summary](docs/PA4/evidence/confirmation-ci-summary.json) |
| Schema-fix CI 46009732306, head `07d14e0` | All five jobs passed. Official: 0 wrongful, 36/37 recall, 1/102 unsure. Own: 0 wrongful, 9/9 recall, 0/20 unsure; 77 calls | 0.045214 | [Summary](docs/PA4/evidence/schema-ci-summary.json), [CI](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/46009732306) |

- **First correction:** `60bfe73` added generic reachable-computation/macro/alias guidance and isolated invalid evidence per finding; malformed overall schema still fails closed. New independent-finding regression failed before/pass after. The residual wrongful result led to `cd3374d`: a bounded second call tries to disprove semantic violations; disagreement/error/timeout/budget refusal yields unsure. Passes/hard rules need no confirmation. Sanitized evidence reconstruction prevents identity/comment leakage. Fake tests cover agreement/disagreement, malformed output, deadlines, errors and wire framing. Existing evidence/verdict assertions remain; intentional call-count assertions became two.
- **Final correction:** Confirmation CI exposed a violation with empty evidence: local validation correctly rejected it, but the provider schema had allowed it. `e592bf7` requires nonempty evidence for decisive verdicts using conditional anyOf/minItems. A pinned test-only JSON Schema validator checks the actual wire schema with counterexamples; regression failed before/pass after. G1/G2 were not relaxed. Source: [Google structured-output documentation](https://ai.google.dev/gemini-api/docs/structured-output). Cases, labels, official files, thresholds and required checks were never loosened.
- **Authorization:** Initial approval covered three dev runs plus one combined evaluation. User then approved up to four further evaluations: three first-fix runs plus confirmation CI, US$0.109932. After that batch was exhausted and the local schema fix explained, "continue" was explicitly scoped to one further combined CI run. No further paid run followed.
- **Final limits:** [PR #27](https://github.com/AWAD-K23/AWAD-LeeTKOB/pull/27) remains draft and requires human review; no merge/deployment. Evidence commit `1b8fdb3` remains local to avoid another paid CI run. One green combined run does not prove latest-method three-run stability or held-out performance; older failures remain disclosed in [PA4 report](PA4.md).

## 2026-11-22 to 25 - PA3 review queue and documentation

- **Tool / task:** Codex; implement TA review, explain system flows and reconcile handoffs.
- **Changed / verified:** Added latest-counted review queue/history/UI with server-side class ownership, exact evidence and TA-selected deductions. User ran the focused PostgreSQL API contract successfully in Ubuntu WSL/Docker; this was not deployed UI proof. Publish was added later in PA4.
- **Human direction / documentation:** User requested general architecture explanations and removal of evaluation scores/costs from overview guides. Consolidated seven guides into [architecture](docs/architecture/README.md), updated README/database references and checked links/source excerpts without presenting planned Publish/sampling as implemented. Sixteen pre-existing historical PA2 broken links were outside that cleanup. No model calls/deployment were needed.

## 2026-11-22 - PA3 Member 2 foundation and evaluation

- **Tool / task:** Codex; checker, Gemini adapter, durable worker accounting, syntax coverage and generic online/offline evaluator.
- **Changed:** Strict schema/exact evidence checks, line-preserving comment removal, identity redaction, bounded requests and sanitized audit records; model off by default. Migration 00005 adds lab-locked reservations before calls, replacing them with known cost or retaining unknown/crashed reservations across restarts. API/sandbox get no key. Five-second cancellable pacing respects the reported 15 RPM quota within this process, not other key users.
- **Rejected / human role:** User entered settings privately and authorized real verification; assistant changed no provider limits. Rejected real unit-test calls, in-memory-only budgets, unreserved retries, treating unknown usage as free, and inferring cap verification from a saved key. The lab ledger tracks this application, not provider-console spending or other applications.
- **Datasets / offline checks:** Pinned 128 lecturer files to `552c811db2d35a23d641bc4fdffc1ac120be1a28`, recorded hashes, and added five independently labelled Accepted cases. Expanded libclang statement/type/macro/alias bans without dev-specific rules. Labels never enter the checker. Model-off official result: 0 wrongful, 15/37 recall, 47/102 unsure; own: 0 wrongful, 40% recall, 50% unsure. Both scorer exits were 1. Fixed tmpfs-output capture/archive buffering without changing data. Go format/vet, fake-provider, PostgreSQL accounting/lease, Docker attacks/syntax and hash checks passed.
- **Real smoke:** Checker smoke (original files retained in Git history through `1b8fdb3`) returned exact evidence/provider ID (US$0.000220); [worker smoke](eval/runs/real-worker/) stored Accepted and separate hard-rule/model passes, removed identity from comment/literal (US$0.000444). Added worker --once; no marks published or worker left polling.
- **Original real failure:** [Dev run](eval/runs/real-dev/) produced 102 predictions: 0 wrongful, 21/37 recall (56.8%), 31/102 unsure (30.4%), scorer exit 1. Of 34 calls, 29 had IDs/known usage and five HTTP 503s had neither. Later budget blocks stayed unsure; peak rate 12/minute. Failures were preserved without retries or label/scorer changes.
- **Costs:** [Summarizer](scripts/summarize-model-run.mjs) records 22,841 input / 2,536 output tokens for 29 known calls: US$0.013204; extrapolating one call per submission to 80 was US$0.036425. Unknown 503 usage is excluded; US$2.392065 retained reservation is a conservative local charge, not a bill. Later numbered-line run: 0 wrongful, 35/37 recall, 2/102 unsure, 39 known-ID calls, US$0.020851 (deployment fix attributed below). PA4 confirmation changes the old one-call estimate. Historical pricing was checked against [Google pricing](https://ai.google.dev/gemini-api/docs/pricing#gemini-3.5-flash-lite) on 22 November.
- **Evidence / limits:** [Progress](docs/PA3/MEMBER2-PROGRESS.md), [handoff](docs/PA3/MEMBER2-HANDOFF.md), [checker](backend/cmd/checker/), [worker accounting](backend/internal/worker/model.go). Local completion did not establish public deployment, Member 3 acceptance or held-out performance.

## 1. Local PostgreSQL port

- **Task:** Start the local PostgreSQL service.
- **Problem:** Docker could not bind host port 5432 during setup.
- **Changed:** The assistant changed the host port to 5433, updated the connection instructions, and confirmed PostgreSQL accepted connections. `AGENTS.md` records the port so the team uses the working setup.
- **Evidence:** `compose.yaml` and `README.md`; the first repository commit preceded this change.

## 2. Example environment file

- **Task:** Prepare local configuration while following the PA2 rubric.
- **Problem:** `.env.example` contained values, while the rubric requires variable names without values.
- **Changed:** The assistant copied the local settings to ignored `.env` and cleared the values from `.env.example`. `AGENTS.md` now states this convention.
- **Evidence:** `.env.example` and `.gitignore`; see the current change history for the setup update.

## 3. Framework choice

- **Task:** Discuss a frontend and backend stack.
- **Rejected:** The user rejected the assistant's recommendation to use Next.js for the backend, preferring a backend-focused design.
- **Changed:** The plan uses React with TypeScript and a separate Go backend.
- **Evidence:** `docs/PA2-PLAN.md`, step 1; repository setup commit `aa16e18`.

## 2026-11-08 — 4. Submission setup: login, enrollment and `POST /labs/{labId}/submissions` (H2, commit `2c7e6ef`, PR #3)
Tool: Codex
Asked for: implement login, enrollment and the submission endpoint with a simple PostgreSQL queue, following `docs/contract/submissions.md`.
Kept: SCS sessions stored in PostgreSQL, versioned migrations, and the submission and its pending job inserted in one transaction; PostgreSQL contract checks run in an isolated temporary schema. The second migration (`00002_source_bytes.sql`) stores source as bytes, because PostgreSQL `text` cannot hold a NUL byte while the contract allows one in `source`; a contract check verifies the bytes are preserved.
Changed: No separate human rewrite is recorded for this task. The assistant's NUL-byte storage correction is included under Kept.
Rejected: None recorded.
By hand: The team selected the endpoint and contract; the available record does not attribute further hand edits to a member.

## 2026-11-09 — 5. H1/H2 contract follow-up (commit `2b18a48`, PR #5)
Tool: Codex
Asked for: recheck H1/H2 against the frozen submission contract at `9f25fab`, and add the missing contract tests.
Kept: the duplicate-Origin finding — the handler accepted two `Origin` values when the first matched, although contract §4 says "exactly the configured application origin". The new `TestOriginMustMatchConfiguredValue` failed first (`api_test.go:67: multiple Origin values accepted`), then `validOrigin` was changed to require exactly one matching value. Also kept the expanded PostgreSQL HTTP tests: every error type's status, type and title, check order, ownership/enrollment, idempotency, rollback, source boundaries and expired sessions. Red and green runs: `docs/PA2-H1-H2-VERIFICATION.md`. The tests and `server.go` changed in the same commit; the tests got stricter (no assertion weakened) and the contract text did not change, because the code was the part that disagreed with it. The original red run from before the endpoint was implemented was never kept; this follow-up does not replace it.
Changed: No separate human rewrite is recorded. The assistant made the Origin fix and added assertions against the existing contract.
Rejected: None recorded.
By hand: The team requested the contract comparison; the available record does not attribute further hand edits to a member.

## 6. PA2 checks and CI

- **Task:** Add the highest-priority reproducible checks and a required-ready quality gate, following the course harness Session 7 guidance.
- **Changed:** Codex added frontend lint/format/type/build commands, a Go format check, CI jobs for Go/PostgreSQL and frontend checks, a full-history secret scan, and a separate Playwright shell check. README and `AGENTS.md` now list the commands CI runs. The lockfile audit exposed a moderate `minimatch` advisory in the lint toolchain; a scoped override removed the reported findings.
- **Changed after CI review:** The first Gitleaks PR run returned HTTP 403 because the security job lacked pull-request read permission. Codex added only `pull-requests: read` to that job; the subsequent CI run passed. This was a workflow-permission failure, not a deliberate contract-break demonstration.
- **Observed:** Running Vite from the requested local `PA#2` path fails because `#` in that filesystem path is parsed as a URL fragment. Fresh-clone checks passed from a path without `#`; the Windows path limitation is recorded in `docs/PA2-QUALITY-GATE-VERIFICATION.md`.
- **Evidence:** Verification commit `90b68ef` passed `gates`, `secret-history`, and `e2e` in [run 42661351986](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661351986); the Go log shows 5 unit tests and 46 contract subtests. Codex then temporarily bypassed the enrollment check in code-only commit `3d030a5`; `gates` failed on the unenrolled/missing-lab cases in [run 42661114861](https://github.com/longtoZ/AWAD-LeeTKOB/actions/runs/42661114861). Restoring the guard in `90b68ef` produced the green run. This demonstrates that the tests catch this defect; it does not prove that GitHub blocked a merge.
- **Status at that time:** Branch protection and public visibility were still unresolved. They were established later on 10 November; see `PA2.md` and PR #10. The proposed US$10 provider budget was not an enforced spending cap.

## 2026-11-10 — H3: protected merge block and deliberate contract break (Pham Bao Huy)
Tool: Claude Code (Claude Opus 5.5).
Asked for: check my PA#2 tasks (H1–H4) against the rubric and finish H3: a red required check that blocks a merge, then the commit that turns it green.
Kept: the finding that every refusal test checked only the problem response, and that the replay test counted jobs for the first submission ID only, so a handler that wrote a second submission and job on a retry and still replayed the stored `201` stayed green. Kept the before/after row-count checks (PR #9, tests only, +30 lines, nothing removed) and the code-only replay break `22030bc` (PR #10). Before opening PR #10 the assistant checked the claim both ways in throwaway postgres:17 / Go 1.26 containers: the old test file passed against the break, the new one failed on `creation_and_replay` and `key_retained_after_23_hours`.
Changed: nothing in the code or tests. I read the diff of PR #9 (30 added lines, no assertion removed) and of `22030bc` before pushing, and checked that PR #10 touched only `submission.go`.
Rejected: repeating the enrollment-guard break from the team handoff (PR #7). That break was already shown red in run 42661114861, so a second one would prove nothing new; the replay break showed a hole the old tests could not catch.
By hand: opened and merged PR #9 and PR #10; checked that the required `gates` check blocked the merge and took the screenshot; chose which break to make. No additional hand work is recorded.

## 2026-11-10 — PA#2 evidence and hand-in preparation
Tool: Codex (this task).
Asked for: resolve the remaining PA#2 gaps against the live rubric and prepare a submission ZIP.
Kept: the existing contract, tests, and red/green PR evidence. A fresh second reader identified the duplicate write at `submission.go:230–235` in the deliberate-break commit against contract §4; it did not run CI. No product implementation was generated in this task.
Changed: corrected stale status in `PA2.md` and the self-assessment, replaced unfilled historical log prompts with explicit "not recorded" statements, and aligned the report with confirmed team IDs and live evidence.
Rejected: treating the draft 80/100, an unverified spending cap, or an uncited product-owner answer as established facts.
By hand: Tran Anh Khoa supplied the team's three student IDs. Further team review of the final score and archive is not recorded here.

## 2026-11-11 — Gemini project spend-cap evidence
Tool: Codex (this task), using screenshots supplied by the user.
Asked for: update the PA#2 permissions-and-secrets evidence based on Google AI Studio screenshots.
Changed: saved screenshots showing the `AWAD-LeeTKOB` project identity, linked billing/paid-tier setup, and a ₫100,000 monthly project spend cap. Updated the README, PA2 evidence page, and self-assessment. The score claims the rubric's top band while disclosing Google's experimental cap and approximate 10-minute billing delay, which can permit overages.
By hand: the screenshots were supplied by the user. No API key value or identifier was added to the repository; model integration is not implemented.

## 2026-11-11 — Remaining PA#2 rubric evidence
Tool: Codex (this task).
Asked for: close the remaining self-assessment gaps for commands, protected CI, and the endpoint contract.
Kept: the README commands and CI workflow, existing red/green PR #10 evidence, active main ruleset, and the frozen submission behavior. No application code, tests, test data, or labels were changed.
Changed: recorded a clean local run of every documented check from a path without `#`; documented the post-protection main history and linked PRs #9–#13; added contract source and cookie-jar notes distinguishing the class brief from team choices and stating that no product-owner FAQ answer is relied on. Updated the draft claims for criteria 3–5 to 10/10, 20/20, and 25/25 respectively. The draft total is 98/100 pending merge and final tag/package.
Verified: frontend install/check, Go format/vet/PostgreSQL tests, dependency audit, and Playwright browser smoke test passed locally. The main branch history and active ruleset were read from GitHub; no external changes were made.

## 2026-11-11 — Self-assessment rubric calibration
Tool: Codex (this task), comparing the report against the supplied PA#2 rubric.
Changed: revised the draft total from 98/100 to 100/100. Criterion 6's full-credit band requires contract-error/ownership/idempotency coverage in required CI and a deliberate red/green break; uncertain-commit fault injection is not listed. Criterion 7's five-point description requires named tools/tasks and specific Changed/Rejected entries, which the log contains. Kept the absent fault-injection test and missing historical H1/H2 evidence disclosed as limits, without using them to deduct points the rubric does not specify.
Rejected: treating the rubric's ±10 self-assessment tolerance as permission to inflate a score; the updated score is based on the written full-credit descriptions. If the instructor awards 98, the 2-point difference is within the rubric's no-adjustment range.

## 2026-11-11 — Clarify browser test setup
Tool: Codex, following read-only review of the PA#2 evidence changes.
Changed: separated the Playwright test command from browser installation in README and clarified that Windows/macOS and Linux CI need different browser setup flags, while both run the same `npm run e2e --prefix frontend` check.
Kept: the existing local test result and the CI workflow command; no application or test files changed.

## 2026-11-21 - PA3 Member 1 backend, UI, sandbox and handoff

- **Tool / task:** Codex and official shadcn CLI; class/lab setup, student submission/status, persistent worker, syntax rules and connected React screens.
- **Changed:** Migrations 00003-00004; admin/TA/class/student ownership; transactional lab/test/constraint creation; lease/fencing worker with restart recovery and one atomic result set. Hard rules run after Accepted; semantic checks were then pending. UI includes role forms, accessible repeated fields, idempotent retry and status polling; student status exposes no deductions.
- **Sandbox:** Non-root/read-only Docker, no network/secrets/socket/mounts; CPU/memory/swap/process/time/output limits and whole-container cleanup. Nine real C++ attacks passed; memory test writes and reads allocation, output bounded to 64 KiB. libclang tests cover calls, comments/strings, namespaces, includes, using declarations and macros. Failed analysis stays unsure without model fallback.
- **Verified:** Local Go format/vet/PostgreSQL/Docker tests passed, including API-to-worker C++, queue restart, stale/duplicate lease refusal and single results. Frontend lint/format/types/build and three Playwright flows passed (two generated-component Fast Refresh warnings). Deployment images/local stack built: frontend HTTP 200, unauthenticated identity HTTP 401. Added CI Docker/PostgreSQL and pa3 triggers with no ignored failures; no remote run/public deployment was claimed then.
- **Evidence / limits:** [Handoff](docs/PA3/MEMBER1-HANDOFF.md), [plan](docs/PA3/PA3-PLAN.md), [check script](scripts/test-pa3.ps1), [deployment](deploy/), [sandbox attacks](backend/internal/sandbox/testdata/). PostgreSQL host port 5433; model off by default. Completed implementation checklist did not imply receiving-member acceptance.

## 2026-11-23 — Repair secret-history license failure

- **Tool:** Codex.
- **Asked for:** Resolve the reported Gitleaks Action failure requiring `GITLEAKS_LICENSE`.
- **Changed:** Replaced the licensed Action wrapper with standalone MIT-licensed Gitleaks CLI v8.30.1. CI verifies the pinned Linux archive checksum, retains full-history checkout and the `secret-history` job name, redacts findings, and fails on leaks or scanner errors. Removed the unused PR-read permission and documented the command in README.
- **Rejected:** Disabling the scanner, suppressing failures, adding a baseline to hide findings, or requiring license provisioning for the standalone scanner. No scan rules or existing tests were weakened.
- **Verified:** Parsed workflow YAML and checked full checkout, scan arguments and failure propagation. Downloaded the official Windows CLI to a temporary folder and verified its release SHA-256. It scanned all 66 local commits with no leaks. A separate temporary Git repository passed when clean and failed with exit 1 after a synthetic token was committed and then removed from HEAD, confirming historical detection. No real credentials were used in the fixture.
- **Limits:** The Linux install step and remote GitHub run have not been executed here; local validation used the same CLI release on Windows. No commit, push or remote CI success is claimed.
- **Sources:** [Action licensing](https://github.com/gitleaks/gitleaks-action#environment-variables), [CLI documentation/license](https://github.com/gitleaks/gitleaks/tree/v8.30.1), [release checksums](https://github.com/gitleaks/gitleaks/releases/tag/v8.30.1).

## 2026-11-24 — PA3 FAQ-09 unsure reasons

- **Tool:** Codex.
- **Asked for:** Add the FAQ-09 cause for each `unsure` finding to the TA review details, add tests, and record the work. Keep FAQ-10 random pass sampling as a separate PA4 follow-up.
- **Changed:** Added locally determined reasons for model uncertainty, invalid response schema, mismatched evidence, provider failure, timeout/cancellation, exhausted lab budget, disabled or missing model configuration, and source preparation failure. Persisted the reason with each constraint finding, returned it only through the TA review queue, and displayed it for unsure findings. Added a schema migration and regression coverage for reason generation, persistence, API visibility, and student non-disclosure. FAQ-10 pass sampling remains a separate PA4 follow-up.
- **Rejected:** Allowing model-generated explanation text; reasons are derived from checker outcomes.
- **Verified:** [GitHub Actions run #43](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45243460493) passed all four jobs (gates, sandbox, end-to-end, and secret-history) for commit `53f66e2`. Gates ran Go format/vet/tests and frontend checks; sandbox ran Docker-backed tests. No separate local Ubuntu run is claimed for this addition.
- **Evidence:** `backend/internal/semantic/check.go`, `backend/internal/semantic/check_test.go`, `backend/internal/worker/model.go`, `backend/internal/worker/worker.go`, `backend/internal/worker/worker_test.go`, `backend/internal/api/review.go`, `backend/internal/api/api_test.go`, `frontend/src/pages/review.tsx`, and `backend/migrations/00007_verdict_reasons.sql`.


## 2026-11-24 — PA3 local smoke documentation and Member 3 handoff

- **Tool:** Codex with GitHub browser; the user ran the Ubuntu WSL/Docker commands and supplied their output.
- **Asked for:** Record the merged-build local smoke, review PR #19, reconcile handoff status and evidence links, and prepare Member 3's remaining work.
- **Changed:** Added the local Compose smoke record, clarified AGENTS/README/PA3 plan documentation, reconciled the dated Member 1 handoff, corrected the durable-worker evidence link, and added the Member 3 handoff.
- **Rejected:** Treating the local `MODEL_OFF=1` run as a public deployment or live-model result, or treating the TA review API test as proof of exactly one persisted worker result set. No provider call, final tag run, or public flow is invented.
- **Verified:** Reviewed the current PA3 Classroom post and original assignment materials, the merged implementation evidence and user-supplied Ubuntu smoke outputs. PR #19 previously passed its four CI jobs on head `3717c62`; the final documentation head requires its own CI result. The local smoke observed Caddy navigation, sign-in, class isolation, TA confirmation, and queued work completing after restart.
- **Limits:** The local two-integer lab is not the rubric's public `sort-merge` lab. Live semantic behavior, final spending-cap check, public walkthrough, official hash comparison, `pa3` tag CI, teammate acceptance, and submission remain open.
- **Evidence:** [PR #19](https://github.com/AWAD-K23/AWAD-LeeTKOB/pull/19), `docs/PA3/LOCAL-COMPOSE-SMOKE-2026-11-24.md`, `docs/PA3/TA-REVIEW-QUEUE.md`, `docs/PA3/MEMBER1-HANDOFF.md`, `docs/PA3/MEMBER3-HANDOFF.md`, and `docs/PA3/PA3-PLAN.md`.

## 2026-11-24 — PA3 tracker reconciliation

- **Tool:** Codex, using the current Classroom assignment and course announcements, repository source, saved evidence, and GitHub Actions.
- **Asked for:** Find stale project trackers and synchronize them in one reviewable change.
- **Changed:** Updated the project overview, PA3 plan, Member 1 handoff follow-up, Member 2 progress snapshot, Member 3 handoff baseline/FAQ note, and TA review CI evidence. Clarified at the top of this log that it records history rather than current status. Kept past entries intact.
- **Verified:** Classroom still lists PA#3 due November 25, 2026 at 22:00 ICT. Main commit `97346d0` passed all four jobs in [Actions run #56](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45326541493). The separate local smoke remains on `1122f7c` with `MODEL_OFF=1`; no public deployment or live-model flow is claimed.
- **Still open:** Member 3 handoff acceptance, public `sort-merge` walkthrough and live-model verification, final provider-cap check, reproduction of the transient “Lab not available” message, official upstream hash comparison, PA3 tag CI, `PA3.md`, five-minute recording, self-assessment, ZIP and submission receipt.

## 2026-11-24 — PA3 Member 3: handoff acceptance, deployment runbook, PA3.md draft (Phạm Bảo Huy)
Tool: Claude Code (Claude Opus 5.5)
Asked for: Do the Member 3 work. Run the Member 1 and 2 handoff checks, rehearse the `sort-merge` flow locally, and prepare the public deployment, `PA3.md` and the tracker updates.
Kept: The acceptance record (`docs/PA3/MEMBER3-ACCEPTANCE.md`): every documented check passed on `df6f621`, run in Docker. The dev set is byte-identical to a fresh clone of the published repository. The marker-style `docker run --network none` checker run exits 1, never 2. `scripts/pa3-walkthrough.mjs` passed every step locally with the worker started after the submissions. `deploy/vm-bootstrap.sh`, `deploy/vm-accounts.sh` and `docs/PA3/PUBLIC-DEPLOYMENT.md`. The `PA3.md` draft.
Changed: None in the files. I directed the work: chose a new cloud VM over the alternatives and the team's AI Studio project key; asked what a stopped VM costs, which led to the deallocate notes in deploy/azure-create-vm.sh.
Rejected: The first VM in Azure eastasia (Hong Kong): Gemini refused every call from it ("User location is not supported"), so I deallocated it and created a new one in Japan East. Found by the assistant's own checks: shellcheck caught an apostrophe in `${DOMAIN:?...}` that broke quoting in the first `vm-bootstrap.sh`. Gitleaks flagged a `password: value` line in a sample PA3.md, so the accounts go in a table instead.
By hand: Installed Azure CLI and signed in; arranged subscription access with the subscription owner; ran azure-create-vm.sh twice (Hong Kong, then Japan East); ran vm-bootstrap.sh and vm-accounts.sh on the VM; put the model key into /etc/leetkob/worker.env with sudoedit and set MODEL_OFF=0; deallocated the Hong Kong VM; opened PR #23.
Limits: the local rehearsal ran with `MODEL_OFF=1`; no public deployment, live-model flow or tag CI is claimed here. Three findings are recorded for their owners: the queue shows only each student's latest submission (FAQ-02); "Lab not available" comes from a stale `["labs"]` cache; hard-rule verdicts keep the reason "Semantic model check has not run."

## 2026-11-24 — Model evidence lines one too low on the deployment (Phạm Bảo Huy)
Tool: Claude Code (Claude Opus 5.5)
Asked for: Find why every model verdict on the public deployment became `unsure`, then fix it with a test first.
Kept: The diagnosis from the call audit rows (right verdicts for `s01`/`s04`/`s09`, every evidence line one too low, the same prompt counted correctly on 22 Nov); the test `TestGeminiSendsNumberedLines` (commit `c05230c`, failing on `main`); the fix in `backend/internal/semantic/gemini.go` that sends the code as `{line, text}` objects while evidence validation stays exact (commit `1cdc486`).
Changed: None in the code. I chose to have the fix made in Member 2's area with my name on the commits instead of handing it back.
Rejected: None (assistant's design choice: alternative considered and not used: accepting evidence whose text matches a different line, because that loosens the check that stops invented evidence)
By hand: Created the git-ignored local .env with the key for the real dev-set run (key never pasted into the chat); approved the push; opened PR #22; redeployed the VM to the fixed commit.

## 2026-11-24 — PA3 review fixes: earlier submissions, lab cache, hard-rule reason (Phạm Bảo Huy)
Tool: Claude Code (Claude Opus 5.5)
Asked for: Fix the three Member 3 acceptance findings, each with a failing test first.
Kept: Read-only `GET /labs/{labId}/earlier-submissions` and the review page section *Earlier submissions (not counted)* (tests `7f9aa78`, code `44523a9`); `refetchType: "all"` before opening a new lab (browser test `bf10683` delays the list request so the stale render is seen every run, fix `2c69b0e`); clearing the pending reason when a ban list decides (test `f9bf22b`, fix `744188a`). Findings 2 and 3 are in files owned by Tô Thanh Long; the commits say so.
Changed: None in the code.
Rejected: Two of the three options for the disappearing submission: only a note in PA3.md (the marker could still read it as missing), and a third student account (the rubric asks for two). Also rejected handing findings 2 and 3 back to Long or leaving them as known issues. (assistant's design choice: not used: putting earlier submissions into the queue itself, which FAQ-02 rules out)
By hand: Chose the read-only "Earlier submissions" design; opened PRs #24 and #25; redeployed the VM twice.

## 2026-11-24 — PA3 self-assessment draft and README model cost (Phạm Bảo Huy)
Tool: Claude Code (Claude Opus 5.5)
Asked for: Draft SELF_ASSESSMENT_REPORT.md for PA#3 against the rubric, checking every claimed piece of evidence.
Kept: The draft claims 92 / 100. Checks behind it: 39/39 logged prompts carry provider ids, none holds a ban-list constraint or comment markers; the sandbox flags; the memory read-back; five labelled own cases. The README gained the measured cost of one 80-submission lab: 0.020851 / 39 × 80 = US$0.043 (`eval/runs/real-dev-numbered/usage-summary.json`).
Changed: none
Rejected: none
By hand: recorded and uploaded the video, took the screenshot
