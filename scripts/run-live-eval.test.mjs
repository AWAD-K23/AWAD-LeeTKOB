import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, readdirSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { spawnSync, execFileSync } from 'node:child_process';

const script = resolve('scripts/run-live-eval.mjs');
const environment = { ...process.env, MODEL_OFF: '0', MODEL_API_KEY: 'fake-test-key', MODEL_NAME: 'fake-model' };

test('live runner refuses offline mode before any checker/provider work', () => {
  const temp = mkdtempSync(join(tmpdir(), 'leetkob-live-mode-'));
  const result = spawnSync(process.execPath, [script, 'unused'], {
    cwd: temp, env: { ...environment, MODEL_OFF: '1' }, encoding: 'utf8',
  });
  assert.notEqual(result.status, 0);
  assert.match(result.stderr, /Explicit MODEL_OFF=0/);
  assert.deepEqual(readdirSync(temp), []);
});

test('live runner refuses an existing output folder instead of overwriting evidence', () => {
  const temp = mkdtempSync(join(tmpdir(), 'leetkob-live-output-'));
  const git = args => execFileSync('git', ['-C', temp, ...args], { stdio: 'ignore' });
  git(['init']);
  git(['-c','user.name=Fixture','-c','user.email=fixture@example.test','commit','--allow-empty','-m','Synthetic fixture']);
  mkdirSync(join(temp, 'existing'));
  const result = spawnSync(process.execPath, [script, 'existing'], {
    cwd: temp, env: environment, encoding: 'utf8',
  });
  assert.notEqual(result.status, 0);
  assert.match(result.stderr, /EEXIST/);
  assert.deepEqual(readdirSync(join(temp, 'existing')), []);
});
