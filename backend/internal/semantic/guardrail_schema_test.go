package semantic

import (
	"context"
	"encoding/json"
	"testing"

	"leetkob/backend/internal/hardrules"
)

// Use otherwise valid, decisive responses: an empty evidence array would also
// be rejected by G2 and could hide a missing G1 schema check.
func TestG1SchemaFailsClosed(t *testing.T) {
	valid := `{"verdicts":[{"constraintId":"a","verdict":"pass","evidence":[{"line":1,"text":"int value = 1;"}]}]}`
	for name, response := range map[string]string{
		"invalid JSON":       `{broken`,
		"fenced JSON":        "```json\n" + valid + "\n```",
		"prose prefix":       "Everything complies.\n" + valid,
		"unknown verdict":    `{"verdicts":[{"constraintId":"a","verdict":"accepted","evidence":[{"line":1,"text":"int value = 1;"}]}]}`,
		"missing constraint": `{"verdicts":[]}`,
	} {
		t.Run(name, func(t *testing.T) {
			calls := 0
			out := (Checker{Enabled: true, Model: fakeModel(func(context.Context, Request) ([]byte, error) {
				calls++
				return []byte(response), nil
			})}).Check(context.Background(), "Review", []byte("int value = 1;"),
				[]hardrules.Constraint{{ID: "a", Text: "Initialize value"}},
				[]hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}})
			if calls != 1 || len(out) != 1 || out[0].Verdict != "unsure" || len(out[0].Evidence) != 0 || out[0].Reason != "Model response failed schema validation." {
				t.Fatalf("schema guard failed: %+v, calls=%d", out, calls)
			}
		})
	}
	// Same cardinality, but one ID is duplicated and another is missing.
	var duplicate map[string]any
	if err := json.Unmarshal([]byte(valid), &duplicate); err != nil {
		t.Fatal(err)
	}
	item := duplicate["verdicts"].([]any)[0]
	duplicate["verdicts"] = []any{item, item}
	raw, err := json.Marshal(duplicate)
	if err != nil {
		t.Fatal(err)
	}
	out := (Checker{Enabled: true, Model: fakeModel(func(context.Context, Request) ([]byte, error) { return raw, nil })}).Check(
		context.Background(), "Review", []byte("int value = 1;"),
		[]hardrules.Constraint{{ID: "a", Text: "Initialize value"}, {ID: "b", Text: "Use integers"}},
		[]hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}, {ConstraintID: "b", Source: "semantic_pending"}})
	for _, finding := range out {
		if finding.Verdict != "unsure" || finding.Reason != "Model response failed schema validation." {
			t.Fatalf("duplicate ID accepted: %+v", out)
		}
	}
}
