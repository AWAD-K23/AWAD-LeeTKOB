// Starts the API, database connection, and login sessions.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/api"
)

func main() {
	// Load local settings before checking the required server configuration.
	if err := loadLocalEnv(); err != nil {
		log.Fatal(err)
	}

	dbURL, origin := os.Getenv("DATABASE_URL"), os.Getenv("APP_ORIGIN")
	if dbURL == "" || origin == "" {
		log.Fatal("DATABASE_URL and APP_ORIGIN are required")
	}

	// Keep a pool of database connections for API requests.
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Ping the database to ensure the connection is valid.
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Keep session data in PostgreSQL instead of server memory, so a server
	// restart does not automatically discard existing login sessions.
	store := pgxstore.New(pool)
	defer store.StopCleanup()
	sessions := scs.New()
	sessions.Store = store
	sessions.Lifetime = 24 * time.Hour

	// HttpOnly blocks browser JavaScript from reading the session cookie.
	// SameSite and Secure settings also control when the browser sends it.
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.SameSite = http.SameSiteLaxMode
	sessions.Cookie.Secure = os.Getenv("COOKIE_SECURE") != "false"
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	// Start the API server and listen for requests.
	server := &http.Server{
		Addr: ":" + port,
		Handler: (&api.Server{
			DB:       pool,
			Sessions: sessions,
			Origin:   origin,
		}).Routes(),
		ReadHeaderTimeout: 5 * time.Second}

	log.Printf("Listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}
