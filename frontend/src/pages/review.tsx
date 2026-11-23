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
};

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
    },
  });
  const maxLine = finding.sourceCode.split("\n").length;

  return (
    <li className="panel form-section">
      <h3>{finding.constraint}</h3>
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
  const queue = useQuery({
    queryKey: ["review-queue", id],
    queryFn: () => api<Finding[]>(`/labs/${id}/review-queue`),
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
        description="Findings stay private until a TA decides. No mark is published here."
      />
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
    </div>
  );
}
