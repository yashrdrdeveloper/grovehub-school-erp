import { readFile, readdir } from 'node:fs/promises';
import path from 'node:path';
import ts from 'typescript';

/** Block native markup in feature presentation; primitives own actual HTML elements. */
const nativeTags = new Set([
  'div',
  'span',
  'p',
  'h1',
  'h2',
  'h3',
  'h4',
  'h5',
  'h6',
  'small',
  'strong',
  'button',
  'input',
  'label',
  'form',
  'a',
  'select',
  'textarea',
  'main',
  'section',
  'header',
  'footer',
  'aside',
  'nav',
  'article',
]);
const errors = [];
const scannedCounter = { current: 0 };
const walk = async (directory) => {
  for (const item of await readdir(directory, { withFileTypes: true })) {
    const filename = path.join(directory, item.name);
    if (item.isDirectory()) {
      if (item.name === '__tests__') continue;
      await walk(filename);
      continue;
    }
    if (!/\.tsx$/.test(item.name) || /\.test\.tsx$/.test(item.name)) continue;
    const text = await readFile(filename, 'utf8');
    const tree = ts.createSourceFile(
      filename,
      text,
      ts.ScriptTarget.Latest,
      true,
      ts.ScriptKind.TSX,
    );
    scannedCounter.current++;
    const visit = (node) => {
      if (ts.isJsxOpeningElement(node) || ts.isJsxSelfClosingElement(node)) {
        const tag = node.tagName.getText(tree);
        if (nativeTags.has(tag)) {
          const at = tree.getLineAndCharacterOfPosition(node.getStart(tree));
          errors.push(
            path.relative(process.cwd(), filename) +
              ':' +
              (at.line + 1) +
              ': raw native <' +
              tag +
              '>',
          );
        }
      }
      ts.forEachChild(node, visit);
    };
    visit(tree);
  }
};
await walk(path.join(process.cwd(), 'packages/domain'));
await walk(path.join(process.cwd(), 'packages/shell'));
for (const error of errors) process.stderr.write(error + '\n');
if (errors.length > 0) process.exitCode = 1;
else process.stdout.write('Semantic UI boundary verified across ' + scannedCounter.current + ' TSX files.\\n');
