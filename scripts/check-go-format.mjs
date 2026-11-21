// Fails if gofmt finds Go files that need formatting.
import { execFileSync } from 'node:child_process';

// Include tracked and untracked Go files, but skip ignored files.
const files = execFileSync('git', ['ls-files', '--cached', '--others', '--exclude-standard', '*.go'], {
  encoding: 'utf8',
}).trim();

if (!files) process.exit(0);

// gofmt -l prints files whose formatting differs from the standard.
const result = execFileSync('gofmt', ['-l', ...files.split(/\r?\n/)], {
  encoding: 'utf8',
});

if (result.trim()) {
  console.error(`Go files need formatting:\n${result}`);
  process.exit(1);
}
