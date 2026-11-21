package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type Server struct {
	DB       *pgxpool.Pool
	Sessions *scs.SessionManager
	Origin   string
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/auth/login", s.login)
	r.Post("/auth/logout", s.logout)
	r.Post("/labs/{labId}/enrollments", s.enroll)
	r.Post("/labs/{labId}/submissions", s.submit)

	// Load the session before each request and save any changes afterwards.
	return s.Sessions.LoadAndSave(r)
}

func (s *Server) validOrigin(r *http.Request) bool {

	// Require exactly one Origin header matching the configured frontend.
	// State-changing handlers use this to reject requests from other origins.
	values := r.Header.Values("Origin")
	return s.Origin != "" && len(values) == 1 && values[0] == s.Origin
}
