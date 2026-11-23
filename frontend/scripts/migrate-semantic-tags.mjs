import { readFile, readdir, writeFile } from 'node:fs/promises';
import path from 'node:path';
import ts from 'typescript';

/**
 * One-time owner-safe migration and permanent read-only gate.
 * Converts semantic JSX tag names to existing typed component constants.
 * Use --check in CI; no source is mutated by the check.
 */
const root = process.cwd();
const checkOnly = process.argv.includes('--check');
const variants = {
  Box: {
    symbol: 'BOX_TAG',
    module: '@Components/components/Box',
    allowed: new Set([
      'div',
      'span',
      'p',
      'small',
      'strong',
      'figure',
      'figcaption',
      'em',
      'i',
      'time',
      'main',
      'section',
      'article',
      'aside',
      'header',
      'footer',
      'nav',
      'h1',
      'h2',
      'h3',
      'h4',
      'h5',
      'h6',
    ]),
  },
  Text: {
    symbol: 'TYPOGRAPHY_TAG',
    module: '@Components/components/Typography',
    allowed: new Set(['span', 'p', 'small', 'strong']),
  },
  Typography: {
    symbol: 'TYPOGRAPHY_TAG',
    module: '@Components/components/Typography',
    allowed: new Set(['h1', 'h2', 'h3', 'h4', 'p', 'span', 'small', 'strong']),
  },
};
const filesUpdatedCounter = { current: 0 };
const changesCounter = { current: 0 };
const violations = [];

const scan = async (folder) => {
  for (const entry of await readdir(folder, { withFileTypes: true })) {
    const filename = path.join(folder, entry.name);
    if (entry.isDirectory()) {
      if (!['node_modules', '.generated', 'dist'].includes(entry.name)) await scan(filename);
      continue;
    }
    if (!entry.name.endsWith('.tsx')) continue;

    const before = await readFile(filename, 'utf8');
    const sf = ts.createSourceFile(
      filename,
      before,
      ts.ScriptTarget.Latest,
      true,
      ts.ScriptKind.TSX,
    );
    const edits = [];
    const imports = new Map();
    const visit = (node) => {
      if (
        (ts.isJsxOpeningElement(node) || ts.isJsxSelfClosingElement(node)) &&
        ts.isIdentifier(node.tagName)
      ) {
        const owner = variants[node.tagName.text];
        if (owner) {
          for (const attribute of node.attributes.properties) {
            if (
              !ts.isJsxAttribute(attribute) ||
              attribute.name.text !== 'as' ||
              !attribute.initializer ||
              !ts.isStringLiteral(attribute.initializer)
            )
              continue;
            const tag = attribute.initializer.text;
            const location = sf.getLineAndCharacterOfPosition(attribute.getStart(sf));
            if (!owner.allowed.has(tag)) {
              violations.push(
                path.relative(root, filename) +
                  ':' +
                  (location.line + 1) +
                  ': unknown semantic tag ' +
                  tag,
              );
              continue;
            }
            if (checkOnly) {
              violations.push(
                path.relative(root, filename) +
                  ':' +
                  (location.line + 1) +
                  ': ' +
                  node.tagName.text +
                  ' as="' +
                  tag +
                  '" requires a typed semantic constant',
              );
              continue;
            }
            edits.push({
              start: attribute.getStart(sf),
              end: attribute.getEnd(),
              text: 'as={' + owner.symbol + '.' + tag.toUpperCase() + '}',
            });
            imports.set(owner.symbol, owner.module);
            changesCounter.current++;
          }
        }
      }
      ts.forEachChild(node, visit);
    };
    visit(sf);
    if (!edits.length) continue;
    // Reuse the original import when it owns the same public component API.
    const importsToAdd = [];
    for (const [symbol, module] of imports) {
      if (new RegExp('\\b' + symbol + '\\b').test(before)) continue;
      const existing = sf.statements.find(
        (statement) =>
          ts.isImportDeclaration(statement) &&
          ts.isStringLiteral(statement.moduleSpecifier) &&
          statement.moduleSpecifier.text === module &&
          statement.importClause?.namedBindings &&
          ts.isNamedImports(statement.importClause.namedBindings),
      );
      if (!existing) {
        importsToAdd.push('import { ' + symbol + " } from '" + module + "';");
        continue;
      }
      const names = existing.importClause.namedBindings;
      edits.push({ start: names.getEnd() - 1, end: names.getEnd() - 1, text: ', ' + symbol });
    }

    // Replacements are nonoverlapping source ranges. Apply backwards to preserve positions.
    const updated = edits
      .sort((a, b) => b.start - a.start)
      .reduce(
        (current, edit) => current.slice(0, edit.start) + edit.text + current.slice(edit.end),
        before,
      );
    const withImports = importsToAdd.length ? importsToAdd.join('\n') + '\n' + updated : updated;
    await writeFile(filename, withImports, 'utf8');
    filesUpdatedCounter.current++;
  }
};

await scan(path.join(root, 'app'));
await scan(path.join(root, 'packages'));
if (violations.length) {
  for (const violation of violations) process.stderr.write(violation + '\n');
  process.exitCode = 1;
}
if (checkOnly && !violations.length) process.stdout.write('Semantic JSX tag constants verified.\n');
if (!checkOnly)
  process.stdout.write(
    'Migrated ' +
      changesCounter.current +
      ' JSX attributes across ' +
      filesUpdatedCounter.current +
      ' files.\n',
  );
