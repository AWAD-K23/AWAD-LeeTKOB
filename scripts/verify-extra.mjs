// Compile and execute our five cases without host mounts or credentials.
import { execFileSync } from 'node:child_process';

const archive = execFileSync('tar', ['-cf', '-', '-C', 'eval/extra', '.']);
const command = `set -eu
tar -xf - -C /work
count=0
for source in /work/solutions/*/*.cpp; do
  problem=$(basename "$(dirname "$source")")
  g++ -std=c++17 -O2 "$source" -o /work/program
  for input in /work/problems/"$problem"/tests/*.in; do
    timeout 3 /work/program < "$input" > /work/actual
    awk '{for(i=1;i<=NF;i++) print $i}' /work/actual > /work/actual.tokens
    expected="$(dirname "$input")/$(basename "$input" .in).out"
    awk '{for(i=1;i<=NF;i++) print $i}' "$expected" > /work/expected.tokens
    cmp /work/actual.tokens /work/expected.tokens
  done
  count=$((count+1))
done
test "$count" -ge 5
echo "Verified $count Accepted extra solutions with g++ -std=c++17 -O2"
`;

const output = execFileSync('docker', [
  'run', '--rm', '--interactive', '--network', 'none',
  '--memory', '512m', '--memory-swap', '512m', '--pids-limit', '64',
  '--cpus', '2', '--read-only', '--cap-drop', 'ALL',
  '--security-opt', 'no-new-privileges', '--user', '65534:65534',
  '--tmpfs', '/work:rw,exec,nosuid,size=64m,uid=65534,gid=65534,mode=700',
  '--tmpfs', '/tmp:rw,noexec,nosuid,size=32m,mode=1777',
  '--entrypoint', 'sh', 'leetkob-checker:pa3', '-c', command,
], { input: archive, timeout: 120_000, encoding: 'utf8', maxBuffer: 1024 * 1024 });
process.stdout.write(output);
