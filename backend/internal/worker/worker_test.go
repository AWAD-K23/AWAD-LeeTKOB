package worker

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/sandbox"
	"leetkob/backend/internal/semantic"
)

type fakeRunner struct{}

type failedAnalyzer struct{}

type fakeSemanticModel func(context.Context, semantic.Request) ([]byte, error)

func (f fakeSemanticModel) Generate(ctx context.Context, request semantic.Request) ([]byte, error) {
	return f(ctx, request)
}

func (failedAnalyzer) Run(_ context.Context, _ []byte, _, mode string) (sandbox.Result, error) {
	if mode == "analyze" {
		return sandbox.Result{}, errors.New("analysis unavailable")
	}

	return sandbox.Result{Verdict: "Executed", Output: "1 3"}, nil
}

func (fakeRunner) Run(_ context.Context, _ []byte, _, mode string) (sandbox.Result, error) {
	if mode == "analyze" {
		return sandbox.Result{Verdict: "Executed", Output: "std::sort\t2\n"}, nil
	}

	return sandbox.Result{Verdict: "Executed", Output: "1 3\n"}, nil
}

func TestDurableWorker(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")

	if url == "" {
		t.Skip("set TEST_DATABASE_URL")
	}

	ctx := context.Background()
	base, err := pgxpool.New(ctx, url)

	if err != nil {
		t.Fatal(err)
	}

	defer base.Close()
	schema := "worker_" + strings.ReplaceAll(uuid.NewString(), "-", "")

	if _, err = base.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}

	defer base.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
	config, err := pgxpool.ParseConfig(url)

	if err != nil {
		t.Fatal(err)
	}

	config.ConnConfig.RuntimeParams["search_path"] = schema
	db, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	files, err := filepath.Glob("../../migrations/*.sql")

	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		sql, e := os.ReadFile(file)

		if e != nil {
			t.Fatal(e)
		}

		if _, e = db.Exec(ctx, strings.Split(string(sql), "-- +goose Down")[0]); e != nil {
			t.Fatal(file, e)
		}
	}

	student, lab := uuid.NewString(), uuid.NewString()

	if _, err = db.Exec(ctx, "INSERT INTO users(id,email,password_hash,role) VALUES($1,'worker@test.invalid','unused','student')", student); err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec(ctx, "INSERT INTO labs(id,title) VALUES($1,'Sorting')", lab); err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec(ctx, "INSERT INTO lab_tests(id,lab_id,input,expected_output) VALUES($1,$2,'','1 3')", uuid.NewString(), lab); err != nil {
		t.Fatal(err)
	}

	for _, text := range []string{"Do not call `std::sort`", "Use merge sort"} {
		if _, err = db.Exec(ctx, "INSERT INTO lab_constraints(id,lab_id,description,deduction) VALUES($1,$2,$3,3)", uuid.NewString(), lab, text); err != nil {
			t.Fatal(err)
		}
	}

	w := Worker{DB: db, Runner: fakeRunner{}}

	if os.Getenv("SANDBOX_TESTS") == "1" {
		w.Runner = sandbox.Runner{}
	}

	enqueue := func() string {
		t.Helper()
		id := uuid.NewString()
		source := []byte("#include <algorithm>\n#include <iostream>\nint main(){int a[]={3,1};std::sort(a,a+2);std::cout<<a[0]<<\" \"<<a[1];}")

		if _, err := db.Exec(ctx, "INSERT INTO submissions(id,lab_id,student_id,source,created_at) VALUES($1,$2,$3,$4,now())", id, lab, student, source); err != nil {
			t.Fatal(err)
		}

		if _, err := db.Exec(ctx, "INSERT INTO jobs(id,submission_id) VALUES($1,$2)", uuid.NewString(), id); err != nil {
			t.Fatal(err)
		}

		return id
	}

	t.Run("queued work survives absent worker and is stored once", func(t *testing.T) {
		id := enqueue()
		var status string

		if err := db.QueryRow(ctx, "SELECT status FROM jobs WHERE submission_id=$1", id).Scan(&status); err != nil || status != "pending" {
			t.Fatal(status, err)
		}

		// A new worker instance represents restarting a stopped worker.
		restarted := Worker{DB: db, Runner: w.Runner}

		if found, err := restarted.ProcessOne(ctx); err != nil || !found {
			t.Fatal(found, err)
		}

		if found, err := restarted.ProcessOne(ctx); err != nil || found {
			t.Fatal(found, err)
		}

		var results, checks, violations, pending, pendingReasons int
		var verdict string

		if err := db.QueryRow(ctx, "SELECT verdict FROM judge_results WHERE submission_id=$1", id).Scan(&verdict); err != nil || verdict != "Accepted" {
			t.Fatal(verdict, err)
		}

		if err := db.QueryRow(ctx, `SELECT (SELECT count(*) FROM judge_results WHERE submission_id=$1),
   count(*),count(*) FILTER(WHERE verdict='violation'),count(*) FILTER(WHERE source='semantic_pending'),
   count(*) FILTER(WHERE source='semantic_pending' AND reason<>'')
   FROM constraint_verdicts WHERE submission_id=$1`, id).Scan(&results, &checks, &violations, &pending, &pendingReasons); err != nil || results != 1 || checks != 2 || violations != 1 || pending != 1 || pendingReasons != 1 {
			t.Fatalf("results=%d checks=%d violations=%d pending=%d pendingReasons=%d err=%v", results, checks, violations, pending, pendingReasons, err)
		}
	})

	t.Run("expired lease fences old worker", func(t *testing.T) {
		id := enqueue()
		old, err := w.Claim(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if found, err := w.ProcessOne(ctx); err != nil || found {
			t.Fatal("live lease claimed twice", found, err)
		}

		if _, err = db.Exec(ctx, "UPDATE jobs SET leased_until=now()-interval '1 second' WHERE id=$1", old.ID); err != nil {
			t.Fatal(err)
		}

		recovered, err := w.Claim(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if recovered.SubmissionID != id || recovered.Token == old.Token {
			t.Fatal("job not recovered")
		}

		if err = w.Complete(ctx, old, "Wrong Answer", "stale", nil); !errors.Is(err, ErrLeaseLost) {
			t.Fatal("stale worker stored result", err)
		}

		if err = w.Complete(ctx, recovered, "Accepted", "recovered", nil); err != nil {
			t.Fatal(err)
		}

		if err = w.Complete(ctx, recovered, "Accepted", "duplicate", nil); !errors.Is(err, ErrLeaseLost) {
			t.Fatal("duplicate completion allowed", err)
		}

		var count int

		if err = db.QueryRow(ctx, "SELECT count(*) FROM judge_results WHERE submission_id=$1", id).Scan(&count); err != nil || count != 1 {
			t.Fatal(count, err)
		}
	})

	t.Run("analysis failure remains unsure without sending bans to model", func(t *testing.T) {
		id := enqueue()
		failed := Worker{DB: db, Runner: failedAnalyzer{}}

		if found, err := failed.ProcessOne(ctx); err != nil || !found {
			t.Fatal(found, err)
		}

		var uncertain, hard, pending int

		if err := db.QueryRow(ctx, `SELECT count(*) FILTER(WHERE verdict='unsure'),
		 count(*) FILTER(WHERE source='hard_rule'),count(*) FILTER(WHERE source='semantic_pending')
		 FROM constraint_verdicts WHERE submission_id=$1`, id).Scan(&uncertain, &hard, &pending); err != nil || uncertain != 2 || hard != 1 || pending != 1 {
			t.Fatal(uncertain, hard, pending, err)
		}
	})

	t.Run("Accepted semantic findings are stored separately with fake model", func(t *testing.T) {
		id := enqueue()
		calls := 0
		integrated := Worker{
			DB: db, Runner: w.Runner,
			SemanticFor: func(_, _ string) semantic.Checker {
				return semantic.Checker{Enabled: true, Model: fakeSemanticModel(
					func(_ context.Context, request semantic.Request) ([]byte, error) {
						calls++
						if len(request.Constraints) != 1 || request.Constraints[0].Text != "Use merge sort" ||
							strings.Contains(request.Source, "worker@test.invalid") {
							t.Fatalf("incorrect semantic request: %+v", request)
						}

						return json.Marshal(map[string]any{"verdicts": []any{map[string]any{
							"constraintId": request.Constraints[0].ID, "verdict": "violation",
							"evidence": []any{map[string]any{
								"line": 3, "text": strings.Split(request.Source, "\n")[2],
							}},
						}}})
					},
				)}
			},
		}

		if found, err := integrated.ProcessOne(ctx); err != nil || !found || calls != 1 {
			t.Fatal(found, calls, err)
		}

		var modelFindings, hardFindings int
		if err := db.QueryRow(ctx, `SELECT
 count(*) FILTER(WHERE source='model' AND verdict='violation'),
 count(*) FILTER(WHERE source='hard_rule')
 FROM constraint_verdicts WHERE submission_id=$1`, id).
			Scan(&modelFindings, &hardFindings); err != nil || modelFindings != 1 || hardFindings != 1 {
			t.Fatal(modelFindings, hardFindings, err)
		}
	})

	t.Run("lab reservations serialize and survive new worker instances", func(t *testing.T) {
		id := enqueue()
		budget := semantic.ReserveMicroUSD + 10
		var group sync.WaitGroup
		results := make(chan string, 2)

		for i := 0; i < 2; i++ {
			group.Add(1)
			go func() {
				defer group.Done()
				callID := uuid.NewString()
				if err := reserveCall(ctx, db, callID, lab, id, budget); err == nil {
					results <- callID
				} else {
					results <- ""
				}
			}()
		}
		group.Wait()
		close(results)

		winner, successes := "", 0
		for result := range results {
			if result != "" {
				winner = result
				successes++
			}
		}
		if successes != 1 {
			t.Fatalf("concurrent budget admitted %d calls", successes)
		}

		if err := recordCall(ctx, db, winner, semantic.CallLog{
			UsageKnown: true, CostMicroUSD: 7,
		}); err != nil {
			t.Fatal(err)
		}

		// Known usage releases only the unused reservation. A new instance can
		// reserve the remaining amount; it cannot reset the database ledger.
		unknown := uuid.NewString()
		if err := reserveCall(ctx, db, unknown, lab, id, budget); err != nil {
			t.Fatal(err)
		}
		if err := recordCall(ctx, db, unknown, semantic.CallLog{Outcome: "timeout"}); err != nil {
			t.Fatal(err)
		}
		if err := reserveCall(ctx, db, uuid.NewString(), lab, id, budget); err == nil {
			t.Fatal("unknown usage did not retain its reservation")
		}
		if err := recordCall(ctx, db, winner, semantic.CallLog{UsageKnown: true}); err == nil {
			t.Fatal("audit replay overwrote an existing charge")
		}

		var charged int64
		if err := db.QueryRow(ctx,
			"SELECT sum(charged_micro_usd)::bigint FROM model_calls WHERE lab_id=$1", lab,
		).Scan(&charged); err != nil || charged != semantic.ReserveMicroUSD+7 {
			t.Fatal(charged, err)
		}
	})
}
