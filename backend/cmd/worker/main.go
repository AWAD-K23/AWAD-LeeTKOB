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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is required")
	}
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	w := worker.Worker{DB: db, Runner: sandbox.Runner{Image: os.Getenv("SANDBOX_IMAGE")}}
	log.Print("Judging worker started")
	for ctx.Err() == nil {
		found, err := w.ProcessOne(ctx)
		if err != nil && ctx.Err() == nil {
			log.Printf("Judging job: %v", err)
		}
		if !found || err != nil {
			select {
			case <-ctx.Done():
			case <-time.After(time.Second):
			}
		}
	}
}
