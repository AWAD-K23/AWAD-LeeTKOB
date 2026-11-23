# LeeTKOB

LeeTKOB is a C++ programming-lab platform for students and teaching assistants.
Students submit code, a background worker judges it against private tests, and
assignment constraints are checked separately using syntax analysis and an
optional AI model.

**Accepted means the output tests passed.** A submission can still violate an
assignment rule. Constraint findings are intended for TA review before any marks
are published; the model never chooses deductions.

## What the app does

| Role | Current features |
| --- | --- |
| Admin | Create classes and assign existing TA/student accounts by email |
| Teaching assistant | Create labs with Markdown statements, private tests, requirements, and deductions |
| Student | Browse enrolled labs, edit highlighted C++, submit work, and view judging status, history, and saved source |

- The frontend includes light/dark themes, responsive pages, drafts, and transitions.
- The backend enforces roles, class membership, and submission ownership.
- PostgreSQL stores sessions, submissions, the durable job queue, and results.
- Hard rules check supported bans without AI. Gemini checks remaining semantic rules.
- Student responses exclude protected constraint findings and deductions.

**Current limits:** TA review UI, mark publication, and appeals are still pending.
Public deployment verification and final submission work are tracked separately.
See [the PA3 plan](docs/PA3/PA3-PLAN.md) for delivery status.

## How it works

```mermaid
flowchart TD
    Browser[React frontend] --> API[Go API]
    API -->|Save submission and pending job| DB[(PostgreSQL)]
    DB -->|Claim job| Worker[Go worker]
    Worker --> Sandbox[Restricted C++ sandbox]
    Sandbox --> Tests{Output tests pass?}
    Tests -->|No| Judge[Execution failure or Wrong Answer]
    Tests -->|Yes| Rules[Hard-rule checks]
    Rules --> Model[Optional semantic model checks]
    Model --> Results[Accepted plus separate constraint findings]
    Judge --> DB
    Results --> DB
    DB -->|Status and judge result| API
    Results -.-> Review[TA review and publication: pending]
```

The server and worker are separate programs connected through database rows.
Student code runs in restricted Docker containers with no network, host mounts,
Docker socket, or application secrets. The trusted worker controls Docker.

## Stack and repository

| Part | Technology | Location |
| --- | --- | --- |
| Frontend | React, TypeScript, Vite, Tailwind CSS, React Router, TanStack Query | `frontend/` |
| API and worker | Go, Chi, SCS sessions, pgx | `backend/` |
| Database | PostgreSQL 17 and Goose migrations | `backend/migrations/`, `compose.yaml` |
| Sandbox and scanner | Docker, Clang/libclang, C++17 | `sandbox/` |
| Semantic checking | Gemini API, validated JSON findings | `backend/internal/semantic/` |
| Evaluation | Generic checker, official dataset, unchanged Node evaluator | `eval/` |
| Deployment setup | Dockerfiles, Compose, Caddy, host-worker service template | `deploy/`, `compose.pa3.yaml` |

## Run locally

Requirements: Git, Node.js 22.13.1 or newer, npm, Go 1.26 or newer, and Docker
with Compose. Start Docker Desktop when using Windows.

### 1. Configure the environment

Copy `.env.example` to an ignored `.env` file at the repository root, then fill
in the settings below. Keep `.env.example` as variable names without values.

| Setting | Local development value |
| --- | --- |
| `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` | Choose local database credentials |
| `DATABASE_URL` | `postgres://USER:PASSWORD@localhost:5433/DATABASE?sslmode=disable`, using the same credentials |
| `API_PORT` | `8080`, matching the Vite proxy |
| `APP_ORIGIN` | `http://localhost:5173` |
| `COOKIE_SECURE` | `false` for local HTTP |
| `SANDBOX_IMAGE` | `leetkob-sandbox:pa3`, or leave empty for that default |
| `MODEL_OFF` | `1` for no model calls; `0` explicitly enables them |

The API, worker, and checker load `.env` automatically from the current directory
or its parent. Existing process environment variables take precedence. Restart
running programs after changing settings.

### 2. Start PostgreSQL and prepare judging

Run from the repository root, checking each command succeeds:

```text
docker compose up -d
docker build -t leetkob-sandbox:pa3 sandbox
```

Apply the schema using the migration container, which receives database settings
from Compose:

```text
docker compose -f compose.pa3.yaml run --build --rm migrate
```

Alternatively, supply `DATABASE_URL` in your shell and run:

```text
go -C backend run ./cmd/migrate
```

PostgreSQL uses host port **5433** and the named `postgres_data` volume. Normal
restarts preserve data. `docker compose down -v` deletes persistent volumes.

### 3. Create accounts

There is no automatic account seed or registration endpoint. To create each
admin, TA, or student account, supply these variables in the command's shell:

| Variable | Value |
| --- | --- |
| `DATABASE_URL` | Local database connection |
| `USER_EMAIL` | Unique email used for login |
| `USER_PASSWORD` | Password of at least 12 bytes |
| `USER_ROLE` | `admin`, `ta`, or `student` |

```text
go -C backend run ./cmd/create-user
```

The migration and account commands require shell settings; they do not load
`.env` automatically. See the [manual walkthrough](docs/PA3/MANUAL-FLOW-TEST.md)
for shell examples, existing local test accounts, and sample lab data. Those
accounts exist only in databases where they were created.

### 4. Start the three programs

Install frontend dependencies once:

```text
npm ci --prefix frontend
```

Run each command in a separate terminal from the repository root:

| Program | Command |
| --- | --- |
| API | `go -C backend run ./cmd/server` |
| Worker | `go -C backend run ./cmd/worker` |
| Frontend | `npm run dev --prefix frontend` |

Open **http://localhost:5173**. On Windows, use `npm.cmd` if PowerShell blocks
the npm script launcher.

Sign in as admin to create a class, then as its TA to create a lab, and finally
as an enrolled student to submit code. The worker must be running to process
queued submissions. It can also process at most one job with `--once`.

Browser API requests use `/api/*`; Vite and Caddy strip that prefix before
forwarding to Go. Direct clients on port 8080 use unprefixed paths such as `/me`.

## Model checking and evaluation

The supported model is `gemini-3.5-flash-lite`. With `MODEL_OFF=1`, hard rules
still run and semantic constraints remain unsure without provider calls.

For real checking, set `MODEL_NAME` and `MODEL_API_KEY` in the ignored `.env`,
verify the intended provider spending cap, then explicitly set `MODEL_OFF=0`.
The key belongs only in the trusted worker or standalone checker environment.

| Control | Behavior |
| --- | --- |
| Input | Statement, semantic rules, and code with comments and known identities removed |
| Response | Validated verdict/evidence per constraint; never a deduction |
| Failure | Invalid answers, timeouts, and blocked calls become unsure |
| Requests | Calls are paced within each process; no automatic provider retries |
| Production budget | Shared by the lab and retained across worker restarts |

The worker reserves budget before contacting the provider. Known usage replaces
the reservation with a token-based estimate; unknown usage keeps the reservation.
The optional `MODEL_LAB_BUDGET_MICRO_USD` setting controls the local cap. Multiple
processes sharing a key need coordinated pacing.

The standalone checker reads statements, constraints, and C++ solutions from a
dataset folder and writes a prediction file. The separate evaluator compares
those predictions with expected labels. The checker never reads labels to choose
its answers. Official data and the evaluator are preserved unchanged; our own
cases are stored separately.

Predictions, provider logs, and evaluator reports are saved for review. See
[evaluation documentation](eval/README.md) and
[Member 2 handoff](docs/PA3/MEMBER2-HANDOFF.md) for detailed results and limitations.

Build and run the standalone model-off checker with no provider access:

```text
docker build -f deploy/Dockerfile.checker -t leetkob-checker:pa3 .
node scripts/verify-eval.mjs
node scripts/run-offline-eval.mjs
```

The script disables container networking and applies resource limits. It saves
predictions and the original evaluator report in a new local run folder. See
[evaluation commands](eval/README.md) for other datasets, real calls, and reruns.

## Checks

Run from the repository root:

```text
node scripts/check-go-format.mjs
go -C backend vet ./...
go -C backend test ./... -count=1
npm ci --prefix frontend
npm run check --prefix frontend
```

- Set `TEST_DATABASE_URL` in the test process to include PostgreSQL integration
  checks. They create and remove only isolated temporary schemas.
- Set `SANDBOX_TESTS=1` and build the sandbox image to include real Docker attacks
  and syntax checks. `scripts/test-pa3.ps1` runs the full backend checks on Windows.
- Unit tests use fake models and never call a real provider.
- Install Chromium once with `npx --prefix frontend playwright install chromium`,
  then run `npm run e2e --prefix frontend` for browser workflows.
- `npm audit --prefix frontend` checks frontend dependencies.
- CI runs `gates`, `secret-history`, `e2e`, and `sandbox` on pull requests, `main`
  pushes, and the `pa2`/`pa3` tags. Offline checker checks run in the sandbox job.
- `secret-history` uses the standalone MIT-licensed Gitleaks CLI v8.30.1,
  downloaded with a pinned SHA-256 checksum. It scans the full fetched Git
  history with `gitleaks git --log-opts="--all" --redact --exit-code=1 .` and
  fails on findings or scanner errors. It needs no `GITLEAKS_LICENSE` secret.

## Documentation

| Guide | Contents |
| --- | --- |
| [Backend structure](docs/BACKEND.md) | Go files, commands, routes, and model controls |
| [Database structure](docs/DATABASE.md) | Tables, relationships, queue state, and cost reservations |
| [Backend connections](docs/PA3/backend-mechanism/BACKEND-CONNECTIONS.md) | Server, worker, sandbox, and hard rules |
| [Model checker flow](docs/PA3/backend-mechanism/MODEL-CHECKER-FLOW.md) | Visual walkthrough with real data and code snippets |
| [Manual test](docs/PA3/MANUAL-FLOW-TEST.md) | Accounts, example labs, submission cases, and database checks |
| [Evaluation](eval/README.md) | Dataset layout, checker image, commands, and reports |
| [Deployment setup](docs/PA3/MEMBER1-HANDOFF.md) | Compose/Caddy setup and separate host worker |
| [PA3 plan](docs/PA3/PA3-PLAN.md) | Completed work, remaining tasks, and ownership |
| [AI log](AI-LOG.md) | Assistant use, changes, rejected suggestions, and verification |

The local deployment stack is configured in `compose.pa3.yaml`; the worker runs
separately on the Docker host. Deployment files and local evidence do not claim
that public hosting, the TA review workflow, or final tag CI are complete.
