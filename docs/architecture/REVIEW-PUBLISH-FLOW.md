# TA review and Publish flow

**Updated:** November 27, 2026. PA4 server sampling and publication are implemented.
The PA3 UI still needs Member 3 integration. The authoritative backend behavior is
the [publication contract](../contract/publication.md); see the
[Member 1 handoff and verification](../PA4/MEMBER1-HANDOFF.md).

## Implemented review path

| Component | Source |
| --- | --- |
| Class authorization, queue, earlier submissions and decisions | [review.go](../../backend/internal/api/review.go) |
| Registered routes | [server.go](../../backend/internal/api/server.go) |
| Current decision and history schema | [00006_reviews.sql](../../backend/migrations/00006_reviews.sql) |
| Finding reason schema | [00007_verdict_reasons.sql](../../backend/migrations/00007_verdict_reasons.sql) |
| TA review screen | [review.tsx](../../frontend/src/pages/review.tsx) |
| Existing verification and FAQ references | [TA review evidence](../PA3/TA-REVIEW-QUEUE.md) |

Direct Go routes below omit the browser's `/api` prefix:

| Route | Behavior |
| --- | --- |
| `GET /labs/{labId}/review-queue` | Assigned TA sees every constraint finding for each student's latest counted Accepted submission, with source, verdict, reason, exact evidence, configured deduction and any TA decision |
| `GET /labs/{labId}/earlier-submissions` | Assigned TA can inspect replaced submissions and their findings read-only |
| `POST /labs/{labId}/submissions/{submissionId}/constraints/{constraintId}/review` | Assigned TA records pass or violation for the latest counted Accepted submission, with valid Origin and strict request fields |

The membership check joins the lab's class with `class_tas`. Non-TA roles receive 403; an unknown lab or a TA from another class receives 404. An admin role does not bypass the TA-only review check.

Latest is selected by `(created_at, id)` within each student/lab. A newer pending or failed submission still replaces the older attempt for counting; the queue only includes it after it is judged Accepted and has findings. Earlier attempts stay readable through the separate endpoint, and the decision route refuses them.

The queue includes all pass findings as well as violations/unsure. Required rows
are explicitly flagged: every violation/unsure and the persisted random 10% pass
sample, rounded up. Other pass findings remain optionally reviewable.

## Saving a TA decision

```mermaid
sequenceDiagram
    actor TA as Assigned TA
    participant API as Review API
    participant DB as PostgreSQL
    TA->>API: Decision and evidence
    API->>DB: Check class membership
    API->>API: Validate Origin, JSON and decision
    API->>DB: Begin transaction and lock current lab revision
    API->>API: Require exact original-source evidence for violation
    API->>DB: Upsert current review with configured deduction
    API->>DB: Insert decision history with TA identity
    API->>DB: Commit
    API-->>TA: Decision and reviewer identity
```

A pass requires empty evidence and stores deduction 0. A violation requires nonempty evidence whose line numbers and text match the submitted source exactly; the deduction comes from `lab_constraints.deduction`. The caller cannot submit a deduction field. Current review and history are saved in one transaction.

`constraint_verdicts` remains the worker's finding. `constraint_reviews` is the TA's current decision; `constraint_review_history` records each API review action and reviewer. The API appends history, but the database schema does not enforce immutability by itself. The queue returns the current decision, not a full history feed.

These writes do not publish a mark. Student submission endpoints do not read review findings, reasons, evidence or deductions. Unsure causes are set locally by checker/worker code and displayed to the assigned TA, including exhausted budget and timeout; they are not free-form model explanations.

## PA4 backend sampling and publication

Migration [00008_publication.sql](../../backend/migrations/00008_publication.sql)
adds persisted samples, review revisions, publications and student result snapshots.
[publication.go](../../backend/internal/api/publication.go) implements server random
sampling, progress, atomic Publish and ownership-protected published results.
Population triggers serialize writes with Publish and invalidate prior reviews and
student visibility. Decision edits require republishing. The
[contract](../contract/publication.md) defines exact routes, errors and transitions;
[publication_test.go](../../backend/internal/api/publication_test.go) covers them
against PostgreSQL, including concurrent mutation and restart persistence.

The following checklist is retained as the original design/acceptance criteria;
backend requirements are now implemented. UI, deployment and live-model evidence
remain separate Member 2/3 work.

### Original design criteria

See [Member 1's completed backend checklist](../PA4/PA4-PLAN.md#member-1--server-sampling-required-reviews-and-publish) and the agreed implementation contract. Member 3 integrates the UI.

1. Define a review generation for the lab's current counted submissions, findings and constraint configuration. Incomplete judging/checking must not silently make a missing result eligible for publication.
2. Require review of every violation and every unsure finding in that generation.
3. From the pass findings, draw `ceil(0.10 × pass count)` randomly on the server, with no duplicates. Zero passes yields zero sampled findings. Persist both the selection and its generation; reading the queue or restarting a service must not redraw it.
4. Preserve optional TA access to other findings and earlier submissions while distinguishing which findings block Publish.
5. Treat new counted submissions, changed constraints or rejudging as changes to the review population. A stale review/sample cannot authorize a different generation. Define the refresh and publication policy in the contract and test it.

Implemented records represent the revision, persisted sample, publication snapshot
and publishing TA/time. The current finding/review distinction and history are
reused; model findings have no mark-writing path.

## Publish state transitions

**Status: implemented on the server; UI integration pending.**

```mermaid
stateDiagram-v2
    [*] --> Preparing: Counted submissions are being judged
    Preparing --> ReviewPending: Findings complete and sample persisted
    ReviewPending --> Ready: Every required finding reviewed
    Ready --> Published: Assigned TA presses Publish and API validates atomically
    Ready --> ReviewPending: Required review invalidated
    ReviewPending --> Preparing: Population changes
    Ready --> Preparing: Population changes
    Published --> Preparing: New generation under agreed revision policy
```

The states describe the logical workflow, not database enums. Population changes
invalidate published visibility until the new revision is reviewed and published;
the prior snapshot remains stored. Partial revisions are never exposed.

Publish checks the authenticated role and database class membership, then verifies
the same revision and every required review inside a transaction coordinated with
population changes. A disabled button is only a UI aid. Pending violations, unsure
findings or sampled passes make direct Publish requests fail; a student and another
class's TA are also refused.

Only a successful Publish makes the reviewed deduction visible to its owning
student. Model findings, TA decisions and published results remain separate. The
worker/checker has no mark-writing path. Repeating Publish for the current revision
upserts one snapshot per counted submission.

## Acceptance evidence still required for PA4

- A fixed-seed test exercises the production sampling path, including size, unique membership, non-first-row selection, edge counts and persistence across reads/restarts.
- Direct API tests deny Publish with each required-review type pending, deny student/other-class TA callers and permit the assigned TA after all required reviews.
- Tests show that newer submissions or changed findings/constraints invalidate stale authorization, including a concurrent update at publication time.
- Integration/browser evidence shows no student deduction before Publish and the correct TA-configured deduction afterward, including malicious model deduction fields/instructions.
- The deployed injection tests remain violation or unsure, never pass, for violating solutions.
- Root `PA4.md` explains **Where we stop it** using the team's own evaluation/deduction numbers: who pays, how mistakes are detected, why the gate is at Publish, and the cost of unsampled passes that hide missed violations.

See the [PA4 plan](../PA4/PA4-PLAN.md) for owners, dated handoffs and remaining release evidence. Full appeals behavior belongs to PA5.
