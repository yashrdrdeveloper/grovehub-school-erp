import { access, readFile, readdir } from 'node:fs/promises';
import path from 'node:path';
import ts from 'typescript';

/**
 * Block misplaced business transport contracts, flat compound primitives, and
 * duplicated form normalization without restricting genuine domain route logic.
 */
const root = path.join(process.cwd(), 'packages');
const failures = [];
const exists = async (filename) => {
  try {
    await access(filename);
    return true;
  } catch {
    return false;
  }
};

const business = path.join(root, 'business');
for (const entry of await readdir(business, { withFileTypes: true })) {
  if (!entry.isDirectory()) {
    continue;
  }

  const owner = path.join(business, entry.name);
  if (await exists(path.join(owner, 'types.ts'))) {
    failures.push(`business/${entry.name}/types.ts: use the owned types/ directory`);
  }

  for (const name of ['constants', 'helpers']) {
    const folder = path.join(owner, name);
    if (!await exists(folder)) {
      continue;
    }

    for (const file of await readdir(folder, { withFileTypes: true })) {
      if (!file.isFile()) {
        continue;
      }

      if (name === 'constants' && /\.(?:endpoints|operations|proofScheme)\.ts$/.test(file.name)) {
        failures.push(`business/${entry.name}/constants/${file.name}: service metadata belongs in services/constants`);
      }

      if (name === 'helpers' && /\.(?:getEndpoint|getProofHeader)\.ts$/.test(file.name)) {
        failures.push(`business/${entry.name}/helpers/${file.name}: service helper belongs in services/helpers`);
      }
    }
  }
}

/**
 * Enforce reusable owner conventions across every package, not only Dialog.
 * A flat parts/components/primitives barrel hides a compound unit's ownership.
 */
const checkOwnerFiles = async (directory) => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const filename = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      if (entry.name !== 'node_modules' && entry.name !== '.generated') {
        await checkOwnerFiles(filename);
      }
      continue;
    }

    const relative = path.relative(root, filename).replaceAll(path.sep, '/');
    if (/\/(?:helpers|constants|renderer)\/index\.ts$/.test(relative)) {
      failures.push(relative + ': helper/constant/renderer barrels are forbidden');
    }

    if (/\/config\//.test(relative)) {
      failures.push(relative + ': move config to owned constants/ or helpers/');
    }

    if (/\/helpers\/(?!types\/|__tests__\/)[^/]+\/.+\.[jt]sx?$/.test(relative) &&
        !relative.startsWith('base/')) {
      failures.push(relative + ': helpers must be flat, with types/ and __tests__/ only');
    }

    if (/^components\/components\/[^/]+\/(?:parts|components|primitives)\.[jt]sx?$/.test(relative)) {
      failures.push(relative + ': compound child components require named owned units');
    }
  }
};

await checkOwnerFiles(root);

const dialogRoot = path.join(root, 'components/components/Dialog');
if (await exists(path.join(dialogRoot, 'parts.ts'))) {
  failures.push('components/Dialog/parts.ts: use one named unit per compound component');
}
for (const [unit, component] of [
  ['root', 'DialogRoot'],
  ['portal', 'DialogPortal'],
  ['overlay', 'DialogOverlay'],
  ['content', 'DialogContent'],
  ['title', 'DialogTitle'],
  ['description', 'DialogDescription'],
  ['close', 'DialogClose'],
]) {
  const folder = path.join(dialogRoot, 'components', unit);
  for (const relative of [`${component}.tsx`, 'index.ts', `types/${unit}.props.ts`]) {
    if (!await exists(path.join(folder, relative))) {
      failures.push(`Dialog/components/${unit}: missing ${relative}`);
    }
  }
}

const dialogConsumers = [
  path.join(root, 'domain'),
  path.join(root, 'shell'),
  path.join(root, 'components'),
  path.join(root, 'business'),
];
const checkImports = async (directory) => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const filename = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      if (entry.name !== '__tests__') {
        await checkImports(filename);
      }
      continue;
    }

    if (!/\.tsx?$/.test(entry.name)) {
      continue;
    }

    const source = await readFile(filename, 'utf8');
    const ast = ts.createSourceFile(filename, source, ts.ScriptTarget.Latest, true);
    for (const statement of ast.statements) {
      if (ts.isImportDeclaration(statement) &&
          ts.isStringLiteral(statement.moduleSpecifier) &&
          /Dialog\/parts$/.test(statement.moduleSpecifier.text)) {
        failures.push(path.relative(root, filename) + ': import Dialog units through Dialog public API');
      }
    }
  }
};
for (const directory of dialogConsumers) {
  await checkImports(directory);
}

for (const unit of ['Input', 'PasswordField']) {
  const filename = path.join(root, 'components/components', unit, `${unit}.tsx`);
  const source = await readFile(filename, 'utf8');

  if (/\.replace\s*\(|\.filter\s*\(\s*Boolean\s*\)/.test(source)) {
    failures.push(`${unit}: use the tested FormField normalization/description helpers`);
  }

  for (const helper of [
    'getFormFieldControlTestId',
    'getFormFieldDescriptionId',
    'getFormFieldCombinedDescription',
  ]) {
    if (!source.includes(helper)) {
      failures.push(`${unit}: missing shared helper ${helper}`);
    }
  }
}

if (failures.length > 0) {
  failures.forEach((item) => process.stderr.write(item + '\n'));
  process.exitCode = 1;
} else {
  process.stdout.write('Business service ownership, Dialog units and shared field helpers verified.\n');
}
