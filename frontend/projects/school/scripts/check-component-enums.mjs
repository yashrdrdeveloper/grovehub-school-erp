import { readFile, readdir } from 'node:fs/promises';
import path from 'node:path';

import { validateComponentEnums } from './lib/component-enum-policy.mjs';

import { packageSourceRoots, packageRelativePath } from './lib/source-layout.mjs';
const inspectedCounter = { current: 0 };
const failures = [];
const walk = async (folder) => {
  for (const entry of await readdir(folder, { withFileTypes: true })) {
    const absolute = path.join(folder, entry.name);
    if (entry.isDirectory()) {
      if (!['__tests__', 'node_modules', '.generated'].includes(entry.name)) await walk(absolute);
      continue;
    }
    if (!entry.name.endsWith('.tsx') || entry.name.endsWith('.test.tsx')) continue;
    const relative = packageRelativePath(absolute).replaceAll(path.sep, '/');
    const source = await readFile(absolute, 'utf8');
    for (const issue of validateComponentEnums(relative, source)) {
      failures.push(
        relative +
          ':' +
          issue.line +
          ': ' +
          issue.component +
          '.' +
          issue.prop +
          '="' +
          issue.value +
          '" is a closed-set literal; import its owner enum',
      );
    }
    inspectedCounter.current++;
  }
};
for (const directory of packageSourceRoots) await walk(directory);
if (failures.length) {
  for (const problem of failures) process.stderr.write(problem + '\n');
  process.exitCode = 1;
} else {
  process.stdout.write(
    'Closed UI prop enums verified across ' +
      inspectedCounter.current +
      ' production TSX modules.\n',
  );
}
