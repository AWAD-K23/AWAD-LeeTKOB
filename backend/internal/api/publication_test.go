package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPassSample(t *testing.T) {
	for _, n := range []int{0, 1, 9, 10, 11, 21, 100} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			rng := rand.New(rand.NewSource(42))
			pick := func(n int) (int, error) { return rng.Intn(n), nil }
			got, err := samplePasses(n, pick)
			if err != nil || len(got) != (n+9)/10 {
				t.Fatalf("sample %v: %v", got, err)
			}
			seen := map[int]bool{}
			for _, i := range got {
				if i < 0 || i >= n || seen[i] {
					t.Fatalf("invalid sample %v", got)
				}
				seen[i] = true
			}
			if n == 21 && !reflect.DeepEqual(got, []int{5, 8, 12}) {
				t.Fatalf("fixed seed sample changed: %v", got)
			}
		})
	}
}

// Each test owns a fresh schema; no persistent project data or real models.
func publicationDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	ctx := context.Background()
	base, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	schema := "publication_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err = base.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		base.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() { base.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE"); base.Close() })
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(db.Close)
	files, err := filepath.Glob("../../migrations/*.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = db.Exec(ctx, strings.Split(string(b), "-- +goose Down")[0]); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
	}
	return db
}

func TestPublicationGate(t *testing.T) {
	db := publicationDB(t)
	ctx := context.Background()
	exec := func(sql string, args ...any) {
		t.Helper()
		if _, err := db.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	student, other, ta, otherTA, class, otherClass, lab, submission := uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString()
	for i, user := range []string{student, other, ta, otherTA} {
		role := "student"
		if i > 1 {
			role = "ta"
		}
		exec(`INSERT INTO users(id,email,password_hash,role) VALUES($1,$2,'unused',$3)`, user, user+"@example.test", role)
	}
	exec(`INSERT INTO classes(id,name) VALUES($1,'Own'),($2,'Other')`, class, otherClass)
	exec(`INSERT INTO class_tas(class_id,ta_id) VALUES($1,$2),($3,$4)`, class, ta, otherClass, otherTA)
	exec(`INSERT INTO labs(id,title,class_id) VALUES($1,'Publication',$2)`, lab, class)
	exec(`INSERT INTO enrollments(lab_id,student_id) VALUES($1,$2)`, lab, student)
	exec(`INSERT INTO submissions(id,lab_id,student_id,source,created_at) VALUES($1,$2,$3,$4,now())`, submission, lab, student, []byte("int main() {}"))
	exec(`INSERT INTO jobs(id,submission_id,status) VALUES($1,$2,'completed')`, uuid.NewString(), submission)
	exec(`INSERT INTO judge_results(submission_id,verdict,detail) VALUES($1,'Accepted','OK')`, submission)
	for i := 0; i < 13; i++ {
		constraint := uuid.NewString()
		verdict := "pass"
		if i == 11 {
			verdict = "violation"
		}
		if i == 12 {
			verdict = "unsure"
		}
		exec(`INSERT INTO lab_constraints(id,lab_id,description,deduction) VALUES($1,$2,'Requirement',3)`, constraint, lab)
		exec(`INSERT INTO constraint_verdicts(submission_id,constraint_id,verdict,source) VALUES($1,$2,$3,'model')`, submission, constraint, verdict)
	}
	store := pgxstore.NewWithCleanupInterval(db, 0)
	t.Cleanup(store.StopCleanup)
	sessions := scs.New()
	sessions.Store = store
	rng := rand.New(rand.NewSource(42))
	s := &Server{DB: db, Sessions: sessions, Origin: "http://localhost:5173", ReviewRandom: func(n int) (int, error) { return rng.Intn(n), nil }}
	h := s.Routes()
	cookies := map[string]*http.Cookie{}
	for _, user := range []string{student, other, ta, otherTA} {
		loaded, err := sessions.Load(ctx, "")
		if err != nil {
			t.Fatal(err)
		}
		sessions.Put(loaded, "user_id", user)
		token, _, err := sessions.Commit(loaded)
		if err != nil {
			t.Fatal(err)
		}
		cookies[user] = &http.Cookie{Name: sessions.Cookie.Name, Value: token}
	}
	request := func(method, path, body, user string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Origin", s.Origin)
		r.Header.Set("Content-Type", "application/json")
		if cookies[user] != nil {
			r.AddCookie(cookies[user])
		}
		out := httptest.NewRecorder()
		h.ServeHTTP(out, r)
		return out
	}
	check := func(out *httptest.ResponseRecorder, code int) {
		t.Helper()
		if out.Code != code {
			t.Fatalf("wanted %d, got %d: %s", code, out.Code, out.Body.String())
		}
	}
	path := "/labs/" + lab
	getProgress := func() publicationProgress {
		t.Helper()
		out := request("GET", path+"/publication", "", ta)
		check(out, 200)
		var p publicationProgress
		if err := json.Unmarshal(out.Body.Bytes(), &p); err != nil {
			t.Fatal(err)
		}
		return p
	}
	p := getProgress()
	if p.Required != 4 || p.Pending != 4 || p.Incomplete != 0 {
		t.Fatalf("progress: %+v", p)
	}
	body := func(revision int64) string { return fmt.Sprintf(`{"reviewRevision":%d}`, revision) }
	for _, route := range []string{"/publication", "/review-queue"} {
		check(request("GET", path+route, "", student), 403)
		check(request("GET", path+route, "", otherTA), 404)
		check(request("GET", path+route, "", ""), 401)
	}
	check(request("POST", path+"/publish", body(p.Revision), student), 403)
	check(request("POST", path+"/publish", body(p.Revision), otherTA), 404)
	check(request("POST", path+"/publish", body(p.Revision), ta), 409)
	check(request("POST", path+"/publish", fmt.Sprintf(`{"reviewRevision":%d,"deduction":0}`, p.Revision), ta), 400)
	resultPath := "/submissions/" + submission + "/published-result"
	check(request("GET", resultPath, "", student), 404)
	first := request("GET", path+"/review-queue", "", ta)
	check(first, 200)
	second := request("GET", path+"/review-queue", "", ta)
	check(second, 200)
	if first.Body.String() != second.Body.String() {
		t.Fatal("queue reload redrew sample")
	}
	// A new server instance must load the persisted sample without using RNG.
	h = (&Server{DB: db, Sessions: sessions, Origin: s.Origin, ReviewRandom: func(int) (int, error) {
		return 0, errors.New("persisted sample must not be redrawn")
	}}).Routes()
	restarted := request("GET", path+"/review-queue", "", ta)
	check(restarted, 200)
	if restarted.Body.String() != first.Body.String() {
		t.Fatal("restart changed sample")
	}
	h = s.Routes()
	var queue []struct {
		Submission string `json:"submissionId"`
		Constraint string `json:"constraintId"`
		Required   bool   `json:"reviewRequired"`
		Sampled    bool   `json:"sampledPass"`
		Reason     string `json:"reviewReason"`
	}
	if err := json.Unmarshal(first.Body.Bytes(), &queue); err != nil {
		t.Fatal(err)
	}
	var sampled []string
	decisionPath := func(constraint string) string {
		return path + "/submissions/" + submission + "/constraints/" + constraint + "/review"
	}
	for _, row := range queue {
		if row.Sampled {
			sampled = append(sampled, row.Constraint)
			if row.Reason != "sampled-pass" {
				t.Fatal(row)
			}
			continue
		}
		if row.Required {
			check(request("POST", decisionPath(row.Constraint), `{"decision":"violation","deduction":999,"evidence":[{"line":1,"text":"int main() {}"}]}`, ta), 400)
			check(request("POST", decisionPath(row.Constraint), `{"decision":"violation","evidence":[{"line":1,"text":"int main() {}"}]}`, ta), 200)
		}
	}
	if len(sampled) != 2 {
		t.Fatalf("sampled %v", sampled)
	}
	if p = getProgress(); p.Pending != 2 {
		t.Fatal(p)
	}
	check(request("POST", path+"/publish", body(p.Revision), ta), 409) // sampled passes are mandatory too
	for _, constraint := range sampled {
		check(request("POST", decisionPath(constraint), `{"decision":"pass"}`, ta), 200)
	}
	if p = getProgress(); p.Pending != 0 || p.Reviewed != 4 {
		t.Fatal(p)
	}
	check(request("GET", resultPath, "", student), 404) // reviews alone never publish
	check(request("POST", path+"/publish", body(p.Revision), ta), 200)
	check(request("POST", path+"/publish", body(p.Revision), ta), 200) // repeat safely
	out := request("GET", resultPath, "", student)
	check(out, 200)
	var result struct {
		Deduction float64 `json:"deduction"`
	}
	if err := json.Unmarshal(out.Body.Bytes(), &result); err != nil || result.Deduction != 6 {
		t.Fatalf("result %s: %v", out.Body, err)
	}
	check(request("GET", resultPath, "", other), 404)
	check(request("GET", resultPath, "", otherTA), 404)
	check(request("GET", resultPath, "", ta), 200)
	check(request("POST", decisionPath(sampled[0]), `{"decision":"pass"}`, ta), 200)
	check(request("GET", resultPath, "", student), 404) // editing a decision requires republishing
	check(request("POST", path+"/publish", body(p.Revision), ta), 200)

	// A committed rejudge invalidates even identical old findings/decisions.
	exec(`UPDATE jobs SET status='processing' WHERE submission_id=$1`, submission)
	if fresh := getProgress(); fresh.Pending != 4 || fresh.Incomplete != 1 || fresh.Published || fresh.Revision == p.Revision {
		t.Fatal(fresh)
	}
	check(request("GET", resultPath, "", student), 404)
	check(request("POST", path+"/publish", body(p.Revision), ta), 409)
	exec(`UPDATE jobs SET status='completed' WHERE submission_id=$1`, submission)

	// Constraint mutation holds the population lock while a real Publish request
	// waits. After commit it must see the new revision and reject the old one.
	p = getProgress()
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `UPDATE lab_constraints SET deduction=7 WHERE id=$1`, sampled[0]); err != nil {
		t.Fatal(err)
	}
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() { done <- request("POST", path+"/publish", body(p.Revision), ta) }()
	deadline := time.Now().Add(5 * time.Second)
	blocked := false
	for time.Now().Before(deadline) {
		if err = db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_stat_activity WHERE wait_event_type='Lock' AND query LIKE 'SELECT review_revision FROM labs%')`).Scan(&blocked); err != nil {
			t.Fatal(err)
		}
		if blocked {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if !blocked {
		t.Fatal("Publish did not wait on the population mutation")
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case out := <-done:
		check(out, 409)
		if !strings.Contains(out.Body.String(), "stale-review-population") {
			t.Fatal(out.Body.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Publish did not finish")
	}
	if fresh := getProgress(); fresh.Pending != 4 || fresh.Published {
		t.Fatal(fresh)
	}
	// Updating an existing finding invalidates reviews even on the same attempt.
	before := getProgress()
	exec(`UPDATE constraint_verdicts SET verdict='unsure' WHERE submission_id=$1 AND constraint_id=$2`, submission, sampled[0])
	if fresh := getProgress(); fresh.Revision == before.Revision || fresh.Reviewed != 0 {
		t.Fatal(fresh)
	}

	// Replacement attempt is counted immediately, even while pending; earlier
	// decisions cannot authorize publication and the old attempt is read-only.
	replacement := uuid.NewString()
	exec(`INSERT INTO submissions(id,lab_id,student_id,source,created_at) VALUES($1,$2,$3,$4,now()+interval '1 second')`, replacement, lab, student, []byte("replacement"))
	exec(`INSERT INTO jobs(id,submission_id) VALUES($1,$2)`, uuid.NewString(), replacement)
	p = getProgress()
	if p.Incomplete != 1 || p.Required != 0 {
		t.Fatal(p)
	}
	check(request("POST", path+"/publish", body(p.Revision), ta), 409)
	check(request("POST", decisionPath(sampled[0]), `{"decision":"pass"}`, ta), 404)
	check(request("GET", resultPath, "", student), 404)
	// Accepted with missing constraints remains incomplete.
	exec(`UPDATE jobs SET status='completed' WHERE submission_id=$1`, replacement)
	exec(`INSERT INTO judge_results(submission_id,verdict,detail) VALUES($1,'Accepted','OK')`, replacement)
	p = getProgress()
	if p.Incomplete != 1 {
		t.Fatal(p)
	}
	check(request("POST", path+"/publish", body(p.Revision), ta), 409)
	exec(`INSERT INTO constraint_verdicts(submission_id,constraint_id,verdict,source)
 SELECT $1,constraint_id,verdict,source FROM constraint_verdicts WHERE submission_id=$2`, replacement, submission)
	p = getProgress()
	if p.Incomplete != 0 || p.Reviewed != 0 || p.Pending == 0 {
		t.Fatal(p)
	}
	check(request("POST", path+"/publish", body(p.Revision), ta), 409)
	var history int
	if err = db.QueryRow(ctx, `SELECT count(*) FROM constraint_review_history WHERE submission_id=$1`, submission).Scan(&history); err != nil || history != 5 {
		t.Fatalf("history %d: %v", history, err)
	}
}
