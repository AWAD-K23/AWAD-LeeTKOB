package semantic

import (
	"context"
	"encoding/json"
	"testing"

	"leetkob/backend/internal/hardrules"
)

// A failed finding must stay unsure without discarding a different constraint's
// exact, valid evidence. Overall schema/ID corruption still rejects the response.
func TestIndependentFindingEvidenceFailures(t *testing.T) {
	for _, tc := range []struct {
		name, reason string
		evidence     []hardrules.Evidence
	}{
		{"missing decisive evidence", "Model response failed schema validation.", []hardrules.Evidence{}},
		{"invented evidence", "Model evidence did not match submitted code.", []hardrules.Evidence{{Line: 1, Text: "fabricated"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := fakeModel(func(_ context.Context, request Request) ([]byte, error) {
				if len(request.ProposedViolations) > 0 {
					return confirmationResponse(request, "violation"), nil
				}
				return json.Marshal(map[string]any{"verdicts": []any{
					map[string]any{"constraintId": "first", "verdict": "violation", "evidence": []hardrules.Evidence{{Line: 1, Text: "int value = 1;"}}},
					map[string]any{"constraintId": "second", "verdict": "pass", "evidence": tc.evidence},
				}})
			})
			out := (Checker{Enabled: true, Model: model}).Check(context.Background(), "Review", []byte("int value = 1;"),
				[]hardrules.Constraint{{ID: "first", Text: "Use floating point"}, {ID: "second", Text: "Initialize value"}},
				[]hardrules.Verdict{{ConstraintID: "first", Source: "semantic_pending"}, {ConstraintID: "second", Source: "semantic_pending"}})
			if out[0].Verdict != "violation" || len(out[0].Evidence) != 1 || out[0].Evidence[0].Text != "int value = 1;" {
				t.Fatalf("valid independent finding discarded: %+v", out)
			}
			if out[1].Verdict != "unsure" || len(out[1].Evidence) != 0 || out[1].Reason != tc.reason {
				t.Fatalf("invalid finding escaped guard: %+v", out)
			}
		})
	}
}
