package semantic

import (
	"context"
	"encoding/json"
	"testing"

	"leetkob/backend/internal/hardrules"
)

func TestG2ExactEvidenceAndOriginalMapping(t *testing.T) {
	source := "// int invented = 9;\nint first = 1; /* explanatory comment */\nint second = 2;"
	for _, tc := range []struct {
		name string
		line int
		text string
	}{
		{"invented line", 99, "int first = 1;"},
		{"invented text", 2, "int fabricated = 1;"},
		{"right text wrong line", 3, "int first = 1;"},
		{"stripped comment only", 1, "// int invented = 9;"},
		{"blank stripped line", 1, "                     "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := evidenceFinding(t, source, func(Request) hardrules.Evidence { return hardrules.Evidence{Line: tc.line, Text: tc.text} })
			if out.Verdict != "unsure" || len(out.Evidence) != 0 || out.Reason != "Model evidence did not match submitted code." {
				t.Fatalf("invalid evidence accepted: %+v", out)
			}
		})
	}
	// Positive control: the model quotes cleaned text, and the saved evidence
	// maps to the exact original line including its trailing comment.
	out := evidenceFinding(t, source, func(request Request) hardrules.Evidence {
		lines := numberLines(request.Source)
		if lines[1].Line != 2 {
			t.Fatal(lines)
		}
		return hardrules.Evidence{Line: 2, Text: lines[1].Text}
	})
	if out.Verdict != "violation" || len(out.Evidence) != 1 || out.Evidence[0].Line != 2 || out.Evidence[0].Text != "int first = 1; /* explanatory comment */" {
		t.Fatalf("original mapping lost: %+v", out)
	}
}

func evidenceFinding(t *testing.T, source string, point func(Request) hardrules.Evidence) hardrules.Verdict {
	t.Helper()
	model := fakeModel(func(_ context.Context, request Request) ([]byte, error) {
		return json.Marshal(map[string]any{"verdicts": []any{map[string]any{"constraintId": "a", "verdict": "violation", "evidence": []hardrules.Evidence{point(request)}}}})
	})
	out := (Checker{Enabled: true, Model: model}).Check(context.Background(), "Review", []byte(source),
		[]hardrules.Constraint{{ID: "a", Text: "Initialize a value"}}, []hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}})
	if len(out) != 1 {
		t.Fatal(out)
	}
	return out[0]
}
