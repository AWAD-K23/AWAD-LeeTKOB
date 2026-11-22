package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/sandbox"
	"leetkob/backend/internal/worker"
)

// main runs judging separately from the HTTP server.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Read connection settings before starting the job loop.
	url := os.Getenv("DATABASE_URL")

	if url == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := pgxpool.New(ctx, url)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Use the real sandbox runner for every claimed submission.
	w := worker.Worker{
		DB: db,
		Runner: sandbox.Runner{
			Image: os.Getenv("SANDBOX_IMAGE"),
		},
	}
	log.Print("Judging worker started")

	// Keep checking for jobs until shutdown is requested.
	for ctx.Err() == nil {
		found, err := w.ProcessOne(ctx)

		if err != nil && ctx.Err() == nil {
			log.Printf("Judging job: %v", err)
		}

		// Avoid repeatedly querying the database when idle or recovering.
		if !found || err != nil {
			select {
			case <-ctx.Done():
			case <-time.After(time.Second):
			}
		}
	}
}
