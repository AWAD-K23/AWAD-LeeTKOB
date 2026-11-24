# Evaluation files

| Folder | Meaning |
|---|---|
| `official/` | Unchanged lecturer dataset and evaluator |
| `extra/` | Our separate labelled cases |
| `runs/` | Saved checker predictions and evaluation reports |

- Upstream: https://github.com/fithcmus/awad-leetkob-eval
- Pinned commit: `552c811db2d35a23d641bc4fdffc1ac120be1a28`.
- `official-sha256.json` records all 128 original files.
- Run `node scripts/verify-eval.mjs` to check the snapshot.
- The checker never reads `labels.jsonl`; only the evaluator reads labels.
- No official file, expected label, or evaluator rule may be edited to improve a result.

## Commands

From the repository root:

```text
docker build -f deploy/Dockerfile.checker -t leetkob-checker:pa3 .
node scripts/verify-eval.mjs
node scripts/run-offline-eval.mjs
node scripts/verify-extra.mjs
node scripts/run-offline-eval.mjs eval/extra eval/runs/local-extra-offline
```

The script creates the output folder and refuses to overwrite saved predictions.
Its default output is ignored `eval/runs/local-offline/`; choose a new folder
argument when repeating a run.
It streams dataset files through stdin and captures predictions through stdout.
The container uses no network, host mount, or credentials and has a ten-minute timeout.
The image needs no Docker socket, database, or credentials for offline checks.
Input solutions are analyzed by libclang; they are never executed by this command.

For any other dataset, pass its folder to the same script, or use the native
command with the host Docker runner:

```text
go -C backend run ./cmd/checker --data ../eval/extra --predictions ../eval/runs/extra-predictions.jsonl
```

## Real-model evaluation

- Keep `MODEL_OFF=1` until final provider-cap verification is confirmed.
- After verification, use `MODEL_OFF=0`, `MODEL_NAME`, and `MODEL_API_KEY` from
  ignored local configuration. Never put the key on the command line.
- Add `--call-log <new-file.jsonl>` to the same checker command.
- Existing prediction and audit files are refused to preserve previous evidence.
- The command has a default US$1 budget per problem for each run. This is separate
  from the production worker's persistent per-lab database ledger. A new evaluation
  run has a new budget; its previous audit file still records earlier spending.
- The provider call log contains reserved and recorded events. Unknown usage keeps
  its full reservation; failures become `unsure` without automatic retries.
- Exit code 1 from the evaluator means valid predictions missed score thresholds.
  PA3 reports those scores honestly. Exit code 2 means invalid input and must be fixed.

## Saved model-off results

| Dataset | Wrongful deductions | Recall | Unsure rate | Evaluator exit |
|---|---:|---:|---:|---:|
| `runs/offline/` | 0 | 40.5% (15/37) | 46.1% (47/102) | 1 |
| `runs/extra-offline/` | 0 | 40.0% (2/5) | 50.0% (5/10) | 1 |

- All 55 official ban-list pairs were decided. All 47 semantic pairs stayed unsure.
- These are actual offline results, not real-model evaluation evidence.
- The offline verification script allows the expected semantic threshold failure,
  but fails on invalid input, missing predictions, wrongful deductions, undecided
  bans, missed ban violations, or invalid ban evidence. It preserves the evaluator's
  original `FAIL` report. It does not claim the full evaluation passed.
- CI runs these offline checks in the existing sandbox job. Remote CI is not
  claimed until the branch/tag actually runs there.

## Saved real-model results

- `runs/real-dev-numbered/` (2026-11-24, after code is sent as numbered lines):
  102 predictions; 0 wrongful deductions, 94.6% recall (35/37), 2.0% unsure;
  unchanged evaluator exit 0 (PASS). 39 calls, all answered, US$0.021 estimated.
- `runs/real-dev/` (2026-11-22, before that change): 102 predictions; 0 wrongful
  deductions, 56.8% recall, 30.4% unsure; unchanged evaluator exit 1.
- `runs/real-smoke/`: one successful real call on an Accepted extra solution.
- `runs/real-worker/`: database proof of Accepted judging and a real model pass.
- The dev run made 34 sequential requests; five returned HTTP 503 and were not
  retried. Budget blocks and failed evidence validation stayed unsure.
- Every successful response includes its provider ID and token usage. The 503
  responses did not supply IDs or usage; none were fabricated.
- `node scripts/summarize-model-run.mjs eval/runs/real-dev` reproduces recorded
  usage totals and the observed request-rate check without making model calls.
- See [Member 2 handoff](../docs/PA3/MEMBER2-HANDOFF.md) for measured cost arithmetic.
