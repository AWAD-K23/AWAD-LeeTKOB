import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { spawnSync } from 'node:child_process';

function fixtures() {
  const root = mkdtempSync(join(tmpdir(), 'leetkob-compare-test-'));
  const folders = [0,1,2].map(i => join(root, String(i)));
  for (let i = 0; i < folders.length; i++) {
    const folder = folders[i];
    mkdirSync(join(folder, 'official'), { recursive: true });
    const json = (name,data) => writeFileSync(join(folder,name),JSON.stringify(data));
    json('metadata.json', { commit:'synthetic-commit',model:'fake-model',host:'example.test',
      budgetMicroUsdPerProblem:'1',sandboxImage:'fake-image',datasets:[{name:'official',sha256:'fake-hash'}] });
    json('official/summary.json',{recall:i===1?0.9:1,unsureRate:i===1?0.1:0,wrongful:0});
    json('official/usage-summary.json',{knownEstimatedCostMicroUsd:1});
    writeFileSync(join(folder,'official/predictions.jsonl'),[
      {solution:'example/one',constraint:'a',verdict:i===1?'unsure':'violation'},
      {solution:'example/two',constraint:'a',verdict:'pass'},
    ].map(JSON.stringify).join('\n')+'\n');
    writeFileSync(join(folder,'official/calls.jsonl'),JSON.stringify({event:'recorded',call:{outcome:'response_received',providerRequestId:'fake-request',modelVersion:'fake-version'}})+'\n');
  }
  return folders;
}
const run = folders => spawnSync(process.execPath,[resolve('scripts/compare-model-runs.mjs'),...folders],{encoding:'utf8'});

test('counts unique changed pairs across all three runs and reports ranges',()=>{
  const result=run(fixtures());
  assert.equal(result.status,0,result.stderr);
  const report=JSON.parse(result.stdout);
  assert.equal(report.changedPairCount,1);
  assert.deepEqual(report.changedPairs,['example/one a']);
  assert.deepEqual(report.recallRange,[0.9,1]);
  assert.deepEqual(report.unsureRateRange,[0,0.1]);
});
test('refuses evidence from different commits',()=>{
  const folders=fixtures();
  writeFileSync(join(folders[1],'metadata.json'),JSON.stringify({commit:'different'}));
  const result=run(folders);
  assert.notEqual(result.status,0);
  assert.match(result.stderr,/Runs differ/);
});
test('refuses a changed pair population instead of silently counting missing pairs',()=>{
  const folders=fixtures();
  writeFileSync(join(folders[2],'official/predictions.jsonl'),'');
  const result=run(folders);
  assert.notEqual(result.status,0);
  assert.match(result.stderr,/Pair populations differ/);
});
