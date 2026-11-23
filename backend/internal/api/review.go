package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type reviewEvidence struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}

// assignedTA checks the database membership, rather than trusting a class ID
// supplied by the browser. An unknown lab and another class have the same 404.
func (s *Server) assignedTA(w http.ResponseWriter, r *http.Request) (string, uuid.UUID, bool) {
	id, role, ok := s.identity(w, r)
	if !ok {
		return "", uuid.UUID{}, false
	}
	if role != "ta" {
		fail(w, 403, "forbidden", "Forbidden", "Only assigned TAs may review.")
		return "", uuid.UUID{}, false
	}
	lab, err := uuid.Parse(chi.URLParam(r, "labId"))
	if err != nil {
		fail(w, 404, "not-found", "Lab not found", "The lab was not found.")
		return "", uuid.UUID{}, false
	}
	var allowed bool
	if err := s.DB.QueryRow(r.Context(), `SELECT EXISTS(
		SELECT 1 FROM labs l JOIN class_tas t ON t.class_id=l.class_id
		WHERE l.id=$1 AND t.ta_id=$2)`, lab, id).Scan(&allowed); err != nil {
		internal(w)
		return "", uuid.UUID{}, false
	}
	if !allowed {
		fail(w, 404, "not-found", "Lab not found", "The lab was not found.")
		return "", uuid.UUID{}, false
	}
	return id, lab, true
}

// reviewQueue includes only the latest submission for each student in a lab
// (FAQ-02). No route used by students reads this data (FAQ-08).
func (s *Server) reviewQueue(w http.ResponseWriter, r *http.Request) {
	_, lab, ok := s.assignedTA(w, r)
	if !ok {
		return
	}
	rows, err := s.DB.Query(r.Context(), `
		SELECT s.id::text,s.student_id::text,s.created_at::text,
		j.verdict,convert_from(s.source,'UTF8'),c.id::text,c.description,
		c.deduction::float8,v.verdict,v.reason,v.evidence,v.source,
		COALESCE(cr.decision,''),COALESCE(cr.evidence,'[]'::jsonb),
		COALESCE(cr.applied_deduction::float8,0)
		FROM submissions s
		JOIN judge_results j ON j.submission_id=s.id
		JOIN constraint_verdicts v ON v.submission_id=s.id
		JOIN lab_constraints c ON c.id=v.constraint_id AND c.lab_id=s.lab_id
		LEFT JOIN constraint_reviews cr ON cr.submission_id=s.id AND cr.constraint_id=c.id
		WHERE s.lab_id=$1 AND j.verdict='Accepted'
		AND NOT EXISTS (SELECT 1 FROM submissions newer WHERE newer.lab_id=s.lab_id
			AND newer.student_id=s.student_id AND
			(newer.created_at, newer.id) > (s.created_at, s.id))
		ORDER BY s.created_at DESC,s.id,c.id`, lab)
	if err != nil {
		internal(w)
		return
	}
	defer rows.Close()
	out := []map[string]any{}
	for rows.Next() {
		var submissionID, studentID, createdAt, judge, code, constraintID, description, verdict, reason, source, decision string
		var deduction, appliedDeduction float64
		var evidence, reviewedEvidence []byte
		if err := rows.Scan(&submissionID, &studentID, &createdAt, &judge, &code, &constraintID, &description, &deduction, &verdict, &reason, &evidence, &source, &decision, &reviewedEvidence, &appliedDeduction); err != nil {
			internal(w)
			return
		}
		out = append(out, map[string]any{
			"submissionId": submissionID, "studentId": studentID, "createdAt": createdAt,
			"judgeVerdict": judge, "sourceCode": code, "constraintId": constraintID,
			"constraint": description, "deduction": deduction, "verdict": verdict,
			"reason": reason, "evidence": json.RawMessage(evidence), "source": source, "decision": decision,
			"reviewEvidence":   json.RawMessage(reviewedEvidence),
			"appliedDeduction": appliedDeduction,
		})
	}
	if rows.Err() != nil {
		internal(w)
		return
	}
	jsonResponse(w, out)
}

func exactEvidence(code string, evidence []reviewEvidence) bool {
	if len(evidence) == 0 {
		return false
	}
	lines := strings.Split(code, "\n")
	for _, item := range evidence {
		if item.Line < 1 || item.Line > len(lines) || lines[item.Line-1] != item.Text {
			return false
		}
	}
	return true
}

// reviewDecision stores the TA's own verdict separately from the model finding.
// Deductions remain configured on lab_constraints; nothing here publishes a mark.
func (s *Server) reviewDecision(w http.ResponseWriter, r *http.Request) {
	id, lab, ok := s.assignedTA(w, r)
	if !ok {
		return
	}
	if !s.validOrigin(r) {
		fail(w, 403, "invalid-origin", "Invalid origin", "The request origin is not allowed.")
		return
	}
	submission, e1 := uuid.Parse(chi.URLParam(r, "submissionId"))
	constraint, e2 := uuid.Parse(chi.URLParam(r, "constraintId"))
	if e1 != nil || e2 != nil {
		fail(w, 404, "not-found", "Submission not found", "The submission was not found.")
		return
	}
	var in struct {
		Decision string           `json:"decision"`
		Evidence []reviewEvidence `json:"evidence"`
	}
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64*1024))
	d.DisallowUnknownFields()
	if d.Decode(&in) != nil || d.Decode(&struct{}{}) != io.EOF ||
		(in.Decision != "pass" && in.Decision != "violation") ||
		(in.Decision == "pass" && len(in.Evidence) != 0) {
		fail(w, 400, "invalid-request", "Invalid request", "Choose pass or violation and valid evidence.")
		return
	}
	// Lock the submission and validate that this is still the counted attempt.
	tx, err := s.DB.Begin(r.Context())
	if err != nil {
		internal(w)
		return
	}
	defer tx.Rollback(r.Context())
	var code []byte
	err = tx.QueryRow(r.Context(), `
		SELECT s.source FROM submissions s
		JOIN judge_results j ON j.submission_id=s.id AND j.verdict='Accepted'
		JOIN constraint_verdicts v ON v.submission_id=s.id AND v.constraint_id=$3
		JOIN lab_constraints c ON c.id=v.constraint_id AND c.lab_id=s.lab_id
		WHERE s.id=$1 AND s.lab_id=$2
		AND NOT EXISTS (SELECT 1 FROM submissions newer WHERE newer.lab_id=s.lab_id
			AND newer.student_id=s.student_id AND
			(newer.created_at,newer.id)>(s.created_at,s.id))
		FOR UPDATE OF s`, submission, lab, constraint).Scan(&code)
	if errors.Is(err, pgx.ErrNoRows) {
		fail(w, 404, "not-found", "Submission not found", "The submission was not found.")
		return
	}
	if err != nil {
		internal(w)
		return
	}
	if in.Decision == "violation" && !exactEvidence(string(code), in.Evidence) {
		fail(w, 400, "invalid-request", "Invalid request", "A violation needs exact submitted code lines.")
		return
	}
	evidence, err := json.Marshal(in.Evidence)
	if err != nil {
		internal(w)
		return
	}
	_, err = tx.Exec(r.Context(), `INSERT INTO constraint_reviews(submission_id,constraint_id,ta_id,decision,evidence,applied_deduction)
		SELECT $1,$2,$3,$4,$5,CASE WHEN $4='violation' THEN c.deduction ELSE 0 END
		FROM lab_constraints c WHERE c.id=$2 AND c.lab_id=$6
		ON CONFLICT(submission_id,constraint_id) DO UPDATE SET
		ta_id=excluded.ta_id,decision=excluded.decision,evidence=excluded.evidence,
		applied_deduction=excluded.applied_deduction,reviewed_at=now()`,
		submission, constraint, id, in.Decision, evidence, lab)
	if err != nil {
		internal(w)
		return
	}
	_, err = tx.Exec(r.Context(), `INSERT INTO constraint_review_history(id,submission_id,constraint_id,ta_id,decision,evidence,applied_deduction)
		SELECT $1,r.submission_id,r.constraint_id,r.ta_id,r.decision,r.evidence,r.applied_deduction
		FROM constraint_reviews r WHERE r.submission_id=$2 AND r.constraint_id=$3`, uuid.New(), submission, constraint)
	if err != nil {
		internal(w)
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		internal(w)
		return
	}
	jsonResponse(w, map[string]string{"decision": in.Decision, "reviewedBy": id})
}
