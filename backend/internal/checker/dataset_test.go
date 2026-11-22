package checker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"leetkob/backend/internal/sandbox"
)

type fakeScanner struct{ fail bool }

func (s fakeScanner) Run(context.Context, []byte, string, string) (sandbox.Result, error) {
	if s.fail {
		return sandbox.Result{}, errors.New("analyzer unavailable")
	}
	return sandbox.Result{Verdict: "Executed", Output: "for\t1\n"}, nil
}

func TestDatasetWithoutLabelsOrProblemSettings(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"problems/unseen-name/statement.md": "Calculate a total.",
		"problems/unseen-name/constraints.json": `[
		 {"id":"ban","text":"Do not use loop statements: ` + "`for`" + `","deduction":3},
		 {"id":"method","text":"Use recursion","deduction":2}]`,
		"solutions/unseen-name/example.cpp": "int main(){for(;;){break;}}",
	}
	for path, text := range files {
		fullPath := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}

	for _, fail := range []bool{false, true} {
		var output bytes.Buffer
		if err := (Checker{Scanner: fakeScanner{fail: fail}}).Run(context.Background(), root, &output); err != nil {
			t.Fatal(err)
		}
		decoder := json.NewDecoder(&output)
		var ban, method Prediction
		if decoder.Decode(&ban) != nil || decoder.Decode(&method) != nil {
			t.Fatal("missing predictions")
		}
		if ban.Solution != "unseen-name/example" || ban.Constraint != "ban" ||
			method.Constraint != "method" || method.Verdict != "unsure" || method.Evidence == nil {
			t.Fatal(ban, method)
		}
		if !fail && (ban.Verdict != "violation" || len(ban.Evidence) != 1 || ban.Evidence[0].Line != 1) {
			t.Fatal("lost syntax evidence", ban)
		}
		if fail && (ban.Verdict != "unsure" || len(ban.Evidence) != 0) {
			t.Fatal("analysis failure became a decision", ban)
		}
	}
}
