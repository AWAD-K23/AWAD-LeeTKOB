-- +goose Up
-- Reservations survive worker restarts and serialize spending for each lab.
CREATE TABLE model_calls (
    id uuid PRIMARY KEY,
    lab_id uuid NOT NULL REFERENCES labs(id),
    submission_id uuid NOT NULL REFERENCES submissions(id),
    charged_micro_usd bigint NOT NULL CHECK (charged_micro_usd >= 0),
    status text NOT NULL DEFAULT 'reserved'
        CHECK (status IN ('reserved', 'recorded')),
    call_log jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX model_calls_lab_idx ON model_calls(lab_id);

-- +goose Down
DROP TABLE model_calls;
