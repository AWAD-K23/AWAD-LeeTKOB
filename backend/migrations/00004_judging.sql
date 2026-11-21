-- +goose Up
CREATE TABLE class_students (
    class_id uuid NOT NULL REFERENCES classes(id),
    student_id uuid NOT NULL REFERENCES users(id),
    PRIMARY KEY (class_id, student_id)
);

ALTER TABLE jobs ADD COLUMN lease_token uuid;
ALTER TABLE jobs ADD COLUMN leased_until timestamptz;
ALTER TABLE jobs ADD COLUMN attempts integer NOT NULL DEFAULT 0;
CREATE INDEX jobs_pending_idx ON jobs(status, created_at);

CREATE TABLE judge_results (
    submission_id uuid PRIMARY KEY REFERENCES submissions(id),
    verdict text NOT NULL,
    detail text NOT NULL,
    judged_at timestamptz NOT NULL DEFAULT now()
);

-- Model verdicts never contain the TA-selected deduction or a published mark.
CREATE TABLE constraint_verdicts (
    submission_id uuid NOT NULL REFERENCES submissions(id),
    constraint_id uuid NOT NULL REFERENCES lab_constraints(id),
    verdict text NOT NULL CHECK (verdict IN ('pass', 'violation', 'unsure')),
    evidence jsonb NOT NULL DEFAULT '[]',
    source text NOT NULL CHECK (source IN ('hard_rule', 'semantic_pending', 'model')),
    PRIMARY KEY (submission_id, constraint_id)
);

-- +goose Down
DROP TABLE constraint_verdicts, judge_results;
DROP INDEX jobs_pending_idx;
ALTER TABLE jobs DROP COLUMN attempts;
ALTER TABLE jobs DROP COLUMN leased_until;
ALTER TABLE jobs DROP COLUMN lease_token;
DROP TABLE class_students;
