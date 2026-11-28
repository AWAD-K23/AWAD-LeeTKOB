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
	Statement          string                 `json:"statement"`
	Constraints        []hardrules.Constraint `json:"constraints"`
	Source             string                 `json:"code"`
	ProposedViolations []ProposedViolation    `json:"proposedViolations,omitempty"`
}

// Only already validated, sanitized model findings enter the confirmation pass.
// No TA decision, mark, identity or private test is included.
type ProposedViolation struct {
	ConstraintID string               `json:"constraintId"`
	Evidence     []hardrules.Evidence `json:"evidence"`
}

// Model is replaceable so guardrail tests never need network access.
// Implementations must honor context cancellation and bound response sizes.
type Model interface {
	Generate(context.Context, Request) ([]byte, error)
}

// Model-response failures are categorized locally so the model cannot invent
// its own uncertainty explanation.
var (
	ErrInvalidModelResponse = errors.New("invalid model response")
	ErrEvidenceMismatch     = errors.New("model evidence mismatch")
	ErrBudgetExceeded       = errors.New("lab model budget exhausted")
)

// Checker is off unless explicitly enabled with a configured model.
type Checker struct {
	Model   Model
	Enabled bool
	Timeout time.Duration
	// Known account identifiers are removed from all prompt text locally.
	Redactions []string
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

	if len(pending) == 0 {
		return out
	}

	for _, constraint := range pending {
		finding := &out[positions[constraint.ID]]
		finding.Verdict = "unsure"
		finding.Evidence = []hardrules.Evidence{}
		finding.Reason = "Semantic model check has not run."
		if !c.Enabled {
			finding.Reason = "Semantic model checking is disabled."
		} else if c.Model == nil {
			finding.Reason = "Semantic model is not configured."
		} else {
			finding.Source = "model"
			finding.Reason = "Model request failed."
		}
	}
	if !c.Enabled || c.Model == nil {
		return out
	}

	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 20 * time.Second
	}

	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cleaned, err := stripComments(string(source))
	if err != nil {
		for _, constraint := range pending {
			out[positions[constraint.ID]].Reason = "Source could not be prepared for model review."
		}
		return out
	}
	cleaned = redact(cleaned, c.Redactions)
	statement = redact(statement, c.Redactions)
	for i := range pending {
		pending[i].Text = redact(pending[i].Text, c.Redactions)
	}

	raw, err := c.Model.Generate(callCtx, Request{
		Statement: statement, Constraints: pending, Source: cleaned,
	})
	if err != nil || callCtx.Err() != nil {
		reason := "Model request failed."
		switch {
		case errors.Is(callCtx.Err(), context.DeadlineExceeded), errors.Is(err, context.DeadlineExceeded):
			reason = "Model request timed out."
		case errors.Is(callCtx.Err(), context.Canceled), errors.Is(err, context.Canceled):
			reason = "Model request was cancelled."
		case errors.Is(err, ErrBudgetExceeded):
			reason = "Lab model budget exhausted."
		}
		for _, constraint := range pending {
			out[positions[constraint.ID]].Reason = reason
		}
		return out
	}

	validated, err := validate(raw, pending, cleaned, string(source))
	if err != nil {
		reason := "Model response failed schema validation."
		if errors.Is(err, ErrEvidenceMismatch) {
			reason = "Model evidence did not match submitted code."
		}
		for _, constraint := range pending {
			out[positions[constraint.ID]].Reason = reason
		}
		return out
	}

	validated = c.confirmViolations(ctx, timeout, Request{
		Statement: statement, Constraints: pending, Source: cleaned,
	}, string(source), validated)
	for _, verdict := range validated {
		if verdict.Verdict == "unsure" && verdict.Reason == "" {
			verdict.Reason = "Model returned unsure."
		}
		out[positions[verdict.ConstraintID]] = verdict
	}

	return out
}

// Preserve newlines and byte positions for subsequent evidence validation.
func redact(text string, identities []string) string {
	for _, identity := range identities {
		if identity != "" && !strings.ContainsAny(identity, "\r\n") {
			text = strings.ReplaceAll(text, identity, strings.Repeat(" ", len(identity)))
		}
	}
	return text
}

// The provider cannot add deductions or invent extra fields, IDs, or evidence.
// Reject malformed response structure or missing/duplicate IDs as a whole.
// Evidence failures make that finding unsure without discarding independently
// valid findings. No decisive result survives without exact nonblank evidence.
func validate(
	raw []byte,
	constraints []hardrules.Constraint,
	cleaned string,
	original string,
) ([]hardrules.Verdict, error) {
	invalid := ErrInvalidModelResponse
	if len(raw) > 256*1024 {
		return nil, invalid
	}
	// encoding/json otherwise accepts duplicate keys and case-insensitive
	// struct field names, neither of which belongs to our response contract.
	if err := strictKeys(json.NewDecoder(bytes.NewReader(raw)), 0); err != nil {
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
			out = append(out, hardrules.Verdict{
				ConstraintID: item.ConstraintID, Verdict: "unsure", Source: "model",
				Evidence: []hardrules.Evidence{}, Reason: "Model response failed schema validation.",
			})
			continue
		}

		evidence := []hardrules.Evidence{}
		invalidEvidence := false
		for _, point := range *item.Evidence {
			if point.Line < 1 || point.Line > len(cleanLines) ||
				strings.TrimSpace(point.Text) == "" || point.Text != cleanLines[point.Line-1] {
				invalidEvidence = true
				break
			}

			// Store exact original lines so the TA sees the submitted source.
			evidence = append(evidence, hardrules.Evidence{
				Line: point.Line, Text: originalLines[point.Line-1],
			})
		}
		if invalidEvidence {
			out = append(out, hardrules.Verdict{
				ConstraintID: item.ConstraintID, Verdict: "unsure", Source: "model",
				Evidence: []hardrules.Evidence{}, Reason: "Model evidence did not match submitted code.",
			})
			continue
		}

		reason := ""
		if item.Verdict == "unsure" {
			reason = "Model returned unsure."
		}
		out = append(out, hardrules.Verdict{
			ConstraintID: item.ConstraintID, Verdict: item.Verdict,
			Reason: reason, Evidence: evidence, Source: "model",
		})
	}

	return out, nil
}

func strictKeys(decoder *json.Decoder, depth int) error {
	invalid := errors.New("invalid JSON keys")
	if depth > 16 {
		return invalid
	}

	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, composite := token.(json.Delim)
	if !composite {
		return nil
	}

	seen := map[string]bool{}
	for decoder.More() {
		if delimiter == '{' {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok || seen[key] {
				return invalid
			}
			seen[key] = true
			switch key {
			case "verdicts", "constraintId", "verdict", "evidence", "line", "text":
			default:
				return invalid
			}
		}

		if err := strictKeys(decoder, depth+1); err != nil {
			return err
		}
	}

	_, err = decoder.Token()
	return err
}
