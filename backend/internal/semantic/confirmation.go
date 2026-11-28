package semantic

import (
	"context"
	"errors"
	"strings"
	"time"

	"leetkob/backend/internal/hardrules"
)

// A second call tries to falsify proposed semantic violations. Disagreement,
// malformed output, timeout or exhausted budget cannot produce a violation.
// All candidates share one bounded confirmation call, rather than one per rule.
func (c Checker) confirmViolations(ctx context.Context, timeout time.Duration, request Request, original string, findings []hardrules.Verdict) []hardrules.Verdict {
	positions := map[string]int{}
	lines := strings.Split(request.Source, "\n")
	for i, finding := range findings {
		if finding.Verdict != "violation" {
			continue
		}
		positions[finding.ConstraintID] = i
		proposal := ProposedViolation{ConstraintID: finding.ConstraintID, Evidence: []hardrules.Evidence{}}
		for _, point := range finding.Evidence {
			// Validated findings carry original text. Rebuild sanitized text so
			// stripped comments and redacted identities cannot leak on recheck.
			proposal.Evidence = append(proposal.Evidence, hardrules.Evidence{Line: point.Line, Text: lines[point.Line-1]})
		}
		request.ProposedViolations = append(request.ProposedViolations, proposal)
	}
	if len(positions) == 0 {
		return findings
	}
	selected := []hardrules.Constraint{}
	for _, constraint := range request.Constraints {
		if _, ok := positions[constraint.ID]; ok {
			selected = append(selected, constraint)
		}
	}
	request.Constraints = selected
	unsure := func(reason string) []hardrules.Verdict {
		for _, i := range positions {
			findings[i].Verdict = "unsure"
			findings[i].Evidence = []hardrules.Evidence{}
			findings[i].Reason = reason
		}
		return findings
	}
	confirmCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	raw, err := c.Model.Generate(confirmCtx, request)
	if err != nil || confirmCtx.Err() != nil {
		reason := "Model violation confirmation failed."
		switch {
		case errors.Is(confirmCtx.Err(), context.DeadlineExceeded), errors.Is(err, context.DeadlineExceeded):
			reason = "Model request timed out."
		case errors.Is(confirmCtx.Err(), context.Canceled), errors.Is(err, context.Canceled):
			reason = "Model request was cancelled."
		case errors.Is(err, ErrBudgetExceeded):
			reason = "Lab model budget exhausted."
		}
		return unsure(reason)
	}
	confirmed, err := validate(raw, selected, request.Source, original)
	if err != nil {
		return unsure("Model confirmation failed schema validation.")
	}
	for _, finding := range confirmed {
		i := positions[finding.ConstraintID]
		if finding.Verdict == "violation" {
			findings[i] = finding
			continue
		}
		findings[i].Verdict = "unsure"
		findings[i].Evidence = []hardrules.Evidence{}
		findings[i].Reason = "Model violation was not confirmed."
		if finding.Verdict == "unsure" && finding.Reason != "" {
			findings[i].Reason = finding.Reason
		}
	}
	return findings
}
