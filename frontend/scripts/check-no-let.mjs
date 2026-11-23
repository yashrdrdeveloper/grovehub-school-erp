import { readFile, readdir } from 'node:fs/promises';
import path from 'node:path';

import { findLetDeclarations } from './lib/no-let-policy.mjs';

const root = process.cwd();
const failures = [];
const files = [];
const walk = async (directory) => {
  for (const item of await readdir(directory, { withFileTypes: true })) {
    const absolute = path.join(directory, item.name);
    if (item.isDirectory()) {
      if (!['node_modules', 'dist', '.generated'].includes(item.name)) await walk(absolute);
      continue;
    }
    if (/\.(?:[cm]?[jt]sx?)$/.test(item.name)) files.push(absolute);
  }
};
for (const folder of ['app', 'packages', 'scripts']) await walk(path.join(root, folder));

for (const absolute of files) {
  const relative = path.relative(root, absolute).replaceAll(path.sep, '/');
  const source = await readFile(absolute, 'utf8');
  for (const issue of findLetDeclarations(relative, source)) {
    failures.push(relative + ':' + issue.line + ': let declarations are forbidden (' +
      issue.variables.join(', ') + '); use const with a pure helper or scoped state owner');
  }
}
if (failures.length > 0) {
  for (const issue of failures) process.stderr.write(issue + '\n');
  process.exitCode = 1;
} else {
  process.stdout.write('No let declarations across ' + files.length + ' frontend source and test files.\n');
}
