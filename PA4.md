# PA4 evidence — draft

Release evidence and Member 3 integration are still pending. No held-out result,
deployed UI behavior, tag CI result or submission is claimed here.

## CI evaluation and forks

The `live-eval` job in [gates.yml](.github/workflows/gates.yml) generates fresh
predictions for the official dev set and all ten own cases. It uses the real
model, runs the unchanged evaluator on each, fails on every nonzero evaluator
exit and uploads predictions, provider logs, metrics and cost even on failure.
The separate sandbox job retains model-off evaluation; its expected semantic
threshold failure does not satisfy this gate. Unit tests use fakes, never a key.

Fork PRs fail `live-eval` explicitly before secret-bearing execution. A maintainer
must review the changes and run them on a same-repository branch before merge;
secretless or model-off runs cannot substitute for a passing live gate. The
workflow does not use `pull_request_target`. The CI key belongs exclusively in
the team repository's `MODEL_API_KEY` Actions secret and the trusted checker
process, never a student execution container.

With explicit user approval, active ruleset `26947349` now requires `live-eval`
alongside the existing gates/secret-history/e2e checks, one approval and CODEOWNERS
review. [Before](docs/PA4/evidence/repository-protection-before.json) and
[after](docs/PA4/evidence/repository-protection-after.json) API snapshots preserve
the settings. Existing restrictions and the empty bypass list were retained.
The selected team repository's `MODEL_API_KEY` Actions secret is configured;
its value was neither printed nor committed.

[Baseline CI run](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45968971182)
failed the required `live-eval` check on [PR #27](https://github.com/AWAD-K23/AWAD-LeeTKOB/pull/27).
The official set passed (0 wrongful, 94.6% recall, 2.0% unsure); the own set
missed one of nine violations (88.9% recall), so the evaluator exited 1 and merge
was blocked. [Exact blocked state](docs/PA4/evidence/baseline-pr-blocked.json) and
[downloaded artifacts](eval/runs/pa4-ci-baseline/run-summary.json) preserve the
result. No data, labels, check failures or branch protection were changed to
obtain this evidence. The generic checker correction is `60bfe73`; its green
required-check evidence is pending validation.

## Guardrails and case history

The PA3 baseline re-scores without model calls at 0 wrongful, 35/37 recall (94.6%)
and 2/102 unsure (2.0%); all 128 official hashes match the manifest. This is a
historical baseline, not a fresh or held-out run. Prior upstream comparison is
retained in [Member 3's PA3 acceptance](docs/PA3/MEMBER3-ACCEPTANCE.md).

Five new Accepted own cases were committed in `f1510fa`, before subsequent checker
changes. All ten compile and pass their output tests. Labels and reasoning are
in [eval/extra](eval/extra/README.md); none are embedded in production prompts or
rules. The two missed pairs in the saved baseline require semantic interpretation;
no case-specific production fix or relaxed evidence matching has been introduced.

| Guard | Focused regression commit | Test |
| --- | --- | --- |
| G1 schema | `ed609fe` | `TestG1SchemaFailsClosed` |
| G2 evidence | `4f85bdf` | `TestG2ExactEvidenceAndOriginalMapping` |
| G3 timeout | `3305f2a` | `TestG3DeadlineDiscardsLateSuccess` |
| G4 lab cap | `a4f7d10` | `TestDurableWorker/G4` |
| G5 deductions | `55b7c85` | `TestG5ModelCannotSupplyMarks`, `TestPublicationGate` |

Each new audit was committed separately with its test. Existing PA3 guards were
already implemented together; no history has been split or rebuilt to claim the
rubric's original guardrail-per-commit condition. All five mutation checks passed:
normal code passed its regression and a temporary copy with the relevant guard
bypassed failed. [Mutation evidence](docs/PA4/evidence/guardrail-mutations.json)
records the exact commit and failure output. Run
`node scripts/audit-guardrails.mjs <new-evidence-file>` with `TEST_DATABASE_URL` to
repeat the audit in a disposable export without editing the checkout.

## Three runs, variability and cost

Three fresh official dev runs used frozen commit `77ef5b0`, identical dataset,
model/version/configuration and sandbox image. All **failed**, rather than being
represented as passing evidence:

| Run | Wrongful | Recall | Unsure | Estimated provider cost |
| --- | ---: | ---: | ---: | ---: |
| [1](eval/runs/pa4-dev-1/official/report.txt) | 1 | 33/37 (89.2%) | 3/102 (2.9%) | US$0.020726 |
| [2](eval/runs/pa4-dev-2/official/report.txt) | 1 | 34/37 (91.9%) | 3/102 (2.9%) | US$0.020148 |
| [3](eval/runs/pa4-dev-3/official/report.txt) | 1 | 35/37 (94.6%) | 0/102 (0%) | US$0.021454 |

Seven unique `(solution, constraint)` pairs changed verdict across any of the
three runs. Lowest/highest recall: **89.2% / 94.6%**; unsure: **0% / 2.9%**.
[Comparison output](docs/PA4/evidence/baseline-variability.json) gives the exact
definition and changed pairs. Every run contains predictions, reports, metadata,
usage summaries and all 39 provider IDs. There were no unknown-usage calls or
unfinished reservations; the three runs total US$0.062328 in estimated usage.
Run `node scripts/compare-model-runs.mjs <run1> <run2> <run3>` to reproduce this
without provider calls.

The baseline full CI evaluation cost **US$0.023860 estimated**: 21,126 microdollars
official plus 2,734 own cases. All 49 successful responses have provider IDs and
known usage; unknown-usage reservations, unfinished calls and retries were zero.
This is provider token-cost estimation, not an invoice or GitHub runner cost.
The uploaded artifact is `pa4-live-eval-45968971182-1`; its metadata records the
Actions merge SHA, while the PR head was `df15711`.

[Diagnostic findings](docs/PA4/evidence/baseline-findings.md) describe generic
failure mechanisms. No labels, official source or thresholds were altered.
These baseline failures require a checker fix and new validation; they do not
support a passing held-out estimate.

The first correction (`60bfe73`, tested on frozen `61d9c89`) adds generic
reachable-computation/rule-scope guidance and isolates invalid evidence per
finding. Its three new dev runs were:

| Run | Wrongful | Recall | Unsure | Estimated provider cost |
| --- | ---: | ---: | ---: | ---: |
| [1](eval/runs/pa4-fixed-dev-1/official/report.txt) | 0 | 36/37 (97.3%) | 1/102 (1.0%) | US$0.022579 |
| [2](eval/runs/pa4-fixed-dev-2/official/report.txt) | 0 | 36/37 (97.3%) | 1/102 (1.0%) | US$0.022260 |
| [3](eval/runs/pa4-fixed-dev-3/official/report.txt) | 1 | 36/37 (97.3%) | 0/102 (0%) | US$0.023739 |

Two pairs changed verdict; recall min/max **97.3% / 97.3%**, unsure **0% / 1.0%**.
[Comparison](docs/PA4/evidence/fixed-variability.json) verifies identical commit,
model/version, dataset and image. All 117 calls had known usage and provider IDs,
costing US$0.068578 total. This series still failed the zero-wrongful requirement;
two passing runs do not justify hiding the third.

The latest correction `cd3374d` adds a separate confirmation call for proposed
semantic violations. It tries to disprove the proposed finding and keeps a
violation only when the second strictly validated response agrees. Disagreement,
invalid evidence/output, timeout or budget failure becomes unsure. Sanitized
proposal text prevents original comments and known identities from re-entering
the prompt. Hard-rule findings and semantic passes do not incur another call.
Both calls use the same rate, audit and budget controls; this is not a retry.
The confirmation tests cover positive/negative decisions, budget and provider
failures, cancellation, wire framing and identity redaction. Existing call-count
assertions intentionally require two calls; their verdict/evidence/isolation
assertions were retained. Model consensus can still share an error and is not
claimed to be a correctness proof.

The original guardrail regressions still pass; the
[second mutation audit](docs/PA4/evidence/guardrail-mutations-after-fix.json) still
catches all five bypasses. The latest
[mutation audit](docs/PA4/evidence/confirmation-guard-mutations.json) also removes
the confirmation gate and requires its regression to fail. There are not three
repeated real runs of this correction within the authorized run budget.

The final approved [confirmation CI run](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/45987925317)
finished with:

| Dataset | Wrongful | Recall | Unsure | Evaluator |
| --- | ---: | ---: | ---: | --- |
| [Official](eval/runs/pa4-ci-confirmation/official/report.txt) | 0 | 35/37 (94.6%) | 5/102 (4.9%) | PASS |
| [Own](eval/runs/pa4-ci-confirmation/extra/report.txt) | 0 | 8/9 (88.9%) | 1/20 (5.0%) | FAIL |

The missed own finding was a violation with **empty evidence** from the primary
model. Local validation correctly made it unsure; it did not reach confirmation.
The required check stayed red and [merge stayed blocked](docs/PA4/evidence/confirmation-ci-blocked.json).
All four ordinary CI jobs passed. Its 75 calls (49 primary plus 26 confirmation)
all have provider IDs/known usage. Cost was **US$0.041354 estimated**, with no
unknown reservations/retries: US$0.035837 official plus US$0.005517 own cases.
The four newly authorized evaluations total **US$0.109932**, and no fifth was run.

Local follow-up `e592bf7` fixes a provider/local schema mismatch: the provider
schema previously allowed empty evidence for decisive verdicts. It now uses
conditional branches requiring at least one evidence item for pass/violation,
while preserving empty-evidence unsure. Google documents these supported
[structured-output features](https://ai.google.dev/gemini-api/docs/structured-output).
A regression checks the actual serialized schema using a real JSON Schema
validator; it failed before the fix and passes after. Existing local validation
is unchanged. [Seven mutation checks](docs/PA4/evidence/provider-schema-guard-mutations.json)
pass normally and fail when each guard is bypassed, including this schema bound.
After the user said “continue,” the single further
[schema-fix CI evaluation](https://github.com/AWAD-K23/AWAD-LeeTKOB/actions/runs/46009732306)
**passed all five jobs**, including required live-eval on both datasets:

| Dataset | Wrongful | Recall | Unsure | Evaluator |
| --- | ---: | ---: | ---: | --- |
| [Official](eval/runs/pa4-ci-schema/official/report.txt) | 0 | 36/37 (97.3%) | 1/102 (1.0%) | PASS |
| [Own](eval/runs/pa4-ci-schema/extra/report.txt) | 0 | 9/9 (100%) | 0/20 (0%) | PASS |

Cost: **US$0.045214 estimated**, comprising US$0.038785 official and US$0.006429
own cases. All **77 calls** (49 primary, 28 confirmation) have provider IDs and
known usage; there were no provider errors, unknown reservations or retries.
The tested PR head was `07d14e0`, Actions merge SHA
`87f1dd8eae8b86f177a48f83b14a397a947a18e6`; provider-schema fix `e592bf7` is included.
[Call/cost summary](docs/PA4/evidence/schema-ci-summary.json),
[all-green checks](docs/PA4/evidence/schema-ci-checks.json) and the complete
[run artifact](eval/runs/pa4-ci-schema/run-summary.json) preserve this result.
The PR remains draft and requires a human owner review; no merge or deployment
was performed. Earlier failed required-check runs remain linked above.

This is **one passing combined CI run** of the latest schema, not three repeated
passing runs or a held-out result. The two earlier same-commit studies and all
their failures remain disclosed. No further paid run was started after this
single follow-up. Final evidence updates are committed locally without a push
that would trigger another paid evaluation. See the
[Member 2 handoff](docs/PA4/MEMBER2-HANDOFF.md).

## Marking command and model host

Allow HTTPS only to **`generativelanguage.googleapis.com:443`**, through the marker's
proxy. The generic checker accepts any folder with `problems/` and `solutions/`:

```text
checker --scanner /usr/local/bin/scan --data /path/to/arbitrary-dataset --predictions /writable/predictions.jsonl --call-log /writable/calls.jsonl
```

Build `deploy/Dockerfile.checker`; this command runs inside the trusted checker
image. It analyzes code without executing it. Model credentials are supplied only
to the checker process; its scanner uses a cleared environment. No Docker socket,
host mount or secrets belong in a student execution sandbox. Go honors proxy
environment variables, and redirects never forward authentication.

[Offline image rehearsal](eval/runs/pa4-marking-model-off/rehearsal.json) used
an arbitrary dataset path, no network/credentials and all twenty own pairs. It
completed in 2.8 seconds with valid output and evaluator exit 1 (expected offline
semantic threshold failure). It proves command validity, not live proxy access
or held-out performance. [Fresh upstream comparison](docs/PA4/evidence/official-upstream-comparison.json)
also confirms all 128 pinned official files match the lecturer's archive.

Private `marking.env` variable names: `MODEL_OFF`, `MODEL_NAME`, `MODEL_API_KEY`,
and optionally `MODEL_LAB_BUDGET_MICRO_USD`. Use a dedicated working key; attach
the file beside the ZIP and revoke the key after marking. No private marking file
is created or committed here. The model must be enabled with `MODEL_OFF=0`.

## Where we stop it

Pending Member 3: explain the error/fairness cost using the team's own measured
numbers and configuration, and demonstrate the deployed Publish/student flow.
The [Member 1 contract](docs/contract/publication.md) and
[handoff](docs/PA4/MEMBER1-HANDOFF.md) describe the implemented server gate.

## Remaining release evidence

Member 3 still owns deployed injection/UI checks, assigned-target red-team work,
the rubric self-assessment, tag CI, ZIP and submission. Current backend handoff
and checker checks do not establish those results.
