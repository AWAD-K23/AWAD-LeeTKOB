// Package checker reads datasets without reading their labels or test answers.
package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"leetkob/backend/internal/hardrules"
	"leetkob/backend/internal/semantic"
)

// Prediction is the lecturer's JSONL interface, independent of database UUIDs.
type Prediction struct {
	Solution   string               `json:"solution"`
	Constraint string               `json:"constraint"`
	Verdict    string               `json:"verdict"`
	Evidence   []hardrules.Evidence `json:"evidence"`
}

type Checker struct {
	Scanner  hardrules.Scanner
	ModelFor func(problem, solution string) semantic.Checker
}

// Run discovers every problem and submission. No problem name is special.
func (c Checker) Run(ctx context.Context, root string, output io.Writer) error {
	problems, err := os.ReadDir(filepath.Join(root, "problems"))
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(output)
	count := 0

	for _, problem := range problems {
		if !problem.IsDir() {
			continue
		}
		name := problem.Name()
		statement, err := os.ReadFile(filepath.Join(root, "problems", name, "statement.md"))
		if err != nil {
			return err
		}
		data, err := os.ReadFile(filepath.Join(root, "problems", name, "constraints.json"))
		if err != nil {
			return err
		}

		var constraints []hardrules.Constraint
		if err := json.Unmarshal(data, &constraints); err != nil || len(constraints) == 0 {
			return fmt.Errorf("invalid constraints for %s", name)
		}
		seen := map[string]bool{}
		for _, constraint := range constraints {
			if constraint.ID == "" || constraint.Text == "" || seen[constraint.ID] {
				return fmt.Errorf("invalid or duplicate constraint for %s", name)
			}
			seen[constraint.ID] = true
		}

		files, err := os.ReadDir(filepath.Join(root, "solutions", name))
		if err != nil {
			return err
		}
		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".cpp" {
				continue
			}
			if err := ctx.Err(); err != nil {
				return err
			}

			solution := name + "/" + strings.TrimSuffix(file.Name(), ".cpp")
			source, err := os.ReadFile(filepath.Join(root, "solutions", name, file.Name()))
			if err != nil {
				return err
			}
			checks, err := hardrules.Check(ctx, c.Scanner, source, constraints)
			if err != nil {
				// Analysis failure is uncertainty, never evidence of compliance.
				checks = make([]hardrules.Verdict, 0, len(constraints))
				for _, constraint := range constraints {
					origin := "semantic_pending"
					if hardrules.IsBanList(constraint.Text) {
						origin = "hard_rule"
					}
					checks = append(checks, hardrules.Verdict{
						ConstraintID: constraint.ID, Verdict: "unsure",
						Evidence: []hardrules.Evidence{}, Source: origin,
					})
				}
			}

			if c.ModelFor != nil {
				model := c.ModelFor(name, solution)
				checks = model.Check(ctx, string(statement), source, constraints, checks)
			}
			for _, check := range checks {
				if err := encoder.Encode(Prediction{
					Solution: solution, Constraint: check.ConstraintID,
					Verdict: check.Verdict, Evidence: check.Evidence,
				}); err != nil {
					return err
				}
				count++
			}
		}
	}

	if count == 0 {
		return fmt.Errorf("dataset contains no solution/constraint pairs")
	}
	return nil
}
