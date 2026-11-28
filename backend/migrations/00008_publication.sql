-- +goose Up
ALTER TABLE labs ADD COLUMN review_revision bigint NOT NULL DEFAULT 1;
ALTER TABLE constraint_reviews ADD COLUMN review_revision bigint NOT NULL DEFAULT 0;
ALTER TABLE constraint_review_history ADD COLUMN review_revision bigint NOT NULL DEFAULT 0;

CREATE TABLE review_populations (
    lab_id uuid NOT NULL REFERENCES labs(id),
    revision bigint NOT NULL,
    PRIMARY KEY (lab_id, revision)
);
CREATE TABLE review_pass_samples (
    lab_id uuid NOT NULL,
    revision bigint NOT NULL,
    submission_id uuid NOT NULL,
    constraint_id uuid NOT NULL,
    PRIMARY KEY (lab_id, revision, submission_id, constraint_id),
    FOREIGN KEY (lab_id, revision) REFERENCES review_populations(lab_id, revision)
);
CREATE TABLE lab_publications (
    lab_id uuid NOT NULL REFERENCES labs(id),
    revision bigint NOT NULL,
    published_by uuid NOT NULL REFERENCES users(id),
    published_at timestamptz NOT NULL DEFAULT now(),
    active boolean NOT NULL DEFAULT true,
    PRIMARY KEY (lab_id, revision)
);
CREATE TABLE published_results (
    lab_id uuid NOT NULL,
    revision bigint NOT NULL,
    submission_id uuid NOT NULL REFERENCES submissions(id),
    deduction numeric NOT NULL CHECK (deduction >= 0),
    PRIMARY KEY (lab_id, revision, submission_id),
    FOREIGN KEY (lab_id, revision) REFERENCES lab_publications(lab_id, revision)
);

-- Lock BEFORE changing population data. This also covers direct SQL/rejudging,
-- rather than relying on all callers remembering an application lock.
-- +goose StatementBegin
CREATE FUNCTION invalidate_review_population() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    old_lab uuid;
    new_lab uuid;
BEGIN
    IF TG_TABLE_NAME IN ('submissions', 'lab_constraints') THEN
        IF TG_OP <> 'INSERT' THEN old_lab := OLD.lab_id; END IF;
        IF TG_OP <> 'DELETE' THEN new_lab := NEW.lab_id; END IF;
    ELSE
        IF TG_OP <> 'INSERT' THEN
            SELECT lab_id INTO old_lab FROM submissions WHERE id=OLD.submission_id;
        END IF;
        IF TG_OP <> 'DELETE' THEN
            SELECT lab_id INTO new_lab FROM submissions WHERE id=NEW.submission_id;
        END IF;
    END IF;
    -- Ordered locking also handles a row being moved between labs.
    PERFORM id FROM labs WHERE id IN (old_lab,new_lab) ORDER BY id FOR UPDATE;
    UPDATE labs SET review_revision=review_revision+1 WHERE id IN (old_lab,new_lab);
    RETURN COALESCE(NEW,OLD);
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER submissions_review_population BEFORE INSERT OR UPDATE OR DELETE ON submissions
    FOR EACH ROW EXECUTE FUNCTION invalidate_review_population();
CREATE TRIGGER constraints_review_population BEFORE INSERT OR UPDATE OR DELETE ON lab_constraints
    FOR EACH ROW EXECUTE FUNCTION invalidate_review_population();
CREATE TRIGGER judges_review_population BEFORE INSERT OR UPDATE OR DELETE ON judge_results
    FOR EACH ROW EXECUTE FUNCTION invalidate_review_population();
CREATE TRIGGER findings_review_population BEFORE INSERT OR UPDATE OR DELETE ON constraint_verdicts
    FOR EACH ROW EXECUTE FUNCTION invalidate_review_population();
CREATE TRIGGER jobs_review_population BEFORE INSERT OR DELETE OR UPDATE OF status ON jobs
    FOR EACH ROW EXECUTE FUNCTION invalidate_review_population();
-- +goose StatementBegin
CREATE FUNCTION invalidate_lab_statement() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.statement IS DISTINCT FROM OLD.statement THEN
        NEW.review_revision := OLD.review_revision+1;
    END IF;
    RETURN NEW;
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER statement_review_population BEFORE UPDATE OF statement ON labs
    FOR EACH ROW EXECUTE FUNCTION invalidate_lab_statement();

-- +goose Down
DROP TRIGGER statement_review_population ON labs;
DROP FUNCTION invalidate_lab_statement();
DROP TRIGGER jobs_review_population ON jobs;
DROP TRIGGER findings_review_population ON constraint_verdicts;
DROP TRIGGER judges_review_population ON judge_results;
DROP TRIGGER constraints_review_population ON lab_constraints;
DROP TRIGGER submissions_review_population ON submissions;
DROP FUNCTION invalidate_review_population();
DROP TABLE published_results,lab_publications,review_pass_samples,review_populations;
ALTER TABLE constraint_review_history DROP COLUMN review_revision;
ALTER TABLE constraint_reviews DROP COLUMN review_revision;
ALTER TABLE labs DROP COLUMN review_revision;
