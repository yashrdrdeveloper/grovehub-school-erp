import { readFile, readdir } from 'node:fs/promises';
import path from 'node:path';

/**
 * CSS Modules are only valid when they use published semantic design tokens.
 * Keeps theme and responsive extraction from introducing undeclared CSS vars.
 */
const root = process.cwd();
const styles = path.join(root, 'packages/styles');
const modules = path.join(root, 'packages');
const declared = new Set();
const violations = [];

const visit = async (directory, predicate, callback) => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const filename = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      await visit(filename, predicate, callback);
      continue;
    }

    if (predicate(filename)) {
      await callback(filename);
    }
  }
};

await visit(styles, (filename) => filename.endsWith('.css'), async (filename) => {
  const css = await readFile(filename, 'utf8');
  for (const match of css.matchAll(/(--gh-[\w-]+)\s*:/g)) {
    declared.add(match[1]);
  }
});

let checked = 0;
await visit(modules, (filename) => filename.endsWith('.module.css'), async (filename) => {
  checked += 1;
  const css = await readFile(filename, 'utf8');

  for (const match of css.matchAll(/var\((--gh-[\w-]+)/g)) {
    if (!declared.has(match[1])) {
      violations.push(`${path.relative(root, filename)}: undefined design token ${match[1]}`);
    }
  }

  if (/:global\(\s*\.\w+\s*\)\s*\{/.test(css)) {
    violations.push(`${path.relative(root, filename)}: unscoped global classes are not a CSS Module migration`);
  }
});

if (violations.length > 0) {
  violations.forEach((violation) => process.stderr.write(violation + '\n'));
  process.exitCode = 1;
} else {
  process.stdout.write(`CSS Module token ownership verified across ${checked} scoped stylesheets.\n`);
}
