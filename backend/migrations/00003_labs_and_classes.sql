-- +goose Up
-- Existing labs remain available while new labs are assigned to a class.
CREATE TABLE classes (
    id uuid PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE class_tas (
    class_id uuid NOT NULL REFERENCES classes(id),
    ta_id uuid NOT NULL REFERENCES users(id),
    PRIMARY KEY (class_id, ta_id)
);

ALTER TABLE labs ADD COLUMN class_id uuid REFERENCES classes(id);
ALTER TABLE labs ADD COLUMN statement text NOT NULL DEFAULT '';

CREATE TABLE lab_tests (
    id uuid PRIMARY KEY,
    lab_id uuid NOT NULL REFERENCES labs(id) ON DELETE CASCADE,
    input text NOT NULL,
    expected_output text NOT NULL
);

CREATE TABLE lab_constraints (
    id uuid PRIMARY KEY,
    lab_id uuid NOT NULL REFERENCES labs(id) ON DELETE CASCADE,
    description text NOT NULL,
    deduction numeric NOT NULL CHECK (deduction >= 0)
);

-- +goose Down
DROP TABLE lab_constraints;
DROP TABLE lab_tests;
ALTER TABLE labs DROP COLUMN statement;
ALTER TABLE labs DROP COLUMN class_id;
DROP TABLE class_tas;
DROP TABLE classes;
