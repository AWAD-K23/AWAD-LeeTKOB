-- +goose Up
ALTER TABLE constraint_verdicts ADD COLUMN reason text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE constraint_verdicts DROP COLUMN reason;
