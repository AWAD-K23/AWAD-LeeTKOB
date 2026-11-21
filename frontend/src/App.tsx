import { useEffect, useState } from "react";
import "./style.css";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";

type Lab = { id: string; title: string; statement: string };
type Status = { id: string; status: string; verdict: string; detail: string };

async function api<T>(path: string, body?: object, headers = {}): Promise<T> {
  const response = await fetch(path, {
    method: body ? "POST" : "GET",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json", ...headers },
    body: body ? JSON.stringify(body) : undefined,
  });
  if (!response.ok) throw new Error((await response.json()).detail);
  return response.status === 204 ? (undefined as T) : response.json();
}

export default function App() {
  const [role, setRole] = useState("");
  const [classes, setClasses] = useState<{ id: string; name: string }[]>([]);
  const [labs, setLabs] = useState<Lab[]>([]);
  const [selected, setSelected] = useState("");
  const [message, setMessage] = useState("");
  const [status, setStatus] = useState<Status | null>(null);
  const [busy, setBusy] = useState(false);
  const [tests, setTests] = useState([{ input: "", expectedOutput: "" }]);
  const [constraints, setConstraints] = useState([
    { description: "", deduction: 0 },
  ]);

  async function refresh() {
    setRole((await api<{ role: string }>("/me")).role);
    setClasses(await api("/classes"));
    setLabs(await api("/labs"));
  }
  useEffect(() => {
    let active = true;
    void Promise.all([
      api<{ role: string }>("/me"),
      api<{ id: string; name: string }[]>("/classes"),
      api<Lab[]>("/labs"),
    ])
      .then(([me, availableClasses, availableLabs]) => {
        if (!active) return;
        setRole(me.role);
        setClasses(availableClasses);
        setLabs(availableLabs);
      })
      .catch(() => {});
    return () => {
      active = false;
    };
  }, []);
  useEffect(() => {
    if (!status || ["completed", "failed"].includes(status.status)) return;
    const timer = setInterval(() => {
      void api<Status>(`/submissions/${status.id}`)
        .then(setStatus)
        .catch((error: Error) => setMessage(error.message));
    }, 1500);
    return () => clearInterval(timer);
  }, [status]);

  async function run(action: () => Promise<void>) {
    setBusy(true);
    setMessage("");
    try {
      await action();
    } catch (error) {
      setMessage((error as Error).message);
    } finally {
      setBusy(false);
    }
  }

  return (
    <main>
      <header>
        <h1>LeeTKOB</h1>
        <span>C++ laboratory</span>
        {role && (
          <Button
            onClick={() =>
              void run(async () => {
                await api("/auth/logout", {});
                setRole("");
                setStatus(null);
              })
            }
          >
            Sign out
          </Button>
        )}
      </header>
      {message && (
        <p role="alert" className="notice">
          {message}
        </p>
      )}
      {!role ? (
        <Card className="mx-auto my-12 max-w-md p-6">
          <h2>Welcome back</h2>
          <p>Sign in to manage labs or submit your work.</p>
          <form
            onSubmit={(event) => {
              event.preventDefault();
              const data = new FormData(event.currentTarget);
              void run(async () => {
                await api("/auth/login", {
                  email: data.get("email"),
                  password: data.get("password"),
                });
                await refresh();
              });
            }}
          >
            <div className="grid gap-2">
              <Label htmlFor={"field-1"}>Email</Label>
              <Input
                id={"field-1"}
                name="email"
                type="email"
                required
                autoComplete="username"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor={"field-2"}>Password</Label>
              <Input
                id={"field-2"}
                name="password"
                type="password"
                required
                autoComplete="current-password"
              />
            </div>
            <Button disabled={busy}>Sign in</Button>
          </form>
        </Card>
      ) : (
        <>
          <p className="mb-4">
            Signed in as <Badge variant="secondary">{role}</Badge>
          </p>
          {role === "admin" && (
            <Card className="mb-6 p-6">
              <h2>Create a class</h2>
              <form
                onSubmit={(event) => {
                  event.preventDefault();
                  const data = new FormData(event.currentTarget);
                  void run(async () => {
                    await api("/classes", {
                      name: data.get("name"),
                      taId: data.get("taId"),
                      studentIds: String(data.get("studentIds"))
                        .split(/[\s,]+/)
                        .filter(Boolean),
                    });
                    await refresh();
                    setMessage("Class created.");
                  });
                }}
              >
                <div className="grid gap-2">
                  <Label htmlFor={"field-3"}>Class name</Label>
                  <Input id={"field-3"} name="name" required />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor={"field-4"}>TA account ID</Label>
                  <Input id={"field-4"} name="taId" required />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor={"field-5"}>
                    Student account IDs (one per line)
                  </Label>
                  <Textarea id={"field-5"} name="studentIds" />
                </div>
                <Button disabled={busy}>Create class</Button>
              </form>
            </Card>
          )}
          {role === "ta" && (
            <Card className="mb-6 p-6">
              <h2>Create a lab</h2>
              <form
                onSubmit={(event) => {
                  event.preventDefault();
                  const data = new FormData(event.currentTarget);
                  void run(async () => {
                    await api(`/classes/${data.get("classId")}/labs`, {
                      title: data.get("title"),
                      statement: data.get("statement"),
                      tests,
                      constraints,
                    });
                    await refresh();
                    setMessage(
                      "Lab created. Your class students can now submit.",
                    );
                  });
                }}
              >
                <div className="grid gap-2">
                  <Label htmlFor={"field-6"}>Class</Label>
                  <select
                    id={"field-6"}
                    name="classId"
                    required
                    defaultValue=""
                  >
                    <option value="" disabled>
                      Choose a class
                    </option>
                    {classes.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor={"field-7"}>Title</Label>
                  <Input id={"field-7"} name="title" required />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor={"field-8"}>Statement</Label>
                  <Textarea id={"field-8"} name="statement" required rows={6} />
                </div>
                <h3>Test cases</h3>
                {tests.map((test, index) => (
                  <fieldset key={index}>
                    <legend>Case {index + 1}</legend>
                    <div className="grid gap-2">
                      <Label htmlFor={`field-9-${index}`}>Input</Label>
                      <Textarea
                        id={`field-9-${index}`}
                        value={test.input}
                        onChange={(e) =>
                          setTests(
                            tests.map((t, i) =>
                              i === index ? { ...t, input: e.target.value } : t,
                            ),
                          )
                        }
                      />
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor={`field-10-${index}`}>
                        Expected output
                      </Label>
                      <Textarea
                        id={`field-10-${index}`}
                        value={test.expectedOutput}
                        onChange={(e) =>
                          setTests(
                            tests.map((t, i) =>
                              i === index
                                ? { ...t, expectedOutput: e.target.value }
                                : t,
                            ),
                          )
                        }
                      />
                    </div>
                    {tests.length > 1 && (
                      <Button
                        type="button"
                        onClick={() =>
                          setTests(tests.filter((_, i) => i !== index))
                        }
                      >
                        Remove case
                      </Button>
                    )}
                  </fieldset>
                ))}
                <Button
                  type="button"
                  onClick={() =>
                    setTests([...tests, { input: "", expectedOutput: "" }])
                  }
                >
                  Add case
                </Button>
                <h3>Constraints</h3>
                {constraints.map((c, index) => (
                  <fieldset key={index}>
                    <legend>Constraint {index + 1}</legend>
                    <div className="grid gap-2">
                      <Label htmlFor={`field-11-${index}`}>Rule</Label>
                      <Textarea
                        id={`field-11-${index}`}
                        value={c.description}
                        required
                        onChange={(e) =>
                          setConstraints(
                            constraints.map((v, i) =>
                              i === index
                                ? { ...v, description: e.target.value }
                                : v,
                            ),
                          )
                        }
                      />
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor={`field-12-${index}`}>Deduction</Label>
                      <Input
                        id={`field-12-${index}`}
                        type="number"
                        min="0"
                        step="0.01"
                        value={c.deduction}
                        required
                        onChange={(e) =>
                          setConstraints(
                            constraints.map((v, i) =>
                              i === index
                                ? { ...v, deduction: Number(e.target.value) }
                                : v,
                            ),
                          )
                        }
                      />
                    </div>
                    {constraints.length > 1 && (
                      <Button
                        type="button"
                        onClick={() =>
                          setConstraints(
                            constraints.filter((_, i) => i !== index),
                          )
                        }
                      >
                        Remove constraint
                      </Button>
                    )}
                  </fieldset>
                ))}
                <Button
                  type="button"
                  onClick={() =>
                    setConstraints([
                      ...constraints,
                      { description: "", deduction: 0 },
                    ])
                  }
                >
                  Add constraint
                </Button>
                <Button
                  className="bg-emerald-500 text-black hover:bg-emerald-400"
                  disabled={busy || !classes.length}
                >
                  Create lab
                </Button>
              </form>
            </Card>
          )}
          {role !== "admin" && (
            <Card className="mb-6 p-6">
              <h2>Your labs</h2>
              {!labs.length && <p>No labs are available yet.</p>}
              <div className="grid gap-2">
                <Label htmlFor={"field-13"}>Lab</Label>
                <select
                  id={"field-13"}
                  value={selected}
                  onChange={(e) => setSelected(e.target.value)}
                >
                  <option value="">Choose a lab</option>
                  {labs.map((l) => (
                    <option key={l.id} value={l.id}>
                      {l.title}
                    </option>
                  ))}
                </select>
              </div>
              {selected && (
                <article>
                  <h3>{labs.find((l) => l.id === selected)?.title}</h3>
                  <p className="statement">
                    {labs.find((l) => l.id === selected)?.statement}
                  </p>
                </article>
              )}
              {role === "student" && selected && (
                <form
                  onSubmit={(event) => {
                    event.preventDefault();
                    const data = new FormData(event.currentTarget);
                    void run(async () => {
                      const result = await api<{ id: string }>(
                        `/labs/${selected}/submissions`,
                        { source: data.get("source") },
                        { "Idempotency-Key": crypto.randomUUID() },
                      );
                      setStatus({
                        id: result.id,
                        status: "pending",
                        verdict: "",
                        detail: "",
                      });
                    });
                  }}
                >
                  <div className="grid gap-2">
                    <Label htmlFor={"field-14"}>C++ source</Label>
                    <Textarea
                      id={"field-14"}
                      className="code"
                      name="source"
                      required
                      rows={18}
                      spellCheck={false}
                    />
                  </div>
                  <Button
                    className="bg-emerald-500 text-black hover:bg-emerald-400"
                    disabled={busy}
                  >
                    Submit
                  </Button>
                </form>
              )}
              {status && (
                <aside role="status">
                  <h3>
                    {status.verdict ||
                      (status.status === "pending" ? "Queued" : status.status)}
                  </h3>
                  <p>{status.detail}</p>
                  <small>Submission {status.id}</small>
                </aside>
              )}
            </Card>
          )}
        </>
      )}
    </main>
  );
}
