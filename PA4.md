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
its value was neither printed nor committed. Actual CI/block/fix links follow.

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
without provider calls. One full official-and-own CI cost remains to be measured.

[Diagnostic findings](docs/PA4/evidence/baseline-findings.md) describe generic
failure mechanisms. No labels, official source or thresholds were altered.
These baseline failures require a checker fix and new validation; they do not
support a passing held-out estimate.

The offline-tested candidate fix is local commit `60bfe73`. It adds generic
reachable-computation/rule-scope guidance and isolates invalid evidence per
finding without accepting any unsupported decisive verdict. All original
guardrail regressions remain unchanged and pass; the
[second mutation audit](docs/PA4/evidence/guardrail-mutations-after-fix.json) still
catches all five bypasses. The candidate is not pushed or claimed to pass real
evaluation; extra paid validation remains unapproved. See the
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
