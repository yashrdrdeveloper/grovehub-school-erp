import { readdir, readFile } from 'node:fs/promises';
import path from 'node:path';
import ts from 'typescript';

/** Route UI and hooks use @Components/router, never direct React Router imports. */
const violations = [];
const root = path.join(process.cwd(), 'packages');

const visit = async (folder) => {
  for (const entry of await readdir(folder, { withFileTypes: true })) {
    const file = path.join(folder, entry.name);

    if (entry.isDirectory()) {
      if (entry.name !== '__tests__') {
        await visit(file);
      }

      continue;
    }

    if (!/\.tsx?$/.test(entry.name) || /\.test\.tsx?$/.test(entry.name)) {
      continue;
    }

    const source = await readFile(file, 'utf8');
    const ast = ts.createSourceFile(file, source, ts.ScriptTarget.Latest, true);
    for (const statement of ast.statements) {
      if (!ts.isImportDeclaration(statement)) {
        continue;
      }

      if (!ts.isStringLiteral(statement.moduleSpecifier)) {
        continue;
      }

      if (statement.moduleSpecifier.text === 'react-router-dom') {
        violations.push(path.relative(process.cwd(), file));
      }
    }
  }
};

await visit(path.join(root, 'domain'));
await visit(path.join(root, 'shell'));

if (violations.length > 0) {
  for (const source of violations) {
    process.stderr.write(source + ': import React Router through @Components/router\\n');
  }

  process.exitCode = 1;
} else {
  process.stdout.write('All domain and shell Router APIs pass through @Components/router.\\n');
}
