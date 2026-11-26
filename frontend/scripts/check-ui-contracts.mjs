import { readdir, readFile } from 'node:fs/promises';
import path from 'node:path';

import { resolvePackageSource } from './lib/source-layout.mjs';

const root = resolvePackageSource('components/components');
const failures = [];
const checkedCounter = { current: 0 };

for (const unit of await readdir(root, { withFileTypes: true })) {
  if (!unit.isDirectory()) continue;

  const folder = path.join(root, unit.name);
  const impl = (await readdir(folder)).find((file) => file === `${unit.name}.tsx`);
  if (!impl) continue;

  checkedCounter.current++;
  const source = await readFile(path.join(folder, impl), 'utf8');
  const forwardsToTypography =
    unit.name === 'Text' &&
    /<Typography\b/.test(source) &&
    /dataTestId=\{dataTestId\}/.test(source);
  const forwardsToInput =
    unit.name === 'TextField' && /<Input\b/.test(source) && /\{\.\.\.props\}/.test(source);
  // Explicit forwarding through shared controls; each destination is separately checked.
  const forwardsToSharedControl =
    /<(?:Input|Button|ButtonWithIcon|Dialog|SelectionPanel)\b/.test(source) &&
    /dataTestId=\{dataTestId\}/.test(source);
  const forwardsToBox = /<Box\b/.test(source) && /dataTestId=\{/.test(source);
  const forwardsToNativeLeaf =
    unit.name === 'Input' &&
    /<InputControl\b/.test(source) &&
    /dataTestId=\{controlTestId\}/.test(source);

  if (!forwardsToInput && !/dataTestId/.test(source)) {
    failures.push(`${unit.name}: missing explicit dataTestId propagation`);
  }

  if (
    !forwardsToInput &&
    !forwardsToNativeLeaf &&
    !forwardsToTypography &&
    !forwardsToBox &&
    !forwardsToSharedControl &&
    !/data-testid/.test(source) &&
    !/['"]data-testid['"]/.test(source)
  ) {
    failures.push(`${unit.name}: dataTestId must reach a native DOM element`);
  }

  const typeDir = path.join(folder, 'types');
  try {
    const contracts = (
      await Promise.all(
        (await readdir(typeDir))
          .filter((file) => file.endsWith('.ts'))
          .map((file) => readFile(path.join(typeDir, file), 'utf8')),
      )
    ).join('\n');
    if (!/dataTestId/.test(contracts) && !/InputProps/.test(contracts)) {
      failures.push(`${unit.name}: public component prop contract must expose dataTestId`);
    }
  } catch {
    failures.push(`${unit.name}: missing typed prop contract`);
  }
}

const brandSource = await readFile(
  path.join(process.cwd(), 'packages/domain/brand/components/schoolMonogram/SchoolMonogram.tsx'),
  'utf8',
);
if (!/<Box\b/.test(brandSource) || !/dataTestId=\{dataTestId\}/.test(brandSource)) {
  failures.push('SchoolMonogram: dataTestId must reach its native root through Box');
}

if (failures.length) {
  for (const failure of failures) process.stderr.write(`${failure}\n`);
  process.exitCode = 1;
} else {
  process.stdout.write(
    `UI automation contracts verified across ${checkedCounter.current} shared components.\n`,
  );
}
