import { access, readFile, readdir } from 'node:fs/promises';
import path from 'node:path';

import { validateRendererSource } from './lib/renderer-policy.mjs';

/** Renderer functions are pure JSX selectors with typed contracts and owner-local tests. */
const root = path.join(process.cwd(), 'packages');
const errors = [];
let checked = 0;
const exists = async (filename) => {
  try {
    await access(filename);
    return true;
  } catch {
    return false;
  }
};
const walk = async (folder) => {
  for (const entry of await readdir(folder, { withFileTypes: true })) {
    const filename = path.join(folder, entry.name);
    if (entry.isDirectory()) {
      if (entry.name !== 'node_modules' && entry.name !== '.generated') await walk(filename);
      continue;
    }
    if (!filename.endsWith('.tsx')) continue;
    const relative = path.relative(root, filename).replaceAll(path.sep, '/');
    if (!relative.includes('/renderer/') || relative.includes('/__tests__/')) continue;
    checked++;
    const basename = path.basename(filename, '.tsx');
    if (!/^[a-z][a-zA-Z0-9]*\.[a-z][a-zA-Z0-9]*$/.test(basename)) {
      errors.push(relative + ': use renderer/<owner>.<useCase>.tsx naming');
    }
    const source = await readFile(filename, 'utf8');
    for (const problem of validateRendererSource(filename, source)) {
      errors.push(relative + ': ' + problem);
    }
    if (!(await exists(path.join(path.dirname(filename), '__tests__', basename + '.test.tsx')))) {
      errors.push(relative + ': missing owner-local renderer test');
    }
  }
};
await walk(root);
for (const issue of errors) process.stderr.write(issue + '\n');
if (errors.length) process.exitCode = 1;
else
  process.stdout.write(
    'Renderer purity, naming and tests verified across ' + checked + ' JSX renderer files.\n',
  );
