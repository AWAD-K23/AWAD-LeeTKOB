# Member 1 handoff

## Delivered parts

| Part | Implementation | Proof |
|---|---|---|
| Class setup | Admin selects existing TA/students by email; class details respect membership | PostgreSQL API tests and admin browser test |
| Lab creation | TA supplies Markdown statement, private tests, constraints, and deductions; preview and unsaved-change protection | API tests and TA browser test |
| Submission | Enrolled student submits C++; API saves submission and job atomically | Existing contract suite and class flow test |
| Student results | Protected status, saved source, history, and 1.5-second active-job polling; C++ highlighting and local source drafts | Ownership tests, student browser test, and rendering test |
| Worker | Separate Go executable; PostgreSQL leases and fencing tokens | Restart/recovery and duplicate-completion tests |
| Sandbox | Fresh non-root Docker container for compilation, execution, and analysis | Nine real attack submissions |
| Hard rules | libclang resolves C++ references; Go matches generic backtick ban lists | Seven syntax cases, including macros and compliant traps |
| CI | Existing jobs plus Docker/PostgreSQL sandbox job; `pa3` trigger | Same commands passed locally; remote tag run belongs to final submission |
| Deployment setup | Web/API/migration images, Compose, Caddy, and worker service template; browser requests now use `/api` | Original local stack smoke passed; [post-redesign local smoke](LOCAL-COMPOSE-SMOKE-2026-11-24.md) recorded November 24; public deployment remains pending |

## Follow-up verification — November 22, 2026

- Go format, vet, and PostgreSQL-backed tests passed after adding browsing endpoints.
- Eight workflow browser tests passed, including sign-in/out, admin account selection,
  TA authoring, unsaved navigation, student judging, retry, transitions, and mobile drafts.
- The additional Markdown/C++ rendering test passed in a targeted run after correcting
  its selector to account for both include-path and string-literal tokens.
- Frontend lint, formatting, type-checking, and production build passed. Two existing
  Fast Refresh lint warnings remain; Vite also reports a bundle-size warning.
- Real Docker attack tests were not rerun as part of these frontend changes.
- Local database inspection found submission `7f6b8932-d056-4d5f-996d-bdd32e4a2f6b`
  Accepted, with a separate violation for the `std::sort` ban and evidence at line 12.
  Its comparison requirement remained `unsure` / `semantic_pending`.
- Accepted does not mean all assignment constraints passed. The student UI intentionally
  excludes protected findings; the TA review endpoint/UI is still Member 3 work.
- Public hosting, final tag CI, teammate acceptance, and final public model-cap
  verification are not claimed by these local checks.
- **Assigned to Member 3:** rebuild and smoke-test the current Docker deployment
  using `compose.pa3.yaml`. Check frontend loading, direct page refreshes, `/api/me`,
  sign-in, class/lab access, and a real worker-processed submission. Record the tested
  commit and results. Member 1 remains responsible for fixes to their supplied
  deployment files and backend components.

## Follow-up verification — November 24, 2026

The post-redesign local Compose and Caddy smoke test ran against merged commit `1122f7c` with `MODEL_OFF=1`. It covered direct page refresh, sign-in, separate class ownership, a TA-reviewed hard-rule deduction, and worker restart. See the [smoke record](LOCAL-COMPOSE-SMOKE-2026-11-24.md). The November 22 assignment above records status at that time. Public deployment, live model checks, teammate acceptance, and final `pa3` tag CI remain pending.



## Fixed interfaces for Members 2 and 3

- Receiving members run the acceptance checks before implementing their part.
- These interfaces are ready for team review; no separate teammate approval is claimed.
- `hardrules.Check(ctx, scanner, source, constraints)` returns one verdict per constraint.
- `Constraint` contains `id` and `text`. Product constraints use UUIDs; the dataset checker can supply dataset IDs.
- `Evidence` contains a one-based `line` and the exact `text` from that source line.
- `Verdict` contains `constraintId`, `verdict`, `evidence`, and `source`.

| Source | Meaning | Next owner |
|---|---|---|
| `hard_rule` | Ban-list decision; analysis failure stays `unsure` | Member 3 reviews it; do not call the model for it |
| `semantic_pending` | Semantic constraint not yet checked by a model | Member 2 adds the real model, schema validation, and timeout |
| `model` | Reserved database value for Member 2's validated model verdict | Member 3 displays it in the TA queue |

- TA-selected deductions live in `lab_constraints.deduction`; they are never set by the checker.
- Judge verdicts live in `judge_results`; constraint verdicts live separately in `constraint_verdicts`.
- There is no publication action or automatic deduction in this build.
- Member 2 integrates semantic checking before `Worker.Complete`, preserving the lease check and transaction.
- Member 3 reads verdicts through a TA/class-protected endpoint they add. Student status must continue to exclude constraint verdicts and deductions.

## API contracts

- The paths below are direct Go API routes. Browser requests use `/api` before
  those paths; Vite and Caddy remove that prefix before forwarding to Go.
- Browser page routes such as `/classes` serve React. A deployed status smoke check
  must call `/api/me`, not `/me`.
- See [frontend browsing contracts](../contract/frontend-browsing.md) for added fields
  and access rules.

| Endpoint | Access | Result |
|---|---|---|
| `GET /me` | Signed-in account | ID and current role |
| `POST /classes` | Admin + valid Origin | Create class, assign existing TA/students |
| `GET /classes` | Signed-in account | Classes the account belongs to; admin sees all |
| `GET /classes/{classId}` | Admin or class member | Class details; roster available only to admin/assigned TA |
| `GET /people?q=...` | Admin | Up to 200 matching accounts; no password data |
| `POST /classes/{classId}/labs` | TA assigned to that class + valid Origin | Lab and related rows saved together; students enrolled |
| `GET /labs` | TA or student | Accessible labs, statements, class IDs, and public requirement descriptions; no private tests/deductions |
| `POST /labs/{labId}/submissions` | Enrolled student + valid Origin | Existing PA2 contract, including retry key |
| `GET /submissions` | Student | Latest 200 submissions belonging to that student |
| `GET /submissions/{submissionId}` | Submission owner or class TA | Job status, judge result, saved source, and lab ID; no constraint verdicts/deductions |

- Lab body example:

```json
{
  "title": "Merge sort",
  "statement": "Sort the input numbers using merge sort.",
  "tests": [{ "input": "3 1", "expectedOutput": "1 3" }],
  "constraints": [
    { "description": "Do not call `std::sort`", "deduction": 3 },
    { "description": "Use merge sort", "deduction": 2 }
  ]
}
```

- Judge status example:

```json
{
  "id": "8d00cb12-af80-4eb4-bbc5-dc216cb26740",
  "status": "completed",
  "verdict": "Accepted",
  "detail": "All test cases passed.",
  "source": "int main() {}",
  "labId": "6554bbe5-af9c-4340-afef-211b513efc80"
}
```

## Run and verify locally

- Use the ignored `.env` described in README. PostgreSQL uses host port **5433**.
- Start the database with `docker compose up -d`.
- Follow [manual frontend testing](MANUAL-FLOW-TEST.md) for accounts, lab data,
  submissions, queue/restart checks, and read-only database evidence queries.
- The API server loads local `.env` automatically. The worker, migration command,
  and account command still need settings loaded into their own shell processes.
- Development uses `http://localhost:5173` as `APP_ORIGIN`, with `COOKIE_SECURE=false`.
- In PowerShell, run the full backend/Docker checks:

```powershell
./scripts/test-pa3.ps1
```

- Or set `TEST_DATABASE_URL` and `SANDBOX_TESTS=1` in the shell, then run:

```text
docker build -t leetkob-sandbox:pa3 sandbox
node scripts/check-go-format.mjs
go -C backend vet ./...
go -C backend test ./... -count=1 -v
npm ci --prefix frontend
npm run check --prefix frontend
npm run e2e --prefix frontend
```

- Install Chromium once with `npx --prefix frontend playwright install chromium` (Linux CI also uses `--with-deps`).
- Integration tests create and remove only their own temporary schemas.
- `TestSubmissionContract/class_setup_and_status_ownership` creates the class and lab via HTTP, submits real code, invokes the worker, and reads Accepted through HTTP when `SANDBOX_TESTS=1`.
- `TestDurableWorker` proves queued work survives an absent worker, a new worker stores one result set, and an expired lease cannot let an old worker overwrite results.

## Deployment environment

- The local deployment is configured for `http://localhost:8088` after starting Compose below.
  Rebuild/smoke-test the current redesign before treating the old deployment result as current evidence.
- Set these values in the ignored `.env`: `APP_ORIGIN=http://localhost:8088`, `COOKIE_SECURE=false`, `PUBLIC_ADDRESS=:80`, and `WEB_PORT=8088`.
- Keep `MODEL_OFF=1` until Member 2 supplies the real model path and limits. This setting is reserved configuration; current Member 1 code never calls a model.

```text
docker build -t leetkob-sandbox:pa3 sandbox
docker compose -f compose.pa3.yaml up -d --build
go -C backend build -o ../bin/worker.exe ./cmd/worker
```

- Load `DATABASE_URL` into the shell without printing it, then run `./bin/worker.exe` in a separate terminal.
- Create test accounts with the existing `cmd/create-user` command. The admin selects
  accounts by email in the class form; the frontend submits their IDs. Labs and
  submissions are created through the app, not seed scripts.
- The worker runs on the Docker host. Its service template is `deploy/leetkob-worker.service`; Linux uses `bin/worker` instead of `worker.exe`.
- A trusted worker controls Docker. Student containers receive neither Docker access nor worker/API environment variables.
- Worker and database may be restarted without deleting volumes. Never use `down -v` to restart this environment.
- On a public Linux host, set the real HTTPS origin, `COOKIE_SECURE=true`, the real domain in `PUBLIC_ADDRESS`, and web ports 80/443. Build the Linux worker with Go 1.26, install it under `/opt/leetkob/bin`, and use the supplied service template.
- Public provisioning, DNS, and final publication are not claimed as completed by the local deployment test. Member 3's public walkthrough and final tag CI remain separate checklist items.

## Spending cap and secrets

- Existing saved evidence documents a **₫100,000 monthly Google AI Studio project cap**;
  see README. Verify the final model project/key and its cap before public testing.
  No fresh provider verification was performed in this documentation review.
  Billing latency can permit overages, as already documented.
- No paid resource or provider setting was changed in this implementation.
- Only the server receives `MODEL_API_KEY`. The frontend build and sandbox receive no key, `.env` file, Docker socket, or host mount.
- Member 2 must implement per-lab model cost enforcement and measure the 80-submission cost before enabling real calls. The Member 1 environment starts with model calls off.

## Scope decisions

- This implementation uses a small custom Docker runner rather than the provisional Judge0 choice in the PA2 plan. The nine requirements have direct attack tests, and a second judging service is not required.
- Output comparison ignores whitespace differences between tokens.
- Runtime limits: 2 CPU seconds, 20 seconds total wall time including compilation, 256 MiB memory with no additional swap, 32 processes, and 64 KiB captured output. Tests use an 8-second wall limit to keep CI short.
- Each test case compiles in a fresh sandbox; caching/compile-once optimization belongs to later performance work.
- Public deployment and human teammate acceptance still require the team's final handoff process; the local environment and interfaces are ready for it.
