package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Read identity from the session and database, never from a request body.
func (s *Server) identity(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	id := s.Sessions.GetString(r.Context(), "user_id")
	var role string
	if id == "" {
		fail(w, 401, "unauthenticated", "Authentication required", "Sign in first.")
		return "", "", false
	}
	err := s.DB.QueryRow(r.Context(), "SELECT role FROM users WHERE id=$1", id).Scan(&role)
	if errors.Is(err, pgx.ErrNoRows) {
		fail(w, 401, "unauthenticated", "Authentication required", "Sign in first.")
		return "", "", false
	}
	if err != nil {
		internal(w)
		return "", "", false
	}
	return id, role, true
}

func jsonResponse(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	id, role, ok := s.identity(w, r)
	if !ok {
		return
	}
	jsonResponse(w, map[string]string{"id": id, "role": role})
}

func (s *Server) createClass(w http.ResponseWriter, r *http.Request) {
	_, role, ok := s.identity(w, r)
	if !ok {
		return
	}
	if role != "admin" {
		fail(w, 403, "forbidden", "Forbidden", "Only administrators may create classes.")
		return
	}
	if !s.validOrigin(r) {
		fail(w, 403, "invalid-origin", "Invalid origin", "The request origin is not allowed.")
		return
	}
	var in struct {
		Name       string   `json:"name"`
		TAID       string   `json:"taId"`
		StudentIDs []string `json:"studentIds"`
	}
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536))
	d.DisallowUnknownFields()
	if d.Decode(&in) != nil || d.Decode(&struct{}{}) != io.EOF || in.Name == "" {
		fail(w, 400, "invalid-request", "Invalid request", "Provide a name, TA ID and student IDs.")
		return
	}
	tx, err := s.DB.Begin(r.Context())
	if err != nil {
		internal(w)
		return
	}
	defer tx.Rollback(r.Context())
	var eligible bool
	if err = tx.QueryRow(r.Context(), "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1 AND role='ta')", in.TAID).Scan(&eligible); err != nil || !eligible {
		fail(w, 400, "invalid-request", "Invalid request", "Provide an existing TA.")
		return
	}
	id := uuid.NewString()
	if _, err = tx.Exec(r.Context(), "INSERT INTO classes(id,name) VALUES($1,$2)", id, in.Name); err != nil {
		internal(w)
		return
	}
	if _, err = tx.Exec(r.Context(), "INSERT INTO class_tas(class_id,ta_id) VALUES($1,$2)", id, in.TAID); err != nil {
		internal(w)
		return
	}
	for _, student := range in.StudentIDs {
		if err = tx.QueryRow(r.Context(), "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1 AND role='student')", student).Scan(&eligible); err != nil || !eligible {
			fail(w, 400, "invalid-request", "Invalid request", "Provide existing students.")
			return
		}
		if _, err = tx.Exec(r.Context(), "INSERT INTO class_students(class_id,student_id) VALUES($1,$2) ON CONFLICT DO NOTHING", id, student); err != nil {
			internal(w)
			return
		}
	}
	if err = tx.Commit(r.Context()); err != nil {
		internal(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	jsonResponse(w, map[string]string{"id": id})
}

func (s *Server) classes(w http.ResponseWriter, r *http.Request) {
	id, role, ok := s.identity(w, r)
	if !ok {
		return
	}
	rows, err := s.DB.Query(r.Context(), `SELECT c.id::text,c.name FROM classes c WHERE $2='admin'
 OR EXISTS(SELECT 1 FROM class_tas t WHERE t.class_id=c.id AND t.ta_id=$1)
 OR EXISTS(SELECT 1 FROM class_students t WHERE t.class_id=c.id AND t.student_id=$1) ORDER BY c.name`, id, role)
	if err != nil {
		internal(w)
		return
	}
	defer rows.Close()
	out := []map[string]string{}
	for rows.Next() {
		var id, name string
		if rows.Scan(&id, &name) != nil {
			internal(w)
			return
		}
		out = append(out, map[string]string{"id": id, "name": name})
	}
	if rows.Err() != nil {
		internal(w)
		return
	}
	jsonResponse(w, out)
}

func (s *Server) labs(w http.ResponseWriter, r *http.Request) {
	id, role, ok := s.identity(w, r)
	if !ok {
		return
	}
	rows, err := s.DB.Query(r.Context(), `SELECT l.id::text,l.title,l.statement,
 COALESCE(l.class_id::text,''),ARRAY(SELECT description FROM lab_constraints
 WHERE lab_id=l.id ORDER BY id) FROM labs l WHERE
 ($2='ta' AND EXISTS(SELECT 1 FROM class_tas t WHERE t.class_id=l.class_id AND t.ta_id=$1))
 OR ($2='student' AND EXISTS(SELECT 1 FROM enrollments e WHERE e.lab_id=l.id AND e.student_id=$1)) ORDER BY l.title`, id, role)
	if err != nil {
		internal(w)
		return
	}
	defer rows.Close()
	out := []map[string]any{}
	for rows.Next() {
		var id, title, statement, classID string
		var requirements []string
		if rows.Scan(&id, &title, &statement, &classID, &requirements) != nil {
			internal(w)
			return
		}
		out = append(out, map[string]any{
			"id": id, "title": title, "statement": statement, "classId": classID,
			"requirements": requirements,
		})
	}
	if rows.Err() != nil {
		internal(w)
		return
	}
	jsonResponse(w, out)
}

func (s *Server) submissionStatus(w http.ResponseWriter, r *http.Request) {
	id, role, ok := s.identity(w, r)
	if !ok {
		return
	}
	submissionID, err := uuid.Parse(chi.URLParam(r, "submissionId"))
	if err != nil {
		fail(w, 404, "not-found", "Submission not found", "The submission was not found.")
		return
	}
	var status, verdict, detail, labID string
	var source []byte
	// Students see judge results only; constraint verdicts are reserved for TA review.
	err = s.DB.QueryRow(r.Context(), `SELECT j.status,COALESCE(x.verdict,''),COALESCE(x.detail,''),s.source,l.id::text FROM submissions s
 JOIN labs l ON l.id=s.lab_id JOIN jobs j ON j.submission_id=s.id LEFT JOIN judge_results x ON x.submission_id=s.id
 WHERE s.id=$1 AND (($3='student' AND s.student_id=$2) OR ($3='ta' AND EXISTS(SELECT 1 FROM class_tas t WHERE t.class_id=l.class_id AND t.ta_id=$2)))`, submissionID, id, role).Scan(&status, &verdict, &detail, &source, &labID)
	if errors.Is(err, pgx.ErrNoRows) {
		fail(w, 404, "not-found", "Submission not found", "The submission was not found.")
		return
	}
	if err != nil {
		internal(w)
		return
	}
	jsonResponse(w, map[string]string{
		"id": submissionID.String(), "status": status, "verdict": verdict,
		"detail": detail, "source": string(source), "labId": labID,
	})
}
