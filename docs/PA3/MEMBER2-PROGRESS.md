# Member 2 progress

## Completed foundation

| Part | Current result |
|---|---|
| Member 1 checks | User confirmed these were already run; not rerun here |
| Model boundary | Replaceable Go interface; no provider or network dependency |
| Offline default | Disabled checker never calls the model |
| Rule routing | Only pending semantic rules enter a model request |
| Response validation | Rejects unknown fields, missing results, invalid verdicts, and invalid evidence |
| Timeout | Default 20 seconds; provider must honor context cancellation |
| Source comments | Blanked while keeping line positions; literals preserved |
| TA evidence | Validated against cleaned code, then mapped to original source lines |
| Gemini adapter | Configured model `gemini-3.5-flash-lite`; fake-transport tests |
| Worker integration | Accepted code reaches semantic checking before fenced completion |
| Identity handling | Comments removed; known account UUID, email, and email local part redacted |
| Durable budget | Lab-row lock and `model_calls` ledger prevent concurrent overspending |
| Usage log | Sanitized prompt, provider response ID, model, tokens, time, estimated cost |
| Environment loading | API and worker share automatic local `.env` loading |

## Next steps

1. Import the official dataset unchanged from the pinned upstream revision below.
2. Extend syntax scanning for loop statements and container/type bans, then add
   the dataset command, offline Docker image, and separate labelled cases.
3. Confirm the final key's project spending cap before enabling real calls.
4. Run the real model and unchanged evaluator; save predictions, metrics, and logs.

## Current limits

- No real model calls or spending occurred during implementation.
- No account metadata, private tests, or deductions enter provider requests.
  Unknown personal names embedded in code cannot be reliably inferred from the
  current account schema, which stores only email and UUID. Known identifiers
  are redacted; arbitrary personal text remains a limitation.
- Line-spliced or ambiguous C++ source becomes `unsure` without a model call.
- Official source found: [fithcmus/awad-leetkob-eval](https://github.com/fithcmus/awad-leetkob-eval).
  Inspected revision: `552c811db2d35a23d641bc4fdffc1ac120be1a28`.
  It contains six problems, 47 solutions, and 102 labelled pairs.
  Files were inspected in a temporary checkout; not yet imported into this repository.
- The current scanner reports calls and declaration references. The official
  constraints also ban loop statements and container types; those need broader
  syntax coverage before claiming complete offline ban-list support.
- The full Member 2 checklist remains open until integration and evidence exist.

## Official dataset format

| Path | Contents |
|---|---|
| `problems/<id>/statement.md` | Problem statement |
| `problems/<id>/constraints.json` | Constraint IDs, text, and deductions |
| `problems/<id>/tests/` | Judge input/output files |
| `solutions/<id>/*.cpp` | Accepted C++ solutions |
| `labels.jsonl` | Expected labels; never used as checker input |
| `eval.mjs` | Unchanged Node 22+ evaluator |

Predictions use one JSONL record per solution/constraint pair, with `solution`,
`constraint`, `verdict`, and `evidence` fields. Evaluate with:

```text
node <official-folder>/eval.mjs --data <official-folder> --predictions <predictions.jsonl> --json <summary.json>
```

Exit 0 means thresholds met; exit 1 means valid input but thresholds missed;
exit 2 means invalid input. PA3 requires honest reporting, not passing the PA4
metric thresholds.

## Local configuration for the next step

Edit the ignored root `.env` file locally. Set `MODEL_NAME` to the exact
provider model ID and `MODEL_API_KEY` to its key. Keep `MODEL_OFF=1` while
preparing the adapter and tests. Never paste the key into chat or commit it.

Verify in the provider console that the final project/key uses the intended
spending cap. Existing screenshot evidence alone does not confirm this.

The worker now loads `.env` automatically from the repository root or `backend/`.
Shell variables take precedence. Restart the worker after changing local settings.

## Cost controls

| Setting or state | Meaning |
|---|---|
| `MODEL_OFF=1` | No model call or budget reservation |
| `MODEL_OFF=0` | Explicitly enables the configured model |
| `MODEL_LAB_BUDGET_MICRO_USD` | Optional integer cap; default `1000000` means US$1 per lab |
| Before each call | Reserve 478,413 microdollars using full documented model limits |
| Known usage | Replace reservation with rounded-up standard-price estimated cost |
| Timeout, crash, unknown usage | Keep the full reservation |
| Insufficient remaining budget | Return `unsure`; do not contact the provider |

- Apply migration `00005_model_calls.sql` before enabling model calls.
- Reservations and charges are cumulative across submissions and worker restarts.
- The reservation is conservative: a call can be blocked even when its likely
  actual cost would fit. The model and its reviewed prices are fixed in code.
- Provider pricing changes require updating the rates and reservation together.
- Call logs are protected database data; no student endpoint exposes them.
- The API container no longer receives model credentials; only the host worker
  needs them. The sandbox never receives worker environment variables.
