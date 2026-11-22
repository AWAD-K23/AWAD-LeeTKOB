// Run the real image with networking disabled and save the unchanged scorer's
// report. Exit 1 from scoring is expected for pending semantic work in this mode.
import { execFileSync, spawnSync } from 'node:child_process';
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { join, resolve } from 'node:path';
import { randomUUID } from 'node:crypto';

const dataset = resolve(process.argv[2] ?? 'eval/official');
const output = resolve(process.argv[3] ?? 'eval/runs/local-offline');
if (existsSync(join(output, 'predictions.jsonl'))) {
  throw new Error('Output already contains predictions; choose a new run folder');
}
const archive = execFileSync('tar', ['-cf', '-', '-C', dataset, '.'], {
  maxBuffer: 64 * 1024 * 1024,
});
const name = 'leetkob-offline-' + randomUUID();
let predictions;

try {
  predictions = execFileSync('docker', [
    'run', '--name', name, '--interactive', '--network', 'none',
    '--memory', '512m', '--memory-swap', '512m', '--pids-limit', '64',
    '--cpus', '2', '--read-only', '--cap-drop', 'ALL',
    '--security-opt', 'no-new-privileges',
    '--tmpfs', '/tmp:rw,nosuid,size=128m,mode=1777', '-e', 'MODEL_OFF=1',
    '--entrypoint', 'sh', 'leetkob-checker:pa3', '-c',
    'mkdir /tmp/data && tar -xf - -C /tmp/data && exec /usr/local/bin/checker --scanner /usr/local/bin/scan --data /tmp/data --predictions -',
  ], { input: archive, timeout: 600_000, maxBuffer: 16 * 1024 * 1024 });
} finally {
  // Removing this task's named container also stops it after a CLI timeout.
  execFileSync('docker', ['rm', '-f', name], { timeout: 10_000, stdio: 'ignore' });
}

mkdirSync(output, { recursive: true });
const predictionPath = join(output, 'predictions.jsonl');
writeFileSync(predictionPath, predictions, { flag: 'wx' });
const scored = spawnSync(process.execPath, [
  'eval/official/eval.mjs', '--data', dataset, '--predictions', predictionPath,
  '--json', join(output, 'summary.json'),
], { encoding: 'utf8' });
writeFileSync(join(output, 'report.txt'), scored.stdout + scored.stderr);
process.stdout.write(scored.stdout);
process.stderr.write(scored.stderr);
if (scored.error || ![0, 1].includes(scored.status)) {
  throw new Error('Official evaluator rejected predictions or failed to run');
}

// Offline checks must decide every supported ban and leave all semantic work
// unsure. The scorer's threshold failure is not hidden or relabelled as PASS.
const rows = readFileSync(predictionPath, 'utf8').trim().split('\n').map(JSON.parse);
const summary = JSON.parse(readFileSync(join(output, 'summary.json')));
if (summary.wrongful !== 0 || summary.missing !== 0) {
  throw new Error('Offline run has wrongful deductions or missing predictions');
}

function isBan(text) {
  return /^(Do not use|Do not call)/.test(text) && /`[^`]+`/.test(text);
}

for (const row of rows) {
  const problem = row.solution.split('/')[0];
  const constraints = JSON.parse(readFileSync(join(dataset, 'problems', problem, 'constraints.json')));
  const text = constraints.find(c => c.id === row.constraint)?.text ?? '';
  const ban = isBan(text);
  if (ban ? !['pass', 'violation'].includes(row.verdict) : row.verdict !== 'unsure') {
    throw new Error('Unexpected offline decision: ' + row.solution + '/' + row.constraint);
  }
}
for (const missed of summary.missedPairs) {
  const [solution, constraint] = missed.split(' ');
  const problem = solution.split('/')[0];
  const constraints = JSON.parse(readFileSync(join(dataset, 'problems', problem, 'constraints.json')));
  if (isBan(constraints.find(c => c.id === constraint)?.text ?? '')) {
    throw new Error('Ban violation was missed or has invalid evidence: ' + missed);
  }
}
console.log(`Saved ${rows.length} offline predictions; evaluator exit ${scored.status}`);
