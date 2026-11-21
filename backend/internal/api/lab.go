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

type labInput struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
	Tests     []struct {
		Input          string `json:"input"`
		ExpectedOutput string `json:"expectedOutput"`
	} `json:"tests"`
	Constraints []struct {
		Description string   `json:"description"`
		Deduction   *float64 `json:"deduction"`
	} `json:"constraints"`
}

func validLab(in labInput) bool {
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.Statement) == "" || len(in.Tests) == 0 || len(in.Constraints) == 0 {
		return false
	}
	for _, c := range in.Constraints {
		if strings.TrimSpace(c.Description) == "" || c.Deduction == nil || *c.Deduction < 0 {
			return false
		}
	}
	return true
}

func (s *Server) createLab(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := s.Sessions.GetString(ctx, "user_id")
	if user == "" {
		fail(w, 401, "unauthenticated", "Authentication required", "A valid session is required.")
		return
	}
	var role string
	if err := s.DB.QueryRow(ctx, "SELECT role FROM users WHERE id=$1", user).Scan(&role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fail(w, 401, "unauthenticated", "Authentication required", "A valid session is required.")
		} else {
			internal(w)
		}
		return
	}
	if role != "ta" {
		fail(w, 403, "forbidden", "Forbidden", "Only TAs may create labs.")
		return
	}
	if !s.validOrigin(r) {
		fail(w, 403, "invalid-origin", "Invalid origin", "The request origin is not allowed.")
		return
	}
	classID, err := uuid.Parse(chi.URLParam(r, "classId"))
	if err != nil {
		fail(w, 404, "not-found", "Class not found", "The class was not found.")
		return
	}
	var member bool
	if err = s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM class_tas WHERE class_id=$1 AND ta_id=$2)", classID, user).Scan(&member); err != nil {
		internal(w)
		return
	}
	if !member {
		fail(w, 404, "not-found", "Class not found", "The class was not found.")
		return
	}
	var in labInput
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024*1024))
	dec.DisallowUnknownFields()
	if dec.Decode(&in) != nil || dec.Decode(&struct{}{}) != io.EOF || !validLab(in) {
		fail(w, 400, "invalid-request", "Invalid request", "Provide a title, statement, tests, and valid constraints.")
		return
	}
	// Keep the lab, test cases, and deductions together: a failure saves none of them.
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		internal(w)
		return
	}
	defer tx.Rollback(ctx)
	id := uuid.New()
	if _, err = tx.Exec(ctx, "INSERT INTO labs (id,title,class_id,statement) VALUES ($1,$2,$3,$4)", id, in.Title, classID, in.Statement); err != nil {
		internal(w)
		return
	}
	for _, test := range in.Tests {
		if _, err = tx.Exec(ctx, "INSERT INTO lab_tests (id,lab_id,input,expected_output) VALUES ($1,$2,$3,$4)", uuid.New(), id, test.Input, test.ExpectedOutput); err != nil {
			internal(w)
			return
		}
	}
	for _, c := range in.Constraints {
		if _, err = tx.Exec(ctx, "INSERT INTO lab_constraints (id,lab_id,description,deduction) VALUES ($1,$2,$3,$4)", uuid.New(), id, c.Description, *c.Deduction); err != nil {
			internal(w)
			return
		}
	}
	// Students assigned to the class can submit to every newly created lab.
	if _, err = tx.Exec(ctx, "INSERT INTO enrollments(lab_id,student_id) SELECT $1,student_id FROM class_students WHERE class_id=$2", id, classID); err != nil {
		internal(w)
		return
	}
	if err = tx.Commit(ctx); err != nil {
		internal(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id.String(), "classId": classID.String()})
}
