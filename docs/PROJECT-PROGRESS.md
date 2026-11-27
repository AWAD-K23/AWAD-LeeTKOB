# LeeTKOB progress — PA1 to PA6

**Updated:** November 25, 2026 (Asia/Saigon).

**Implementation reviewed:** [`ce82a9d`](https://github.com/AWAD-K23/AWAD-LeeTKOB/commit/ce82a9d897754547a898bac6e06056f9a015d5dc), with the final PA3 report and saved evaluation evidence.

**Current checkpoint:** PA4, due December 2 at 22:00; PA3 is complete per the team. See the [PA4 plan](PA4/PA4-PLAN.md).

This tracker compares the repository with the [class project brief](<LeeTKOB — class project brief.md>), the full PA2–PA6 rubrics, implementation, and saved evidence. The PA4–PA6 sections were revised after reading the complete [PA4 rubric](<PA4/PA#4 rubric — Eval, guardrails, human gate and red-team.md>), [PA5 rubric](<PA5/PA#5 rubric — Final build and operations.md>), and [PA6 rubric](<PA6/PA#6 rubric — Final demo and individual oral.md>). Checked work means implemented or supported by saved evidence, not an instructor mark, predicted score, or confirmation of submission.

| Mark | Meaning |
| --- | --- |
| ✅ | Finished, with implementation or recorded evidence |
| ⬜ | Still open, partially implemented, or awaiting required verification |
| ❔ | Historical completion cannot be established from this repository |

## Checkpoint overview

All dates below are in 2026. PA2–PA5 deadlines are at 22:00, as stated in the brief.

| Checkpoint | Due / session | Goal | Current position | Main remaining work |
| --- | --- | --- | --- | --- |
| PA1 | Not specified in the supplied brief | Original team proposal | ❔ Historical checkpoint; proposal/submission evidence not located | Confirm the team's original proposal record if this tracker must cover submission history |
| PA2 | November 11 | Harness, quality gate, endpoint contract and tests | ✅ Core engineering delivered; submission evidence recorded separately | Keep the current checkout's checks green; see [PA2 evidence](../PA2.md) |
| PA3 | November 29 (extended from November 25 per PA3 plan) | Working queue/sandbox/model pipeline ending in the TA review queue | ✅ Complete per team; public model-on walkthrough, report and self-assessment recorded | Keep deployment available through PA4 marking; retain the disclosed 8:48 recording limitation |
| PA4 | December 2 | Real-model CI evaluation, held-out run, five guardrails, human gate, red-team report | ⬜ Saved numbered-line dev run passes all three thresholds; PA4 integration and evidence remain | Follow the [PA4 plan](PA4/PA4-PLAN.md): live required eval gate, ten custom cases, repeated runs, Publish/random sampling and red-team report |
| PA5 | December 16 | Final build running and observed | ⬜ Submission updates and deployment scaffolding already exist | Publication/appeals, cross-queue traces, dashboard/SLO/fired alert, reproducible load proof and incident evidence |
| PA6 | December 24 / December 31, 08:20–11:50; assigned slot | Final demo (40 shared points) and individual oral (60 points) | ⬜ Existing PA5 build must be demonstrated; individual readiness unverified | Commit/post OWNERS.md by December 21 at 22:00; rehearse all eleven features, failure cases and explanations |

**Overall assessment:** PA3 is complete per the team, with a recorded public model-on flow and a saved dev run meeting all three thresholds. PA4 adds a required live evaluation gate, Publish, random pass sampling, repeated-run evidence and the red-team report; a single passing dev run does not establish held-out performance. PA6 reuses the PA5 build and also scores a live demonstration, not just the oral; any beyond-brief feature claimed for PA6 must already be in the submitted PA5 build.

| Checkpoint | Rubric criteria and weights | Course weight |
| --- | --- | --- |
| PA4 | CI eval 20; held-out 20; guardrails 20; human gate 15; red-team 20; AI log 5 | 7% |
| PA5 | Deployment 10; publishing/appeals/results 20; observability 20; load 20; post-mortem 20; AI log 10 | 3% |
| PA6 | B1 eleven features 22; B2 submitted build/persistence 4; B3 two failure cases 6; B4 authorization 4; B5 beyond brief 4; individual explanation 60 | 15% |

These weights are not claimed marks. The current [self-assessment](../SELF_ASSESSMENT_REPORT.md) covers PA3; prepare a rubric-specific assessment for PA4 and PA5.

## PA1 — proposal

The brief says each team's original PA1 proposal keeps its mark, while every team builds LeeTKOB from PA2 onward. The named LeeTKOB product-owner authors are not evidence that our team authored that proposal.

| Status | Job | Evidence / next action |
| --- | --- | --- |
| ❔ | Prepare and submit our team's original PA1 proposal | No original proposal, deadline, receipt, or mark was located in the reviewed files |
| ✅ | Adopt LeeTKOB as the shared product for PA2 onward | [README](../README.md), [PA2 plan](PA2/PA2-PLAN.md), and current application |

## PA2 — harness and quality gate

These rows describe completed engineering and historical verification. They do not mean all checks still pass on today's checkout.

| Status | Job | Evidence / remaining qualification |
| --- | --- | --- |
| ✅ | Set up React, TypeScript, Vite, npm, Go and PostgreSQL | [README](../README.md), `frontend/`, `backend/`, `compose.yaml` |
| ✅ | Document project rules, permission tiers, commands and prohibited actions | [AGENTS.md](../AGENTS.md) |
| ✅ | Keep local configuration ignored and list environment names without values | `.gitignore`, `.env.example`; setup lessons recorded in [AI log](../AI-LOG.md) |
| ✅ | Write the submission endpoint contract | [Submission contract](contract/submissions.md) |
| ✅ | Implement authentication, enrollment checks and submission creation | `backend/internal/api/` |
| ✅ | Store submissions and pending jobs atomically; enforce idempotency | `backend/internal/api/submission.go`, contract tests |
| ✅ | Test contract errors, access rules, replay and concurrency behavior | `backend/internal/api/api_test.go`, [verification](PA2/PA2-H1-H2-VERIFICATION.md) |
| ✅ | Configure local checks and CI for lint, format, types, build and backend tests | [CI workflow](../.github/workflows/gates.yml), [quality-gate verification](PA2/PA2-QUALITY-GATE-VERIFICATION.md) |
| ✅ | Record a failing required check blocking merge and the subsequent fix | [PA2 report](../PA2.md), [quality-gate verification](PA2/PA2-QUALITY-GATE-VERIFICATION.md); historical evidence, not a fresh remote-settings check |
| ✅ | Record assistant use and an evidence-backed PA2 self-assessment | [AI log](../AI-LOG.md), [self-assessment](../SELF_ASSESSMENT_REPORT.md) |
| ✅ | Record later green checks after PA2 | Historical run #56 passed the checks below; this documentation update does not claim a fresh full-suite run on the current checkout |

## PA3 — prototype with the LLM feature running

Detailed ownership and delivery history remain in the [PA3 plan](PA3/PA3-PLAN.md). Some handoff checklists predate completion; use [PA3.md](../PA3.md) and the [PA3 self-assessment](../SELF_ASSESSMENT_REPORT.md) for final evidence and disclosed limitations.

| Status | Job | Evidence / next action |
| --- | --- | --- |
| ✅ | Support admin, TA and student accounts with role/class access checks | `backend/internal/api/`, `frontend/src/pages/` |
| ✅ | Let a TA create a lab with statement, private tests, constraints and deductions | `backend/internal/api/lab.go`, `frontend/src/pages/labs.tsx` |
| ✅ | Let enrolled students submit C++ and view history, source and judging results | `frontend/src/pages/workspace.tsx`, submission/status handlers |
| ✅ | Implement a durable queue and separate worker with recovery and result fencing | `backend/internal/worker/worker.go`, worker tests; [Member 1 handoff](PA3/MEMBER1-HANDOFF.md) |
| ✅ | Compile and run C++ inside a restricted sandbox | `backend/internal/sandbox/runner.go`, `sandbox/Dockerfile` |
| ✅ | Implement no network; CPU, wall, memory/swap, process and output limits | Sandbox runner and `backend/internal/sandbox/testdata/` |
| ✅ | Run as non-root with no host mounts, Docker socket or application secrets | Sandbox container policy and attack fixtures |
| ✅ | Add sandbox attack tests and a Docker-enabled CI job | [Main run #56](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45326541493) passed the sandbox job on `97346d0` |
| ✅ | Run syntax-aware hard rules before semantic checking | `backend/internal/hardrules/`, `sandbox/scan.cpp` |
| ✅ | Test hard-rule traps such as comments, strings and a student's own sort function | `backend/internal/hardrules/check_test.go` |
| ✅ | Call the real model for pending semantic constraints on Accepted code | `backend/internal/semantic/`, [saved worker result](../eval/runs/real-worker/result.json) |
| ✅ | Remove comments and known account identities from model input | `semantic/source.go`, `semantic/check.go`, worker setup; arbitrary unknown personal names remain a documented limitation |
| ✅ | Validate response schema and evidence; turn malformed responses/timeouts into unsure | `backend/internal/semantic/check.go`, fake-model tests |
| ✅ | Keep model findings, TA-set deductions and judge results separate | `backend/migrations/00004_judging.sql`, worker completion logic |
| ✅ | Record model calls and enforce a persistent per-lab application budget | `backend/internal/worker/model.go`, migration 00005, [Member 2 handoff](PA3/MEMBER2-HANDOFF.md) |
| ✅ | Provide a generic checker command and an offline model-off Docker image | `backend/cmd/checker/`, `deploy/Dockerfile.checker`, [eval guide](../eval/README.md) |
| ✅ | Preserve the official dataset/evaluator and add five separate labelled cases | `eval/official-sha256.json`, `eval/extra/`, integrity scripts |
| ✅ | Save a real-model dev-set run, predictions, call logs, metrics and cost estimate | Passing [numbered-line run](../eval/runs/real-dev-numbered/report.txt); original `eval/runs/real-dev/` retained |
| ✅ | Supply deployment files and host-worker instructions | `compose.pa3.yaml`, `deploy/`, [Member 1 handoff](PA3/MEMBER1-HANDOFF.md) |
| ✅ | Complete Member 3 acceptance of both handoffs | [Acceptance record](PA3/MEMBER3-ACCEPTANCE.md): checks on `df6f621`, local walkthrough and upstream hash comparison |
| ✅ | Expose findings through a TA/class-protected review endpoint | `backend/internal/api/review.go`, `server.go`, and [TA review evidence](PA3/TA-REVIEW-QUEUE.md) |
| ✅ | Build the TA review queue with source, constraint, verdict, evidence and TA-set deduction | Review UI, class-scoped API, TA decisions, configured deductions and FAQ-09 unsure reasons; [TA review evidence](PA3/TA-REVIEW-QUEUE.md) |
| ✅ | Confirm and cite product-owner behavior for pre-publication visibility and unsure handling | [FAQ-08, FAQ-09, FAQ-10 and FAQ-12](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md) are cited in the TA review evidence; no Publish endpoint exists in this build |
| ✅ | Rebuild and smoke-test the local deployment after the frontend/proxy changes | [November 24 local smoke](PA3/LOCAL-COMPOSE-SMOKE-2026-11-24.md), model off; subsequent stale lab-cache, earlier-submission and hard-rule-reason fixes are recorded in [AI-LOG](../AI-LOG.md) |
| ✅ | Verify the final public model project/key spending cap | [PA3 model section](../PA3.md#model) and [November 25 cap evidence](PA3/evidence/ai-studio-cap.png); per-lab application budget remains separate |
| ✅ | Demonstrate compliant, violating and injection submissions on the deployment | [Public walkthrough](PA3/evidence/public-walk.json), model on at `2a61adf`: 22/22 checks, including TA2/student access denial |
| ✅ | Add the PA3 rubric's assistant-mistake examples to the rules file's never list | Tracked complete in [PA3 plan](PA3/PA3-PLAN.md); preserve the official data and evaluator |
| ✅ | Write root `PA3.md` and update the self-assessment for PA3 | [PA3 report](../PA3.md), [self-assessment](../SELF_ASSESSMENT_REPORT.md), claimed 92/100 |
| ✅ | Record the deployment walkthrough and verify public availability | [Recording linked in PA3.md](../PA3.md#deployment) is 8:48, exceeding the five-minute requirement; limitation disclosed in self-assessment. Keep deployment running through PA4 marking |
| ✅ | Tag the submitted commit `pa3`, verify tag CI, prepare ZIP and confirm submission | Complete per team; [PA3 report](../PA3.md#ci-for-the-tag) links tag runs. No Classroom receipt or fresh remote CI verification was inspected in this update |

## PA4 — evaluation, guardrails, human gate and red-team

**November 27 backend update:** Member 1's sampling, review requirements and atomic
Publish gate are implemented and tested locally against PostgreSQL. See the
[contract](contract/publication.md) and [handoff](PA4/MEMBER1-HANDOFF.md). Member 3
UI acceptance/deployment and Member 2/3 evaluation/release evidence remain pending;
the November 25 checkpoint overview above is the original planning baseline.

Source: [full PA4 rubric](<PA4/PA#4 rubric — Eval, guardrails, human gate and red-team.md>). Execution and member handoffs: [PA4 plan](PA4/PA4-PLAN.md). Existing guardrail code is useful progress; its complete test matrix, UI behavior and history requirements still need verification.

| Status | Job | Evidence / next action |
| --- | --- | --- |
| ✅ | Implement schema, source-evidence and timeout validation foundations | `backend/internal/semantic/check.go` and fake-model tests; complete rubric coverage remains open below |
| ✅ | Enforce a persistent application cost cap per lab | Worker budget ledger and tests; persisted cap-exhaustion reason is displayed to the TA; full PA4 acceptance remains below |
| ✅ | Reject model deduction fields and keep configured deductions separate | Unknown fields are rejected; model request excludes deductions; TA queue displays configured and TA-applied deductions |
| ✅ | Produce valid generic prediction JSONL without reading labels to decide answers | Checker command and saved evaluation artifacts |
| ✅ | Demonstrate zero wrongful deductions on the saved official dev-set run | 0 wrongful deductions in the passing [numbered-line run](../eval/runs/real-dev-numbered/report.txt); one run, not proof for unseen data |
| ✅ | Add five own labelled/tagged Accepted cases separately from official data | `eval/extra/`; PA4 requires at least ten and the history/coverage checks below |
| ⬜ | Generate fresh real-model predictions inside CI for official and own cases | Every eval exit other than 0 must fail; committed predictions and model-off runs do not meet this requirement; unit tests still use fakes |
| ⬜ | Make the live eval job a required check on main and run it for PRs | Document fork PR handling because forks receive no secrets; keep CI key in repository secrets |
| ⬜ | Prove a red eval blocks merging, then fix the checker | Link the blocked PR and fixing commit; do not change labels/data or suppress failures |
| ⬜ | Extend own cases from five to at least ten, covering trap, injection and disguised | Each case must precede the checker change that handles it; inspect real history without reconstructing it |
| ⬜ | Record three real dev-set runs of the same commit | In PA4.md: changed-verdict pair count, lowest/highest recall and unsure rate, and cost of one CI run |
| ⬜ | Meet all three thresholds on the lecturer's held-out run | 0 wrongful, recall ≥90%, unsure ≤20%; ≥95% recall reaches the top of the held-out band; dev results are only an estimate |
| ⬜ | Validate reliability after the numbered-line evidence fix | Latest saved run passes; earlier provider 503/budget/evidence failures remain useful regression evidence. Preserve safe unsure fallback and confirm stability across fresh runs |
| ⬜ | Audit G1–G5 with meaningful fake-model/no-network tests | Cover fenced/prose JSON, unknown verdict, missing/duplicate constraints, invented/misnumbered/comment-only evidence, timeout, exhausted cap, and injected marks; tests must fail if their guardrail is removed |
| ✅ | Show the TA why a check became unsure, including cap exhaustion | FAQ-09 local reasons are persisted, returned through the TA-only API and displayed in review; [TA review evidence](PA3/TA-REVIEW-QUEUE.md). Complete the G1–G5 audit above |
| ⬜ | Verify no model instruction changes the TA-set deduction on screen | Prove configured deduction survives injected fields/instructions; checker must have no mark-writing path |
| ⬜ | Verify guardrail/test commit history and deployed injection behavior | Each guardrail/test in its own real commit, test no later than code; injections must become violation or unsure, never pass; record historical limitations honestly |
| ✅ | Persist TA decisions independently of model findings | `constraint_reviews` and `constraint_review_history`, reviewer identity and configured deductions; [TA review evidence](PA3/TA-REVIEW-QUEUE.md) |
| ✅ | Require review of every violation and unsure finding | Server Publish gate and direct PostgreSQL API tests; [handoff](PA4/MEMBER1-HANDOFF.md) |
| ✅ | Randomly sample 10% of pass verdicts, rounded up, on the server | Persisted samples, fixed-seed edge tests and restart persistence; [handoff](PA4/MEMBER1-HANDOFF.md) |
| ✅ | Implement Publish as the only action exposing a deduction to students | Atomic revision gate, pending/role/class denials and student visibility tested locally; Member 3 UI/deployment pending; [contract](contract/publication.md) |
| ⬜ | Write PA4.md section “Where we stop it” | Use own numbers: who pays when wrong, how error is detected, where it is stopped; explain Publish placement and the cost of unsampled passes |
| ⬜ | Complete at least eight reproducible red-team attacks across five categories | `docs/red-team/<target-student-IDs>.md`; include at least one code-derived attack, held and broken results, exact input/time/output, rule, impact and proposed fix; lecturer replays two |
| ⬜ | Follow the assigned target and rules of engagement | At most about 1 request/second; one resource attack per kind; trigger cap once; stop on live key/real personal data and notify lecturer/target that day; report an unavailable target by November 30 with timed attempts |
| ⬜ | Rehearse the generic held-out invocation under marking restrictions | 30-minute limit; network only through a proxy to the README's model host; valid predictions required; do not special-case dev solution IDs |
| ⬜ | Prepare private marking credentials and their lifecycle | `marking.env` beside the ZIP, never in Git or inside ZIP; only checker-required variables; dedicated key/endpoint credentials; revoke after PA4 marks; without working access the lecturer uses model-off |
| ⬜ | Complete PA4 evidence, AI log and submission | Root PA4.md: tag CI, blocked PR/fix, repeated runs/cost, human-gate explanation, model host and fork policy; actual Changed/Rejected entries; public `pa4` tag, running deployment, rubric self-assessment and ZIP including .git |
| ⬜ | Ensure the submitted pa4 tag has a CI run | Workflow currently lists only pa2/pa3 tags; extend tag coverage as part of the CI work |

The rubric also highlights CODEOWNERS review for tests, `.github/` and eval data as protection against moving the target; this is not currently configured. Ordinary checks/guardrail tests must still work offline after dependencies are installed within 15 minutes; the separately authorized live evaluation job is the exception that contacts the provider.

### Saved real-model evaluation

Latest saved run: [report](../eval/runs/real-dev-numbered/report.txt), [metrics](../eval/runs/real-dev-numbered/summary.json), [usage](../eval/runs/real-dev-numbered/usage-summary.json), and [call log](../eval/runs/real-dev-numbered/calls.jsonl), recorded November 24 on checker code `1cdc486`. The unchanged evaluator reports **PASS (exit 0)** for all 102 pairs, with no missing pairs.

| Metric | Saved result | Required | Status |
| --- | --- | --- | --- |
| Wrongful deductions | 0 | 0 | ✅ Meets threshold in this run |
| Violations caught with correct evidence | 35/37 = 94.6% | At least 90% | ✅ Meets threshold in this run |
| Unsure verdicts | 2/102 = 2.0% | At most 20% | ✅ Meets threshold in this run |

All 39 calls succeeded and have provider request IDs and known usage: 37,838 input tokens, 3,795 output tokens, estimated **US$0.020851**. This is the cost of this official dev run; PA4 must measure the complete CI run including own cases. Extrapolating 80 similar calls gives approximately US$0.043, as reported in PA3.

The checker now sends source as numbered `{line, text}` objects (PR #22, test `c05230c`, fix `1cdc486`), avoiding model line-count drift while keeping exact evidence validation. Two labelled violations remain missed: `dynamic-array/s05 C2` is unsure and `dynamic-array/s07 C3` is predicted pass. These IDs describe evaluation findings only; they must not appear in production prompts or hard rules.

The earlier [November 22 report](../eval/runs/real-dev/report.txt) remains unchanged: 0 wrongful, 21/37 recall (56.8%), 31/102 unsure (30.4%), evaluator exit 1. Its provider 503 responses lacked request IDs/usage and are documented separately from successful calls; budget blocks and rejected evidence also contributed. The newer passing run is neither the required three runs of one commit nor a live CI or held-out result. PA4 still needs all of that evidence.

## PA5 — final build, running and observed

Source: [full PA5 rubric](<PA5/PA#5 rubric — Final build and operations.md>). Keep the submitted `pa5` deployment available over HTTPS from December 17 through December 23 at 22:00, and then through the assigned PA6 slot. Post-deadline bug fixes do not count as the submitted build.

| Status | Job | Evidence / next action |
| --- | --- | --- |
| ✅ | Update student judging results automatically | Active jobs are polled every 1.5 seconds in `frontend/src/pages/workspace.tsx` |
| ✅ | Provide basic worker logs and detailed model-call audit records | Foundation only; console output alone does not meet the full observability criterion |
| ⬜ | Show the pa5 commit hash through a version endpoint or footer | Deployment and checked-out source must match the submitted tag |
| ⬜ | Document deployment, environment configuration locations and rollback | Root PA5.md; run HTTPS deployment with marking accounts working for every role |
| ⬜ | Complete the six-step deployed publishing/appeals walkthrough | TA creates a lab with hard/semantic constraints; two students submit; results update without reload; TA reviews/publishes; both students see their marks/deductions |
| ⬜ | Implement appeals and documented rejection errors | Refuse appeals before publication and on another student's submission; cite the API contract and relevant product-owner FAQ entries |
| ⬜ | Persist appeal decisions and their consequences | Rejection reason visible to student; acceptance updates published mark; retain who/when/why; reopen only if FAQ allows it |
| ⬜ | Search structured logs for all stages by one submission/correlation ID | Include level, time, service and ID; prevent API keys or student names/IDs inside logged prompts from appearing in accessible logs |
| ⬜ | Trace HTTP → queue → worker → sandbox → hard rules → each model call → stored verdict | Preserve context across the durable queue; duration for each stage; model span identifies model and tokens or cost |
| ⬜ | Provide a dashboard and read-only marking access | Request rate/errors; queue depth/oldest age; judge time; model call count/latency/errors/timeouts; unsure share; spend versus cap |
| ⬜ | Define an SLO and prove its alert has fired | PA5.md names objective and alert; save alert history/notification timestamp matching logs; an untriggered alert is insufficient |
| ⬜ | Build `load/` script and Dockerfile, requiring only URL and credentials | Must run in lecturer's container with no mounts and only network access; exact command in PA5.md |
| ⬜ | Submit once as each of 80 distinct students to one lab within 10 seconds | Keep every sandbox protection enabled; sequential requests or repeated use of one account do not meet the load specification |
| ⬜ | Report complete load metrics and achieve the replay target | Accepted requests, final results, failures/statuses, duplicates, result p50/p95, queue drain time, model reach/cost; target 80/80 final within 15 minutes, zero duplicates/5xx; retried 429 acceptable only with exactly one resulting submission |
| ⬜ | Write `load/RESULT.md` with reproducible measured evidence | Date/time, URL/commit, script output, during-run dashboard screenshot, first saturation point, resulting change/commit and any cap-driven unsure; lecturer reruns once during marking |
| ⬜ | Resolve measured throughput bottlenecks | Current worker serializes jobs and recompiles per test; coordinate per-key provider pacing if scaling workers; explain major p95 differences between own and marking runs |
| ⬜ | Write a real incident post-mortem in `docs/postmortem/YYYY-MM-DD-slug.md` | Seven sections: summary, impact, timeline, detection, root cause, fix, follow-up; evidence for each timeline entry; timezone; wrong-mark impact; fix commit and verification |
| ⬜ | Complete post-mortem follow-ups with ownership and detection improvement | At least one follow-up done and linked; one changes detection; real staging/CI/load events count, intentional fault injection must be declared on the first line; do not invent/rewrite history |
| ⬜ | Add PA5-specific AI-log entries and measured lab cost | Deployment, dashboards/alerts, load script and post-mortem; real Changed/Rejected entries including a caught assistant error; production model and actual lab cost from metrics, not only estimates |
| ⬜ | Record a continuous public-deployment walkthrough after tagging pa5 | At most ten minutes, all six product-flow steps; PA5.md links it; this is the specified downtime fallback |
| ⬜ | Prepare ACCESS.md privately and exclude it from Git | Add to .gitignore; include only in ZIP: marking-role accounts, second student, read-only observability; no real personal passwords; not currently ignored or present |
| ⬜ | Finish PA5 release/submission evidence | Public pa5 tag; PA5.md links URL/version, deployment/rollback, FAQ, tools, SLO/alert, load command/report, post-mortem and recording; updated rubric self-assessment and correctly named ZIP including .git |
| ⬜ | Retain final tests/CI evidence and monitor availability through marking | Workflow tag coverage currently stops at pa3; address it if relying on tag CI; deployment is tested twice on different days, at least six hours apart |

## PA6 — final demo and individual oral

Source: [full PA6 rubric](<PA6/PA#6 rubric — Final demo and individual oral.md>). Each member receives the same **40 build points** plus their own **60 explanation points**. The build must be the one submitted for PA5; this is not a new implementation milestone. There is no PA6 ZIP, self-assessment, or self-scored total. The only posted deliverable is the commit link for OWNERS.md.

### Eleven-feature demonstration checklist (B1: 22 points)

Each feature earns 2 if working live, 1 if partial/second-try/recording-only, or 0 if missing/faked. The statuses below describe readiness, not claimed B1 points. Implemented foundations remain checked above; a feature stays open here until its full required demo is verified.

| Status | Feature | Current evidence / remaining demo work |
| --- | --- | --- |
| ⬜ | F1 — Accounts, roles and classes | Existing access checks; verify role screens plus another class's lab/review-queue denial after review integration |
| ⬜ | F2 — Lab with tests, hard/semantic constraints and TA deductions | Authoring exists; create the demo lab live in the TA's class on pa5 |
| ⬜ | F3 — Queue, separate judge worker and real-time results | Implemented; verify queued-to-final without reload on the final deployment |
| ⬜ | F4 — Sandbox | Nine protections/attack fixtures and CI job exist; retain current CI evidence and rehearse attack submissions |
| ⬜ | F5 — Syntax-aware hard rules first | Existing implementation/trap tests; show a trap and call log proving no model call for decided bans |
| ⬜ | F6 — Real semantic model, privacy, evidence, no model deduction | Real-call evidence exists; verify live pa5 input sanitization, findings and unchanged TA deductions |
| ⬜ | F7 — Guardrails and reasons for unsure | Validation/budget foundations exist; TA reasons, deployed injection behavior and full five-guardrail test coverage remain open |
| ⬜ | F8 — Review queue and Publish | Missing queue/API/UI, pass sampling and server Publish gate |
| ⬜ | F9 — Published marks and appeals | Publication and appeals pending |
| ⬜ | F10 — Checker and required eval gate | Generic command/model-off exist; required live-model CI gate missing |
| ⬜ | F11 — Logs, metrics, traces, SLO and alert | Basic logs/model ledger exist; full observability pending; show trace of a submission created during the demo |

### Team preparation and build checks

| Status | Job | Completion evidence / requirement |
| --- | --- | --- |
| ⬜ | Commit OWNERS.md and post its commit link by December 21 at 22:00 | Two tables: subsystem → actual member IDs → commits/PRs; beyond brief → where visible → checking test; authorship table guides questions, not points |
| ⬜ | Decide whether to target B5's four beyond-brief points before PA5 freezes | One working, tested extra feature in submitted pa5; identify an existing candidate or implement only if core work permits; FAQ-required behavior is not an extra; empty table means no B5 claim |
| ⬜ | Keep the pa5 deployment running through the assigned oral slot | Slot/table posted by December 21; December 24 or 31, 08:20–11:50, cs2:G202; arrive at session start |
| ⬜ | Verify the submitted build and real persistence for B2 | Version hash matches pa5 and laptop `git log -1 --format=%H`; reload at two requested points preserves data; Network tab shows actual API calls; disclose configuration/seed changes |
| ⬜ | Rehearse the live main flow within eight minutes | Fresh lab, compliant/violating submissions, early appeal denial, pending-review Publish denial, review/publish, marks, appeal decision, logs and trace; demo data created in front of examiners, existing accounts/classes may be seeded |
| ⬜ | Rehearse two examiner-selected server failure cases for B3 | T1, Q1, M1–M9, D1, G4/G5, H1–H3, A1–A3; exact refusal/error and unchanged state after reload; G1–G3 are demonstrated by tests rather than selected running-build cases |
| ⬜ | Rehearse direct role/class/student ownership denial for B4 | Attempt URL or curl request against another role/class/student; API must refuse |
| ⬜ | Prepare role windows, CLI/tests and equipment | Admin, own-class TA, two students and other-class TA in separate signed-in windows; observability one click away; pa5 checkout, tests/eval/checker ready; charger/HDMI adapter, backup laptop, phone hotspot and student ID |
| ⬜ | Rehearse a cold start and the documented fallback | Switch failed laptop; deployment has three minutes, then questions proceed and one retry at session end; retry may run pa5 locally; if still failed, PA5 recording plus enforcement code/tests yields only limited build credit |

### Individual explanation (60 points per member)

Documentation already exists in the [backend guide](architecture/BACKEND.md), [database guide](architecture/DATABASE.md), and [mechanism walkthroughs](architecture/BACKEND-CONNECTIONS.md). Each member must find and explain the actual source/evidence independently; reading documentation is not completion evidence.

| Status | Job | Completion evidence to collect |
| --- | --- | --- |
| ⬜ | A — Explain own code and commits, plus the demo path | Locate files/lines/diffs, describe what they actually do and what breaks if removed; cover teammate components too |
| ⬜ | B — Defend decisions and actual alternatives | Use own measurements/limits/incidents for sandbox, polling, contract, queue and FAQ choices; explain costs and when the decision would change |
| ⬜ | C — Show exact evidence and trace the incident | Locate a meaningful test/CI/eval/log/trace; explain what it catches, assistant mistakes caught, and detection-to-root-cause evidence |
| ⬜ | D — Predict and run the checker on an examiner-provided input | Explain hard-rule routing, model constraints, verdict/evidence, guardrail, review queue and TA action; run once with own key and explain any difference |
| ⬜ | Explain why the human gate sits at Publish and who pays when wrong | Use actual deductions, lab/course weight, unsampled-pass risk and own evaluation numbers |
| ⬜ | Rehearse independent navigation and answers without the assistant | Files/search/git/tests/eval/dashboards may stay open; assistant closed during the member's turn; teammates remain quiet; no writing code or memorizing source required |

For three members, plan about 23 minutes total: one-minute version check, eight-minute flow/features, two-minute failure/ownership checks, then about four minutes of questions per member. Load testing is scored in PA5, not rerun in the PA6 slot. The separate seminar on December 10/17 still needs its own preparation. Oral explanation affects individual PA6 points, not earlier PA2–PA5 marks.

## Current verification snapshot

The table retains the explicitly recorded GitHub Actions [run #56](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45326541493) at `97346d0` as historical automated verification. It is not the latest-head or `pa3` tag result. Final PA3 evidence is in [PA3.md](../PA3.md) and the [self-assessment](../SELF_ASSESSMENT_REPORT.md); that report links tag runs. This documentation update inspected local source and saved artifacts, without freshly querying remote CI or the deployment. The current workflow still uses model-off evaluation; PA4's required live gate remains open.

| Check | Result | Qualification / action |
| --- | --- | --- |
| Go format, vet and unit tests | ✅ Passed | `gates` job in run #56 |
| Frontend lint, formatting, type-check and production build | ✅ Passed | `npm run check --prefix frontend` in the `gates` job |
| Official dataset integrity | ✅ Passed | `verify-eval` step in run #56; local manifest lists 128 official files |
| Docker sandbox attack tests and offline checker | ✅ Passed | `sandbox` job in run #56; checker ran offline with semantic results left unsure |
| Browser end-to-end workflows | ✅ Passed | `e2e` job in run #56 |
| Full-history secret scan and dependency audit | ✅ Passed | `secret-history` job in run #56 |
| Local Compose smoke | ✅ Recorded | Separate user-run smoke on `1122f7c`, `MODEL_OFF=1`; see [smoke record](PA3/LOCAL-COMPOSE-SMOKE-2026-11-24.md) |
| Public deployment and live-model walkthrough | ✅ Recorded | [Public API walkthrough](PA3/evidence/public-walk.json), 22/22 checks at `2a61adf`; URL, recording and cap evidence in [PA3.md](../PA3.md). No fresh availability probe in this update |
| Saved real-model dev evaluation | ✅ Pass | `real-dev-numbered`: 0 wrongful, 94.6% recall, 2.0% unsure; 39 successful provider-identified calls; one saved run |
| PA3 release checks | ✅ Complete per team | Root report/self-assessment and tag-run links recorded; 8:48 video exceeds the limit. Submission receipt was not independently inspected |
| PA4 live CI / Publish / red-team | ⬜ Open | See [PA4 plan](PA4/PA4-PLAN.md); no completion inferred from PA3 evidence |

## Recommended work order

| Priority | Work package | Finish when |
| --- | --- | --- |
| 1 | Establish PA4 contracts and history before fixes | Agree Publish/sampling behavior, audit G1–G5, commit own cases before checker changes, obtain the assigned target |
| 2 | Complete PA4 trust and evidence requirements | Follow [member handoffs and dates](PA4/PA4-PLAN.md): fresh required real-model CI; ten cases/three runs; guardrail audit; Publish/random sampling; eight-attack report; held-out credentials and evidence package |
| 3 | Complete PA5 workflows and reproducible operations | Publication/appeals; traces across queue; dashboard, SLO and fired alert; 80 students/10 seconds/15-minute result; incident follow-ups; versioned HTTPS release and private access package |
| 4 | Prepare PA6 demo and explanations before PA5 freezes | Choose/test any B5 extra in pa5; OWNERS.md posted by December 21; live F1–F11, failure cases, persistence and independent evidence-based answers rehearsed |

**Immediate schedule:** PA4 is due in one week, December 2 at 22:00. Use the dated [PA4 delivery plan](PA4/PA4-PLAN.md#delivery-order-and-checkpoints), with an integrated candidate and repeated-run evidence by December 1. Report an unavailable assigned red-team target by November 30. Provider variability, publication concurrency and deployment integration remain the immediate uncertainties; PA5 operations and PA6 rehearsal follow after this checkpoint.
