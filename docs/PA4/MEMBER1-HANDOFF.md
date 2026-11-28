# PA4 Member 1 backend handoff

Implemented November 27, 2026 (Asia/Saigon). Ready for Member 3 integration;
teammate acceptance and public deployment have not been performed.

Local commit breakdown: `6b40ca2` defines the contract; `f287fec` contains the
backend implementation, migration and tests together. Handoff and tracker updates
are grouped separately. These commits organize completed workspace work and do
not claim a reconstructed test-first history.

## Contract and migration

Read [the API contract](../contract/publication.md) before wiring the UI.
Apply migration `00008_publication.sql` using the existing migration command:

```text
go -C backend run ./cmd/migrate
```

Supply `DATABASE_URL` through the local environment as documented in README.
The Goose migration was applied successfully to a fresh temporary PostgreSQL 17
database. It adds lab revisions, persisted populations/samples, review/history
revision columns, publications and result snapshots. PA3 reviews stay in history
but revision zero cannot satisfy the new gate. Population changes conservatively
invalidate all reviews and published visibility for the lab. No production
database has been migrated during this implementation.

## Member 3 integration

1. Load the existing queue plus `GET /labs/{labId}/publication`. The queue remains
   an array and includes all findings. Highlight `reviewRequired`, `reviewReason`
   and `sampledPass`; use status counts for progress, including an empty queue.
2. Save TA decisions using the existing review route and evidence format. Reload
   progress after each decision. Evidence and configured deductions remain
   independent of model findings. A saved decision invalidates published visibility
   until Publish is repeated.
3. Send `POST /labs/{labId}/publish` with the exact status revision:

   ```json
   {"reviewRevision":123}
   ```

   Include the session cookie and configured Origin. An example ready response is:

   ```json
   {"reviewRevision":123,"required":4,"reviewed":4,"pending":0,"incompleteSubmissions":0,"published":true}
   ```

   Display server 409 `pending-reviews` and `stale-review-population` errors and
   reload current state. A disabled button does not replace this server gate.
4. The owning student's result view calls
   `GET /submissions/{submissionId}/published-result`. A 404 means no accessible
   current publication. Success includes the configured TA deduction total,
   publisher, timestamp and revision. Submission status/history keep their old
   responses; completing a review does not expose deductions there.

Preparing (incomplete judging) → review pending → ready → published. A new
submission, rejudge/status transition, finding change, constraint edit or statement
change returns the lab to a new revision. An earlier attempt cannot receive new
decisions. Publish serializes validation and snapshot writes against all population
writers through a shared lab-row lock in database triggers.

The new endpoints use the same class/ownership policy as PA3: students cannot
review or publish; another class's TA receives 404. Students can read only their
own current published result. Administrators have no TA bypass.

## Verification

Passed locally against a temporary PostgreSQL 17 Docker container on loopback
port 55434, isolated from the project database:

- `node scripts/check-go-format.mjs`
- `go -C backend vet ./...`
- `go -C backend test ./... -count=1` with `TEST_DATABASE_URL` set
- `go -C backend test ./internal/api -run 'TestPublicationGate|TestPassSample' -count=1 -v`
- `node scripts/verify-eval.mjs`: all 128 official files unchanged
- Goose `go -C backend run ./cmd/migrate`: migrated through version 8

`TestPassSample` covers pass counts 0, 1, 9, 10, 11, 21 and 100, ceiling size,
unique in-range membership and the fixed-seed non-first-row sample `[5,8,12]`
for 21 passes. `TestPublicationGate` exercises production sampling persistence
across reads and a new server instance whose RNG fails if called. It checks all
violation/unsure requirements plus two of eleven passes, pending/sample/role/class
denials, strict mark-field rejection, successful publication (deduction 6 from
two configured deductions of 3), owning-student visibility, repeated Publish,
decision edits, rejudging, finding/constraint changes and replacement attempts.
The concurrency test observes Publish waiting on PostgreSQL's population lock
while a constraint mutation is uncommitted, then checks the stale-revision denial.
History remains intact and incomplete Accepted results cannot be published.

Frontend `npm ci --prefix frontend` completed. Type-checking and the production
build passed separately. `npm run check --prefix frontend`
passed lint with two existing Fast Refresh warnings, then stopped on existing
Prettier formatting in untouched `frontend/src/style.css`. The extra-case verifier
failed on the CRLF working-tree test outputs. It passed all five cases against an
unchanged temporary Git export made with `-c core.autocrlf=false` (the normal
Windows Git export also converts these outputs). Neither repository fixtures nor
the verifier were changed. Official hashes remain unchanged.

No real model requests, public deployment, browser integration, teammate sign-off
or push are claimed. Members 2 and 3 still own live evaluation, guardrail/
injection evidence, the UI and release work. The backend tests use synthetic data
and fake models only. Check detailed remaining evidence in [the PA4 plan](PA4-PLAN.md).
