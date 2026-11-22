import { readdir, access } from 'node:fs/promises';
import path from 'node:path';

const root = path.join(process.cwd(), 'packages');
const failures = [];

const exists = async (file) => {
  try { await access(file); return true; } catch { return false; }
};
const sourceExists = async (folder, name) =>
  await exists(path.join(folder, `${name}.ts`))
  || await exists(path.join(folder, `${name}.tsx`));

const units = [
  ['base/utilities', true],
  ['base/readers', true],
  ['auth/pages/login/helpers', true],
  ['auth/pages/forgotPassword/helpers', true],
  ['auth/pages/recovery/helpers', true],
  ['base/events', false],
  ['auth/pages/login/hooks', false],
  ['auth/pages/forgotPassword/hooks', false],
  ['auth/pages/recovery/hooks', false],
  ['shell/hooks', false],
  ['business/auth/services', false],
];

let checked = 0;
for (const [relative, requirePureTests] of units) {
  const parent = path.join(root, relative);
  if (!await exists(parent)) continue;
  for (const item of await readdir(parent, { withFileTypes: true })) {
    if (!item.isDirectory()) continue;
    const folder = path.join(parent, item.name);
    const name = item.name;
    // Enforce explicit callable contracts for named implementation units.
    if (!await sourceExists(folder, name)) continue;
    checked++;
    if (!await exists(path.join(folder, 'index.ts'))) {
      failures.push(`${relative}/${name}: missing index.ts public API`);
    }
    if (!await exists(path.join(folder, 'types', `${name}.ts`))) {
      failures.push(`${relative}/${name}: missing types/${name}.ts callable contract`);
    }
    if (requirePureTests && !await exists(path.join(folder, '__tests__', `${name}.test.ts`))) {
      failures.push(`${relative}/${name}: missing __tests__/${name}.test.ts`);
    }
  }
}

const guardedPages = {
  'auth/pages/login': 'Login.tsx',
  'auth/pages/forgotPassword': 'ForgotPassword.tsx',
  'auth/pages/recovery': 'Recovery.tsx',
  'shell': 'Shell.tsx',
};
for (const [relative, entryName] of Object.entries(guardedPages)) {
  const directory = path.join(root, relative);
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    if (entry.isDirectory() || entry.name === 'index.ts' || entry.name === entryName || entry.name.endsWith('.md')) {
      continue;
    }
    if (/\.[jt]sx?$/.test(entry.name)) {
      failures.push(`${relative}: root code ${entry.name} must live in a named subunit`);
    }
  }
}

if (failures.length > 0) {
  failures.forEach(failure => process.stderr.write(`${failure}\n`));
  process.exitCode = 1;
} else {
  process.stdout.write(`Self-contained unit boundaries verified across ${checked} operation units.\n`);
}
