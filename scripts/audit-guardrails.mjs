// Mutation checks run in a disposable Git export, never the working checkout.
import { execFileSync, spawnSync } from 'node:child_process';
import { mkdtempSync, readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve, dirname } from 'node:path';

if (!process.env.TEST_DATABASE_URL) throw new Error('TEST_DATABASE_URL is required for G4');
const output = resolve(process.argv[2] ?? 'docs/PA4/evidence/guardrail-mutations.json');
const root = process.cwd();
const temp = mkdtempSync(join(tmpdir(), 'leetkob-guardrails-'));
const git = args => execFileSync('git', ['-c', `safe.directory=${root.replaceAll('\\', '/')}`, ...args], { encoding: 'utf8' }).trim();
const commit = git(['rev-parse', 'HEAD']);
git(['-c', 'core.autocrlf=false', 'archive', '--format=tar', `--output=${join(temp, 'source.tar')}`, 'HEAD']);
execFileSync('tar', ['-xf', join(temp, 'source.tar'), '-C', temp]);
const semantic = 'backend/internal/semantic/check.go';
const worker = 'backend/internal/worker/model.go';
const replace = (source, from, to) => {
  if (!source.includes(from)) throw new Error('Mutation anchor missing');
  return source.replace(from, to);
};
const mutations = [
  {
    guard: 'G1 schema', file: semantic, pkg: './internal/semantic', test: '^TestG1',
    mutate: source => replace(source, 'invalid := ErrInvalidModelResponse', `
    // Mutation: accept every response as a decisive pass, skipping schema validation.
    if true {
        out := []hardrules.Verdict{}
        for _, constraint := range constraints {
            out = append(out, hardrules.Verdict{ConstraintID: constraint.ID, Verdict: "pass", Source: "model"})
        }
        return out, nil
    }
    invalid := ErrInvalidModelResponse`),
  },
  {
    guard: 'G2 exact numbered evidence', file: semantic, pkg: './internal/semantic', test: '^TestG2',
    mutate: source => replace(source, 'point.Text != cleanLines[point.Line-1]', 'false'),
  },
  {
    guard: 'G3 discard late success', file: semantic, pkg: './internal/semantic', test: '^TestG3',
    mutate: source => replace(source, 'err != nil || callCtx.Err() != nil', 'err != nil'),
  },
  {
    guard: 'G4 persistent reservation', file: worker, pkg: './internal/worker', test: 'TestDurableWorker/G4',
    mutate: source => replace(source,
      'if err := reserveCall(ctx, m.db, callID, m.lab, m.submission, m.config.BudgetMicroUSD); err != nil {\n\t\treturn nil, err\n\t}', ''),
  },
  {
    guard: 'G5 unknown mark fields', file: semantic, pkg: './internal/semantic', test: '^TestG5',
    mutate: source => replace(replace(source,
      'if err := strictKeys(json.NewDecoder(bytes.NewReader(raw)), 0); err != nil {',
      'if false {'), 'decoder.DisallowUnknownFields()', '// Mutation: allow unknown fields'),
  },
  {
    guard: 'Semantic violation confirmation', file: semantic, pkg: './internal/semantic', test: '^TestViolationConfirmation',
    mutate: source => replace(source,
      'validated = c.confirmViolations(ctx, timeout, Request{\n\t\tStatement: statement, Constraints: pending, Source: cleaned,\n\t}, string(source), validated)',
      '// Mutation: skip confirmation of semantic violations'),
  },
];
const env = { ...process.env, MODEL_OFF: '1', MODEL_API_KEY: '' };
const results = [];
for (const item of mutations) {
  const file = join(temp, item.file);
  const original = readFileSync(file, 'utf8');
  const run = () => spawnSync('go', ['-C', join(temp, 'backend'), 'test', item.pkg, '-run', item.test, '-count=1', '-timeout=15s'], {
    env, encoding: 'utf8', timeout: 120_000, maxBuffer: 1024 * 1024,
  });
  const baseline = run();
  if (baseline.error || baseline.status !== 0) throw new Error(`Baseline failed for ${item.guard}: ${baseline.stdout}${baseline.stderr}`);
  writeFileSync(file, item.mutate(original));
  const mutated = run();
  writeFileSync(file, original);
  if (mutated.error || mutated.status !== 1 || !mutated.stdout.includes('--- FAIL:')) {
    throw new Error(`Mutation was not caught by a failing test for ${item.guard}: ${mutated.stdout}${mutated.stderr}`);
  }
  results.push({ guard: item.guard, test: item.test, baselineExit: baseline.status, mutatedExit: mutated.status,
    failure: mutated.stdout.trim(), stderr: mutated.stderr.trim() });
  console.log(`${item.guard}: baseline passed; removed guard failed its test`);
}
mkdirSync(dirname(output), { recursive: true });
writeFileSync(output, JSON.stringify({ commit, generatedAt: new Date().toISOString(), results }, null, 2) + '\n', { flag: 'wx' });
