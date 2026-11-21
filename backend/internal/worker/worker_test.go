package worker

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/sandbox"
)

type fakeRunner struct{}

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
		var results, checks, violations, pending int
		var verdict string
		if err := db.QueryRow(ctx, "SELECT verdict FROM judge_results WHERE submission_id=$1", id).Scan(&verdict); err != nil || verdict != "Accepted" {
			t.Fatal(verdict, err)
		}
		if err := db.QueryRow(ctx, `SELECT (SELECT count(*) FROM judge_results WHERE submission_id=$1),
   count(*),count(*) FILTER(WHERE verdict='violation'),count(*) FILTER(WHERE source='semantic_pending')
   FROM constraint_verdicts WHERE submission_id=$1`, id).Scan(&results, &checks, &violations, &pending); err != nil || results != 1 || checks != 2 || violations != 1 || pending != 1 {
			t.Fatalf("results=%d checks=%d violations=%d pending=%d err=%v", results, checks, violations, pending, err)
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
}
