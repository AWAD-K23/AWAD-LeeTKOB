package semantic

import (
	"context"
	"strings"
	"testing"
	"time"

	"leetkob/backend/internal/hardrules"
)

type fakeModel func(context.Context, Request) ([]byte, error)

func (f fakeModel) Generate(ctx context.Context, request Request) ([]byte, error) {
	return f(ctx, request)
}

func TestCheckerBoundary(t *testing.T) {
	constraints := []hardrules.Constraint{
		{ID: "ban", Text: "Do not call `std::sort`"},
		{ID: "algorithm", Text: "Use a comparison"},
	}
	checks := []hardrules.Verdict{
		{ConstraintID: "ban", Verdict: "violation", Source: "hard_rule"},
		{ConstraintID: "algorithm", Verdict: "unsure", Source: "semantic_pending"},
	}
	source := []byte("// student identity\nint main() {} // private comment")
	calls := 0
	model := fakeModel(func(_ context.Context, _ Request) ([]byte, error) {
		calls++
		t.Fatal("offline checker invoked model")
		return nil, nil
	})

	// Offline mode must not call even an available model.
	offline := (Checker{Model: model}).Check(context.Background(), "Compare", source, constraints, checks)
	if calls != 0 || offline[1].Source != "semantic_pending" ||
		offline[1].Reason != "Semantic model checking is disabled." {
		t.Fatalf("offline mode called model or lost its reason: %+v", offline[1])
	}

	// Build exact evidence from the cleaned source rather than relying on spaces.
	model = fakeModel(func(_ context.Context, request Request) ([]byte, error) {
		calls++
		if len(request.Constraints) != 1 || request.Constraints[0].ID != "algorithm" ||
			strings.Contains(request.Source, "private") || strings.Contains(request.Source, "identity") {
			t.Fatalf("unexpected model request: %+v", request)
		}
		return []byte(`{"verdicts":[{"constraintId":"algorithm","verdict":"violation","evidence":[{"line":2,"text":"` +
			strings.Split(request.Source, "\n")[1] + `"}]}]}`), nil
	})
	out := (Checker{Enabled: true, Model: model}).Check(
		context.Background(), "Compare", source, constraints, checks,
	)
	if calls != 1 || out[0].Verdict != "violation" || out[0].Source != "hard_rule" ||
		out[1].Verdict != "violation" || out[1].Source != "model" ||
		out[1].Evidence[0].Text != strings.Split(string(source), "\n")[1] {
		t.Fatalf("unexpected verdicts: %+v", out)
	}

	// The caller's input must not be mutated by the semantic checker.
	if checks[1].Source != "semantic_pending" {
		t.Fatal("input findings changed")
	}
}

func TestInvalidResponsesBecomeUnsure(t *testing.T) {
	responses := map[string]string{
		"free text":        "the student passes",
		"missing result":   `{"verdicts":[]}`,
		"unknown ID":       `{"verdicts":[{"constraintId":"other","verdict":"pass","evidence":[]}]}`,
		"deduction":        `{"verdicts":[{"constraintId":"a","verdict":"pass","evidence":[],"deduction":3}]}`,
		"bad verdict":      `{"verdicts":[{"constraintId":"a","verdict":"accepted","evidence":[]}]}`,
		"no evidence":      `{"verdicts":[{"constraintId":"a","verdict":"violation","evidence":[]}]}`,
		"missing field":    `{"verdicts":[{"constraintId":"a","verdict":"unsure"}]}`,
		"null evidence":    `{"verdicts":[{"constraintId":"a","verdict":"unsure","evidence":null}]}`,
		"duplicate key":    `{"verdicts":[{"constraintId":"a","verdict":"pass","verdict":"unsure","evidence":[]}]}`,
		"wrong key case":   `{"verdicts":[{"constraintId":"a","Verdict":"unsure","evidence":[]}]}`,
		"invented line":    `{"verdicts":[{"constraintId":"a","verdict":"violation","evidence":[{"line":99,"text":"code"}]}]}`,
		"invented text":    `{"verdicts":[{"constraintId":"a","verdict":"violation","evidence":[{"line":1,"text":"invented"}]}]}`,
		"extra document":   `{"verdicts":[{"constraintId":"a","verdict":"unsure","evidence":[]}]} {}`,
		"oversized output": strings.Repeat(" ", 256*1024+1),
	}

	for name, response := range responses {
		t.Run(name, func(t *testing.T) {
			model := fakeModel(func(context.Context, Request) ([]byte, error) {
				return []byte(response), nil
			})
			out := (Checker{Enabled: true, Model: model}).Check(
				context.Background(), "Statement", []byte("int main() {}"),
				[]hardrules.Constraint{{ID: "a", Text: "Use recursion"}},
				[]hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}},
			)
			if out[0].Verdict != "unsure" || len(out[0].Evidence) != 0 || out[0].Reason == "" {
				t.Fatalf("invalid output accepted or unexplained: %+v", out[0])
			}
		})
	}
}

func TestTimeoutBecomesUnsure(t *testing.T) {
	model := fakeModel(func(ctx context.Context, _ Request) ([]byte, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	})
	checker := Checker{Enabled: true, Model: model, Timeout: 5 * time.Millisecond}
	out := checker.Check(
		context.Background(), "Statement", []byte("int main() {}"),
		[]hardrules.Constraint{{ID: "a", Text: "Use recursion"}},
		[]hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}},
	)
	if out[0].Verdict != "unsure" || out[0].Source != "model" ||
		out[0].Reason != "Model request timed out." {
		t.Fatalf("timeout was not uncertain or explained: %+v", out)
	}
}

func TestUncertaintyReasons(t *testing.T) {
	constraints := []hardrules.Constraint{{ID: "a", Text: "Use recursion"}}
	checks := []hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}}
	source := []byte("int main() {}")
	for _, tc := range []struct {
		name, response, reason string
		err                    error
	}{
		{name: "model unsure", response: `{"verdicts":[{"constraintId":"a","verdict":"unsure","evidence":[]}]}`, reason: "Model returned unsure."},
		{name: "invalid schema", response: "not json", reason: "Model response failed schema validation."},
		{name: "evidence mismatch", response: `{"verdicts":[{"constraintId":"a","verdict":"violation","evidence":[{"line":1,"text":"invented"}]}]}`, reason: "Model evidence did not match submitted code."},
		{name: "budget exhausted", err: ErrBudgetExceeded, reason: "Lab model budget exhausted."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := fakeModel(func(context.Context, Request) ([]byte, error) { return []byte(tc.response), tc.err })
			out := (Checker{Enabled: true, Model: model}).Check(context.Background(), "Statement", source, constraints, checks)
			if out[0].Verdict != "unsure" || out[0].Reason != tc.reason {
				t.Fatalf("unexpected uncertain finding: %+v", out[0])
			}
		})
	}
}

func TestCommentRemoval(t *testing.T) {
	cases := []string{
		"// hidden identity\nint x; /* hidden\nidentity */",
		`const char* x = "https://example.test/*literal*/"; // hidden identity`,
		`auto x = R"tag(/*literal*/ // literal)tag"; /* hidden identity */`,
		"int x = 1'000; // hidden identity",
		"int x = 0xAB'CD; // hidden identity ' with extra quote",
		`char x = '\''; /* hidden identity */`,
	}

	for _, source := range cases {
		clean, err := stripComments(source)
		if err != nil || strings.Contains(clean, "hidden") || strings.Contains(clean, "identity") ||
			len(clean) != len(source) || strings.Count(clean, "\n") != strings.Count(source, "\n") {
			t.Fatalf("comment removal failed: %q, %v", clean, err)
		}
		if strings.Contains(source, "literal") && !strings.Contains(clean, "literal") {
			t.Fatal("literal was damaged")
		}
	}

	for _, source := range []string{"// hidden\\\nidentity", "/* unclosed", `R"(unclosed`} {
		if _, err := stripComments(source); err == nil {
			t.Fatalf("ambiguous source accepted: %q", source)
		}
	}
}

func TestKnownIdentityRedaction(t *testing.T) {
	model := fakeModel(func(_ context.Context, request Request) ([]byte, error) {
		data := request.Statement + request.Source + request.Constraints[0].Text
		if strings.Contains(data, "alice") || strings.Contains(data, "student-123") {
			t.Fatal("known identity leaked into prompt")
		}
		return []byte(`{"verdicts":[{"constraintId":"a","verdict":"unsure","evidence":[]}]}`), nil
	})
	checker := Checker{
		Enabled: true, Model: model,
		Redactions: []string{"alice@example.test", "alice", "student-123"},
	}
	out := checker.Check(
		context.Background(), "alice's assignment", []byte(`auto id = "student-123";`),
		[]hardrules.Constraint{{ID: "a", Text: "alice: use recursion"}},
		[]hardrules.Verdict{{ConstraintID: "a", Source: "semantic_pending"}},
	)
	if out[0].Source != "model" || out[0].Verdict != "unsure" {
		t.Fatal(out)
	}
}
