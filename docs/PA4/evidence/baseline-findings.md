# Fresh PA4 baseline findings

The three original-budget dev runs use frozen commit
`77ef5b0aa57c5e019ce0d6d34712151ac3b648c9`, the same model/configuration and
official dataset hash. Run folders preserve the original predictions and audits.
This record is diagnostic, not a passing-performance claim.

All three runs found one wrongful violation each. They caught 33/37 (89.2%),
34/37 (91.9%) and 35/37 (94.6%) violations. Unsure was 3/102, 3/102 and 0/102.
All 39 responses in each run had provider IDs and known usage. No labels,
official source, thresholds or existing saved runs were changed.

Generic failure mechanisms observed:

- A declaration was treated as proof of a prohibited computation without following
  how the input was transformed. Mere presence of a storage variable is not proof
  of extra working storage.
- An unused compliant implementation distracted the model from the computation
  actually called. Names and dead code are insufficient algorithm evidence.
- A rule about explicit allocation/release was confused with other storage
  lifetimes. Each constraint's scope and exceptions must be assessed independently.
- One invalid finding currently makes all findings in the same response unsure,
  including a separate constraint with valid evidence. A possible improvement is
  to preserve independently valid findings while keeping every invalid finding
  unsure and continuing to reject malformed/duplicate/missing response structure.

Any checker fix must remain generic and must not include official problem IDs,
solution IDs or quoted dev code in production prompts/rules. Exact evidence and
schema checks remain mandatory. Proposed improvements require their own tests
and fresh evidence; they are not claimed to solve the observed failures yet.

The human review/Publish gate remains separate: these are model findings, not
student deductions. Deployed UI acceptance is still pending.
