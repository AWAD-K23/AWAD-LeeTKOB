# PA#3 — Prototype with the LLM feature running

Team: 23120131 Phạm Bảo Huy · 23120135 Trần Anh Khoa · 23120143 Tô Thanh Long
Repository (public): https://github.com/AWAD-K23/AWAD-LeeTKOB · submitted tag: `pa3` → `TODO: commit`

## Deployment

**URL:** `TODO: https://<dns-name>` — serves `TODO: commit` (the `pa3` tag). It stays up until the
PA#4 marks are out. How it is built and run: [public deployment runbook](docs/PA3/PUBLIC-DEPLOYMENT.md).

Test accounts. They hold no real data and anyone who reads this repository can use them; the
model project has a monthly spending cap and each lab a US$1 model budget (see *Model* below).

| Role | Class | Email | Password |
|---|---|---|---|
| Admin | — | `pa3-admin@leetkob.test` | `TODO` |
| TA 1 | Class A — Sorting lab | `pa3-ta1@leetkob.test` | `TODO` |
| TA 2 | Class B — other TA | `pa3-ta2@leetkob.test` | `TODO` |
| Student 1 | Class A | `pa3-student1@leetkob.test` | `TODO` |
| Student 2 | Class A | `pa3-student2@leetkob.test` | `TODO` |

**Walking through criterion 1.** Sign in as TA 1 → *Classes* → Class A → *New lab*: paste the
`sort-merge` statement, its tests and its two constraints with deductions 3 and 4. Sign in as a
student → *Labs* → the lab → submit. The status goes from *Queued* to the judge result without a
reload. As TA 1 → the lab → *Review* shows each counted submission, each constraint, the verdict
and source (`hard_rule` / `model`), the evidence lines, the reason for any `unsure`, and the
deduction the TA chooses. TA 2 gets "not found" for the lab and its queue (direct request:
`GET /api/labs/<labId>/review-queue` → 404).

> **Which submissions the queue shows.** Following [FAQ-02](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md#faq-02),
> only each student's latest submission counts and enters the queue; an earlier submission of the
> same student leaves it. `TODO: say how the marker sees all three cases (see open decision in the
> Member 3 handoff).`

The same flow through the API, with the results we got: `node scripts/pa3-walkthrough.mjs setup|walk`
(usage in the runbook) → `TODO: docs/PA3/evidence/public-walk.json`.

**Recording (≤ 5 min, the flow of criterion 1 on the deployment):** `TODO: link`

## CI for the tag

`TODO: link to the GitHub Actions run for tag pa3` — jobs `gates`, `secret-history`, `e2e`,
`sandbox`. The `sandbox` job builds the sandbox image and runs `TestAttacks` (one attack submission
per minimum: network, cpu, wall, memory, process, output, user, mounts, secrets; sources in
`backend/internal/sandbox/testdata/`), `TestDurableWorker` (worker stopped → submit → restart →
one stored result set), and the model-off checker on the dev set and on our own cases.

## Checker command

```text
docker build -f deploy/Dockerfile.checker -t leetkob-checker:pa3 .
docker run --rm --network none -e MODEL_OFF=1 -v <dataset-folder>:/data:ro leetkob-checker:pa3 > predictions.jsonl
```

- The dataset folder has the dev-set layout (`problems/<id>/statement.md`, `constraints.json`,
  `solutions/<id>/*.cpp`). The checker reads only statements, constraints and solutions, never
  `labels.jsonl`, and needs no per-problem setting.
- **Model-off switch:** `MODEL_OFF=1` (the image's default). Ban-list constraints are decided by the
  hard rules; every semantic constraint is written as `unsure`; no network connection is opened.
- Without Docker: `go -C backend run ./cmd/checker --data <dataset-folder> --predictions <file.jsonl>`.
- Real model: `MODEL_OFF=0` with `MODEL_NAME` and `MODEL_API_KEY` from the ignored `.env`, plus
  `--call-log <new-file.jsonl>` ([eval/README.md](eval/README.md#real-model-evaluation)).

## Committed dev-set run (real model)

Folder: [`eval/runs/real-dev/`](eval/runs/real-dev/) — `predictions.jsonl`, the unchanged
`eval.mjs` report `report.txt` and `summary.json`, and the call log `calls.jsonl` (provider request
id, model, tokens, time and estimated cost per call; the five HTTP 503 calls have no id and are
marked unknown, not invented).

| Wrongful deductions | Violations caught | Unsure | `eval.mjs` exit |
|---:|---:|---:|---:|
| 0 | 56.8% (21/37) | 30.4% (31/102) | 1 (thresholds missed; the output is valid) |

Re-scored on 2026-11-24 with a fresh copy of the published `eval.mjs`: the same three numbers
([acceptance record](docs/PA3/MEMBER3-ACCEPTANCE.md#the-markers-checker-run-reproduced)). PA#3 has no
pass mark on these; they do not meet the PA#4 thresholds yet.

Dev set: `eval/official/` is byte-identical to `github.com/fithcmus/awad-leetkob-eval`
(`eval.mjs` SHA-256 `c836edd8713521bb…`, `labels.jsonl` `262cda7b0ef5612e…`; all 128 files in
`eval/official-sha256.json`, checked by `node scripts/verify-eval.mjs` in CI). Our own cases are in
[`eval/extra/`](eval/extra/), same layout, labelled and tagged (trap and injection included).

## Model

`gemini-3.5-flash-lite` (Google AI Studio, project AWAD-LeeTKOB, monthly cap ₫100,000 —
`TODO: re-check date`). It receives only the statement, the semantic constraints and the code with
comments and known identities removed, and returns `pass` / `violation` / `unsure` with evidence
per constraint; it never sets a deduction. One lab of 80 submissions costs about **US$0.036**
(US$0.013204 for 29 measured calls ÷ 29 × 80; arithmetic in the
[Member 2 handoff](docs/PA3/MEMBER2-HANDOFF.md#measured-cost)).

## Where each criterion's evidence is

| # | Criterion | Evidence |
|---|---|---|
| 1 | Flow on the deployment | URL and accounts above; recording; `docs/PA3/evidence/public-walk.json` |
| 2 | Queue and sandbox | `backend/internal/sandbox/runner_test.go` (`TestAttacks`), `backend/internal/worker/worker_test.go` (`TestDurableWorker`), tag CI run |
| 3 | Hard rules first | `backend/internal/hardrules/`, `sandbox/scan.cpp` (libclang), `TestSyntaxTraps`; model-off run [`eval/runs/offline/`](eval/runs/offline/) (0 wrongful, all 55 ban-list pairs decided) |
| 4 | Model on semantic constraints | `backend/internal/semantic/`, [`eval/runs/real-dev/`](eval/runs/real-dev/) |
| 5 | Limits and review queue | `TestInvalidResponsesBecomeUnsure`, `TestTimeoutBecomesUnsure` (fake model, no network); [TA review queue](docs/PA3/TA-REVIEW-QUEUE.md); no Publish route exists |
| 6 | Checker and first eval | *Checker command* above; `deploy/Dockerfile.checker`; `eval/extra/` |
| 7 | Tests first, one part per commit | `git log` on `main`; `AI-LOG.md` |
| 8 | AI-LOG | [`AI-LOG.md`](AI-LOG.md) |
