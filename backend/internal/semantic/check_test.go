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
	if calls != 0 || offline[1].Source != "semantic_pending" {
		t.Fatal("offline mode called model or changed pending state")
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
			if out[0].Verdict != "unsure" || len(out[0].Evidence) != 0 {
				t.Fatalf("invalid output accepted: %+v", out)
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
	if out[0].Verdict != "unsure" || out[0].Source != "model" {
		t.Fatalf("timeout was not uncertain: %+v", out)
	}
}

func TestCommentRemoval(t *testing.T) {
	cases := []string{
		"// hidden identity\nint x; /* hidden\nidentity */",
		`const char* x = "https://example.test/*literal*/"; // hidden identity`,
		`auto x = R"tag(/*literal*/ // literal)tag"; /* hidden identity */`,
		"int x = 1'000; // hidden identity",
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
