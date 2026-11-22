# Manual frontend flow test

## 1. Setup

- Start Docker Desktop. Run commands from the repository root in PowerShell.
- Keep existing database credentials in `.env`.
- Set `API_PORT=8080`, `APP_ORIGIN=http://localhost:5173`, and `COOKIE_SECURE=false`.
- `DATABASE_URL` must use host port `5433`.

Load environment settings in every terminal running migrations, account creation,
or the worker:

```powershell
Get-Content .env |
  Where-Object { $_ -match '^[A-Za-z_][A-Za-z0-9_]*=' } |
  ForEach-Object {
    $name, $value = $_ -split '=', 2
    [Environment]::SetEnvironmentVariable(
      $name, $value.Trim().Trim('"', "'"), 'Process'
    )
  }
```

Run each setup command successfully before continuing:

```powershell
docker compose up -d
docker build -t leetkob-sandbox:pa3 sandbox
go -C backend run ./cmd/migrate
npm.cmd ci --prefix frontend
```

## 2. Accounts

- Login uses the email address as the username.
- These local test accounts were created for the manual test:

| Role | Username / email | Password |
| --- | --- | --- |
| Admin | `manual-admin@example.test` | `manual-test-password-123` |
| Teaching assistant | `manual-ta@example.test` | `manual-test-password-123` |
| Assigned student | `manual-student@example.test` | `manual-test-password-123` |
| Unassigned student | `manual-outsider@example.test` | `manual-test-password-123` |

- These are local test credentials, not production credentials.
- For a fresh database, run this once after loading `.env`:

```powershell
foreach ($role in @('admin', 'ta', 'student', 'outsider')) {
  $env:USER_EMAIL = "manual-$role@example.test"
  $env:USER_PASSWORD = 'manual-test-password-123'
  $env:USER_ROLE = if ($role -eq 'outsider') { 'student' } else { $role }

  go -C backend run ./cmd/create-user
  if ($LASTEXITCODE -ne 0) { throw 'Account creation failed; check duplicate emails.' }
}

Remove-Item Env:USER_EMAIL, Env:USER_PASSWORD, Env:USER_ROLE
```

## 3. Start the app

| Terminal | Command from repository root |
| --- | --- |
| A: API | `go -C backend run ./cmd/server` |
| B: Frontend | `npm.cmd run dev --prefix frontend` |
| C: Worker | Leave stopped initially |

Open **http://localhost:5173**. Stop older API/worker instances before testing.

## 4. Create class and lab

1. Admin: `manual-admin@example.test` → **Classes → Create class**.
2. Name: `Manual Flow 01`; TA: `manual-ta@example.test`.
3. Assign only `manual-student@example.test`; leave the outsider unassigned.
4. Sign out. TA signs in → class → **Create lab**.

| Lab field | Value |
| --- | --- |
| Title | `Order two integers` |
| Statement | Read two integers. Print them in ascending order, separated by a space. |
| Test 1 | Input `3 1`; output `1 3` |
| Test 2 | Input `-2 5`; output `-2 5` |
| Test 3 | Input `4 4`; output `4 4` |
| Rule 1 | ``Do not call `std::sort` ``; deduction `3` |
| Rule 2 | `Use a comparison to order the numbers`; deduction `2` |

Create the lab. Sign out, then sign in as `manual-student@example.test`.

## 5. Queue and correct result

Open the lab and submit:

```cpp
#include <iostream>

int main() {
    int a, b;
    if (!(std::cin >> a >> b)) return 1;

    // Order the numbers without a banned function.
    if (a > b) {
        int temporary = a;
        a = b;
        b = temporary;
    }

    std::cout << a << " " << b << "\n";
}
```

- With no worker running: expect **Queued**.
- In terminal C, load `.env` using step 1, then run:

```powershell
$env:SANDBOX_IMAGE = 'leetkob-sandbox:pa3'
go -C backend run ./cmd/worker
```

- Expect **Accepted** automatically. Check **My submissions**.

## 6. Failure cases

Replace the source and submit each program:

| Test | Source | Expected |
| --- | --- | --- |
| Wrong output | `#include <iostream>` followed by `int main() { std::cout << "999 999"; }` | Wrong Answer |
| Invalid C++ | `int main() { this is not valid C++; }` | Compilation Error |
| Endless loop | `int main() { for (volatile unsigned long long n = 0;; ++n) {} }` | Time Limit Exceeded |

## 7. Correct output with a banned call

Submit this last, before running the database query:

```cpp
#include <algorithm>
#include <iostream>

int main() {
    int values[2];
    if (!(std::cin >> values[0] >> values[1])) return 1;

    // Passes judge tests but violates the named ban.
    std::sort(values, values + 2);
    std::cout << values[0] << " " << values[1] << "\n";
}
```

| Where | Expected |
| --- | --- |
| Student UI | Accepted |
| Database: Rule 1 | `violation`, source `hard_rule`, evidence containing the `std::sort` line |
| Database: Rule 2 | `unsure`, source `semantic_pending` |

Read the latest matching submission's findings from a free terminal at the root:

```powershell
@'
WITH latest AS (
    SELECT s.id
    FROM submissions s
    JOIN users u ON u.id = s.student_id
    JOIN labs l ON l.id = s.lab_id
    JOIN classes c ON c.id = l.class_id
    WHERE u.email = 'manual-student@example.test'
      AND c.name = 'Manual Flow 01'
      AND l.title = 'Order two integers'
    ORDER BY s.created_at DESC LIMIT 1
)
SELECT lc.description, cv.verdict, cv.source, cv.evidence
FROM constraint_verdicts cv
JOIN lab_constraints lc ON lc.id = cv.constraint_id
JOIN latest ON latest.id = cv.submission_id;
'@ | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
```

- Accepted means output tests passed, not that all rules passed.
- No real AI check, TA review UI, automatic deduction, or mark publication exists in this flow.

## 8. Access and restart checks

| Check | Expected |
| --- | --- |
| Copy lab/submission URLs; sign in as `manual-outsider@example.test`; open them | Access denied / not found; lab absent from list |
| Stop terminal C with Ctrl+C; submit correct code as the assigned student | Queued |
| Restart `go -C backend run ./cmd/worker` in terminal C | Existing submission becomes Accepted without resubmitting |

- Keep existing data. Do not use `docker compose down -v`.
- For an independent rerun, use a new class name and update the query filter.
