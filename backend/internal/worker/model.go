package worker

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/semantic"
)

// ModelConfig is worker-only configuration, never supplied to the sandbox.
type ModelConfig struct {
	Enabled        bool
	Name           string
	Key            string
	BudgetMicroUSD int64
}

// Checker binds a call to its lab budget without putting IDs in the prompt.
func (config ModelConfig) Checker(db *pgxpool.Pool, lab, submission string) semantic.Checker {
	return semantic.Checker{
		Enabled: config.Enabled,
		Model: &budgetedModel{
			db: db, config: config, lab: lab, submission: submission,
		},
	}
}

type budgetedModel struct {
	db              *pgxpool.Pool
	config          ModelConfig
	lab, submission string
	// Tests inject an offline provider behind the actual persistent budget gate.
	newModel func(func(semantic.CallLog) error) (semantic.Model, error)
}

func (m *budgetedModel) Generate(ctx context.Context, request semantic.Request) ([]byte, error) {
	// Validate the provider configuration before reserving any budget.
	callID := uuid.NewString()
	factory := m.newModel
	if factory == nil {
		factory = func(record func(semantic.CallLog) error) (semantic.Model, error) {
			return semantic.NewGemini(m.config.Name, m.config.Key, record)
		}
	}
	model, err := factory(func(entry semantic.CallLog) error {
		// Audit writes need their own short deadline after a provider timeout.
		auditCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return recordCall(auditCtx, m.db, callID, entry)
	})
	if err != nil {
		return nil, err
	}

	if err := reserveCall(ctx, m.db, callID, m.lab, m.submission, m.config.BudgetMicroUSD); err != nil {
		return nil, err
	}

	return model.Generate(ctx, request)
}

// Lock the lab row before summing reservations. All workers use this lock, so
// simultaneous submissions cannot both spend the same remaining budget.
func reserveCall(
	ctx context.Context,
	db *pgxpool.Pool,
	id, lab, submission string,
	budget int64,
) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var locked string
	if err := tx.QueryRow(ctx, "SELECT id::text FROM labs WHERE id=$1 FOR UPDATE", lab).Scan(&locked); err != nil {
		return err
	}

	var charged int64
	if err := tx.QueryRow(ctx,
		"SELECT COALESCE(sum(charged_micro_usd),0)::bigint FROM model_calls WHERE lab_id=$1", lab,
	).Scan(&charged); err != nil {
		return err
	}

	if budget < semantic.ReserveMicroUSD || charged > budget-semantic.ReserveMicroUSD {
		return semantic.ErrBudgetExceeded
	}

	_, err = tx.Exec(ctx, `INSERT INTO model_calls
 (id,lab_id,submission_id,charged_micro_usd) VALUES($1,$2,$3,$4)`,
		id, lab, submission, semantic.ReserveMicroUSD,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Unknown usage keeps the full reservation, including failed requests and
// crashed workers. Audit failure also leaves the reservation untouched.
func recordCall(ctx context.Context, db *pgxpool.Pool, id string, entry semantic.CallLog) error {
	charged := semantic.ReserveMicroUSD
	if entry.UsageKnown {
		charged = entry.CostMicroUSD
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, `UPDATE model_calls
 SET charged_micro_usd=$2,status='recorded',call_log=$3
 WHERE id=$1 AND status='reserved'`, id, charged, data)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("model audit reservation missing or already recorded")
	}

	return nil
}
