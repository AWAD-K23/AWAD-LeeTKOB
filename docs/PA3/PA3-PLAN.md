# PA3 plan — three-member handoff

**Due:** Sunday, November 29, 2026, at 22:00 (moved from November 25; Classroom announcement of November 24).  
**Goal:** Run the complete LeeTKOB flow on the deployment and provide the evidence needed for marking.

## How the handoff works

- Member 1 is you. Assign names to Members 2 and 3 before starting. Member 1 owns the largest share, including the hard rules.
- Work in the order below. The next member starts their implementation after the previous handoff passes its checks.
- Each member owns their part through its tests, documentation, and fixes. At handoff, demonstrate the listed result and agree on the interfaces. If a later issue appears, return it to the owner; do not silently redo their work.
- Keep tests before or in the same commit as the code. Make separate commits for sandbox, queue, hard rules, model call, and limits. Do not rewrite history to simulate this order.
- Record real assistant use and actual changed or rejected suggestions in `AI-LOG.md` throughout the work.

## Member 1 — submission, queue, sandbox, hard rules, and deployment setup

- [x] Inspect the current implementation and publish fixed job, result, and constraint interfaces for Members 2 and 3. Receiving members confirm acceptance in their handoff tasks below.
- [x] Ensure a TA can create a lab with its statement, tests, constraints, and deductions in their class.
- [x] Ensure a student can submit C++ and see the judging status and result; enforce class access on the server.
- [x] Add a persistent queue and separate worker. Test stopping the worker, submitting, restarting it, and storing exactly one set of verdicts.
- [x] Compile and run code inside the sandbox with all nine protections: no network; CPU, wall time, memory including swap, process, and output limits; non-root user; no host mounts or Docker socket; no secrets.
- [x] Add a meaningful attack submission for each protection. Read allocated memory back in the memory test; ensure wall timeout kills the container.
- [x] Set up CI for the project and sandbox tests on the `pa3` tag; document local commands.
- [x] Set up the local deployment environment for the application, queue, worker, sandbox, database, and server-side model configuration; document deployment steps.
- [ ] Verify that the documented provider spending cap applies to the final model project/key used for public testing before enabling real calls. Existing saved cap evidence does not prove final public provisioning is complete.
- [x] Write hard-rule tests, including comments, strings, a student's own `sort`, and an `<algorithm>` include used for `std::swap`.
- [x] Implement syntax-aware ban-list checks for constraints starting with “Do not use” or “Do not call”. Do not call the model for constraints decided by hard rules.
- [x] **Handoff:** Show a real submission moving from queued to judged, the worker restart and hard-rule tests, and a working deployment environment. Provide the Accepted submission, constraint data, hard-rule interface, and deployment instructions.

**Evidence:** See [Member 1 handoff](MEMBER1-HANDOFF.md) for code, tests, interfaces, and local deployment commands. The original PostgreSQL/Docker checks passed during Member 1 delivery. Later verification passed Go format/vet/PostgreSQL integration checks, eight frontend workflow tests, and the additional Markdown/C++ rendering test in a targeted run. Frontend lint, formatting, type-checking, and production build passed. Docker attack tests were not rerun during the frontend changes. CI is configured for the final `pa3` tag; no remote tag run or public hosting is claimed.

**Frontend follow-up completed:** Dedicated role-aware pages, account selection by email, class details, lab authoring, submission history/source, Markdown statements, C++ highlighting, and page transitions are implemented. See [manual flow test](MANUAL-FLOW-TEST.md). A local database inspection confirmed an Accepted submission also had a separate `hard_rule` violation with exact source evidence. TA review UI and real model checking were delivered by Members 3 and 2 respectively; public deployment verification remains open.

**Handoff ownership:** Member 1 has published the interfaces and verification steps. Actual acceptance by Members 2 and 3 remains in their own first checklist items; no teammate agreement is claimed here.

**Deployment follow-up owner:** Member 3 owns rebuilding and smoke-testing the current
Docker deployment after the frontend redesign. Member 1 retains ownership of fixes
to the deployment files or backend components they supplied.

## Member 2 — model checking and evaluation

**Local delivery complete:** See [Member 2 handoff](MEMBER2-HANDOFF.md) for real
checker/worker evidence, usage, costs, and interfaces. The real evaluation returned
0 wrongful deductions, 56.8% recall, and 30.4% unsure; the original evaluator
returned 1. These scores do not meet PA4 thresholds. Provider 503 errors, budget
blocks, and rejected evidence are recorded honestly. Member 3 acceptance and final
public provisioning remain separate tasks.

**Offline evidence:** The checker image ran with `--network none`, memory/process
limits, and `MODEL_OFF=1`. It produced 102 official predictions, decided all 55
ban-list pairs, and caused zero wrongful deductions. All 47 semantic pairs stayed
unsure. The unchanged evaluator returned 1 (40.5% recall, 46.1% unsure), honestly
recorded under `eval/runs/offline/`. Five separate labelled/tagged cases were
verified Accepted with `g++ -std=c++17 -O2` and checked with the same command.

- [x] Accept Member 1's handoff and run its checks before building on it. User confirmed the checks were already run.
- [x] Add a real model call for semantic constraints. Send only the statement, constraints, and code with comments and identity removed. Return a verdict and evidence per constraint; never let the model set deductions.
- [x] Add schema validation and a time limit. Invalid output and timeouts become `unsure`; test both with a fake model and no network.
- [x] Log provider request IDs, model, tokens, time, and cost per call without exposing secrets. Fake-transport tested and saved real logs supplied; HTTP 503 responses without IDs/usage are explicitly recorded as unknown.
- [x] Add a checker command and Dockerfile. With `MODEL_OFF=1`, decide ban lists, return `unsure` for semantic constraints, and make no network calls. Accept a dataset folder without problem-specific setup.
- [x] Preserve the published dev set and evaluator unchanged. Add at least five labelled, tagged cases separately, including a trap and an injection.
- [x] Run the real model on the dev set and save predictions, the evaluator report, and call log. Report all three metrics honestly; calculate the cost for 80 submissions in `README` and `AI-LOG.md`.
- [x] **Handoff:** Demonstrate the checker on a real Accepted submission, the offline model-off image, fake-model guardrail tests, and the saved real-run evidence. Provide the verdict and evidence interface to Member 3. Local demonstrations are recorded; Member 3 acceptance remains below.

## Member 3 — review, deployment walkthrough, and submission (Phạm Bảo Huy, 23120131)

**Handoff status and remaining work:** [Member 3 handoff](MEMBER3-HANDOFF.md).


- [x] Accept both handoffs and run their documented checks before integration. All checks passed on `df6f621`, plus a local `sort-merge` walkthrough; three findings are listed for their owners. See [Member 3 acceptance](MEMBER3-ACCEPTANCE.md).
- [x] Rebuild and smoke-test `compose.pa3.yaml` after the frontend redesign and `/api` proxy change. The November 24 local smoke on `1122f7c` checked Caddy page refreshes, `/api/me`, sign-in, class isolation, TA review, and a real worker-processed submission/restart with `MODEL_OFF=1`; see [smoke record](LOCAL-COMPOSE-SMOKE-2026-11-24.md). This is not public deployment or live-model evidence. Reproduce the transient “Lab not available” message.
- [x] Connect verdicts to the TA review queue. Show each counted submission and constraint with the checker verdict, model evidence, TA decision, and applied deduction; enforce class ownership on the server. See [TA review queue verification](TA-REVIEW-QUEUE.md).
- [x] Keep model findings out of student marks unless a TA reviews them. The current build has no Publish or marks route; constraint reviews stay separate from student submission responses. If publishing is added, block it while any verdict is unreviewed. Follow [FAQ-08, FAQ-09, FAQ-10, and FAQ-12](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md). See [TA review queue verification](TA-REVIEW-QUEUE.md). PA4's required random 10% pass sampling is still a follow-up and is not implemented in this PA3 build.
- [ ] Run the deployed end-to-end flow with the dev-set `sort-merge` lab and compliant, violation, and injection submissions. Verify another class's TA cannot access the lab or queue.
- [ ] Verify final branch checks and the submitted `pa3` tag CI include sandbox attacks. GitHub Actions [run #56](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45326541493) passed all four jobs on main `97346d0`, including gates, browser end-to-end, full-history secret scan, and Docker sandbox/offline checker. A `pa3` tag run is still required.
- [x] Add the rubric's “Catching your assistant” mistakes to the rules file's never list.
- [x] Compare the official dev-set and evaluator hashes with the supplied materials and record the result. `eval.mjs` and `labels.jsonl` match the rubric prefixes and a fresh clone of the published repository differs in no file; see [Member 3 acceptance](MEMBER3-ACCEPTANCE.md#dev-set-hashes-against-the-published-copy).
- [ ] Create `PA3.md` (draft at the root; URL, passwords, tag CI and recording still TODO) with the deployment URL, required test accounts, tag CI link, checker and Docker commands, model-off switch, real-run folder, and recording link.
- [ ] Record a deployment walkthrough of at most five minutes. Keep the deployment running until PA4 marks are released.
- [ ] Update `SELF_ASSESSMENT_REPORT.md` with one evidence-backed row per criterion, claimed points, and unfinished work. Agree on an honest total with the team.
- [ ] Tag the submitted commit `pa3`, push the tag, and confirm the repository is public and the tag CI passed.
- [ ] Create the submission ZIP including `.git`, named with student IDs in ascending order and the claimed total. One member submits it before the deadline.

## Final team check

- [ ] Review the complete deployed flow and evidence together before submission.
- [ ] Confirm no model verdict becomes a student deduction without TA review, and schema validation and timeout are active.
- [ ] Confirm no secrets are committed or exposed to the browser or sandbox. Set a spending cap for the public test accounts.
- [ ] Confirm the ZIP, assessment total, tag, deployment, recording, and CI links all match the submitted version.




**Current main CI:** [Run #56](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45326541493) passed all four jobs on `97346d0` (gates, browser end-to-end, secret history, and Docker sandbox/offline checker). **Local Compose smoke:** [2026-11-24 record](LOCAL-COMPOSE-SMOKE-2026-11-24.md) covers merged commit `1122f7c`, Caddy, sign-in, class ownership, TA review, and worker restart with `MODEL_OFF=1`. It is not the rubric’s public `sort-merge` model-on flow. The brief “Lab not available” message after creation needs reproduction.
