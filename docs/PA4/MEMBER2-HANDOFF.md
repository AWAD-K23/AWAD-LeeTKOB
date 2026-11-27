# PA4 Member 2 handoff — candidate, live fix validation pending

November 27, 2026 (Asia/Saigon). This is a partial acceptance record with explicit
failing performance evidence. Do not present the checker candidate as passing.

## Delivered

- Re-scored the saved PA3 baseline without model calls and freshly compared all
  128 official files with the lecturer's pinned archive: zero mismatches.
- Ten Accepted own solutions with twenty labelled/tagged pairs. New cases were
  committed in `f1510fa` before checker changes. All ten compile and pass their
  problem tests. Existing cases and labels are unchanged.
- G1–G5 regression audits in separate genuine commits, listed in [PA4.md](../../PA4.md).
  Production guards already existed in mixed PA3 history; no history was rebuilt.
  Both [baseline mutation evidence](evidence/guardrail-mutations.json) and
  [candidate mutation evidence](evidence/guardrail-mutations-after-fix.json) show
  all five normal regressions passing and their bypassed guards causing failures.
- Fresh official-and-own live CI for PR/main/pa4, with failure propagation and
  provider/cost artifacts. Model-off checking and fake-model unit tests remain
  separate. Fork PRs fail rather than silently count as passing without secrets.
- Active main ruleset requires live-eval, the original three checks, one approval
  and CODEOWNERS review. The key is in the approved team Actions secret, never
  in Git or a student execution container. CODEOWNERS takes effect on main once
  its file is merged; its reviewers all have write access. Settings snapshots are
  saved under [evidence](evidence/repository-protection-after.json).
- Three same-commit real dev runs, valid generic image rehearsal, variability
  comparison tools, model-host/proxy and private marking-variable instructions.

## Actual results and candidate fix

All three dev runs of `77ef5b0` failed: one wrongful violation each, recall
89.2–94.6%, unsure 0–2.9%. Seven unique pairs changed verdict. All 117 successful
responses have provider IDs and known usage, costing an estimated US$0.062328.
Raw predictions, reports, provider logs and configuration are preserved under
`eval/runs/pa4-dev-{1,2,3}/`. No prior run was overwritten.

[PR #27](https://github.com/AWAD-K23/AWAD-LeeTKOB/pull/27) preserves the failing
baseline. Its original-budget combined live evaluation is the fourth approved
evaluation. The candidate fix is local commit `60bfe73`, currently unpushed:
generic guidance follows actual reachable computation and independent rule scope,
and invalid evidence no longer discards a separate constraint's valid evidence.
Malformed overall schema/IDs still reject the response; every invalid individual
finding stays unsure. The new regression failed before the fix and passes after
it. Existing guardrail tests were not weakened. Production prompts contain no
official problem/solution IDs or quoted dev code.

Go format/vet/full PostgreSQL tests passed after the fix. The complete Docker
sandbox suite and frontend check also passed during this work. Five Node tooling
tests verify comparison correctness, metadata/population rejection, explicit
live mode and refusal to overwrite evidence. No fresh passing real-model result
for `60bfe73` is claimed; additional live runs are awaiting explicit approval.

## Acceptance commands

Run the usual format/vet/backend/frontend checks from README. With an isolated
`TEST_DATABASE_URL`, repeat mutations into a new evidence file:

```text
node scripts/audit-guardrails.mjs <new-evidence.json>
node --test scripts/*.test.mjs
node scripts/verify-eval.mjs
node scripts/verify-extra.mjs
```

On Windows, Git's LF attributes now keep CSS and own cases consistent in fresh
checkouts. The earlier extra-case verifier failure was CRLF conversion of test
outputs; verification used the unchanged committed bytes in a temporary Git
export. Neither labels nor the verifier was changed to make that check pass.

For an authorized real run, set the model variables privately and use a new folder:

```text
node scripts/run-live-eval.mjs eval/runs/<new-run> eval/official
node scripts/run-live-eval.mjs eval/runs/<new-combined-run>
node scripts/compare-model-runs.mjs <run1> <run2> <run3>
```

The runner refuses dirty checker/tooling/case state and existing output folders.
It logs the exact commit, model, budget, image and dataset hash. Run evaluations
sequentially so separate processes do not compete for the same provider RPM cap.
The trusted native checker holds the key; isolated syntax-analysis containers
receive no key, host mount, Docker socket or network. Unit tests never use it.

## Remaining

1. Preserve the actual combined CI result, artifact/cost and merge-block evidence.
2. After extra-run approval, push the candidate fix and obtain a green required
   evaluation; record three fresh runs of one fixed commit if authorized. Failing
   baselines must remain disclosed; never change data/labels or weaken the gate.
3. Supply a dedicated private marking key/credentials outside Git and the ZIP;
   verify proxy access in the marking environment and revoke after marking.
4. Member 3 accepts the handoff and supplies deployed UI/injection checks, target
   red-team report, final rubric report and release evidence. No such acceptance
   or deployment is claimed here.
