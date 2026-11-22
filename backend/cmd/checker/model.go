package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"leetkob/backend/internal/checker"
	"leetkob/backend/internal/semantic"
)

type runModel func(context.Context, semantic.Request) ([]byte, error)

func (m runModel) Generate(ctx context.Context, request semantic.Request) ([]byte, error) {
	return m(ctx, request)
}

// Each evaluation run owns a new audit file and a per-problem run budget.
// Production lab budgets instead use the worker's persistent database ledger.
func configureModel(engine *checker.Checker, path string) (func(), error) {
	mode := os.Getenv("MODEL_OFF")
	if mode == "" || mode == "1" {
		return func() {}, nil
	}
	if mode != "0" {
		return nil, fmt.Errorf("MODEL_OFF must be 0 or 1")
	}
	name, key := os.Getenv("MODEL_NAME"), os.Getenv("MODEL_API_KEY")
	if name != semantic.SupportedModel || key == "" || path == "" {
		return nil, fmt.Errorf("enabled checker requires supported MODEL_NAME, MODEL_API_KEY and --call-log")
	}

	budget := int64(1000000)
	if value := os.Getenv("MODEL_LAB_BUDGET_MICRO_USD"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < 0 {
			return nil, fmt.Errorf("invalid MODEL_LAB_BUDGET_MICRO_USD")
		}
		budget = parsed
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, err
	}

	// Sync reservations before calling the provider; a partial run still records
	// uncertain spending. Existing audit files cannot be silently overwritten.
	appendLog := func(record any) error {
		if err := json.NewEncoder(file).Encode(record); err != nil {
			return err
		}
		return file.Sync()
	}
	spent := map[string]int64{}
	engine.ModelFor = func(problem, solution string) semantic.Checker {
		return semantic.Checker{Enabled: true, Model: runModel(
			func(ctx context.Context, request semantic.Request) ([]byte, error) {
				if budget < semantic.ReserveMicroUSD || spent[problem] > budget-semantic.ReserveMicroUSD {
					return nil, fmt.Errorf("evaluation problem budget exhausted")
				}
				id := uuid.NewString()
				if err := appendLog(map[string]any{
					"event": "reserved", "callId": id, "problem": problem, "solution": solution,
					"reservedMicroUsd": semantic.ReserveMicroUSD,
				}); err != nil {
					return nil, err
				}
				spent[problem] += semantic.ReserveMicroUSD

				model, err := semantic.NewGemini(name, key, func(entry semantic.CallLog) error {
					if err := appendLog(map[string]any{
						"event": "recorded", "callId": id, "problem": problem,
						"solution": solution, "call": entry,
					}); err != nil {
						return err
					}
					if entry.UsageKnown {
						spent[problem] -= semantic.ReserveMicroUSD - entry.CostMicroUSD
					}
					return nil
				})
				if err != nil {
					return nil, err
				}
				return model.Generate(ctx, request)
			},
		)}
	}

	return func() { _ = file.Close() }, nil
}
