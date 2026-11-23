/**
 * One-time measured migration: wrap plain display copy in bare global __().
 *
 * Only JSX literal text/known message attributes and literal JSX expressions are
 * changed. Domain copy constants, API protocol strings, and dynamic text are
 * deliberately left intact; those require owner-aware refactoring.
 */
import { readFile, readdir, writeFile } from 'node:fs/promises';
import path from 'node:path';
import ts from 'typescript';

import { normalizeJsxText, VISIBLE_TEXT_ATTRIBUTES } from './lib/i18n-policy.mjs';

const root = process.cwd();
const directories = ['packages/domain', 'packages/shell', 'packages/components'];
const skipped = new Set([
  'packages/shell/containers/shellContainer/components/i18nBootstrapFallback/I18nBootstrapFallback.tsx',
]);
const hasLetters = (message) => /[A-Za-z\u00c0-\u024f]/.test(message);
const hasEntity = (message) => /&(?:[a-zA-Z][a-zA-Z0-9]+|#[xX]?[0-9a-fA-F]+);/.test(message);
const allFiles = [];
const scan = async (directory) => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const absolute = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      if (entry.name !== '__tests__' && entry.name !== 'node_modules') await scan(absolute);
      continue;
    }
    if (!entry.name.endsWith('.tsx') || entry.name.endsWith('.test.tsx')) continue;
    const relative = path.relative(root, absolute).replaceAll(path.sep, '/');
    if (!skipped.has(relative)) allFiles.push({ absolute, relative });
  }
};
for (const folder of directories) await scan(path.join(root, folder));

let changedFiles = 0;
let replacedNodes = 0;
let skippedEntities = 0;
for (const { absolute, relative } of allFiles) {
  const source = await readFile(absolute, 'utf8');
  const syntax = ts.createSourceFile(relative, source, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX);
  const edits = [];
  const insert = (node, text) => {
    edits.push({ start: node.pos, end: node.end, text });
  };
  const visit = (node) => {
    if (ts.isJsxText(node)) {
      const text = normalizeJsxText(node.getText(syntax));
      if (hasEntity(text)) skippedEntities++;
      if (hasLetters(text) && !hasEntity(text)) {
        insert(node, '{__(' + JSON.stringify(text) + ')}');
      }
    } else if (ts.isJsxAttribute(node) &&
        VISIBLE_TEXT_ATTRIBUTES.has(node.name.text) &&
        node.initializer && ts.isStringLiteral(node.initializer)) {
      const text = node.initializer.text;
      if (hasLetters(text)) {
        edits.push({
          start: node.initializer.getStart(syntax), end: node.initializer.end,
          text: '{__(' + JSON.stringify(text) + ')}',
        });
      }
    } else if (ts.isJsxExpression(node) && node.expression &&
        ts.isStringLiteral(node.expression) && hasLetters(node.expression.text)) {
      // Only convert literal string expressions; changing dynamic expressions can
      // accidentally alter API identifiers and is out of scope.
      edits.push({
        start: node.expression.getStart(syntax), end: node.expression.end,
        text: '__(' + JSON.stringify(node.expression.text) + ')',
      });
    }
    ts.forEachChild(node, visit);
  };
  ts.forEachChild(syntax, visit);
  if (!edits.length) continue;
  let output = source;
  for (const edit of edits.sort((a, b) => b.start - a.start)) {
    output = output.slice(0, edit.start) + edit.text + output.slice(edit.end);
  }
  await writeFile(absolute, output, 'utf8');
  changedFiles++;
  replacedNodes += edits.length;
}
process.stdout.write(
  'Wrapped ' + replacedNodes + ' user-facing JSX expressions in ' + changedFiles +
    ' files; left ' + skippedEntities + ' entity-encoded fragments for manual review.\n',
);
