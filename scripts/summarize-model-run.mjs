// Summarize recorded provider usage without changing predictions or labels.
import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';

const folder = process.argv[2];
if (!folder) throw new Error('Usage: node scripts/summarize-model-run.mjs <run-folder>');
const readRows = name => readFileSync(join(folder, name), 'utf8').trim()
  .split('\n').filter(Boolean).map(JSON.parse);
const events = readRows('calls.jsonl');
const reserved = events.filter(e => e.event === 'reserved');
const recorded = events.filter(e => e.event === 'recorded');
const known = recorded.filter(e => e.call.usageKnown);
const total = field => known.reduce((sum, e) => sum + e.call[field], 0);
const predictions = readRows('predictions.jsonl');
const starts = recorded.map(e => Date.parse(e.call.startedAt)).sort((a,b) => a-b);
let peak = 0;
for (let i = 0; i < starts.length; i++) {
  peak = Math.max(peak, starts.filter(t => t >= starts[i] && t < starts[i]+60_000).length);
}

const knownCost = total('estimatedCostMicroUsd');
const summary = {
  model: recorded.find(e => e.call.model)?.call.model,
  reservedCalls: reserved.length,
  recordedCalls: recorded.length,
  knownUsageCalls: known.length,
  unknownUsageCalls: recorded.length-known.length,
  unfinishedReservations: reserved.length-recorded.length,
  successfulResponses: recorded.filter(e => e.call.outcome === 'response_received').length,
  errors: recorded.filter(e => e.call.outcome !== 'response_received')
    .map(e => ({ solution:e.solution, outcome:e.call.outcome })),
  providerRequestIds: recorded.filter(e => e.call.providerRequestId).length,
  inputTokens: total('inputTokens'),
  outputTokens: total('outputTokens'),
  thinkingTokens: total('thinkingTokens'),
  knownEstimatedCostMicroUsd: knownCost,
  unknownUsageReservedMicroUsd: (reserved.length-known.length)*478413,
  estimatedUsdFor80SimilarCalls: known.length ? knownCost/known.length*80/1_000_000 : null,
  minimumStartGapMs: starts.length>1 ? Math.min(...starts.slice(1).map((t,i)=>t-starts[i])) : null,
  peakRequestsIn60Seconds: peak,
  predictionPairs: predictions.length,
  verdictCounts: Object.fromEntries(['pass','violation','unsure']
    .map(v=>[v,predictions.filter(p=>p.verdict===v).length])),
};

if (peak > 15) throw new Error('Recorded request rate exceeded the user limit');
writeFileSync(join(folder, 'usage-summary.json'), JSON.stringify(summary,null,2)+'\n');
console.log(JSON.stringify(summary,null,2));
