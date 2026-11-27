package semantic

import (
	"context"
	"testing"
	"time"

	"leetkob/backend/internal/hardrules"
)

func TestG3DeadlineDiscardsLateSuccess(t *testing.T) {
	calls := 0
	model := fakeModel(func(ctx context.Context, _ Request) ([]byte, error) {
		calls++
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("provider has no deadline")
		}
		<-ctx.Done()
		// Even an otherwise valid result arriving after cancellation is unsafe.
		return []byte(`{"verdicts":[{"constraintId":"a","verdict":"pass","evidence":[{"line":1,"text":"int value = 1;"}]}]}`), nil
	})
	start := time.Now()
	out := (Checker{Enabled: true, Model: model, Timeout: 10 * time.Millisecond}).Check(context.Background(), "Review", []byte("int value = 1;"),
		[]hardrules.Constraint{{ID: "a", Text: "Initialize value"}}, []hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}})
	if calls != 1 || time.Since(start) > time.Second || out[0].Verdict != "unsure" || len(out[0].Evidence) != 0 || out[0].Reason != "Model request timed out." {
		t.Fatalf("late result escaped timeout guard: %+v", out)
	}
}
