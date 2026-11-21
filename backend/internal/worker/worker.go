// Package worker consumes durable jobs independently of HTTP request handlers.
package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/hardrules"
	"leetkob/backend/internal/sandbox"
)

var ErrLeaseLost = errors.New("job lease no longer belongs to this worker")

type Job struct{ ID, SubmissionID, Token string }
type Worker struct {
	DB     *pgxpool.Pool
	Runner hardrules.Scanner
}

func (w Worker) Claim(ctx context.Context) (Job, error) {
	token := uuid.NewString()
	var job Job
	// SKIP LOCKED lets multiple workers claim different rows. Expired leases
	// recover work after a crash, without deleting the queued submission.
	err := w.DB.QueryRow(ctx, `UPDATE jobs SET status='processing', lease_token=$1,
 leased_until=now()+interval '60 seconds', attempts=attempts+1 WHERE id=(
 SELECT id FROM jobs WHERE status='pending' OR (status='processing' AND leased_until<now())
 ORDER BY created_at,id FOR UPDATE SKIP LOCKED LIMIT 1)
 RETURNING id::text,submission_id::text,lease_token::text`, token).Scan(&job.ID, &job.SubmissionID, &job.Token)
	return job, err
}

func (w Worker) Complete(ctx context.Context, job Job, verdict, detail string, checks []hardrules.Verdict) error {
	tx, err := w.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var token string
	err = tx.QueryRow(ctx, "SELECT lease_token::text FROM jobs WHERE id=$1 AND status='processing' FOR UPDATE", job.ID).Scan(&token)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && token != job.Token) {
		return ErrLeaseLost
	}
	if err != nil {
		return err
	}
	// The lease token fences out an old worker after another worker recovers
	// its job. Unique keys ensure one stored result and one verdict per pair.
	_, err = tx.Exec(ctx, `INSERT INTO judge_results(submission_id,verdict,detail) VALUES($1,$2,$3)
 ON CONFLICT(submission_id) DO UPDATE SET verdict=excluded.verdict,detail=excluded.detail,judged_at=now()`, job.SubmissionID, verdict, detail)
	if err != nil {
		return err
	}
	for _, check := range checks {
		evidence, err := json.Marshal(check.Evidence)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO constraint_verdicts(submission_id,constraint_id,verdict,evidence,source)
 VALUES($1,$2,$3,$4,$5) ON CONFLICT(submission_id,constraint_id) DO UPDATE
 SET verdict=excluded.verdict,evidence=excluded.evidence,source=excluded.source`, job.SubmissionID, check.ConstraintID, check.Verdict, evidence, check.Source); err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, "UPDATE jobs SET status='completed',leased_until=NULL WHERE id=$1", job.ID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (w Worker) ProcessOne(ctx context.Context) (bool, error) {
	job, err := w.Claim(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	workCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Refresh the lease while compilation and multiple test cases run. A failed
	// heartbeat stops this worker; completion still checks the fencing token.
	done := make(chan struct{})
	defer close(done)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-workCtx.Done():
				return
			case <-ticker.C:
				tag, e := w.DB.Exec(workCtx, "UPDATE jobs SET leased_until=now()+interval '60 seconds' WHERE id=$1 AND lease_token=$2 AND status='processing'", job.ID, job.Token)
				if e != nil || tag.RowsAffected() != 1 {
					cancel()
					return
				}
			}
		}
	}()
	err = w.judge(workCtx, job)
	if err != nil {
		// Infrastructure errors retry at most three times; never label a student
		// wrong just because Docker or the database was temporarily unavailable.
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		var attempts int
		if e := w.DB.QueryRow(cleanup, "SELECT attempts FROM jobs WHERE id=$1", job.ID).Scan(&attempts); e == nil && attempts >= 3 {
			_ = w.Complete(cleanup, job, "Judge Error", "Judging could not be completed. Contact the TA.", nil)
		} else {
			_, _ = w.DB.Exec(cleanup, "UPDATE jobs SET status='pending',leased_until=NULL WHERE id=$1 AND lease_token=$2 AND status='processing'", job.ID, job.Token)
		}
	}
	return true, err
}

func (w Worker) judge(ctx context.Context, job Job) error {
	var source []byte
	var lab string
	if err := w.DB.QueryRow(ctx, "SELECT source,lab_id::text FROM submissions WHERE id=$1", job.SubmissionID).Scan(&source, &lab); err != nil {
		return err
	}
	rows, err := w.DB.Query(ctx, "SELECT input,expected_output FROM lab_tests WHERE lab_id=$1 ORDER BY id", lab)
	if err != nil {
		return err
	}
	type testCase struct{ input, expected string }
	tests := []testCase{}
	for rows.Next() {
		var test testCase
		if err = rows.Scan(&test.input, &test.expected); err != nil {
			rows.Close()
			return err
		}
		tests = append(tests, test)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	if len(tests) == 0 {
		return errors.New("lab has no test cases")
	}
	for i, test := range tests {
		result, err := w.Runner.Run(ctx, source, test.input, "run")
		if err != nil {
			return err
		}
		if result.Verdict != "Executed" {
			return w.Complete(ctx, job, result.Verdict, fmt.Sprintf("Test case %d did not complete successfully.", i+1), nil)
		}
		if normalize(result.Output) != normalize(test.expected) {
			return w.Complete(ctx, job, "Wrong Answer", fmt.Sprintf("Output did not match test case %d.", i+1), nil)
		}
	}
	rows, err = w.DB.Query(ctx, "SELECT id::text,description FROM lab_constraints WHERE lab_id=$1 ORDER BY id", lab)
	if err != nil {
		return err
	}
	constraints := []hardrules.Constraint{}
	for rows.Next() {
		var c hardrules.Constraint
		if err = rows.Scan(&c.ID, &c.Text); err != nil {
			rows.Close()
			return err
		}
		constraints = append(constraints, c)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	checks, err := hardrules.Check(ctx, w.Runner, source, constraints)
	if err != nil {
		// A failed syntax check is uncertain, never an automatic pass or deduction.
		checks = []hardrules.Verdict{}
		for _, c := range constraints {
			checks = append(checks, hardrules.Verdict{ConstraintID: c.ID, Verdict: "unsure", Evidence: []hardrules.Evidence{}, Source: "semantic_pending"})
		}
	}
	return w.Complete(ctx, job, "Accepted", "All test cases passed.", checks)
}

func normalize(text string) string { return strings.Join(strings.Fields(text), " ") }

var _ hardrules.Scanner = sandbox.Runner{}
