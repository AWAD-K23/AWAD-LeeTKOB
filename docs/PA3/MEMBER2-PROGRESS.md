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

## Next steps

1. Configure the real provider, model name, and API key locally.
2. Add the provider adapter, prompt, request/usage logging, and cost limits.
3. Connect the checker to the worker before its existing result transaction.
4. Add the dataset command, offline Docker image, and separate labelled cases.
5. Run the unchanged official evaluator and save real-run evidence.

## Current limits

- This package is not connected to the worker yet; application behavior is unchanged.
- No real model calls or spending occurred during this step.
- Requests have no account fields, private tests, or deductions. Full identity
  handling for the real provider remains part of the next step; arbitrary
  identifying text inside code literals is not removed by the comment cleaner.
- Line-spliced or ambiguous C++ source becomes `unsure` without a model call.
- No official dataset or evaluator was found in the repository file listing.
  Their original files are needed before building the dataset adapter.
- The full Member 2 checklist remains open until integration and evidence exist.

## Local configuration for the next step

Edit the ignored root `.env` file locally. Set `MODEL_NAME` to the exact
provider model ID and `MODEL_API_KEY` to its key. Keep `MODEL_OFF=1` while
preparing the adapter and tests. Never paste the key into chat or commit it.

Verify in the provider console that the final project/key uses the intended
spending cap. Existing screenshot evidence alone does not confirm this.

The worker currently needs environment settings in its own process; automatic
loading for this configuration will be addressed during worker integration.
