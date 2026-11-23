-- +goose Up
-- A model finding has no grade effect. A named TA makes a separate decision.
CREATE TABLE constraint_reviews (
    submission_id uuid NOT NULL,
    constraint_id uuid NOT NULL,
    ta_id uuid NOT NULL REFERENCES users(id),
    decision text NOT NULL CHECK (decision IN ('pass', 'violation')),
    applied_deduction numeric NOT NULL CHECK (applied_deduction >= 0),
    evidence jsonb NOT NULL DEFAULT '[]',
    reviewed_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (submission_id, constraint_id),
    FOREIGN KEY (submission_id, constraint_id)
        REFERENCES constraint_verdicts(submission_id, constraint_id)
);
CREATE TABLE constraint_review_history (
    id uuid PRIMARY KEY,
    submission_id uuid NOT NULL,
    constraint_id uuid NOT NULL,
    ta_id uuid NOT NULL REFERENCES users(id),
    decision text NOT NULL CHECK (decision IN ('pass', 'violation')),
    applied_deduction numeric NOT NULL CHECK (applied_deduction >= 0),
    evidence jsonb NOT NULL,
    reviewed_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (submission_id, constraint_id)
        REFERENCES constraint_verdicts(submission_id, constraint_id)
);

-- +goose Down
DROP TABLE constraint_review_history;
DROP TABLE constraint_reviews;
