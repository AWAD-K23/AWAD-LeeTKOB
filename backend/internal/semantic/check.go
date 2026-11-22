// Package semantic validates model findings before they can reach TA review.
package semantic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"leetkob/backend/internal/hardrules"
)

// Request deliberately excludes account IDs, deductions, and private tests.
// Source must have comments removed before it is passed to a provider.
type Request struct {
	Statement   string                 `json:"statement"`
	Constraints []hardrules.Constraint `json:"constraints"`
	Source      string                 `json:"code"`
}

// Model is replaceable so guardrail tests never need network access.
// Implementations must honor context cancellation and bound response sizes.
type Model interface {
	Generate(context.Context, Request) ([]byte, error)
}

// Checker is off unless explicitly enabled with a configured model.
type Checker struct {
	Model   Model
	Enabled bool
	Timeout time.Duration
}

// Check replaces only pending semantic findings. Hard-rule findings survive
// unchanged, including uncertain hard-rule results after an analysis failure.
func (c Checker) Check(
	ctx context.Context,
	statement string,
	source []byte,
	constraints []hardrules.Constraint,
	checks []hardrules.Verdict,
) []hardrules.Verdict {
	out := append([]hardrules.Verdict(nil), checks...)
	pending := []hardrules.Constraint{}
	positions := map[string]int{}

	for i, check := range out {
		if check.Source != "semantic_pending" {
			continue
		}

		for _, constraint := range constraints {
			if constraint.ID == check.ConstraintID && !hardrules.IsBanList(constraint.Text) {
				pending = append(pending, constraint)
				positions[constraint.ID] = i
				break
			}
		}
	}

	if !c.Enabled || c.Model == nil || len(pending) == 0 {
		return out
	}

	// Failed attempts remain unsure but are distinguished from offline work.
	for _, constraint := range pending {
		out[positions[constraint.ID]] = hardrules.Verdict{
			ConstraintID: constraint.ID,
			Verdict:      "unsure",
			Evidence:     []hardrules.Evidence{},
			Source:       "model",
		}
	}

	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 20 * time.Second
	}

	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cleaned, err := stripComments(string(source))
	if err != nil {
		return out
	}

	raw, err := c.Model.Generate(callCtx, Request{
		Statement: statement, Constraints: pending, Source: cleaned,
	})
	if err != nil || callCtx.Err() != nil {
		return out
	}

	validated, err := validate(raw, pending, cleaned, string(source))
	if err != nil {
		return out
	}

	for _, verdict := range validated {
		out[positions[verdict.ConstraintID]] = verdict
	}

	return out
}

// The provider cannot add deductions or invent extra fields, IDs, or evidence.
// Reject the entire response if any constraint is missing or malformed.
func validate(
	raw []byte,
	constraints []hardrules.Constraint,
	cleaned string,
	original string,
) ([]hardrules.Verdict, error) {
	invalid := errors.New("invalid model response")
	if len(raw) > 256*1024 {
		return nil, invalid
	}

	var response struct {
		Verdicts []struct {
			ConstraintID string                `json:"constraintId"`
			Verdict      string                `json:"verdict"`
			Evidence     *[]hardrules.Evidence `json:"evidence"`
		} `json:"verdicts"`
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&response); err != nil {
		return nil, invalid
	}

	var extra any
	if decoder.Decode(&extra) != io.EOF || len(response.Verdicts) != len(constraints) {
		return nil, invalid
	}

	remaining := map[string]bool{}
	for _, constraint := range constraints {
		remaining[constraint.ID] = true
	}

	cleanLines := strings.Split(cleaned, "\n")
	originalLines := strings.Split(original, "\n")
	out := []hardrules.Verdict{}

	for _, item := range response.Verdicts {
		if !remaining[item.ConstraintID] || item.Evidence == nil {
			return nil, invalid
		}
		delete(remaining, item.ConstraintID)

		if item.Verdict != "pass" && item.Verdict != "violation" && item.Verdict != "unsure" {
			return nil, invalid
		}
		if item.Verdict != "unsure" && len(*item.Evidence) == 0 {
			return nil, invalid
		}

		evidence := []hardrules.Evidence{}
		for _, point := range *item.Evidence {
			if point.Line < 1 || point.Line > len(cleanLines) ||
				strings.TrimSpace(point.Text) == "" || point.Text != cleanLines[point.Line-1] {
				return nil, invalid
			}

			// Store exact original lines so the TA sees the submitted source.
			evidence = append(evidence, hardrules.Evidence{
				Line: point.Line, Text: originalLines[point.Line-1],
			})
		}

		out = append(out, hardrules.Verdict{
			ConstraintID: item.ConstraintID, Verdict: item.Verdict,
			Evidence: evidence, Source: "model",
		})
	}

	return out, nil
}
