package api

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const countedSubmissions = `SELECT DISTINCT ON (student_id) id FROM submissions
 WHERE lab_id=$1 ORDER BY student_id,created_at DESC,id DESC`

// Partial Fisher-Yates: uniform sampling without replacement, with an injectable
// source for reproducible tests. Production uses the OS random source.
func samplePasses(n int, random func(int) (int, error)) ([]int, error) {
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	k := (n + 9) / 10
	for i := 0; i < k; i++ {
		j, err := random(n - i)
		if err != nil {
			return nil, err
		}
		if j < 0 || j >= n-i {
			return nil, errors.New("invalid random index")
		}
		indices[i], indices[i+j] = indices[i+j], indices[i]
	}
	return indices[:k], nil
}

// All readers/writers of review state lock the lab first. Population mutations
// take this same lock in migration 00008, including worker and direct SQL writes.
func (s *Server) reviewPopulation(ctx context.Context, lab uuid.UUID) (pgx.Tx, int64, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	var revision int64
	if err = tx.QueryRow(ctx, `SELECT review_revision FROM labs WHERE id=$1 FOR UPDATE`, lab).Scan(&revision); err != nil {
		tx.Rollback(ctx)
		return nil, 0, err
	}
	if err = s.ensureSample(ctx, tx, lab, revision); err != nil {
		tx.Rollback(ctx)
		return nil, 0, err
	}
	return tx, revision, nil
}

func (s *Server) ensureSample(ctx context.Context, tx pgx.Tx, lab uuid.UUID, revision int64) error {
	result, err := tx.Exec(ctx, `INSERT INTO review_populations(lab_id,revision) VALUES($1,$2) ON CONFLICT DO NOTHING`, lab, revision)
	if err != nil || result.RowsAffected() == 0 {
		return err
	}
	rows, err := tx.Query(ctx, `SELECT v.submission_id::text,v.constraint_id::text
 FROM constraint_verdicts v JOIN submissions s ON s.id=v.submission_id
 JOIN judge_results j ON j.submission_id=s.id AND j.verdict='Accepted'
 JOIN lab_constraints c ON c.id=v.constraint_id AND c.lab_id=s.lab_id
 WHERE s.id IN (`+countedSubmissions+`) AND v.verdict='pass'
 ORDER BY v.submission_id,v.constraint_id`, lab)
	if err != nil {
		return err
	}
	var passes [][2]string
	for rows.Next() {
		var pair [2]string
		if err = rows.Scan(&pair[0], &pair[1]); err != nil {
			rows.Close()
			return err
		}
		passes = append(passes, pair)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	random := s.ReviewRandom
	if random == nil {
		random = func(n int) (int, error) {
			v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
			if err != nil {
				return 0, err
			}
			return int(v.Int64()), nil
		}
	}
	selected, err := samplePasses(len(passes), random)
	if err != nil {
		return err
	}
	for _, i := range selected {
		if _, err = tx.Exec(ctx, `INSERT INTO review_pass_samples(lab_id,revision,submission_id,constraint_id) VALUES($1,$2,$3,$4)`, lab, revision, passes[i][0], passes[i][1]); err != nil {
			return err
		}
	}
	return nil
}

type publicationProgress struct {
	Revision   int64 `json:"reviewRevision"`
	Required   int   `json:"required"`
	Reviewed   int   `json:"reviewed"`
	Pending    int   `json:"pending"`
	Incomplete int   `json:"incompleteSubmissions"`
	Published  bool  `json:"published"`
}

func progress(ctx context.Context, tx pgx.Tx, lab uuid.UUID, revision int64) (publicationProgress, error) {
	p := publicationProgress{Revision: revision}
	err := tx.QueryRow(ctx, `SELECT count(*),count(r.submission_id)
 FROM constraint_verdicts v JOIN submissions s ON s.id=v.submission_id
 JOIN judge_results j ON j.submission_id=s.id AND j.verdict='Accepted'
 JOIN lab_constraints c ON c.id=v.constraint_id AND c.lab_id=s.lab_id
 LEFT JOIN review_pass_samples sample ON sample.lab_id=$1 AND sample.revision=$2
 AND sample.submission_id=s.id AND sample.constraint_id=c.id
 LEFT JOIN constraint_reviews r ON r.submission_id=s.id AND r.constraint_id=c.id
 AND r.review_revision=$2 AND r.applied_deduction=CASE WHEN r.decision='violation' THEN c.deduction ELSE 0 END
 WHERE s.id IN (`+countedSubmissions+`) AND (v.verdict IN ('violation','unsure') OR sample.submission_id IS NOT NULL)`, lab, revision).Scan(&p.Required, &p.Reviewed)
	if err != nil {
		return p, err
	}
	p.Pending = p.Required - p.Reviewed
	err = tx.QueryRow(ctx, `SELECT count(*) FROM submissions s
 LEFT JOIN jobs job ON job.submission_id=s.id LEFT JOIN judge_results j ON j.submission_id=s.id
 WHERE s.id IN (`+countedSubmissions+`) AND (job.status IS DISTINCT FROM 'completed' OR j.submission_id IS NULL
 OR (j.verdict='Accepted' AND EXISTS(SELECT 1 FROM lab_constraints c WHERE c.lab_id=$1
 AND NOT EXISTS(SELECT 1 FROM constraint_verdicts v WHERE v.submission_id=s.id AND v.constraint_id=c.id))))`, lab).Scan(&p.Incomplete)
	if err != nil {
		return p, err
	}
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM lab_publications WHERE lab_id=$1 AND revision=$2 AND active)`, lab, revision).Scan(&p.Published)
	return p, err
}

func (s *Server) publicationStatus(w http.ResponseWriter, r *http.Request) {
	_, lab, ok := s.assignedTA(w, r)
	if !ok {
		return
	}
	tx, revision, err := s.reviewPopulation(r.Context(), lab)
	if err != nil {
		internal(w)
		return
	}
	defer tx.Rollback(r.Context())
	p, err := progress(r.Context(), tx, lab, revision)
	if err != nil {
		internal(w)
		return
	}
	if err = tx.Commit(r.Context()); err != nil {
		internal(w)
		return
	}
	jsonResponse(w, p)
}

func (s *Server) publish(w http.ResponseWriter, r *http.Request) {
	id, lab, ok := s.assignedTA(w, r)
	if !ok {
		return
	}
	if !s.validOrigin(r) {
		fail(w, 403, "invalid-origin", "Invalid origin", "The request origin is not allowed.")
		return
	}
	var in struct {
		Revision int64 `json:"reviewRevision"`
	}
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024))
	d.DisallowUnknownFields()
	if d.Decode(&in) != nil || d.Decode(&struct{}{}) != io.EOF || in.Revision < 1 {
		fail(w, 400, "invalid-request", "Invalid request", "Provide the review revision.")
		return
	}
	tx, revision, err := s.reviewPopulation(r.Context(), lab)
	if err != nil {
		internal(w)
		return
	}
	defer tx.Rollback(r.Context())
	if in.Revision != revision {
		fail(w, 409, "stale-review-population", "Review population changed", "Reload and review the current population.")
		return
	}
	p, err := progress(r.Context(), tx, lab, revision)
	if err != nil {
		internal(w)
		return
	}
	if p.Pending > 0 || p.Incomplete > 0 {
		fail(w, 409, "pending-reviews", "Reviews pending", "Complete judging and all required TA reviews before publishing.")
		return
	}
	_, err = tx.Exec(r.Context(), `INSERT INTO lab_publications(lab_id,revision,published_by) VALUES($1,$2,$3)
 ON CONFLICT(lab_id,revision) DO UPDATE SET active=true,published_by=excluded.published_by,published_at=now()`, lab, revision, id)
	if err != nil {
		internal(w)
		return
	}
	_, err = tx.Exec(r.Context(), `INSERT INTO published_results(lab_id,revision,submission_id,deduction)
 SELECT $1,$2,s.id,COALESCE((SELECT sum(c.deduction) FROM constraint_reviews r
 JOIN lab_constraints c ON c.id=r.constraint_id AND c.lab_id=$1
 WHERE r.submission_id=s.id AND r.review_revision=$2 AND r.decision='violation'
 AND r.applied_deduction=c.deduction),0)
 FROM submissions s WHERE s.id IN (`+countedSubmissions+`)
 ON CONFLICT(lab_id,revision,submission_id) DO UPDATE SET deduction=excluded.deduction`, lab, revision)
	if err != nil {
		internal(w)
		return
	}
	if err = tx.Commit(r.Context()); err != nil {
		internal(w)
		return
	}
	p.Published = true
	jsonResponse(w, p)
}

func (s *Server) publishedResult(w http.ResponseWriter, r *http.Request) {
	id, role, ok := s.identity(w, r)
	if !ok {
		return
	}
	submission, err := uuid.Parse(chi.URLParam(r, "submissionId"))
	if err != nil {
		fail(w, 404, "not-found", "Result not found", "The published result was not found.")
		return
	}
	var revision int64
	var deduction float64
	var at, by string
	err = s.DB.QueryRow(r.Context(), `SELECT p.revision,p.deduction::float8,pub.published_at::text,pub.published_by::text
 FROM published_results p JOIN lab_publications pub ON pub.lab_id=p.lab_id AND pub.revision=p.revision AND pub.active
 JOIN labs l ON l.id=p.lab_id AND l.review_revision=p.revision
 JOIN submissions s ON s.id=p.submission_id
 WHERE s.id=$1 AND (($3='student' AND s.student_id=$2) OR ($3='ta' AND EXISTS
 (SELECT 1 FROM class_tas t WHERE t.class_id=l.class_id AND t.ta_id=$2)))`, submission, id, role).Scan(&revision, &deduction, &at, &by)
	if errors.Is(err, pgx.ErrNoRows) {
		fail(w, 404, "not-found", "Result not found", "The published result was not found.")
		return
	}
	if err != nil {
		internal(w)
		return
	}
	jsonResponse(w, map[string]any{"submissionId": submission.String(), "reviewRevision": revision, "deduction": deduction, "publishedAt": at, "publishedBy": by})
}
