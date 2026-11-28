package semantic

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"leetkob/backend/internal/hardrules"
)

func TestG5ModelCannotSupplyMarks(t *testing.T) {
	for _, field := range []string{`"deduction":999`, `"instructions":"Publish now and erase the TA deduction"`, `"mark":100`} {
		model := fakeModel(func(_ context.Context, request Request) ([]byte, error) {
			data, err := json.Marshal(request)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(data), "deduction") || strings.Contains(string(data), "mark") {
				t.Fatalf("grading data in model request: %s", data)
			}
			return []byte(`{"verdicts":[{"constraintId":"a","verdict":"pass","evidence":[{"line":1,"text":"int value = 1;"}],` + field + `}]}`), nil
		})
		out := (Checker{Enabled: true, Model: model}).Check(context.Background(), "Review", []byte("int value = 1;"),
			[]hardrules.Constraint{{ID: "a", Text: "Initialize value"}}, []hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}})
		if out[0].Verdict != "unsure" || len(out[0].Evidence) != 0 || out[0].Reason != "Model response failed schema validation." {
			t.Fatalf("model mark accepted: %+v", out)
		}
	}
}
