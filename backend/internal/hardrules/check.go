// Package hardrules handles ban lists without any model dependency.
package hardrules

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"leetkob/backend/internal/sandbox"
)

type Constraint struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type Evidence struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}
type Verdict struct {
	ConstraintID string     `json:"constraintId"`
	Verdict      string     `json:"verdict"`
	Evidence     []Evidence `json:"evidence"`
	Source       string     `json:"source"`
}
type Scanner interface {
	Run(context.Context, []byte, string, string) (sandbox.Result, error)
}

var backticks = regexp.MustCompile("`([^`]+)`")

func bannedNames(text string) []string {
	if !strings.HasPrefix(text, "Do not use") && !strings.HasPrefix(text, "Do not call") {
		return nil
	}
	names := []string{}
	for _, match := range backticks.FindAllStringSubmatch(text, -1) {
		names = append(names, strings.TrimPrefix(match[1], "::"))
	}
	return names
}

// All unresolved semantic work is explicit. Member 2 can call a model only for
// semantic_pending records, keeping hard-rule decisions out of model prompts.
func Check(ctx context.Context, scanner Scanner, source []byte, constraints []Constraint) ([]Verdict, error) {
	needsScan := false
	for _, c := range constraints {
		if len(bannedNames(c.Text)) > 0 {
			needsScan = true
		}
	}
	calls := map[string][]int{}
	if needsScan {
		result, err := scanner.Run(ctx, source, "", "analyze")
		if err != nil {
			return nil, err
		}
		if result.Verdict != "Executed" {
			return nil, fmt.Errorf("syntax analysis failed: %s", result.Verdict)
		}
		for _, line := range strings.Split(result.Output, "\n") {
			fields := strings.Split(line, "\t")
			if len(fields) != 2 {
				continue
			}
			number, err := strconv.Atoi(fields[1])
			if err == nil {
				calls[fields[0]] = append(calls[fields[0]], number)
			}
		}
	}
	lines := strings.Split(string(source), "\n")
	out := make([]Verdict, 0, len(constraints))
	for _, c := range constraints {
		v := Verdict{ConstraintID: c.ID, Verdict: "unsure", Evidence: []Evidence{}, Source: "semantic_pending"}
		names := bannedNames(c.Text)
		if len(names) > 0 {
			v.Verdict = "pass"
			v.Source = "hard_rule"
			found := map[int]bool{}
			for _, name := range names {
				for called, numbers := range calls {
					short := called
					if i := strings.LastIndex(short, "::"); i >= 0 {
						short = short[i+2:]
					}
					if called != name && (strings.Contains(name, "::") || short != name) {
						continue
					}
					for _, number := range numbers {
						if number >= 1 && number <= len(lines) && !found[number] {
							v.Verdict = "violation"
							v.Evidence = append(v.Evidence, Evidence{Line: number, Text: lines[number-1]})
							found[number] = true
						}
					}
				}
			}
		}
		out = append(out, v)
	}
	return out, nil
}
