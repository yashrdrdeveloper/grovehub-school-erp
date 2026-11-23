import { readdir, readFile } from 'node:fs/promises';
import path from 'node:path';
import ts from 'typescript';

/** Blocking AST rule for production application code; type annotations are not ternaries. */
const violations = [];
const checkedCounter = { current: 0 };
const visitDir = async (directory) => {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const filename = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      if (
        entry.name === '__tests__' ||
        entry.name === 'node_modules' ||
        entry.name === '.generated'
      )
        continue;
      await visitDir(filename);
      continue;
    }
    if (!/\.tsx?$/.test(entry.name) || /\.test\.tsx?$/.test(entry.name)) continue;
    const source = await readFile(filename, 'utf8');
    const tree = ts.createSourceFile(
      filename,
      source,
      ts.ScriptTarget.Latest,
      true,
      filename.endsWith('.tsx') ? ts.ScriptKind.TSX : ts.ScriptKind.TS,
    );
    checkedCounter.current++;
    const check = (node) => {
      if (ts.isConditionalExpression(node)) {
        const location = tree.getLineAndCharacterOfPosition(node.getStart(tree));
        violations.push(path.relative(process.cwd(), filename) + ':' + (location.line + 1));
      }
      ts.forEachChild(node, check);
    };
    check(tree);
  }
};
await visitDir(path.join(process.cwd(), 'app'));
await visitDir(path.join(process.cwd(), 'packages'));
if (violations.length > 0) {
  for (const item of violations)
    process.stderr.write(item + ': conditional expression is forbidden\\n');
  process.exitCode = 1;
} else {
  process.stdout.write(
    'No production ternary expressions in ' + checkedCounter.current + ' source files.\\n',
  );
}
