# PA3 local Compose smoke — 2026-11-24

This is a local Ubuntu WSL / Docker Desktop smoke run of merged commit `1122f7c7865da11257b50091641eaca17e6cd136` (PR #18). It is **not** evidence of public deployment or of live model checking. The checkout was `~/leetkob-pa3-smoke-1122f7c`; Compose project `leetkobpa3smoke` used web port 18088 and PostgreSQL host port 5433. The worker ran separately in the Ubuntu terminal with `MODEL_OFF=1` and `SANDBOX_IMAGE=leetkob-sandbox:pa3smoke-1122f7c`.

## Observed results

| Check | Observation |
| --- | --- |
| Stack | PostgreSQL, API and Caddy web containers healthy; migration exited successfully. `/` and `/classes` returned 200; unauthenticated `/api/me` returned 401. |
| Caddy navigation | Direct refresh of `/labs` loaded the app. Admin, TA and student test accounts signed in through the web UI. |
| Class ownership | Admin assigned two different TAs to separate classes. The other-class TA saw only their own class; a direct request to the first class's lab review page returned “The lab was not found.” |
| Lab and judging | The first TA created a lab with two private tests and two constraints (ban `std::sort`, deduction 3; use a comparison, deduction 2). Student submissions moved from Queued to Accepted after the separate worker started. |
| Hard-rule review | Latest counted `std::sort` submission `669782e6-2920-4b75-9c3f-006993e7ed32` showed a hard-rule violation with exact source line 7; the semantic constraint was `unsure (semantic_pending)` because the model was off. Before TA action, the selected deduction was 0. The TA selected line 7 and confirmed the violation; the recorded TA deduction became 3. |
| Student boundary | Before and after the TA review, the student submission page showed judge Accepted and no constraint deduction or published mark. The build has no Publish/marks route; a future publication gate is untested. |
| Worker restart | With the worker stopped, submission `a36e05d3-f456-44df-b078-b651caed9892` remained Queued. On manual worker restart it reached Accepted. This UI observation does not prove exactly one verdict set in the database; the separate [TestDurableWorker integration test](MEMBER1-HANDOFF.md) covers queued recovery, one persisted result set, and stale-lease refusal. |
| Injection fixture | A compliant program with an injection comment also reached Accepted under `MODEL_OFF=1`. No live-model prompt-injection resistance is claimed. |

The browser briefly showed “Lab not available” immediately after lab creation, then the lab appeared in the TA list and the student could open it. This transient behavior needs reproduction before attributing a cause or closing it.

## Still needed for PA3 marking

- Run the rubric's exact `sort-merge` dev-set lab on a **public deployment**, with compliant, violation and injection submissions and another-class TA access check. The local two-integer lab above is a smoke fixture, not that lab.
- Verify the final provider project/key spending cap before `MODEL_OFF=0`; run and inspect semantic findings on the deployment. The committed [real dev run](../../eval/runs/real-dev/) is separate evidence and reports 0 wrongful deductions, 56.8% recall and 30.4% unsure (evaluator exit 1).
- Record a tag `pa3` CI run including sandbox attacks; publish the deployment URL, test accounts and recording in `PA3.md`; finish evidence-backed self-assessment and ZIP submission.

Product behavior follows [FAQ-02, FAQ-08, FAQ-09, FAQ-10 and FAQ-12](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md). The PA3 rubric's severe boundary applies: a model finding must never become a student deduction or published mark without TA review.
