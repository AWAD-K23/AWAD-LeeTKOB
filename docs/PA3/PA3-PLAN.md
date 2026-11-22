# PA3 plan — three-member handoff

**Due:** Wednesday, November 25, 2026, at 22:00.  
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

**Frontend follow-up completed:** Dedicated role-aware pages, account selection by email, class details, lab authoring, submission history/source, Markdown statements, C++ highlighting, and page transitions are implemented. See [manual flow test](MANUAL-FLOW-TEST.md). A local database inspection confirmed an Accepted submission also had a separate `hard_rule` violation with exact source evidence. TA review UI and real model checking remain Members 3 and 2 work respectively.

**Handoff ownership:** Member 1 has published the interfaces and verification steps. Actual acceptance by Members 2 and 3 remains in their own first checklist items; no teammate agreement is claimed here.

**Deployment follow-up owner:** Member 3 owns rebuilding and smoke-testing the current
Docker deployment after the frontend redesign. Member 1 retains ownership of fixes
to the deployment files or backend components they supplied.

## Member 2 — model checking and evaluation

**In progress:** The user confirmed Member 1 acceptance checks were already run.
The semantic boundary, Gemini adapter, worker integration, fake-model guardrails,
and persistent per-lab cost ledger are implemented. Real calls remain off.
See [Member 2 progress](MEMBER2-PROGRESS.md) for limits and pending dataset/evaluation
work. Checklist items requiring real-call evidence remain open.

- [ ] Accept Member 1's handoff and run its checks before building on it.
- [ ] Add a real model call for semantic constraints. Send only the statement, constraints, and code with comments and identity removed. Return a verdict and evidence per constraint; never let the model set deductions.
- [ ] Add schema validation and a time limit. Invalid output and timeouts become `unsure`; test both with a fake model and no network.
- [ ] Log provider request IDs, model, tokens, time, and cost per call without exposing secrets.
- [ ] Add a checker command and Dockerfile. With `MODEL_OFF=1`, decide ban lists, return `unsure` for semantic constraints, and make no network calls. Accept a dataset folder without problem-specific setup.
- [ ] Preserve the published dev set and evaluator unchanged. Add at least five labelled, tagged cases separately, including a trap and an injection.
- [ ] Run the real model on the dev set and save predictions, the evaluator report, and call log. Report all three metrics honestly; calculate the cost for 80 submissions in `README` and `AI-LOG.md`.
- [ ] **Handoff:** Demonstrate the checker on a real Accepted submission, the offline model-off image, fake-model guardrail tests, and the saved real-run evidence. Provide the verdict and evidence interface to Member 3.

## Member 3 — review, deployment walkthrough, and submission

- [ ] Accept both handoffs and run their documented checks before integration.
- [ ] Rebuild and smoke-test `compose.pa3.yaml` after the frontend redesign and `/api` proxy change. Check frontend loading and direct page refreshes through Caddy, `/api/me` authentication, sign-in, class/lab access, and a real submission processed by the separate worker. Record the tested commit and results; return Member 1 implementation defects to Member 1.
- [ ] Connect verdicts to the TA review queue. Show each submission, constraint, verdict, evidence, and TA-selected deduction. Enforce class ownership on the server.
- [ ] Keep verdicts out of published marks. If Publish exists, block it while any verdict is unreviewed. Follow and cite the product owners' FAQ for student visibility before publication.
- [ ] Run the deployed end-to-end flow with the dev-set `sort-merge` lab and compliant, violation, and injection submissions. Verify another class's TA cannot access the lab or queue.
- [ ] Run the documented checks and verify tag CI includes the sandbox attacks. Check that offline checks work without network and that the checker image works with `MODEL_OFF=1`.
- [ ] Add the rubric's “Catching your assistant” mistakes to the rules file's never list. Check that the dev set and evaluator remain unchanged.
- [ ] Create `PA3.md` with the deployment URL, required test accounts, tag CI link, checker and Docker commands, model-off switch, real-run folder, and recording link.
- [ ] Record a deployment walkthrough of at most five minutes. Keep the deployment running until PA4 marks are released.
- [ ] Update `SELF_ASSESSMENT_REPORT.md` with one evidence-backed row per criterion, claimed points, and unfinished work. Agree on an honest total with the team.
- [ ] Tag the submitted commit `pa3`, push the tag, and confirm the repository is public and the tag CI passed.
- [ ] Create the submission ZIP including `.git`, named with student IDs in ascending order and the claimed total. One member submits it before the deadline.

## Final team check

- [ ] Review the complete deployed flow and evidence together before submission.
- [ ] Confirm no model verdict becomes a student deduction without TA review, and schema validation and timeout are active.
- [ ] Confirm no secrets are committed or exposed to the browser or sandbox. Set a spending cap for the public test accounts.
- [ ] Confirm the ZIP, assessment total, tag, deployment, recording, and CI links all match the submitted version.


