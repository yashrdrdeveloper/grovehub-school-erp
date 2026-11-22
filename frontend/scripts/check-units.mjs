import { readdir, access, readFile } from 'node:fs/promises';
import path from 'node:path';

const root = path.join(process.cwd(), 'packages');
const failures = [];

const exists = async (file) => {
  try {
    await access(file);
    return true;
  } catch {
    return false;
  }
};
const sourceExists = async (folder, name) =>
  (await exists(path.join(folder, `${name}.ts`))) ||
  (await exists(path.join(folder, `${name}.tsx`)));

const units = [
  ['base/utilities', true],
  ['base/readers', true],
  ['base/events', false],
  ['domain/auth/pages/login/hooks', false],
  ['domain/auth/pages/forgotPassword/hooks', false],
  ['domain/auth/pages/recovery/hooks', false],
  ['domain/dashboard/pages/dashboard/hooks', false],
  ['shell/hooks', false],
  ['business/auth/services', false],
];

let checked = 0;
for (const [relative, requirePureTests] of units) {
  const parent = path.join(root, relative);
  if (!(await exists(parent))) continue;
  for (const item of await readdir(parent, { withFileTypes: true })) {
    if (!item.isDirectory()) continue;
    const folder = path.join(parent, item.name);
    const name = item.name;
    // Enforce explicit callable contracts for named implementation units.
    if (!(await sourceExists(folder, name))) continue;
    checked++;
    if (!(await exists(path.join(folder, 'index.ts')))) {
      failures.push(`${relative}/${name}: missing index.ts public API`);
    }
    if (!(await exists(path.join(folder, 'types', `${name}.ts`)))) {
      failures.push(`${relative}/${name}: missing types/${name}.ts callable contract`);
    }
    if (requirePureTests && !(await exists(path.join(folder, '__tests__', `${name}.test.ts`)))) {
      failures.push(`${relative}/${name}: missing __tests__/${name}.test.ts`);
    }
  }
}

/** All helper, constant and JSX renderer folders follow flat source ownership. */
const flatUnits = [];
const collectFlatUnits = async (directory, relative = '') => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    if (!entry.isDirectory()) continue;
    const child = path.join(directory, entry.name);
    const childRelative = path.join(relative, entry.name);
    if (['helpers', 'constants', 'renderer'].includes(entry.name)) {
      flatUnits.push(childRelative);
    }
    await collectFlatUnits(child, childRelative);
  }
};
await collectFlatUnits(root);

for (const relative of flatUnits) {
  const directory = path.join(root, relative);
  const entries = await readdir(directory, { withFileTypes: true });
  for (const item of entries) {
    if (item.name === 'index.ts') {
      failures.push(relative + ': flat directory must not have an index.ts barrel');
    }
    if (item.isDirectory() && item.name !== 'types' && item.name !== '__tests__') {
      failures.push(
        relative + '/' + item.name + ': helpers/constants/renderer nesting is forbidden',
      );
    }
  }
  const tests = path.join(directory, '__tests__');
  if (!(await exists(tests))) continue;
  for (const item of await readdir(tests, { withFileTypes: true })) {
    if (!item.isFile() || !/\.test\.tsx?$/.test(item.name)) continue;
    const file = await readFile(path.join(tests, item.name), 'utf8');
    const descriptions = [...file.matchAll(/\bdescribe\(\s*['"]([^'"]+)['"]/g)];
    if (descriptions.length !== 1) {
      failures.push(
        relative + '/__tests__/' + item.name + ': exactly one top-level describe is required',
      );
      continue;
    }
    const exportedName = descriptions[0][1];
    const matchingImport = [...file.matchAll(/import\s*\{([^}]+)\}\s*from/g)].some((match) =>
      match[1].split(',').some((name) => name.trim() === exportedName),
    );
    if (!matchingImport) {
      failures.push(
        relative + '/__tests__/' + item.name + ': describe must match a named imported export',
      );
    }
    if (/\bjest\.mock\(/.test(file)) {
      failures.push(
        relative + '/__tests__/' + item.name + ': pure helper tests must not mock dependencies',
      );
    }
  }
}

const guardedPages = {
  'domain/auth/pages/login': 'Login.tsx',
  'domain/auth/pages/forgotPassword': 'ForgotPassword.tsx',
  'domain/auth/pages/recovery': 'Recovery.tsx',
  'domain/dashboard/pages/dashboard': 'Dashboard.tsx',
  shell: 'Shell.tsx',
};
for (const [relative, entryName] of Object.entries(guardedPages)) {
  const directory = path.join(root, relative);
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    if (
      entry.isDirectory() ||
      entry.name === 'index.ts' ||
      entry.name === entryName ||
      entry.name.endsWith('.md')
    ) {
      continue;
    }
    if (/\.[jt]sx?$/.test(entry.name)) {
      failures.push(`${relative}: root code ${entry.name} must live in a named subunit`);
    }
  }
}

if (failures.length > 0) {
  failures.forEach((failure) => process.stderr.write(`${failure}\n`));
  process.exitCode = 1;
} else {
  process.stdout.write(
    `Self-contained unit boundaries verified across ${checked} operation units.\n`,
  );
}
