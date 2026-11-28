import { useState } from "react";
import { useParams } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import {
  Breadcrumb,
  EmptyState,
  Feedback,
  Loading,
  PageHeader,
} from "@/components/layout";

type Evidence = { line: number; text: string };
type Finding = {
  submissionId: string;
  studentId: string;
  createdAt: string;
  judgeVerdict: string;
  sourceCode: string;
  constraintId: string;
  constraint: string;
  deduction: number;
  verdict: "pass" | "violation" | "unsure";
  reason: string;
  evidence: Evidence[];
  source: string;
  decision: "" | "pass" | "violation";
  reviewEvidence: Evidence[];
  appliedDeduction: number;
  reviewRevision: number;
  reviewRequired: boolean;
  sampledPass: boolean;
  reviewReason: string;
};

type EarlierSubmission = {
  submissionId: string;
  studentId: string;
  createdAt: string;
  judgeVerdict: string;
  sourceCode: string;
  findings: {
    constraint: string;
    verdict: "pass" | "violation" | "unsure";
    source: string;
    reason: string;
    evidence: Evidence[];
  }[];
};

// Only each student's latest submission counts (FAQ-02); the TA can still open
// the earlier ones (FAQ-10). They are shown read-only: no decision applies.
function EarlierSubmissions({ items }: { items: EarlierSubmission[] }) {
  return (
    <section className="panel form-section" aria-labelledby="earlier-heading">
      <h2 id="earlier-heading">Earlier submissions (not counted)</h2>
      <p className="muted">
        Only each student&apos;s latest submission counts and enters the queue.
        These earlier ones are shown for reference; they cannot be decided.
      </p>
      {items.map((item) => (
        <details key={item.submissionId}>
          <summary>
            Submission {item.submissionId} · Student {item.studentId} ·{" "}
            {item.createdAt} · {item.judgeVerdict || "Not judged"}
          </summary>
          <pre className="mono">{item.sourceCode}</pre>
          <ul className="review-rules">
            {item.findings.map((finding, index) => (
              <li key={index}>
                {finding.constraint}:{" "}
                <strong>
                  {finding.verdict} ({finding.source})
                </strong>
                {finding.verdict === "unsure" && finding.reason
                  ? ` · ${finding.reason}`
                  : ""}
                {finding.evidence.map((line) => (
                  <div key={line.line}>
                    Line {line.line}: <code>{line.text}</code>
                  </div>
                ))}
              </li>
            ))}
          </ul>
        </details>
      ))}
    </section>
  );
}

function FindingDecision({
  labId,
  finding,
}: {
  labId: string;
  finding: Finding;
}) {
  const cache = useQueryClient();
  const [line, setLine] = useState("");
  const decide = useMutation({
    mutationFn: (decision: "pass" | "violation") => {
      const number = Number(line);
      const codeLine = finding.sourceCode.split("\n")[number - 1];
      return api(
        `/labs/${labId}/submissions/${finding.submissionId}/constraints/${finding.constraintId}/review`,
        {
          decision,
          ...(decision === "violation"
            ? { evidence: [{ line: number, text: codeLine ?? "" }] }
            : {}),
        },
      );
    },
    onSuccess: () => {
      void cache.invalidateQueries({ queryKey: ["review-queue", labId] });
      void cache.invalidateQueries({ queryKey: ["publication", labId] });
    },
  });
  const maxLine = finding.sourceCode.split("\n").length;

  return (
    <li className="panel form-section">
      <h3>{finding.constraint}</h3>
      <p>
        {finding.sampledPass
          ? "Sampled pass - review required"
          : finding.reviewRequired
            ? "Review required"
            : "Optional review"}
      </p>
      <p>
        Checker: <strong>{finding.verdict}</strong> ({finding.source}) ·
        Configured deduction: <strong>{finding.deduction}</strong>
      </p>
      <p>
        TA decision: <strong>{finding.decision || "Pending review"}</strong>.
        TA-selected deduction: <strong>{finding.appliedDeduction}</strong>
      </p>
      {finding.verdict === "unsure" && (
        <p role="status">
          Why unsure: {finding.reason || "Reason not recorded."}
        </p>
      )}
      {finding.evidence.length > 0 && (
        <ul>
          {finding.evidence.map((item, index) => (
            <li key={index}>
              Suggested line {item.line}: <code>{item.text}</code>
            </li>
          ))}
        </ul>
      )}
      {finding.reviewEvidence.length > 0 && (
        <ul aria-label="TA-confirmed evidence">
          {finding.reviewEvidence.map((item, index) => (
            <li key={index}>
              Confirmed line {item.line}: <code>{item.text}</code>
            </li>
          ))}
        </ul>
      )}
      <Feedback error={decide.error} />
      <label>
        Evidence line for a violation
        <input
          type="number"
          min={1}
          max={maxLine}
          value={line}
          onChange={(event) => setLine(event.target.value)}
        />
      </label>
      <div className="form-footer">
        <button
          className="btn"
          disabled={decide.isPending}
          onClick={() => decide.mutate("pass")}
        >
          Confirm pass
        </button>
        <button
          className="btn btn-primary"
          disabled={
            decide.isPending ||
            !Number.isInteger(Number(line)) ||
            Number(line) < 1 ||
            Number(line) > maxLine
          }
          onClick={() => decide.mutate("violation")}
        >
          Confirm violation
        </button>
      </div>
    </li>
  );
}

export function ReviewQueuePage() {
  const { id = "" } = useParams();
  const cache = useQueryClient();
  const publication = useQuery({
    queryKey: ["publication", id],
    queryFn: () =>
      api<{
        reviewRevision: number;
        required: number;
        reviewed: number;
        pending: number;
        incompleteSubmissions: number;
        published: boolean;
      }>(`/labs/${id}/publication`),
    refetchInterval: 5000,
  });
  const publish = useMutation({
    mutationFn: () =>
      api(`/labs/${id}/publish`, {
        reviewRevision: publication.data?.reviewRevision,
      }),
    onSettled: () => {
      void cache.invalidateQueries({ queryKey: ["publication", id] });
      void cache.invalidateQueries({ queryKey: ["review-queue", id] });
    },
  });
  const queue = useQuery({
    queryKey: ["review-queue", id],
    queryFn: () => api<Finding[]>(`/labs/${id}/review-queue`),
    refetchInterval: 5000,
  });
  const earlier = useQuery({
    queryKey: ["earlier-submissions", id],
    queryFn: () => api<EarlierSubmission[]>(`/labs/${id}/earlier-submissions`),
    refetchInterval: 5000,
  });
  if (queue.isPending) return <Loading />;
  const submissions = (queue.data ?? []).reduce<Record<string, Finding[]>>(
    (groups, item) => {
      (groups[item.submissionId] ??= []).push(item);
      return groups;
    },
    {},
  );

  return (
    <div className="page-container">
      <Breadcrumb parent="Labs" to="/labs" current="TA review" />
      <PageHeader
        title="TA review queue"
        description="Review the required findings, then Publish to release deductions to students."
      />
      <section className="panel form-section" aria-label="Publication">
        <h2>Publication</h2>
        <Feedback
          error={publication.error}
          retry={() => void publication.refetch()}
        />
        {publication.isPending && <Loading />}
        {publication.data && !publication.error && (
          <>
            <p role="status">
              {publication.data.reviewed} / {publication.data.required} required
              reviews complete; {publication.data.pending} pending;{" "}
              {publication.data.incompleteSubmissions} submissions awaiting
              complete judging.
            </p>
            <p>
              {publication.data.published
                ? "Results published."
                : "Results are not published."}
            </p>
          </>
        )}
        <Feedback error={publish.error} />
        <button
          className="btn btn-primary"
          disabled={
            !publication.data ||
            !!publication.error ||
            publication.isFetching ||
            publish.isPending ||
            publication.data.published ||
            publication.data.pending > 0 ||
            publication.data.incompleteSubmissions > 0
          }
          onClick={() => publish.mutate()}
        >
          Publish results
        </button>
      </section>
      <Feedback error={queue.error} retry={() => void queue.refetch()} />
      {!queue.error && !queue.data?.length && (
        <EmptyState title="No counted Accepted submissions yet">
          The latest Accepted submission for each student appears after judging.
        </EmptyState>
      )}
      {Object.entries(submissions).map(([submissionId, findings]) => (
        <section className="panel form-section" key={submissionId}>
          <h2>Submission {submissionId}</h2>
          <p className="muted">
            Student {findings?.[0]?.studentId} · {findings?.[0]?.createdAt} ·{" "}
            {findings?.[0]?.judgeVerdict}
          </p>
          <pre className="mono">{findings?.[0]?.sourceCode}</pre>
          <ul className="review-rules">
            {findings?.map((finding) => (
              <FindingDecision
                key={finding.constraintId}
                labId={id}
                finding={finding}
              />
            ))}
          </ul>
        </section>
      ))}
      <Feedback error={earlier.error} retry={() => void earlier.refetch()} />
      {!!earlier.data?.length && <EarlierSubmissions items={earlier.data} />}
    </div>
  );
}
