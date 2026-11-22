import { useState } from "react";
import { LabMarkdown } from "@/components/markdown";
import { useUnsavedChanges } from "@/lib/use-unsaved-changes";
import { useNavigate, useParams } from "react-router-dom";
import { Link } from "@/components/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { BookOpen, Code2, Plus, Search, Trash2 } from "lucide-react";
import { api, type Classroom, type Lab, type User } from "@/lib/api";
import {
  Breadcrumb,
  EmptyState,
  Feedback,
  Loading,
  PageHeader,
} from "@/components/layout";

export function LabsPage({ user }: { user: User }) {
  const [search, setSearch] = useState("");
  const [classId, setClassId] = useState("");
  const query = useQuery({
    queryKey: ["labs"],
    queryFn: () => api<Lab[]>("/labs"),
  });
  const classes = useQuery({
    queryKey: ["classes"],
    queryFn: () => api<Classroom[]>("/classes"),
  });
  const filtered = query.data?.filter(
    (lab) =>
      lab.title.toLowerCase().includes(search.toLowerCase()) &&
      (!classId || lab.classId === classId),
  );

  return (
    <div className="page-container">
      <PageHeader
        eyebrow="PRACTICE WITH PURPOSE"
        title={user.role === "ta" ? "Your labs" : "Explore your labs"}
        description={
          user.role === "ta"
            ? "Shape the next challenge for your class."
            : "Every solution starts with a little curiosity."
        }
        action={
          user.role === "ta" && (
            <Link className="btn btn-primary" to="/labs/new">
              <Plus size={16} />
              Create lab
            </Link>
          )
        }
      />
      <div className="lab-banner">
        <div className="banner-icon">
          <Code2 size={32} />
        </div>
        <div>
          <span className="eyebrow">C++ LABORATORY</span>
          <h2>
            {user.role === "ta"
              ? "Turn ideas into challenges."
              : "Think it through. Code it out."}
          </h2>
          <p>
            {user.role === "ta"
              ? "Define the task. Write the tests. Help your students grow."
              : "Read the problem, build your solution, and learn from the result."}
          </p>
        </div>
        <span className="banner-braces" aria-hidden="true">
          {"{ }"}
        </span>
      </div>
      <div className="section-toolbar">
        <label className="search-field">
          <Search size={17} />
          <input
            aria-label="Search labs"
            value={search}
            placeholder="Search labs…"
            onChange={(e) => setSearch(e.target.value)}
          />
        </label>
        <select
          aria-label="Filter by class"
          value={classId}
          onChange={(e) => setClassId(e.target.value)}
        >
          <option value="">All classes</option>
          {classes.data?.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>
        <span className="muted">{filtered?.length ?? 0} labs</span>
      </div>
      <Feedback
        error={query.error || classes.error}
        retry={() => {
          void query.refetch();
          void classes.refetch();
        }}
      />
      {query.isPending ? (
        <Loading />
      ) : (
        <div className="panel table-wrap">
          <table>
            <thead>
              <tr>
                <th>Lab</th>
                <th>Class</th>
                <th>Language</th>
                <th>
                  <span className="sr-only">Action</span>
                </th>
              </tr>
            </thead>
            <tbody>
              {filtered?.map((lab) => (
                <tr key={lab.id}>
                  <td>
                    <Link className="lab-title" to={`/labs/${lab.id}`}>
                      <BookOpen size={17} />
                      {lab.title}
                    </Link>
                  </td>
                  <td className="muted">
                    {classes.data?.find((item) => item.id === lab.classId)
                      ?.name || "Assigned lab"}
                  </td>
                  <td>
                    <span className="small-tag">C++</span>
                  </td>
                  <td>
                    <Link className="text-link" to={`/labs/${lab.id}`}>
                      Open lab
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {!filtered?.length && !query.error && (
            <EmptyState
              title={
                search
                  ? "No matching labs"
                  : "Your next challenge is on its way"
              }
            >
              {search
                ? "Try another search or class filter."
                : user.role === "ta"
                  ? "Create a lab for an assigned class to get started."
                  : "Labs appear here when your TA creates them for your class."}
            </EmptyState>
          )}
        </div>
      )}
    </div>
  );
}

type TestCase = { input: string; expectedOutput: string; key: string };
type Constraint = { description: string; deduction: number; key: string };

export function LabForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const cache = useQueryClient();
  const classes = useQuery({
    queryKey: ["classes"],
    queryFn: () => api<Classroom[]>("/classes"),
  });
  const [tests, setTests] = useState<TestCase[]>([
    { input: "", expectedOutput: "", key: crypto.randomUUID() },
  ]);
  const [constraints, setConstraints] = useState<Constraint[]>([
    { description: "", deduction: 0, key: crypto.randomUUID() },
  ]);
  const [dirty, setDirty] = useState(false);
  const [preview, setPreview] = useState(false);
  const [statement, setStatement] = useState("");
  const [title, setTitle] = useState("");

  const allowNavigation = useUnsavedChanges(dirty);

  const save = useMutation({
    mutationFn: ({ classId, title }: { classId: string; title: string }) =>
      api<{ id: string }>(`/classes/${classId}/labs`, {
        title,
        statement,
        tests: tests.map(({ input, expectedOutput }) => ({
          input,
          expectedOutput,
        })),
        constraints: constraints.map(({ description, deduction }) => ({
          description,
          deduction,
        })),
      }),
    onSuccess: async (result) => {
      allowNavigation();
      setDirty(false);
      await cache.invalidateQueries({ queryKey: ["labs"] });
      navigate(`/labs/${result.id}`, {
        state: { created: true },
        viewTransition: true,
      });
    },
  });

  return (
    <div className="page-container authoring-page">
      <Breadcrumb parent="Labs" to="/labs" current="New lab" />
      <PageHeader
        title="Create a lab"
        description="A thoughtful challenge starts with a clear brief."
      />
      <Feedback
        error={save.error || classes.error}
        retry={classes.error ? () => void classes.refetch() : undefined}
      />
      <div className="authoring-layout">
        <aside className="section-nav" aria-label="Form sections">
          <a href="#details">Details</a>
          <a href="#test-cases">
            Test cases <span>{tests.length}</span>
          </a>
          <a href="#constraints">
            Constraints <span>{constraints.length}</span>
          </a>
          <a href="#review-lab">Review</a>
          <div className="authoring-note">
            <BookOpen size={20} />
            <p>Clear requirements lead to better solutions.</p>
          </div>
        </aside>
        <form
          className="form-layout"
          onChange={() => setDirty(true)}
          onSubmit={(event) => {
            event.preventDefault();
            const data = new FormData(event.currentTarget);
            save.mutate({
              classId: String(data.get("classId")),
              title: String(data.get("title")),
            });
          }}
        >
          <section id="details" className="panel form-section">
            <h2>Lab details</h2>
            <p className="muted">
              Give students the context they need to begin.
            </p>
            <label>
              Class
              <select
                aria-label="Class"
                name="classId"
                required
                defaultValue={id || ""}
              >
                <option value="" disabled>
                  Choose a class
                </option>
                {classes.data?.map((item) => (
                  <option key={item.id} value={item.id}>
                    {item.name}
                  </option>
                ))}
              </select>
            </label>
            {!classes.isPending && !classes.data?.length && (
              <p className="help-text">
                Ask an administrator to assign you to a class first.
              </p>
            )}
            <label>
              Title
              <input
                name="title"
                placeholder="e.g. Implement merge sort"
                required
                maxLength={200}
                value={title}
                onChange={(event) => setTitle(event.target.value)}
              />
            </label>
            <div className="inline-heading">
              <label htmlFor="statement">Statement</label>
              <button
                className="btn btn-small"
                type="button"
                onClick={() => setPreview(!preview)}
              >
                {preview ? "Edit" : "Preview"}
              </button>
            </div>
            {preview && (
              <div className="statement-preview">
                <LabMarkdown>
                  {statement || "Your statement preview will appear here."}
                </LabMarkdown>
              </div>
            )}
            <textarea
              id="statement"
              name="statement"
              rows={8}
              required
              value={statement}
              className={preview ? "visually-hidden-input" : ""}
              placeholder="Describe the task, input format, output format, and examples…"
              onChange={(e) => setStatement(e.target.value)}
            />
            <p className="help-text">
              Markdown supported: headings, lists, tables, links, and fenced
              code blocks.
            </p>
          </section>
          <section id="test-cases" className="panel form-section">
            <h2>Test cases</h2>
            <p className="muted">
              Private inputs and expected outputs used by the judge.
            </p>
            {tests.map((item, index) => (
              <fieldset className="repeat-section" key={item.key}>
                <legend>Case {index + 1}</legend>
                <div className="two-columns">
                  <label>
                    Input
                    <textarea
                      rows={4}
                      className="mono"
                      value={item.input}
                      onChange={(e) =>
                        setTests(
                          tests.map((test) =>
                            test.key === item.key
                              ? { ...test, input: e.target.value }
                              : test,
                          ),
                        )
                      }
                    />
                  </label>
                  <label>
                    Expected output
                    <textarea
                      rows={4}
                      className="mono"
                      value={item.expectedOutput}
                      onChange={(e) =>
                        setTests(
                          tests.map((test) =>
                            test.key === item.key
                              ? { ...test, expectedOutput: e.target.value }
                              : test,
                          ),
                        )
                      }
                    />
                  </label>
                </div>
                {tests.length > 1 && (
                  <button
                    className="btn btn-small"
                    type="button"
                    onClick={() => {
                      setTests(tests.filter((test) => test.key !== item.key));
                      setDirty(true);
                    }}
                  >
                    <Trash2 size={14} />
                    Remove case
                  </button>
                )}
              </fieldset>
            ))}
            <button
              className="btn"
              type="button"
              onClick={() => {
                setTests([
                  ...tests,
                  { input: "", expectedOutput: "", key: crypto.randomUUID() },
                ]);
                setDirty(true);
              }}
            >
              <Plus size={16} />
              Add case
            </button>
          </section>
          <section id="constraints" className="panel form-section">
            <h2>Assignment constraints</h2>
            <p className="muted">
              You choose deductions. Checker findings are reviewed separately.
            </p>
            {constraints.map((item, index) => (
              <fieldset className="repeat-section" key={item.key}>
                <legend>Constraint {index + 1}</legend>
                <label>
                  Rule
                  <textarea
                    value={item.description}
                    required
                    rows={2}
                    placeholder="e.g. Do not call `std::sort`"
                    onChange={(e) =>
                      setConstraints(
                        constraints.map((rule) =>
                          rule.key === item.key
                            ? { ...rule, description: e.target.value }
                            : rule,
                        ),
                      )
                    }
                  />
                </label>
                <label>
                  Deduction
                  <input
                    type="number"
                    min="0"
                    step="0.01"
                    required
                    value={item.deduction}
                    onChange={(e) =>
                      setConstraints(
                        constraints.map((rule) =>
                          rule.key === item.key
                            ? { ...rule, deduction: Number(e.target.value) }
                            : rule,
                        ),
                      )
                    }
                  />
                </label>
                {constraints.length > 1 && (
                  <button
                    className="btn btn-small"
                    type="button"
                    onClick={() => {
                      setConstraints(
                        constraints.filter((rule) => rule.key !== item.key),
                      );
                      setDirty(true);
                    }}
                  >
                    <Trash2 size={14} />
                    Remove constraint
                  </button>
                )}
              </fieldset>
            ))}
            <button
              className="btn"
              type="button"
              onClick={() => {
                setConstraints([
                  ...constraints,
                  { description: "", deduction: 0, key: crypto.randomUUID() },
                ]);
                setDirty(true);
              }}
            >
              <Plus size={16} />
              Add constraint
            </button>
          </section>
          <section id="review-lab" className="panel form-section">
            <h2>Review your lab</h2>
            <p className="muted">
              Check the statement, private tests, and deductions before
              creating.
            </p>
            <div className="statement-preview">
              <h3>{title || "Untitled lab"}</h3>
              <LabMarkdown>
                {statement || "Your statement will appear here."}
              </LabMarkdown>
            </div>
            <p>
              {tests.length} private test cases · {constraints.length}{" "}
              assignment constraints
            </p>
            <ul className="review-rules">
              {constraints.map((rule) => (
                <li key={rule.key}>
                  <span>{rule.description || "Add a requirement above"}</span>
                  <span className="small-tag">Deduction: {rule.deduction}</span>
                </li>
              ))}
            </ul>
          </section>
          <div className="form-footer sticky-footer">
            <span className="muted">
              {tests.length} tests · {constraints.length} constraints
            </span>
            <button
              className="btn"
              type="button"
              onClick={() => {
                navigate("/labs");
              }}
            >
              Cancel
            </button>
            <button
              className="btn btn-primary"
              disabled={save.isPending || !classes.data?.length}
            >
              {save.isPending ? "Creating…" : "Create lab"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
