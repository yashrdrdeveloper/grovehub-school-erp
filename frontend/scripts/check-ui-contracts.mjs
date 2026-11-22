import { readdir, readFile } from 'node:fs/promises';
import path from 'node:path';

const root = path.join(process.cwd(), 'packages/components/components');
const failures = [];
let checked = 0;

for (const unit of await readdir(root, { withFileTypes: true })) {
  if (unit.isFile() && unit.name === 'SchoolMonogram.tsx') {
    const source = await readFile(path.join(root, unit.name), 'utf8');
    checked++;
    if (!/dataTestId/.test(source) || !/data-testid/.test(source)) {
      failures.push('SchoolMonogram: dataTestId must reach its native root');
    }
    continue;
  }
  if (!unit.isDirectory()) continue;

  const folder = path.join(root, unit.name);
  const impl = (await readdir(folder)).find(file => file === `${unit.name}.tsx`);
  if (!impl) continue;

  checked++;
  const source = await readFile(path.join(folder, impl), 'utf8');
  const forwardsToInput = unit.name === 'TextField'
    && /<Input\b/.test(source) && /\{\.\.\.props\}/.test(source);

  if (!forwardsToInput && !/dataTestId/.test(source)) {
    failures.push(`${unit.name}: missing explicit dataTestId propagation`);
  }

  if (!forwardsToInput && !/data-testid/.test(source)
      && !/['"]data-testid['"]/.test(source)) {
    failures.push(`${unit.name}: dataTestId must reach a native DOM element`);
  }

  const typeDir = path.join(folder, 'types');
  try {
    const contracts = (await Promise.all((await readdir(typeDir))
      .filter(file => file.endsWith('.ts'))
      .map(file => readFile(path.join(typeDir, file), 'utf8')))).join('\n');
    if (!/dataTestId/.test(contracts) && !/InputProps/.test(contracts)) {
      failures.push(`${unit.name}: public component prop contract must expose dataTestId`);
    }
  } catch {
    failures.push(`${unit.name}: missing typed prop contract`);
  }
}

if (failures.length) {
  for (const failure of failures) process.stderr.write(`${failure}\n`);
  process.exitCode = 1;
} else {
  process.stdout.write(`UI automation contracts verified across ${checked} shared components.\n`);
}
