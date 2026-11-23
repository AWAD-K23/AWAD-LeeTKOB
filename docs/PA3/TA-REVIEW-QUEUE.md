# TA review queue verification

## Behavior

- The queue is restricted to a TA assigned to the class that owns the lab. The server checks membership from the database; an unknown lab and a lab from another TA's class both return 404.
- The queue shows each student's latest counted submission, its judge verdict, each constraint finding, the checker verdict and evidence, the configured deduction, and the TA decision with the applied deduction.
- A model verdict is only a finding. A TA must choose pass (deduction 0) or violation (the configured constraint deduction). For a violation, submitted source lines must be supplied as exact evidence. Each decision records the TA and a history row.
- Student submission endpoints do not read review rows. The focused contract test checks that neither model findings nor TA review details appear in student responses before publication.
- No Publish or marks endpoint exists in the current backend routes. This verifies the current build has no route that can turn a finding into a published mark; it does not claim a future Publish flow is gated.

## Verification evidence

The repository owner ran this focused API contract test from Ubuntu WSL, using a temporary PostgreSQL 17 Docker container and the Go 1.26 image:

```text
=== RUN   TestSubmissionContract/TA_review_is_class_scoped_and_cannot_leak_into_student_marks
--- PASS: TestSubmissionContract/TA_review_is_class_scoped_and_cannot_leak_into_student_marks
PASS
ok      leetkob/backend/internal/api    0.575s
Test exit status: 0
```

The test covers class ownership, latest submission selection, rejection of fabricated evidence, no review row before TA action, a TA-attributed violation decision and deduction, and no review-data leak through student submission endpoints. It is a focused API integration test; it does not cover a deployed end-to-end flow, frontend build, or a Publish action.

## Product-owner FAQ references

- [FAQ-02](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md): which submission counts.
- [FAQ-08](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md): student visibility before publication.
- [FAQ-09](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md): handling unsure findings and evidence.
- [FAQ-10](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md): findings sent to review and TA changes recorded.
- [FAQ-12](https://github.com/tkob-team/LeeTKOB-FAQ/blob/main/FAQ.md): when publication is allowed.
