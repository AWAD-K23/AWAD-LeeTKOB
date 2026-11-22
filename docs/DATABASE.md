# Current database structure

- **Database:** PostgreSQL 17, accessed from Go through `pgx`.
- **Schema source:** `backend/migrations/00001` through `00005`.
- **Local port:** `5433`. Migrations preserve the existing local data.
- **Storage:** Named Docker volume `postgres_data` in Compose; survives normal container restarts. `docker compose down -v` deletes it.
- **Backend guide:** See [BACKEND.md](BACKEND.md) for commands, files, and routes.

## Relationships

```mermaid
erDiagram
    users ||--o{ class_tas : teaches
    users ||--o{ class_students : joins
    classes ||--o{ class_tas : assigns
    classes ||--o{ class_students : contains
    classes o|--o{ labs : owns
    labs ||--o{ lab_tests : tests
    labs ||--o{ lab_constraints : requires
    users ||--o{ enrollments : has
    labs ||--o{ enrollments : permits
    users ||--o{ submissions : submits
    labs ||--o{ submissions : receives
    submissions ||--o| jobs : queues
    submissions ||--o| judge_results : receives
    submissions ||--o{ constraint_verdicts : receives
    submissions ||--o{ model_calls : checks
    labs ||--o{ model_calls : budgets
    lab_constraints ||--o{ constraint_verdicts : checks
    users ||--o{ idempotency_keys : retries
    labs ||--o{ idempotency_keys : scopes
```

- **Primary key (PK):** Identifies a row, like an object's unique `id` in JavaScript.
- **Foreign key (FK):** Requires a linked row to exist in another table.
- **Unique key:** Prevents duplicate values or duplicate combinations.
- `sessions` keeps its user ID inside session data, so it has no direct user FK.

## Tables and important fields

| Table | Fields | Purpose and rules |
|---|---|---|
| `users` | `id`, `email`, `password_hash`, `role` | UUID PK; unique email; role is `admin`, `ta`, or `student`. No plain-text passwords. |
| `sessions` | `token`, `data`, `expiry` | Token PK; stores server-side SCS login sessions. Index on expiry supports cleanup. |
| `classes` | `id`, `name` | UUID PK; groups TAs, students, and labs. |
| `class_tas` | `class_id`, `ta_id` | Combined PK; links classes to users. The API requires the assigned user to be a TA. |
| `class_students` | `class_id`, `student_id` | Combined PK; links classes to users. The API requires the assigned user to be a student. |
| `labs` | `id`, `title`, `class_id`, `statement` | UUID PK; class FK; statement stored as plain text containing Markdown, rendered by the frontend. New labs require a TA's class. Old PA2 labs may have a null class and an empty statement. |
| `lab_tests` | `id`, `lab_id`, `input`, `expected_output` | UUID PK; test inputs and expected outputs. Removed with their lab through `ON DELETE CASCADE`. |
| `lab_constraints` | `id`, `lab_id`, `description`, `deduction` | UUID PK; TA-written constraint and nonnegative numeric deduction. Removed with their lab through `ON DELETE CASCADE`. |
| `enrollments` | `lab_id`, `student_id` | Combined PK; controls whether a student may submit to a lab. Class students are enrolled when the TA creates a lab. |
| `submissions` | `id`, `lab_id`, `student_id`, `source`, `created_at` | UUID PK; source is `bytea` so submitted bytes, including NUL, are preserved. |
| `jobs` | `id`, `submission_id`, `status`, `created_at`, `lease_token`, `leased_until`, `attempts` | UUID PK; submission FK is unique. Tracks durable judging work and recovery after a worker crash. Index on status and creation time supports job claims. |
| `judge_results` | `submission_id`, `verdict`, `detail`, `judged_at` | Submission FK is also the PK: at most one result per submission. Stores Accepted, Wrong Answer, resource errors, or Judge Error. |
| `constraint_verdicts` | `submission_id`, `constraint_id`, `verdict`, `evidence`, `source` | Combined PK: at most one verdict per submission/constraint. Verdict is `pass`, `violation`, or `unsure`. Evidence is JSONB. No deduction or published mark is stored here. |
| `model_calls` | `id`, `lab_id`, `submission_id`, `charged_micro_usd`, `status`, `call_log`, `created_at` | Persistent model budget reservations and protected JSONB audit records. `reserved` retains the maximum estimated cost; `recorded` uses known usage or retains the reservation if usage is unknown. Added by migration 00005. |
| `idempotency_keys` | `student_id`, `lab_id`, `key`, `fingerprint`, `response`, `created_at` | Combined PK on student/lab/key. Stores a SHA-256 request fingerprint and original response so retries do not create duplicate submissions. |
| `goose_db_version` | `id`, `version_id`, `is_applied`, `tstamp` | Migration tool's version history, created by Goose when migrations run; not an application table defined in migrations 00001–00005. |

- User foreign keys verify that an account exists; API checks enforce TA/student roles.
- Cascades on test/rule lab links do not mean every lab can be deleted: other foreign
  keys, including submissions and constraint findings, can block deletion.

## Job and checker states

| Field | Value | Meaning |
|---|---|---|
| `jobs.status` | `pending` | Waiting for a worker. The create-submission response calls this `queued`. |
| | `processing` | Claimed by a worker with a renewable 60-second lease. |
| | `completed` | Judge result has been stored. This can include a non-Accepted verdict. |
| | `failed` | Reserved by the schema; current bounded infrastructure retries finish with a stored Judge Error result. |
| `constraint_verdicts.source` | `hard_rule` | Syntax-aware ban-list decision; a failed analysis remains unsure. |
| | `semantic_pending` | Semantic work left unsure while model calls are off. |
| | `model` | Validated model finding, or unsure after a failed, invalid, timed-out, or budget-blocked attempt. |
| `model_calls.status` | `reserved` | Maximum estimated cost reserved before provider access; may remain after a crash |
| | `recorded` | Audit saved; known usage replaces the reservation, while unknown usage retains it |

- The worker renews its lease every 10 seconds.
- When idle or after a processing error, the worker waits one second before looking
  for another job. After success, it checks again immediately.
- A new worker may claim a job after its lease expires.
- Before storing results, the worker checks the lease token inside a transaction. An old worker cannot overwrite the new worker's result.
- Database primary keys also prevent duplicate result and verdict rows.

## Example flow

1. Admin creates **Sorting A**, assigning TA Linh and student An.
2. Linh creates **Merge sort** with tests and two constraints. The server checks Linh's class membership, then saves the lab, tests, constraints, and student enrollments in one transaction.
3. An submits C++ code with a retry key. The API saves the submission, pending job, and retry response in one transaction.
4. The worker claims the job, increasing `attempts` and setting a lease token and expiry.
5. The sandbox compiles and runs the code. If every test output matches, the worker proceeds to constraint checking; the judge result will be **Accepted** regardless of separate constraint findings.
6. Hard rules check ban lists. Semantic constraints stay `unsure` with `semantic_pending` when model calls are off; enabled calls use the lab budget ledger and validated Gemini responses.
7. The worker stores the result and constraint verdicts, then marks the job completed in one transaction.
8. An's status endpoint returns job status, judge result, saved source, and lab ID. Constraint verdicts and TA deductions stay separate for the TA review work.

## Model budget and call records

| Action | Database effect |
| --- | --- |
| Before a real call | Lock the lab row, sum previous charges, and insert a `reserved` model-call row if budget remains |
| Successful response with usage | Save audit JSON and replace the reservation with the rounded token-based estimate |
| Error with unknown usage | Save available audit data and retain the full reservation |
| Worker crash or audit-write failure | Reservation remains charged; restarting cannot reset the lab budget |
| Budget exhausted | No provider request; semantic finding becomes unsure |

- `charged_micro_usd` stores the budget charge in microdollars. Reservations
  use a conservative estimate, not a provider bill.
- All submissions in the lab share the budget. It does not reset per student,
  job, or worker restart.
- A submission can have multiple call records after recovery or later processing;
  the unique constraint verdict key still allows only one saved finding per rule.
- Audit records contain sanitized prompts and raw provider responses. A recorded
  response may still fail verdict/evidence validation and produce unsure.
- Provider IDs and usage can be absent on errors; `recorded` does not mean the
  provider request succeeded.
- The standalone dataset checker writes JSONL files instead of these database rows.

## Current data visibility

| API data | What is returned |
| --- | --- |
| Accessible lab list | Statement, class ID, and public requirement descriptions; no private tests or deductions |
| Student submission history | Latest 200 submissions owned by the signed-in student |
| Submission details | Owner or class TA can read status, judge result, saved source, and lab ID |
| Class details | Admin/assigned TA can read member accounts; students receive no roster |
| Admin account search | Up to 200 matching account IDs/emails/roles; no password hashes |
| Model call audits | No public/student endpoint; protected backend records |

- The frontend redesign added these reads without new schema migrations.
- An Accepted judge result can coexist with a `violation` in `constraint_verdicts`.
- TA-configured deductions are not automatically applied or published.

- Example constraint evidence:

```json
{
  "verdict": "violation",
  "source": "hard_rule",
  "evidence": [
    { "line": 12, "text": "    std::sort(a, a + n);" }
  ]
}
```

## Example model result

| Field | Meaning |
| --- | --- |
| Judge result | Accepted |
| Sort ban | Pass, source `hard_rule` |
| Comparison requirement | Pass, source `model`, with the matching source line |
| Call status | Recorded, with a provider request ID and known token usage |

- Execution results, constraint findings, and provider audits are stored separately.
- See [Member 2 handoff](PA3/MEMBER2-HANDOFF.md) for the saved local example.

## Remaining work

- Real-model checking and evaluation evidence are delivered; semantic accuracy
  still needs improvement before PA4. See [Member 2 handoff](PA3/MEMBER2-HANDOFF.md).
- There are no TA-review decisions, published marks, or appeals tables yet.
- The API enforces roles and ownership; user foreign keys alone do not enforce a role.
- See [Member 1 handoff](PA3/MEMBER1-HANDOFF.md) for API interfaces, verification commands, and deployment steps.
