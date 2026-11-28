# PA4 Member 2 handoff - subsequent live CI failure

> Current update: CI run 46036066453 failed with one wrongful official finding. A local independent-confirmation candidate is offline-tested only; see the [latest report](../../PA4.md#subsequent-ci-failure-and-local-confirmation-change) and [compact audit](evidence/confirmation-repeat-failure.json). Earlier green evidence below is historical.

November 27, 2026 (Asia/Saigon). The latest checker passes the required combined
CI evaluation. Earlier failures remain preserved; latest-code repeated-run
stability, held-out performance and Member 3 acceptance are not claimed.

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
  [Confirmation audit](evidence/confirmation-guard-mutations.json) repeats these
  five checks and also verifies that bypassing the new confirmation gate fails.
  [Provider-schema audit](evidence/provider-schema-guard-mutations.json) verifies
  all seven guards, including the local follow-up's decisive-evidence bound.
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
baseline. Its original-budget combined live evaluation failed on own cases at
8/9 recall (88.9%), while the official set passed. The required check blocked
merge; [saved state](evidence/baseline-pr-blocked.json) and
[CI artifacts](../../eval/runs/pa4-ci-baseline/run-summary.json) preserve proof.
All 49 successful calls had provider IDs/known usage and cost an estimated
US$0.023860, with no unknown reservations or retries.

The first correction is commit `60bfe73`:
generic guidance follows actual reachable computation and independent rule scope,
and invalid evidence no longer discards a separate constraint's valid evidence.
Malformed overall schema/IDs still reject the response; every invalid individual
finding stays unsure. The new regression failed before the fix and passes after
it. Existing guardrail tests were not weakened. Production prompts contain no
official problem/solution IDs or quoted dev code.

Go format/vet/full PostgreSQL tests passed after the fix. The complete Docker
sandbox suite and frontend check also passed during this work. Five Node tooling
tests verify comparison correctness, metadata/population rejection, explicit
live mode and refusal to overwrite evidence.

The three authorized dev runs on `61d9c89` finished: two passed, but one still had
a wrongful violation. All caught 36/37 violations; unsure was 1/102, 1/102 and
0/102. Two pairs changed verdict; all 117 calls had known usage/provider IDs,
costing US$0.068578. [Comparison](evidence/fixed-variability.json) and raw run
folders preserve the failure. These are not three passing runs.

The latest correction `cd3374d` requires an additional strictly validated call
to confirm semantic violations. It attempts to disprove the proposal; disagreement
or any confirmation failure becomes unsure. Both calls use the existing cap,
audit and pacing controls; comments/known identities remain removed in the
confirmation request. Positive, disagreement, malformed, cap/provider error,
timeout, wire and identity tests pass, as does the full Docker/PostgreSQL suite.
The final approved [CI run](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45987925317)
passed official metrics (0 wrongful, 35/37 recall, 5/102 unsure) but failed own
recall (8/9, with 0 wrongful and 1/20 unsure). A primary model violation lacked
evidence; the checker correctly left it unsure. Merge remains blocked. All other
jobs passed. [Saved artifacts](../../eval/runs/pa4-ci-confirmation/run-summary.json)
contain 75 known-usage responses with provider IDs: 49 primary and 26 confirmation.
Estimated cost: US$0.041354; unknown reservations/retries zero. The newly approved
four-run batch cost US$0.109932. No fifth evaluation was started.

Local follow-up `e592bf7` aligns the provider schema with the unchanged local
evidence guard: pass/violation require a nonempty evidence array, while unsure
can stay empty. A real JSON Schema engine tests the actual wire schema and
counterexamples offline. The test failed before the correction and passes now;
all seven guard bypasses fail their regressions. Go format/vet/unit tests pass.
After the four-run batch, the user said “continue” and one further CI evaluation
was executed. [Run 46009732306](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/46009732306)
passed all five jobs. Official metrics: 0 wrongful, 36/37 recall (97.3%), 1/102
unsure (1.0%). Own metrics: 0 wrongful, 9/9 recall (100%), 0/20 unsure.
[Saved call/cost summary](evidence/schema-ci-summary.json) and
[complete artifact](../../eval/runs/pa4-ci-schema/run-summary.json) contain all
77 provider IDs/known usages: 49 primary calls and 28 confirmation calls. Cost
US$0.045214 estimated; unknown reservations, provider errors and retries zero.
Tested head `07d14e0` includes schema fix `e592bf7`; Actions used merge SHA
`87f1dd8eae8b86f177a48f83b14a397a947a18e6`. No cases, labels or thresholds changed.

PR #27 remains a draft requiring human review. No merge, deployment or additional
paid run followed this result. This is one passing run of the latest method,
not three repeated passing runs or the lecturer's held-out result. Final evidence
updates are committed locally without a push to avoid an unauthorized extra CI.

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

1. Member 3 reviews and accepts the supplied evidence. The required live gate is
   green, while human review/draft state still blocks merging.
2. If repeated-run stability of the latest method is required, authorize new
   runs first. Existing same-commit series remain honestly reported, including
   their wrongful results. Do not alter labels or weaken the gate.
3. Supply a dedicated private marking key/credentials outside Git and the ZIP;
   verify proxy access in the marking environment and revoke after marking.
4. Member 3 accepts the handoff and supplies deployed UI/injection checks, target
   red-team report, final rubric report and release evidence. No such acceptance
   or deployment is claimed here.
