package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// people exposes account selectors only to administrators, never password hashes.
func (s *Server) people(w http.ResponseWriter, r *http.Request) {
	_, role, ok := s.identity(w, r)
	if !ok {
		return
	}

	if role != "admin" {
		fail(w, 403, "forbidden", "Forbidden", "Only administrators may browse accounts.")
		return
	}

	rows, err := s.DB.Query(r.Context(), `
		SELECT id::text, email, role FROM users
		WHERE email ILIKE '%' || $1 || '%'
		ORDER BY email LIMIT 200`, r.URL.Query().Get("q"))
	if err != nil {
		internal(w)
		return
	}
	defer rows.Close()

	out := []map[string]string{}
	for rows.Next() {
		var id, email, role string
		if rows.Scan(&id, &email, &role) != nil {
			internal(w)
			return
		}

		out = append(out, map[string]string{"id": id, "email": email, "role": role})
	}

	if rows.Err() != nil {
		internal(w)
		return
	}

	jsonResponse(w, out)
}

// classDetail checks membership before returning class members or lab links.
func (s *Server) classDetail(w http.ResponseWriter, r *http.Request) {
	id, role, ok := s.identity(w, r)
	if !ok {
		return
	}

	classID, err := uuid.Parse(chi.URLParam(r, "classId"))
	if err != nil {
		fail(w, 404, "not-found", "Class not found", "The class was not found.")
		return
	}

	var name string
	err = s.DB.QueryRow(r.Context(), `
		SELECT name FROM classes c WHERE c.id=$1 AND (
		 $3='admin' OR EXISTS(SELECT 1 FROM class_tas WHERE class_id=c.id AND ta_id=$2)
		 OR EXISTS(SELECT 1 FROM class_students WHERE class_id=c.id AND student_id=$2))`,
		classID, id, role).Scan(&name)
	if errors.Is(err, pgx.ErrNoRows) {
		fail(w, 404, "not-found", "Class not found", "The class was not found.")
		return
	}

	if err != nil {
		internal(w)
		return
	}

	// Students do not receive their classmates' account details.
	members := []map[string]string{}
	if role != "student" {
		rows, err := s.DB.Query(r.Context(), `
			SELECT u.id::text,u.email,u.role FROM users u WHERE
			EXISTS(SELECT 1 FROM class_tas WHERE class_id=$1 AND ta_id=u.id)
			OR EXISTS(SELECT 1 FROM class_students WHERE class_id=$1 AND student_id=u.id)
			ORDER BY u.role,u.email`, classID)
		if err != nil {
			internal(w)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var memberID, email, memberRole string
			if rows.Scan(&memberID, &email, &memberRole) != nil {
				internal(w)
				return
			}

			members = append(members, map[string]string{
				"id": memberID, "email": email, "role": memberRole,
			})
		}

		if rows.Err() != nil {
			internal(w)
			return
		}
	}

	jsonResponse(w, map[string]any{"id": classID, "name": name, "members": members})
}

// submissionHistory returns the signed-in student's own submissions only.
func (s *Server) submissionHistory(w http.ResponseWriter, r *http.Request) {
	id, role, ok := s.identity(w, r)
	if !ok {
		return
	}

	if role != "student" {
		fail(w, 403, "forbidden", "Forbidden", "Only students may view their submission history.")
		return
	}

	rows, err := s.DB.Query(r.Context(), `
		SELECT s.id::text,l.id::text,l.title,s.created_at::text,j.status,
		COALESCE(x.verdict,'') FROM submissions s
		JOIN labs l ON l.id=s.lab_id JOIN jobs j ON j.submission_id=s.id
		LEFT JOIN judge_results x ON x.submission_id=s.id
		WHERE s.student_id=$1 ORDER BY s.created_at DESC LIMIT 200`, id)
	if err != nil {
		internal(w)
		return
	}
	defer rows.Close()

	out := []map[string]string{}
	for rows.Next() {
		var submissionID, labID, title, createdAt, status, verdict string
		if rows.Scan(&submissionID, &labID, &title, &createdAt, &status, &verdict) != nil {
			internal(w)
			return
		}

		out = append(out, map[string]string{
			"id": submissionID, "labId": labID, "title": title,
			"createdAt": createdAt, "status": status, "verdict": verdict,
		})
	}

	if rows.Err() != nil {
		internal(w)
		return
	}

	jsonResponse(w, out)
}
