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

Repository settings, actual CI runs and blocked-merge/fixing-commit links will be
recorded after inspection and execution. CODEOWNERS covers tests, `.github/`,
scripts and eval data; the file alone is not proof that required review is active.

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

Pending: freeze one commit and record three real official dev runs plus one full
official-and-own CI-equivalent run. The user confirmed the provider cap and
authorized these runs. Every run uses a new folder and preserves prior evidence.
Compare them with `node scripts/compare-model-runs.mjs <run1> <run2> <run3>`.
Changed pairs means unique `(solution, constraint)` pairs with differing verdicts
across any of the three runs. Commit, model/configuration and dataset must match.
Unknown provider usage and unfinished reservations remain separately disclosed.

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
