# Frontend browsing endpoints

The browser calls `/api/*`. Vite and Caddy strip `/api` before proxying to Go.
Direct API clients continue to use the unprefixed routes on the Go server.
Browser page routes such as `/classes` and `/labs/:id` serve the React application.

All endpoints below require a valid server session. Failures use the existing
`application/problem+json` format. Authentication failures return 401.

| Endpoint | Access | Response |
| --- | --- | --- |
| `GET /people?q=...` | Admin only; others receive 403 | Up to 200 email matches, ordered by email: `{id,email,role}[]`. Password data is never returned. |
| `GET /classes/{classId}` | Admin or class member; missing/inaccessible classes return 404 | `{id,name,members}`. Admin and assigned TA see member ID/email/role; students receive an empty member array. |
| `GET /submissions` | Student only; other roles receive 403 | Latest 200 submissions owned by the session account: `{id,labId,title,createdAt,status,verdict}[]`. |

`GET /labs` additionally returns `classId` and `requirements` (constraint description
strings). Tests, checker findings, deductions, and unpublished marks are excluded.
Existing TA/student lab access rules remain enforced.

`GET /submissions/{submissionId}` additionally returns saved `source` and `labId`.
Its existing ownership rule remains: submission owner or TA assigned to its class.
Source is read from the database byte column and encoded as a JSON string. The
response continues to exclude deductions and protected constraint verdicts.

These additions require no database migration. Account provisioning, judging,
constraint review, publication, and appeals are separate capabilities.
