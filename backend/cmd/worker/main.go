package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/localenv"
	"leetkob/backend/internal/sandbox"
	"leetkob/backend/internal/semantic"
	"leetkob/backend/internal/worker"
)

// main runs judging separately from the HTTP server.
func main() {
	once := flag.Bool("once", false, "process at most one job and exit")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := localenv.Load(); err != nil {
		log.Fatal(err)
	}

	modelConfig, err := readModelConfig()
	if err != nil {
		log.Fatal(err)
	}

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
		SemanticFor: func(lab, submission string) semantic.Checker {
			return modelConfig.Checker(db, lab, submission)
		},
	}
	log.Print("Judging worker started")

	// Keep checking for jobs until shutdown is requested.
	for ctx.Err() == nil {
		found, err := w.ProcessOne(ctx)

		if err != nil && ctx.Err() == nil {
			log.Printf("Judging job: %v", err)
		}

		// A bounded run supports handoff smoke checks without leaving a worker
		// polling or accidentally processing additional jobs afterward.
		if *once {
			if err != nil {
				log.Fatal(err)
			}
			return
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

// Calls remain off unless explicitly enabled. The local per-lab limit is
// cumulative across all submissions and persists in the database ledger.
func readModelConfig() (worker.ModelConfig, error) {
	config := worker.ModelConfig{
		Name: os.Getenv("MODEL_NAME"), Key: os.Getenv("MODEL_API_KEY"),
		BudgetMicroUSD: 1000000, // Default local cap: US$1 per lab.
	}

	switch os.Getenv("MODEL_OFF") {
	case "", "1":
		return config, nil
	case "0":
		config.Enabled = true
	default:
		return config, fmt.Errorf("MODEL_OFF must be 0 or 1")
	}

	if value := os.Getenv("MODEL_LAB_BUDGET_MICRO_USD"); value != "" {
		budget, err := strconv.ParseInt(value, 10, 64)
		if err != nil || budget < 0 {
			return config, fmt.Errorf("MODEL_LAB_BUDGET_MICRO_USD must be a nonnegative integer")
		}
		config.BudgetMicroUSD = budget
	}

	if config.Name != semantic.SupportedModel || config.Key == "" {
		return config, fmt.Errorf("enabled model requires supported MODEL_NAME and MODEL_API_KEY")
	}

	return config, nil
}
