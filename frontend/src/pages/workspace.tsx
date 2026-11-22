import { useRef, useState } from "react";
import { useLocation, useParams } from "react-router-dom";
import { Link } from "@/components/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  BookOpen,
  Check,
  Code2,
  Copy,
  RotateCcw,
  Send,
  Search,
} from "lucide-react";
import { api, type Lab, type Submission, type User } from "@/lib/api";
import {
  Breadcrumb,
  EmptyState,
  Feedback,
  Loading,
  PageHeader,
  Verdict,
} from "@/components/layout";

const starter =
  "#include <iostream>\n\nint main() {\n    // Write your solution here.\n\n    return 0;\n}\n";

function SourceEditor({
  value,
  change,
  readOnly = false,
}: {
  value: string;
  change?: (value: string) => void;
  readOnly?: boolean;
}) {
  const gutter = useRef<HTMLPreElement>(null);

  return (
    <div className="source-editor">
      <pre ref={gutter} className="line-numbers" aria-hidden="true">
        {value
          .split("\n")
          .map((_, i) => i + 1)
          .join("\n")}
      </pre>
      <textarea
        aria-label={readOnly ? "Submitted C++ source" : "C++ source"}
        value={value}
        readOnly={readOnly}
        spellCheck={false}
        autoCapitalize="off"
        autoComplete="off"
        wrap="off"
        onChange={(event) => change?.(event.target.value)}
        onScroll={(event) => {
          if (gutter.current)
            gutter.current.scrollTop = event.currentTarget.scrollTop;
        }}
      />
    </div>
  );
}

function Result({ id }: { id: string }) {
  const query = useQuery({
    queryKey: ["submission", id],
    queryFn: () => api<Submission>(`/submissions/${id}`),
    // Poll only active jobs; failed connections expose a manual retry.
    refetchInterval: (state) =>
      state.state.error ||
      ["completed", "failed"].includes(state.state.data?.status ?? "")
        ? false
        : 1500,
  });

  return (
    <div className="result-content" role="status">
      <Feedback error={query.error} retry={() => void query.refetch()} />
      {query.isPending && <p className="muted">Checking submission…</p>}
      {query.data && (
        <>
          <Verdict status={query.data.status} verdict={query.data.verdict} />
          <p>
            {query.data.detail ||
              "Your submission is in the judging queue. Results will update automatically."}
          </p>
          {query.data.verdict === "Accepted" && (
            <p className="help-text">
              Judge tests passed. This is not a published mark.
            </p>
          )}
          <Link className="text-link" to={`/submissions/${id}`}>
            View submission
          </Link>
        </>
      )}
    </div>
  );
}

export function LabWorkspace({ user }: { user: User }) {
  const { id } = useParams();
  const location = useLocation();
  const query = useQuery({
    queryKey: ["labs"],
    queryFn: () => api<Lab[]>("/labs"),
  });
  const lab = query.data?.find((item) => item.id === id);

  if (query.isPending) return <Loading />;
  if (query.error)
    return (
      <div className="page-container">
        <Feedback error={query.error} retry={() => void query.refetch()} />
      </div>
    );
  if (!lab)
    return (
      <div className="page-container">
        <EmptyState title="Lab not available">
          This lab does not exist or is not assigned to your account.
        </EmptyState>
      </div>
    );

  return (
    <WorkspaceBody
      key={`${user.id}-${lab.id}`}
      lab={lab}
      user={user}
      created={Boolean(location.state?.created)}
    />
  );
}

function WorkspaceBody({
  lab,
  user,
  created,
}: {
  lab: Lab;
  user: User;
  created: boolean;
}) {
  const draftKey = `leetkob-source:${user.id}:${lab.id}`;
  const [source, setSource] = useState(
    () => localStorage.getItem(draftKey) ?? starter,
  );
  const [tab, setTab] = useState("description");
  const [submissionId, setSubmissionId] = useState("");
  const [copied, setCopied] = useState(false);
  const [copyError, setCopyError] = useState<Error | null>(null);
  const [draftError, setDraftError] = useState<Error | null>(null);
  const [split, setSplit] = useState(43);
  const cache = useQueryClient();
  const attempt = useRef({ source: "", key: "" });

  const submit = useMutation({
    mutationFn: () => {
      // Reuse a retry key for identical source after an uncertain network result.
      if (attempt.current.source !== source || !attempt.current.key) {
        attempt.current = { source, key: crypto.randomUUID() };
      }

      return api<{ id: string }>(
        `/labs/${lab.id}/submissions`,
        { source },
        attempt.current.key,
      );
    },
    onSuccess: (result) => {
      setSubmissionId(result.id);
      setTab("result");
      attempt.current = { source: "", key: "" };
      void cache.invalidateQueries({ queryKey: ["submissions"] });
    },
  });

  function updateSource(value: string) {
    setSource(value);

    try {
      localStorage.setItem(draftKey, value);
    } catch {
      setDraftError(
        new Error(
          "Local draft storage is unavailable. Keep a copy before leaving this page.",
        ),
      );
    }
  }

  return (
    <div className="workspace-container">
      <div className="workspace-breadcrumb">
        <Breadcrumb parent="Labs" to="/labs" current={lab.title} />
        <span className="small-tag">C++ laboratory</span>
      </div>
      {created && (
        <div className="success-banner" role="alert">
          <Check size={17} />
          Lab created. Your class students can now submit.
        </div>
      )}
      <div
        className="workspace-mobile-tabs"
        role="tablist"
        aria-label="Lab workspace"
      >
        {["description", "code", "result"].map((item) => (
          <button
            role="tab"
            key={item}
            aria-selected={tab === item}
            className={tab === item ? "active" : ""}
            onClick={() => setTab(item)}
          >
            {item}
          </button>
        ))}
      </div>
      <div
        className="workspace-grid"
        style={{ gridTemplateColumns: `${split}% minmax(0, 1fr)` }}
      >
        <section
          className={`panel description-panel mobile-${tab === "description" ? "visible" : "hidden"}`}
        >
          <div className="panel-title">
            <BookOpen size={16} />
            Description
          </div>
          <div className="problem-content">
            <span className="eyebrow">C++ LAB</span>
            <h1>{lab.title}</h1>
            <div className="problem-tags">
              <span className="small-tag">C++</span>
              <span className="small-tag">Class assignment</span>
            </div>
            <div className="statement">
              {lab.statement || "No statement was provided."}
            </div>
            {!!lab.requirements?.length && (
              <section className="assignment-rules">
                <h2>Assignment requirements</h2>
                <ul>
                  {lab.requirements.map((requirement, index) => (
                    <li key={index}>{requirement}</li>
                  ))}
                </ul>
              </section>
            )}
            <div className="problem-note">
              <Code2 size={20} />
              <p>
                Read the full statement before submitting. Your code is compiled
                and tested in an isolated environment.
              </p>
            </div>
            {user.role === "student" && (
              <p className="help-text">
                Submissions may also undergo assignment-constraint review. Marks
                are released separately by your TA.
              </p>
            )}
          </div>
        </section>
        <div
          className={`coding-column mobile-${tab !== "description" ? "visible" : "hidden"}`}
        >
          <section
            className={`panel editor-panel mobile-${tab === "code" ? "visible" : "hidden"}`}
          >
            <div className="panel-title">
              <Code2 size={16} />
              Code<span className="editor-language">C++</span>
              <button
                className="icon-button"
                aria-label="Copy source"
                title="Copy source"
                onClick={() => {
                  void navigator.clipboard
                    .writeText(source)
                    .then(() => {
                      setCopied(true);
                      window.setTimeout(() => setCopied(false), 1800);
                    })
                    .catch(() =>
                      setCopyError(
                        new Error(
                          "Could not copy source. Select and copy it manually.",
                        ),
                      ),
                    );
                }}
              >
                {copied ? <Check size={15} /> : <Copy size={15} />}
              </button>
              <button
                className="icon-button"
                aria-label="Reset source"
                title="Reset source"
                onClick={() => {
                  if (
                    window.confirm(
                      "Reset your code to the starter? Your saved local draft will be replaced.",
                    )
                  )
                    updateSource(starter);
                }}
              >
                <RotateCcw size={15} />
              </button>
            </div>
            <SourceEditor value={source} change={updateSource} />
            <Feedback error={draftError || copyError} />
            <div className="editor-footer">
              <span>
                {draftError ? "Draft not saved" : "Draft stored on this device"}
              </span>
              {user.role === "student" && (
                <button
                  className="btn btn-submit"
                  disabled={!source.trim() || submit.isPending}
                  onClick={() => submit.mutate()}
                >
                  <Send size={15} />
                  {submit.isPending ? "Submitting…" : "Submit"}
                </button>
              )}
            </div>
          </section>
          <section
            className={`panel result-panel mobile-${tab === "result" ? "visible" : "hidden"}`}
          >
            <div className="panel-title">
              <Check size={16} />
              Submission result
            </div>
            <Feedback error={submit.error} retry={() => submit.mutate()} />
            {submissionId ? (
              <Result id={submissionId} />
            ) : (
              <div className="result-placeholder">
                <Code2 size={24} />
                <p>
                  {user.role === "student"
                    ? "Your result will appear here after you submit."
                    : "Students can submit from their own accounts."}
                </p>
              </div>
            )}
          </section>
        </div>
      </div>
      <label className="resize-control">
        Panel width
        <input
          type="range"
          min="30"
          max="65"
          value={split}
          onChange={(event) => setSplit(Number(event.target.value))}
        />
      </label>
    </div>
  );
}

export function SubmissionsPage() {
  const [search, setSearch] = useState("");
  const query = useQuery({
    queryKey: ["submissions"],
    queryFn: () => api<Submission[]>("/submissions"),
  });
  const filtered = query.data?.filter((item) =>
    item.title?.toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <div className="page-container">
      <PageHeader
        eyebrow="YOUR JOURNEY"
        title="My submissions"
        description="Every attempt is a step forward. Revisit your solutions and results."
      />
      <div className="section-toolbar">
        <label className="search-field">
          <Search size={17} />
          <input
            aria-label="Search submissions"
            placeholder="Search by lab…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </label>
        <button
          className="btn"
          onClick={() => void query.refetch()}
          disabled={query.isFetching}
        >
          Refresh
        </button>
      </div>
      <Feedback error={query.error} retry={() => void query.refetch()} />
      {query.isPending ? (
        <Loading />
      ) : (
        <div className="panel table-wrap">
          <table>
            <thead>
              <tr>
                <th>Lab</th>
                <th>Result</th>
                <th>Submitted</th>
                <th>Language</th>
              </tr>
            </thead>
            <tbody>
              {filtered?.map((item) => (
                <tr key={item.id}>
                  <td>
                    <Link className="lab-title" to={`/submissions/${item.id}`}>
                      {item.title}
                    </Link>
                  </td>
                  <td>
                    <Verdict status={item.status} verdict={item.verdict} />
                  </td>
                  <td className="muted">
                    {item.createdAt
                      ? new Date(item.createdAt).toLocaleString()
                      : "—"}
                  </td>
                  <td>C++</td>
                </tr>
              ))}
            </tbody>
          </table>
          {!filtered?.length && !query.error && (
            <EmptyState
              title="Your solutions will live here"
              action={
                <Link className="btn btn-primary" to="/labs">
                  Explore labs
                </Link>
              }
            >
              Submit a solution to start your history.
            </EmptyState>
          )}
        </div>
      )}
      <p className="help-text">Showing your latest 200 submissions.</p>
    </div>
  );
}

export function SubmissionPage() {
  const { id } = useParams();
  const query = useQuery({
    queryKey: ["submission", id],
    queryFn: () => api<Submission>(`/submissions/${id}`),
    refetchInterval: (state) =>
      state.state.error ||
      ["completed", "failed"].includes(state.state.data?.status ?? "")
        ? false
        : 1500,
  });

  return (
    <div className="page-container">
      <Breadcrumb parent="Labs" to="/labs" current="Submission" />
      <PageHeader
        title="Submission details"
        description="Your saved solution and its judge result."
      />
      <Feedback error={query.error} retry={() => void query.refetch()} />
      {query.isPending ? (
        <Loading />
      ) : (
        query.data && (
          <>
            <section className="panel form-section" role="status">
              <Verdict
                status={query.data.status}
                verdict={query.data.verdict}
              />
              <p>
                {query.data.detail ||
                  "Judging is in progress. This page updates automatically."}
              </p>
              <p className="help-text">
                Judge results and published marks are separate.
              </p>
              <Link className="text-link" to={`/labs/${query.data.labId}`}>
                Return to lab
              </Link>
            </section>
            <section className="panel saved-source">
              <div className="panel-title">
                <Code2 size={16} />
                Submitted source<span className="editor-language">C++</span>
              </div>
              <SourceEditor value={query.data.source ?? ""} readOnly />
            </section>
          </>
        )
      )}
    </div>
  );
}
