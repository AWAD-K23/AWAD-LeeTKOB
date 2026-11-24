# Member 3 acceptance of the Member 1 and Member 2 handoffs — 2026-11-24

**Commit checked:** `df6f621` (main after PR #21), exported with `git -c core.autocrlf=false archive`
so the files are the committed LF bytes. **Where:** Windows 11, Docker Desktop; Go checks ran in
`golang:1.26` with the host Docker socket (the sandbox runner pipes code through stdin and uses
no bind mounts), PostgreSQL 17 in a separate container, frontend checks in
`mcr.microsoft.com/playwright:v1.65.0-noble`.

## Results

| Check (from the handoff) | Command | Result |
|---|---|---|
| Go format | `gofmt -l backend` | no files listed |
| Go vet | `go -C backend vet ./...` | passed |
| Go tests with PostgreSQL and real Docker attacks | `TEST_DATABASE_URL=… SANDBOX_TESTS=1 go -C backend test ./... -count=1 -v` | all packages `ok`; no `--- FAIL` |
| Nine sandbox minimums | `TestAttacks/{network,cpu,wall,memory,process,output,user,mounts,secrets}` | passed (17.3 s) |
| Queue: worker stopped, submit, restart, stored once | `TestDurableWorker/queued_work_survives_absent_worker_and_is_stored_once`, `…/expired_lease_fences_old_worker` | passed |
| Hard-rule traps | `TestSyntaxTraps`, `TestStatementAndTypeBans`, `TestBanListFormat` | passed |
| Schema and timeout limits with a fake model | `TestInvalidResponsesBecomeUnsure`, `TestTimeoutBecomesUnsure`, `TestUncertaintyReasons` | passed |
| Prompt redaction | `TestCommentRemoval`, `TestKnownIdentityRedaction` | passed |
| Model off opens no provider | `TestOfflineNeverConfiguresProvider` | passed |
| TA review isolation | `TestSubmissionContract` (incl. TA review subtest) | passed |
| Official dev set unchanged | `node scripts/verify-eval.mjs` | "Verified 128 unchanged official files" |
| Own cases compile and are Accepted | `node scripts/verify-extra.mjs` | 5 solutions with `g++ -std=c++17 -O2` |
| Checker image, model off, official set | `node scripts/run-offline-eval.mjs eval/official …` | 102 predictions; wrongful 0, recall 40.5% (15/37), unsure 46.1%; evaluator exit 1 (thresholds, not format) — same as `eval/runs/offline/` |
| Checker image, model off, own cases | `node scripts/run-offline-eval.mjs eval/extra …` | 10 predictions; evaluator exit 1 (semantic pairs stay unsure) |
| Frontend lint, format, types, build | `npm run check --prefix frontend` | passed (bundle-size warning only) |
| Browser workflows | `npm run e2e --prefix frontend` | 9 passed |

## Dev-set hashes against the published copy

- `eval/official/eval.mjs` SHA-256 starts `c836edd8713521bb` and `eval/official/labels.jsonl` starts
  `262cda7b0ef5612e`: both equal the prefixes printed in the PA#3 rubric ("Catching your assistant").
- A fresh clone of `github.com/fithcmus/awad-leetkob-eval` (2026-11-24) and `eval/official/` at
  `df6f621` differ in no file (`diff -rq`, `.git` excluded).

## Local rehearsal of the criterion 1 flow

`scripts/pa3-walkthrough.mjs` ran against `compose.pa3.yaml` built from `df6f621`
(Caddy on `http://localhost:18088`, worker from `cmd/worker`, `MODEL_OFF=1`), with five throwaway accounts:

- `setup`: the admin created class A (TA1 and two students) and class B (TA2).
- `walk`: TA1 created the `sort-merge` lab from `eval/official/problems/sort-merge` (statement, 7 tests,
  C1 deduction 3, C2 deduction 4). With **no worker running**, student1 submitted `s01.cpp` (compliant),
  student2 `s04.cpp` (calls `std::sort`), student1 `s09.cpp` (injection); all three jobs stayed `pending`.
  The worker was started 25 s later; all three became Accepted. Each submission has exactly one
  `judge_results` row and two `constraint_verdicts` rows.
- TA1's queue: `s04` C1 `violation/hard_rule` with evidence line 11 `sort(a.begin(), a.end());`;
  C2 `unsure` (model off); `s09` C1 `pass/hard_rule`, C2 `unsure`. Student status responses carry no
  finding or deduction fields.
- Refusals: TA2 → queue 404, lab not in TA2's list, TA2 creating a lab in class A → 404;
  a student → queue 403; student2 reading student1's submission → 404.

## Findings to resolve before the public walkthrough

1. **The queue hides a student's earlier submission.** Under FAQ-02 only the counted (latest)
   submission per student enters the queue, so `s01` (student1) disappeared when student1 submitted
   `s09`. The rubric's marker submits three solutions with the two student accounts it is given and
   then looks for "each submission". FAQ-10 also lets a TA open every submission outside the queue.
2. **"Lab not available" after creating a lab — reproduced by reading the code.** `LabForm`
   (`frontend/src/pages/labs.tsx`) calls `invalidateQueries({ queryKey: ["labs"] })`, which refetches
   only *active* queries; the form page has no observer of `["labs"]`, so the cached list stays stale.
   `LabWorkspace` (`frontend/src/pages/workspace.tsx`) then renders from that stale list, finds no lab,
   shows "Lab not available", and recovers when its own refetch returns. It happens only when the
   list was cached earlier in the session (the TA visited Labs first).
3. **A hard-rule verdict shows the reason "Semantic model check has not run."** (e.g. `s09` C1
   `pass/hard_rule`). The text is misleading for a TA on a constraint the model never handles.
   Cause: `backend/internal/hardrules/check.go` starts every verdict with that reason and keeps it
   when a ban list decides the constraint (`v.Source = "hard_rule"` does not clear `v.Reason`).

Neither handoff is blocked by these: every documented check passes. Items 1–3 are recorded here
and raised with their owners before any fix.

## The marker's checker run, reproduced

As criterion 6 describes, with a fresh clone of the published dev set mounted over the image's `/data`:

```text
docker build -f deploy/Dockerfile.checker -t leetkob-checker:pa3 .
docker run --rm --network none --memory 1g --pids-limit 128 -e MODEL_OFF=1 \
  -v <published-dev-set>:/data:ro leetkob-checker:pa3 > predictions.jsonl
node <published-dev-set>/eval.mjs --data <published-dev-set> --predictions predictions.jsonl
```

102 predictions in 14 s; the published `eval.mjs` reports 0 wrongful, recall 40.5%, unsure 46.1%,
exit **1** (thresholds), not 2 (format). Re-scoring the committed `eval/runs/real-dev/predictions.jsonl`
with the published `eval.mjs` gives 0 wrongful, 56.8% (21/37), 30.4% (31/102), exit 1: the numbers we report.
