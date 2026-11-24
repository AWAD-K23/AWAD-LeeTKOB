#!/usr/bin/env node
// PA#3 criterion 1 walkthrough through the product's own API, against any
// deployment (local Compose or the public URL). It does what the marker does:
//   setup  admin creates class A (TA1 + both students) and class B (TA2)
//   walk   TA1 creates the dev-set sort-merge lab; students submit a compliant,
//          a violation and an injection solution; the script waits for judging,
//          reads TA1's review queue, and checks that TA2 and a student are refused.
// Credentials come only from environment variables; nothing is written to disk
// except the optional JSON evidence file given with --out.
//
// Usage:
//   LK_URL=https://host LK_ADMIN=email:password LK_TA1=... LK_TA2=... \
//   LK_STUDENT1=... LK_STUDENT2=... node scripts/pa3-walkthrough.mjs setup|walk [--out file.json]
import { mkdirSync, readFileSync, readdirSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { randomUUID } from "node:crypto";

const [mode, ...rest] = process.argv.slice(2);
const outIndex = rest.indexOf("--out");
const outFile = outIndex >= 0 ? rest[outIndex + 1] : "";
const base = (process.env.LK_URL ?? "").replace(/\/$/, "");
if (!base || !["setup", "walk"].includes(mode)) {
  console.error("Usage: LK_URL=... LK_ADMIN=... LK_TA1=... LK_TA2=... LK_STUDENT1=... LK_STUDENT2=... node scripts/pa3-walkthrough.mjs setup|walk [--out file.json]");
  process.exit(2);
}
const origin = new URL(base).origin;
const dataset = "eval/official";
const problem = "sort-merge";
// One compliant, one plain violation (std::sort), one injection case, as the
// PA#3 rubric's criterion 1 asks. Labels live in eval/official/labels.jsonl.
const cases = [
  { who: "student1", file: "s01.cpp", kind: "compliant" },
  { who: "student2", file: "s04.cpp", kind: "violation (calls std::sort)" },
  { who: "student1", file: "s09.cpp", kind: "injection (bubble sort + comment aimed at the model)" },
];

function account(name) {
  const raw = process.env[`LK_${name.toUpperCase()}`] ?? "";
  const at = raw.indexOf(":");
  if (at < 1) throw new Error(`Set LK_${name.toUpperCase()}=email:password`);
  return { email: raw.slice(0, at), password: raw.slice(at + 1) };
}

const log = [];
function record(step, ok, detail) {
  log.push({ step, ok, detail });
  console.log(`${ok ? "PASS" : "FAIL"}  ${step}${detail ? ` — ${detail}` : ""}`);
}

async function session(name) {
  const { email, password } = account(name);
  const res = await fetch(`${base}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json", Origin: origin },
    body: JSON.stringify({ email, password }),
  });
  if (res.status !== 204) throw new Error(`${name} sign-in returned ${res.status}`);
  const cookie = res.headers.getSetCookie().map((c) => c.split(";")[0]).join("; ");
  const call = async (method, path, body, extra = {}) => {
    const r = await fetch(`${base}/api${path}`, {
      method,
      headers: { Cookie: cookie, Origin: origin, ...(body ? { "Content-Type": "application/json" } : {}), ...extra },
      body: body ? JSON.stringify(body) : undefined,
    });
    const text = await r.text();
    let json = null;
    try { json = text ? JSON.parse(text) : null; } catch { json = text; }
    return { status: r.status, json };
  };
  const me = await call("GET", "/me");
  return { name, email, id: me.json?.id, role: me.json?.role, call };
}

async function setup() {
  const admin = await session("admin");
  record("admin signs in", admin.role === "admin", `role=${admin.role}`);
  const ids = {};
  for (const name of ["ta1", "ta2", "student1", "student2"]) {
    const { email } = account(name);
    const found = await admin.call("GET", `/people?q=${encodeURIComponent(email)}`);
    const person = (found.json ?? []).find((p) => p.email === email.toLowerCase());
    if (!person) throw new Error(`No account for ${name}; create it with cmd/create-user first`);
    ids[name] = person.id;
  }
  const a = await admin.call("POST", "/classes", { name: "PA3 Class A — Sorting lab", taId: ids.ta1, studentIds: [ids.student1, ids.student2] });
  record("admin creates class A (TA1, two students)", a.status === 201, `status=${a.status}`);
  const b = await admin.call("POST", "/classes", { name: "PA3 Class B — other TA", taId: ids.ta2, studentIds: [] });
  record("admin creates class B (TA2)", b.status === 201, `status=${b.status}`);
}

function labFromDataset() {
  const dir = join(dataset, "problems", problem);
  const statement = readFileSync(join(dir, "statement.md"), "utf8");
  const constraints = JSON.parse(readFileSync(join(dir, "constraints.json"), "utf8"));
  const testDir = join(dir, "tests");
  const tests = readdirSync(testDir).filter((f) => f.endsWith(".in")).sort().map((f) => ({
    input: readFileSync(join(testDir, f), "utf8"),
    expectedOutput: readFileSync(join(testDir, f.replace(/\.in$/, ".out")), "utf8"),
  }));
  return {
    title: `Sorting lab: merge sort (${new Date().toISOString().slice(0, 16)}Z)`,
    statement,
    tests,
    constraints: constraints.map((c) => ({ description: c.text, deduction: c.deduction })),
  };
}

async function walk() {
  const ta1 = await session("ta1");
  const ta2 = await session("ta2");
  const students = { student1: await session("student1"), student2: await session("student2") };
  record("TA1, TA2 and two students sign in", [ta1, ta2].every((s) => s.role === "ta") && Object.values(students).every((s) => s.role === "student"));

  const classes = await ta1.call("GET", "/classes");
  const classA = (classes.json ?? [])[0];
  if (!classA) throw new Error("TA1 has no class; run setup first");
  const lab = await ta1.call("POST", `/classes/${classA.id}/labs`, labFromDataset());
  record("TA1 creates the sort-merge lab with its tests and two constraints", lab.status === 201, `status=${lab.status} lab=${lab.json?.id}`);
  const labId = lab.json?.id;

  const submitted = [];
  for (const c of cases) {
    const s = students[c.who];
    const visible = await s.call("GET", "/labs");
    const sees = (visible.json ?? []).some((l) => l.id === labId);
    const source = readFileSync(join(dataset, "solutions", problem, c.file), "utf8");
    const r = await s.call("POST", `/labs/${labId}/submissions`, { source }, { "Idempotency-Key": randomUUID() });
    record(`${c.who} sees the lab and submits ${problem}/${c.file} (${c.kind})`, sees && r.status === 201, `status=${r.status} submission=${r.json?.id} job=${r.json?.status}`);
    submitted.push({ ...c, id: r.json?.id });
  }

  const deadline = Date.now() + 10 * 60 * 1000;
  for (const sub of submitted) {
    let last;
    while (Date.now() < deadline) {
      last = await students[sub.who].call("GET", `/submissions/${sub.id}`);
      if (["completed", "failed"].includes(last.json?.status)) break;
      await new Promise((r) => setTimeout(r, 2000));
    }
    sub.judge = last.json?.verdict;
    // Check field names, not values: the source code itself may mention "constraint".
    const leaked = Object.keys(last.json ?? {}).some((k) => /constraint|deduction|evidence|review|reason|finding/i.test(k));
    record(`${sub.who} sees ${sub.file} judged`, last.json?.status === "completed", `status=${last.json?.status} verdict=${last.json?.verdict}`);
    record(`${sub.who}'s status response carries no findings or deductions`, !leaked);
  }

  const queue = await ta1.call("GET", `/labs/${labId}/review-queue`);
  record("TA1 opens the review queue", queue.status === 200, `rows=${(queue.json ?? []).length}`);
  const rows = (queue.json ?? []).map((q) => ({
    submission: q.submissionId,
    file: submitted.find((s) => s.id === q.submissionId)?.file ?? "(earlier submission)",
    constraint: q.constraint.slice(0, 40),
    verdict: q.verdict,
    source: q.source,
    reason: q.reason,
    evidence: q.evidence,
    deduction: q.deduction,
  }));
  for (const row of rows) console.log(`      ${row.file}  ${row.constraint}…  ${row.verdict}/${row.source}${row.reason ? ` (${row.reason})` : ""}  evidence=${JSON.stringify(row.evidence)}  deduction=${row.deduction}`);

  // FAQ-02 keeps only each student's latest submission in the queue; every
  // submission that a newer one replaced must still be open to the TA (FAQ-10).
  const earlierRes = await ta1.call("GET", `/labs/${labId}/earlier-submissions`);
  const earlierIds = (earlierRes.json ?? []).map((e) => e.submissionId);
  const replaced = submitted.filter((s) => submitted.some((n) => n.who === s.who && submitted.indexOf(n) > submitted.indexOf(s)));
  record(
    "TA1 sees each replaced submission under earlier submissions (not counted)",
    earlierRes.status === 200 && replaced.every((s) => earlierIds.includes(s.id)) && rows.every((r) => !earlierIds.includes(r.submission)),
    `status=${earlierRes.status} earlier=${replaced.map((s) => s.file).join(",")}`,
  );
  const earlier = (earlierRes.json ?? []).map((e) => ({
    submission: e.submissionId,
    file: submitted.find((s) => s.id === e.submissionId)?.file ?? "(other)",
    judge: e.judgeVerdict,
    findings: e.findings.map((f) => ({ constraint: f.constraint.slice(0, 40), verdict: f.verdict, source: f.source, reason: f.reason, evidence: f.evidence })),
  }));
  for (const e of earlier) for (const f of e.findings) console.log(`      ${e.file} (earlier)  ${f.constraint}…  ${f.verdict}/${f.source}${f.reason ? ` (${f.reason})` : ""}  evidence=${JSON.stringify(f.evidence)}`);

  const other = await ta2.call("GET", `/labs/${labId}/review-queue`);
  record("TA2 (other class) is refused the review queue", other.status === 404, `status=${other.status}`);
  const otherLabs = await ta2.call("GET", "/labs");
  record("TA2 does not see the lab in their list", !(otherLabs.json ?? []).some((l) => l.id === labId));
  const otherLab = await ta2.call("POST", `/classes/${classA.id}/labs`, labFromDataset());
  record("TA2 cannot create a lab in class A", otherLab.status === 404, `status=${otherLab.status}`);
  const otherEarlier = await ta2.call("GET", `/labs/${labId}/earlier-submissions`);
  record("TA2 (other class) is refused the earlier submissions", otherEarlier.status === 404, `status=${otherEarlier.status}`);
  const asStudent = await students.student1.call("GET", `/labs/${labId}/review-queue`);
  record("a student is refused the review queue", asStudent.status === 403, `status=${asStudent.status}`);
  const studentEarlier = await students.student1.call("GET", `/labs/${labId}/earlier-submissions`);
  record("a student is refused the earlier submissions", studentEarlier.status === 403, `status=${studentEarlier.status}`);
  const peek = await students.student2.call("GET", `/submissions/${submitted[0].id}`);
  record("a student cannot read another student's submission", peek.status === 404 || peek.status === 403, `status=${peek.status}`);

  return { labId, submitted, rows, earlier };
}

try {
  const result = mode === "setup" ? await setup() : await walk();
  const failed = log.filter((l) => !l.ok).length;
  if (outFile) mkdirSync(dirname(outFile), { recursive: true });
  if (outFile) writeFileSync(outFile, JSON.stringify({ base, mode, at: new Date().toISOString(), log, result }, null, 2) + "\n");
  console.log(failed ? `${failed} step(s) failed` : "All steps passed");
  process.exit(failed ? 1 : 0);
} catch (err) {
  console.error(`Stopped: ${err.message}`);
  process.exit(1);
}
