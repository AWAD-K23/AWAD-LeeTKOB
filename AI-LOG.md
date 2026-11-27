# AI log

## 2026-11-27 — PA4 G4 persistent cap audit

- **Tool / task:** Codex; audit G4 for Member 2.
- **Changed:** Added an unexported provider factory seam behind the real worker reservation gate. The regression persists a charge, constructs a fresh budgeted model with a fake provider and checks that exhausted budget yields zero provider calls, unchanged ledger balance and local budget reasons for all pending constraints.
- **Verified:** PostgreSQL-backed `go -C backend test ./internal/worker -run 'TestDurableWorker/G4' -count=1` passed with no provider/network calls. Production still constructs the same Gemini provider when the test seam is absent. Existing concurrent reservation/unknown-usage tests remain unchanged.
- **Limits:** No budget limit or paid-service setting was changed. The user confirmed the provider cap and authorized the upcoming three dev runs and one combined evaluation; those runs are not yet claimed.

## 2026-11-27 — PA4 G3 timeout regression audit

- **Tool / task:** Codex; audit G3 for Member 2.
- **Changed:** Added a cancelling fake model that returns valid decisive JSON after its deadline. The checker must discard that late success, finish within a bounded time and provide the local timeout reason.
- **Verified:** `go -C backend test ./internal/semantic -run TestG3 -count=1` passed offline. The production timeout logic needed no change.

## 2026-11-27 — PA4 G2 evidence regression audit

- **Tool / task:** Codex; audit G2 for Member 2.
- **Changed:** Added fake-model tests for invented text/lines, correct text at the wrong line, stripped-comment-only evidence and blank evidence. Added a positive control proving cleaned evidence maps back to the exact original line and trailing comment.
- **Verified:** `go -C backend test ./internal/semantic -run TestG2 -count=1` passed offline. No evidence matching rule or existing label was relaxed; production checks already enforce this boundary.

## 2026-11-27 — PA4 G1 schema regression audit

- **Tool / task:** Codex; audit G1 for Member 2.
- **Changed:** Added fake-model tests for fenced/prose/invalid JSON, unknown verdicts, missing constraints and duplicate IDs at otherwise correct cardinality. Decisive response fixtures include valid exact evidence so an unrelated missing-evidence rejection cannot hide a broken schema guard.
- **Verified:** `go -C backend test ./internal/semantic -run TestG1 -count=1` passed without network. Existing production schema checks needed no change. Mutation evidence is recorded in the Member 2 handoff after the audit.
- **Limits:** Existing PA3 guardrail history is unchanged; this commit adds G1-specific regression coverage, not a claim that the original guard was introduced separately.

## 2026-11-27 — PA4 Member 2 own cases before checker changes

- **Tool:** Codex.
- **Asked for:** Move to Member 2's PA4 tasks.
- **Changed:** Added five separate own Accepted C++ solutions and ten new labels covering comment/raw-string injection, a macro-hidden loop, an unused-recursion decoy and valid accumulator recursion. Documented label rationale. The original five cases and labels are unchanged; no official files were edited.
- **Verified:** All ten own solutions compiled with C++17/O2 and passed all problem tests in the restricted Docker checker. Unchanged fixtures came from a temporary Git export with CRLF conversion disabled. The official 128-file manifest passed. Re-scoring the saved numbered-line predictions without model calls reproduced 0 wrongful, 35/37 recall and 2/102 unsure.
- **Limits:** These are new cases committed before subsequent checker changes, not evidence of new real-model performance. Original PA3 mixed/case history is retained, not reconstructed. Live runs and the target repository choice are awaiting the user's response.

## 2026-11-27 — PA4 Member 1 handoff and check evidence

- **Tool:** Codex.
- **Asked for:** Finish Member 1's handoff and break down the AI log and commits into focused groups.
- **Changed:** Added `docs/PA4/MEMBER1-HANDOFF.md` with migration steps, API examples, state transitions, integration instructions and verification evidence. Updated README, the Member 1 checklist, project progress and the review/Publish architecture guide. Split the original combined AI-log entry into contract, backend/tests and handoff entries; retained older history. Contract: `6b40ca2`; implementation and tests: `f287fec`.
- **Verified:** Frontend dependencies installed; lint passed with two existing Fast Refresh warnings, and separate type-checking and production build passed. All 128 official file hashes verified. Extra-case verification passed all five unchanged cases in a temporary Git export with CRLF conversion disabled. Git whitespace checks passed before committing.
- **Limits:** The combined frontend check stops at pre-existing formatting in untouched `frontend/src/style.css`. Direct extra-case verification fails on Windows CRLF working-tree outputs; neither fixtures nor the verifier were changed. Member 3 acceptance, UI integration, deployment, live evaluation and release evidence remain pending. These are local commits; no push was performed.

## 2026-11-27 — PA4 Member 1 backend and regression tests

- **Tool:** Codex.
- **Asked for:** Implement Member 1's server sampling, required reviews and Publish gate; group implementation and tests in one commit.
- **Changed:** Added migration 00008, persisted random pass samples, review progress and revision invalidation. Added atomic TA/class-authorized Publish and owning-student published-result access. Reused separate TA decisions, configured deductions, reviewer attribution and append-only API history. Added fixed-seed sampling and direct PostgreSQL API tests alongside the implementation; updated the existing API fixture to load the new migration without changing its assertions.
- **Design decision:** Replaced submission-only locking with a shared lab-row lock and database triggers so worker writes, rejudging, constraint edits and new submissions participate in the same publication boundary. A decision edit hides the publication until Publish is repeated. No user-rejected suggestion is claimed.
- **Verified:** Go format, vet and the full backend suite passed against an isolated PostgreSQL 17 Docker database. Publication tests passed three additional consecutive runs. Coverage includes edge sample counts, restart persistence, pending/sample/role/class denials, strict mark-field rejection, configured deductions, student visibility, stale reviews, incomplete findings and a concurrent mutation that blocks Publish before producing a stale-revision denial. The real Goose migration command applied through version 8.
- **Limits:** No real model call, public deployment or UI integration was performed. The temporary database was isolated from project data and its container was stopped after verification. Tests and implementation are committed together; no earlier test-first history is claimed.

## 2026-11-27 — PA4 Member 1 publication contract

- **Tool:** Codex.
- **Asked for:** Read PROJECT-PROGRESS.md, implement Member 1 in PA4-PLAN.md, then break down the AI log and commits.
- **Changed:** Wrote `docs/contract/publication.md` before implementation. Defined counted submissions, persisted random pass sampling, mandatory reviews, revision invalidation, authorized Publish and student result visibility. Specified request/response shapes and errors for Member 3's UI integration.
- **Design decision:** Population changes invalidate the whole lab's reviews and publication conservatively. Student status/history responses keep their existing shape; a separate protected route exposes published deductions.
- **Verified:** Compared the contract with the existing PA3 review routes and the completed backend implementation.
- **Limits:** Member 3 acceptance and UI integration are pending. These commits group work already completed in the workspace; no test-first commit history is claimed or reconstructed. Backend implementation/tests and handoff evidence follow in separate commits.

> This is a dated history of assistant-supported work, not a live completion tracker. Later entries may supersede earlier status statements. For current PA status, use [the project tracker](docs/PROJECT-PROGRESS.md) and [the PA3 plan](docs/PA3/PA3-PLAN.md); verify each completion claim against its cited commit, test run, or saved evidence.


## 2026-11-24 — TA review queue implementation

- **Tool:** Codex.
- **Asked for:** Finish the PA3 TA review queue, show constraint verdicts and evidence with the deduction chosen through TA review, enforce class ownership, and verify model findings cannot become student deductions or published marks without TA review.
- **Kept:** Model findings separate from TA decisions and student-facing submission results. A TA decision is either pass (zero deduction) or violation (the configured constraint deduction); violation evidence must exactly match submitted source lines. Class access is checked from server-side TA membership.
- **Changed:** Added a latest-counted-submission review queue, TA review records and decision history, a review UI with checker evidence and the TA-confirmed evidence text, and project-plan status plus verification notes. The queue applies a deduction only after a TA decision. Cited FAQ-02, FAQ-08, FAQ-09, FAQ-10, and FAQ-12 in the verification note.
- **Verified:** The user ran the PostgreSQL-backed `TestSubmissionContract/TA_review_is_class_scoped_and_cannot_leak_into_student_marks` test in Ubuntu WSL using Docker; it passed with exit status 0. This verifies the focused API contract, not the deployed end-to-end flow or frontend build. The current backend has no Publish or marks endpoint; no future publish gate is claimed.
- **By hand:** The user supplied the Docker run output. The assistant edited the queue implementation and UI based on the repository rubric and FAQ requirements.


## 2026-11-25 — Centralized architecture and mechanism documentation

- **Tool:** Codex.
- **Asked for:** Consolidate architecture/mechanism guides into one folder and update them for the completed PA3 baseline and upcoming PA4.
- **Changed:** Moved seven guides into `docs/architecture/`, added an index and review/Publish mechanism guide, rebased links, documented review schema through migration 00007, current TA routes/history/reasons, worker environment loading, numbered model input and the newer passing evaluation. Linked the central folder from README.
- **Rejected:** Treating the planned PA4 sampling or Publish design as implemented, duplicating current guides per milestone, or rewriting historical handoff evidence as fresh verification.
- **Verified:** Compared documentation with current Go routes, worker/checker source, migrations and saved evaluation artifacts. Architecture/README/progress/PA4 local links and code fences passed, and Git whitespace checks passed. A wider scan found 16 existing broken links in historical PA2 documents, outside this consolidation. No application code, schema, provider calls or deployment was changed.
- **Evidence:** [Architecture index](docs/architecture/README.md), [review/Publish flow](docs/architecture/REVIEW-PUBLISH-FLOW.md). Older entries below may use the guides' former paths; the index records the consolidation.

## 2026-11-22 — Project overview and reference documentation refresh

- **Tool:** Codex.
- **Asked for:** Check BACKEND.md and DATABASE.md against the current implementation, update outdated details, and rewrite README.md to describe the project as a whole.
- **Kept:** Simple tables and short explanations; separate judge results, constraint findings, and future TA review; general control flow with detailed results retained in evaluation and handoff documents.
- **Changed:** Reorganized README around roles, architecture, setup, judging, evaluation, checks, and documentation. Updated backend environment/CLI/model controls and moved mechanism links. Expanded database budget/audit states. Following the user's correction, removed model scores, measured usage/costs, and quota figures from these overview documents; kept the logic and links to detailed evidence.
- **Rejected:** Keeping SQL-only lab creation instructions, an incomplete backend tree, outdated quota claims, or presenting retained cost reservations as provider bills.
- **By hand:** The user requested the documentation refresh, then clarified that the guides should describe general logic rather than evaluation figures and requested a commit amend. Their separately staged `.env.example` change is outside this documentation task.
- **Verified:** Reviewed command entry points, frontend scripts/proxy, deployment Dockerfile, migration 00005, and saved model evidence. All 29 revised local links resolved; overview evaluation figures were removed; Compose option and Git whitespace checks passed. Documentation-only changes; no model calls or database writes were needed.
- **Evidence:** `README.md`, `docs/BACKEND.md`, and `docs/DATABASE.md`.

## 2026-11-22 — Member 2 flow explanation

- **Tool:** Codex.
- **Asked for:** Explain Member 2 as a visual flow matching the backend-mechanism documents, using Mermaid diagrams, tables, simple examples, and exact code excerpts.
- **Kept:** The distinction between normal judging, hard rules, model findings, and TA review; actual saved smoke/evaluation data rather than invented successful results.
- **Changed:** Added `MODEL-CHECKER-FLOW.md` with the runtime and evaluation paths, source links, budget/rate controls, failure behavior, and step-by-step snippets. Linked it from the existing backend connections guide and clarified that guide's Member 1 scope.
- **Rejected:** Showing TA review as already implemented, treating model findings as deductions, or suggesting that the checker reads labels.
- **By hand:** The user requested a visual explanation similar to the existing mechanism documents.
- **Verified:** All local links resolved; all 16 Go excerpts matched current source after whitespace normalization. The saved worker result confirmed the example's model evidence at line 8. Git whitespace checks passed. Documentation-only change; no model calls or code changes were required.
- **Evidence:** `docs/PA3/backend-mechanism/MODEL-CHECKER-FLOW.md` and `BACKEND-CONNECTIONS.md`.

## 2026-11-22 — Real checker smoke test

- **Tool:** Codex.
- **Asked for:** Verify the configured real model with one smoke request before running the dev set, while respecting the reported 15 RPM limit.
- **Kept:** One sequential request, strict response/evidence validation, and no automatic retries.
- **Changed:** Saved the prediction, provider audit, and usage summary for an Accepted extra solution.
- **Rejected:** Starting the full evaluation before verifying that the configured provider returned a usable result.
- **By hand:** The user entered the API key locally, enabled `MODEL_OFF=0`, and reported the account quota. The assistant did not change provider settings.
- **Verified:** The semantic constraint returned pass with exact evidence and a provider request ID. The estimated cost was US$0.000220.
- **Evidence:** `eval/runs/real-smoke/`.

## 2026-11-22 — Real dev-set evaluation

- **Tool:** Codex.
- **Asked for:** Run the real model on the official dev set within the 15 RPM limit and preserve predictions, call logs, and the unchanged evaluator's report.
- **Kept:** Sequential calls, strict evidence checks, no automatic retries, and retained reservations for unknown usage.
- **Changed:** Saved all 102 predictions and the original evaluation output in `eval/runs/real-dev/`.
- **Rejected:** Retrying 503s to hide failures, replacing predictions with hand-written answers, or modifying labels, evaluator rules, or guardrail tests to improve scores.
- **By hand:** The user authorized continued real calls and supplied the account limit; no provider quota or spending setting was changed.
- **Verified:** The evaluator returned 1: 0 wrongful deductions, 56.8% recall (21/37), and 30.4% unsure (31/102). Of 34 requests, 29 succeeded with provider IDs and known usage; five returned HTTP 503 without IDs or usage. Later budget blocks remained unsure. Peak request rate was 12 per rolling minute.
- **Evidence:** `eval/runs/real-dev/predictions.jsonl`, `calls.jsonl`, `report.txt`, and `summary.json`. The original FAIL report is preserved.

## 2026-11-22 — Model usage and cost accounting

- **Tool:** Codex.
- **Asked for:** Record measured model usage and calculate the estimated cost of 80 comparable submissions in README and AI-LOG.
- **Kept:** Known token usage and unknown error usage are reported separately; retained reservations are not presented as provider bills.
- **Changed:** Added `scripts/summarize-model-run.mjs` and saved reproducible usage summaries.
- **Rejected:** Treating failed calls with missing usage as free or claiming the estimate is a measured bill for all requests.
- **By hand:** No provider billing setting was changed; the user supplied the model configuration.
- **Verified:** The 29 known-usage dev calls used 22,841 input and 2,536 output tokens, with zero reported thinking tokens. At US$0.30/US$2.50 per million respectively, rounded per-call estimates total US$0.013204. For 80 comparable calls, assuming one per submission: `0.013204 / 29 × 80 = US$0.036425`. Unknown usage on five 503 calls is excluded; its US$2.392065 retained reservation is a conservative budget charge. Across all 36 real requests, including both smoke tests, peak usage remained 12 per rolling minute.
- **Evidence:** `scripts/summarize-model-run.mjs`, `eval/runs/real-dev/usage-summary.json`, and README's model/cost section.

## 2026-11-22 — Real worker integration and identity removal

- **Tool:** Codex.
- **Asked for:** Demonstrate real semantic checking on an Accepted submission through the worker and verify that known student identity is removed before provider access.
- **Kept:** Real sandbox execution, existing lease-checked result storage, and separate judge and constraint verdicts.
- **Changed:** Added bounded `worker --once` operation and queued one labelled local smoke fixture for the enrolled manual-test student. Saved a read-only query and its result.
- **Rejected:** Leaving a worker polling after the smoke test or treating model findings as published deductions.
- **By hand:** No separate user edits were recorded for this step; the assistant prepared and processed the local fixture within the authorized verification work.
- **Verified:** The worker stored Accepted, a hard-rule pass, and a model pass with exact comparison evidence. The provider prompt removed known identity text from both a comment and a literal. Estimated call cost was US$0.000444. No marks were published.
- **Evidence:** `backend/cmd/worker/main.go` and `eval/runs/real-worker/`.

## 2026-11-22 — Member 2 documentation and local handoff

- **Tool:** Codex.
- **Asked for:** Document the completed Member 2 work, maintain honest limitations, and check finished tasks in PA3-PLAN.
- **Kept:** Actual metrics and failure evidence, with PA4 accuracy targets distinct from PA3 delivery requirements.
- **Changed:** Added `MEMBER2-HANDOFF.md`; updated README, backend/progress documentation, and the remaining Member 2 checklist items.
- **Rejected:** Claiming PA4-level accuracy, Member 3 acceptance, final public cap verification, or remote CI results from local checks.
- **By hand:** The user requested regular commits and completed-task updates. No teammate acceptance was supplied.
- **Verified:** Final Go format/vet, PostgreSQL-backed tests with fake models, official file hashes, and credential checks passed. No configured API key was found in saved evidence. Sandbox attacks passed in the preceding delivery and were not repeated for the pacing/CLI change.
- **Evidence:** `docs/PA3/MEMBER2-HANDOFF.md`, `docs/PA3/PA3-PLAN.md`, `docs/PA3/MEMBER2-PROGRESS.md`, and README.

## 2026-11-22 — Model request pacing

- **Tool:** Codex.
- **Asked for:** Continue real-model verification after the user enabled `MODEL_OFF=0`, while respecting the account's 15 requests-per-minute limit.
- **Kept:** Sequential calls and no automatic retries. Provider requests in each process are spaced at least five seconds apart, below 15 RPM.
- **Changed:** Added a shared cancellable limiter across Gemini instances and tests for spacing and cancellation. Audit start times now refer to the provider request after any rate wait.
- **Rejected:** Launching parallel real-model tests or assuming the local limiter controls other applications using the same API key.
- **By hand:** The user enabled model calls and reported the provider quota. No quota or provider setting was changed by the assistant.
- **Verified:** Semantic package tests passed, including fake-transport provider tests, rate spacing, cancellation, schema validation, and timeout behavior. No real call had been made at this commit.
- **Evidence:** `backend/internal/semantic/rate.go`, `rate_test.go`, and real-run evidence to follow.

## 2026-11-22 — Generic checker command and model-off image

- **Tool:** Codex.
- **Asked for:** Implement the next Member 2 checker and evaluation tasks sequentially, keep regular commits, and check finished tasks in `docs/PA3/PA3-PLAN.md`.
- **Kept:** Generic discovery from statement, constraint, and source files; labels are never checker input. The offline image has no network, database, Docker socket, host mounts, or credentials. Student source is analyzed, not executed.
- **Changed:** Added `cmd/checker`, dataset and scanner packages, an isolated checker Dockerfile, offline scripts, saved predictions/reports, and CI steps in the existing jobs. Added a separate online audit-file path and per-problem evaluation-run cap; the production worker retains its durable database cap. Checked handoff acceptance, guardrails, logging implementation, checker, and dataset tasks; real-call and final handoff tasks remain open.
- **Rejected:** Counting model-off scores as real-model results, reading labels to produce predictions, hiding the original evaluator's FAIL, or passing CI when bans are missed or cause wrongful deductions.
- **By hand:** Provider-cap verification remains pending; no real model call was made. The user confirmed Member 1 checks were already run.
- **Verified:** The offline Docker image generated 102 official predictions with 0 wrongful deductions, 40.5% recall (15/37), and 46.1% unsure (47/102). All 55 ban pairs were decided; every missed violation was semantic and unsure. The unchanged scorer returned 1. The five extra cases produced 10 predictions, 0 wrongful, 40.0% recall, and 50.0% unsure, also exit 1. Final Go format/vet, PostgreSQL-backed tests, all nine sandbox attacks, syntax tests, checker tests, and official hashes passed. The final image and stricter offline verifier passed on both datasets; remote CI was not run. An initial harness attempt lost tmpfs output after container exit; stdout capture fixed it. A tar buffer limit was then corrected without changing data or labels.
- **Evidence:** `backend/cmd/checker/`, `backend/internal/checker/`, `deploy/Dockerfile.checker`, `eval/runs/`, and `scripts/run-offline-eval.mjs`.

## 2026-11-22 — Pinned official evaluation data and independent cases

- **Tool:** Codex.
- **Asked for:** Continue Member 2 using the lecturer's published dataset, preserve its files unchanged, add independent labelled cases, and update completed PA3 checklist items.
- **Kept:** All 128 upstream files at revision `552c811db2d35a23d641bc4fdffc1ac120be1a28`, including the original labels and evaluator.
- **Changed:** Imported `eval/official`, recorded SHA-256 hashes and Git byte-preservation attributes, and added five separately labelled/tagged recursive-total solutions with test inputs. Cases include compliant recursion, a loop violation, a formula, a text trap, and prompt injection.
- **Rejected:** Modifying official labels, using labels as checker input, or assuming extra cases are Accepted without executing them.
- **By hand:** The user found the official source and requested completed tasks be checked in PA3-PLAN.
- **Verified:** All 128 imported files matched the temporary upstream checkout byte for byte. All five additional solutions compiled with `g++ -std=c++17 -O2` and passed their three input/output cases in a restricted network-free container with no host mounts or credentials.
- **Evidence:** `eval/official-sha256.json`, `scripts/verify-eval.mjs`, and `eval/extra/`.

## 2026-11-22 — Generic statement and type ban coverage

- **Tool:** Codex.
- **Asked for:** Continue Member 2 against the official dev set, address scanner coverage gaps, and check completed PA3 plan tasks only after verification.
- **Kept:** Syntax-based checks with no model calls for ban lists and exact source evidence.
- **Changed:** Added statement cursors for loops and goto, type/template references, and canonical declaration types for aliases and inferred containers. The scanner accepts its actual source path for standalone use.
- **Rejected:** Searching source text for keywords or treating included headers as student use of a banned facility. No problem IDs or solution text were added to production rules.
- **By hand:** The user supplied the official repository URL and requested continued work and checklist updates.
- **Verified:** Rebuilt the sandbox image; all seven original syntax cases and ten new statement/type cases passed using real libclang in restricted Docker containers. No labels or evaluator files were modified.
- **Evidence:** `sandbox/scan.cpp` and `backend/internal/hardrules/check_test.go`.

## 2026-11-22 — Official evaluation source inspection

- **Tool:** Codex.
- **Asked for:** Verify the user-supplied evaluation repository and compare its dataset and evaluator interface with the Member 2 implementation.
- **Kept:** The lecturer's unchanged dataset and evaluator as the required evaluation source; no substitute labels or scorer were created.
- **Changed:** Documented the upstream URL, inspected revision, file layout, JSONL prediction contract, and evaluator exit codes in the Member 2 progress record.
- **Rejected:** Claiming complete ban-list coverage from call-reference scanning alone: the official data also includes loop-statement and container-type bans.
- **By hand:** The user found and supplied the official GitHub URL.
- **Verified:** Inspected upstream README, evaluator source, all six constraint files, and the local scanner. Downloaded a temporary checkout at `552c811db2d35a23d641bc4fdffc1ac120be1a28`. No dataset files were modified, no evaluation was run, and no real model calls were made.
- **Evidence:** https://github.com/fithcmus/awad-leetkob-eval and `docs/PA3/MEMBER2-PROGRESS.md`.

## 2026-11-22 — PA3 Member 2 worker integration and durable limits

- **Tool:** Codex.
- **Asked for:** Continue Member 2 tasks sequentially after local key configuration, integrate semantic checking with the durable worker, and maintain tested commits and formal implementation records.
- **Kept:** Model calls off by default; a US$1 local cap per lab; conservative reservations before network access; existing job lease fencing and result transactions. Verdicts never select deductions.
- **Changed:** Added migration 00005 and a lab-locked model-call ledger. Known usage replaces reserved cost; unknown usage and crashes retain the reservation. Added worker integration, known-identity redaction, shared automatic environment loading, stricter JSON validation, and concurrency/accounting tests. Removed model credentials from the API container configuration.
- **Rejected:** In-memory-only spending counters, resetting budgets on worker restart, giving the sandbox credentials, and treating a saved key as proof of final provider-cap verification.
- **By hand:** The user supplied local provider settings. No provider console setting was changed, and no real model call was made. The assistant started the existing local PostgreSQL service and applied additive migration 00005; Goose reported schema version 5.
- **Verified:** Final Go format, vet, and PostgreSQL-backed tests passed, including fake-model worker integration, concurrent budget admission, retained unknown usage, audit replay refusal, duplicate JSON-key rejection, and comment-cleaning cases. Database tests used temporary schemas. Docker attack tests were not rerun. No real provider calls were made.
- **Cost estimate:** For `gemini-3.5-flash-lite`, 80 calls with 2,000 input and 500 output tokens including thinking cost `80 × (2000 × 0.30 + 500 × 2.50) / 1000000 = US$0.148`. This is an illustrative standard-price estimate, not measured evaluation usage. The real dev-set run remains pending.
- **Evidence:** `backend/internal/worker/model.go`, worker and semantic tests, migration `00005_model_calls.sql`, and `docs/PA3/MEMBER2-PROGRESS.md`.

## 2026-11-22 — PA3 Member 2 Gemini adapter

- **Tool:** Codex.
- **Asked for:** Continue the Member 2 implementation after local model settings were supplied, maintain formal task records, and commit each tested part separately.
- **Kept:** The configured `gemini-3.5-flash-lite` model and `MODEL_OFF=1`. A fixed Google HTTPS endpoint, header authentication, one bounded request, provider-side JSON schema, and mandatory local validation.
- **Changed:** Added the Gemini REST adapter and fake-transport tests. Audit records include the sanitized prompt, response ID, model version, token counts, duration, and standard-price estimated cost including thinking tokens. Provider error bodies and authentication headers are excluded.
- **Rejected:** Real provider calls in tests, automatic retries without another cost reservation, and unaudited model names with unknown rates.
- **By hand:** The user entered local model settings; the key was not displayed. Final provider spending-cap verification is not inferred from the presence of a key.
- **Verified:** Semantic package tests and vet passed. Provider tests use a fake HTTP transport, including successful usage accounting, rejected output, oversized responses, cancellation, and audit failure. No real model call has been made.
- **Sources:** Google [generateContent reference](https://ai.google.dev/api/generate-content) and [standard pricing](https://ai.google.dev/gemini-api/docs/pricing#gemini-3.5-flash-lite), checked November 22, 2026. Text input is $0.30 and output including thinking is $2.50 per million tokens.
- **Evidence:** `backend/internal/semantic/gemini.go` and `gemini_test.go`.

## 2026-11-22 — PA3 Member 2 semantic-checker foundation

- **Tool:** Codex.
- **Asked for:** Implement the Member 2 tasks sequentially, reuse the handoff checks already performed by the user, and pause with local configuration instructions when an API key or other manual provider setup is required. Maintain this log during implementation and commit tested parts separately.
- **Kept:** A provider-independent checker, disabled by default, that preserves hard-rule findings and requests only pending semantic constraints. Model findings remain separate from TA deductions.
- **Changed:** Added response and evidence validation, a context timeout, comment removal with source-line preservation, and fake-model tests. Added `docs/PA3/MEMBER2-PROGRESS.md` to distinguish the tested foundation from pending provider, worker, dataset, and evaluation work.
- **Rejected:** Real model calls in unit tests and marking the entire Member 2 integration complete before provider and evaluation evidence exist.
- **By hand:** The user confirmed Member 1 checks had already run and requested regular commits. No provider setup or key entry is claimed.
- **Verified:** Go formatting, full Go vet, and `go -C backend test ./... -count=1` passed. The semantic tests use fakes and no network. PostgreSQL/Docker integration cases require their opt-in settings and were not rerun; the user supplied Member 1 acceptance. The first vet attempt hit a local cache permission denial; running the checks with cache access succeeded.
- **Evidence:** `backend/internal/semantic/check.go`, `source.go`, `check_test.go`, and `docs/PA3/MEMBER2-PROGRESS.md`.

Tool used: Codex, for planning and initial project setup.

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

## 2026-11-21 — PA3 lab creation backend

- **Tool:** Codex.
- **Asked for:** Implement the Member 1 lab creation backend so a TA can create a lab in their assigned class, including the statement, test cases, constraints, and TA-selected deductions.
- **Kept:** A versioned migration for classes, TA membership, lab statements, test cases, and constraints; a class-scoped lab creation endpoint; server-side TA membership checks; and one transaction for the lab and its related rows.
- **Changed:** Extended the existing lab structure while preserving local data and the existing submission contract. Added PostgreSQL tests for student refusal, another class's refusal, invalid input, and saved tests and constraints.
- **Rejected:** None recorded for this completed part.
- **Verified:** `go -C backend test ./... -count=1` passed with `TEST_DATABASE_URL` pointing to local PostgreSQL on port 5433. The suite ran in an isolated temporary schema, including the new lab creation test.
- **Evidence:** `backend/migrations/00003_labs_and_classes.sql`, `backend/internal/api/lab.go`, `backend/internal/api/server.go`, and `backend/internal/api/api_test.go`.

## 2026-11-21 — PA3 class setup and submission status

- **Tool:** Codex.
- **Asked for:** Implement the Member 1 class and submission backend, including class assignment, student enrollment, and protected judging-status access.
- **Kept:** Admin-only class creation, server-side class membership, automatic enrollment in newly created labs, class/lab lists, session identity lookup, and student-owned submission status. The status endpoint exposes judging results without constraint deductions.
- **Changed:** Added class-student membership, persistent judging result storage, job lease fields, and a separate constraint-verdict table for the worker handoff.
- **Rejected:** None recorded for this completed part.
- **Verified:** The complete Go test suite passed against local PostgreSQL with an isolated temporary schema, including class setup, automatic enrollment, status ownership, class-list privacy, and the existing submission contract.
- **Evidence:** `backend/internal/api/workspace.go`, `backend/internal/api/api_test.go`, and `backend/migrations/00004_judging.sql`.

## 2026-11-21 — PA3 shadcn/ui and Tailwind screens

- **Tool:** Codex and the official shadcn CLI.
- **Asked for:** Build the Member 1 frontend screens using shadcn/ui and Tailwind CSS while retaining the Go backend integration.
- **Kept:** Generated shadcn components, Tailwind's Vite integration, a dark theme, role-based class and lab forms, editable test cases and constraints, student C++ submission, and protected status polling.
- **Changed:** Replaced the placeholder page and temporary JSON forms with the connected application screens. Added the Vite API proxy paths needed by class lists and submission polling.
- **Rejected:** None recorded for this completed part.
- **Verified:** Frontend lint, formatting, TypeScript checks, production build, and the Playwright login-screen check passed. Lint reported two non-blocking Fast Refresh warnings in generated shadcn component exports.
- **Evidence:** `frontend/src/App.tsx`, `frontend/src/components/ui/`, `frontend/src/style.css`, `frontend/components.json`, and `frontend/e2e/app-shell.spec.ts`.

## 2026-11-21 — PA3 isolated C++ sandbox

- **Tool:** Codex.
- **Asked for:** Implement and verify a sandbox that compiles and executes C++ submissions while enforcing all nine PA3 isolation and resource requirements.
- **Kept:** A non-root, read-only Docker container with no network, host mounts, Docker socket, or application secrets; explicit CPU, memory/swap, process, wall-time, and output limits; streamed source input; and forced whole-container cleanup.
- **Changed:** Added a dedicated runner image and nine real C++ attack submissions. Compilation and syntax analysis also run inside the restricted container.
- **Rejected:** None recorded for this completed part.
- **Verified:** All nine Docker attack tests passed. A further wall-time check confirmed the timed-out container is removed; the memory test writes and reads its allocation, and output capture is bounded to 64 KiB.
- **Evidence:** `sandbox/Dockerfile`, `backend/internal/sandbox/runner.go`, and `backend/internal/sandbox/testdata/`.

## 2026-11-21 — PA3 syntax-aware hard rules

- **Tool:** Codex.
- **Asked for:** Implement syntax-aware ban-list decisions before semantic model checking, with tests for compliant traps and evidence lines.
- **Kept:** A libclang-based scanner running in the sandbox, generic backtick ban-list parsing, resolved call targets, and exact source-line evidence. Semantic constraints return `unsure` with `semantic_pending`; this package has no model dependency.
- **Changed:** Added a reusable Go checker interface for the worker and Member 2. The checker analyzes all ban-list constraints together and separates their decisions from semantic work.
- **Rejected:** None recorded for this completed part.
- **Verified:** Real Docker syntax tests passed for actual calls, comments, strings, a student's namespace, an include used only for `std::swap`, using declarations, and macros. Ban-list format tests also passed.
- **Evidence:** `backend/internal/hardrules/check.go`, `backend/internal/hardrules/check_test.go`, and `sandbox/scan.cpp`.

## 2026-11-21 — PA3 persistent queue worker

- **Tool:** Codex.
- **Asked for:** Add a separate judging worker that recovers queued submissions after restart and stores exactly one result set.
- **Kept:** PostgreSQL row-lock claims, renewable leases, fencing tokens, bounded infrastructure retries, sandbox judging, and hard-rule checks only after Accepted. Semantic work remains explicitly pending for Member 2.
- **Changed:** Added `cmd/worker` and a worker package. Results and constraint verdicts are stored in a transaction guarded by the current lease token and database unique keys.
- **Rejected:** None recorded for this completed part.
- **Verified:** PostgreSQL/Docker integration tests passed for queued work with no worker, a new worker processing real C++ to Accepted, one stored result set, no second claim of a live lease, recovery of an expired lease, and refusal of stale or duplicate completion.
- **Evidence:** `backend/cmd/worker/main.go`, `backend/internal/worker/worker.go`, and `backend/internal/worker/worker_test.go`.

## 2026-11-21 — PA3 form accessibility and browser flows

- **Tool:** Codex.
- **Asked for:** Verify the Member 1 shadcn forms and student result flow in the browser.
- **Kept:** Explicit, unique label/control IDs, including repeated test-case and constraint fields, and separate labels so select options do not become label text.
- **Changed:** Added browser tests that verify the TA's submitted tests and deductions, the student's retry key, and the polled Accepted result.
- **Rejected:** None recorded for this completed part.
- **Verified:** Frontend lint, formatting, types, production build, and all three Playwright tests passed. Generated-component Fast Refresh warnings remain non-blocking.
- **Evidence:** `frontend/src/App.tsx` and `frontend/e2e/member-one.spec.ts`.

## 2026-11-21 — PA3 CI and full pipeline checks

- **Tool:** Codex.
- **Asked for:** Run the Member 1 checks together and configure sandbox and worker checks for the `pa3` tag.
- **Kept:** The existing required jobs, with a separate Docker/PostgreSQL sandbox job and no ignored failures. PRs, main, and `pa3` run the checks.
- **Changed:** Extended the API test to submit real code through HTTP and process it through the worker when Docker tests are enabled. Kept failed hard-rule analysis as `unsure` without routing bans to a model.
- **Rejected:** None recorded for this completed part.
- **Verified:** Local Go format/vet and the full PostgreSQL/Docker suite passed, including all nine attacks, syntax traps, worker recovery, and API-to-worker judging. The follow-up failure-handling test also passed. Remote tag CI is configured; no remote run is claimed here.
- **Evidence:** `.github/workflows/gates.yml`, `backend/internal/api/api_test.go`, and `backend/internal/worker/worker_test.go`.

## 2026-11-21 — PA3 deployment environment and handoff

- **Tool:** Codex.
- **Asked for:** Prepare and verify the Member 1 deployment environment and document the worker, sandbox, and next-member interfaces.
- **Kept:** Separate web/API/migration images, PostgreSQL on host port 5433, Caddy routing, a host-worker service template, ignored secret configuration, and model calls off by default. Retained the provider project cap already documented in README.
- **Changed:** Added deployment Compose and Docker files, the full-check PowerShell script, server-only model configuration names, and the Member 1 handoff guide. Documented the custom sandbox choice and how the next members consume verdicts.
- **Rejected:** None recorded for this completed part.
- **Verified:** Deployment images built and the local stack started successfully. The frontend returned HTTP 200 and the protected identity endpoint returned HTTP 401 without a session. The standalone worker built. The documented backend script passed Go format/vet, PostgreSQL checks, all nine Docker attacks, syntax traps, and worker tests.
- **Evidence:** `compose.pa3.yaml`, `deploy/`, `scripts/test-pa3.ps1`, and `docs/PA3/MEMBER1-HANDOFF.md`. This verifies local setup; no public-host provisioning or remote tag CI run is claimed.

## 2026-11-21 — PA3 database reference and completed Member 1 checklist

- **Tool:** Codex.
- **Asked for:** Update the current database reference and check the completed Member 1 plan items against the implemented work and passing checks.
- **Kept:** A schema reference covering all four migrations, separate judge results and constraint verdicts, and a step-by-step submission/lease flow.
- **Changed:** Updated `DATABASE.md` and checked the eleven Member 1 implementation/handoff items with linked evidence. Clarified that publishing the fixed interfaces belongs to Member 1, while receiving-member acceptance remains in Members 2 and 3's own tasks.
- **Rejected:** None recorded for this completed part.
- **Verified:** Reviewed the documented tables and fields against migrations `00001`–`00004`, and linked the completed checklist to passing backend, Docker, browser, and local deployment checks. Public hosting and the remote `pa3` CI run are explicitly not claimed.
- **Evidence:** `docs/DATABASE.md`, `docs/PA3/PA3-PLAN.md`, and `docs/PA3/MEMBER1-HANDOFF.md`.

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
