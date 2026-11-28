// Explicit live evaluation. The key stays in the trusted checker process;
// its Docker scanner receives no host environment, mounts, or network.
import { spawn, execFileSync } from 'node:child_process';
import { existsSync, mkdirSync, readFileSync, writeFileSync, readdirSync, mkdtempSync } from 'node:fs';
import { resolve, join, basename } from 'node:path';
import { tmpdir } from 'node:os';
import { createHash } from 'node:crypto';

const output = process.argv[2];
if (!output) throw new Error('Usage: node scripts/run-live-eval.mjs <new-run-folder> [dataset ...]');
if (existsSync('.env')) process.loadEnvFile('.env');
if (process.env.MODEL_OFF !== '0' || !process.env.MODEL_API_KEY || !process.env.MODEL_NAME) {
  throw new Error('Explicit MODEL_OFF=0, MODEL_NAME and MODEL_API_KEY are required');
}
const root = process.cwd();
const git = args => execFileSync('git', ['-c', `safe.directory=${root.replaceAll('\\', '/')}`, ...args], { encoding: 'utf8' }).trim();
const dirty = git(['status', '--porcelain', '--', 'backend', 'sandbox', 'scripts', '.github', 'eval/extra', '.gitattributes']);
if (dirty) throw new Error('Commit checker, scripts, cases and CI changes before recording a live run');
const datasets = (process.argv.length > 3 ? process.argv.slice(3) : ['eval/official', 'eval/extra']).map(folder => resolve(folder));
const names = datasets.map(folder => basename(folder));
if (new Set(names).size !== names.length) throw new Error('Dataset folder names must be unique');
const out = resolve(output);
mkdirSync(out); // Refuse any existing run folder, including partially failed runs.
const write = (path, value) => writeFileSync(path, JSON.stringify(value, null, 2) + '\n', { flag: 'wx' });
function hashDataset(folder) {
  const hash = createHash('sha256');
  function visit(relative = '') {
    for (const entry of readdirSync(join(folder, relative), { withFileTypes: true }).sort((a,b) => a.name.localeCompare(b.name))) {
      const name = relative + entry.name;
      if (entry.isDirectory()) visit(name + '/');
      else { hash.update(name + '\0'); hash.update(readFileSync(join(folder, name))); }
    }
  }
  visit();
  return hash.digest('hex');
}
const run = (command, args, timeout = 1_800_000) => new Promise((accept, reject) => {
  const child = spawn(command, args, { env: process.env, stdio: ['ignore', 'pipe', 'pipe'], windowsHide: true });
  let stdout = '', stderr = '';
  child.stdout.on('data', data => { stdout += data; });
  child.stderr.on('data', data => { stderr += data; });
  const timer = setTimeout(() => child.kill(), timeout);
  child.on('error', error => { clearTimeout(timer); reject(error); });
  child.on('close', status => { clearTimeout(timer); accept({ status, stdout, stderr }); });
});
const metadata = {
  commit: git(['rev-parse', 'HEAD']), model: process.env.MODEL_NAME,
  host: 'generativelanguage.googleapis.com', budgetMicroUsdPerProblem: process.env.MODEL_LAB_BUDGET_MICRO_USD ?? '1000000',
  startedAt: new Date().toISOString(), scanner: 'isolated Docker; no model credentials',
  sandboxImage: execFileSync('docker', ['image', 'inspect', 'leetkob-sandbox:pa3', '--format', '{{.Id}}'], { encoding: 'utf8' }).trim(),
  datasets: datasets.map((folder, i) => ({ name: names[i], sha256: hashDataset(folder) })),
};
write(join(out, 'metadata.json'), metadata);
const temp = mkdtempSync(join(tmpdir(), 'leetkob-live-'));
const binary = join(temp, process.platform === 'win32' ? 'checker.exe' : 'checker');
const built = await run('go', ['-C', 'backend', 'build', '-trimpath', '-o', binary, './cmd/checker']);
if (built.status !== 0) throw new Error(`Checker build failed: ${built.stderr}`);
let exit = 0;
const results = [];
for (let i = 0; i < datasets.length; i++) {
  const folder = join(out, names[i]);
  mkdirSync(folder);
  console.log(`Generating fresh real-model predictions for ${names[i]}`);
  const checked = await run(binary, ['--data', datasets[i], '--scanner', 'docker', '--predictions', join(folder, 'predictions.jsonl'), '--call-log', join(folder, 'calls.jsonl')]);
  writeFileSync(join(folder, 'checker.txt'), checked.stdout + checked.stderr, { flag: 'wx' });
  if (checked.status !== 0) { exit = checked.status || 1; results.push({ dataset: names[i], checkerExit: checked.status }); continue; }
  const scored = await run(process.execPath, ['eval/official/eval.mjs', '--data', datasets[i], '--predictions', join(folder, 'predictions.jsonl'), '--json', join(folder, 'summary.json')]);
  writeFileSync(join(folder, 'report.txt'), scored.stdout + scored.stderr, { flag: 'wx' });
  process.stdout.write(scored.stdout);
  process.stderr.write(scored.stderr);
  const usage = await run(process.execPath, ['scripts/summarize-model-run.mjs', folder]);
  if (usage.status !== 0) { exit = usage.status || 1; throw new Error('Call audit could not be summarized'); }
  const summary = JSON.parse(readFileSync(join(folder, 'usage-summary.json')));
  if (summary.successfulResponses === 0 || summary.providerRequestIds < summary.successfulResponses) {
    exit = 1;
    console.error('No audited real responses, or successful responses lack provider IDs');
  }
  if (scored.status !== 0) exit = scored.status || 1;
  results.push({ dataset: names[i], evaluatorExit: scored.status, ...summary });
}
write(join(out, 'run-summary.json'), {
  commit: metadata.commit, finishedAt: new Date().toISOString(), results,
  knownEstimatedCostMicroUsd: results.reduce((n,r) => n + (r.knownEstimatedCostMicroUsd ?? 0), 0),
  unknownUsageReservedMicroUsd: results.reduce((n,r) => n + (r.unknownUsageReservedMicroUsd ?? 0), 0),
  exit,
});
process.exitCode = exit; // Every nonzero evaluator exit fails this live gate.
