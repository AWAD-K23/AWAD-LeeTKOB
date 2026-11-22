// Verify the pinned official snapshot without rewriting any published file.
import { readFileSync, readdirSync } from 'node:fs';
import { createHash } from 'node:crypto';

const root = new URL('../eval/official/', import.meta.url);
const manifest = JSON.parse(readFileSync(new URL('../eval/official-sha256.json', import.meta.url)));
const found = [];

function visit(relative = '') {
  for (const entry of readdirSync(new URL(relative, root), { withFileTypes: true })) {
    const name = relative + entry.name;
    if (entry.isDirectory()) visit(name + '/');
    else {
      const hash = createHash('sha256').update(readFileSync(new URL(name, root))).digest('hex');
      if (manifest[name] !== hash) throw new Error(`Official file changed: ${name}`);
      found.push(name);
    }
  }
}

visit();
if (found.length !== Object.keys(manifest).length) throw new Error('Official files are missing');
console.log(`Verified ${found.length} unchanged official files`);
