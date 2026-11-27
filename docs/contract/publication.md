# PA4 review and publication contract

Member 1 backend contract, November 27, 2026. Member 3 UI acceptance remains separate.

The counted submission is the latest `(created_at, id)` for each student in a lab,
including pending and failed attempts. Earlier attempts remain read-only. Accepted
counted attempts contribute one finding per current lab constraint. Every violation
and unsure requires a TA decision. Of all pass findings in that lab population,
the server uniformly samples `ceil(passCount / 10)` without replacement.

The lab has a monotonically increasing `reviewRevision`. Database triggers serialize
population changes on its lab row: submissions, judge results, constraint findings,
job creation/status transitions, and constraint edits invalidate the whole lab's
sample, decisions and publication. Lab statement changes also invalidate them.
This conservative policy requires fresh reviews even when an unrelated student's
attempt changes. Historical decisions and publications remain recorded. Existing
PA3 decisions have revision zero and must be reviewed again after migration.
Rejudging must transition the job status. Constraint edits require fresh TA
decisions against the changed requirements (the API has no constraint edit route).

`GET /labs/{labId}/review-queue` retains its array response. Each row adds
`reviewRevision`, `reviewRequired`, `sampledPass`, and `reviewReason` (`violation`,
`unsure`, `sampled-pass`, or empty). The sample is persisted, including empty
populations, and is never redrawn on reload within the same revision.

`GET /labs/{labId}/publication` returns `reviewRevision`, `required`, `reviewed`,
`pending`, `incompleteSubmissions`, and `published`. Both reads require the assigned
TA. Incomplete means a counted attempt lacks a completed job/judge result or an
Accepted attempt lacks a finding for every current constraint.

`POST /labs/{labId}/submissions/{submissionId}/constraints/{constraintId}/review`
keeps the PA3 body: `{"decision":"pass"}` or
`{"decision":"violation","evidence":[{"line":1,"text":"exact source line"}]}`.
Only current Accepted attempts can be reviewed. The server captures the current
revision, configured deduction and named reviewer, and appends history. A changed
decision invalidates an existing publication; Publish must be repeated.

`POST /labs/{labId}/publish` takes `{"reviewRevision":123}` (from the queue or
publication status). It requires the assigned TA and configured Origin. A stale
revision returns 409 `stale-review-population`; pending required reviews or incomplete
judging returns 409 `pending-reviews`. No client-supplied mark is accepted. A single
transaction locks the lab, checks completeness and reviews, then snapshots every
counted result. Deductions are the sum of current configured constraint deductions
for valid TA violation decisions, including optional reviewed passes; unreviewed
unsampled passes contribute zero. Repeating Publish for the same revision is safe.

`GET /submissions/{submissionId}/published-result` requires the owning student or
assigned TA. Before publication, for an earlier attempt, or after invalidation it
returns 404. Otherwise it returns `submissionId`, `reviewRevision`, `deduction`,
`publishedAt`, and `publishedBy`. Existing status/history responses keep their PA3
shape and expose no deductions. Model findings never write published results.

All population writers acquire the same lab row through triggers, including direct
SQL and worker writes. Publish and review use that lock before reading the
population. A concurrent writer either commits first and is included, or follows
Publish and immediately invalidates it. Student result visibility uses one SQL
snapshot joining the publication revision to the current lab revision.
