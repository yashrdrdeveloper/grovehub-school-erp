import { readFile, readdir, writeFile } from 'node:fs/promises';
import path from 'node:path';

import { analyzeI18nSource } from './lib/i18n-policy.mjs';

const root = process.cwd();
const baselineFile = path.join(root, 'scripts/i18n-allowlist.json');
const writeBaseline = process.argv.includes('--write-baseline');
const issues = [];
const copies = {};
const sourceCountCounter = { current: 0 };
const copyCountCounter = { current: 0 };

const bootstrapInstaller =
  'packages/shell/containers/shellContainer/helpers/shellContainer.installI18n.ts';
const bootFallback =
  'packages/shell/containers/shellContainer/components/i18nBootstrapFallback/I18nBootstrapFallback.tsx';

const walk = async (directory) => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const abs = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      if (!['__tests__', 'node_modules', '.generated'].includes(entry.name)) {
        await walk(abs);
      }
      continue;
    }
    if (!/\.(?:tsx?|jsx?)$/.test(entry.name) || /\.test\.[jt]sx?$/.test(entry.name)) continue;
    const relative = path.relative(root, abs).replaceAll(path.sep, '/');
    const source = await readFile(abs, 'utf8');
    const options = {
      allowRuntimeGlobal: relative === bootstrapInstaller,
      allowResolverImport: relative === bootstrapInstaller || relative.startsWith('packages/base/'),
      allowEagerCopy: relative.includes('/constants/'),
      checkVisibleCopy: relative !== bootFallback && relative.endsWith('.tsx'),
    };
    const result = analyzeI18nSource(relative, source, options);
    for (const issue of result.issues) {
      issues.push(relative + ':' + issue.line + ': ' + issue.code + ' ' + issue.label);
    }
    if (result.rawCopy.length) {
      copies[relative] = result.rawCopy.length;
      copyCountCounter.current += result.rawCopy.length;
    }
    sourceCountCounter.current++;
  }
};
for (const folder of [
  'packages/base',
  'packages/business',
  'packages/components',
  'packages/domain',
  'packages/shell',
]) {
  await walk(path.join(root, folder));
}

if (issues.length) {
  for (const issue of issues) process.stderr.write(issue + '\n');
  process.exitCode = 1;
} else if (writeBaseline) {
  const baseline = {
    description:
      'Initial measured JSX copy debt after the Phase 2.1 source-message codemod. Increase is forbidden.',
    maxTotal: copyCountCounter.current,
    counts: Object.fromEntries(Object.entries(copies).sort(([a], [b]) => a.localeCompare(b))),
  };
  await writeFile(baselineFile, JSON.stringify(baseline, null, 2) + '\n');
  process.stdout.write(
    'Recorded ' +
      copyCountCounter.current +
      ' legacy JSX copy findings across ' +
      Object.keys(copies).length +
      ' files. ' +
      sourceCountCounter.current +
      ' source modules checked.\n',
  );
} else {
  const readBaseline = async () => {
    try {
      return JSON.parse(await readFile(baselineFile, 'utf8'));
    } catch {
      process.stderr.write(
        'Missing scripts/i18n-allowlist.json; measure the current migration first.\n',
      );
      process.exitCode = 1;
      return null;
    }
  };
  const baseline = await readBaseline();
  if (baseline) {
    if (copyCountCounter.current > baseline.maxTotal) {
      issues.push(
        'RAW_UI_COPY total increased from ' + baseline.maxTotal + ' to ' + copyCountCounter.current,
      );
    }
    for (const [file, current] of Object.entries(copies)) {
      const allowed = baseline.counts[file] ?? 0;
      if (current > allowed) {
        issues.push('RAW_UI_COPY ' + file + ' increased from ' + allowed + ' to ' + current);
      }
    }
    if (issues.length) {
      for (const issue of issues) process.stderr.write(issue + '\n');
      process.exitCode = 1;
    } else {
      process.stdout.write(
        'Global __() ownership verified across ' +
          sourceCountCounter.current +
          ' source modules; ' +
          copyCountCounter.current +
          ' legacy JSX messages remain within the locked per-file baseline.\n',
      );
    }
  }
}
