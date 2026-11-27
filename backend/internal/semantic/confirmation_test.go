package semantic

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"leetkob/backend/internal/hardrules"
)

func TestViolationConfirmationFailsClosed(t *testing.T) {
	for _, tc := range []struct {
		name, decision, reason string
		err                    error
	}{
		{"confirmed", "violation", "", nil},
		{"disagreed", "pass", "Model violation was not confirmed.", nil},
		{"uncertain", "unsure", "Model returned unsure.", nil},
		{"malformed", "broken", "Model confirmation failed schema validation.", nil},
		{"budget", "", "Lab model budget exhausted.", ErrBudgetExceeded},
		{"provider error", "", "Model violation confirmation failed.", errors.New("provider failed")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			model := fakeModel(func(_ context.Context, request Request) ([]byte, error) {
				calls++
				if calls == 1 {
					if len(request.ProposedViolations) != 0 {
						t.Fatal("primary request includes a proposal")
					}
					return confirmationResponse(request, "violation"), nil
				}
				if len(request.ProposedViolations) != 1 || len(request.Constraints) != 1 {
					t.Fatalf("wrong confirmation population: %+v", request)
				}
				if strings.Contains(request.ProposedViolations[0].Evidence[0].Text, "private") || strings.Contains(request.Source, "private") {
					t.Fatal("original comment leaked into confirmation")
				}
				if tc.err != nil {
					return nil, tc.err
				}
				if tc.decision == "broken" {
					return []byte("not JSON"), nil
				}
				return confirmationResponse(request, tc.decision), nil
			})
			out := (Checker{Enabled: true, Model: model}).Check(context.Background(), "Review", []byte("int quantity = 1; // private comment"),
				[]hardrules.Constraint{{ID: "a", Text: "Use a floating point quantity"}}, []hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}})
			if tc.decision == "violation" {
				if out[0].Verdict != "violation" || len(out[0].Evidence) != 1 || out[0].Evidence[0].Text != "int quantity = 1; // private comment" {
					t.Fatal(out)
				}
			} else if out[0].Verdict != "unsure" || len(out[0].Evidence) != 0 || out[0].Reason != tc.reason {
				t.Fatalf("unconfirmed violation escaped: %+v", out)
			}
			if calls != 2 {
				t.Fatalf("confirmation calls=%d", calls)
			}
		})
	}
}

func TestConfirmationOnlyReviewsViolations(t *testing.T) {
	calls := 0
	model := fakeModel(func(_ context.Context, request Request) ([]byte, error) {
		calls++
		return confirmationResponse(request, "pass"), nil
	})
	out := (Checker{Enabled: true, Model: model}).Check(context.Background(), "Review", []byte("int quantity = 1;"),
		[]hardrules.Constraint{{ID: "a", Text: "Use a quantity"}}, []hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}})
	if calls != 1 || out[0].Verdict != "pass" {
		t.Fatalf("unneeded confirmation: %d %+v", calls, out)
	}
}

func TestConfirmationDiscardsLateSuccess(t *testing.T) {
	model := fakeModel(func(ctx context.Context, request Request) ([]byte, error) {
		if len(request.ProposedViolations) > 0 {
			<-ctx.Done()
		}
		return confirmationResponse(request, "violation"), nil
	})
	out := (Checker{Enabled: true, Model: model, Timeout: 5 * time.Millisecond}).Check(context.Background(), "Review", []byte("int quantity = 1;"),
		[]hardrules.Constraint{{ID: "a", Text: "Use a floating point quantity"}}, []hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}})
	if out[0].Verdict != "unsure" || out[0].Reason != "Model request timed out." {
		t.Fatal(out)
	}
}

func TestConfirmationDoesNotRestoreRedactedIdentity(t *testing.T) {
	calls := 0
	model := fakeModel(func(_ context.Context, request Request) ([]byte, error) {
		calls++
		encoded, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(encoded), "known-account-42") {
			t.Fatal("redacted identity restored in model request")
		}
		return confirmationResponse(request, "violation"), nil
	})
	out := (Checker{Enabled: true, Model: model, Redactions: []string{"known-account-42"}}).Check(context.Background(), "Review", []byte(`const char* tag = "known-account-42";`),
		[]hardrules.Constraint{{ID: "a", Text: "Use a numeric tag"}}, []hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}})
	if calls != 2 || out[0].Verdict != "violation" {
		t.Fatalf("confirmation not exercised: %d %+v", calls, out)
	}
	if !strings.Contains(out[0].Evidence[0].Text, "known-account-42") {
		t.Fatal("TA's original source mapping lost")
	}
}

func confirmationResponse(request Request, decision string) []byte {
	items := []any{}
	for _, constraint := range request.Constraints {
		evidence := []hardrules.Evidence{}
		if decision != "unsure" {
			evidence = append(evidence, hardrules.Evidence{Line: 1, Text: strings.Split(request.Source, "\n")[0]})
		}
		items = append(items, map[string]any{"constraintId": constraint.ID, "verdict": decision, "evidence": evidence})
	}
	raw, _ := json.Marshal(map[string]any{"verdicts": items})
	return raw
}
