-- Read-only evidence for the single real-model worker smoke submission.
SELECT json_build_object(
  'submissionId',s.id,'lab',l.title,'status',j.status,
  'judgeVerdict',r.verdict,
  'checks',(SELECT json_agg(json_build_object(
      'constraint',c.description,'verdict',v.verdict,
      'source',v.source,'evidence',v.evidence))
    FROM constraint_verdicts v JOIN lab_constraints c ON c.id=v.constraint_id
    WHERE v.submission_id=s.id),
  'calls',(SELECT json_agg(json_build_object(
      'callId',m.id,'chargedMicroUsd',m.charged_micro_usd,
      'status',m.status,'log',m.call_log))
    FROM model_calls m WHERE m.submission_id=s.id)
)
FROM submissions s JOIN labs l ON l.id=s.lab_id
JOIN jobs j ON j.submission_id=s.id
JOIN judge_results r ON r.submission_id=s.id
WHERE s.id='e7ad52ed-686a-4d50-abc7-b426fb8d1931';
