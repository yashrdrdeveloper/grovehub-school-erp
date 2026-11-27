import { readFile, readdir } from 'node:fs/promises';
import path from 'node:path';

import { packageSourceRoots, resolvePackageSource } from './lib/source-layout.mjs';

/**
 * CSS Modules are only valid when they use published semantic design tokens.
 * Keeps theme and responsive extraction from introducing undeclared CSS vars.
 */
const root = process.cwd();
const styleOwners = [resolvePackageSource('styles'), path.join(root, 'app/theme')];
const declared = new Set();
const violations = [];

const visit = async (directory, predicate, callback) => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const filename = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      if (!['node_modules', 'dist', '.generated'].includes(entry.name)) {
        await visit(filename, predicate, callback);
      }
      continue;
    }

    if (predicate(filename)) {
      await callback(filename);
    }
  }
};

for (const styles of styleOwners)
  await visit(
    styles,
    (filename) => filename.endsWith('.css'),
    async (filename) => {
      const css = await readFile(filename, 'utf8');
      for (const match of css.matchAll(/(--gh-[\w-]+)\s*:/g)) {
        declared.add(match[1]);
      }
    },
  );

const checkedCounter = { current: 0 };
for (const modules of packageSourceRoots)
  await visit(
    modules,
    (filename) => filename.endsWith('.module.css'),
    async (filename) => {
      checkedCounter.current += 1;
      const css = await readFile(filename, 'utf8');

      for (const match of css.matchAll(/var\((--gh-[\w-]+)/g)) {
        if (!declared.has(match[1])) {
          violations.push(`${path.relative(root, filename)}: undefined design token ${match[1]}`);
        }
      }

      if (/:global\(\s*\.\w+\s*\)\s*\{/.test(css)) {
        violations.push(
          `${path.relative(root, filename)}: unscoped global classes are not a CSS Module migration`,
        );
      }

      if (/#[0-9a-f]{3,8}\b|\brgba?\(|\bhsla?\(/i.test(css)) {
        violations.push(
          `${path.relative(root, filename)}: use semantic theme tokens instead of color literals`,
        );
      }

      if (/:root\s*\{/.test(css)) {
        violations.push(
          `${path.relative(root, filename)}: root design tokens belong to packages/styles`,
        );
      }
    },
  );

if (violations.length > 0) {
  violations.forEach((violation) => process.stderr.write(violation + '\n'));
  process.exitCode = 1;
} else {
  process.stdout.write(
    `CSS Module token ownership verified across ${checkedCounter.current} scoped stylesheets.\n`,
  );
}
