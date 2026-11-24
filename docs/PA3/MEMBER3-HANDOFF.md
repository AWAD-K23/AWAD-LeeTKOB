# Member 3 handoff — PA#3

**Status checked:** November 24, 2026
**Current baseline:** [`1122f7c`](https://github.com/AWAD-K23/AWAD-LeeTKOB/commit/1122f7c7865da11257b50091641eaca17e6cd136) (PR #18)
**PA#3 deadline:** Wednesday, November 25, 2026 at 22:00 ICT

This guide shows what is ready, what the evidence proves, and what still needs to happen before submission. It does not claim that Member 3 accepted the handoff, that a public deployment is running, or that PA#3 has been submitted.

## Read this first

The core implementation and local checks are in place. The key release evidence is still missing: the public model-on walkthrough, final `pa3` tag CI, `PA3.md`, final self-assessment, and submission ZIP.

| Area | What is verified | What this does not prove |
| --- | --- | --- |
| TA review queue | [Review queue evidence](TA-REVIEW-QUEUE.md): class-scoped API and UI, per-constraint evidence, TA decisions, and deductions. | No Publish or marks endpoint exists yet. |
| Student visibility | In the local smoke, students saw judge status but no constraint findings or deductions before or after TA review. | This is not evidence for a future Publish flow. |
| Worker and sandbox | [Member 1 handoff](MEMBER1-HANDOFF.md) records the durable worker and nine sandbox protections. The [local smoke](LOCAL-COMPOSE-SMOKE-2026-11-24.md) observed queued work complete after a worker restart. | The browser observation alone does not prove exactly one stored result set; the worker integration test covers persistence. |
| Local deployment | The smoke record covers sign-in, page refresh, two classes, cross-class denial, TA review, and worker restart. | It used `MODEL_OFF=1`, local ports, and a small test lab. It is not the required public `sort-merge` flow or a live-model demonstration. |
| Real model evaluation | [Saved dev-set run](../../eval/runs/real-dev/) reports 0 wrongful deductions, 56.8% recall, and 30.4% unsure; the evaluator exited 1. | These are honest results, but they do not meet the brief's model-quality thresholds. |

## Finish in this order

### 1. Accept the implementation handoffs

Read the [PA3 plan](PA3-PLAN.md), [Member 1 handoff](MEMBER1-HANDOFF.md), [Member 2 handoff](MEMBER2-HANDOFF.md), [TA review evidence](TA-REVIEW-QUEUE.md), and the [local smoke record](LOCAL-COMPOSE-SMOKE-2026-11-24.md). Run their documented checks against the final commit and record the results before marking either handoff accepted.

### 2. Prepare the public deployment

Follow the deployment steps in the [README](../../README.md). Set up public HTTPS, persistent storage, the API, web app, and trusted worker. Verify the spending cap on the exact provider project and key before enabling `MODEL_OFF=0`. Keep the key on the server and out of the browser and student-code sandbox.

### 3. Run the required public walkthrough

Use the required admin, two TAs in different classes, and two students in the first TA's class. Have the TA create the dev-set `sort-merge` lab with its statement, tests, two constraints, and deductions. Submit compliant, violating, and injection solutions. Check the verdicts and evidence, TA decisions and deductions, student visibility, page refresh, sign-in, and cross-class denial through a direct request. Check worker restart recovery with the integration test or database evidence.

The local smoke recorded a brief “Lab not available” message immediately after lab creation. Reproduce it and record the result before calling it resolved.

### 4. Verify the live model path

After the provider cap is confirmed, test model-on submissions and inspect the prompt, provider call log, schema/evidence validation, timeout behavior, and TA-visible `unsure` reasons. Keep the existing real dev run and report its results without presenting it as a passing evaluation. The random 10% sample of `pass` findings is tracked as a PA#4 follow-up in the current plan.

### 5. Finish release and submission evidence

- Compare the official dev-set and evaluator files with the published copies; record the full hash results.
- Run final checks and confirm the `pa3` tag CI includes sandbox attacks and passes.
- Add root-level `PA3.md` with the public URL, required test accounts, tag CI link, checker and Docker commands, model-off switch, committed dev-run folder, and recording link.
- Record a deployment walkthrough of no more than five minutes.
- Update `SELF_ASSESSMENT_REPORT.md` with all eight criteria, evidence, an honest total, and “what we did not manage”.
- Tag the submitted commit `pa3`, push the tag, and verify the repository is public and tag CI passed.
- Create the repository ZIP including `.git`; use student IDs in ascending order and a filename total that matches the report. One member submits it from their student account.
- Keep the deployment running until PA#4 marks are out.

## Final review before submission

- Confirm the public walkthrough, deployment, accounts, recording, report, tag, CI run, and ZIP all refer to the same commit.
- Confirm the model key and test accounts are protected, and the public spending cap is active.
- Confirm findings stay in the TA review flow; no model verdict becomes a student deduction without TA review.

## Product behavior references

The current implementation follows [FAQ-02, FAQ-08, FAQ-09, FAQ-10, and FAQ-12](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md). In particular, students do not see constraint findings before publication, `unsure` findings require TA action, and the review queue includes a random sample of `pass` findings. The sample is tracked for PA#4 in the current plan.
