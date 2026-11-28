// Compare exactly three frozen-commit official dev runs without provider calls.
import { readFileSync } from 'node:fs';
import { join } from 'node:path';

const folders = process.argv.slice(2);
if (folders.length !== 3) throw new Error('Usage: node scripts/compare-model-runs.mjs <run1> <run2> <run3>');
const read = file => JSON.parse(readFileSync(file, 'utf8'));
const rows = file => readFileSync(file, 'utf8').trim().split(/\r?\n/).filter(Boolean).map(JSON.parse);
const metadata = folders.map(folder => read(join(folder, 'metadata.json')));
const signature = item => JSON.stringify({ commit: item.commit, model: item.model, host: item.host,
  budget: item.budgetMicroUsdPerProblem, image: item.sandboxImage,
  data: item.datasets?.find(dataset => dataset.name === 'official')?.sha256 });
if (metadata.some(item => signature(item) !== signature(metadata[0]))) throw new Error('Runs differ in commit/model/configuration/dataset');
const predictions = folders.map(folder => {
  const map = new Map();
  for (const row of rows(join(folder, 'official', 'predictions.jsonl'))) {
    const key = row.solution + ' ' + row.constraint;
    if (map.has(key)) throw new Error('Duplicate prediction pair');
    map.set(key, row.verdict);
  }
  return map;
});
if (predictions.some(map => map.size !== predictions[0].size || [...map.keys()].some(key => !predictions[0].has(key)))) throw new Error('Pair populations differ');
const changed = [...predictions[0].keys()].filter(key => new Set(predictions.map(map => map.get(key))).size > 1).sort();
const summaries = folders.map(folder => read(join(folder, 'official', 'summary.json')));
const versions = folders.map(folder => [...new Set(rows(join(folder, 'official', 'calls.jsonl')).filter(e => e.event === 'recorded' && e.call.outcome === 'response_received').map(e => e.call.modelVersion))].sort());
if (versions.some(items => JSON.stringify(items) !== JSON.stringify(versions[0]))) throw new Error('Provider model versions differ between runs');
const report = {
  commit: metadata[0].commit, model: metadata[0].model, modelVersions: versions[0],
  definition: 'Unique (solution, constraint) pairs with more than one verdict across all three runs',
  changedPairCount: changed.length, changedPairs: changed,
  recallRange: [Math.min(...summaries.map(s => s.recall)), Math.max(...summaries.map(s => s.recall))],
  unsureRateRange: [Math.min(...summaries.map(s => s.unsureRate)), Math.max(...summaries.map(s => s.unsureRate))],
  runs: folders.map((folder,i) => ({ folder, wrongful: summaries[i].wrongful, recall: summaries[i].recall,
    unsureRate: summaries[i].unsureRate, ...read(join(folder, 'official', 'usage-summary.json')) })),
};
console.log(JSON.stringify(report, null, 2));
